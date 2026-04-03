package notebook

import (
	"context"
	"fmt"
	"gin-vue-admin/global"
	apisixModel "gin-vue-admin/model/apisix"
	apisixReq "gin-vue-admin/model/apisix/request"
	"gin-vue-admin/model/consts"
	imageModel "gin-vue-admin/model/image"
	nbModel "gin-vue-admin/model/notebook"
	"gin-vue-admin/model/notebook/request"
	"gin-vue-admin/model/notebook/response"
	orderModel "gin-vue-admin/model/order"
	orderReq "gin-vue-admin/model/order/request"
	pipeReq "gin-vue-admin/model/pipe/request"
	productModel "gin-vue-admin/model/product"
	productRes "gin-vue-admin/model/product/response"
	pvcModelPkg "gin-vue-admin/model/pvc"
	pvcModel "gin-vue-admin/model/pvc/request"
	secretModelPkg "gin-vue-admin/model/secret"
	secretModel "gin-vue-admin/model/secret/request"
	systemModel "gin-vue-admin/model/system"
	tensorboardModel "gin-vue-admin/model/tensorboard"
	tensorboardReq "gin-vue-admin/model/tensorboard/request"
	orderService "gin-vue-admin/service/order"
	"gin-vue-admin/service/piper"
	"gin-vue-admin/service/podgroup"
	"gin-vue-admin/service/product"
	"gin-vue-admin/service/pvc"
	"gin-vue-admin/service/secret"
	"gin-vue-admin/service/tensorboard"
	helper "gin-vue-admin/utils/k8s"
	"gin-vue-admin/utils/password"
	"gin-vue-admin/utils/ssh"
	"gin-vue-admin/utils/validator"
	"io"
	"path"
	"strings"

	nbv1 "github.com/kubeflow/kubeflow/components/notebook-controller/api/v1"
	"github.com/pkg/errors"
	"github.com/zeromicro/go-zero/core/logx"
	"gorm.io/gorm"
	corev1 "k8s.io/api/core/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
	"k8s.io/client-go/kubernetes"
	v1lister "k8s.io/client-go/listers/apps/v1"
	v1podlister "k8s.io/client-go/listers/core/v1"
	"k8s.io/utils/pointer"
	ctrlclient "sigs.k8s.io/controller-runtime/pkg/client"
)

type NotebookManager interface {
	CreateNotebook(ctx context.Context, req *request.AddNoteBookReq) error
	DeleteNotebook(ctx context.Context, req *request.DeleteNoteBookReq) error
	UpdateNotebook(ctx context.Context, req *request.UpdateNoteBookReq) error
	StopNotebook(ctx context.Context, id uint) error
	StartNotebook(ctx context.Context, id uint) error
	GetNotebookList(ctx context.Context, req *request.GetNotebookListReq) (*response.GetNotebookListResp, error)
	GetNotebookDetail(ctx context.Context, req *request.GetNotebookDetailReq) (*response.NotebookItem, error)
	GetNotebookLogs(ctx context.Context, req *request.GetNotebookLogsReq) (io.ReadCloser, error)
	GetNotebookPods(ctx context.Context, req *request.GetNotebookPodsReq) ([]response.PodInfoResp, error)
	GetTerminalInfo(ctx context.Context, req *request.HandleTerminalReq) (*response.TerminalInfoResp, error)
}

var _ NotebookManager = &NotebookService{}

type notebookApisixService interface {
	CreateRoute(ctx context.Context, req *apisixReq.CreateRouteReq) error
	DeleteRoute(ctx context.Context, req *apisixReq.DeleteRouteReq) error
	CreateStreamRoute(ctx context.Context, req *apisixReq.CreateStreamRouteReq) error
}

type NotebookService struct {
	apisixSvc notebookApisixService
}

type Cleanups []func(ctx context.Context)

type notebookCreateState struct {
	req               *request.AddNoteBookReq
	userInfo          *systemModel.SysUser
	productInfo       productRes.ProductDetailResponse
	imageAddr         string
	sshPublicKey      string
	cluster           *global.ClusterClientInfo
	pvcManager        *pvc.K8sPVCService
	secretManager     *secret.K8sSecretService
	nbRef             *nbModel.Notebook
	privKeySecretName string
}

type notebookOrderSpec struct {
	userID      uint
	notebookID  uint
	productID   uint
	imageID     uint
	chargeType  int64
	displayName string
	remark      string
	area        string
	clusterID   uint
}

type notebookResponseLookup struct {
	imageNames map[uint]string
	products   map[uint]productModel.Product
}

var NotebookServiceApp = new(NotebookService)

func (c *Cleanups) Add(fn func(ctx context.Context)) {
	*c = append(*c, fn)
}

func (c Cleanups) Run(ctx context.Context) {
	for i := len(c) - 1; i >= 0; i-- {
		c[i](ctx)
	}
}

// SetApisixService 设置 Apisix 服务
func (nb *NotebookService) SetApisixService(svc notebookApisixService) {
	nb.apisixSvc = svc
}

// validateCreateRequest 校验创建请求参数
func (nb *NotebookService) validateCreateRequest(req *request.AddNoteBookReq) error {
	// 1. 校验名称格式
	if err := helper.ValidateNotebookName(req.DisplayName); err != nil {
		logx.Error("名称验证失败", err)
		return err
	}

	// 2. 校验挂载路径
	mountsPaths := make(map[string]bool)
	// 默认系统盘路径
	mountsPaths[nbModel.DefaultWorkspacePath] = true

	for _, mount := range req.VolumeMounts {
		// 如果用户没有填写，默认挂载到 /home/notebook/neptune-fs
		if mount.MountsPath == "" {
			mount.MountsPath = nbModel.DefaultDataMountPath
		}
		// 如果是相对路径（不以 / 开头），加上 workspace 前缀
		// 如果是绝对路径（以 / 开头），则保持原样，尊重用户的选择
		if !strings.HasPrefix(mount.MountsPath, "/") {
			mount.MountsPath = path.Join(nbModel.DefaultWorkspacePath, mount.MountsPath)
		}

		// 1. 调用通用路径验证器
		if err := validator.ValidateMountPath(mount.MountsPath); err != nil {
			logx.Error("挂载路径非法", err)
			return err
		}

		// 2. 重复校验
		if mountsPaths[mount.MountsPath] {
			logx.Error("挂载路径重复", mount.MountsPath)
			return errors.New("挂载路径重复")
		}
		mountsPaths[mount.MountsPath] = true
	}
	// 验证TensorBoard日志路径安全性（相对路径）
	if req.TensorBoard && req.TensorBoardLogPath != "" {
		if err := validator.ValidateSubPath(req.TensorBoardLogPath); err != nil {
			logx.Error("TensorBoard路径验证失败", err)
			return err
		}
	}

	return nil
}

// CreateNotebook 创建Notebook
func (nb *NotebookService) CreateNotebook(ctx context.Context, req *request.AddNoteBookReq) (err error) {
	// 检查参数
	if err = nb.validateCreateRequest(req); err != nil {
		logx.Error("参数验证失败", err)
		return err
	}

	// 资源清理器：收集需要回滚的操作，失败时统一执行
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

	var state *notebookCreateState
	if state, err = nb.buildCreateNotebookState(ctx, req, &cleanups); err != nil {
		return err
	}

	if err = nb.ensureNotebookNamespace(ctx, state.cluster.ClientSet, state.nbRef.Namespace); err != nil {
		return err
	}

	if err = nb.createNotebookPVCResources(ctx, state, &cleanups); err != nil {
		return err
	}

	if err = nb.createNotebookSecrets(ctx, state, &cleanups); err != nil {
		return err
	}

	// 4. 构建 Notebook 对象
	notebookObj := buildNotebook(state.nbRef, state.sshPublicKey)

	// 5. 写入数据库
	// 注意：GORM 会自动级联创建 VolumeMounts（通过外键关联），无需手动创建
	if err = global.GVA_DB.Create(state.nbRef).Error; err != nil {
		logx.Error("保存Notebook失败", err)
		return err
	}

	// 6. 创建 K8s Notebook
	if err = state.cluster.NotebookClient.Create(ctx, notebookObj); err != nil {
		logx.Error("K8s创建Notebook失败", err)
		return err
	}

	// 7. 创建订单（预付费在CreateOrder时已扣费，按量付费金额为0）
	orderID, orderErr := nb.createNotebookOrder(ctx, notebookOrderSpec{
		userID:      req.UserId,
		notebookID:  state.nbRef.ID,
		productID:   req.ProductId,
		imageID:     req.ImageId,
		chargeType:  req.PayType,
		displayName: state.nbRef.DisplayName,
		remark:      "创建容器实例",
		area:        state.cluster.Area,
		clusterID:   state.nbRef.ClusterID,
	}, &cleanups)
	if orderErr != nil {
		return errors.Wrap(orderErr, "创建订单失败")
	}

	// 更新Notebook的OrderId
	if err = nb.updateNotebookOrderID(state.nbRef.ID, orderID); err != nil {
		logx.Error("更新Notebook订单ID失败", err)
	}
	state.nbRef.OrderId = orderID

	// === 以下资源创建失败不回滚主流程 ===
	nb.createOptionalResources(
		ctx,
		state.nbRef,
		state.sshPublicKey,
		state.privKeySecretName,
		state.cluster.TensorboardClient,
		state.cluster.RuntimeClient,
		state.cluster.ClientSet,
	)

	// 8. 创建Apisix路由
	return nb.createNotebookAccessRoute(ctx, state.nbRef, true)
}

