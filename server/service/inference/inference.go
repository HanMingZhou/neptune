package inference

import (
	"context"
	"encoding/json"
	"fmt"
	"gin-vue-admin/global"
	apisixReq "gin-vue-admin/model/apisix/request"
	"gin-vue-admin/model/consts"
	imageModel "gin-vue-admin/model/image"
	inferenceModel "gin-vue-admin/model/inference"
	inferenceReq "gin-vue-admin/model/inference/request"
	"gin-vue-admin/model/inference/response"
	"gin-vue-admin/service/inference/builder"
	productService "gin-vue-admin/service/product"
	terminalService "gin-vue-admin/service/terminal"
	helper "gin-vue-admin/utils/k8s"
	"io"
	"strings"

	"github.com/pkg/errors"
	"github.com/zeromicro/go-zero/core/logx"
	corev1 "k8s.io/api/core/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
	"k8s.io/client-go/kubernetes"
)

// InferenceServiceManager 推理服务管理器接口
type InferenceServiceManager interface {
	CreateInferenceService(ctx context.Context, req *inferenceReq.CreateInferenceServiceReq) (*response.CreateInferenceServiceResp, error)
	DeleteInferenceService(ctx context.Context, req *inferenceReq.DeleteInferenceServiceReq) error
	StopInferenceService(ctx context.Context, req *inferenceReq.StopInferenceServiceReq) error
	StartInferenceService(ctx context.Context, req *inferenceReq.StartInferenceServiceReq) error
	GetInferenceServiceList(ctx context.Context, req *inferenceReq.GetInferenceServiceListReq) (*response.GetInferenceServiceListResp, error)
	GetInferenceServiceDetail(ctx context.Context, req *inferenceReq.GetInferenceServiceDetailReq) (*response.InferenceServiceDetail, error)
	GetInferenceServiceLogs(ctx context.Context, req *inferenceReq.GetInferenceServiceLogsReq) (io.ReadCloser, error)
	GetInferenceServicePods(ctx context.Context, id uint) ([]terminalService.PodInfo, error)
}

var _ InferenceServiceManager = (*InferenceServiceService)(nil)

// InferenceServiceService 推理服务实现
type InferenceServiceService struct {
	apisixSvc inferenceApisixService
}

var InferenceServiceServiceApp = new(InferenceServiceService)

// SetApisixService 注入 APISIX 服务依赖
func (s *InferenceServiceService) SetApisixService(svc inferenceApisixService) {
	s.apisixSvc = svc
}

// CreateInferenceService 创建推理服务
func (s *InferenceServiceService) CreateInferenceService(ctx context.Context, req *inferenceReq.CreateInferenceServiceReq) (resp *response.CreateInferenceServiceResp, err error) {
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

	state, err := s.buildCreateState(ctx, req)
	if err != nil {
		return nil, err
	}

	if err = s.persistInferenceService(state.service); err != nil {
		return nil, err
	}

	if err = s.reserveInferenceCapacity(ctx, state.plan, &cleanups); err != nil {
		s.updateStatus(state.service.ID, consts.InferenceStatusFailed, "锁定资源失败")
		return nil, err
	}

	if err = s.saveInferenceAllocations(ctx, state.service, state.plan, &cleanups); err != nil {
		s.updateStatus(state.service.ID, consts.InferenceStatusFailed, err.Error())
		return nil, err
	}

	if err = s.saveInferenceMounts(state.service.ID, state.req.Mounts); err != nil {
		s.updateStatus(state.service.ID, consts.InferenceStatusFailed, err.Error())
		return nil, err
	}

	if err = s.saveInferenceEnvs(state.service.ID, state.req.Envs); err != nil {
		s.updateStatus(state.service.ID, consts.InferenceStatusFailed, err.Error())
		return nil, err
	}

	if err = s.loadCreateCluster(state); err != nil {
		return nil, err
	}

	if err = s.createInferenceOrder(ctx, s.createOrderSpecFromCreateState(state), &cleanups); err != nil {
		s.updateStatus(state.service.ID, consts.InferenceStatusFailed, err.Error())
		return nil, err
	}

	if err = s.createInferenceRuntime(ctx, state, &cleanups); err != nil {
		return nil, err
	}

	if routeErr := s.ensureInferenceRoute(ctx, state.service); routeErr != nil {
		logx.Error("创建推理服务APISIX路由失败", routeErr)
	}

	s.updateStatus(state.service.ID, consts.InferenceStatusPending, "")

	return &response.CreateInferenceServiceResp{
		ID:           state.service.ID,
		InstanceName: state.service.InstanceName,
		Status:       consts.InferenceStatusPending,
		AccessUrl:    buildGatewayUrl(state.user.Namespace, state.service.InstanceName),
	}, nil
}

