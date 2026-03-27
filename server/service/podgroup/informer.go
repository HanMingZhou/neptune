package podgroup

import (
	"context"
	"fmt"
	"gin-vue-admin/model/consts"
	podgroupModel "gin-vue-admin/model/podgroup"
	"strings"
	"time"

	"github.com/pkg/errors"
	"github.com/zeromicro/go-zero/core/logx"
	"gorm.io/gorm"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/util/runtime"
	"k8s.io/apimachinery/pkg/util/wait"
	"k8s.io/client-go/tools/cache"
	"k8s.io/client-go/util/workqueue"
	vcv1beta1 "volcano.sh/apis/pkg/apis/scheduling/v1beta1"
	vcinformers "volcano.sh/apis/pkg/client/informers/externalversions/scheduling/v1beta1"
	vclisters "volcano.sh/apis/pkg/client/listers/scheduling/v1beta1"
)

// 事件类型常量
const (
	eventTypeAdd    = "add"
	eventTypeUpdate = "update"
	eventTypeDelete = "delete"
)

// 资源标签常量
const (
	labelResourceCPU      = "resource.cpu"
	labelResourceMemory   = "resource.memory"
	labelResourceGPU      = "resource.gpu"
	labelResourceGPUModel = "resource.gpu-model"
)

// Table 名称常量
const (
	tableNameTrainingJobs      = "training_jobs"
	tableNameNotebooks         = "notebooks"
	tableNameInferenceServices = "inferences"
)

// 业务资源配置映射
var businessConfigs = map[string]struct {
	tableName string
	orderType int
	statusMap map[string]string
}{
	consts.TrainingInstance: {
		tableName: tableNameTrainingJobs,
		orderType: consts.OrderTypeTraining,
		statusMap: map[string]string{
			consts.VolcanoPhasePending:   consts.TrainingStatusPending,
			consts.VolcanoPhaseInqueue:   consts.TrainingStatusPending,
			consts.VolcanoPhaseRunning:   consts.TrainingStatusRunning,
			consts.VolcanoPhaseSucceeded: consts.TrainingStatusSucceeded,
			consts.VolcanoPhaseFailed:    consts.TrainingStatusFailed,
			consts.VolcanoPhaseUnknown:   consts.TrainingStatusFailed,
		},
	},
	consts.NotebookInstance: {
		tableName: tableNameNotebooks,
		orderType: consts.OrderTypeNotebook,
		statusMap: map[string]string{
			consts.VolcanoPhasePending:   consts.NotebookStatusPending,
			consts.VolcanoPhaseInqueue:   consts.NotebookStatusPending,
			consts.VolcanoPhaseRunning:   consts.NotebookStatusRunning,
			consts.VolcanoPhaseSucceeded: consts.NotebookStatusStopped,
			consts.VolcanoPhaseFailed:    consts.NotebookStatusFailed,
			consts.VolcanoPhaseUnknown:   consts.NotebookStatusFailed,
		},
	},
	consts.InferenceInstance: {
		tableName: tableNameInferenceServices,
		orderType: consts.OrderTypeInference,
		statusMap: map[string]string{
			consts.VolcanoPhasePending:   consts.InferenceStatusPending,
			consts.VolcanoPhaseInqueue:   consts.InferenceStatusPending,
			consts.VolcanoPhaseRunning:   consts.InferenceStatusRunning,
			consts.VolcanoPhaseSucceeded: consts.InferenceStatusStopped,
			consts.VolcanoPhaseFailed:    consts.InferenceStatusFailed,
			consts.VolcanoPhaseUnknown:   consts.InferenceStatusFailed,
		},
	},
}

// OrderManager 计费接口
type OrderManager interface {
	StartOrder(ctx context.Context, resourceId int64, resourceType int) error
	StopOrder(ctx context.Context, resourceId int64, resourceType int) error
}

// ProductManager 产品资源管理接口
type ProductManager interface {
	ReleaseCapacity(ctx context.Context, productId uint, count int64) error
	ReleaseCapacityAuto(ctx context.Context, productId uint) error
}

type podGroupEvent struct {
	eventType string
	podGroup  *vcv1beta1.PodGroup
	key       string
}