func (nb *NotebookService) buildCreateNotebookState(ctx context.Context, req *request.AddNoteBookReq, cleanups *Cleanups) (*notebookCreateState, error) {
	req.InstanceName = helper.GenerateInstanceName(consts.NotebookInstance)

	userInfo, err := nb.getNotebookUser(req.UserId)
	if err != nil {
		return nil, err
	}

	imageAddr, err := nb.getNotebookImageAddr(req.ImageId)
	if err != nil {
		return nil, err
	}

	productInfo, err := nb.reserveNotebookProduct(ctx, req, cleanups)
	if err != nil {
		return nil, err
	}

	sshPublicKey, err := nb.getNotebookSSHPublicKey(req.SSHKeyId)
	if err != nil {
		return nil, err
	}

	cluster, err := nb.getNotebookCluster(productInfo.ClusterId)
	if err != nil {
		return nil, err
	}

	nbRef, err := nb.buildNotebookRecord(req, userInfo, productInfo, imageAddr)
	if err != nil {
		return nil, err
	}

	volumeMounts, err := nb.buildNotebookVolumeMounts(req, nbRef.InstanceName)
	if err != nil {
		return nil, err
	}
	nbRef.VolumeMounts = volumeMounts

	return &notebookCreateState{
		req:           req,
		userInfo:      userInfo,
		productInfo:   productInfo,
		imageAddr:     imageAddr,
		sshPublicKey:  sshPublicKey,
		cluster:       cluster,
		pvcManager:    pvc.NewK8sPVCManager(cluster.ClientSet),
		secretManager: secret.NewK8sSecretManager(cluster.ClientSet),
		nbRef:         nbRef,
	}, nil
}

func (nb *NotebookService) getNotebookUser(userID uint) (*systemModel.SysUser, error) {
	var userInfo systemModel.SysUser
	if err := global.GVA_DB.Where("id = ?", userID).First(&userInfo).Error; err != nil {
		logx.Error("用户不存在", err)
		return nil, errors.New("用户不存在")
	}
	return &userInfo, nil
}

func (nb *NotebookService) getNotebookImageAddr(imageID uint) (string, error) {
	if imageID == 0 {
		return "", nil
	}

	var imageInfo imageModel.Image
	if err := global.GVA_DB.Where("id = ?", imageID).First(&imageInfo).Error; err != nil {
		logx.Error("镜像不存在", err)
		return "", errors.New("镜像不存在")
	}

	return imageInfo.ImageAddr, nil
}

func (nb *NotebookService) reserveNotebookProduct(ctx context.Context, req *request.AddNoteBookReq, cleanups *Cleanups) (productRes.ProductDetailResponse, error) {
	var productInfo productRes.ProductDetailResponse
	if req.ProductId == 0 {
		return productInfo, nil
	}

	orderSvc := &orderService.OrderService{}
	if balanceErr := orderSvc.CheckBalanceSufficient(ctx, req.UserId, req.ProductId, req.PayType, req.Quantity); balanceErr != nil {
		return productInfo, balanceErr
	}

	productSvc := &product.ProductService{}
	reserve, reserveErr := productSvc.ReserveCapacity(ctx, req.ProductId)
	if reserveErr != nil {
		logx.Error("产品资源锁定失败", reserveErr)
		return productInfo, reserveErr
	}

	cleanups.Add(func(ctx context.Context) {
		if releaseErr := productSvc.ReleaseCapacity(ctx, req.ProductId, reserve.ResourceCount); releaseErr != nil {
			logx.Error("回滚时释放资源失败", releaseErr)
		}
	})

	return reserve.Product, nil
}

func (nb *NotebookService) getNotebookSSHPublicKey(sshKeyID uint) (string, error) {
	if sshKeyID == 0 {
		return "", nil
	}

	var sshKey secretModelPkg.SSHKey
	if err := global.GVA_DB.Where("id = ?", sshKeyID).First(&sshKey).Error; err != nil {
		logx.Error("SSH密钥不存在", err)
		return "", errors.New("SSH密钥不存在")
	}

	return sshKey.PublicKey, nil
}

func (nb *NotebookService) getNotebookCluster(clusterID uint) (*global.ClusterClientInfo, error) {
	cluster := global.GVA_K8S_CLUSTER_MANAGER.GetCluster(clusterID)
	if cluster == nil {
		logx.Error("没有可用的K8s集群")
		return nil, errors.New("没有可用的K8s集群")
	}
	return cluster, nil
}

func (nb *NotebookService) buildNotebookRecord(
	req *request.AddNoteBookReq,
	userInfo *systemModel.SysUser,
	productInfo productRes.ProductDetailResponse,
	imageAddr string,
) (*nbModel.Notebook, error) {
	nbRef := &nbModel.Notebook{
		DisplayName:        req.DisplayName,
		InstanceName:       req.InstanceName,
		Namespace:          userInfo.Namespace,
		Image:              imageAddr,
		ImageId:            req.ImageId,
		CPU:                productInfo.CPU,
		Memory:             productInfo.Memory,
		GPU:                productInfo.GPUCount,
		GPUModel:           productInfo.GPUModel,
		VGPUNumber:         productInfo.VGPUNumber,
		VGPUMemory:         productInfo.VGPUMemory,
		VGPUCores:          productInfo.VGPUCores,
		StorageSize:        nbModel.DefaultWorkspaceSize,
		Status:             consts.NotebookStatusCreating,
		UserId:             req.UserId,
		ClusterID:          productInfo.ClusterId,
		ProductId:          req.ProductId,
		PayType:            int(req.PayType),
		SSHKeyId:           req.SSHKeyId,
		EnableSSHPassword:  req.EnableSSHPassword,
		EnableTensorboard:  req.TensorBoard,
		TensorboardLogPath: req.TensorBoardLogPath,
	}

	if req.EnableSSHPassword {
		sshPassword, err := password.GenerateRandomPassword(8)
		if err != nil {
			logx.Error("生成SSH密码失败", err)
			return nil, err
		}
		nbRef.SSHPassword = sshPassword
		logx.Info("生成SSH密码成功")
	}

	return nbRef, nil
}

func (nb *NotebookService) buildNotebookVolumeMounts(req *request.AddNoteBookReq, instanceName string) ([]nbModel.NotebookVolume, error) {
	dbVolumeMounts := []nbModel.NotebookVolume{
		{
			Name:       nbModel.Workspace,
			MountsPath: nbModel.DefaultWorkspacePath,
			Size:       nbModel.DefaultWorkspaceSize,
			Type:       nbModel.Workspace,
			PVCName:    fmt.Sprintf("%s-%s", instanceName, nbModel.Workspace),
		},
	}

	for _, mount := range req.VolumeMounts {
		if mount.PVCId == 0 {
			continue
		}

		var vol pvcModelPkg.Volume
		if err := global.GVA_DB.Where("id = ?", mount.PVCId).First(&vol).Error; err != nil {
			logx.Error("数据盘不存在", err)
			return nil, errors.New("数据盘不存在")
		}

		dbVolumeMounts = append(dbVolumeMounts, nbModel.NotebookVolume{
			Name:       vol.Name,
			MountsPath: mount.MountsPath,
			Size:       vol.Size,
			Type:       pvcModelPkg.VolumeTypeToString[vol.Type],
			PVCId:      mount.PVCId,
			PVCName:    vol.PVCName,
		})
	}

	return dbVolumeMounts, nil
}

func (nb *NotebookService) ensureNotebookNamespace(ctx context.Context, clientSet *kubernetes.Clientset, namespace string) error {
	if _, err := clientSet.CoreV1().Namespaces().Get(ctx, namespace, metav1.GetOptions{}); err != nil {
		if !apierrors.IsNotFound(err) {
			logx.Error("获取namespace失败", err)
			return err
		}
		if err := helper.EnsureNamespace(ctx, clientSet, namespace); err != nil {
			logx.Error("创建namespace失败", err)
			return err
		}
	}
	return nil
}

func (nb *NotebookService) createNotebookPVCResources(ctx context.Context, state *notebookCreateState, cleanups *Cleanups) error {
	var pvcMounts []*pvcModel.VolumeMountReq
	for _, mount := range state.nbRef.VolumeMounts {
		if mount.Type != nbModel.Workspace {
			continue
		}
		pvcMounts = append(pvcMounts, &pvcModel.VolumeMountReq{
			Name:       mount.Name,
			MountsPath: mount.MountsPath,
			Size:       mount.Size,
			PVCId:      0,
			Type:       mount.Type,
			PVCName:    mount.PVCName,
		})
	}

	pvcReq := &pvcModel.AddPVCReq{
		InstanceName: state.nbRef.InstanceName,
		Namespace:    state.nbRef.Namespace,
		ClusterId:    state.nbRef.ClusterID,
		VolumeMounts: pvcMounts,
	}
	if err := state.pvcManager.CreatePVCs(ctx, pvcReq, consts.NotebookInstance); err != nil {
		logx.Error("创建PVC失败", err)
		return err
	}

	cleanups.Add(func(cleanupCtx context.Context) {
		_ = state.pvcManager.DeletePVCs(cleanupCtx, &pvcModel.DeletePVCReq{
			InstanceName: state.nbRef.InstanceName,
			Namespace:    state.nbRef.Namespace,
		}, consts.NotebookInstance)
	})

	return nil
}

func (nb *NotebookService) createNotebookSecrets(ctx context.Context, state *notebookCreateState, cleanups *Cleanups) error {
	if err := nb.createNotebookSSHSecrets(ctx, state, cleanups); err != nil {
		return err
	}

	if err := nb.createNotebookPasswordSecret(ctx, state, cleanups); err != nil {
		return err
	}

	return nil
}

