package training

import (
	"context"
	"fmt"
	"gin-vue-admin/global"
	apisixReq "gin-vue-admin/model/apisix/request"
	clusterModel "gin-vue-admin/model/cluster"
	"gin-vue-admin/model/consts"
	imageModel "gin-vue-admin/model/image"
	orderModel "gin-vue-admin/model/order"
	orderReq "gin-vue-admin/model/order/request"
	productModel "gin-vue-admin/model/product"
	pvcModel "gin-vue-admin/model/pvc"
	trainingModel "gin-vue-admin/model/training"
	trainingReq "gin-vue-admin/model/training/request"
	trainingResp "gin-vue-admin/model/training/response"
	orderService "gin-vue-admin/service/order"
	productService "gin-vue-admin/service/product"
	terminalService "gin-vue-admin/service/terminal"
	"gin-vue-admin/service/training/builder"
	helper "gin-vue-admin/utils/k8s"
	"gin-vue-admin/utils/timer"
	"gin-vue-admin/utils/validator"
	"io"
	"time"

	"github.com/pkg/errors"
	"github.com/zeromicro/go-zero/core/logx"
	"gorm.io/gorm"
	corev1 "k8s.io/api/core/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/client-go/kubernetes"
)

// TrainingJobManager 训练任务管理器接口
type TrainingJobManager interface {
	CreateTrainingJob(ctx context.Context, req *trainingReq.CreateTrainingJobReq) (*trainingResp.CreateTrainingJobResp, error)
	DeleteTrainingJob(ctx context.Context, req *trainingReq.DeleteTrainingJobReq) error
	StopTrainingJob(ctx context.Context, req *trainingReq.StopTrainingJobReq) error
	GetTrainingJobList(ctx context.Context, req *trainingReq.GetTrainingJobListReq) (*trainingResp.GetTrainingJobListResp, error)
	GetTrainingJobDetail(ctx context.Context, req *trainingReq.GetTrainingJobDetailReq) (*trainingResp.TrainingJobDetail, error)
	GetTrainingJobLogs(ctx context.Context, req *trainingReq.GetTrainingJobLogsReq) (io.ReadCloser, error)
	GetTrainingJobPods(ctx context.Context, id uint) ([]terminalService.PodInfo, error)
}

var _ TrainingJobManager = (*TrainingJobService)(nil)

type trainingApisixService interface {
	CreateRoute(ctx context.Context, req *apisixReq.CreateRouteReq) error
	DeleteRoute(ctx context.Context, req *apisixReq.DeleteRouteReq) error
}

// TrainingJobService 训练任务服务实现
type TrainingJobService struct {
	apisixSvc trainingApisixService
}

// Cleanups 资源清理函数列表
type Cleanups []func(ctx context.Context)

type trainingCreateState struct {
	req         *trainingReq.CreateTrainingJobReq
	jobName     string
	displayName string
	image       *imageModel.Image
	product     *productModel.Product
	plan        *productService.AllocationPlan
	cluster     *global.ClusterClientInfo
	trainingJob *trainingModel.TrainingJob
}

type trainingListLookup struct {
	images   map[uint]imageModel.Image
	clusters map[uint]clusterModel.K8sCluster
}

var TrainingJobServiceApp = new(TrainingJobService)

// SetApisixService 设置 Apisix 服务依赖
func (s *TrainingJobService) SetApisixService(svc trainingApisixService) {
	s.apisixSvc = svc
}

func (c *Cleanups) Add(fn func(ctx context.Context)) {
	*c = append(*c, fn)
}

func (c Cleanups) Run(ctx context.Context) {
	for i := len(c) - 1; i >= 0; i-- {
		c[i](ctx)
	}
}

// validateCreateRequest 校验创建请求参数
func (s *TrainingJobService) validateCreateRequest(req *trainingReq.CreateTrainingJobReq) error {
	// 1. 校验挂载路径
	mountsPaths := make(map[string]bool)
	for _, mount := range req.Mounts {
		if mount.PvcId == 0 {
			return errors.New("挂载配置必须指定 PVC ID")
		}

		if err := validator.ValidateMountPath(mount.MountPath); err != nil {
			logx.Error("挂载路径验证失败", logx.Field("path", mount.MountPath), logx.Field("err", err))
			return err
		}

		if mountsPaths[mount.MountPath] {
			logx.Error("挂载路径重复", logx.Field("path", mount.MountPath))
			return errors.New("挂载路径重复: " + mount.MountPath)
		}
		mountsPaths[mount.MountPath] = true
	}

	// 2. 验证 TensorBoard 日志路径安全性
	if req.EnableTensorboard && req.TensorboardLogPath != "" {
		if err := validator.ValidateSubPath(req.TensorboardLogPath); err != nil {
			logx.Error("TensorBoard路径验证失败", logx.Field("path", req.TensorboardLogPath), logx.Field("err", err))
			return err
		}
	}

	return nil
}

// CreateTrainingJob 创建训练任务
func (s *TrainingJobService) CreateTrainingJob(ctx context.Context, req *trainingReq.CreateTrainingJobReq) (resp *trainingResp.CreateTrainingJobResp, err error) {
	if err = s.validateCreateRequest(req); err != nil {
		return nil, err
	}

	cleanups := make(Cleanups, 0)
	defer func() {
		r := recover()
		if r != nil || err != nil {
			if r != nil {
				logx.Error("Panic recovered", r)
				err = fmt.Errorf("panic: %v", r)
			}
			cleanups.Run(context.Background())
		}
	}()

	state, err := s.buildCreateTrainingState(ctx, req, &cleanups)
	if err != nil {
		return nil, err
	}

	if err = s.persistTrainingJob(ctx, state); err != nil {
		return nil, err
	}

	if err = s.saveTrainingAllocations(ctx, state, &cleanups); err != nil {
		return nil, err
	}

	if err = s.saveTrainingJobMounts(ctx, state); err != nil {
		return nil, err
	}

	if err = s.saveTrainingJobEnvs(ctx, state.trainingJob.ID, state.req.Envs); err != nil {
		return nil, err
	}

	if err = s.submitTrainingJob(ctx, state, &cleanups); err != nil {
		return nil, err
	}

	if err = s.createTrainingOrder(ctx, state, &cleanups); err != nil {
		return nil, err
	}

	go s.createOptionalResources(context.Background(), state.trainingJob, state.req, state.cluster)

	logx.Info("创建训练任务成功",
		logx.Field("job", state.jobName),
		logx.Field("displayName", state.displayName),
		logx.Field("namespace", state.req.Namespace),
		logx.Field("framework", state.req.FrameworkType),
		logx.Field("tensorboard", state.req.EnableTensorboard),
	)

	return &trainingResp.CreateTrainingJobResp{
		ID:          state.trainingJob.ID,
		DisplayName: state.displayName,
	}, nil
}

