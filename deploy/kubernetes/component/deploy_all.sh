#!/bin/bash

# ==============================================================================
# Neptune 一键部署脚本
# 包含组件: Volcano, Kubeflow, APISIX, SSHPiper
# 支持重复运行（幂等）
# ==============================================================================

set -euo pipefail

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
HELM_CACHE_DIR="${HELM_CACHE_DIR:-$BASE_DIR/.helm-cache}"
HELM_RETRY_COUNT="${HELM_RETRY_COUNT:-3}"
HELM_RETRY_DELAY_SECONDS="${HELM_RETRY_DELAY_SECONDS:-5}"
APPLY_RETRY_COUNT="${APPLY_RETRY_COUNT:-5}"
APPLY_RETRY_DELAY_SECONDS="${APPLY_RETRY_DELAY_SECONDS:-10}"
CERT_MANAGER_WEBHOOK_TIMEOUT_SECONDS="${CERT_MANAGER_WEBHOOK_TIMEOUT_SECONDS:-180}"
VOLCANO_CHART_VERSION="${VOLCANO_CHART_VERSION:-1.14.1}"
VOLCANO_CHART_FILE="${VOLCANO_CHART_FILE:-}"
APISIX_CHART_VERSION="${APISIX_CHART_VERSION:-2.13.0}"
APISIX_CHART_FILE="${APISIX_CHART_FILE:-}"
APISIX_ADMIN_KEY="${APISIX_ADMIN_KEY:-edd1c9f034335f136f87ad84b625c8f1}"
KUBEFLOW_MANIFESTS_REPO="${KUBEFLOW_MANIFESTS_REPO:-https://github.com/kubeflow/manifests.git}"
KUBEFLOW_MANIFESTS_REF="${KUBEFLOW_MANIFESTS_REF:-v1.10.0}"
SSHPIPER_CRD_URL="${SSHPIPER_CRD_URL:-https://raw.githubusercontent.com/tg123/sshpiper/7ce7b52e6a71f167ee78fd439a19d016e610d1d2/plugin/kubernetes/crd.yaml}"
SSHPIPER_SERVER_KEY_FILE="${SSHPIPER_SERVER_KEY_FILE:-$BASE_DIR/sshpiper/ssh_host_ed25519_key}"

function ensure_namespace() {
    local namespace="$1"
    if kubectl get namespace "$namespace" >/dev/null 2>&1; then
        log_info "namespace ${namespace} 已存在，继续复用。"
    else
        log_info "正在创建 namespace ${namespace} ..."
        kubectl create namespace "$namespace" >/dev/null
    fi
}

function ensure_sshpiper_secret() {
    if kubectl get secret sshpiper-server-key -n kubeflow >/dev/null 2>&1; then
        log_info "已检测到 kubeflow/sshpiper-server-key。"
        return 0
    fi

    if [ -f "$SSHPIPER_SERVER_KEY_FILE" ]; then
        if [ "$SSHPIPER_SERVER_KEY_FILE" = "$BASE_DIR/sshpiper/ssh_host_ed25519_key" ]; then
            log_warn "未检测到 kubeflow/sshpiper-server-key，将使用仓库内默认私钥自动创建。该私钥只建议用于开发测试环境。"
        else
            log_info "未检测到 kubeflow/sshpiper-server-key，将使用指定私钥文件自动创建。"
        fi

        kubectl create secret generic sshpiper-server-key \
          -n kubeflow \
          --from-file=server_key="$SSHPIPER_SERVER_KEY_FILE" >/dev/null
        log_info "kubeflow/sshpiper-server-key 已自动创建。"
        return 0
    fi

    log_error "缺少 kubeflow/sshpiper-server-key，且未找到可用私钥文件。"
    echo "请先创建服务端私钥 Secret，例如："
    echo "  kubectl create secret generic sshpiper-server-key -n kubeflow --from-file=server_key=$BASE_DIR/sshpiper/ssh_host_ed25519_key"
    echo "或者在执行脚本前指定自定义私钥文件："
    echo "  SSHPIPER_SERVER_KEY_FILE=/path/to/ssh_host_ed25519_key ./deploy_all.sh"
    exit 1
}

function helm_release_exists() {
    local release_name="$1"
    local namespace="$2"
    helm status "$release_name" -n "$namespace" >/dev/null 2>&1
}

