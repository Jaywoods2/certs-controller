package server

import (
	"github.com/Jaywoods/certs-controller/pkg/webhook/server/ingress/mutating"
	"github.com/operator-framework/operator-sdk/pkg/k8sutil"
	"k8s.io/apimachinery/pkg/types"
	"os"
	"sigs.k8s.io/controller-runtime/pkg/manager"
	logf "sigs.k8s.io/controller-runtime/pkg/runtime/log"
	"sigs.k8s.io/controller-runtime/pkg/webhook"
)

var (
	log = logf.Log.WithName("webhook_server")
	// HandlerMap contains all admission webhook handlers.
)

func Add(mgr manager.Manager) error {
	ns, err := k8sutil.GetOperatorNamespace()
	if err != nil {
		ns = "test"
	}
	secretName := os.Getenv("SECRET_NAME")
	if len(secretName) == 0 {
		secretName = "cert-controller-webhook-secret"
	}
	bootstrapOptions := &webhook.BootstrapOptions{
		MutatingWebhookConfigName: "cert-controller-mutating-webhook",
		//ValidatingWebhookConfigName: "cert-controller-validating-webhook",
	}
	bootstrapOptions.Service = &webhook.Service{
		Name:      "cert-controller-server",
		Namespace: ns,
		Selectors: map[string]string{
			"name": "certs-controller",
		},
	}
	//var host string = "10.72.42.5"
	//bootstrapOptions.Host = &host
	bootstrapOptions.Secret = &types.NamespacedName{
		Namespace: ns,
		Name:      secretName,
	}

	var webhookPort int32 = 9876
	svr, err := webhook.NewServer("cert-controller-admission-server", mgr, webhook.ServerOptions{
		Port:             webhookPort,
		CertDir:          "/tmp/k8s-webhook-server/serving-certs",
		BootstrapOptions: bootstrapOptions,
	})
	if err != nil {
		return err
	}
	wh, err := mutating.IngressBuilder.Handlers(&mutating.IngressCreateHandle{}).WithManager(mgr).Build()
	if err != nil {
		log.Info("webhook create failed")
		return err
	}
	return svr.Register(wh)
}