// resourceIdentity identifyResource 的返回结果
type resourceIdentity struct {
	InstanceType string // 资源类型: notebook / training / inference
	OwnerID      uint   // 业务表主键 ID
	InstanceName string // K8s 实例名称 (instance_name / job_name)
}

// Found 是否成功识别到资源
func (r *resourceIdentity) Found() bool {
	return r.OwnerID > 0 && r.InstanceType != ""
}

// PodGroupInformerFactory 核心工厂
type PodGroupInformerFactory struct {
	informer       vcinformers.PodGroupInformer
	lister         vclisters.PodGroupLister
	synced         cache.InformerSynced
	queue          workqueue.TypedRateLimitingInterface[podGroupEvent]
	stopCh         chan struct{}
	clusterID      uint
	db             *gorm.DB
	orderManager   OrderManager
	productManager ProductManager
}

func NewPodGroupInformerFactory(
	informer vcinformers.PodGroupInformer,
	clusterID uint,
	db *gorm.DB,
	orderManager OrderManager,
) *PodGroupInformerFactory {
	return &PodGroupInformerFactory{
		informer:       informer,
		lister:         informer.Lister(),
		synced:         informer.Informer().HasSynced,
		queue:          workqueue.NewTypedRateLimitingQueue(workqueue.DefaultTypedControllerRateLimiter[podGroupEvent]()),
		stopCh:         make(chan struct{}),
		clusterID:      clusterID,
		db:             db,
		orderManager:   orderManager,
		productManager: nil, // 延迟注入
	}
}

// SetOrderManager 设置计费管理器（用于延迟注入）
func (f *PodGroupInformerFactory) SetOrderManager(bm OrderManager) {
	f.orderManager = bm
}

// SetProductManager 设置产品管理器（用于延迟注入）
func (f *PodGroupInformerFactory) SetProductManager(pm ProductManager) {
	f.productManager = pm
}

func (f *PodGroupInformerFactory) Start(ctx context.Context) error {
	f.informer.Informer().AddEventHandler(cache.ResourceEventHandlerFuncs{
		AddFunc: func(obj any) {
			pg := obj.(*vcv1beta1.PodGroup)
			f.queue.Add(podGroupEvent{eventType: eventTypeAdd, podGroup: pg})
		},
		UpdateFunc: func(old, new any) {
			pg := new.(*vcv1beta1.PodGroup)
			f.queue.Add(podGroupEvent{eventType: eventTypeUpdate, podGroup: pg})
		},
		DeleteFunc: func(obj any) {
			pg, ok := obj.(*vcv1beta1.PodGroup)
			if !ok {
				tombstone, ok := obj.(cache.DeletedFinalStateUnknown)
				if !ok {
					runtime.HandleError(fmt.Errorf("error decoding object, invalid type"))
					return
				}
				pg, ok = tombstone.Obj.(*vcv1beta1.PodGroup)
				if !ok {
					runtime.HandleError(fmt.Errorf("error decoding object tombstone, invalid type"))
					return
				}
			}
			key, _ := cache.MetaNamespaceKeyFunc(pg)
			f.queue.Add(podGroupEvent{eventType: eventTypeDelete, podGroup: pg, key: key})
		},
	})

	go func() {
		logx.Info("等待 PodGroup Informer 缓存同步...")
		if !cache.WaitForCacheSync(f.stopCh, f.synced) {
			select {
			case <-f.stopCh:
				logx.Info("PodGroup Informer 已停止,clusterId", f.clusterID)
			default:
				logx.Error("PodGroup Informer 缓存同步失败,clusterId", f.clusterID)
			}
			return
		}
		logx.Info("PodGroup Informer 缓存同步完成,clusterId", f.clusterID)
		wait.Until(f.runWorker, time.Second, f.stopCh)
	}()

	return nil
}

func (f *PodGroupInformerFactory) Stop() {
	logx.Info("停止 PodGroup Controller,clusterId", f.clusterID)
	close(f.stopCh)
	f.queue.ShutDown()
}

func (f *PodGroupInformerFactory) runWorker() {
	for f.processNextItem() {
	}
}