// validateCreateRequest 校验创建请求
func (s *InferenceServiceService) validateCreateRequest(req *inferenceReq.CreateInferenceServiceReq) error {
	if req.DeployType != consts.DeployTypeDistributed && req.DeployType != consts.DeployTypeStandalone {
		return errors.New("部署类型必须为 DISTRIBUTED 或 STANDALONE")
	}
	// 启动命令必填
	if strings.TrimSpace(req.Command) == "" {
		return errors.New("启动命令不能为空")
	}
	// 分布式模式必须指定框架（用于 NCCL 环境变量注入）
	if req.DeployType == consts.DeployTypeDistributed {
		if req.Framework != consts.FrameworkSGLang && req.Framework != consts.FrameworkVLLM {
			return errors.New("分布式模式下框架必须为 SGLANG 或 VLLM")
		}
		if req.InstanceCount < 2 {
			return errors.New("分布式模式实例数量必须大于 1（总节点数，包含 head）")
		}
	}

	// 校验用户挂载路径不能与模型挂载路径冲突
	modelMount := req.ModelMountPath
	if modelMount == "" {
		modelMount = "/model"
	}
	for _, m := range req.Mounts {
		if m.MountPath == modelMount || strings.HasPrefix(m.MountPath, modelMount+"/") {
			return fmt.Errorf("挂载路径 %s 与模型挂载路径 %s 冲突，请更换挂载路径", m.MountPath, modelMount)
		}
	}

	return nil
}