func (nb *NotebookService) createNotebookSSHSecrets(ctx context.Context, state *notebookCreateState, cleanups *Cleanups) error {
	if state.sshPublicKey == "" {
		return nil
	}

	privKey, pubKey, err := ssh.GenerateSSHKeyPair()
	if err != nil {
		logx.Error("生成SSH密钥失败", err)
		return err
	}

	addPubKeySecretReq := &secretModel.AddSecretReq{
		InstanceName: state.nbRef.InstanceName,
		Namespace:    state.nbRef.Namespace,
		Content: map[string]string{
			"ssh-publickey": strings.TrimSpace(state.sshPublicKey) + "\n" + pubKey,
		},
		InstanceType: consts.NotebookInstance,
	}
	if err := state.secretManager.CreateSSHSecret(ctx, addPubKeySecretReq); err != nil {
		logx.Error("创建用户+系统SSH公钥Secret失败", err)
		return err
	}
	cleanups.Add(func(cleanupCtx context.Context) {
		_ = state.secretManager.DeleteSSHSecret(cleanupCtx, &secretModel.DeleteSecretReq{
			InstanceName: state.nbRef.InstanceName,
			Namespace:    state.nbRef.Namespace,
		})
	})

	addPrivKeySecretReq := &secretModel.AddSecretReq{
		InstanceName: state.nbRef.InstanceName,
		Namespace:    consts.SSHPiperNamespace,
		Content:      map[string]string{"ssh-privatekey": privKey},
		InstanceType: consts.NotebookInstance,
	}
	privKeySecretName, err := state.secretManager.CreateSSHPrivateKeySecret(ctx, addPrivKeySecretReq)
	if err != nil {
		logx.Error("创建平台私钥SSH私钥Secret失败", err)
		return err
	}
	state.privKeySecretName = privKeySecretName

	cleanups.Add(func(cleanupCtx context.Context) {
		_ = state.secretManager.DeleteSSHPrivateKeySecret(cleanupCtx, &secretModel.DeleteSecretReq{
			InstanceName: state.nbRef.InstanceName,
			Namespace:    consts.SSHPiperNamespace,
		})
	})

	return nil
}

func (nb *NotebookService) createNotebookPasswordSecret(ctx context.Context, state *notebookCreateState, cleanups *Cleanups) error {
	if !state.req.EnableSSHPassword || state.nbRef.SSHPassword == "" {
		return nil
	}

	sshPasswordSecretReq := &secretModel.AddSecretReq{
		InstanceName: state.nbRef.InstanceName,
		Namespace:    state.nbRef.Namespace,
		Content:      map[string]string{"password": state.nbRef.SSHPassword},
		InstanceType: consts.NotebookInstance,
	}
	if err := state.secretManager.CreateSSHPasswordSecret(ctx, sshPasswordSecretReq); err != nil {
		logx.Error("创建SSH密码Secret失败", err)
		return err
	}

	cleanups.Add(func(cleanupCtx context.Context) {
		_ = state.secretManager.DeleteSSHPasswordSecret(cleanupCtx, &secretModel.DeleteSecretReq{
			InstanceName: state.nbRef.InstanceName,
			Namespace:    state.nbRef.Namespace,
		})
	})

	logx.Info("创建SSH密码Secret成功")
	return nil
}

func (nb *NotebookService) createNotebookOrder(ctx context.Context, spec notebookOrderSpec, cleanups *Cleanups) (uint, error) {
	orderSvc := &orderService.OrderService{}
	order, err := orderSvc.CreateOrder(ctx, &orderReq.CreateOrderRequest{
		UserId:       spec.userID,
		ProduceType:  orderModel.ProductTypeCompute,
		ResourceType: orderModel.OrderTypeNotebook,
		ResourceId:   spec.notebookID,
		ProductId:    spec.productID,
		ImageId:      spec.imageID,
		PayType:      consts.PayMethodToInt64[consts.PayMethodBalance],
		ChargeType:   spec.chargeType,
		Quantity:     1,
		Area:         spec.area,
		ClusterId:    spec.clusterID,
		Remark:       fmt.Sprintf("%s: %s", spec.remark, spec.displayName),
	})
	if err != nil {
		logx.Error("创建订单失败", err)
		return 0, err
	}

	cleanups.Add(func(ctx context.Context) {
		_ = global.GVA_DB.Model(&orderModel.Order{}).
			Where("id = ?", order.ID).
			Update("status", orderModel.OrderStatusStopped).Error
	})

	return order.ID, nil
}

func (nb *NotebookService) updateNotebookOrderID(notebookID, orderID uint) error {
	return global.GVA_DB.Model(&nbModel.Notebook{}).
		Where("id = ?", notebookID).
		Update("order_id", orderID).Error
}

func (nb *NotebookService) createNotebookAccessRoute(ctx context.Context, nbRef *nbModel.Notebook, required bool) error {
	if nb.apisixSvc == nil {
		if required {
			logx.Error("Apisix服务未初始化")
			return errors.New("Apisix服务未初始化")
		}
		return nil
	}

	baseDomain := strings.TrimSpace(global.GVA_CONFIG.Apisix.BaseDomain)

	authEnabled := global.GVA_CONFIG.Apisix.AuthEnabled
	authUri := global.GVA_CONFIG.Apisix.AuthUri
	if authEnabled && authUri == "" {
		logx.Error("auth-enabled 为 true 但 auth-uri 未配置，跳过认证")
		return errors.New("auth-enabled 为 true 但 auth-uri 未配置，跳过认证")
	}

	createRouteReq := &apisixReq.CreateRouteReq{
		Name:        fmt.Sprintf("%s-%s", apisixModel.RoutePrefix, nbRef.InstanceName),
		Namespace:   nbRef.Namespace,
		ClusterId:   nbRef.ClusterID,
		Host:        baseDomain,
		Path:        fmt.Sprintf("/notebook/%s/%s/*", nbRef.Namespace, nbRef.InstanceName),
		ServiceName: nbRef.InstanceName,
		ServicePort: 80,
		Labels: map[string]string{
			consts.LabelCreatedBy: consts.LabelValuePlatform,
			consts.LabelInstance:  nbRef.InstanceName,
			consts.LabelType:      consts.NotebookInstance,
		},
		Websocket:  true,
		EnableAuth: authEnabled,
		AuthUri:    authUri,
	}

	if err := nb.apisixSvc.CreateRoute(ctx, createRouteReq); err != nil {
		logx.Error("创建Apisix路由失败", err)
	}

	return nil
}

func (nb *NotebookService) getNotebookTensorboardLogsPath(nbRef *nbModel.Notebook) string {
	subPath := strings.TrimPrefix(nbRef.TensorboardLogPath, "/")
	if subPath == "" {
		subPath = consts.DefaultTensorBoardLogsPath
	}
	return fmt.Sprintf("pvc://%s-%s/%s", nbRef.InstanceName, nbModel.Workspace, subPath)
}

// createOptionalResources 创建可选资源（Tensorboard、SSH Pipe），失败不影响主流程
func (nb *NotebookService) createOptionalResources(
	ctx context.Context,
	nbRef *nbModel.Notebook,
	sshPublicKey string,
	privKeySecretName string,
	tensorboardClient *global.TensorboardClient,
	runtimeClient ctrlclient.Client,
	clusterSet *kubernetes.Clientset,
) {
	if nbRef.EnableTensorboard {
		nb.createNotebookTensorboardResources(ctx, nbRef, tensorboardClient)
	}

	if sshPublicKey != "" || nbRef.EnableSSHPassword {
		nb.createNotebookSSHResources(ctx, nbRef, sshPublicKey, privKeySecretName, runtimeClient, clusterSet)
	}
}

func (nb *NotebookService) createNotebookTensorboardResources(
	ctx context.Context,
	nbRef *nbModel.Notebook,
	tensorboardClient *global.TensorboardClient,
) {
	logsPath := nb.getNotebookTensorboardLogsPath(nbRef)
	tbInstanceName := fmt.Sprintf("%s-tb", nbRef.InstanceName)
	tbManager := tensorboard.NewTensorboardManager(tensorboardClient)
	if err := tbManager.CreateTensorboard(ctx, &tensorboardReq.AddTensorBoardReq{
		InstanceName:      tbInstanceName,
		Namespace:         nbRef.Namespace,
		LogsPath:          logsPath,
		EnableTensorboard: true,
	}); err != nil {
		logx.Error("创建Tensorboard失败", err)
	}

	tensorBoardRef := &tensorboardModel.Tensorboard{
		InstanceName: tbInstanceName,
		Namespace:    nbRef.Namespace,
		OwnerType:    consts.NotebookInstance,
		OwnerID:      nbRef.ID,
		LogsPath:     logsPath,
		Status:       consts.InferenceStatusCreating,
		UserId:       nbRef.UserId,
		ClusterID:    nbRef.ClusterID,
	}
	if err := global.GVA_DB.Create(tensorBoardRef).Error; err != nil {
		logx.Error("保存Tensorboard失败", err)
	}
	if err := global.GVA_DB.Model(&nbModel.Notebook{}).Where("id = ?", nbRef.ID).Update("tensorboard_id", int64(tensorBoardRef.ID)).Error; err != nil {
		logx.Error("更新Notebook TensorboardID失败", err)
	}

	nb.createTensorboardAccessRoute(ctx, nbRef)
}

func (nb *NotebookService) createTensorboardAccessRoute(ctx context.Context, nbRef *nbModel.Notebook) {
	if nb.apisixSvc == nil {
		return
	}

	baseDomain := strings.TrimSpace(global.GVA_CONFIG.Apisix.BaseDomain)

	authEnabled := global.GVA_CONFIG.Apisix.AuthEnabled
	authUri := global.GVA_CONFIG.Apisix.AuthUri
	if authEnabled && authUri == "" {
		logx.Error("auth-enabled 为 true 但 auth-uri 未配置，跳过 TensorBoard 认证")
		authEnabled = false
	}

	tbRouteReq := &apisixReq.CreateRouteReq{
		Name:          fmt.Sprintf("%s-tb-%s", apisixModel.RoutePrefix, nbRef.InstanceName),
		Namespace:     nbRef.Namespace,
		ClusterId:     nbRef.ClusterID,
		Host:          baseDomain,
		Path:          fmt.Sprintf("/tensorboard/%s/%s/*", nbRef.Namespace, nbRef.InstanceName),
		RewriteRegex:  fmt.Sprintf("^/tensorboard/%s/%s/(.*)", nbRef.Namespace, nbRef.InstanceName),
		RewriteTarget: "/$1",
		ServiceName:   fmt.Sprintf("%s-tb", nbRef.InstanceName),
		ServicePort:   80,
		Labels: map[string]string{
			consts.LabelInstance: nbRef.InstanceName,
			consts.LabelType:     consts.TensorBoardInstance,
		},
		Websocket:  false,
		EnableAuth: authEnabled,
		AuthUri:    authUri,
	}

	if err := nb.apisixSvc.CreateRoute(ctx, tbRouteReq); err != nil {
		logx.Error("创建TensorBoard Apisix路由失败", err)
		return
	}

	logx.Info("创建TensorBoard Apisix路由成功")
}