func (s *TrainingJobService) buildCreateTrainingState(
	ctx context.Context,
	req *trainingReq.CreateTrainingJobReq,
	cleanups *Cleanups,
) (*trainingCreateState, error) {
	if err := normalizeTrainingCreateRequest(req); err != nil {
		return nil, err
	}

	jobName, displayName := s.prepareTrainingCreateRequest(req)

	image, err := s.getTrainingImage(ctx, req.ImageId)
	if err != nil {
		return nil, err
	}

	product, err := s.getTrainingProduct(ctx, req.TemplateProductId)
	if err != nil {
		return nil, err
	}

	if err := s.checkTrainingBalance(ctx, req); err != nil {
		return nil, err
	}

	computeNodes, err := normalizeTrainingComputeNodes(req)
	if err != nil {
		return nil, err
	}

	plan, err := (&productService.ProductService{}).PlanAllocations(ctx, product.ID, computeNodes, req.ScheduleStrategy)
	if err != nil {
		return nil, err
	}

	if err := s.reserveTrainingCapacity(ctx, plan, cleanups); err != nil {
		return nil, err
	}

	cluster, err := s.getTrainingCluster(product.ClusterId)
	if err != nil {
		return nil, err
	}

	trainingJob := s.buildTrainingJobRecord(req, jobName, displayName, product, computeNodes)
	return &trainingCreateState{
		req:         req,
		jobName:     jobName,
		displayName: displayName,
		image:       image,
		product:     product,
		plan:        plan,
		cluster:     cluster,
		trainingJob: trainingJob,
	}, nil
}

func (s *TrainingJobService) prepareTrainingCreateRequest(req *trainingReq.CreateTrainingJobReq) (string, string) {
	jobName := helper.GenerateInstanceName(consts.TrainingInstance)
	displayName := req.Name
	if displayName == "" {
		displayName = jobName
	}
	if req.Namespace == "" {
		req.Namespace = consts.DefaultNamespace
	}
	return jobName, displayName
}

func (s *TrainingJobService) getTrainingImage(ctx context.Context, imageID uint) (*imageModel.Image, error) {
	var image imageModel.Image
	if err := global.GVA_DB.WithContext(ctx).Where("id = ?", imageID).First(&image).Error; err != nil {
		return nil, errors.Wrap(err, "获取镜像信息失败")
	}
	return &image, nil
}

func (s *TrainingJobService) getTrainingProduct(ctx context.Context, productID uint) (*productModel.Product, error) {
	if productID == 0 {
		return nil, errors.New("必须指定模板商品ID (templateProductId)")
	}

	var product productModel.Product
	if err := global.GVA_DB.WithContext(ctx).Where("id = ?", productID).First(&product).Error; err != nil {
		return nil, errors.Wrap(err, "获取产品规格失败")
	}
	if product.Status != productModel.ProductStatusEnabled {
		return nil, errors.New("所选产品已下架，无法使用")
	}

	return &product, nil
}

func (s *TrainingJobService) checkTrainingBalance(ctx context.Context, req *trainingReq.CreateTrainingJobReq) error {
	orderSvc := &orderService.OrderService{}
	return orderSvc.CheckBalanceSufficient(ctx, req.UserId, req.TemplateProductId, req.PayType, req.InstanceCount)
}

func normalizeTrainingComputeNodes(req *trainingReq.CreateTrainingJobReq) (int64, error) {
	switch req.FrameworkType {
	case consts.FrameworkPyTorchDDP:
		if req.InstanceCount < 2 {
			return 0, errors.New("DDP 模式实例数量必须大于 1")
		}
		req.WorkerCount = req.InstanceCount
		return req.InstanceCount, nil
	case consts.FrameworkMPI:
		if req.InstanceCount < 2 {
			return 0, errors.New("MPI 模式实例数量必须大于 1")
		}
		req.WorkerCount = req.InstanceCount
		return req.InstanceCount, nil
	default:
		req.InstanceCount = 1
		req.WorkerCount = 1
		return 1, nil
	}
}

func calculateTrainingTotalGPU(product *productModel.Product, computeNodes int64) int64 {
	var gpuPerNode int64
	if product.IsGPUOnly() {
		gpuPerNode = product.GPUCount
	} else if product.IsVGPU() {
		gpuPerNode = product.VGPUNumber
	}
	return gpuPerNode * computeNodes
}

func (s *TrainingJobService) reserveTrainingCapacity(
	ctx context.Context,
	plan *productService.AllocationPlan,
	cleanups *Cleanups,
) error {
	productSvc := &productService.ProductService{}
	for _, reservation := range plan.Reservations {
		reserve, err := productSvc.ReserveCapacityWithCount(ctx, reservation.ProductID, reservation.Count)
		if err != nil {
			return errors.Wrap(err, "锁定产品资源失败")
		}

		productID := reservation.ProductID
		resourceCount := reserve.ResourceCount
		cleanups.Add(func(ctx context.Context) {
			_ = (&productService.ProductService{}).ReleaseCapacity(ctx, productID, resourceCount)
		})
	}
	return nil
}

func (s *TrainingJobService) getTrainingCluster(clusterID uint) (*global.ClusterClientInfo, error) {
	cluster := global.GVA_K8S_CLUSTER_MANAGER.GetCluster(clusterID)
	if cluster == nil {
		return nil, errors.New("集群不存在")
	}
	return cluster, nil
}

func (s *TrainingJobService) buildTrainingJobRecord(
	req *trainingReq.CreateTrainingJobReq,
	jobName, displayName string,
	product *productModel.Product,
	computeNodes int64,
) *trainingModel.TrainingJob {
	return &trainingModel.TrainingJob{
		UserId:             req.UserId,
		DisplayName:        displayName,
		JobName:            jobName,
		Namespace:          req.Namespace,
		ClusterID:          product.ClusterId,
		FrameworkType:      req.FrameworkType,
		ImageId:            req.ImageId,
		StartupCommand:     req.StartupCommand,
		TotalGPUCount:      calculateTrainingTotalGPU(product, computeNodes),
		ProductId:          req.TemplateProductId,
		InstanceCount:      computeNodes,
		ScheduleStrategy:   req.ScheduleStrategy,
		WorkerCount:        computeNodes,
		Status:             consts.TrainingStatusSubmitted,
		EnableTensorboard:  req.EnableTensorboard,
		TensorboardLogPath: req.TensorboardLogPath,
		PayType:            req.PayType,
		Price:              product.GetPrice(req.PayType),
	}
}

