package namespace

import (
	"context"
	"fmt"
	"github.com/Jaywoods/certs-controller/pkg/apis/app/v1alpha1"
	"github.com/Jaywoods/certs-controller/pkg/resources"
	"github.com/operator-framework/operator-sdk/pkg/k8sutil"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	"sigs.k8s.io/controller-runtime/pkg/manager"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
	logf "sigs.k8s.io/controller-runtime/pkg/runtime/log"
	"sigs.k8s.io/controller-runtime/pkg/source"
)

var log = logf.Log.WithName("controller_namespace")

func Add(mgr manager.Manager) error {
	return add(mgr, newReconciler(mgr))
}

func newReconciler(mgr manager.Manager) reconcile.Reconciler {
	return &ReconcileNamespace{client: mgr.GetClient(), scheme: mgr.GetScheme()}
}

func add(mgr manager.Manager, r reconcile.Reconciler) error {
	c, err := controller.New("controller_namespace", mgr, controller.Options{Reconciler: r})
	if err != nil {
		return err
	}
	err = c.Watch(&source.Kind{
		Type: &corev1.Namespace{},
	}, &handler.EnqueueRequestForObject{})
	if err != nil {
		return err
	}
	return nil
}

type ReconcileNamespace struct {
	client client.Client
	scheme *runtime.Scheme
}

var _ reconcile.Reconciler = &ReconcileNamespace{}

func (r *ReconcileNamespace) Reconcile(request reconcile.Request) (reconcile.Result, error) {
	reqLogger := log.WithValues("Request.Name", request.Name)
	reqLogger.Info("Reconciling Namespace")
	ns := &corev1.Namespace{}
	nserr := r.client.Get(context.Background(), client.ObjectKey{
		Namespace: request.Namespace,
		Name:      request.Name,
	}, ns)
	if nserr != nil {
		if errors.IsNotFound(nserr) {
			reqLogger.Info(fmt.Sprintf("Namespace %s 已删除", request.Name))
			return reconcile.Result{}, nil
		}
		return reconcile.Result{}, nil
	}
	// 判断ns的状态
	if !(ns.Status.Phase == "Active") {
		reqLogger.Info(fmt.Sprintf("Namespace %s 状态为%s", request.Name, ns.Status.Phase))
		return reconcile.Result{}, nil
	}
	reqLogger.Info(fmt.Sprintf("Namespace %s 触发", request.Name))

	operatorNs, err := k8sutil.GetOperatorNamespace()
	if err != nil {
		return reconcile.Result{}, err
	}
	operatorName, err := k8sutil.GetOperatorName()
	if err != nil {
		return reconcile.Result{}, err
	}
	newNs := request.Name
	reqLogger.Info(fmt.Sprintf("operatorNs: %s |operatorName: %s ", operatorNs, operatorName))
	instance := &v1alpha1.CertSecret{}
	reqLogger.Info(fmt.Sprintf("operatorNs: %s |operatorName: %s ", operatorNs, operatorName))
	if err := r.client.Get(context.Background(), client.ObjectKey{
		Namespace: operatorNs,
		Name:      operatorName,
	}, instance); err != nil && errors.IsNotFound(err) {
		return reconcile.Result{}, nil
	}
	if instance.DeletionTimestamp != nil {
		return reconcile.Result{}, err
	}
	tls := instance.Spec.Tls
	for _, tl := range tls {
		// 判断该ns是否存在secret
		secret := &corev1.Secret{}
		if errs := r.client.Get(context.Background(), client.ObjectKey{
			Namespace: newNs,
			Name:      tl.Name,
		}, secret); errs != nil && errors.IsNotFound(errs) {
			reqLogger.Info(fmt.Sprintf("1|查询secret %s/%s 报错，不存在", newNs, tl.Name))
			// 不存在创建
			reqLogger.Info(fmt.Sprintf("2|创建secret:%s/%s", newNs, tl.Name))
			newSecret := resources.NewSecret(instance, tl, newNs)
			if errc := r.client.Create(context.Background(), newSecret); errc != nil {
				reqLogger.Error(errc, "创建secret失败")
				// 重新入队
				return reconcile.Result{}, errc
			}
		}

	}

	return reconcile.Result{}, nil
}