func (nb *NotebookService) createNotebookSSHResources(
	ctx context.Context,
	nbRef *nbModel.Notebook,
	sshPublicKey string,
	privKeySecretName string,
	runtimeClient ctrlclient.Client,
	clusterSet *kubernetes.Clientset,
) {
	sshServiceName := fmt.Sprintf("%s-ssh", nbRef.InstanceName)
	nb.createNotebookSSHService(ctx, nbRef, clusterSet, sshServiceName)
	nb.createNotebookSSHPipe(ctx, nbRef, sshPublicKey, privKeySecretName, runtimeClient, sshServiceName)
	nb.createNotebookSSHStreamRoute(ctx, nbRef)
}

func (nb *NotebookService) createNotebookSSHService(
	ctx context.Context,
	nbRef *nbModel.Notebook,
	clusterSet *kubernetes.Clientset,
	sshServiceName string,
) {
	sshService := &corev1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name:      sshServiceName,
			Namespace: nbRef.Namespace,
			Labels: map[string]string{
				consts.LabelApp:       "ssh",
				consts.LabelOwnerType: consts.NotebookInstance,
			},
		},
		Spec: corev1.ServiceSpec{
			Selector: map[string]string{
				"statefulset": nbRef.InstanceName,
			},
			Ports: []corev1.ServicePort{
				{
					Name:       "ssh",
					Port:       22,
					TargetPort: intstr.FromInt(22),
					Protocol:   corev1.ProtocolTCP,
				},
			},
		},
	}
	if _, err := clusterSet.CoreV1().Services(nbRef.Namespace).Create(ctx, sshService, metav1.CreateOptions{}); err != nil {
		if !apierrors.IsAlreadyExists(err) {
			logx.Error("创建SSH Service失败", err)
		}
		return
	}

	logx.Info("创建SSH Service成功")
}

func (nb *NotebookService) createNotebookSSHPipe(
	ctx context.Context,
	nbRef *nbModel.Notebook,
	sshPublicKey string,
	privKeySecretName string,
	runtimeClient ctrlclient.Client,
	sshServiceName string,
) {
	sshMgr := piper.NewK8sSSHPiperManager(runtimeClient)
	addPipeReq := &pipeReq.AddPipeReq{
		InstanceName:         nbRef.InstanceName,
		Namespace:            consts.SSHPiperNamespace,
		TargetNamespace:      nbRef.Namespace,
		TargetHost:           fmt.Sprintf("%s.%s.svc.cluster.local:22", sshServiceName, nbRef.Namespace),
		TargetUsername:       "root",
		UserSSHKey:           sshPublicKey,
		PrivateKeySecretName: privKeySecretName,
		EnablePasswordAuth:   nbRef.EnableSSHPassword,
		Labels: map[string]string{
			consts.LabelApp:       nbRef.InstanceName,
			consts.LabelOwnerType: consts.NotebookInstance,
		},
	}
	if err := sshMgr.CreatePipe(ctx, addPipeReq); err != nil {
		logx.Error("创建SSH Pipe失败", err)
	}
}

func (nb *NotebookService) createNotebookSSHStreamRoute(ctx context.Context, nbRef *nbModel.Notebook) {
	if nb.apisixSvc == nil {
		return
	}

	defaultPort := global.GVA_CONFIG.SSHPiper.Port
	if defaultPort == 0 {
		defaultPort = apisixModel.DefaultSSHIngressPort
	}

	streamRouteReq := &apisixReq.CreateStreamRouteReq{
		Name:             fmt.Sprintf("%s-sshpiper", apisixModel.StreamRoutePrefix),
		Namespace:        consts.SSHPiperNamespace,
		ClusterId:        nbRef.ClusterID,
		IngressPort:      defaultPort,
		ServiceName:      "sshpiper",
		ServiceNamespace: "",
		ServicePort:      22,
		Labels: map[string]string{
			consts.LabelType: consts.LabelValueSSH,
		},
	}
	if err := nb.apisixSvc.CreateStreamRoute(ctx, streamRouteReq); err != nil {
		if !strings.Contains(err.Error(), "already exists") {
			logx.Error("创建Apisix SSH Stream Route失败（可能已存在）", err)
		}
		return
	}

	logx.Info("创建Apisix SSH Stream Route成功",
		logx.Field("ingressPort", defaultPort),
		logx.Field("backend", "sshpiper"),
	)
}

// cleanupNotebookResources 清理 Notebook 关联的可选 K8s 资源
// 包括：TensorBoard、SSH Pipe、Apisix 路由、SSH Service
// 被 DeleteNotebook 和 StopNotebook 共同调用
func (nb *NotebookService) cleanupNotebookResources(ctx context.Context, dbNb *nbModel.Notebook, cluster *global.ClusterClientInfo) error {
	var errs []error

	nb.appendCleanupError(&errs, nb.deleteNotebookTensorboardResources(ctx, dbNb, cluster.TensorboardClient))
	nb.appendCleanupError(&errs, nb.deleteNotebookSSHPipe(ctx, dbNb, cluster.RuntimeClient))
	nb.appendCleanupError(&errs, nb.deleteNotebookAccessRoutes(ctx, dbNb))
	nb.appendCleanupError(&errs, nb.deleteNotebookSSHService(ctx, dbNb, cluster.ClientSet))

	if len(errs) > 0 {
		return fmt.Errorf("cleanup notebook resources failed with %d errors", len(errs))
	}
	return nil
}

func (nb *NotebookService) appendCleanupError(errs *[]error, err error) {
	if err != nil {
		*errs = append(*errs, err)
	}
}

func (nb *NotebookService) deleteNotebookTensorboardResources(
	ctx context.Context,
	dbNb *nbModel.Notebook,
	tensorboardClient *global.TensorboardClient,
) error {
	if !dbNb.EnableTensorboard {
		return nil
	}

	var errs []error
	tbMgr := tensorboard.NewTensorboardManager(tensorboardClient)
	if err := tbMgr.DeleteTensorboard(ctx, &tensorboardReq.DeleteTensorBoardReq{
		InstanceName: fmt.Sprintf("%s-tb", dbNb.InstanceName),
		Namespace:    dbNb.Namespace,
	}); err != nil && !apierrors.IsNotFound(err) {
		logx.Error("删除TensorBoard K8s资源失败", err)
		errs = append(errs, err)
	}

	if err := global.GVA_DB.Where("owner_id = ? AND owner_type = ?", dbNb.ID, consts.NotebookInstance).Delete(&tensorboardModel.Tensorboard{}).Error; err != nil {
		logx.Error("删除TensorBoard数据库记录失败", err)
		errs = append(errs, err)
	}

	if len(errs) > 0 {
		return fmt.Errorf("delete tensorboard resources failed with %d errors", len(errs))
	}
	return nil
}

func (nb *NotebookService) deleteNotebookSSHPipe(
	ctx context.Context,
	dbNb *nbModel.Notebook,
	runtimeClient ctrlclient.Client,
) error {
	sshMgr := piper.NewK8sSSHPiperManager(runtimeClient)
	pipeName := fmt.Sprintf("pipe-%s-%s", dbNb.Namespace, dbNb.InstanceName)
	if err := sshMgr.DeletePipe(ctx, &pipeReq.DeletePipeReq{
		InstanceName: pipeName,
		Namespace:    consts.SSHPiperNamespace,
	}); err != nil && !apierrors.IsNotFound(err) {
		logx.Error("删除SSH Pipe失败", err)
		return err
	}
	return nil
}

func (nb *NotebookService) deleteNotebookAccessRoutes(ctx context.Context, dbNb *nbModel.Notebook) error {
	if nb.apisixSvc == nil {
		return nil
	}

	var errs []error
	if err := nb.apisixSvc.DeleteRoute(ctx, &apisixReq.DeleteRouteReq{
		Name:      fmt.Sprintf("%s-%s", apisixModel.RoutePrefix, dbNb.InstanceName),
		Namespace: dbNb.Namespace,
		ClusterId: dbNb.ClusterID,
	}); err != nil && !apierrors.IsNotFound(err) {
		logx.Error("删除Notebook Apisix路由失败", err)
		errs = append(errs, err)
	}

	if dbNb.EnableTensorboard {
		if err := nb.apisixSvc.DeleteRoute(ctx, &apisixReq.DeleteRouteReq{
			Name:      fmt.Sprintf("%s-tb-%s", apisixModel.RoutePrefix, dbNb.InstanceName),
			Namespace: dbNb.Namespace,
			ClusterId: dbNb.ClusterID,
		}); err != nil && !apierrors.IsNotFound(err) {
			logx.Error("删除TensorBoard Apisix路由失败", err)
			errs = append(errs, err)
		}
	}

	if len(errs) > 0 {
		return fmt.Errorf("delete notebook access routes failed with %d errors", len(errs))
	}
	return nil
}

func (nb *NotebookService) deleteNotebookSSHService(
	ctx context.Context,
	dbNb *nbModel.Notebook,
	clusterSet *kubernetes.Clientset,
) error {
	sshServiceName := fmt.Sprintf("%s-ssh", dbNb.InstanceName)
	if err := clusterSet.CoreV1().Services(dbNb.Namespace).Delete(ctx, sshServiceName, metav1.DeleteOptions{
		PropagationPolicy: func() *metav1.DeletionPropagation { p := metav1.DeletePropagationBackground; return &p }(),
	}); err != nil && !apierrors.IsNotFound(err) {
		logx.Error("删除SSH Service失败", err)
		return err
	}
	return nil
}