func normalizeTrainingCreateRequest(req *trainingReq.CreateTrainingJobReq) error {
	if req.TemplateProductId == 0 {
		req.TemplateProductId = req.ResourceId
	}
	if req.ScheduleStrategy == "" {
		req.ScheduleStrategy = consts.ScheduleStrategyBalanced
	}
	if req.InstanceCount <= 0 {
		if req.WorkerCount > 0 {
			req.InstanceCount = req.WorkerCount
		} else {
			req.InstanceCount = 1
		}
	}
	return nil
}

func (s *TrainingJobService) saveTrainingAllocations(
	ctx context.Context,
	state *trainingCreateState,
	cleanups *Cleanups,
) error {
	allocations := buildTrainingAllocations(state.trainingJob, state.plan)
	if len(allocations) == 0 {
		return nil
	}
	if err := (&productService.ProductService{}).SaveResourceAllocations(ctx, allocations); err != nil {
		return errors.Wrap(err, "保存训练资源分配失败")
	}

	instanceID := state.trainingJob.ID
	cleanups.Add(func(ctx context.Context) {
		_ = (&productService.ProductService{}).DeleteResourceAllocations(ctx, consts.TrainingInstance, instanceID)
	})
	return nil
}

func buildTrainingAllocations(job *trainingModel.TrainingJob, plan *productService.AllocationPlan) []productModel.ResourceAllocation {
	if plan == nil {
		return nil
	}
	allocations := make([]productModel.ResourceAllocation, 0, job.InstanceCount)
	replicaIndex := 0
	for _, reservation := range plan.Reservations {
		for count := int64(0); count < reservation.Count; count++ {
			role := consts.TaskSpecWorker
			if job.FrameworkType == consts.FrameworkPyTorchDDP && replicaIndex == 0 {
				role = consts.TaskSpecMaster
			}
			if job.FrameworkType == consts.FrameworkMPI {
				role = consts.TaskSpecMPIWorker
			}
			allocations = append(allocations, productModel.ResourceAllocation{
				InstanceType:      consts.TrainingInstance,
				InstanceID:        job.ID,
				ClusterID:         job.ClusterID,
				TemplateProductID: job.ProductId,
				ProductID:         reservation.ProductID,
				NodeName:          reservation.NodeName,
				ScheduleStrategy:  plan.ScheduleStrategy,
				ReplicaIndex:      replicaIndex,
				TaskRole:          role,
				ReservedCount:     1,
			})
			replicaIndex++
		}
	}
	return allocations
}

func allocationNodesForInstance(ctx context.Context, instanceType string, instanceID uint) []string {
	nodes := make([]string, 0)
	_ = global.GVA_DB.WithContext(ctx).
		Table("resource_allocations").
		Distinct("node_name").
		Where("instance_type = ? AND instance_id = ?", instanceType, instanceID).
		Pluck("node_name", &nodes).Error
	return nodes
}

func (s *TrainingJobService) persistTrainingJob(ctx context.Context, state *trainingCreateState) error {
	if err := global.GVA_DB.WithContext(ctx).Create(state.trainingJob).Error; err != nil {
		return errors.Wrap(err, "创建训练任务记录失败")
	}
	return nil
}

func (s *TrainingJobService) saveTrainingJobMounts(ctx context.Context, state *trainingCreateState) error {
	jobMounts, err := s.buildTrainingJobMounts(ctx, state.trainingJob.ID, state.req)
	if err != nil {
		return err
	}
	if len(jobMounts) == 0 {
		return nil
	}

	if err := global.GVA_DB.WithContext(ctx).Create(&jobMounts).Error; err != nil {
		return errors.Wrap(err, "保存挂载配置失败")
	}
	return nil
}

func (s *TrainingJobService) buildTrainingJobMounts(
	ctx context.Context,
	jobID uint,
	req *trainingReq.CreateTrainingJobReq,
) ([]trainingModel.TrainingJobMount, error) {
	pvcMap, err := s.loadTrainingPVCs(ctx, req.Mounts)
	if err != nil {
		return nil, err
	}

	jobMounts := make([]trainingModel.TrainingJobMount, 0, len(req.Mounts))
	for i := range req.Mounts {
		mount := &req.Mounts[i]
		vol, ok := pvcMap[mount.PvcId]
		if !ok {
			return nil, errors.Errorf("PVC不存在: %d", mount.PvcId)
		}

		mount.PvcName = vol.PVCName
		jobMounts = append(jobMounts, trainingModel.TrainingJobMount{
			JobId:     jobID,
			MountType: mount.MountType,
			SourceId:  mount.SourceId,
			PvcId:     mount.PvcId,
			PvcName:   mount.PvcName,
			SubPath:   mount.SubPath,
			MountPath: mount.MountPath,
			ReadOnly:  mount.ReadOnly,
		})
	}

	return jobMounts, nil
}

func (s *TrainingJobService) loadTrainingPVCs(
	ctx context.Context,
	mounts []trainingReq.CreateTrainingJobMountReq,
) (map[uint]pvcModel.Volume, error) {
	pvcIDs := collectTrainingPVCIDs(mounts)
	if len(pvcIDs) == 0 {
		return map[uint]pvcModel.Volume{}, nil
	}

	var volumes []pvcModel.Volume
	if err := global.GVA_DB.WithContext(ctx).Where("id IN ?", pvcIDs).Find(&volumes).Error; err != nil {
		return nil, errors.Wrap(err, "获取PVC信息失败")
	}

	pvcMap := make(map[uint]pvcModel.Volume, len(volumes))
	for _, vol := range volumes {
		pvcMap[vol.ID] = vol
	}
	return pvcMap, nil
}

func collectTrainingPVCIDs(mounts []trainingReq.CreateTrainingJobMountReq) []uint {
	seen := make(map[uint]struct{})
	ids := make([]uint, 0, len(mounts))
	for _, mount := range mounts {
		if mount.PvcId == 0 {
			continue
		}
		if _, ok := seen[mount.PvcId]; ok {
			continue
		}
		seen[mount.PvcId] = struct{}{}
		ids = append(ids, mount.PvcId)
	}
	return ids
}