// buildInferenceSpec 构建推理服务规格（复用于 create 和 restart）
func (s *InferenceServiceService) buildInferenceSpec(ctx context.Context, service *inferenceModel.Inference, imageName, modelPvcName string) *builder.InferenceSpec {
	spec := &builder.InferenceSpec{
		Name:             service.InstanceName,
		Namespace:        service.Namespace,
		InstanceName:     service.InstanceName,
		ServiceID:        service.ID,
		Framework:        service.Framework,
		Image:            imageName,
		ModelMountPath:   service.ModelMountPath,
		ModelPvcName:     modelPvcName,
		TensorParallel:   service.TensorParallel,
		PipelineParallel: service.PipelineParallel,
		WorkerCount:      service.WorkerCount,
		Product: helper.ProductSpec{
			CPU:        service.CPU,
			Memory:     service.Memory,
			GPUCount:   service.GPU,
			GPUModel:   service.GPUModel,
			VGPUNumber: service.VGPUNumber,
			VGPUMemory: service.VGPUMemory,
			VGPUCores:  service.VGPUCores,
		},
		ServicePort:  service.ServicePort,
		AllowedNodes: allocationNodesForInstance(ctx, consts.InferenceInstance, service.ID),
		StrictSpread: service.ScheduleStrategy == consts.ScheduleStrategyStrict,
	}

	// 解析 Command/Args
	// Command 直接存储为原始字符串
	spec.Command = service.Command
	if service.Args != "" {
		var args []string
		json.Unmarshal([]byte(service.Args), &args)
		spec.Args = args
	}

	// 加载用户自定义挂载
	// 如果用户挂载的 PVC 与模型 PVC 相同，复用 model-volume 避免同一 PV 被多次挂载导致 kubelet 超时
	var mounts []inferenceModel.InferenceMount
	global.GVA_DB.Where("service_id = ?", service.ID).Find(&mounts)
	for _, m := range mounts {
		if m.PvcName == modelPvcName {
			// 复用模型 Volume，仅添加 VolumeMount
			vm := corev1.VolumeMount{
				Name:      "model-volume",
				MountPath: m.MountPath,
				ReadOnly:  m.ReadOnly,
			}
			if m.SubPath != "" {
				vm.SubPath = m.SubPath
			}
			spec.UserVolumeMounts = append(spec.UserVolumeMounts, vm)
		} else {
			volName := fmt.Sprintf("user-vol-%d", m.ID)
			spec.UserVolumes = append(spec.UserVolumes, corev1.Volume{
				Name: volName,
				VolumeSource: corev1.VolumeSource{
					PersistentVolumeClaim: &corev1.PersistentVolumeClaimVolumeSource{
						ClaimName: m.PvcName,
						ReadOnly:  m.ReadOnly,
					},
				},
			})
			vm := corev1.VolumeMount{
				Name:      volName,
				MountPath: m.MountPath,
				ReadOnly:  m.ReadOnly,
			}
			if m.SubPath != "" {
				vm.SubPath = m.SubPath
			}
			spec.UserVolumeMounts = append(spec.UserVolumeMounts, vm)
		}
	}

	// 加载用户自定义环境变量
	var envs []inferenceModel.InferenceEnv
	global.GVA_DB.Where("service_id = ?", service.ID).Find(&envs)
	for _, e := range envs {
		spec.UserEnvVars = append(spec.UserEnvVars, corev1.EnvVar{
			Name:  e.Name,
			Value: e.Value,
		})
	}

	return spec
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

// createInferenceRuntime 创建推理服务运行时资源
func (s *InferenceServiceService) createInferenceRuntime(ctx context.Context, state *inferenceCreateState, cleanups *Cleanups) error {
	switch state.req.DeployType {
	case consts.DeployTypeDistributed:
		if err := s.createVCJob(ctx, state.cluster, state.service, state.image, state.modelPVC.PVCName); err != nil {
			s.updateStatus(state.service.ID, consts.InferenceStatusFailed, err.Error())
			return err
		}

		cleanups.Add(func(ctx context.Context) {
			_ = state.cluster.VolcanoClient.Clientset().BatchV1alpha1().Jobs(state.service.Namespace).Delete(
				ctx, state.service.InstanceName, metav1.DeleteOptions{},
			)
		})
	case consts.DeployTypeStandalone:
		if err := s.createDeployment(ctx, state.cluster.ClientSet, state.service, state.image, state.modelPVC.PVCName); err != nil {
			s.updateStatus(state.service.ID, consts.InferenceStatusFailed, err.Error())
			return err
		}

		cleanups.Add(func(ctx context.Context) {
			_ = state.cluster.ClientSet.AppsV1().Deployments(state.service.Namespace).Delete(
				ctx, state.service.InstanceName, metav1.DeleteOptions{},
			)
		})
	default:
		err := errors.New("不支持的部署类型")
		s.updateStatus(state.service.ID, consts.InferenceStatusFailed, err.Error())
		return err
	}

	if err := s.createK8sService(ctx, state.cluster.ClientSet, state.service); err != nil {
		logx.Error("创建K8s Service失败", err)
		return nil
	}

	cleanups.Add(func(ctx context.Context) {
		_ = state.cluster.ClientSet.CoreV1().Services(state.service.Namespace).Delete(
			ctx, state.service.InstanceName, metav1.DeleteOptions{},
		)
	})
	return nil
}

// createVCJob 创建分布式推理 VCJob
func (s *InferenceServiceService) createVCJob(ctx context.Context, cluster *global.ClusterClientInfo, service *inferenceModel.Inference, image *imageModel.Image, modelPvcName string) error {
	if cluster.VolcanoClient == nil {
		return errors.New("集群未配置 Volcano，无法创建分布式推理服务")
	}

	spec := s.buildInferenceSpec(ctx, service, image.ImageAddr, modelPvcName)

	jobBuilder := builder.NewInferenceBuilder(service.Framework)
	job, err := jobBuilder.BuildVCJob(spec)
	if err != nil {
		return errors.Wrap(err, "构建VCJob失败")
	}

	// 通过 VolcanoClient 创建 VCJob
	createdJob, err := cluster.VolcanoClient.Clientset().BatchV1alpha1().Jobs(service.Namespace).Create(ctx, job, metav1.CreateOptions{})
	if err != nil {
		return errors.Wrap(err, "创建VCJob失败")
	}

	// 更新 K8s 资源 UID
	global.GVA_DB.Model(&inferenceModel.Inference{}).Where("id = ?", service.ID).
		Update("k8s_resource_uid", string(createdJob.UID))

	return nil
}

// createDeployment 创建单体推理 Deployment
//
// Deployment 使用 Volcano 调度器（schedulerName: volcano），
// Volcano 会根据 Pod annotations 自动创建 PodGroup，
// PodGroup 删除时 Informer 统一处理资源释放和订单停止。
func (s *InferenceServiceService) createDeployment(ctx context.Context, clientSet *kubernetes.Clientset, service *inferenceModel.Inference, image *imageModel.Image, modelPvcName string) error {
	spec := s.buildInferenceSpec(ctx, service, image.ImageAddr, modelPvcName)

	jobBuilder := builder.NewInferenceBuilder(service.Framework)
	deploy, err := jobBuilder.BuildDeployment(spec)
	if err != nil {
		return errors.Wrap(err, "构建Deployment失败")
	}

	createdDeploy, err := clientSet.AppsV1().Deployments(service.Namespace).Create(ctx, deploy, metav1.CreateOptions{})
	if err != nil {
		return errors.Wrap(err, "创建Deployment失败")
	}

	// 更新 K8s 资源 UID
	global.GVA_DB.Model(&inferenceModel.Inference{}).Where("id = ?", service.ID).
		Update("k8s_resource_uid", string(createdDeploy.UID))

	return nil
}

// createK8sService 创建 K8s Service
func (s *InferenceServiceService) createK8sService(ctx context.Context, clientSet *kubernetes.Clientset, service *inferenceModel.Inference) error {
	selector := map[string]string{
		consts.LabelApp:       service.InstanceName,
		"neptune.io/instance": service.InstanceName,
	}

	// 分布式模式仅指向 head Pod
	if service.DeployType == consts.DeployTypeDistributed {
		selector["role"] = "head"
	}

	svc := &corev1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name:      service.InstanceName,
			Namespace: service.Namespace,
			Labels: map[string]string{
				consts.LabelApp:          consts.InferenceInstance,
				consts.LabelInstanceType: consts.InferenceInstance,
				"neptune.io/instance":    service.InstanceName,
			},
		},
		Spec: corev1.ServiceSpec{
			Selector: selector,
			Ports: []corev1.ServicePort{
				{
					Name:       "http",
					Protocol:   corev1.ProtocolTCP,
					Port:       80,
					TargetPort: intstr.FromInt(service.ServicePort),
				},
			},
			Type: corev1.ServiceTypeClusterIP,
		},
	}

	_, err := clientSet.CoreV1().Services(service.Namespace).Create(ctx, svc, metav1.CreateOptions{})
	if err != nil && !apierrors.IsAlreadyExists(err) {
		return err
	}
	return nil
}

