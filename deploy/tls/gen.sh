rm -rf ca.pem ca-key.pem cert.pem cert.csr cert-key.pem ca.srl

openssl genrsa -out ca-key.pem 2048

openssl req -x509 -new -nodes \
      -days 36500 \
      -key ca-key.pem \
      -config kube-openssl.cnf \
      -subj "/CN=kubernetes" \
      -extensions v3_ca \
      -out ca.pem

openssl genrsa -out cert-key.pem 2048

openssl req -new -key cert-key.pem -subj "/CN=kubernetes" -out cert.csr

openssl x509 -req -CA ca.pem -CAkey ca-key.pem -days 36500 -in  cert.csr -CAcreateserial -extensions v3_req_peer -extfile kube-openssl.cnf -out cert.pem

caBundle=$(cat ca.pem | base64)

key=$(cat cert-key.pem | base64)

cert=$(cat cert.pem | base64)

sed -i "" "s/ca.crt:.*/ca.crt:\ $caBundle/g" ../webhook/secret.yaml
sed -i "" "s/tls.crt:.*/tls.crt:\ $cert/g" ../webhook/secret.yaml
sed -i "" "s/tls.key:.*/tls.key:\ $key/g" ../webhook/secret.yaml
sed -i "" "s/caBundle:.*/caBundle:\ $caBundle/g" ../webhook/webhookconfiguration.yaml

rm -rf /tmp/k8s-webhook-server/serving-certs/*

cp -rf ca.pem ca-key.pem cert.pem cert.csr cert-key.pem ca.srl /tmp/k8s-webhook-server/serving-certs

mv /tmp/k8s-webhook-server/serving-certs/cert-key.pem /tmp/k8s-webhook-server/serving-certs/key.pem

rm -rf ca.pem ca-key.pem cert.pem cert.csr cert-key.pem ca.srl

kubectl delete mutatingwebhookconfiguration  cert-controller-mutating-webhook

kubectl delete secret cert-controller-webhook-secret -n test