func (s *TrainingJobService) saveTrainingJobEnvs(
	ctx context.Context,
	jobID uint,
	envs []trainingReq.CreateTrainingJobEnvReq,
) error {
	if len(envs) == 0 {
		return nil
	}

	jobEnvs := make([]trainingModel.TrainingJobEnv, 0, len(envs))
	for _, env := range envs {
		jobEnvs = append(jobEnvs, trainingModel.TrainingJobEnv{
			JobId: jobID,
			Name:  env.Name,
			Value: env.Value,
		})
	}

	if err := global.GVA_DB.WithContext(ctx).Create(&jobEnvs).Error; err != nil {
		return errors.Wrap(err, "保存环境变量失败")
	}
	return nil
}

func (s *TrainingJobService) submitTrainingJob(
	ctx context.Context,
	state *trainingCreateState,
	cleanups *Cleanups,
) error {
	jobSpec := s.buildJobSpec(ctx, state.trainingJob, state.product, state.image, state.req)
	jobBuilder := builder.NewJobBuilder(state.req.FrameworkType)
	volcanoJob, err := jobBuilder.Build(jobSpec)
	if err != nil {
		s.updateJobStatus(state.trainingJob.ID, consts.TrainingStatusFailed, err.Error())
		return errors.Wrap(err, "构建 Volcano Job 失败")
	}

	volcanoJob.Name = state.jobName
	volcanoJob.Namespace = state.req.Namespace

	createdJob, err := state.cluster.VolcanoClient.Clientset().BatchV1alpha1().Jobs(state.req.Namespace).Create(ctx, volcanoJob, metav1.CreateOptions{})
	if err != nil {
		s.updateJobStatus(state.trainingJob.ID, consts.TrainingStatusFailed, err.Error())
		return errors.Wrap(err, "创建 Volcano Job 失败")
	}

	cleanups.Add(func(cleanupCtx context.Context) {
		if delErr := state.cluster.VolcanoClient.Clientset().BatchV1alpha1().Jobs(state.req.Namespace).Delete(
			cleanupCtx, state.jobName, metav1.DeleteOptions{},
		); delErr != nil {
			logx.Error("回滚时删除 Volcano Job 失败", logx.Field("err", delErr), logx.Field("job", state.jobName))
		} else {
			logx.Info("回滚时删除 Volcano Job 成功", logx.Field("job", state.jobName))
		}
	})

	if err := s.markTrainingJobSubmitted(ctx, state.trainingJob, string(createdJob.UID), state.req); err != nil {
		return err
	}
	return nil
}

func (s *TrainingJobService) markTrainingJobSubmitted(
	ctx context.Context,
	job *trainingModel.TrainingJob,
	jobUID string,
	req *trainingReq.CreateTrainingJobReq,
) error {
	if err := global.GVA_DB.WithContext(ctx).Model(&trainingModel.TrainingJob{}).Where("id = ?", job.ID).Updates(map[string]interface{}{
		"k8s_job_uid":          jobUID,
		"status":               consts.TrainingStatusPending,
		"enable_tensorboard":   req.EnableTensorboard,
		"tensorboard_log_path": req.TensorboardLogPath,
	}).Error; err != nil {
		return errors.Wrap(err, "更新训练任务状态失败")
	}

	job.K8sJobUid = jobUID
	job.Status = consts.TrainingStatusPending
	job.EnableTensorboard = req.EnableTensorboard
	job.TensorboardLogPath = req.TensorboardLogPath
	return nil
}

func (s *TrainingJobService) createTrainingOrder(
	ctx context.Context,
	state *trainingCreateState,
	cleanups *Cleanups,
) error {
	orderSvc := &orderService.OrderService{}
	order, err := orderSvc.CreateOrder(ctx, &orderReq.CreateOrderRequest{
		UserId:       state.req.UserId,
		ProduceType:  orderModel.ProductTypeCompute,
		ResourceType: orderModel.OrderTypeTraining,
		ResourceId:   state.trainingJob.ID,
		ProductId:    state.req.TemplateProductId,
		ImageId:      state.req.ImageId,
		PayType:      consts.PayMethodToInt64[consts.PayMethodBalance],
		ChargeType:   state.req.PayType,
		Quantity:     int(state.req.InstanceCount),
		Area:         state.cluster.Area,
		ClusterId:    state.product.ClusterId,
		Remark:       fmt.Sprintf("创建训练任务: %s", state.displayName),
	})
	if err != nil {
		logx.Error("创建训练任务订单失败", err)
		return errors.Wrap(err, "创建订单失败")
	}

	if err := s.updateTrainingJobOrderID(ctx, state.trainingJob.ID, order.ID); err != nil {
		return err
	}
	state.trainingJob.OrderId = order.ID

	cleanups.Add(func(cleanupCtx context.Context) {
		_ = global.GVA_DB.Model(&orderModel.Order{}).
			Where("id = ?", order.ID).
			Update("status", orderModel.OrderStatusStopped).Error
	})

	return nil
}

func (s *TrainingJobService) updateTrainingJobOrderID(ctx context.Context, jobID, orderID uint) error {
	if err := global.GVA_DB.WithContext(ctx).Model(&trainingModel.TrainingJob{}).
		Where("id = ?", jobID).
		Update("order_id", orderID).Error; err != nil {
		return errors.Wrap(err, "更新训练任务订单ID失败")
	}
	return nil
}

