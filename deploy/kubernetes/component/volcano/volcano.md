# 安装最新版本volcano
## 通过 Helm 安装
```bash
helm repo add volcano-sh https://volcano-sh.github.io/helm-charts

helm repo update

helm install volcano volcano-sh/volcano -n volcano-system --create-namespace
```
## 验证 Volcano 组件的状态
```bash
kubectl get all -n volcano-system
```