// DeleteNotebook 删除Notebook（彻底销毁，包括 PVC、Secret 和 DB 记录）
// 注意：按量计费的结算由 PodGroup Informer 自动处理
// 删除 Notebook CR 会触发 PodGroup 删除，Informer 会捕获该事件并调用 StopOrder
func (nb *NotebookService) DeleteNotebook(ctx context.Context, req *request.DeleteNoteBookReq) (err error) {
	if req.Id == 0 {
		return errors.New("Notebook ID不能为空")
	}

	nbRef := &nbModel.Notebook{}
	if err = global.GVA_DB.Where("id = ?", req.Id).First(nbRef).Error; err != nil {
		logx.Error("根据ID查询Notebook失败", err)
		return err
	}

	cluster := global.GVA_K8S_CLUSTER_MANAGER.GetCluster(nbRef.ClusterID)
	if cluster == nil {
		return errors.New("没有可用的K8s集群")
	}

	// 资源清理器：使用 Background Context 确保即使请求取消也能完成清理
	defer func() {
		cleanupCtx := context.Background()

		// 清理关联资源（TensorBoard、SSH Pipe、Apisix 路由、SSH Service）
		if err := nb.cleanupNotebookResources(cleanupCtx, nbRef, cluster); err != nil {
			logx.Error("清理 Notebook 关联资源失败", err)
		}

		// 注意：资源配额释放由 PodGroup Informer 的 processDelete 统一处理
		// 删除 Notebook CR → PodGroup 被删除 → Informer 捕获 → ReleaseCapacity
	}()

	// 删除 PVCs（Stop 时保留，Delete 时彻底删除）
	pvcMgr := pvc.NewK8sPVCManager(cluster.ClientSet)
	_ = pvcMgr.DeletePVCs(ctx, &pvcModel.DeletePVCReq{
		InstanceName: nbRef.InstanceName,
		Namespace:    nbRef.Namespace,
	}, consts.NotebookInstance)

	// 删除 SSH Secret（Stop 时保留，Delete 时彻底删除）
	secretMgr := secret.NewK8sSecretManager(cluster.ClientSet)
	if nbRef.SSHKeyId > 0 {
		_ = secretMgr.DeleteSSHSecret(ctx, &secretModel.DeleteSecretReq{
			InstanceName: nbRef.InstanceName,
			Namespace:    nbRef.Namespace,
		})
		_ = secretMgr.DeleteSSHPrivateKeySecret(ctx, &secretModel.DeleteSecretReq{
			InstanceName: nbRef.InstanceName,
			Namespace:    consts.SSHPiperNamespace,
		})
	}
	if nbRef.EnableSSHPassword {
		_ = secretMgr.DeleteSSHPasswordSecret(ctx, &secretModel.DeleteSecretReq{
			InstanceName: nbRef.InstanceName,
			Namespace:    nbRef.Namespace,
		})
	}

	// 更新 DB 状态为 Deleting
	nbRef.Status = consts.NotebookStatusDeleting
	if err = global.GVA_DB.Save(nbRef).Error; err != nil {
		logx.Error("更新数据库状态失败", err)
		return err
	}

	// 删除 K8s Notebook CR
	notebook := &nbv1.Notebook{}
	notebook.Name = nbRef.InstanceName
	notebook.Namespace = nbRef.Namespace
	if err = cluster.NotebookClient.Delete(ctx, notebook, &ctrlclient.DeleteOptions{
		PropagationPolicy: func() *metav1.DeletionPropagation { p := metav1.DeletePropagationBackground; return &p }(),
	}); err != nil {
		if !apierrors.IsNotFound(err) {
			logx.Error("K8s删除Notebook失败", err)
			nbRef.Status = consts.NotebookStatusDeleteFailed
			_ = global.GVA_DB.Save(nbRef).Error
			return err
		}
	}

	// 从数据库软删除（使用主键ID，避免DisplayName重复导致误删）
	if err = global.GVA_DB.Where("id = ?", nbRef.ID).Delete(&nbModel.Notebook{}).Error; err != nil {
		logx.Error("数据库删除失败", err)
		return err
	}

	return nil
}

// UpdateNotebook 更新Notebook（更新镜像或配置）
func (nb *NotebookService) UpdateNotebook(ctx context.Context, req *request.UpdateNoteBookReq) (err error) {
	// 1. 根据ID查询现有Notebook
	nbRef := &nbModel.Notebook{}
	if err = global.GVA_DB.Where("id = ?", req.Id).First(nbRef).Error; err != nil {
		logx.Error("查询Notebook失败", err)
		return err
	}

	// 2. 构建更新字段
	updates := map[string]interface{}{
		"display_name": req.DisplayName,
	}

	// 校验：修改配置（镜像或规格）必须先停止实例
	configChanged := (req.ImageId > 0 && req.ImageId != nbRef.ImageId) || (req.ProductId > 0 && req.ProductId != nbRef.ProductId)
	if configChanged && nbRef.Status != consts.NotebookStatusStopped {
		return errors.New("请先停止实例后再修改配置")
	}

	// 如果更新了镜像
	if req.ImageId > 0 && req.ImageId != nbRef.ImageId {
		var imageInfo imageModel.Image
		if err = global.GVA_DB.Where("id = ?", req.ImageId).First(&imageInfo).Error; err != nil {
			return errors.New("镜像不存在")
		}
		updates["image_id"] = req.ImageId
		updates["image"] = imageInfo.ImageAddr
	}

	// 如果更新了产品（即资源配置）
	if req.ProductId > 0 && req.ProductId != nbRef.ProductId {
		var productInfo productModel.Product
		if err = global.GVA_DB.Where("id = ? AND status = ?", req.ProductId, productModel.ProductStatusEnabled).First(&productInfo).Error; err != nil {
			return errors.New("产品不存在或已下架")
		}
		updates["product_id"] = req.ProductId
		updates["cpu"] = productInfo.CPU
		updates["memory"] = productInfo.Memory
		updates["gpu"] = fmt.Sprintf("%d", productInfo.GPUCount)
		updates["gpu_model"] = productInfo.GPUModel
	}

	// 如果更新了付费类型
	if req.ChargeType > 0 {
		updates["charge_type"] = req.ChargeType
	}

	// 3. 更新数据库
	if err = global.GVA_DB.Model(&nbModel.Notebook{}).Where("id = ?", req.Id).Updates(updates).Error; err != nil {
		logx.Error("更新数据库失败", err)
		return err
	}

	// TODO: 如果需要更新K8s资源（比如重启Pod应用新配置），可以在这里添加逻辑
	// 目前只更新数据库记录

	return nil
}

// GetNotebookList 获取Notebook列表（从数据库获取，并补充实时状态）
func (nb *NotebookService) GetNotebookList(ctx context.Context, req *request.GetNotebookListReq) (resp *response.GetNotebookListResp, err error) {
	dbNotebooks, total, err := nb.listNotebookRecords(ctx, req)
	if err != nil {
		return nil, err
	}

	lookup := nb.loadNotebookResponseLookup(ctx, dbNotebooks)
	notebooks := make([]response.NotebookItem, 0, len(dbNotebooks))
	for _, dbNb := range dbNotebooks {
		stsLister, podLister := nb.getNotebookClusterListers(dbNb.ClusterID)
		item, convertErr := convertDBModelToResponse(dbNb, lookup, stsLister, podLister)
		if convertErr != nil {
			logx.Error("转换Notebook失败, 使用数据库状态", convertErr)
		}
		notebooks = append(notebooks, *item)
	}

	return &response.GetNotebookListResp{
		List:  notebooks,
		Total: total,
	}, nil
}

func (nb *NotebookService) buildNotebookListQuery(ctx context.Context, req *request.GetNotebookListReq) *gorm.DB {
	db := global.GVA_DB.WithContext(ctx).Model(&nbModel.Notebook{})

	if req.UserId > 0 {
		db = db.Where("user_id = ?", req.UserId)
	}
	if req.DisplayName != "" {
		db = db.Where("display_name LIKE ?", "%"+req.DisplayName+"%")
	}
	if req.Status != "" {
		db = db.Where("status = ?", req.Status)
	}

	return db
}

func (nb *NotebookService) listNotebookRecords(ctx context.Context, req *request.GetNotebookListReq) ([]nbModel.Notebook, int64, error) {
	var (
		dbNotebooks []nbModel.Notebook
		total       int64
	)

	db := nb.buildNotebookListQuery(ctx, req)
	if err := db.Count(&total).Error; err != nil {
		logx.Error("查询总数失败", err)
		return nil, 0, err
	}

	if err := db.Preload("VolumeMounts").Scopes(req.Paginate()).Order("id desc").Find(&dbNotebooks).Error; err != nil {
		logx.Error("查询列表失败", err)
		return nil, 0, err
	}

	return dbNotebooks, total, nil
}

func (nb *NotebookService) loadNotebookResponseLookup(ctx context.Context, notebooks []nbModel.Notebook) notebookResponseLookup {
	return notebookResponseLookup{
		imageNames: nb.loadNotebookImageNames(ctx, notebooks),
		products:   nb.loadNotebookProducts(ctx, notebooks),
	}
}

func (nb *NotebookService) loadNotebookImageNames(ctx context.Context, notebooks []nbModel.Notebook) map[uint]string {
	imageIDs := collectNotebookImageIDs(notebooks)
	if len(imageIDs) == 0 {
		return map[uint]string{}
	}

	var images []imageModel.Image
	if err := global.GVA_DB.WithContext(ctx).Where("id IN ?", imageIDs).Find(&images).Error; err != nil {
		logx.Error("批量查询镜像失败", err)
		return map[uint]string{}
	}

	imageNames := make(map[uint]string, len(images))
	for _, image := range images {
		imageNames[image.ID] = image.Name
	}

	return imageNames
}