// buildJobSpec 构建 Job 规格（资源配置从产品获取，不依赖 TrainingJob 冗余字段）
func (s *TrainingJobService) buildJobSpec(ctx context.Context, job *trainingModel.TrainingJob, product *productModel.Product, image *imageModel.Image, req *trainingReq.CreateTrainingJobReq) *trainingReq.TrainingJobSpec {
	// 解析启动命令
	// 使用 /bin/sh -c 方式运行，以支持多行脚本和 shell 特性（如重定向、管道等）
	var command, args []string
	if job.StartupCommand != "" {
		command = []string{"/bin/sh", "-c"}
		args = []string{job.StartupCommand}
	}

	// 构建卷和挂载
	volumes := []corev1.Volume{}
	volumeMounts := []corev1.VolumeMount{}
	for i, mount := range req.Mounts {
		volName := fmt.Sprintf("vol-%d", i)
		volumes = append(volumes, corev1.Volume{
			Name: volName,
			VolumeSource: corev1.VolumeSource{
				PersistentVolumeClaim: &corev1.PersistentVolumeClaimVolumeSource{
					ClaimName: mount.PvcName,
				},
			},
		})
		volumeMounts = append(volumeMounts, corev1.VolumeMount{
			Name:      volName,
			MountPath: mount.MountPath,
			SubPath:   mount.SubPath,
			ReadOnly:  mount.ReadOnly,
		})
	}

	// 构建环境变量
	envs := []corev1.EnvVar{}
	for _, env := range req.Envs {
		envs = append(envs, corev1.EnvVar{
			Name:  env.Name,
			Value: env.Value,
		})
	}

	// 计算节点数（用于 SHM 大小）
	computeNodes := job.WorkerCount
	if job.FrameworkType != consts.FrameworkPyTorchDDP && job.FrameworkType != consts.FrameworkMPI {
		computeNodes = 1
	}

	return &trainingReq.TrainingJobSpec{
		Name:         job.JobName,
		Namespace:    job.Namespace,
		Framework:    job.FrameworkType,
		Image:        image.ImageAddr,
		Command:      command,
		Args:         args,
		WorkerCount:  job.WorkerCount,
		AllowedNodes: allocationNodesForInstance(ctx, consts.TrainingInstance, job.ID),
		StrictSpread: job.ScheduleStrategy == consts.ScheduleStrategyStrict,
		Product: trainingReq.ProductSpec{
			CPU:        product.CPU,
			Memory:     product.Memory,
			GPUModel:   product.GPUModel,
			GPUCount:   product.GPUCount,
			VGPUNumber: product.VGPUNumber,
			VGPUMemory: product.VGPUMemory,
			VGPUCores:  product.VGPUCores,
		},
		Volumes:      volumes,
		VolumeMounts: volumeMounts,
		Envs:         envs,
		UseSHM:       true,
		SHMSize:      product.Memory * computeNodes,
		Labels: map[string]string{
			consts.LabelApp:          consts.TrainingInstance,
			consts.LabelInstanceType: consts.TrainingInstance,
			consts.LabelJobID:        fmt.Sprintf("%d", job.ID),
			"neptune.io/instance":    job.JobName,
		},
		MaxRetry: 3,
	}
}

// DeleteTrainingJob 删除训练任务
// 注意：按量计费的结算和资源释放由 PodGroup Informer 自动处理
// 删除 Volcano Job 会触发 PodGroup 删除，Informer 会捕获该事件并调用 StopOrder 和释放资源
func (s *TrainingJobService) DeleteTrainingJob(ctx context.Context, req *trainingReq.DeleteTrainingJobReq) error {
	job, cluster, err := s.getTrainingJobAndCluster(ctx, req.ID)
	if err != nil {
		return err
	}

	defer s.cleanupTrainingJobOptionalResources(job, cluster)

	if err := s.deleteTrainingVolcanoJob(ctx, cluster, job); err != nil {
		return errors.Wrap(err, "删除 Volcano Job 失败")
	}

	global.GVA_DB.WithContext(ctx).Model(&trainingModel.TrainingJob{}).Where("id = ?", job.ID).Updates(map[string]interface{}{
		"status": consts.TrainingStatusKilled,
	})

	if err := global.GVA_DB.WithContext(ctx).Delete(job).Error; err != nil {
		return errors.Wrap(err, "删除任务记录失败")
	}

	logx.Info("删除训练任务成功", logx.Field("job", job.JobName))
	return nil
}

// StopTrainingJob 停止训练任务
// 注意：按量计费的结算和资源释放由 PodGroup Informer 自动处理
// 删除 Volcano Job 会触发 PodGroup 删除，Informer 会捕获该事件并调用 StopOrder 和释放资源
func (s *TrainingJobService) StopTrainingJob(ctx context.Context, req *trainingReq.StopTrainingJobReq) error {
	job, cluster, err := s.getTrainingJobAndCluster(ctx, req.ID)
	if err != nil {
		return err
	}

	global.GVA_DB.WithContext(ctx).Model(&trainingModel.TrainingJob{}).Where("id = ?", job.ID).Update("status", consts.TrainingStatusKilling)
	defer s.cleanupTrainingJobOptionalResources(job, cluster)

	if err := s.deleteTrainingVolcanoJob(ctx, cluster, job); err != nil {
		return errors.Wrap(err, "停止任务失败")
	}

	now := time.Now()
	global.GVA_DB.WithContext(ctx).Model(&trainingModel.TrainingJob{}).Where("id = ?", job.ID).Updates(map[string]interface{}{
		"status":      consts.TrainingStatusKilled,
		"finished_at": &now,
	})

	logx.Info("停止训练任务成功", logx.Field("job", job.JobName))
	return nil
}

func (s *TrainingJobService) getTrainingJobAndCluster(
	ctx context.Context,
	jobID uint,
) (*trainingModel.TrainingJob, *global.ClusterClientInfo, error) {
	job := &trainingModel.TrainingJob{}
	if err := global.GVA_DB.WithContext(ctx).Where("id = ?", jobID).First(job).Error; err != nil {
		return nil, nil, errors.Wrap(err, "任务不存在")
	}

	cluster, err := s.getTrainingCluster(job.ClusterID)
	if err != nil {
		return nil, nil, err
	}

	return job, cluster, nil
}

func (s *TrainingJobService) cleanupTrainingJobOptionalResources(
	job *trainingModel.TrainingJob,
	cluster *global.ClusterClientInfo,
) {
	cleanupCtx := context.Background()
	if err := s.deleteOptionalResources(cleanupCtx, job, cluster); err != nil {
		logx.Error("清理可选资源失败", logx.Field("err", err))
	}
}

func (s *TrainingJobService) deleteTrainingVolcanoJob(
	ctx context.Context,
	cluster *global.ClusterClientInfo,
	job *trainingModel.TrainingJob,
) error {
	err := cluster.VolcanoClient.Clientset().BatchV1alpha1().Jobs(job.Namespace).Delete(ctx, job.JobName, metav1.DeleteOptions{
		PropagationPolicy: func() *metav1.DeletionPropagation { p := metav1.DeletePropagationBackground; return &p }(),
	})
	if err != nil && !apierrors.IsNotFound(err) {
		return err
	}
	return nil
}

// GetTrainingJobList 获取训练任务列表
func (s *TrainingJobService) GetTrainingJobList(ctx context.Context, req *trainingReq.GetTrainingJobListReq) (*trainingResp.GetTrainingJobListResp, error) {
	jobs, total, err := s.listTrainingJobs(ctx, req)
	if err != nil {
		return nil, err
	}

	lookup := s.loadTrainingListLookup(ctx, jobs)
	items := make([]trainingResp.TrainingJobItem, 0, len(jobs))
	for i := range jobs {
		s.syncJobStatus(ctx, &jobs[i])
		items = append(items, s.convertToListItem(&jobs[i], lookup))
	}

	return &trainingResp.GetTrainingJobListResp{
		List:  items,
		Total: total,
	}, nil
}

