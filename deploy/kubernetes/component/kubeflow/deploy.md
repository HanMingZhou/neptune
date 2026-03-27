# Download kubeflow Repo
```bash
git clone https://github.com/kubeflow/manifests.git
```

# Install kubeflow components
## Install Cert Manager and Kubeflow-Issuer
```bash
kubectl apply -k manifests/common/cert-manager/base

# 等待 CRD 和 Deployment 就绪
kubectl wait --for=condition=Established crd/certificates.cert-manager.io --timeout=120s
kubectl wait --for=condition=Available deployment/cert-manager -n cert-manager --timeout=120s
kubectl wait --for=condition=Available deployment/cert-manager-webhook -n cert-manager --timeout=120s

kubectl apply -k manifests/common/cert-manager/kubeflow-issuer/base
```

## Install Istio CRDs
```bash
kubectl apply -k manifests/common/istio/istio-crds/base
```

## Install Notebook Controller and Tensorboard Controller
```bash
kubectl apply -k manifests/applications/jupyter/notebook-controller/upstream/overlays/kubeflow
kubectl apply -k manifests/applications/tensorboard/tensorboard-controller/upstream/overlays/kubeflow

# 等待 Controller 就绪
kubectl wait --for=condition=Available deployment/notebook-controller-deployment -n kubeflow --timeout=180s
kubectl wait --for=condition=Available deployment/tensorboard-controller-deployment -n kubeflow --timeout=180s
```