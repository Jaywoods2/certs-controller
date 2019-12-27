package certsecret

import (
	"context"
	"fmt"
	"github.com/Jaywoods/certs-controller/pkg/apis/app/v1alpha1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/util/workqueue"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/event"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
	logf "sigs.k8s.io/controller-runtime/pkg/runtime/log"
)

var _ handler.EventHandler = &enqueueRequestForNs{}
var loge = logf.Log.WithName("namespace_event_hander")

type enqueueRequestForNs struct {
	client client.Client
}

func (n *enqueueRequestForNs) Create(evt event.CreateEvent, q workqueue.RateLimitingInterface) {
	ns, ok := evt.Object.(*corev1.Namespace)
	if !ok {
		return
	}
	loge.Info(fmt.Sprintf("Create Namespace: %s", ns.Name))
	n.addNamespace(q, evt.Object, ns)

}

func (n *enqueueRequestForNs) Delete(evt event.DeleteEvent, q workqueue.RateLimitingInterface) {
	ns, ok := evt.Object.(*corev1.Namespace)
	if !ok {
		return
	}
	loge.Info(fmt.Sprintf("Delete Namespace: %s", ns.Name))
}

func (n *enqueueRequestForNs) Generic(evt event.GenericEvent, q workqueue.RateLimitingInterface) {
	return
}

func (n *enqueueRequestForNs) Update(evt event.UpdateEvent, q workqueue.RateLimitingInterface) {
	return
}

func (n *enqueueRequestForNs) addNamespace(q workqueue.RateLimitingInterface, obj runtime.Object, ns *corev1.Namespace) {
	certSecrets := &v1alpha1.CertSecretList{}
	if err := n.client.List(context.Background(), &client.ListOptions{}, certSecrets); err != nil {
		loge.Error(err, "查询CRD:CertSecret 列表失败")
	}
	for _, certSecret := range certSecrets.Items {
		scope := certSecret.Spec.Scope
		switch scope {
		case ScopeCluster:
			loge.Info(fmt.Sprintf("触发入队 , CRD: CertSecret,Name: %s", certSecret.Name))
			q.Add(reconcile.Request{NamespacedName:
			types.NamespacedName{
				Name: certSecret.Name,
			}})
		case ScopeNamespaced:
			exist := false
			for _, n := range certSecret.Spec.Namespaces {
				if n == ns.Name {
					exist = true
					break
				}
			}
			if exist {
				loge.Info(fmt.Sprintf("触发入队 , CRD: CertSecret,Name: %s", certSecret.Name))
				q.Add(reconcile.Request{NamespacedName:
				types.NamespacedName{
					Name: certSecret.Name,
				}})
			}
			loge.Info(fmt.Sprintf("新增Namespace:%s 不在%s范围内", ns.Name, scope))

		default:
			loge.Info("scope not in [\"Cluster\",\"Namespaced\"]")
		}

	}
}