func (nb *NotebookService) loadNotebookProducts(ctx context.Context, notebooks []nbModel.Notebook) map[uint]productModel.Product {
	productIDs := collectNotebookProductIDs(notebooks)
	if len(productIDs) == 0 {
		return map[uint]productModel.Product{}
	}

	var products []productModel.Product
	if err := global.GVA_DB.WithContext(ctx).Where("id IN ?", productIDs).Find(&products).Error; err != nil {
		logx.Error("批量查询产品失败", err)
		return map[uint]productModel.Product{}
	}

	productMap := make(map[uint]productModel.Product, len(products))
	for _, product := range products {
		productMap[product.ID] = product
	}

	return productMap
}

func collectNotebookImageIDs(notebooks []nbModel.Notebook) []uint {
	return collectNotebookRelationIDs(notebooks, func(notebook nbModel.Notebook) uint {
		return notebook.ImageId
	})
}

func collectNotebookProductIDs(notebooks []nbModel.Notebook) []uint {
	return collectNotebookRelationIDs(notebooks, func(notebook nbModel.Notebook) uint {
		return notebook.ProductId
	})
}

func collectNotebookRelationIDs(notebooks []nbModel.Notebook, getID func(nbModel.Notebook) uint) []uint {
	seen := make(map[uint]struct{})
	ids := make([]uint, 0, len(notebooks))
	for _, notebook := range notebooks {
		id := getID(notebook)
		if id == 0 {
			continue
		}
		if _, ok := seen[id]; ok {
			continue
		}
		seen[id] = struct{}{}
		ids = append(ids, id)
	}
	return ids
}

func (nb *NotebookService) getNotebookClusterListers(clusterID uint) (v1lister.StatefulSetLister, v1podlister.PodLister) {
	cluster := global.GVA_K8S_CLUSTER_MANAGER.GetCluster(clusterID)
	if cluster == nil {
		return nil, nil
	}
	return cluster.StsLister, cluster.PodLister
}

// convertDBModelToResponse 将数据库模型转换为响应结构,补充实时状态
func convertDBModelToResponse(
	dbNb nbModel.Notebook,
	lookup notebookResponseLookup,
	stsLister v1lister.StatefulSetLister,
	podLister v1podlister.PodLister,
) (*response.NotebookItem, error) {
	item := buildNotebookResponseItem(dbNb, lookup)
	if err := refreshNotebookRuntimeStatus(item, dbNb, stsLister, podLister); err != nil {
		return item, err
	}
	return item, nil
}

func buildNotebookResponseItem(dbNb nbModel.Notebook, lookup notebookResponseLookup) *response.NotebookItem {
	item := &response.NotebookItem{
		ID:                 dbNb.ID,
		DisplayName:        dbNb.DisplayName,
		InstanceName:       dbNb.InstanceName,
		Namespace:          dbNb.Namespace,
		Image:              dbNb.Image,
		CPU:                dbNb.CPU,
		Memory:             dbNb.Memory,
		CreationTimestamp:  dbNb.CreatedAt.Format("2006-01-02T15:04:05Z"),
		Status:             dbNb.Status,
		GPUCount:           dbNb.GPU,
		GPUModel:           dbNb.GPUModel,
		PayType:            dbNb.PayType,
		EnableTensorboard:  dbNb.EnableTensorboard,
		TensorboardLogPath: dbNb.TensorboardLogPath,
		VolumeMounts:       dbNb.VolumeMounts,
		CreatedAt:          dbNb.CreatedAt.Format("2006-01-02 15:04:05"),
	}

	enrichNotebookDisplayInfo(item, dbNb, lookup)
	enrichNotebookAccessInfo(item, dbNb)
	return item
}

func enrichNotebookDisplayInfo(item *response.NotebookItem, dbNb nbModel.Notebook, lookup notebookResponseLookup) {
	if imageName, ok := lookup.imageNames[dbNb.ImageId]; ok && imageName != "" {
		item.ImageName = imageName
	}
	if item.ImageName == "" {
		item.ImageName = dbNb.Image
	}

	if product, ok := lookup.products[dbNb.ProductId]; ok {
		item.Price = product.GetPrice(int64(dbNb.PayType))
	}
}

func enrichNotebookAccessInfo(item *response.NotebookItem, dbNb nbModel.Notebook) {
	item.JupyterUrl = fmt.Sprintf("/notebook/%s/%s/lab", dbNb.Namespace, dbNb.InstanceName)
	if dbNb.EnableTensorboard {
		item.TensorboardUrl = fmt.Sprintf("/tensorboard/%s/%s/", dbNb.Namespace, dbNb.InstanceName)
	}

	sshUser := fmt.Sprintf("%s-%s", dbNb.Namespace, dbNb.InstanceName)
	sshHost := global.GVA_CONFIG.SSHPiper.Host
	if sshHost == "" {
		sshHost = global.GVA_CONFIG.Apisix.BaseDomain
	}

	sshPort := global.GVA_CONFIG.SSHPiper.Port
	if sshPort == 0 {
		sshPort = 22
	}

	if sshHost != "" && dbNb.SSHKeyId > 0 {
		item.SSHKeyCommand = fmt.Sprintf("ssh -i ~/.ssh/id_rsa -p %d %s@%s", sshPort, sshUser, sshHost)
	}
	if sshHost != "" && dbNb.EnableSSHPassword && dbNb.SSHPassword != "" {
		item.SSHCommand = fmt.Sprintf("ssh -p %d %s@%s", sshPort, sshUser, sshHost)
		item.SSHPassword = dbNb.SSHPassword
	}
}

func refreshNotebookRuntimeStatus(
	item *response.NotebookItem,
	dbNb nbModel.Notebook,
	stsLister v1lister.StatefulSetLister,
	podLister v1podlister.PodLister,
) error {
	if stsLister == nil {
		return nil
	}

	sts, err := stsLister.StatefulSets(dbNb.Namespace).Get(dbNb.InstanceName)
	if apierrors.IsNotFound(err) {
		logx.Info("StatefulSet不存在，使用数据库状态")
		return nil
	}
	if err != nil {
		logx.Error("获取StatefulSet失败", err)
		return err
	}

	if sts.Status.ReadyReplicas > 0 && podLister != nil {
		return refreshNotebookStatusFromPod(item, dbNb, podLister, sts.Status.ReadyReplicas, sts.Status.Replicas)
	}

	if sts.Status.Replicas > 0 {
		item.Status = consts.NotebookStatusPending
	} else {
		item.Status = consts.NotebookStatusStopped
	}

	return nil
}

func refreshNotebookStatusFromPod(
	item *response.NotebookItem,
	dbNb nbModel.Notebook,
	podLister v1podlister.PodLister,
	readyReplicas int32,
	replicas int32,
) error {
	podName := fmt.Sprintf("%s-0", dbNb.InstanceName)
	pod, err := podLister.Pods(dbNb.Namespace).Get(podName)
	if err != nil {
		if readyReplicas == replicas && replicas > 0 {
			item.Status = consts.NotebookStatusRunning
		} else if replicas > 0 {
			item.Status = consts.NotebookStatusPending
		}
		return nil
	}

	switch pod.Status.Phase {
	case corev1.PodRunning:
		item.Status = consts.NotebookStatusRunning
	case corev1.PodPending:
		item.Status = consts.NotebookStatusPending
	case corev1.PodSucceeded:
		item.Status = consts.NotebookStatusSucceeded
	case corev1.PodFailed:
		item.Status = consts.NotebookStatusFailed
	default:
		item.Status = string(pod.Status.Phase)
	}

	if pod.Status.Phase != corev1.PodRunning {
		return nil
	}

	allReady := true
	for _, containerStatus := range pod.Status.ContainerStatuses {
		if containerStatus.Ready {
			continue
		}
		allReady = false
		if containerStatus.State.Waiting != nil {
			item.Status = fmt.Sprintf("Waiting: %s", containerStatus.State.Waiting.Reason)
		} else {
			item.Status = consts.NotebookStatusPending
		}
		break
	}
	if allReady {
		item.Status = consts.NotebookStatusRunning
	}

	return nil
}

