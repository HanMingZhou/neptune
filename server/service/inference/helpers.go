package inference

import (
	"context"
	"encoding/json"
	"fmt"
	"gin-vue-admin/global"
	"gin-vue-admin/model/consts"
	imageModel "gin-vue-admin/model/image"
	inferenceModel "gin-vue-admin/model/inference"
	inferenceReq "gin-vue-admin/model/inference/request"
	orderModel "gin-vue-admin/model/order"
	orderReq "gin-vue-admin/model/order/request"
	productModel "gin-vue-admin/model/product"
	pvcModel "gin-vue-admin/model/pvc"
	systemModel "gin-vue-admin/model/system"
	orderService "gin-vue-admin/service/order"
	productService "gin-vue-admin/service/product"
	helper "gin-vue-admin/utils/k8s"

	"github.com/pkg/errors"
	"github.com/zeromicro/go-zero/core/logx"
	"gorm.io/gorm"
)

// Cleanups 资源清理函数列表
type Cleanups []func(ctx context.Context)

type inferenceCreateState struct {
	req      *inferenceReq.CreateInferenceServiceReq
	user     *systemModel.SysUser
	image    *imageModel.Image
	product  *productModel.Product
	plan     *productService.AllocationPlan
	modelPVC *pvcModel.Volume
	cluster  *global.ClusterClientInfo
	service  *inferenceModel.Inference
}

type inferenceOrderSpec struct {
	serviceID  uint
	userID     uint
	productID  uint
	imageID    uint
	chargeType int64
	quantity   int
	area       string
	clusterID  uint
	remark     string
}

func (c *Cleanups) Add(fn func(ctx context.Context)) {
	*c = append(*c, fn)
}

func (c Cleanups) Run(ctx context.Context) {
	for i := len(c) - 1; i >= 0; i-- {
		c[i](ctx)
	}
}