// ensureInferenceRoute 确保推理服务的 APISIX 路由存在（幂等）
// 直接尝试创建，如果已存在（AlreadyExists）则视为成功
func (s *InferenceServiceService) ensureInferenceRoute(ctx context.Context, service *inferenceModel.Inference) error {
	if s.apisixSvc == nil || !global.GVA_CONFIG.Apisix.Enabled {
		return nil
	}

	baseDomain := strings.TrimSpace(global.GVA_CONFIG.Apisix.BaseDomain)
	authUri := global.GVA_CONFIG.Apisix.AuthUri
	if authUri == "" {
		return errors.New("auth-uri 未配置，跳过推理路由创建")
	}

	routeName := fmt.Sprintf("inference-%s", service.InstanceName)
	pathMatch := fmt.Sprintf("/inference/%s/%s/*", service.Namespace, service.InstanceName)
	rewriteRegex := fmt.Sprintf("^/inference/%s/%s/(.*)", service.Namespace, service.InstanceName)

	routeReq := &apisixReq.CreateRouteReq{
		Name:          routeName,
		Namespace:     service.Namespace,
		ClusterId:     service.ClusterID,
		Host:          baseDomain,
		Path:          pathMatch,
		RewriteRegex:  rewriteRegex,
		RewriteTarget: "/$1",
		ServiceName:   service.InstanceName,
		ServicePort:   80,
		Labels: map[string]string{
			consts.LabelInstanceType: consts.InferenceInstance,
			"neptune.io/instance":    service.InstanceName,
		},
		Websocket:  false,
		EnableAuth: true,
		AuthUri:    authUri,
	}

	err := s.apisixSvc.CreateRoute(ctx, routeReq)
	if err != nil && apierrors.IsAlreadyExists(err) {
		// 路由已存在，视为成功
		return nil
	}
	return err
}

// deleteInferenceRoute 删除推理服务的 APISIX 路由
func (s *InferenceServiceService) deleteInferenceRoute(ctx context.Context, service *inferenceModel.Inference) {
	if s.apisixSvc == nil || !global.GVA_CONFIG.Apisix.Enabled {
		return
	}

	if err := s.apisixSvc.DeleteRoute(ctx, &apisixReq.DeleteRouteReq{
		Name:      fmt.Sprintf("inference-%s", service.InstanceName),
		Namespace: service.Namespace,
		ClusterId: service.ClusterID,
	}); err != nil {
		logx.Error("删除推理服务APISIX路由失败", err)
	}
}

// buildGatewayUrl 构建推理服务的 APISIX 网关完整访问地址
func buildGatewayUrl(namespace, instanceName string) string {
	return fmt.Sprintf("/inference/%s/%s", namespace, instanceName)
}

// updateStatus 更新服务状态
func (s *InferenceServiceService) updateStatus(id uint, status, errMsg string) {
	updates := map[string]interface{}{"status": status}
	if errMsg != "" {
		updates["error_msg"] = errMsg
	}
	global.GVA_DB.Model(&inferenceModel.Inference{}).Where("id = ?", id).Updates(updates)
}