// buildNotebook 构建原生 Notebook 对象
func buildNotebook(nbRef *nbModel.Notebook, sshPublicKey string) *nbv1.Notebook {
	// 使用nbRef中已处理好的配置
	image := nbRef.Image

	// 处理环境变量
	prefix := "/notebook/" + nbRef.Namespace + "/" + nbRef.InstanceName + "/"
	envVars := []corev1.EnvVar{
		{
			// Kubeflow 官方镜像的 entrypoint 读取此变量设置 base_url
			Name:  "NB_PREFIX",
			Value: prefix,
		},
		{
			// Jupyter Server 2.0+ 通过 traitlets 环境变量机制自动识别
			Name:  "JUPYTER_SERVERAPP_BASE_URL",
			Value: prefix,
		},
		{
			// Jupyter Docker Stacks 官方镜像的 start-notebook.sh 读取此变量
			Name:  "JUPYTERHUB_SERVICE_PREFIX",
			Value: prefix,
		},
		{
			// 设置 HOME 目录为工作空间，确保 SSH 登录和 Jupyter terminal 的默认目录一致
			Name:  "HOME",
			Value: nbModel.DefaultWorkspacePath,
		},
	}

	// 处理挂载卷
	var volumeMounts []corev1.VolumeMount
	var volumes []corev1.Volume

	// 1. 处理所有挂载 (包括 workspace)
	for _, m := range nbRef.VolumeMounts {
		pvcName := m.PVCName
		if pvcName == "" {
			pvcName = fmt.Sprintf("%s-%s", nbRef.InstanceName, m.Name)
		}
		volumes = append(volumes, corev1.Volume{
			Name: m.Name,
			VolumeSource: corev1.VolumeSource{
				PersistentVolumeClaim: &corev1.PersistentVolumeClaimVolumeSource{
					ClaimName: pvcName,
				},
			},
		})
		volumeMounts = append(volumeMounts, corev1.VolumeMount{
			Name:      m.Name,
			MountPath: m.MountsPath,
		})
	}

	// 2. 处理 dshm
	volumeMounts = append(volumeMounts, corev1.VolumeMount{
		Name:      "dshm",
		MountPath: "/dev/shm",
	})
	volumes = append(volumes, corev1.Volume{
		Name: "dshm",
		VolumeSource: corev1.VolumeSource{
			EmptyDir: &corev1.EmptyDirVolumeSource{
				Medium: "Memory",
			},
		},
	})

	// 3. 处理 SSH Key (挂载 Secret 到临时目录)
	if sshPublicKey != "" {
		secretName := fmt.Sprintf("%s-ssh-key", nbRef.InstanceName)
		mode := int32(0644) // 临时目录可以宽松一点，sshd-run 会复制并修改权限
		volumeMounts = append(volumeMounts, corev1.VolumeMount{
			Name:      "ssh-key",
			MountPath: "/tmp/ssh-keys", // 挂载到临时目录
			ReadOnly:  true,
		})
		volumes = append(volumes, corev1.Volume{
			Name: "ssh-key",
			VolumeSource: corev1.VolumeSource{
				Secret: &corev1.SecretVolumeSource{
					SecretName:  secretName,
					DefaultMode: &mode,
				},
			},
		})
	}

	// 4. 处理SSH密码Secret（如果启用SSH密码登录）
	if nbRef.EnableSSHPassword && nbRef.SSHPassword != "" {
		passwordSecretName := fmt.Sprintf("%s-ssh-password", nbRef.InstanceName)
		mode := int32(0400) // 密码文件只读
		volumeMounts = append(volumeMounts, corev1.VolumeMount{
			Name:      "ssh-password",
			MountPath: "/etc/ssh-password",
			ReadOnly:  true,
		})
		volumes = append(volumes, corev1.Volume{
			Name: "ssh-password",
			VolumeSource: corev1.VolumeSource{
				Secret: &corev1.SecretVolumeSource{
					SecretName:  passwordSecretName,
					DefaultMode: &mode,
				},
			},
		})
	}

	// 处理资源限制（使用统一的 BuildResources）
	productSpec := &helper.ProductSpec{
		CPU:        nbRef.CPU,
		Memory:     nbRef.Memory,
		GPUCount:   nbRef.GPU,
		VGPUNumber: nbRef.VGPUNumber,
		VGPUMemory: nbRef.VGPUMemory,
		VGPUCores:  nbRef.VGPUCores,
	}
	resourceReqs := helper.BuildResources(productSpec)

	// 构建 Labels
	labels := podgroup.BuildVolcanoLabels(
		nbRef.InstanceName,
		"notebook",
		consts.NotebookInstance,
		nbRef.ID,
	)

	return &nbv1.Notebook{
		ObjectMeta: metav1.ObjectMeta{
			Name:      nbRef.InstanceName,
			Namespace: nbRef.Namespace,
			Labels:    labels,
		},
		Spec: nbv1.NotebookSpec{
			Template: nbv1.NotebookTemplateSpec{
				Spec: corev1.PodSpec{
					// 使用 Volcano 调度器
					SchedulerName: podgroup.VolcanoSchedulerName,
					// Notebook关闭自动挂载k8s凭证
					AutomountServiceAccountToken: pointer.Bool(false),
					Containers: []corev1.Container{
						{
							Name:            nbRef.InstanceName,
							Image:           image,
							ImagePullPolicy: corev1.PullIfNotPresent,
							WorkingDir:      nbModel.DefaultWorkspacePath, // 工作目录
							SecurityContext: &corev1.SecurityContext{
								RunAsUser:  pointer.Int64(0),
								RunAsGroup: pointer.Int64(0),
							},
							Env: envVars,
							Ports: []corev1.ContainerPort{
								{Name: "http", ContainerPort: 8888, Protocol: corev1.ProtocolTCP}, // notebook controller 默认容器监听8888端口
								{Name: "ssh", ContainerPort: 22, Protocol: corev1.ProtocolTCP},    // 镜像自带 sshd 监听 22 端口（由 s6-supervise 管理）
							},
							Resources:    resourceReqs,
							VolumeMounts: volumeMounts,
							Lifecycle:    buildLifecycle(nbRef),
						},
					},
					Volumes: volumes,
				},
			},
		},
	}
}

// sshSetupScript 配置 sshd（密码认证、root 登录、密码设置、公钥安装）
const sshSetupScript = `
mkdir -p /run/sshd /root/.ssh >/dev/null 2>&1 || true
if [ -f /etc/ssh/sshd_config ]; then
  sed -i 's/^#*\s*PasswordAuthentication\s*.*/PasswordAuthentication yes/;s/^#*\s*PermitRootLogin\s*.*/PermitRootLogin yes/' /etc/ssh/sshd_config >/dev/null 2>&1 || true
  grep -q "^PasswordAuthentication" /etc/ssh/sshd_config >/dev/null 2>&1 || echo "PasswordAuthentication yes" >> /etc/ssh/sshd_config || true
  grep -q "^PermitRootLogin" /etc/ssh/sshd_config >/dev/null 2>&1 || echo "PermitRootLogin yes" >> /etc/ssh/sshd_config || true
fi
if [ -f /etc/ssh-password/password ] && command -v chpasswd >/dev/null 2>&1; then
  echo "root:$(cat /etc/ssh-password/password)" | chpasswd >/dev/null 2>&1 || true
fi
if [ -d /tmp/ssh-keys ]; then
  cat /tmp/ssh-keys/* > /root/.ssh/authorized_keys 2>/dev/null || true
  chmod 700 /root/.ssh >/dev/null 2>&1 || true
  chmod 600 /root/.ssh/authorized_keys >/dev/null 2>&1 || true
fi
`

// sshInstallScript 当镜像没有 sshd 时，后台安装 openssh-server 并启动
const sshInstallScript = `
export DEBIAN_FRONTEND=noninteractive
if command -v apt-get >/dev/null 2>&1; then
  apt-get update -qq >/dev/null 2>&1 || true
  apt-get install -y -qq openssh-server >/dev/null 2>&1 || true
fi
mkdir -p /var/run/sshd >/dev/null 2>&1 || true
if command -v ssh-keygen >/dev/null 2>&1; then
  [ ! -f /etc/ssh/ssh_host_rsa_key ] && ssh-keygen -t rsa -f /etc/ssh/ssh_host_rsa_key -N '' >/dev/null 2>&1 || true
  [ ! -f /etc/ssh/ssh_host_ed25519_key ] && ssh-keygen -t ed25519 -f /etc/ssh/ssh_host_ed25519_key -N '' >/dev/null 2>&1 || true
fi
`

// buildLifecycle 构建容器 PostStart 钩子，确保 sshd 可用
// 镜像自带 sshd → 直接配置并重载；没有 sshd → 后台安装后启动（避免 PostStart 超时）
func buildLifecycle(nbRef *nbModel.Notebook) *corev1.Lifecycle {
	if nbRef.SSHKeyId == 0 && !nbRef.EnableSSHPassword {
		return nil
	}

	script := fmt.Sprintf(`
(
  set +e
  SSH_BIN="$(command -v sshd 2>/dev/null || true)"
  HAS_PGREP=0
  if command -v pgrep >/dev/null 2>&1; then
    HAS_PGREP=1
  fi
  HAS_S6=0
  if [ -d /run/s6 ] || [ -d /var/run/s6 ] || [ -d /etc/cont-init.d ] || [ -x /init ]; then
    HAS_S6=1
  fi

  if [ -n "$SSH_BIN" ]; then
    %s

    # 已存在 sshd 时不重复拉起；非 s6 镜像可做轻量 HUP
    if [ "$HAS_PGREP" = "1" ] && pgrep -x sshd >/dev/null 2>&1; then
      if [ "$HAS_S6" != "1" ] && command -v pkill >/dev/null 2>&1; then
        pkill -HUP -x sshd >/dev/null 2>&1 || true
      fi
    elif [ "$HAS_S6" = "1" ]; then
      # s6 镜像由自身服务编排接管，不在 PostStart 里启动额外 sshd
      :
    else
      "$SSH_BIN" -D >/dev/null 2>&1 &
    fi
  elif [ "$HAS_S6" = "1" ]; then
    # s6 镜像且无 sshd 二进制：跳过安装/启动，避免与镜像初始化冲突
    :
  else
    (
      %s
      %s
      if [ -x /usr/sbin/sshd ]; then
        nohup /usr/sbin/sshd -D >/dev/null 2>&1 &
      fi
    ) >/tmp/neptune-ssh-init.log 2>&1 &
  fi
) >/tmp/neptune-poststart.log 2>&1 || true
exit 0`, sshSetupScript, sshInstallScript, sshSetupScript)

	return &corev1.Lifecycle{
		PostStart: &corev1.LifecycleHandler{
			Exec: &corev1.ExecAction{
				Command: []string{"/bin/sh", "-c", script},
			},
		},
	}
}

// GetNotebookDetail 获取Notebook详情
func (nb *NotebookService) GetNotebookDetail(ctx context.Context, req *request.GetNotebookDetailReq) (*response.NotebookItem, error) {
	var dbNb nbModel.Notebook
	if err := global.GVA_DB.WithContext(ctx).Preload("VolumeMounts").Where("id = ?", req.ID).First(&dbNb).Error; err != nil {
		return nil, errors.Wrap(err, "Notebook不存在")
	}

	lookup := nb.loadNotebookResponseLookup(ctx, []nbModel.Notebook{dbNb})
	stsLister, podLister := nb.getNotebookClusterListers(dbNb.ClusterID)
	return convertDBModelToResponse(dbNb, lookup, stsLister, podLister)
}

// GetNotebookLogs 获取Notebook日志
func (nb *NotebookService) GetNotebookLogs(ctx context.Context, req *request.GetNotebookLogsReq) (io.ReadCloser, error) {
	var dbNb nbModel.Notebook
	if err := global.GVA_DB.Where("id = ?", req.ID).First(&dbNb).Error; err != nil {
		return nil, errors.Wrap(err, "Notebook不存在")
	}

	cluster := global.GVA_K8S_CLUSTER_MANAGER.GetCluster(dbNb.ClusterID)
	if cluster == nil {
		return nil, errors.New("集群不存在")
	}

	podName := fmt.Sprintf("%s-0", dbNb.InstanceName)
	container := req.Container
	if container == "" {
		container = dbNb.InstanceName
	}

	logOptions := &corev1.PodLogOptions{
		Container:  container,
		Follow:     req.Follow,
		Timestamps: req.Timestamps,
	}
	if req.TailLines > 0 {
		logOptions.TailLines = &req.TailLines
	}

	logReq := cluster.ClientSet.CoreV1().Pods(dbNb.Namespace).GetLogs(podName, logOptions)
	return logReq.Stream(ctx)
}

