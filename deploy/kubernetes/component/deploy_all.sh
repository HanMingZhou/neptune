#!/bin/bash

# ==============================================================================
# Neptune 一键部署脚本
# 包含组件: Volcano, Kubeflow, APISIX, SSHPiper
# 支持重复运行（幂等）
# ==============================================================================

set -e

# 颜色定义
red_color='\033[0;31m'
green_color='\033[0;32m'
yellow_color='\033[1;33m'
nc_color='\033[0m'

function log_info() {
    echo -e "${green_color}[INFO] $1${nc_color}"
}

function log_warn() {
    echo -e "${yellow_color}[WARN] $1${nc_color}"
}

function log_error() {
    echo -e "${red_color}[ERROR] $1${nc_color}"
}

# 检查必要工具
for cmd in helm kubectl kustomize git; do
    if ! command -v $cmd &> /dev/null; then
        log_error "未找到 $cmd 工具，请先安装。"
        exit 1
    fi
done

BASE_DIR="$(cd "$(dirname "$0")" && pwd)"

# ==============================================================================
# 1. 部署 Volcano
# ==============================================================================
log_info "正在部署 Volcano..."
helm repo add volcano-sh https://volcano-sh.github.io/helm-charts
helm repo update
helm upgrade --install volcano volcano-sh/volcano -n volcano-system --create-namespace
log_info "Volcano 部署完成。"

# ==============================================================================
# 2. 部署 Kubeflow 组件
# ==============================================================================
log_info "正在部署 Kubeflow 组件..."
cd "$BASE_DIR/kubeflow"

if [ ! -d "manifests" ]; then
    log_info "正在克隆 Kubeflow manifests 仓库..."
    git clone https://github.com/kubeflow/manifests.git
fi

MANIFESTS_DIR="$BASE_DIR/kubeflow/manifests"

log_info "1/4 安装 Cert Manager..."
kubectl apply -k "$MANIFESTS_DIR/common/cert-manager/base"
log_info "等待 Cert Manager CRD 注册..."
kubectl wait --for=condition=Established crd/certificates.cert-manager.io --timeout=120s 2>/dev/null || true
kubectl wait --for=condition=Established crd/issuers.cert-manager.io --timeout=120s 2>/dev/null || true
kubectl wait --for=condition=Established crd/clusterissuers.cert-manager.io --timeout=120s 2>/dev/null || true
log_info "等待 Cert Manager Deployment 就绪..."
kubectl wait --for=condition=Available deployment/cert-manager -n cert-manager --timeout=120s 2>/dev/null || true
kubectl wait --for=condition=Available deployment/cert-manager-cainjector -n cert-manager --timeout=120s 2>/dev/null || true
kubectl wait --for=condition=Available deployment/cert-manager-webhook -n cert-manager --timeout=120s 2>/dev/null || true

log_info "2/4 安装 Kubeflow Issuer..."
kubectl apply -k "$MANIFESTS_DIR/common/cert-manager/kubeflow-issuer/base"

log_info "3/4 安装 Istio CRDs..."
kubectl apply -k "$MANIFESTS_DIR/common/istio/istio-crds/base"

log_info "4/4 安装 Notebook & Tensorboard 控制器..."
kubectl apply -k "$MANIFESTS_DIR/applications/jupyter/notebook-controller/upstream/overlays/kubeflow"
kubectl apply -k "$MANIFESTS_DIR/applications/tensorboard/tensorboard-controller/upstream/overlays/kubeflow"

log_info "等待 Notebook Controller 就绪..."
kubectl wait --for=condition=Available deployment/notebook-controller-deployment -n kubeflow --timeout=180s 2>/dev/null || true
log_info "等待 TensorBoard Controller 就绪..."
kubectl wait --for=condition=Available deployment/tensorboard-controller-deployment -n kubeflow --timeout=180s 2>/dev/null || true

log_info "Kubeflow 组件部署完成。"

# ==============================================================================
# 3. 部署 APISIX
# ==============================================================================
log_info "正在部署 APISIX..."
helm repo add apisix https://apache.github.io/apisix-helm-chart
helm repo update

