package resources

import (
	"github.com/Jaywoods/certs-controller/pkg/apis/app/v1alpha1"
	corev1 "k8s.io/api/core/v1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"
)

func NewSecret(cs *v1alpha1.CertSecret, data v1alpha1.TlsData, ns string) *corev1.Secret {
	var ownerReference []v1.OwnerReference
	if cs.Spec.Cascade {
		ownerReference = []v1.OwnerReference{
			*v1.NewControllerRef(cs, schema.GroupVersionKind{
				Group:   v1alpha1.SchemeGroupVersion.Group,
				Version: v1alpha1.SchemeGroupVersion.Version,
				Kind:    "CertSecret",
			})}
	}
	return &corev1.Secret{
		TypeMeta: v1.TypeMeta{
			Kind:       "Secret",
			APIVersion: "v1",
		},
		ObjectMeta: v1.ObjectMeta{
			Name:            data.Name,
			Namespace:       ns,
			OwnerReferences: ownerReference,
		},
		StringData: newTlsData(data),
		Type:       "kubernetes.io/tls",
	}
}

func UpdateSecret(data v1alpha1.TlsData, s *corev1.Secret) *corev1.Secret {
	s.StringData = newTlsData(data)
	return s
}

func newTlsData(data v1alpha1.TlsData) map[string]string {
	sd := make(map[string]string)
	sd["tls.key"] = data.Key
	sd["tls.crt"] = data.Crt
	return sd
}
