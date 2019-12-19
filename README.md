> 用于管理K8S个namespaes证书的创建、更新

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
  scope: "Namespaced" # Namespaced 或 Cluster
  namespaces:
    - default
    - kube-system
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