// syncPodGroupStatus 通过 podgroups 表同步推理服务状态
//
// 所有推理服务（单机/分布式）都由 Volcano 调度，PodGroup Informer 持续更新 podgroups 表。
// 此方法仅做"终态补偿"：如果 podgroups 表中已无对应记录（PodGroup 已被删除），
// 且当前状态不是终态，则标记为 STOPPED。
// 中间状态（Pending → Running）完全信任 Informer，不做覆盖，避免状态冲突。
func (s *InferenceServiceService) syncPodGroupStatus(service *inferenceModel.Inference) {
	// 终态不需要同步
	if service.Status == consts.InferenceStatusStopped ||
		service.Status == consts.InferenceStatusFailed ||
		service.Status == consts.InferenceStatusDeleting {
		return
	}

	// 查 podgroups 表：Informer 已经把 PodGroup 状态写入了这张表
	var pg struct {
		Phase string
	}
	err := global.GVA_DB.Table("podgroups").
		Select("phase").
		Where("instance_name = ? AND instance_type = ? AND deleted_at IS NULL", service.InstanceName, consts.InferenceInstance).
		Order("created_at DESC").
		Limit(1).
		Scan(&pg).Error

	if err != nil || pg.Phase == "" {
		// podgroups 表中无记录：PodGroup 已被清理，标记为 STOPPED
		s.updateStatus(service.ID, consts.InferenceStatusStopped, "")
		service.Status = consts.InferenceStatusStopped
		return
	}

	// 将 PodGroup phase 映射为业务状态（与 Informer 的 businessConfigs 一致）
	statusMap := map[string]string{
		"Pending": consts.InferenceStatusPending,
		"Inqueue": consts.InferenceStatusPending,
		"Running": consts.InferenceStatusRunning,
		"Unknown": consts.InferenceStatusFailed,
	}

	newStatus, ok := statusMap[pg.Phase]
	if !ok {
		return
	}

	if newStatus != service.Status {
		s.updateStatus(service.ID, newStatus, "")
		service.Status = newStatus
	}
}

// DeleteInferenceService 删除推理服务
func (s *InferenceServiceService) DeleteInferenceService(ctx context.Context, req *inferenceReq.DeleteInferenceServiceReq) error {
	service, err := s.getInferenceService(req.ID)
	if err != nil {
		return err
	}

	cluster := global.GVA_K8S_CLUSTER_MANAGER.GetCluster(service.ClusterID)
	if cluster != nil {
		// 删除计算资源
		switch service.DeployType {
		case consts.DeployTypeDistributed:
			if cluster.VolcanoClient != nil {
				_ = cluster.VolcanoClient.Clientset().BatchV1alpha1().Jobs(service.Namespace).Delete(ctx, service.InstanceName, metav1.DeleteOptions{
					PropagationPolicy: func() *metav1.DeletionPropagation { p := metav1.DeletePropagationBackground; return &p }(),
				})
			}
		case consts.DeployTypeStandalone:
			_ = cluster.ClientSet.AppsV1().Deployments(service.Namespace).Delete(ctx, service.InstanceName, metav1.DeleteOptions{})
		}
		// 删除 K8s Service
		_ = cluster.ClientSet.CoreV1().Services(service.Namespace).Delete(ctx, service.InstanceName, metav1.DeleteOptions{})
	}

	// 配额释放和订单停止由 PodGroup Informer 监听 PodGroup 删除事件统一处理
	// （Deployment 删除 → Pod 删除 → Volcano 自动删除 PodGroup → Informer 捕获）

	// 删除 APISIX 路由
	s.deleteInferenceRoute(ctx, service)

	// 删除关联数据
	global.GVA_DB.Where("service_id = ?", service.ID).Delete(&inferenceModel.InferenceApiKey{})
	global.GVA_DB.Where("service_id = ?", service.ID).Delete(&inferenceModel.InferenceMount{})
	global.GVA_DB.Where("service_id = ?", service.ID).Delete(&inferenceModel.InferenceEnv{})
	global.GVA_DB.Where("service_id = ?", service.ID).Delete(&inferenceModel.InferenceApiKeyPolicy{})

	return global.GVA_DB.Delete(&inferenceModel.Inference{}, service.ID).Error
}