function log_helm_release_state() {
    local release_name="$1"
    local namespace="$2"

    if helm_release_exists "$release_name" "$namespace"; then
        log_info "检测到 Helm Release ${release_name} (${namespace}) 已存在，本次会执行 upgrade。"
    else
        log_info "未检测到 Helm Release ${release_name} (${namespace})，本次会执行 install。"
    fi
}

function log_deployment_apply_state() {
    local description="$1"
    local namespace="$2"
    local deployment_name="$3"

    if kubectl get deployment "$deployment_name" -n "$namespace" >/dev/null 2>&1; then
        log_info "${description} 已存在，本次会执行 apply/update。"
    else
        log_info "${description} 不存在，本次会创建。"
    fi
}

function resolve_chart_package() {
    local chart_ref="$1"
    local chart_version="$2"
    local override_chart_file="$3"
    local chart_name="${chart_ref##*/}"
    local cached_chart_file="$HELM_CACHE_DIR/${chart_name}-${chart_version}.tgz"
    local attempt

    mkdir -p "$HELM_CACHE_DIR"

    if [ -n "$override_chart_file" ]; then
        if [ -f "$override_chart_file" ]; then
            log_info "使用本地 Helm Chart 包: $override_chart_file"
            RESOLVED_CHART_PACKAGE="$override_chart_file"
            return 0
        fi

        log_error "指定的 Helm Chart 包不存在: $override_chart_file"
        exit 1
    fi

    if [ -f "$cached_chart_file" ]; then
        log_info "使用本地缓存 Helm Chart 包: $cached_chart_file"
        RESOLVED_CHART_PACKAGE="$cached_chart_file"
        return 0
    fi

    for attempt in $(seq 1 "$HELM_RETRY_COUNT"); do
        log_info "正在下载 Helm Chart ${chart_ref}:${chart_version} (第 ${attempt}/${HELM_RETRY_COUNT} 次)..."
        rm -f "$cached_chart_file"

        if helm pull "$chart_ref" --version "$chart_version" --destination "$HELM_CACHE_DIR" >/dev/null 2>&1 && [ -f "$cached_chart_file" ]; then
            log_info "Helm Chart 下载完成: $cached_chart_file"
            RESOLVED_CHART_PACKAGE="$cached_chart_file"
            return 0
        fi

        if [ "$attempt" -lt "$HELM_RETRY_COUNT" ]; then
            log_warn "下载 ${chart_ref}:${chart_version} 失败，${HELM_RETRY_DELAY_SECONDS}s 后重试。"
            sleep "$HELM_RETRY_DELAY_SECONDS"
        fi
    done

    log_error "下载 Helm Chart ${chart_ref}:${chart_version} 失败。"
    echo "可选做法："
    echo "  1. 稍后重试脚本"
    echo "  2. 在浏览器中手动下载 tgz 包后，使用本地文件执行"
    echo "     VOLCANO_CHART_FILE=/path/to/volcano-${VOLCANO_CHART_VERSION}.tgz ./deploy_all.sh"
    echo "     APISIX_CHART_FILE=/path/to/apisix-${APISIX_CHART_VERSION}.tgz ./deploy_all.sh"
    exit 1
}

function helm_upgrade_install_with_retry() {
    local release_name="$1"
    local namespace="$2"
    local chart_ref="$3"
    local chart_version="$4"
    local override_chart_file="$5"
    local chart_package

    resolve_chart_package "$chart_ref" "$chart_version" "$override_chart_file"
    chart_package="$RESOLVED_CHART_PACKAGE"
    log_helm_release_state "$release_name" "$namespace"
    shift 5
    helm upgrade --install "$release_name" "$chart_package" "$@"
}

function wait_for_cert_manager_webhook_ready() {
    local timeout_seconds="$1"
    local interval_seconds=5
    local elapsed=0
    local ca_bundle=""

    log_info "等待 cert-manager webhook CA 注入完成..."

    while [ "$elapsed" -lt "$timeout_seconds" ]; do
        ca_bundle="$(kubectl get validatingwebhookconfiguration cert-manager-webhook -o jsonpath='{.webhooks[0].clientConfig.caBundle}' 2>/dev/null || true)"
        if [ -n "$ca_bundle" ]; then
            log_info "cert-manager webhook CA bundle 已就绪。"
            return 0
        fi

        sleep "$interval_seconds"
        elapsed=$((elapsed + interval_seconds))
    done

    log_error "等待 cert-manager webhook CA bundle 超时 (${timeout_seconds}s)。"
    exit 1
}

