[![](https://img.shields.io/badge/go-1.12.9-brightgreen.svg)](https://img.shields.io/badge/go-1.12.9-brightgreen.svg)
[![](https://img.shields.io/badge/OperatorSdk-v0.13.1-yellow.svg)](https://img.shields.io/badge/OperatorSdk-v0.13.1-yellow.svg)
[![](https://img.shields.io/badge/ClientGo-v11.0.0-orange.svg)](https://img.shields.io/badge/ClientGo-v11.0.0-orange.svg)
[![](https://img.shields.io/badge/mode-HA-blue.svg)](https://img.shields.io/badge/mode-HA-blue.svg)

## 功能

- 根据CertSecret配置创建更新Namespace下的Secret
- 监听Namespace事件创建该Ns下的Secret
- 创建Ingress时注入Tls配置

## 部署

1、 创建CSR请求证书，用于ApiServer与该组件Webhook交互。

示例部署到`test`Namespace下，登录到Master节点执行如下命令

```shell script
# 创建namespace
kubectl create ns test
# 该命令用于创建CSR请求签名Service地址，创建证书Secret对象
sh deploy/webhook/webhook-create-signed-cert.sh --service cert-controller-server --secret cert-controller-webhook-secret --namespace test
```

2、 创建CRD对象及相关角色权限配置。
> 由于组件在运行时会自动创建MutatingWebhookConfiguration对象，因此授予cluster-admin角色

```shell script
kubectl create -f deploy/service_account.yaml
kubectl create -f deploy/role_binding.yaml
kubectl create -f deploy/crds/app_v1alpha1_certsecret_crd.yaml
kubectl create -f deploy/operator.yaml
```
3、 配置CR对象

> 注意: tls.name 与证书签名根域名保持一致。

```yaml
apiVersion: app.pakchoi.top/v1alpha1
kind: CertSecret
metadata:
  name: my-certsecret
spec:
  tls:
    - name: pakchoi.top
      key: "-----BEGIN RSA PRIVATE KEY-----
            XANKSJHDKWS填写证书内容NKDUASDNA
            -----END RSA PRIVATE KEY-----"
      crt: "-----BEGIN CERTIFICATE-----
            XANKSJHDKWS填写证书内容NKDUASDNA
            -----END CERTIFICATE-----"
    - name: saas.pakchoi.top
      key: "-----BEGIN RSA PRIVATE KEY-----
            XANKSJHDKWS填写证书内容NKDUASDNA
            -----END RSA PRIVATE KEY-----"
      crt: "-----BEGIN CERTIFICATE-----
            XANKSJHDKWS填写证书内容NKDUASDNA2
            -----END CERTIFICATE-----"
  scope: "Cluster" # Namespaced 或 Cluster
# Namespaced时填写
  namespaces: []
  cascade: true # 测试时，删除cr级联删除关联的secret
```

应用CR配置,跟踪cert-controller的日志。

```shell script
kubectl create -f deploy/crds/app_v1alpha1_certsecret_cr.yaml -n test
```

查看Secret
```shell script
kubectl get secret --all-namespaces | grep pakchoi.top
```

## 测试

1、 创建Namespace `test2`,跟踪日志。
```shell script
kubectl create ns test2
```
查看secret
```shell script
kubectl  get secret -n test2
```

2、 创建Ingress,并添加注释`pakchoi.top/inject-cert: "true"`
```shell script
cat <<EOF | kubectl create -f -
apiVersion: extensions/v1beta1
kind: Ingress
metadata:
  name: mynginx
  annotations:
    pakchoi.top/inject-cert: "true"
spec:
  rules:
  - host: mynginx.pakchoi.top
    http:
      paths:
      - backend:
          serviceName: mynginx
          servicePort: 80
EOF
```

检查Ingress注入了Tls证书配置。

## Todo

- 对比secret内容判断更新