// StopInferenceService 停止推理服务
func (s *InferenceServiceService) StopInferenceService(ctx context.Context, req *inferenceReq.StopInferenceServiceReq) error {
	service, err := s.getInferenceService(req.ID)
	if err != nil {
		return err
	}

	if service.Status == consts.InferenceStatusStopped {
		return errors.New("服务已处于停止状态")
	}

	cluster := global.GVA_K8S_CLUSTER_MANAGER.GetCluster(service.ClusterID)
	if cluster == nil {
		return errors.New("集群不存在")
	}

	switch service.DeployType {
	case consts.DeployTypeDistributed:
		// 分布式模式：删除 VCJob（无法缩容）
		if cluster.VolcanoClient != nil {
			if err := cluster.VolcanoClient.Clientset().BatchV1alpha1().Jobs(service.Namespace).Delete(
				ctx, service.InstanceName, metav1.DeleteOptions{
					PropagationPolicy: func() *metav1.DeletionPropagation { p := metav1.DeletePropagationBackground; return &p }(),
				},
			); err != nil && !apierrors.IsNotFound(err) {
				return errors.Wrap(err, "删除VCJob失败")
			}
		}
	case consts.DeployTypeStandalone:
		// 单体模式：缩容到 0（Pod 删除后 Volcano 自动删除 PodGroup，Informer 捕获并释放资源）
		replicas := int32(0)
		deploy, err := cluster.ClientSet.AppsV1().Deployments(service.Namespace).Get(ctx, service.InstanceName, metav1.GetOptions{})
		if err != nil {
			return errors.Wrap(err, "获取Deployment失败")
		}
		deploy.Spec.Replicas = &replicas
		if _, err = cluster.ClientSet.AppsV1().Deployments(service.Namespace).Update(ctx, deploy, metav1.UpdateOptions{}); err != nil {
			return errors.Wrap(err, "缩容Deployment失败")
		}
	}

	s.updateStatus(service.ID, consts.InferenceStatusStopped, "")

	// 删除 APISIX 路由（与 Notebook 一致：停止时清理路由，启动时重建）
	s.deleteInferenceRoute(ctx, service)

	return nil
}

// StartInferenceService 启动推理服务
func (s *InferenceServiceService) StartInferenceService(ctx context.Context, req *inferenceReq.StartInferenceServiceReq) (err error) {
	service, err := s.getInferenceService(req.ID)
	if err != nil {
		return err
	}

	if service.Status != consts.InferenceStatusStopped && service.Status != consts.InferenceStatusFailed {
		return errors.New("服务未处于停止/失败状态，无法启动")
	}

	cleanups := make(Cleanups, 0)
	defer func() {
		if err != nil {
			cleanups.Run(context.Background())
		}
	}()

	plan, err := (&productService.ProductService{}).PlanAllocations(ctx, service.ProductId, int64(service.InstanceCount), service.ScheduleStrategy)
	if err != nil {
		return errors.Wrap(err, "启动服务规划资源失败")
	}

	if err = s.reserveInferenceCapacity(ctx, plan, &cleanups); err != nil {
		return errors.Wrap(err, "启动服务锁定资源失败")
	}

	if err = s.saveInferenceAllocations(ctx, service, plan, &cleanups); err != nil {
		return err
	}

	cluster, err := s.getInferenceCluster(service.ClusterID)
	if err != nil {
		return err
	}

	if err = s.createInferenceOrder(ctx, s.createOrderSpecFromService(service, cluster), &cleanups); err != nil {
		return err
	}

	switch service.DeployType {
	case consts.DeployTypeDistributed:
		// 分布式模式：重新创建 VCJob
		if cluster.VolcanoClient == nil {
			return errors.New("集群未配置 Volcano")
		}

		image, imageErr := s.getInferenceImage(service.ImageId)
		if imageErr != nil {
			return imageErr
		}
		modelPVC, volumeErr := s.getInferenceVolume(service.ModelPvcId, "获取模型PVC信息失败")
		if volumeErr != nil {
			return volumeErr
		}

		// 先尝试删除残留的旧 Job
		_ = cluster.VolcanoClient.Clientset().BatchV1alpha1().Jobs(service.Namespace).Delete(
			ctx, service.InstanceName, metav1.DeleteOptions{},
		)

		if err = s.createVCJob(ctx, cluster, service, image, modelPVC.PVCName); err != nil {
			return errors.Wrap(err, "重新创建VCJob失败")
		}

	case consts.DeployTypeStandalone:
		// 单体模式：扩容到 1（Volcano 调度器会自动为新 Pod 创建 PodGroup）
		replicas := int32(1)
		deploy, getErr := cluster.ClientSet.AppsV1().Deployments(service.Namespace).Get(ctx, service.InstanceName, metav1.GetOptions{})
		if getErr != nil {
			return errors.Wrap(getErr, "获取Deployment失败")
		}
		deploy.Spec.Replicas = &replicas
		if _, err = cluster.ClientSet.AppsV1().Deployments(service.Namespace).Update(ctx, deploy, metav1.UpdateOptions{}); err != nil {
			return errors.Wrap(err, "扩容Deployment失败")
		}
	}

	// 确保 APISIX 路由存在（首次创建可能失败，或被意外删除）
	if routeErr := s.ensureInferenceRoute(ctx, service); routeErr != nil {
		logx.Error("启动时确保APISIX路由失败", routeErr)
	}

	// 重置重启计数（如果是从 FAILED 恢复）
	updates := map[string]interface{}{
		"status":        consts.InferenceStatusPending,
		"error_msg":     "",
		"restart_count": 0,
	}
	global.GVA_DB.Model(&inferenceModel.Inference{}).Where("id = ?", service.ID).Updates(updates)
	return nil
}