function apply_kustomize_with_retry() {
    local description="$1"
    local path="$2"
    local attempt

    for attempt in $(seq 1 "$APPLY_RETRY_COUNT"); do
        if kubectl apply -k "$path"; then
            return 0
        fi

        if [ "$attempt" -lt "$APPLY_RETRY_COUNT" ]; then
            log_warn "${description} 应用失败，第 ${attempt}/${APPLY_RETRY_COUNT} 次重试后等待 ${APPLY_RETRY_DELAY_SECONDS}s。"
            sleep "$APPLY_RETRY_DELAY_SECONDS"
        fi
    done

    log_error "${description} 应用失败，已达到最大重试次数。"
    exit 1
}

function resolve_kubeflow_component_paths() {
    if [ -d "$MANIFESTS_DIR/apps" ]; then
        APPS_DIR="$MANIFESTS_DIR/apps"
    elif [ -d "$MANIFESTS_DIR/applications" ]; then
        APPS_DIR="$MANIFESTS_DIR/applications"
    else
        log_error "未找到 Kubeflow apps/applications 目录，请检查 manifests 版本。"
        exit 1
    fi

    if [ -d "$MANIFESTS_DIR/common/istio/istio-crds/base" ]; then
        ISTIO_CRDS_DIR="$MANIFESTS_DIR/common/istio/istio-crds/base"
    elif [ -d "$MANIFESTS_DIR/common/istio-1-24/istio-crds/base" ]; then
        ISTIO_CRDS_DIR="$MANIFESTS_DIR/common/istio-1-24/istio-crds/base"
    elif [ -d "$MANIFESTS_DIR/common/istio-cni-1-24/istio-crds/base" ]; then
        ISTIO_CRDS_DIR="$MANIFESTS_DIR/common/istio-cni-1-24/istio-crds/base"
    else
        log_error "未找到 Istio CRDs 目录，请检查 manifests 版本。"
        exit 1
    fi

    NOTEBOOK_CONTROLLER_DIR="$APPS_DIR/jupyter/notebook-controller/upstream/overlays/kubeflow"
    TENSORBOARD_CONTROLLER_DIR="$APPS_DIR/tensorboard/tensorboard-controller/upstream/overlays/kubeflow"

    for path in \
        "$NOTEBOOK_CONTROLLER_DIR" \
        "$TENSORBOARD_CONTROLLER_DIR" \
        "$ISTIO_CRDS_DIR" \
        "$MANIFESTS_DIR/common/cert-manager/base" \
        "$MANIFESTS_DIR/common/cert-manager/kubeflow-issuer/base"; do
        if [ ! -d "$path" ]; then
            log_error "缺少 Kubeflow 组件目录: $path"
            exit 1
        fi
    done
}

function cleanup_legacy_apisix_gatewayproxy_conflict() {
    local legacy_name="apisix"
    local managed_name="apisix-config"
    local service_name=""
    local service_port=""
    local admin_key=""

    if ! kubectl get gatewayproxy "$legacy_name" -n apisix >/dev/null 2>&1; then
        return 0
    fi

    service_name="$(kubectl get gatewayproxy "$legacy_name" -n apisix -o jsonpath='{.spec.provider.controlPlane.service.name}' 2>/dev/null || true)"
    service_port="$(kubectl get gatewayproxy "$legacy_name" -n apisix -o jsonpath='{.spec.provider.controlPlane.service.port}' 2>/dev/null || true)"
    admin_key="$(kubectl get gatewayproxy "$legacy_name" -n apisix -o jsonpath='{.spec.provider.controlPlane.auth.adminKey.value}' 2>/dev/null || true)"

    if [ "$service_name" = "apisix-admin" ] && [ "$service_port" = "9180" ] && [ "$admin_key" = "$APISIX_ADMIN_KEY" ]; then
        log_warn "检测到旧版 GatewayProxy apisix/apisix 与 Helm 管理的 apisix-config 冲突，正在自动清理旧资源。"
        kubectl delete gatewayproxy "$legacy_name" -n apisix >/dev/null
        return 0
    fi

    log_error "检测到自定义 GatewayProxy apisix/apisix，与 Helm 管理的 apisix-config 存在潜在冲突，脚本不会自动删除。"
    echo "请先人工确认以下资源是否可删除："
    echo "  kubectl get gatewayproxy $legacy_name -n apisix -o yaml"
    echo "如果确认它是旧版默认配置，可执行："
    echo "  kubectl delete gatewayproxy $legacy_name -n apisix"
    exit 1
}

