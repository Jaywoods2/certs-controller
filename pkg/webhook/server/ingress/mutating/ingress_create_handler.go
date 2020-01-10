package mutating

import (
	"context"
	"github.com/Jaywoods/certs-controller/pkg/apis/app/v1alpha1"
	"k8s.io/api/extensions/v1beta1"
	"net/http"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/runtime/inject"
	logf "sigs.k8s.io/controller-runtime/pkg/runtime/log"
	"sigs.k8s.io/controller-runtime/pkg/webhook/admission"
	"sigs.k8s.io/controller-runtime/pkg/webhook/admission/types"
)

var log = logf.Log.WithName("ingress_create_handler")

type IngressCreateHandle struct {
	Client  client.Client
	Decoder types.Decoder
}

var _ admission.Handler = &IngressCreateHandle{}

func (i *IngressCreateHandle) Handle(ctx context.Context, req types.Request) types.Response {
	obj := &v1beta1.Ingress{}
	err := i.Decoder.Decode(req, obj)
	if err != nil {
		return admission.ErrorResponse(http.StatusBadRequest, err)
	}
	copy := obj.DeepCopy()
	log.Info("接收到Ingress事件", copy.Namespace, copy.Name)
	// todo: 处理对象
	if v, ok := copy.Annotations["pakchoi.top/inject-cert"]; ok {
		if v == "true" {
			certSecrets := &v1alpha1.CertSecretList{}
			if err := i.Client.List(context.Background(), &client.ListOptions{}, certSecrets); err != nil {
				log.Error(err, "查询CRD:CertSecret 列表失败")
			}

			var tlsMaps = make(map[string][]string)
			var secretNames []string
			for _, rule := range copy.Spec.Rules {
				suffixHost := split(rule.Host)
				for _, cs := range certSecrets.Items {
					for _, tls := range cs.Spec.Tls {
						if suffixHost == tls.Name {
							if v, ok := tlsMaps[tls.Name]; ok {
								v = append(v, rule.Host)
								tlsMaps[tls.Name] = v
							} else {
								tlsMaps[tls.Name] = []string{rule.Host}
							}
							secretNames = append(secretNames, tls.Name)
						}
					}
				}
			}
			var tls []v1beta1.IngressTLS
			for _, sn := range secretNames {
				ingressTls := v1beta1.IngressTLS{
					Hosts:      tlsMaps[sn],
					SecretName: sn,
				}
				tls = append(tls, ingressTls)
			}
			copy.Spec.TLS = tls
			log.Info("注入TLS配置", copy.Namespace, copy.Name)
		} else {
			log.Info("Ingress Annotation值不正确")
		}

	} else {
		log.Info("正常接收Ingress对象，不作处理", copy.Namespace, copy.Name)
	}
	return admission.PatchResponse(obj, copy)
}

var _ inject.Client = &IngressCreateHandle{}

// InjectClient injects the client into the IngressCreateHandle
func (i *IngressCreateHandle) InjectClient(c client.Client) error {
	i.Client = c
	return nil
}

var _ inject.Decoder = &IngressCreateHandle{}

// InjectDecoder injects the decoder into the IngressCreateHandle
func (i *IngressCreateHandle) InjectDecoder(d types.Decoder) error {
	i.Decoder = d
	return nil
}

func split(str string) string {
	var index int
	for i, x := range str {
		if string(x) == "." {
			index = i + 1
			break
		}
	}
	return str[index:]
}
