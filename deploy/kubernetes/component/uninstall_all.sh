#!/bin/bash

# ==============================================================================
# Neptune 组件一键卸载脚本
# 默认卸载 deploy_all.sh 安装的组件:
#   SSHPiper, APISIX, Kubeflow 最小控制器, Volcano
# 默认不删除 kubeflow namespace，避免误删用户资源。
# 可选环境变量:
#   DELETE_COMPONENT_NAMESPACES=true  删除 apisix/cert-manager/volcano-system namespace
#   DELETE_KUBEFLOW_NAMESPACE=true    删除 kubeflow namespace
#   REMOVE_SSHPIPER_SECRET=true       删除 kubeflow/sshpiper-server-key（默认 true）
#   REMOVE_LOCAL_MANIFESTS=true       删除本地 kubeflow/manifests 目录（默认 false）
# ==============================================================================

set -e

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

for cmd in helm kubectl git; do
    if ! command -v "$cmd" >/dev/null 2>&1; then
        log_error "未找到 $cmd 工具，请先安装。"
        exit 1
    fi
done

BASE_DIR="$(cd "$(dirname "$0")" && pwd)"
KUBEFLOW_MANIFESTS_REPO="${KUBEFLOW_MANIFESTS_REPO:-https://github.com/kubeflow/manifests.git}"
KUBEFLOW_MANIFESTS_REF="${KUBEFLOW_MANIFESTS_REF:-v1.10.0}"
DELETE_COMPONENT_NAMESPACES="${DELETE_COMPONENT_NAMESPACES:-false}"
DELETE_KUBEFLOW_NAMESPACE="${DELETE_KUBEFLOW_NAMESPACE:-false}"
REMOVE_SSHPIPER_SECRET="${REMOVE_SSHPIPER_SECRET:-true}"
REMOVE_LOCAL_MANIFESTS="${REMOVE_LOCAL_MANIFESTS:-false}"

