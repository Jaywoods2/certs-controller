package mutating

import "sigs.k8s.io/controller-runtime/pkg/webhook/admission/builder"

var (
	IngressBuilder = &builder.WebhookBuilder{}
)