// prepareCreateRequest 处理创建请求默认值
func (s *InferenceServiceService) prepareCreateRequest(req *inferenceReq.CreateInferenceServiceReq) error {
	if req.TemplateProductId == 0 {
		req.TemplateProductId = req.ProductId
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
	req.WorkerCount = req.InstanceCount
	if req.ModelMountPath == "" {
		req.ModelMountPath = "/model"
	}
	if req.ServicePort == 0 {
		req.ServicePort = 30000
	}
	if req.TensorParallel == 0 {
		req.TensorParallel = 1
	}
	if req.MaxRestarts == 0 {
		req.MaxRestarts = 3
	}
	if req.AuthType == 0 {
		return errors.New("请选择鉴权方式")
	}
	return s.validateCreateRequest(req)
}

// buildCreateState 构建创建推理服务所需的上下文
func (s *InferenceServiceService) buildCreateState(ctx context.Context, req *inferenceReq.CreateInferenceServiceReq) (*inferenceCreateState, error) {
	if err := s.prepareCreateRequest(req); err != nil {
		return nil, err
	}

	userInfo, err := s.getInferenceUser(req.UserId)
	if err != nil {
		return nil, err
	}

	image, err := s.getInferenceImage(req.ImageId)
	if err != nil {
		return nil, err
	}

	product, err := s.getEnabledProduct(req.TemplateProductId)
	if err != nil {
		return nil, err
	}

	if err = (&orderService.OrderService{}).CheckBalanceSufficient(ctx, req.UserId, req.TemplateProductId, req.PayType, int64(req.InstanceCount)); err != nil {
		return nil, err
	}

	plan, err := (&productService.ProductService{}).PlanAllocations(ctx, product.ID, int64(req.InstanceCount), req.ScheduleStrategy)
	if err != nil {
		return nil, err
	}

	modelPVC, err := s.getInferenceVolume(req.ModelPvcId, "获取模型PVC信息失败")
	if err != nil {
		return nil, err
	}

	instanceName := helper.GenerateInstanceName(consts.InferenceInstance)
	displayName := req.Name
	if displayName == "" {
		displayName = instanceName
	}

	service, err := s.buildInferenceServiceRecord(req, userInfo, product, instanceName, displayName)
	if err != nil {
		return nil, err
	}

	return &inferenceCreateState{
		req:      req,
		user:     userInfo,
		image:    image,
		product:  product,
		plan:     plan,
		modelPVC: modelPVC,
		service:  service,
	}, nil
}

// buildInferenceServiceRecord 构建推理服务数据库记录
func (s *InferenceServiceService) buildInferenceServiceRecord(
	req *inferenceReq.CreateInferenceServiceReq,
	userInfo *systemModel.SysUser,
	product *productModel.Product,
	instanceName, displayName string,
) (*inferenceModel.Inference, error) {
	argsJSON, err := marshalInferenceArgs(req.Args)
	if err != nil {
		return nil, errors.Wrap(err, "序列化启动参数失败")
	}

	return &inferenceModel.Inference{
		DisplayName:      displayName,
		InstanceName:     instanceName,
		UserId:           req.UserId,
		Namespace:        userInfo.Namespace,
		ClusterID:        product.ClusterId,
		DeployType:       req.DeployType,
		Framework:        req.Framework,
		ModelMountPath:   req.ModelMountPath,
		ModelPvcId:       req.ModelPvcId,
		ImageId:          req.ImageId,
		TensorParallel:   req.TensorParallel,
		PipelineParallel: req.PipelineParallel,
		InstanceCount:    req.InstanceCount,
		WorkerCount:      req.InstanceCount,
		ScheduleStrategy: req.ScheduleStrategy,
		ProductId:        req.TemplateProductId,
		CPU:              product.CPU,
		Memory:           product.Memory,
		GPU:              product.GPUCount,
		GPUModel:         product.GPUModel,
		VGPUNumber:       product.VGPUNumber,
		VGPUMemory:       product.VGPUMemory,
		VGPUCores:        product.VGPUCores,
		ServicePort:      req.ServicePort,
		Command:          req.Command,
		Args:             argsJSON,
		AutoRestart:      req.AutoRestart,
		MaxRestarts:      req.MaxRestarts,
		AuthType:         req.AuthType,
		Status:           consts.InferenceStatusCreating,
		PayType:          req.PayType,
		Price:            product.GetPrice(req.PayType),
	}, nil
}

// persistInferenceService 保存推理服务记录
func (s *InferenceServiceService) persistInferenceService(service *inferenceModel.Inference) error {
	if err := global.GVA_DB.Create(service).Error; err != nil {
		return errors.Wrap(err, "创建推理服务记录失败")
	}
	return nil
}

// reserveInferenceCapacity 锁定推理服务资源配额
func (s *InferenceServiceService) reserveInferenceCapacity(ctx context.Context, plan *productService.AllocationPlan, cleanups *Cleanups) error {
	productSvc := &productService.ProductService{}
	for _, reservation := range plan.Reservations {
		reserve, err := productSvc.ReserveCapacityWithCount(ctx, reservation.ProductID, reservation.Count)
		if err != nil {
			return errors.Wrap(err, "锁定产品资源失败")
		}

		productID := reservation.ProductID
		resourceCount := reserve.ResourceCount
		cleanups.Add(func(ctx context.Context) {
			_ = productSvc.ReleaseCapacity(ctx, productID, resourceCount)
		})
	}
	return nil
}

func (s *InferenceServiceService) saveInferenceAllocations(
	ctx context.Context,
	service *inferenceModel.Inference,
	plan *productService.AllocationPlan,
	cleanups *Cleanups,
) error {
	allocations := buildInferenceAllocations(service, plan)
	if len(allocations) == 0 {
		return nil
	}
	if err := (&productService.ProductService{}).SaveResourceAllocations(ctx, allocations); err != nil {
		return errors.Wrap(err, "保存推理资源分配失败")
	}

	instanceID := service.ID
	cleanups.Add(func(ctx context.Context) {
		_ = (&productService.ProductService{}).DeleteResourceAllocations(ctx, consts.InferenceInstance, instanceID)
	})
	return nil
}

// saveInferenceMounts 保存挂载配置
func (s *InferenceServiceService) saveInferenceMounts(serviceID uint, mounts []inferenceReq.CreateInferenceMountReq) error {
	if len(mounts) == 0 {
		return nil
	}

	return global.GVA_DB.Transaction(func(tx *gorm.DB) error {
		for _, mount := range mounts {
			pvc, err := s.getInferenceVolume(mount.PvcId, "获取挂载PVC信息失败")
			if err != nil {
				return err
			}

			record := &inferenceModel.InferenceMount{
				ServiceId: serviceID,
				MountType: mount.MountType,
				PvcId:     mount.PvcId,
				PvcName:   pvc.PVCName,
				SubPath:   mount.SubPath,
				MountPath: mount.MountPath,
				ReadOnly:  mount.ReadOnly,
			}
			if err = tx.Create(record).Error; err != nil {
				return errors.Wrap(err, "保存挂载配置失败")
			}
		}
		return nil
	})
}

// saveInferenceEnvs 保存环境变量
func (s *InferenceServiceService) saveInferenceEnvs(serviceID uint, envs []inferenceReq.CreateInferenceEnvReq) error {
	if len(envs) == 0 {
		return nil
	}

	records := make([]inferenceModel.InferenceEnv, 0, len(envs))
	for _, env := range envs {
		records = append(records, inferenceModel.InferenceEnv{
			ServiceId: serviceID,
			Name:      env.Name,
			Value:     env.Value,
		})
	}

	if err := global.GVA_DB.Create(&records).Error; err != nil {
		return errors.Wrap(err, "保存环境变量失败")
	}
	return nil
}

// loadCreateCluster 加载创建流程所需集群
func (s *InferenceServiceService) loadCreateCluster(state *inferenceCreateState) error {
	cluster := global.GVA_K8S_CLUSTER_MANAGER.GetCluster(state.product.ClusterId)
	if cluster == nil {
		s.updateStatus(state.service.ID, consts.InferenceStatusFailed, "集群不存在")
		return errors.New("集群不存在")
	}

	state.cluster = cluster
	return nil
}

// createInferenceOrder 创建推理服务计费订单
func (s *InferenceServiceService) createInferenceOrder(ctx context.Context, spec *inferenceOrderSpec, cleanups *Cleanups) error {
	orderSvc := &orderService.OrderService{}
	order, err := orderSvc.CreateOrder(ctx, &orderReq.CreateOrderRequest{
		UserId:       spec.userID,
		ProduceType:  orderModel.ProductTypeCompute,
		ResourceType: orderModel.OrderTypeInference,
		ResourceId:   spec.serviceID,
		ProductId:    spec.productID,
		ImageId:      spec.imageID,
		PayType:      consts.PayMethodToInt64[consts.PayMethodBalance],
		ChargeType:   spec.chargeType,
		Quantity:     spec.quantity,
		Area:         spec.area,
		ClusterId:    spec.clusterID,
		Remark:       spec.remark,
	})
	if err != nil {
		logx.Error("创建推理服务订单失败", err)
		return errors.Wrap(err, "创建订单失败")
	}

	cleanups.Add(func(ctx context.Context) {
		_ = global.GVA_DB.Model(&orderModel.Order{}).
			Where("id = ?", order.ID).
			Update("status", orderModel.OrderStatusStopped).Error
	})

	if err = global.GVA_DB.Model(&inferenceModel.Inference{}).
		Where("id = ?", spec.serviceID).
		Update("order_id", order.ID).Error; err != nil {
		return errors.Wrap(err, "更新推理服务订单关联失败")
	}
	return nil
}

// createOrderSpecFromCreateState 构建创建订单参数
func (s *InferenceServiceService) createOrderSpecFromCreateState(state *inferenceCreateState) *inferenceOrderSpec {
	return &inferenceOrderSpec{
		serviceID:  state.service.ID,
		userID:     state.req.UserId,
		productID:  state.req.TemplateProductId,
		imageID:    state.req.ImageId,
		chargeType: state.req.PayType,
		quantity:   state.req.InstanceCount,
		area:       state.cluster.Area,
		clusterID:  state.product.ClusterId,
		remark:     fmt.Sprintf("创建推理服务: %s", state.service.DisplayName),
	}
}

// createOrderSpecFromService 构建重启订单参数
func (s *InferenceServiceService) createOrderSpecFromService(service *inferenceModel.Inference, cluster *global.ClusterClientInfo) *inferenceOrderSpec {
	return &inferenceOrderSpec{
		serviceID:  service.ID,
		userID:     service.UserId,
		productID:  service.ProductId,
		imageID:    service.ImageId,
		chargeType: service.PayType,
		quantity:   service.InstanceCount,
		area:       cluster.Area,
		clusterID:  service.ClusterID,
		remark:     fmt.Sprintf("重启推理服务: %s", service.DisplayName),
	}
}

// getInferenceUser 获取用户信息
func (s *InferenceServiceService) getInferenceUser(userID uint) (*systemModel.SysUser, error) {
	var userInfo systemModel.SysUser
	if err := global.GVA_DB.Where("id = ?", userID).First(&userInfo).Error; err != nil {
		logx.Error("用户不存在", err)
		return nil, errors.New("用户不存在")
	}
	return &userInfo, nil
}

// getInferenceService 获取推理服务
func (s *InferenceServiceService) getInferenceService(id uint) (*inferenceModel.Inference, error) {
	var service inferenceModel.Inference
	if err := global.GVA_DB.Where("id = ?", id).First(&service).Error; err != nil {
		return nil, errors.Wrap(err, "服务不存在")
	}
	return &service, nil
}

// getInferenceImage 获取镜像信息
func (s *InferenceServiceService) getInferenceImage(imageID uint) (*imageModel.Image, error) {
	var image imageModel.Image
	if err := global.GVA_DB.Where("id = ?", imageID).First(&image).Error; err != nil {
		return nil, errors.Wrap(err, "获取镜像信息失败")
	}
	return &image, nil
}

// getEnabledProduct 获取可用产品规格
func (s *InferenceServiceService) getEnabledProduct(productID uint) (*productModel.Product, error) {
	var product productModel.Product
	if err := global.GVA_DB.Where("id = ?", productID).First(&product).Error; err != nil {
		return nil, errors.Wrap(err, "获取产品规格失败")
	}
	if product.Status != productModel.ProductStatusEnabled {
		return nil, errors.New("所选产品已下架，无法使用")
	}
	return &product, nil
}

// getInferenceVolume 获取推理服务使用的 PVC
func (s *InferenceServiceService) getInferenceVolume(volumeID uint, errMsg string) (*pvcModel.Volume, error) {
	var volume pvcModel.Volume
	if err := global.GVA_DB.Where("id = ?", volumeID).First(&volume).Error; err != nil {
		return nil, errors.Wrap(err, errMsg)
	}
	return &volume, nil
}

// getInferenceCluster 获取集群客户端
func (s *InferenceServiceService) getInferenceCluster(clusterID uint) (*global.ClusterClientInfo, error) {
	cluster := global.GVA_K8S_CLUSTER_MANAGER.GetCluster(clusterID)
	if cluster == nil {
		return nil, errors.New("集群不存在")
	}
	return cluster, nil
}

// marshalInferenceArgs 序列化启动参数
func marshalInferenceArgs(args []string) (string, error) {
	if len(args) == 0 {
		return "", nil
	}

	data, err := json.Marshal(args)
	if err != nil {
		return "", err
	}
	return string(data), nil
}

// calculateInferenceQuotaCount 计算推理服务占用配额
func calculateInferenceQuotaCount(deployType string, workerCount int) int64 {
	if deployType == consts.DeployTypeDistributed {
		return int64(workerCount)
	}
	return 1
}

func buildInferenceAllocations(service *inferenceModel.Inference, plan *productService.AllocationPlan) []productModel.ResourceAllocation {
	if plan == nil {
		return nil
	}
	allocations := make([]productModel.ResourceAllocation, 0, service.InstanceCount)
	replicaIndex := 0
	for _, reservation := range plan.Reservations {
		for count := int64(0); count < reservation.Count; count++ {
			role := "standalone"
			if service.DeployType == consts.DeployTypeDistributed {
				if replicaIndex == 0 {
					role = "head"
				} else {
					role = "worker"
				}
			}
			allocations = append(allocations, productModel.ResourceAllocation{
				InstanceType:      consts.InferenceInstance,
				InstanceID:        service.ID,
				ClusterID:         service.ClusterID,
				TemplateProductID: service.ProductId,
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
