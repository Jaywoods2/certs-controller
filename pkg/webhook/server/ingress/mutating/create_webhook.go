package mutating

import (
	admissionregistrationv1beta1 "k8s.io/api/admissionregistration/v1beta1"
	"k8s.io/api/extensions/v1beta1"
	"sigs.k8s.io/controller-runtime/pkg/webhook/admission/builder"
)

func init() {
	IngressBuilder = builder.NewWebhookBuilder().
		Name("mutating-create-ingress.pakchoi.top").
		Path("/mutating-create-ingress").
		Mutating().
		Operations(admissionregistrationv1beta1.Create).
		FailurePolicy(admissionregistrationv1beta1.Fail).
		ForType(&v1beta1.Ingress{})
}
