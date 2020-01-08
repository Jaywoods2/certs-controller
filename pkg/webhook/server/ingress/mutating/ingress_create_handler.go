package mutating

import (
	"context"
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
	log.Info("接收到Ingress创建事件", copy.Namespace, copy.Name)
	// todo: 处理对象
	return admission.PatchResponse(obj, copy)
}

var _ inject.Client = &IngressCreateHandle{}

// InjectClient injects the client into the PodCreateHandler
func (h *IngressCreateHandle) InjectClient(c client.Client) error {
	h.Client = c
	return nil
}

var _ inject.Decoder = &IngressCreateHandle{}

// InjectDecoder injects the decoder into the PodCreateHandler
func (h *IngressCreateHandle) InjectDecoder(d types.Decoder) error {
	h.Decoder = d
	return nil
}