// GetNotebookPods 获取Notebook的Pod列表
func (nb *NotebookService) GetNotebookPods(ctx context.Context, req *request.GetNotebookPodsReq) ([]response.PodInfoResp, error) {
	var dbNb nbModel.Notebook
	if err := global.GVA_DB.Where("id = ?", req.ID).First(&dbNb).Error; err != nil {
		return nil, errors.Wrap(err, "Notebook不存在")
	}

	cluster := global.GVA_K8S_CLUSTER_MANAGER.GetCluster(dbNb.ClusterID)
	if cluster == nil {
		return nil, errors.New("集群不存在")
	}

	labelSelector := fmt.Sprintf("app=%s", dbNb.InstanceName)
	podList, err := cluster.ClientSet.CoreV1().Pods(dbNb.Namespace).List(ctx, metav1.ListOptions{
		LabelSelector: labelSelector,
	})
	if err != nil {
		return nil, errors.Wrap(err, "获取 Pod 列表失败")
	}

	result := make([]response.PodInfoResp, 0, len(podList.Items))
	for _, pod := range podList.Items {
		containers := make([]string, 0, len(pod.Spec.Containers))
		for _, container := range pod.Spec.Containers {
			containers = append(containers, container.Name)
		}

		result = append(result, response.PodInfoResp{
			Name:       pod.Name,
			Namespace:  pod.Namespace,
			Status:     getNotebookPodDisplayStatus(&pod),
			HostIP:     pod.Status.HostIP,
			PodIP:      pod.Status.PodIP,
			Containers: containers,
		})
	}
	return result, nil
}

func getNotebookPodDisplayStatus(pod *corev1.Pod) string {
	if pod == nil {
		return ""
	}

	if pod.Status.Phase != corev1.PodRunning {
		return string(pod.Status.Phase)
	}

	if isNotebookPodReady(pod) {
		return string(corev1.PodRunning)
	}

	return string(corev1.PodPending)
}

func isNotebookPodReady(pod *corev1.Pod) bool {
	for _, condition := range pod.Status.Conditions {
		if condition.Type == corev1.PodReady {
			return condition.Status == corev1.ConditionTrue
		}
	}

	if len(pod.Status.ContainerStatuses) == 0 {
		return false
	}

	for _, containerStatus := range pod.Status.ContainerStatuses {
		if !containerStatus.Ready {
			return false
		}
	}

	return true
}

// StartNotebook 启动Notebook
func (nb *NotebookService) StartNotebook(ctx context.Context, id uint) (err error) {
	var dbNb nbModel.Notebook
	if err = global.GVA_DB.Preload("VolumeMounts").Where("id = ?", id).First(&dbNb).Error; err != nil {
		return errors.Wrap(err, "Notebook不存在")
	}

	if dbNb.Status != consts.NotebookStatusStopped {
		return errors.New("Notebook不是停止状态，无法启动")
	}

	cluster := global.GVA_K8S_CLUSTER_MANAGER.GetCluster(dbNb.ClusterID)
	if cluster == nil {
		return errors.New("集群不存在")
	}

	// 定义清理函数栈，用于发生错误时回滚
	cleanups := make(Cleanups, 0)
	defer func() {
		if err != nil {
			cleanups.Run(context.Background())
		}
	}()

	// 1. 检查并锁定资源（所有产品类型统一锁定 1 个实例配额）
	if dbNb.ProductId > 0 {
		productSvc := &product.ProductService{}
		reserve, reserveErr := productSvc.ReserveCapacity(ctx, dbNb.ProductId)
		if reserveErr != nil {
			return errors.Wrap(reserveErr, "锁定资源失败")
		}
		lockedCount := reserve.ResourceCount

		// 添加回滚逻辑：释放资源
		cleanups.Add(func(ctx context.Context) {
			_ = productSvc.ReleaseCapacity(ctx, dbNb.ProductId, lockedCount)
		})
	}

	// 2. 创建订单（仅针对按量付费的Notebook）
	if dbNb.PayType == productModel.ChargeTypeHourly {
		orderID, orderErr := nb.createNotebookOrder(ctx, notebookOrderSpec{
			userID:      dbNb.UserId,
			notebookID:  dbNb.ID,
			productID:   dbNb.ProductId,
			imageID:     dbNb.ImageId,
			chargeType:  int64(dbNb.PayType),
			displayName: dbNb.DisplayName,
			remark:      "重启容器实例",
			area:        cluster.Area,
			clusterID:   dbNb.ClusterID,
		}, &cleanups)
		if orderErr != nil {
			return errors.Wrap(orderErr, "创建订单失败")
		}

		// 更新Notebook的OrderId
		if err = nb.updateNotebookOrderID(dbNb.ID, orderID); err != nil {
			logx.Error("更新Notebook订单ID失败", err)
			return errors.Wrap(err, "更新Notebook订单ID失败")
		}
		dbNb.OrderId = orderID
	}

	// 3. 获取 SSH 公钥
	sshPublicKey, sshErr := nb.getNotebookSSHPublicKey(dbNb.SSHKeyId)
	if sshErr != nil {
		return sshErr
	}

	// 4. 构建并创建 Notebook CR
	notebookObj := buildNotebook(&dbNb, sshPublicKey)
	if err = cluster.NotebookClient.Create(ctx, notebookObj); err != nil {
		return errors.Wrap(err, "创建Notebook失败")
	}
	// 如果需要，可以在这里添加删除 Notebook CR 的回滚逻辑
	// 但由于后续步骤（可选资源、Apisix）失败通常不阻断流程，所以这里暂时不需要

	// 5. 创建可选资源
	var privKeySecretName string
	if sshPublicKey != "" {
		privKeySecretName = fmt.Sprintf("%s-ssh-private-key", dbNb.InstanceName)
	}
	nb.createOptionalResources(ctx, &dbNb, sshPublicKey, privKeySecretName, cluster.TensorboardClient, cluster.RuntimeClient, cluster.ClientSet)

	// 6. 创建 Apisix 路由
	if err = nb.createNotebookAccessRoute(ctx, &dbNb, false); err != nil {
		return err
	}

	// 7. 更新数据库状态
	return global.GVA_DB.Model(&dbNb).Update("status", consts.NotebookStatusPending).Error
}

// StopNotebook 停止Notebook（保留 PVC 和 Secret，可重新启动）
// 注意：按量计费的结算由 PodGroup Informer 自动处理
// 删除 Notebook CR 会触发 PodGroup 删除，Informer 会捕获该事件并调用 StopOrder
func (nb *NotebookService) StopNotebook(ctx context.Context, id uint) error {
	var dbNb nbModel.Notebook
	if err := global.GVA_DB.Where("id = ?", id).First(&dbNb).Error; err != nil {
		return errors.Wrap(err, "Notebook不存在")
	}

	if dbNb.Status == consts.NotebookStatusStopped {
		return nil
	}

	cluster := global.GVA_K8S_CLUSTER_MANAGER.GetCluster(dbNb.ClusterID)
	if cluster == nil {
		return errors.New("集群不存在")
	}

	// 资源清理器：使用 Background Context 确保即使请求取消也能完成清理
	defer func() {
		cleanupCtx := context.Background()

		// 清理关联资源（TensorBoard、SSH Pipe、Apisix 路由、SSH Service）
		if err := nb.cleanupNotebookResources(cleanupCtx, &dbNb, cluster); err != nil {
			logx.Error("停止 Notebook 时清理关联资源失败", err)
		}

		// 注意：资源配额释放由 PodGroup Informer 的 processDelete 统一处理
		// 删除 Notebook CR → PodGroup 被删除 → Informer 捕获 → ReleaseCapacity
	}()

	// 更新数据库状态
	if err := global.GVA_DB.Model(&dbNb).Update("status", consts.NotebookStatusStopped).Error; err != nil {
		logx.Error("更新数据库状态失败", err)
		return err
	}

	// 删除 K8s Notebook CR
	notebook := &nbv1.Notebook{
		ObjectMeta: metav1.ObjectMeta{
			Name:      dbNb.InstanceName,
			Namespace: dbNb.Namespace,
		},
	}
	if err := cluster.NotebookClient.Delete(ctx, notebook, &ctrlclient.DeleteOptions{
		PropagationPolicy: func() *metav1.DeletionPropagation { p := metav1.DeletePropagationBackground; return &p }(),
	}); err != nil && !apierrors.IsNotFound(err) {
		logx.Error("删除Notebook CR失败", err)
		return errors.Wrap(err, "停止Notebook失败")
	}

	return nil
}

// GetTerminalInfo 获取 Terminal 连接所需信息（DB查询+状态校验）
func (nb *NotebookService) GetTerminalInfo(ctx context.Context, req *request.HandleTerminalReq) (*response.TerminalInfoResp, error) {
	var notebook nbModel.Notebook
	if err := global.GVA_DB.Where("id = ?", req.ID).First(&notebook).Error; err != nil {
		return nil, errors.Wrap(err, "Notebook不存在")
	}
	if notebook.Status != consts.NotebookStatusRunning {
		return nil, errors.New("Notebook未在运行中，无法连接终端")
	}
	containerName := req.Container
	if containerName == "" {
		containerName = notebook.InstanceName
	}
	return &response.TerminalInfoResp{
		Namespace:    notebook.Namespace,
		PodName:      notebook.InstanceName + "-0",
		Container:    containerName,
		ClusterID:    notebook.ClusterID,
		InstanceName: notebook.InstanceName,
	}, nil
}