helm upgrade --install apisix apisix/apisix \
  --create-namespace \
  --namespace apisix \
  --set ingress-controller.enabled=true \
  --set ingress-controller.gatewayProxy.createDefault=true \
  --set ingress-controller.gatewayProxy.provider.controlPlane.service.name=apisix-admin \
  --set ingress-controller.gatewayProxy.provider.controlPlane.service.port=9180 \
  --set ingress-controller.gatewayProxy.provider.controlPlane.auth.adminKey.value=edd1c9f034335f136f87ad84b625c8f1 \
  --set ingress-controller.config.apisix.serviceNamespace=apisix \
  --set 'apisix.proxy_mode=http&stream' \
  --set "apisix.stream_proxy.tcp[0]=9100"

log_info "等待 APISIX 就绪..."
kubectl wait --for=condition=Available deployment/apisix -n apisix --timeout=180s 2>/dev/null || true

# 修复 proxy_mode（helm 的 & 转义可能导致 stream 模式未生效）
log_info "检查并修复 proxy_mode 配置..."
CURRENT_MODE=$(kubectl get configmap apisix -n apisix -o jsonpath='{.data.config\.yaml}' 2>/dev/null | grep proxy_mode || echo "")
if echo "$CURRENT_MODE" | grep -q 'http&stream'; then
    log_info "proxy_mode 已正确设置为 http&stream"
else
    log_warn "proxy_mode 未包含 stream，正在修复..."
    # 使用 python/perl 避免 sed 中 & 的特殊含义问题
    kubectl get configmap apisix -n apisix -o jsonpath='{.data.config\.yaml}' \
      | perl -pe 's/proxy_mode:\s*http\s*$/proxy_mode: "http\&stream"/' \
      | kubectl create configmap apisix -n apisix --from-file=config.yaml=/dev/stdin --dry-run=client -o yaml \
      | kubectl apply -f -
    kubectl rollout restart deployment apisix -n apisix
    kubectl wait --for=condition=Available deployment/apisix -n apisix --timeout=120s 2>/dev/null || true
fi

# 添加 SSH 端口
log_info "为 APISIX Gateway 添加 SSH 端口..."
kubectl patch svc apisix-gateway -n apisix --type='json' \
  -p='[{"op":"add","path":"/spec/ports/-","value":{"name":"ssh","port":22,"targetPort":9100,"protocol":"TCP"}}]' 2>/dev/null || log_warn "SSH 端口可能已存在，跳过。"

log_info "APISIX 部署完成。"

# ==============================================================================
# 4. 部署 SSHPiper
# ==============================================================================
log_info "正在部署 SSHPiper..."

log_info "安装 SSHPiper CRD..."
kubectl apply -f https://raw.githubusercontent.com/tg123/sshpiper/master/plugin/kubernetes/crd.yaml
log_info "等待 CRD 注册..."
kubectl wait --for=condition=Established crd/pipes.sshpiper.com --timeout=60s 2>/dev/null || true

if [ -f "$BASE_DIR/sshpiper/deploy.yml" ]; then
    kubectl apply -f "$BASE_DIR/sshpiper/deploy.yml"
    log_info "等待 SSHPiper 就绪..."
    kubectl wait --for=condition=Available deployment/sshpiper -n kubeflow --timeout=120s 2>/dev/null || true
    log_info "SSHPiper 部署完成。"
else
    log_error "未找到 $BASE_DIR/sshpiper/deploy.yml，SSHPiper 部署失败。"
    exit 1
fi

# ==============================================================================
# 总结
# ==============================================================================
echo ""
echo "=========================================="
log_info "所有组件部署任务已提交！"
echo "=========================================="
echo "验证命令："
echo "  kubectl get all -n volcano-system"
echo "  kubectl get pods -n kubeflow"
echo "  kubectl get pods -n apisix"
echo "  kubectl get pods -n kubeflow -l app=sshpiper"
echo "  kubectl get pipe -n kubeflow"
echo "=========================================="