function fix_apisix_proxy_mode_config() {
    local config_content=""
    local fixed_content=""

    config_content="$(kubectl get configmap apisix -n apisix -o jsonpath='{.data.config\.yaml}' 2>/dev/null || true)"
    if [ -z "$config_content" ]; then
        log_error "未找到 apisix/configmap，无法检查 proxy_mode 配置。"
        exit 1
    fi

    if printf '%s' "$config_content" | grep -Eq '^[[:space:]]*proxy_mode:[[:space:]]*"http&stream"[[:space:]]*$'; then
        log_info "proxy_mode 已正确设置为独立的 http&stream 行。"
        return 0
    fi

    log_warn "检测到 APISIX proxy_mode 配置异常，正在自动修复..."
    fixed_content="$(
        printf '%s' "$config_content" \
          | perl -0pe 's/^(\s*proxy_mode:\s*\"?(?:http&stream|http)\"?)\s+stream_proxy:/$1\n  stream_proxy:/m; s/^(\s*proxy_mode:)\s*http\s*$/$1 "http\&stream"/m'
    )"

    printf '%s' "$fixed_content" \
      | kubectl create configmap apisix -n apisix --from-file=config.yaml=/dev/stdin --dry-run=client -o yaml \
      | kubectl apply -f -

    kubectl rollout restart deployment apisix -n apisix
    kubectl wait --for=condition=Available deployment/apisix -n apisix --timeout=120s 2>/dev/null || true
}

# ==============================================================================
# 1. 部署 Volcano
# ==============================================================================
log_info "正在部署 Volcano..."
helm repo add volcano-sh https://volcano-sh.github.io/helm-charts --force-update
helm repo update
helm_upgrade_install_with_retry "volcano" "volcano-system" "volcano-sh/volcano" "$VOLCANO_CHART_VERSION" "$VOLCANO_CHART_FILE" \
  -n volcano-system \
  --create-namespace
log_info "Volcano 部署完成。"

# ==============================================================================
# 2. 部署 Kubeflow 组件
# ==============================================================================
log_info "正在部署 Kubeflow 组件..."
cd "$BASE_DIR/kubeflow"
ensure_namespace kubeflow

if [ ! -d "manifests/.git" ]; then
    log_info "正在克隆 Kubeflow manifests 仓库..."
    if [ -e "manifests" ]; then
        log_error "$BASE_DIR/kubeflow/manifests 已存在但不是 git 仓库，请先手动处理后再执行脚本。"
        exit 1
    fi
    git clone "$KUBEFLOW_MANIFESTS_REPO" manifests
fi

MANIFESTS_DIR="$BASE_DIR/kubeflow/manifests"
log_info "固定 Kubeflow manifests 版本到 $KUBEFLOW_MANIFESTS_REF ..."
git -C "$MANIFESTS_DIR" fetch --tags origin "$KUBEFLOW_MANIFESTS_REF"
git -C "$MANIFESTS_DIR" checkout --detach FETCH_HEAD
resolve_kubeflow_component_paths

log_info "1/4 安装 Cert Manager..."
log_deployment_apply_state "Cert Manager" "cert-manager" "cert-manager"
kubectl apply -k "$MANIFESTS_DIR/common/cert-manager/base"
log_info "等待 Cert Manager CRD 注册..."
kubectl wait --for=condition=Established crd/certificates.cert-manager.io --timeout=120s 2>/dev/null || true
kubectl wait --for=condition=Established crd/issuers.cert-manager.io --timeout=120s 2>/dev/null || true
kubectl wait --for=condition=Established crd/clusterissuers.cert-manager.io --timeout=120s 2>/dev/null || true
log_info "等待 Cert Manager Deployment 就绪..."
kubectl wait --for=condition=Available deployment/cert-manager -n cert-manager --timeout=120s 2>/dev/null || true
kubectl wait --for=condition=Available deployment/cert-manager-cainjector -n cert-manager --timeout=120s 2>/dev/null || true
kubectl wait --for=condition=Available deployment/cert-manager-webhook -n cert-manager --timeout=120s 2>/dev/null || true
wait_for_cert_manager_webhook_ready "$CERT_MANAGER_WEBHOOK_TIMEOUT_SECONDS"