// GetInferenceServiceList 获取推理服务列表
func (s *InferenceServiceService) GetInferenceServiceList(ctx context.Context, req *inferenceReq.GetInferenceServiceListReq) (*response.GetInferenceServiceListResp, error) {
	if req.Page <= 0 {
		req.Page = 1
	}
	if req.PageSize <= 0 {
		req.PageSize = 10
	}

	db := global.GVA_DB.Model(&inferenceModel.Inference{})

	if req.UserId > 0 {
		db = db.Where("user_id = ?", req.UserId)
	}
	if req.Name != "" {
		db = db.Where("display_name LIKE ?", "%"+req.Name+"%")
	}
	if req.Status != "" {
		db = db.Where("status = ?", req.Status)
	}
	if req.DeployType != "" {
		db = db.Where("deploy_type = ?", req.DeployType)
	}
	if req.Framework != "" {
		db = db.Where("framework = ?", req.Framework)
	}

	var total int64
	if err := db.Count(&total).Error; err != nil {
		return nil, err
	}

	var services []inferenceModel.Inference
	if err := db.Order("created_at DESC").Offset((req.Page - 1) * req.PageSize).Limit(req.PageSize).Find(&services).Error; err != nil {
		return nil, err
	}

	list := make([]response.InferenceServiceListItem, len(services))
	for i, svc := range services {
		list[i] = response.InferenceServiceListItem{
			ID:           svc.ID,
			DisplayName:  svc.DisplayName,
			InstanceName: svc.InstanceName,
			DeployType:   svc.DeployType,
			Framework:    svc.Framework,
			Status:       svc.Status,
			GPU:          svc.GPU,
			GPUModel:     svc.GPUModel,
			CreatedAt:    svc.CreatedAt,
			StartedAt:    svc.StartedAt,
		}
	}

	return &response.GetInferenceServiceListResp{
		Total: total,
		List:  list,
	}, nil
}

// GetInferenceServiceDetail 获取推理服务详情
func (s *InferenceServiceService) GetInferenceServiceDetail(ctx context.Context, req *inferenceReq.GetInferenceServiceDetailReq) (*response.InferenceServiceDetail, error) {
	service, err := s.getInferenceService(req.ID)
	if err != nil {
		return nil, err
	}

	// 通过 PodGroup 同步状态（统一的 Volcano 状态来源）
	s.syncPodGroupStatus(service)

	// 获取镜像名称
	var image imageModel.Image
	if err = global.GVA_DB.Where("id = ?", service.ImageId).First(&image).Error; err != nil {
		logx.Error("获取镜像信息失败", err)
	}

	// 获取挂载配置
	var mounts []inferenceModel.InferenceMount
	if err = global.GVA_DB.Where("service_id = ?", service.ID).Find(&mounts).Error; err != nil {
		return nil, errors.Wrap(err, "获取挂载配置失败")
	}

	mountItems := make([]response.InferenceMountItem, len(mounts))
	for i, m := range mounts {
		mountItems[i] = response.InferenceMountItem{
			MountType: m.MountType,
			PvcName:   m.PvcName,
			SubPath:   m.SubPath,
			MountPath: m.MountPath,
			ReadOnly:  m.ReadOnly,
		}
	}

	// 获取环境变量
	var envs []inferenceModel.InferenceEnv
	if err = global.GVA_DB.Where("service_id = ?", service.ID).Find(&envs).Error; err != nil {
		return nil, errors.Wrap(err, "获取环境变量失败")
	}

	envItems := make([]response.InferenceEnvItem, len(envs))
	for i, e := range envs {
		envItems[i] = response.InferenceEnvItem{
			Name:  e.Name,
			Value: e.Value,
		}
	}

	// Command 直接是原始字符串，Args 需要 JSON 反序列化
	command := service.Command
	var args []string
	if service.Args != "" {
		json.Unmarshal([]byte(service.Args), &args)
	}

	return &response.InferenceServiceDetail{
		ID:               service.ID,
		DisplayName:      service.DisplayName,
		InstanceName:     service.InstanceName,
		Namespace:        service.Namespace,
		DeployType:       service.DeployType,
		Framework:        service.Framework,
		ModelMountPath:   service.ModelMountPath,
		ImageName:        image.Name,
		TensorParallel:   service.TensorParallel,
		PipelineParallel: service.PipelineParallel,
		WorkerCount:      service.WorkerCount,
		CPU:              service.CPU,
		Memory:           service.Memory,
		GPU:              service.GPU,
		GPUModel:         service.GPUModel,
		ServicePort:      service.ServicePort,
		Command:          command,
		Args:             args,
		AutoRestart:      service.AutoRestart,
		RestartCount:     service.RestartCount,
		MaxRestarts:      service.MaxRestarts,
		AuthType:         service.AuthType,
		Status:           service.Status,
		ErrorMsg:         service.ErrorMsg,
		AccessUrl:        buildGatewayUrl(service.Namespace, service.InstanceName),
		CreatedAt:        service.CreatedAt,
		StartedAt:        service.StartedAt,
		Mounts:           mountItems,
		Envs:             envItems,
	}, nil
}