func (f *PodGroupInformerFactory) processNextItem() bool {
	event, quit := f.queue.Get()
	if quit {
		return false
	}
	defer f.queue.Done(event)

	var err error
	switch event.eventType {
	case eventTypeAdd, eventTypeUpdate:
		err = f.processUpdate(event.podGroup)
	case eventTypeDelete:
		err = f.processDelete(event.podGroup)
	}

	if err == nil {
		f.queue.Forget(event)
		return true
	}

	if f.queue.NumRequeues(event) < 5 {
		logx.Error("处理 PodGroup 失败，重新入队",
			logx.Field("name", event.podGroup.Name),
			logx.Field("namespace", event.podGroup.Namespace),
			logx.Field("error", err))
		f.queue.AddRateLimited(event)
		return true
	}

	f.queue.Forget(event)
	logx.Error("处理 PodGroup 失败，放弃重试",
		logx.Field("name", event.podGroup.Name),
		logx.Field("namespace", event.podGroup.Namespace),
		logx.Field("error", err))
	return true
}

func (f *PodGroupInformerFactory) processUpdate(pg *vcv1beta1.PodGroup) error {
	dbPG := f.convertToDBModel(pg)

	var existing podgroupModel.PodGroup
	err := f.db.Where("name = ? AND namespace = ?", pg.Name, pg.Namespace).First(&existing).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return f.db.Create(dbPG).Error
		}
		return err
	}

	updates := map[string]any{
		"phase":            dbPG.Phase,
		"status":           dbPG.Status,
		"running":          dbPG.Running,
		"succeeded":        dbPG.Succeeded,
		"failed":           dbPG.Failed,
		"resource_version": dbPG.ResourceVersion,
		"generation":       dbPG.Generation,
	}

	if dbPG.Status == consts.PodGroupStatusRunning && existing.StartTime == nil {
		now := time.Now()
		updates["start_time"] = &now
	}

	if err := f.db.Model(&existing).Updates(updates).Error; err != nil {
		return err
	}

	f.updateBusinessStatus(pg)
	return nil
}

func (f *PodGroupInformerFactory) updateBusinessStatus(pg *vcv1beta1.PodGroup) {
	res := f.identifyResource(pg)
	config, ok := businessConfigs[res.InstanceType]
	if !ok || !res.Found() {
		return
	}

	status, ok := config.statusMap[string(pg.Status.Phase)]
	if !ok {
		return
	}

	updates := map[string]any{"status": status}
	now := time.Now()

	if res.InstanceType == consts.TrainingInstance {
		switch status {
		case consts.TrainingStatusRunning:
			f.db.Table(config.tableName).
				Where("id = ? AND started_at IS NULL", res.OwnerID).
				Update("started_at", &now)
		case consts.TrainingStatusSucceeded, consts.TrainingStatusFailed:
			updates["finished_at"] = &now
		}
	}

	// 资源进入运行状态时，触发按量计费
	if pg.Status.Phase == vcv1beta1.PodGroupRunning && config.orderType > 0 {
		if err := f.orderManager.StartOrder(context.Background(), int64(res.OwnerID), config.orderType); err != nil {
			logx.Error("开启按量计费失败",
				logx.Field("err", err),
				logx.Field("ownerID", res.OwnerID))
		}
	}

	f.db.Table(config.tableName).Where("id = ? AND status != ?", res.OwnerID, status).Updates(updates)
}

// identifyResource 从 PodGroup 的 Label 或 OwnerReferences 中识别关联的业务资源
// 返回 resourceIdentity，调用方通过 .Found() 判断是否识别成功
func (f *PodGroupInformerFactory) identifyResource(pg *vcv1beta1.PodGroup) resourceIdentity {
	// 优先使用 Label 中的显式标识（由平台创建资源时注入）
	if labelType := pg.Labels[consts.LabelInstanceType]; labelType != "" {
		if labelID := pg.Labels[consts.LabelJobID]; labelID != "" {
			id := parseUint(labelID)
			if id > 0 {
				return resourceIdentity{
					InstanceType: labelType,
					OwnerID:      id,
					InstanceName: pg.Labels[consts.LabelApp],
				}
			}
		}
	}

	// Label 不完整时，通过 OwnerReferences 或 PodGroup 名称推断
	var ownerName string
	if len(pg.OwnerReferences) > 0 {
		ownerName = pg.OwnerReferences[0].Name
	} else {
		ownerName = pg.Name
	}
	if ownerName == "" {
		return resourceIdentity{}
	}

	// instance name 命名规则: {type}-{6位hex}，例如 notebook-1112c5、training-8a86dd、inference-336ae5
	// 单机推理/Notebook 的 ReplicaSet 名称为 {instance_name}-{hash}，需要去掉后缀
	// 通过前缀判断资源类型，避免逐表盲查
	// 注意：使用 Limit(1).Find() 而非 First()，避免 record not found 被 GORM logger 打成 error
	switch {
	case strings.HasPrefix(ownerName, consts.NotebookInstance+"-"):
		return f.lookupByInstanceName(tableNameNotebooks, "instance_name", ownerName, consts.NotebookInstance)

	case strings.HasPrefix(ownerName, consts.TrainingInstance+"-"):
		return f.lookupByInstanceName(tableNameTrainingJobs, "job_name", ownerName, consts.TrainingInstance)

	case strings.HasPrefix(ownerName, consts.InferenceInstance+"-"):
		return f.lookupByInstanceName(tableNameInferenceServices, "instance_name", ownerName, consts.InferenceInstance)
	}

	return resourceIdentity{}
}

