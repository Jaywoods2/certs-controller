> 用于管理K8S各namespaes证书秘钥的创建、更新

## 参数

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

## 部署

```shell script
kubectl create -f deploy/service_account.yaml
kubectl create -f deploy/role.yaml
kubectl create -f deploy/role_binding.yaml
kubectl create -f deploy/crds/app_v1alpha1_certsecret_crd.yaml
kubectl create -f deploy/operator.yaml
kubectl create -f deploy/crds/app_v1alpha1_certsecret_cr.yaml
```

## Todo

- 对比secret内容判断更新