func (s *TrainingJobService) buildTrainingJobListQuery(ctx context.Context, req *trainingReq.GetTrainingJobListReq) *gorm.DB {
	db := global.GVA_DB.WithContext(ctx).Model(&trainingModel.TrainingJob{}).Where("user_id = ?", req.UserId)

	if req.Name != "" {
		db = db.Where("job_name LIKE ?", "%"+req.Name+"%")
	}
	if req.Status != "" {
		db = db.Where("status = ?", req.Status)
	}
	if req.ClusterId > 0 {
		db = db.Where("cluster_id = ?", req.ClusterId)
	}

	return db
}

func (s *TrainingJobService) listTrainingJobs(
	ctx context.Context,
	req *trainingReq.GetTrainingJobListReq,
) ([]trainingModel.TrainingJob, int64, error) {
	var (
		jobs  []trainingModel.TrainingJob
		total int64
	)

	db := s.buildTrainingJobListQuery(ctx, req)
	if err := db.Count(&total).Error; err != nil {
		return nil, 0, errors.Wrap(err, "查询任务总数失败")
	}

	page, pageSize := normalizeTrainingListPage(req)
	offset := (page - 1) * pageSize
	if err := db.Order("created_at DESC").Offset(offset).Limit(pageSize).Find(&jobs).Error; err != nil {
		return nil, 0, errors.Wrap(err, "查询任务列表失败")
	}

	return jobs, total, nil
}

func normalizeTrainingListPage(req *trainingReq.GetTrainingJobListReq) (int, int) {
	page := req.Page
	pageSize := req.PageSize
	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = 10
	}
	if pageSize > 100 {
		pageSize = 100
	}
	return page, pageSize
}

func (s *TrainingJobService) loadTrainingListLookup(ctx context.Context, jobs []trainingModel.TrainingJob) trainingListLookup {
	return trainingListLookup{
		images:   s.loadTrainingImages(ctx, jobs),
		clusters: s.loadTrainingClusters(ctx, jobs),
	}
}

func (s *TrainingJobService) loadTrainingImages(ctx context.Context, jobs []trainingModel.TrainingJob) map[uint]imageModel.Image {
	imageIDs := collectTrainingJobImageIDs(jobs)
	if len(imageIDs) == 0 {
		return map[uint]imageModel.Image{}
	}

	var images []imageModel.Image
	if err := global.GVA_DB.WithContext(ctx).Where("id IN ?", imageIDs).Find(&images).Error; err != nil {
		logx.Error("批量查询训练任务镜像失败", logx.Field("err", err))
		return map[uint]imageModel.Image{}
	}

	imageMap := make(map[uint]imageModel.Image, len(images))
	for _, image := range images {
		imageMap[image.ID] = image
	}
	return imageMap
}

func (s *TrainingJobService) loadTrainingClusters(ctx context.Context, jobs []trainingModel.TrainingJob) map[uint]clusterModel.K8sCluster {
	clusterIDs := collectTrainingJobClusterIDs(jobs)
	if len(clusterIDs) == 0 {
		return map[uint]clusterModel.K8sCluster{}
	}

	var clusters []clusterModel.K8sCluster
	if err := global.GVA_DB.WithContext(ctx).Where("id IN ?", clusterIDs).Find(&clusters).Error; err != nil {
		logx.Error("批量查询训练任务集群失败", logx.Field("err", err))
		return map[uint]clusterModel.K8sCluster{}
	}

	clusterMap := make(map[uint]clusterModel.K8sCluster, len(clusters))
	for _, cluster := range clusters {
		clusterMap[cluster.ID] = cluster
	}
	return clusterMap
}

func collectTrainingJobImageIDs(jobs []trainingModel.TrainingJob) []uint {
	seen := make(map[uint]struct{})
	ids := make([]uint, 0, len(jobs))
	for _, job := range jobs {
		if job.ImageId == 0 {
			continue
		}
		if _, ok := seen[job.ImageId]; ok {
			continue
		}
		seen[job.ImageId] = struct{}{}
		ids = append(ids, job.ImageId)
	}
	return ids
}

func collectTrainingJobClusterIDs(jobs []trainingModel.TrainingJob) []uint {
	seen := make(map[uint]struct{})
	ids := make([]uint, 0, len(jobs))
	for _, job := range jobs {
		if job.ClusterID == 0 {
			continue
		}
		if _, ok := seen[job.ClusterID]; ok {
			continue
		}
		seen[job.ClusterID] = struct{}{}
		ids = append(ids, job.ClusterID)
	}
	return ids
}

// convertToListItem 转换为列表项
func (s *TrainingJobService) convertToListItem(job *trainingModel.TrainingJob, lookup trainingListLookup) trainingResp.TrainingJobItem {
	item := trainingResp.TrainingJobItem{
		ID:            job.ID,
		DisplayName:   job.DisplayName,
		JobName:       job.JobName,
		Namespace:     job.Namespace,
		FrameworkType: job.FrameworkType,
		Status:        job.Status,
		TotalGPUCount: job.TotalGPUCount,
		WorkerCount:   job.WorkerCount,
		CreatedAt:     job.CreatedAt,
		StartedAt:     job.StartedAt,
		FinishedAt:    job.FinishedAt,
		Duration:      buildTrainingDuration(job.StartedAt, job.FinishedAt),
		ErrorMsg:      job.ErrorMsg,
	}

	if image, ok := lookup.images[job.ImageId]; ok {
		item.Image = image.ImageAddr
		item.ImageName = image.Name
	}
	if cluster, ok := lookup.clusters[job.ClusterID]; ok {
		item.ClusterName = cluster.Name
	}
	if job.EnableTensorboard {
		item.TensorboardUrl = buildTrainingTensorboardURL(job.Namespace, job.JobName)
	}

	return item
}

func buildTrainingDuration(startedAt, finishedAt *time.Time) string {
	if startedAt == nil {
		return ""
	}

	endTime := time.Now()
	if finishedAt != nil {
		endTime = *finishedAt
	}
	return timer.FormatDuration(endTime.Sub(*startedAt))
}

func buildTrainingTensorboardURL(namespace, jobName string) string {
	return fmt.Sprintf("/tensorboard/%s/%s/", namespace, jobName)
}