// lookupByInstanceName 根据实例名称查询业务表，支持去掉 ReplicaSet hash 后缀重试
func (f *PodGroupInformerFactory) lookupByInstanceName(table, column, ownerName, instanceType string) resourceIdentity {
	var id uint

	// 先精确匹配
	f.db.Table(table).Where(column+" = ?", ownerName).Select("id").Limit(1).Scan(&id)
	if id > 0 {
		return resourceIdentity{InstanceType: instanceType, OwnerID: id, InstanceName: ownerName}
	}

	// 去掉 ReplicaSet hash 后缀重试: {instance_name}-{hash} → {instance_name}
	// 前缀长度 = len("notebook-") 等，确保不会截到类型前缀本身
	if idx := strings.LastIndex(ownerName, "-"); idx > len(instanceType) {
		candidate := ownerName[:idx]
		f.db.Table(table).Where(column+" = ?", candidate).Select("id").Limit(1).Scan(&id)
		if id > 0 {
			return resourceIdentity{InstanceType: instanceType, OwnerID: id, InstanceName: candidate}
		}
	}

	return resourceIdentity{}
}

func (f *PodGroupInformerFactory) processDelete(pg *vcv1beta1.PodGroup) error {
	res := f.identifyResource(pg)
	config, ok := businessConfigs[res.InstanceType]

	// 1. 停止计费
	if ok && res.Found() && config.orderType > 0 {
		if stopErr := f.orderManager.StopOrder(context.Background(), int64(res.OwnerID), config.orderType); stopErr != nil {
			logx.Error("Informer StopOrder 失败",
				logx.Field("instanceType", res.InstanceType),
				logx.Field("ownerID", res.OwnerID),
				logx.Field("orderType", config.orderType),
				logx.Field("err", stopErr))
		} else {
			logx.Info("Informer StopOrder 成功",
				logx.Field("instanceType", res.InstanceType),
				logx.Field("ownerID", res.OwnerID))
		}
	}

	// 2. 释放资源配额
	if ok && res.Found() && f.productManager != nil {
		f.releaseProductCapacity(res)
	}

	// 3. 更新 PodGroup 状态
	now := time.Now()
	return f.db.Model(&podgroupModel.PodGroup{}).
		Where("name = ? AND namespace = ?", pg.Name, pg.Namespace).
		Updates(map[string]any{
			"status":   consts.PodGroupStatusDeleted,
			"end_time": &now,
		}).Error
}