log_info "2/4 安装 Kubeflow Issuer..."
apply_kustomize_with_retry "Kubeflow Issuer" "$MANIFESTS_DIR/common/cert-manager/kubeflow-issuer/base"

log_info "3/4 安装 Istio CRDs..."
kubectl apply -k "$ISTIO_CRDS_DIR"

log_info "4/4 安装 Notebook & Tensorboard 控制器..."
log_deployment_apply_state "Notebook Controller" "kubeflow" "notebook-controller-deployment"
kubectl apply -k "$NOTEBOOK_CONTROLLER_DIR"
log_deployment_apply_state "TensorBoard Controller" "kubeflow" "tensorboard-controller-deployment"
kubectl apply -k "$TENSORBOARD_CONTROLLER_DIR"

log_info "等待 Notebook Controller 就绪..."
kubectl wait --for=condition=Available deployment/notebook-controller-deployment -n kubeflow --timeout=180s 2>/dev/null || true
log_info "等待 TensorBoard Controller 就绪..."
kubectl wait --for=condition=Available deployment/tensorboard-controller-deployment -n kubeflow --timeout=180s 2>/dev/null || true

log_info "Kubeflow 组件部署完成。"

# ==============================================================================
# 3. 部署 APISIX
# ==============================================================================
log_info "正在部署 APISIX..."
helm repo add apisix https://apache.github.io/apisix-helm-chart --force-update
helm repo update
cleanup_legacy_apisix_gatewayproxy_conflict

helm_upgrade_install_with_retry "apisix" "apisix" "apisix/apisix" "$APISIX_CHART_VERSION" "$APISIX_CHART_FILE" \
  --create-namespace \
  --namespace apisix \
  --set ingress-controller.enabled=true \
  --set ingress-controller.gatewayProxy.createDefault=true \
  --set ingress-controller.gatewayProxy.provider.controlPlane.service.name=apisix-admin \
  --set ingress-controller.gatewayProxy.provider.controlPlane.service.port=9180 \
  --set ingress-controller.gatewayProxy.provider.controlPlane.auth.adminKey.value="$APISIX_ADMIN_KEY" \
  --set ingress-controller.config.apisix.serviceNamespace=apisix \
  --set 'apisix.proxy_mode=http&stream' \
  --set "apisix.stream_proxy.tcp[0]=9100"

log_info "等待 APISIX 就绪..."
kubectl wait --for=condition=Available deployment/apisix -n apisix --timeout=180s 2>/dev/null || true

# 修复 proxy_mode，避免 YAML 被错误合并成单行
log_info "检查并修复 proxy_mode 配置..."
fix_apisix_proxy_mode_config

# 添加 SSH 端口
log_info "为 APISIX Gateway 添加 SSH 端口..."
kubectl patch svc apisix-gateway -n apisix --type='json' \
  -p='[{"op":"add","path":"/spec/ports/-","value":{"name":"ssh","port":22,"targetPort":9100,"protocol":"TCP"}}]' 2>/dev/null || log_warn "SSH 端口可能已存在，跳过。"

log_info "APISIX 部署完成。"

# ==============================================================================
# 4. 部署 SSHPiper
# ==============================================================================
log_info "正在部署 SSHPiper..."
ensure_sshpiper_secret

log_info "安装 SSHPiper CRD..."
kubectl apply -f "$SSHPIPER_CRD_URL"
log_info "等待 CRD 注册..."
kubectl wait --for=condition=Established crd/pipes.sshpiper.com --timeout=60s 2>/dev/null || true

if [ -f "$BASE_DIR/sshpiper/deploy.yml" ]; then
    log_deployment_apply_state "SSHPiper" "kubeflow" "sshpiper"
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