function ensure_manifests_checkout() {
    local manifests_root="$BASE_DIR/kubeflow/manifests"

    if [ ! -d "$manifests_root/.git" ]; then
        if [ -e "$manifests_root" ]; then
            log_error "$manifests_root 已存在但不是 git 仓库，请先手动处理后再执行卸载脚本。"
            exit 1
        fi

        log_info "未找到 Kubeflow manifests，本地正在拉取 $KUBEFLOW_MANIFESTS_REF ..."
        mkdir -p "$BASE_DIR/kubeflow"
        git clone "$KUBEFLOW_MANIFESTS_REPO" "$manifests_root"
    fi

    log_info "固定 Kubeflow manifests 版本到 $KUBEFLOW_MANIFESTS_REF ..."
    git -C "$manifests_root" fetch --tags origin "$KUBEFLOW_MANIFESTS_REF"
    git -C "$manifests_root" checkout --detach FETCH_HEAD >/dev/null
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

function delete_kustomize_path() {
    local description="$1"
    local path="$2"

    log_info "正在卸载 ${description}..."
    kubectl delete -k "$path" --ignore-not-found=true >/dev/null 2>&1 || log_warn "${description} 卸载时返回非零，请检查该目录下资源是否来自当前 manifests 版本。"
}

function uninstall_helm_release() {
    local release_name="$1"
    local namespace="$2"

    if helm status "$release_name" -n "$namespace" >/dev/null 2>&1; then
        log_info "正在卸载 Helm Release ${release_name} (${namespace}) ..."
        helm uninstall "$release_name" -n "$namespace" >/dev/null
    else
        log_warn "Helm Release ${release_name} (${namespace}) 不存在，跳过。"
    fi
}

function delete_namespace_if_requested() {
    local namespace="$1"

    if [ "$DELETE_COMPONENT_NAMESPACES" = "true" ]; then
        log_info "正在删除 namespace ${namespace} ..."
        kubectl delete namespace "$namespace" --ignore-not-found=true >/dev/null 2>&1 || log_warn "namespace ${namespace} 删除时返回非零，可能仍有资源在终止中。"
    else
        log_info "保留 namespace ${namespace}。如需一并删除，请设置 DELETE_COMPONENT_NAMESPACES=true。"
    fi
}

function delete_kubeflow_namespace_if_requested() {
    if [ "$DELETE_KUBEFLOW_NAMESPACE" = "true" ]; then
        log_warn "正在删除 kubeflow namespace，这会一并删除其中剩余资源。"
        kubectl delete namespace kubeflow --ignore-not-found=true >/dev/null 2>&1 || log_warn "kubeflow namespace 删除时返回非零，可能仍有资源在终止中。"
    else
        log_info "保留 kubeflow namespace。"
    fi
}

function remove_local_manifests_if_requested() {
    local manifests_root="$BASE_DIR/kubeflow/manifests"

    if [ "$REMOVE_LOCAL_MANIFESTS" = "true" ]; then
        if [ -d "$manifests_root" ]; then
            log_info "正在删除本地目录 $manifests_root ..."
            rm -rf "$manifests_root"
        else
            log_warn "本地目录 $manifests_root 不存在，跳过。"
        fi
    else
        log_info "保留本地目录 $manifests_root。"
    fi
}

# ==============================================================================
# 1. 卸载 SSHPiper
# ==============================================================================
log_info "正在卸载 SSHPiper ..."
if [ -f "$BASE_DIR/sshpiper/deploy.yml" ]; then
    kubectl delete -f "$BASE_DIR/sshpiper/deploy.yml" --ignore-not-found=true >/dev/null 2>&1 || log_warn "SSHPiper 基础资源删除时返回非零，可能资源已不存在。"
else
    log_warn "未找到 $BASE_DIR/sshpiper/deploy.yml，跳过 SSHPiper 基础资源删除。"
fi

kubectl delete pipes.sshpiper.com --all -A --ignore-not-found=true >/dev/null 2>&1 || log_warn "Pipe CR 清理返回非零，可能 CRD 已不存在。"
kubectl delete crd pipes.sshpiper.com --ignore-not-found=true >/dev/null 2>&1 || log_warn "SSHPiper CRD 删除返回非零，可能已不存在。"

if [ "$REMOVE_SSHPIPER_SECRET" = "true" ]; then
    log_info "正在删除 kubeflow/sshpiper-server-key ..."
    kubectl delete secret sshpiper-server-key -n kubeflow --ignore-not-found=true >/dev/null 2>&1 || log_warn "sshpiper-server-key 删除返回非零，可能已不存在。"
else
    log_info "保留 kubeflow/sshpiper-server-key。"
fi

# ==============================================================================
# 2. 卸载 APISIX
# ==============================================================================
uninstall_helm_release "apisix" "apisix"
kubectl delete gatewayproxy apisix -n apisix --ignore-not-found=true >/dev/null 2>&1 || log_warn "旧版 GatewayProxy apisix 删除返回非零，可能已不存在。"
kubectl delete gatewayproxy apisix-config -n apisix --ignore-not-found=true >/dev/null 2>&1 || log_warn "GatewayProxy apisix-config 删除返回非零，可能已不存在。"
delete_namespace_if_requested "apisix"

# ==============================================================================
# 3. 卸载 Kubeflow 最小控制器
# ==============================================================================
ensure_manifests_checkout

MANIFESTS_DIR="$BASE_DIR/kubeflow/manifests"
resolve_kubeflow_component_paths
delete_kustomize_path "TensorBoard Controller" "$TENSORBOARD_CONTROLLER_DIR"
delete_kustomize_path "Notebook Controller" "$NOTEBOOK_CONTROLLER_DIR"
delete_kustomize_path "Istio CRDs" "$ISTIO_CRDS_DIR"
delete_kustomize_path "Kubeflow Issuer" "$MANIFESTS_DIR/common/cert-manager/kubeflow-issuer/base"
delete_kustomize_path "Cert Manager" "$MANIFESTS_DIR/common/cert-manager/base"

delete_namespace_if_requested "cert-manager"
delete_kubeflow_namespace_if_requested
remove_local_manifests_if_requested

# ==============================================================================
# 4. 卸载 Volcano
# ==============================================================================
uninstall_helm_release "volcano" "volcano-system"
delete_namespace_if_requested "volcano-system"

# ==============================================================================
# 总结
# ==============================================================================
echo ""
echo "=========================================="
log_info "组件卸载命令已执行完成。"
echo "=========================================="
echo "验证命令："
echo "  helm ls -A"
echo "  kubectl get pods -n kubeflow"
echo "  kubectl get pods -n apisix"
echo "  kubectl get pods -n cert-manager"
echo "  kubectl get pods -n volcano-system"
echo "=========================================="