// releaseProductCapacity 根据业务资源信息释放产品配额
func (f *PodGroupInformerFactory) releaseProductCapacity(res resourceIdentity) {
	var result podgroupModel.ResourceResult
	config := businessConfigs[res.InstanceType]

	switch res.InstanceType {
	case consts.NotebookInstance:
		f.db.Table(config.tableName).Where("id = ?", res.OwnerID).Select("product_id").Scan(&result)
	case consts.TrainingInstance:
		f.db.Table(config.tableName).Where("id = ?", res.OwnerID).Select("product_id, worker_count, framework_type").Scan(&result)
	case consts.InferenceInstance:
		f.db.Table(config.tableName).Where("id = ?", res.OwnerID).Select("product_id, worker_count, deploy_type").Scan(&result)
	}

	if result.ProductId == 0 {
		return
	}

	releaseCount := f.calculateReleaseCount(res.InstanceType, result.WorkerCount, result.FrameworkType, result.DeployType)
	if releaseCount <= 0 {
		return
	}

	if err := f.productManager.ReleaseCapacity(context.Background(), result.ProductId, releaseCount); err != nil {
		logx.Error("Informer 释放资源失败",
			logx.Field("instanceType", res.InstanceType),
			logx.Field("ownerID", res.OwnerID),
			logx.Field("productId", result.ProductId),
			logx.Field("releaseCount", releaseCount),
			logx.Field("err", err))
	} else {
		logx.Info("Informer 成功释放资源",
			logx.Field("instanceType", res.InstanceType),
			logx.Field("ownerID", res.OwnerID),
			logx.Field("productId", result.ProductId),
			logx.Field("releaseCount", releaseCount))
	}
}

func (f *PodGroupInformerFactory) convertToDBModel(pg *vcv1beta1.PodGroup) *podgroupModel.PodGroup {
	res := f.identifyResource(pg)
	dbPG := &podgroupModel.PodGroup{
		Name:            pg.Name,
		Namespace:       pg.Namespace,
		InstanceName:    res.InstanceName,
		InstanceType:    res.InstanceType,
		OwnerID:         res.OwnerID,
		OwnerType:       res.InstanceType, // 当前业务中 OwnerType 始终等于 InstanceType
		ClusterID:       f.clusterID,
		MinMember:       pg.Spec.MinMember,
		Queue:           pg.Spec.Queue,
		Phase:           string(pg.Status.Phase),
		Status:          f.getStatusFromPhase(pg.Status.Phase),
		Running:         pg.Status.Running,
		Succeeded:       pg.Status.Succeeded,
		Failed:          pg.Status.Failed,
		ResourceVersion: pg.ResourceVersion,
		Generation:      pg.Generation,
	}

	if pg.Labels != nil {
		dbPG.CPU = pg.Labels[labelResourceCPU]
		dbPG.Memory = pg.Labels[labelResourceMemory]
		dbPG.GPU = pg.Labels[labelResourceGPU]
		dbPG.GPUModel = pg.Labels[labelResourceGPUModel]
	}
	return dbPG
}

func (f *PodGroupInformerFactory) getStatusFromPhase(phase vcv1beta1.PodGroupPhase) string {
	switch phase {
	case vcv1beta1.PodGroupPending, vcv1beta1.PodGroupInqueue:
		return consts.PodGroupStatusPending
	case vcv1beta1.PodGroupRunning:
		return consts.PodGroupStatusRunning
	default:
		return consts.PodGroupStatusFailed
	}
}

func (f *PodGroupInformerFactory) ListPodGroups(namespace string) ([]*vcv1beta1.PodGroup, error) {
	if namespace == "" {
		namespace = corev1.NamespaceAll
	}
	return f.lister.PodGroups(namespace).List(labels.Everything())
}

// calculateReleaseCount 计算需要释放的资源数量（实例数）
// 训练任务：根据框架类型和工作节点数计算（与创建时的 quotaCount 一致）
// Notebook：释放 1 个实例
// Inference：单机释放 1 个实例，分布式释放 WorkerCount 个实例
func (f *PodGroupInformerFactory) calculateReleaseCount(instanceType string, workerCount int, frameworkType, deployType string) int64 {
	switch instanceType {
	case consts.TrainingInstance:
		switch frameworkType {
		case consts.FrameworkPyTorchDDP:
			return int64(workerCount) // DDP: WorkerCount 包含 Master
		case consts.FrameworkMPI:
			return int64(workerCount) // MPI: Launcher 不消耗计算资源
		default:
			return 1 // Standalone
		}
	case consts.InferenceInstance:
		if deployType == consts.DeployTypeDistributed && workerCount > 1 {
			return int64(workerCount)
		}
		return 1
	default:
		return 1
	}
}

// parseUint 安全解析字符串为 uint，失败返回 0
func parseUint(s string) uint {
	var n uint
	for _, c := range s {
		if c < '0' || c > '9' {
			return 0
		}
		n = n*10 + uint(c-'0')
	}
	return n
}