// GetTrainingJobDetail 获取训练任务详情
func (s *TrainingJobService) GetTrainingJobDetail(ctx context.Context, req *trainingReq.GetTrainingJobDetailReq) (*trainingResp.TrainingJobDetail, error) {
	job := &trainingModel.TrainingJob{}
	if err := global.GVA_DB.WithContext(ctx).Where("id = ?", req.ID).First(job).Error; err != nil {
		return nil, errors.Wrap(err, "任务不存在")
	}

	s.syncJobStatus(ctx, job)

	detail := &trainingResp.TrainingJobDetail{
		ID:                 job.ID,
		DisplayName:        job.DisplayName,
		JobName:            job.JobName,
		Namespace:          job.Namespace,
		ClusterID:          job.ClusterID,
		FrameworkType:      job.FrameworkType,
		ImageId:            job.ImageId,
		StartupCommand:     job.StartupCommand,
		Status:             job.Status,
		TotalGPUCount:      job.TotalGPUCount,
		ProductId:          job.ProductId,
		WorkerCount:        job.WorkerCount,
		K8sJobUid:          job.K8sJobUid,
		ErrorMsg:           job.ErrorMsg,
		CreatedAt:          job.CreatedAt,
		StartedAt:          job.StartedAt,
		FinishedAt:         job.FinishedAt,
		Duration:           buildTrainingDuration(job.StartedAt, job.FinishedAt),
		TensorboardLogPath: job.TensorboardLogPath,
		EnableTensorboard:  job.EnableTensorboard,
		PayType:            job.PayType,
		Price:              job.Price,
	}

	if err := s.fillTrainingClusterDetail(ctx, detail, job.ClusterID); err != nil {
		logx.Error("查询训练任务集群信息失败", logx.Field("err", err), logx.Field("job", job.JobName))
	}
	if err := s.fillTrainingProductDetail(ctx, detail, job); err != nil {
		logx.Error("查询训练任务产品信息失败", logx.Field("err", err), logx.Field("job", job.JobName))
	}
	if err := s.fillTrainingImageDetail(ctx, detail, job.ImageId); err != nil {
		logx.Error("查询训练任务镜像信息失败", logx.Field("err", err), logx.Field("job", job.JobName))
	}
	if detail.ImageName == "" {
		detail.ImageName = detail.Image
	}
	if mounts, err := s.loadTrainingJobMountDetails(ctx, job.ID); err == nil {
		detail.Mounts = mounts
	} else {
		logx.Error("查询训练任务挂载信息失败", logx.Field("err", err), logx.Field("job", job.JobName))
	}
	if envs, err := s.loadTrainingJobEnvDetails(ctx, job.ID); err == nil {
		detail.Envs = envs
	} else {
		logx.Error("查询训练任务环境变量失败", logx.Field("err", err), logx.Field("job", job.JobName))
	}
	if job.EnableTensorboard {
		detail.TensorboardUrl = buildTrainingTensorboardURL(job.Namespace, job.JobName)
	}

	return detail, nil
}

func (s *TrainingJobService) fillTrainingClusterDetail(
	ctx context.Context,
	detail *trainingResp.TrainingJobDetail,
	clusterID uint,
) error {
	var cluster clusterModel.K8sCluster
	if err := global.GVA_DB.WithContext(ctx).Where("id = ?", clusterID).First(&cluster).Error; err != nil {
		return err
	}
	detail.ClusterName = cluster.Name
	detail.Area = cluster.Area
	return nil
}

func (s *TrainingJobService) fillTrainingProductDetail(
	ctx context.Context,
	detail *trainingResp.TrainingJobDetail,
	job *trainingModel.TrainingJob,
) error {
	if job.ProductId == 0 {
		return nil
	}

	var product productModel.Product
	if err := global.GVA_DB.WithContext(ctx).Where("id = ?", job.ProductId).First(&product).Error; err != nil {
		return err
	}

	detail.CPU = product.CPU
	detail.Memory = product.Memory
	detail.GPUModel = product.GPUModel
	if product.IsGPUOnly() {
		detail.GPUType = consts.NvidiaGPUType
		detail.WorkerGPU = product.GPUCount
	} else if product.IsVGPU() {
		detail.GPUType = consts.VGPUType
		detail.WorkerGPU = product.VGPUNumber
	}
	if detail.Price == 0 {
		detail.Price = product.GetPrice(int64(job.PayType))
	}

	return nil
}

func (s *TrainingJobService) fillTrainingImageDetail(
	ctx context.Context,
	detail *trainingResp.TrainingJobDetail,
	imageID uint,
) error {
	var image imageModel.Image
	if err := global.GVA_DB.WithContext(ctx).Where("id = ?", imageID).First(&image).Error; err != nil {
		return err
	}
	detail.Image = image.ImageAddr
	detail.ImageName = image.Name
	if detail.ImageName == "" {
		detail.ImageName = detail.Image
	}
	return nil
}

func (s *TrainingJobService) loadTrainingJobMountDetails(
	ctx context.Context,
	jobID uint,
) ([]trainingResp.TrainingJobMountDetail, error) {
	var mounts []trainingModel.TrainingJobMount
	if err := global.GVA_DB.WithContext(ctx).Where("job_id = ?", jobID).Find(&mounts).Error; err != nil {
		return nil, err
	}

	pvcNames, err := s.loadTrainingMountPVCNames(ctx, mounts)
	if err != nil {
		logx.Error("查询训练任务PVC展示名失败", logx.Field("err", err), logx.Field("jobId", jobID))
		pvcNames = map[uint]string{}
	}

	details := make([]trainingResp.TrainingJobMountDetail, 0, len(mounts))
	for _, mount := range mounts {
		pvcName := mount.PvcName
		if displayName, ok := pvcNames[mount.PvcId]; ok && displayName != "" {
			pvcName = displayName
		}
		details = append(details, trainingResp.TrainingJobMountDetail{
			MountType: mount.MountType,
			SourceId:  mount.SourceId,
			PvcName:   pvcName,
			SubPath:   mount.SubPath,
			MountPath: mount.MountPath,
			ReadOnly:  mount.ReadOnly,
		})
	}

	return details, nil
}

func (s *TrainingJobService) loadTrainingMountPVCNames(
	ctx context.Context,
	mounts []trainingModel.TrainingJobMount,
) (map[uint]string, error) {
	pvcIDs := collectTrainingMountPVCIDs(mounts)
	if len(pvcIDs) == 0 {
		return map[uint]string{}, nil
	}

	var volumes []pvcModel.Volume
	if err := global.GVA_DB.WithContext(ctx).Where("id IN ?", pvcIDs).Find(&volumes).Error; err != nil {
		return nil, err
	}

	pvcNames := make(map[uint]string, len(volumes))
	for _, vol := range volumes {
		pvcNames[vol.ID] = vol.Name
	}
	return pvcNames, nil
}