// GetInferenceServiceLogs 获取推理服务日志
func (s *InferenceServiceService) GetInferenceServiceLogs(ctx context.Context, req *inferenceReq.GetInferenceServiceLogsReq) (io.ReadCloser, error) {
	service, err := s.getInferenceService(req.ID)
	if err != nil {
		return nil, err
	}

	cluster, err := s.getInferenceCluster(service.ClusterID)
	if err != nil {
		return nil, err
	}

	var podName string
	if req.PodName != "" {
		// 直接指定 Pod 名称
		podName = req.PodName
	} else {
		switch service.DeployType {
		case consts.DeployTypeDistributed:
			// VCJob Pod 命名规则: {jobName}-{taskName}-{index}
			taskName := req.TaskName
			if taskName == "" {
				taskName = "head"
			}
			podIndex := 0
			if req.PodIndex != nil {
				podIndex = *req.PodIndex
			}
			podName = fmt.Sprintf("%s-%s-%d", service.InstanceName, taskName, podIndex)
		case consts.DeployTypeStandalone:
			// Deployment: 通过 label 查找
			pods, err := cluster.ClientSet.CoreV1().Pods(service.Namespace).List(ctx, metav1.ListOptions{
				LabelSelector: fmt.Sprintf("%s=%s,neptune.io/instance=%s", consts.LabelApp, service.InstanceName, service.InstanceName),
			})
			if err != nil || len(pods.Items) == 0 {
				return nil, errors.New("未找到运行中的Pod")
			}
			podName = pods.Items[0].Name
		}
	}

	// 容器名称与 builder 中定义的一致
	container := req.Container
	if container == "" {
		if service.DeployType == consts.DeployTypeDistributed {
			container = service.InstanceName + "-head"
		} else {
			container = service.InstanceName
		}
	}

	tailLines := req.TailLines
	if tailLines == 0 {
		tailLines = 100
	}

	return cluster.ClientSet.CoreV1().Pods(service.Namespace).GetLogs(podName, &corev1.PodLogOptions{
		Container: container,
		TailLines: &tailLines,
		Follow:    req.Follow,
	}).Stream(ctx)
}

// GetInferenceServicePods 获取推理服务 Pod 列表
func (s *InferenceServiceService) GetInferenceServicePods(ctx context.Context, id uint) ([]terminalService.PodInfo, error) {
	service, err := s.getInferenceService(id)
	if err != nil {
		return nil, err
	}

	cluster, err := s.getInferenceCluster(service.ClusterID)
	if err != nil {
		return nil, err
	}

	pods, err := cluster.ClientSet.CoreV1().Pods(service.Namespace).List(ctx, metav1.ListOptions{
		LabelSelector: fmt.Sprintf("%s=%s,neptune.io/instance=%s", consts.LabelApp, service.InstanceName, service.InstanceName),
	})
	if err != nil {
		return nil, err
	}

	result := make([]terminalService.PodInfo, len(pods.Items))
	for i, pod := range pods.Items {
		result[i] = terminalService.PodInfo{
			Name:      pod.Name,
			Namespace: pod.Namespace,
			Status:    string(pod.Status.Phase),
		}
	}

	return result, nil
}

// TerminalInfo 终端连接所需的服务信息
type TerminalInfo struct {
	InstanceName string
	Namespace    string
	ClusterID    uint
	DeployType   string
}

// GetInferenceServiceForTerminal 获取推理服务终端连接信息
func (s *InferenceServiceService) GetInferenceServiceForTerminal(ctx context.Context, id uint) (*TerminalInfo, error) {
	svc, err := s.getInferenceService(id)
	if err != nil {
		return nil, err
	}
	return &TerminalInfo{
		InstanceName: svc.InstanceName,
		Namespace:    svc.Namespace,
		ClusterID:    svc.ClusterID,
		DeployType:   svc.DeployType,
	}, nil
}
