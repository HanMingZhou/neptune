package training

import (
	"context"
	"fmt"
	"gin-vue-admin/global"
	apisixModel "gin-vue-admin/model/apisix"
	apisixReq "gin-vue-admin/model/apisix/request"
	"gin-vue-admin/model/consts"
	tensorboardModel "gin-vue-admin/model/tensorboard"
	tensorboardReq "gin-vue-admin/model/tensorboard/request"
	trainingModel "gin-vue-admin/model/training"
	trainingReq "gin-vue-admin/model/training/request"
	"gin-vue-admin/service/tensorboard"
	"path"
	"strings"

	"github.com/zeromicro/go-zero/core/logx"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// createOptionalResources 创建可选资源
// 这些资源创建失败不影响主流程
func (s *TrainingJobService) createOptionalResources(
	ctx context.Context,
	job *trainingModel.TrainingJob,
	req *trainingReq.CreateTrainingJobReq,
	cluster *global.ClusterClientInfo,
) {
	if req.EnableTensorboard {
		s.createTensorboard(ctx, job, req, cluster)
	}
}

// createTensorboard 创建 TensorBoard 相关资源
func (s *TrainingJobService) createTensorboard(
	ctx context.Context,
	job *trainingModel.TrainingJob,
	req *trainingReq.CreateTrainingJobReq,
	cluster *global.ClusterClientInfo,
) {
	logsPath, err := buildTrainingTensorboardLogsPath(req)
	if err != nil {
		logx.Error("创建TensorBoard失败", logx.Field("err", err), logx.Field("job", job.JobName))
		return
	}

	tbManager := tensorboard.NewTensorboardManager(cluster.TensorboardClient)
	tbReq := &tensorboardReq.AddTensorBoardReq{
		InstanceName:      buildTrainingTensorboardName(job.JobName),
		Namespace:         job.Namespace,
		LogsPath:          logsPath,
		EnableTensorboard: true,
	}
	if err := tbManager.CreateTensorboard(ctx, tbReq); err != nil {
		logx.Error("创建TensorBoard失败", logx.Field("err", err), logx.Field("job", job.JobName))
		return
	}

	tbRecord := &tensorboardModel.Tensorboard{
		InstanceName: buildTrainingTensorboardName(job.JobName),
		Namespace:    job.Namespace,
		OwnerType:    consts.TrainingInstance,
		OwnerID:      job.ID,
		LogsPath:     logsPath,
		Status:       consts.StatusCreating,
		UserId:       job.UserId,
		ClusterID:    job.ClusterID,
	}
	if err := global.GVA_DB.Create(tbRecord).Error; err != nil {
		logx.Error("保存TensorBoard记录失败", logx.Field("err", err))
	}

	if err := global.GVA_DB.Model(&trainingModel.TrainingJob{}).
		Where("id = ?", job.ID).
		Update("tensorboard_id", int64(tbRecord.ID)).Error; err != nil {
		logx.Error("更新TrainingJob TensorboardId失败", logx.Field("err", err))
	}

	job.TensorboardId = tbRecord.ID
	s.createTensorboardRoute(ctx, job)

	logx.Info("训练任务TensorBoard创建成功",
		logx.Field("job", job.JobName),
		logx.Field("path", logsPath),
	)
}

// createTensorboardRoute 创建 TensorBoard 的 Apisix 路由
func (s *TrainingJobService) createTensorboardRoute(ctx context.Context, job *trainingModel.TrainingJob) {
	if s.apisixSvc == nil {
		return
	}

	tbRouteReq := buildTrainingTensorboardRouteReq(job)
	if err := s.apisixSvc.CreateRoute(ctx, tbRouteReq); err != nil {
		logx.Error("创建TensorBoard Apisix路由失败", logx.Field("err", err))
		return
	}

	logx.Info("创建TensorBoard Apisix路由成功", logx.Field("path", tbRouteReq.Path))
}

func buildTrainingTensorboardLogsPath(req *trainingReq.CreateTrainingJobReq) (string, error) {
	if len(req.Mounts) == 0 {
		return "", fmt.Errorf("任务未挂载任何PVC，无法存放日志")
	}

	mount := req.Mounts[0]
	logRelPath := req.TensorboardLogPath
	if logRelPath == "" {
		logRelPath = consts.DefaultTensorBoardLogsPath
	}

	innerPath := path.Join(mount.SubPath, logRelPath)
	innerPath = strings.Trim(innerPath, "/")
	return fmt.Sprintf("pvc://%s/%s", mount.PvcName, innerPath), nil
}

func buildTrainingTensorboardName(jobName string) string {
	return fmt.Sprintf("%s-tb", jobName)
}

func buildTrainingTensorboardRouteReq(job *trainingModel.TrainingJob) *apisixReq.CreateRouteReq {
	baseDomain := strings.TrimSpace(global.GVA_CONFIG.Apisix.BaseDomain)

	tbPathMatch := fmt.Sprintf("/tensorboard/%s/%s/*", job.Namespace, job.JobName)
	tbPathRegex := fmt.Sprintf("^/tensorboard/%s/%s/(.*)", job.Namespace, job.JobName)
	tbRouteName := fmt.Sprintf("%s-tb-%s", apisixModel.RoutePrefix, job.JobName)

	authEnabled := global.GVA_CONFIG.Apisix.AuthEnabled
	authUri := global.GVA_CONFIG.Apisix.AuthUri
	if authEnabled && authUri == "" {
		logx.Error("auth-enabled 为 true 但 auth-uri 未配置，跳过 TensorBoard 认证")
		authEnabled = false
	}

	return &apisixReq.CreateRouteReq{
		Name:          tbRouteName,
		Namespace:     job.Namespace,
		ClusterId:     job.ClusterID,
		Host:          baseDomain,
		Path:          tbPathMatch,
		RewriteRegex:  tbPathRegex,
		RewriteTarget: "/$1",
		ServiceName:   fmt.Sprintf("%s-tb", job.JobName),
		ServicePort:   80, // Kubeflow TensorBoard Service 默认对外暴露 80 端口
		Labels: map[string]string{
			consts.LabelInstance: job.JobName,
			consts.LabelType:     consts.TensorBoardInstance,
		},
		Websocket:  false,
		EnableAuth: authEnabled,
		AuthUri:    authUri,
	}
}

// deleteOptionalResources 删除可选资源
func (s *TrainingJobService) deleteOptionalResources(
	ctx context.Context,
	job *trainingModel.TrainingJob,
	cluster *global.ClusterClientInfo,
) error {
	// 删除 TensorBoard 资源
	if job.EnableTensorboard {
		return s.deleteTensorboardResources(ctx, job, cluster)
	}
	return nil
}

// deleteTensorboardResources 删除 TensorBoard 相关资源
func (s *TrainingJobService) deleteTensorboardResources(
	ctx context.Context,
	job *trainingModel.TrainingJob,
	cluster *global.ClusterClientInfo,
) error {
	tbName := buildTrainingTensorboardName(job.JobName)

	s.deleteTensorboardService(ctx, cluster, job.Namespace, tbName)
	s.deleteTensorboardCR(ctx, cluster, job.Namespace, tbName)
	s.deleteTensorboardRoute(ctx, job)

	if err := global.GVA_DB.Where("owner_type = ? AND owner_id = ?", consts.TrainingInstance, job.ID).Delete(&tensorboardModel.Tensorboard{}).Error; err != nil {
		logx.Error("删除 TensorBoard 数据库记录失败", logx.Field("err", err))
		return err
	}

	logx.Info("删除训练任务TensorBoard资源成功", logx.Field("job", job.JobName))
	return nil
}

func (s *TrainingJobService) deleteTensorboardService(
	ctx context.Context,
	cluster *global.ClusterClientInfo,
	namespace, serviceName string,
) {
	if err := cluster.ClientSet.CoreV1().Services(namespace).Delete(ctx, serviceName, metav1.DeleteOptions{
		PropagationPolicy: func() *metav1.DeletionPropagation { p := metav1.DeletePropagationBackground; return &p }(),
	}); err != nil && !isTrainingResourceNotFound(err) {
		logx.Error("删除 TensorBoard Service 失败", logx.Field("err", err))
	}
}

func (s *TrainingJobService) deleteTensorboardCR(
	ctx context.Context,
	cluster *global.ClusterClientInfo,
	namespace, instanceName string,
) {
	tbManager := tensorboard.NewTensorboardManager(cluster.TensorboardClient)
	if err := tbManager.DeleteTensorboard(ctx, &tensorboardReq.DeleteTensorBoardReq{
		InstanceName: instanceName,
		Namespace:    namespace,
	}); err != nil && !isTrainingResourceNotFound(err) {
		logx.Error("删除 TensorBoard CR 失败", logx.Field("err", err))
	}
}

func (s *TrainingJobService) deleteTensorboardRoute(ctx context.Context, job *trainingModel.TrainingJob) {
	if s.apisixSvc == nil {
		return
	}

	routeName := fmt.Sprintf("%s-tb-%s", apisixModel.RoutePrefix, job.JobName)
	if err := s.apisixSvc.DeleteRoute(ctx, &apisixReq.DeleteRouteReq{
		Name:      routeName,
		Namespace: job.Namespace,
		ClusterId: job.ClusterID,
	}); err != nil && !isTrainingResourceNotFound(err) {
		logx.Error("删除 TensorBoard Apisix 路由失败", logx.Field("err", err))
	}
}

func isTrainingResourceNotFound(err error) bool {
	if err == nil {
		return false
	}
	if apierrors.IsNotFound(err) {
		return true
	}
	return strings.Contains(strings.ToLower(err.Error()), "not found")
}
