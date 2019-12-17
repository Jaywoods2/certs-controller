package resources

import (
	"github.com/Jaywoods/certs-controller/pkg/apis/app/v1alpha1"
	corev1 "k8s.io/api/core/v1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"
)

//func NewSecretList(cs *v1alpha1.CertSecret) *corev1.SecretList {
//	var items = make([]corev1.Secret, len(cs.Spec.Tls))
//	for _, d := range cs.Spec.Tls {
//		items = append(items, NewSecret(cs, d))
//	}
//	return &corev1.SecretList{
//		TypeMeta: v1.TypeMeta{
//			Kind:       "List",
//			APIVersion: "v1",
//		},
//		Items:    items,
//	}
//}

func NewSecret(cs *v1alpha1.CertSecret, data v1alpha1.TlsData, ns string) *corev1.Secret {
	return &corev1.Secret{
		TypeMeta: v1.TypeMeta{
			Kind:       "Secret",
			APIVersion: "v1",
		},
		ObjectMeta: v1.ObjectMeta{
			Name:      data.Name,
			Namespace: ns,
			OwnerReferences: []v1.OwnerReference{
				*v1.NewControllerRef(cs, schema.GroupVersionKind{
					Group:   v1alpha1.SchemeGroupVersion.Group,
					Version: v1alpha1.SchemeGroupVersion.Version,
					Kind:    "CertSecret",
				}),
			},
		},
		StringData: newTlsData(data),
		Type:       "kubernetes.io/tls",
	}
}

func newTlsData(data v1alpha1.TlsData) map[string]string {
	sd := make(map[string]string)
	sd["tls.key"] = data.Key
	sd["tls.crt"] = data.Crt
	return sd
}