func collectTrainingMountPVCIDs(mounts []trainingModel.TrainingJobMount) []uint {
	seen := make(map[uint]struct{})
	ids := make([]uint, 0, len(mounts))
	for _, mount := range mounts {
		if mount.PvcId == 0 {
			continue
		}
		if _, ok := seen[mount.PvcId]; ok {
			continue
		}
		seen[mount.PvcId] = struct{}{}
		ids = append(ids, mount.PvcId)
	}
	return ids
}

func (s *TrainingJobService) loadTrainingJobEnvDetails(
	ctx context.Context,
	jobID uint,
) ([]trainingResp.TrainingJobEnvDetail, error) {
	var envs []trainingModel.TrainingJobEnv
	if err := global.GVA_DB.WithContext(ctx).Where("job_id = ?", jobID).Find(&envs).Error; err != nil {
		return nil, err
	}

	details := make([]trainingResp.TrainingJobEnvDetail, 0, len(envs))
	for _, env := range envs {
		details = append(details, trainingResp.TrainingJobEnvDetail{
			Name:  env.Name,
			Value: env.Value,
		})
	}
	return details, nil
}

// GetTrainingJobLogs 获取训练任务日志
func (s *TrainingJobService) GetTrainingJobLogs(ctx context.Context, req *trainingReq.GetTrainingJobLogsReq) (io.ReadCloser, error) {
	job, cluster, err := s.getTrainingJobAndCluster(ctx, req.ID)
	if err != nil {
		return nil, err
	}

	taskName := req.TaskName
	if taskName == "" {
		switch job.FrameworkType {
		case consts.FrameworkMPI:
			taskName = consts.TaskSpecMPIMaster
		default:
			taskName = consts.TaskSpecMaster
		}
	}

	podName, err := s.findPodName(ctx, cluster.ClientSet, job.Namespace, job.JobName, taskName, req.PodIndex)
	if err != nil {
		return nil, err
	}

	container := req.Container
	if container == "" {
		container = job.JobName + "-" + taskName
	}

	logOptions := &corev1.PodLogOptions{
		Container:  container,
		Follow:     req.Follow,
		Timestamps: req.Timestamps,
	}
	if req.TailLines > 0 {
		logOptions.TailLines = &req.TailLines
	}

	logReq := cluster.ClientSet.CoreV1().Pods(job.Namespace).GetLogs(podName, logOptions)
	return logReq.Stream(ctx)
}

// findPodName 查找 Pod 名称
func (s *TrainingJobService) findPodName(ctx context.Context, clientSet *kubernetes.Clientset, namespace, jobName, taskName string, podIndex *int) (string, error) {
	labelSelector := fmt.Sprintf("%s=%s", consts.LabelVolcanoJobName, jobName)
	if taskName != "" {
		labelSelector += fmt.Sprintf(",%s=%s", consts.LabelVolcanoTaskSpec, taskName)
	}

	pods, err := clientSet.CoreV1().Pods(namespace).List(ctx, metav1.ListOptions{
		LabelSelector: labelSelector,
	})
	if err != nil {
		return "", errors.Wrap(err, "查找 Pod 失败")
	}

	if len(pods.Items) == 0 {
		return "", errors.New("未找到 Pod")
	}

	// 如果指定了索引，返回对应索引的 Pod
	if podIndex != nil && *podIndex < len(pods.Items) {
		return pods.Items[*podIndex].Name, nil
	}

	// 默认返回第一个 Pod
	return pods.Items[0].Name, nil
}

// GetTrainingJobPods 获取训练任务的 Pod 列表
func (s *TrainingJobService) GetTrainingJobPods(ctx context.Context, id uint) ([]terminalService.PodInfo, error) {
	var job trainingModel.TrainingJob
	if err := global.GVA_DB.Where("id = ?", id).First(&job).Error; err != nil {
		return nil, errors.Wrap(err, "任务不存在")
	}

	labelSelector := fmt.Sprintf("%s=%s", consts.LabelVolcanoJobName, job.JobName)
	return terminalService.TerminalServiceApp.GetPodList(ctx, job.ClusterID, job.Namespace, labelSelector)
}

// updateJobStatus 更新任务状态
func (s *TrainingJobService) updateJobStatus(jobId uint, status, errMsg string) {
	updates := map[string]interface{}{
		"status": status,
	}
	if errMsg != "" {
		updates["error_msg"] = errMsg
	}
	global.GVA_DB.Model(&trainingModel.TrainingJob{}).Where("id = ?", jobId).Updates(updates)
}

// syncJobStatus 同步任务状态，确保 master pod 的运行状态能及时反映到界面上
// 注意：资源释放由 PodGroup Informer 自动处理，这里只同步状态
func (s *TrainingJobService) syncJobStatus(ctx context.Context, job *trainingModel.TrainingJob) {
	if job.Status == consts.TrainingStatusSucceeded || job.Status == consts.TrainingStatusFailed || job.Status == consts.TrainingStatusKilled {
		return
	}

	cluster := global.GVA_K8S_CLUSTER_MANAGER.GetCluster(job.ClusterID)
	if cluster == nil || cluster.PodLister == nil {
		return
	}

	taskName := ""
	switch job.FrameworkType {
	case consts.FrameworkStandalone:
		taskName = consts.TaskSpecWorker
	default:
		taskName = consts.TaskSpecMaster
	}

	selector := labels.SelectorFromSet(labels.Set{
		consts.LabelVolcanoJobName:  job.JobName,
		consts.LabelVolcanoTaskSpec: taskName,
	})

	pods, err := cluster.PodLister.Pods(job.Namespace).List(selector)
	if err != nil || len(pods) == 0 {
		return
	}

	// 取第一个 pod 检查状态
	pod := pods[0]
	newStatus := ""
	switch pod.Status.Phase {
	case corev1.PodSucceeded:
		newStatus = consts.TrainingStatusSucceeded
	case corev1.PodFailed:
		newStatus = consts.TrainingStatusFailed
	}

	if newStatus != "" && newStatus != job.Status {
		job.Status = newStatus
		now := time.Now()
		updates := map[string]interface{}{
			"status":      newStatus,
			"finished_at": now,
		}
		if newStatus == consts.TrainingStatusFailed {
			updates["error_msg"] = pod.Status.Message
			job.ErrorMsg = pod.Status.Message
		}
		job.FinishedAt = &now
		global.GVA_DB.Model(&trainingModel.TrainingJob{}).Where("id = ?", job.ID).Updates(updates)
	}
}
