package certsecret

import (
	"context"
	"encoding/json"
	"fmt"
	appv1alpha1 "github.com/Jaywoods/certs-controller/pkg/apis/app/v1alpha1"
	"github.com/Jaywoods/certs-controller/pkg/resources"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"reflect"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	"sigs.k8s.io/controller-runtime/pkg/manager"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
	logf "sigs.k8s.io/controller-runtime/pkg/runtime/log"
	"sigs.k8s.io/controller-runtime/pkg/source"
	"time"
)

var log = logf.Log.WithName("controller_certsecret")

const (
	ScopeNamespaced string = "Namespaced"
	ScopeCluster    string = "Cluster"
)

/**
* USER ACTION REQUIRED: This is a scaffold file intended for the user to modify with their own Controller
* business logic.  Delete these comments after modifying this file.*
 */

// Add creates a new CertSecret Controller and adds it to the Manager. The Manager will set fields on the Controller
// and Start it when the Manager is Started.
func Add(mgr manager.Manager) error {
	return add(mgr, newReconciler(mgr))
}

// newReconciler returns a new reconcile.Reconciler
func newReconciler(mgr manager.Manager) reconcile.Reconciler {
	return &ReconcileCertSecret{client: mgr.GetClient(), scheme: mgr.GetScheme()}
}

// add adds a new Controller to mgr with r as the reconcile.Reconciler
func add(mgr manager.Manager, r reconcile.Reconciler) error {
	// Create a new controller
	c, err := controller.New("certsecret-controller", mgr, controller.Options{Reconciler: r})
	if err != nil {
		return err
	}

	// Watch for changes to primary resource CertSecret
	err = c.Watch(&source.Kind{Type: &appv1alpha1.CertSecret{}}, &handler.EnqueueRequestForObject{})
	if err != nil {
		return err
	}

	// Watch for changes to secondary resource Pods and requeue the owner CertSecret
	err = c.Watch(&source.Kind{Type: &appsv1.Deployment{}}, &handler.EnqueueRequestForOwner{
		IsController: true,
		OwnerType:    &appv1alpha1.CertSecret{},
	})
	if err != nil {
		return err
	}
	// 监听namespace变化
	err = c.Watch(&source.Kind{Type: &corev1.Namespace{}}, &enqueueRequestForNs{client: mgr.GetClient()})
	if err != nil {
		return err
	}

	return nil
}

// blank assignment to verify that ReconcileCertSecret implements reconcile.Reconciler
var _ reconcile.Reconciler = &ReconcileCertSecret{}

// ReconcileCertSecret reconciles a CertSecret object
type ReconcileCertSecret struct {
	// This client, initialized using mgr.Client() above, is a split client
	// that reads objects from the cache and writes to the apiserver
	client client.Client
	scheme *runtime.Scheme
}

// Reconcile reads that state of the cluster for a CertSecret object and makes changes based on the state read
// and what is in the CertSecret.Spec
// TODO(user): Modify this Reconcile function to implement your Controller logic.  This example creates
// a Pod as an example
// Note:
// The Controller will requeue the Request to be processed again if the returned error is non-nil or
// Result.Requeue is true, otherwise upon completion it will remove the work from the queue.
func (r *ReconcileCertSecret) Reconcile(request reconcile.Request) (reconcile.Result, error) {
	reqLogger := log.WithValues("Request.Namespace", request.Namespace, "Request.Name", request.Name)
	reqLogger.Info("Reconciling CertSecret")
	// Fetch the CertSecret instance
	instance := &appv1alpha1.CertSecret{}
	err := r.client.Get(context.Background(), request.NamespacedName, instance)
	if err != nil {
		if errors.IsNotFound(err) {
			// Request object not found, could have been deleted after reconcile request.
			// Owned objects are automatically garbage collected. For additional cleanup logic use finalizers.
			// Return and don't requeue
			return reconcile.Result{}, nil
		}
		// Error reading the object - requeue the request.
		return reconcile.Result{}, err
	}

	if instance.DeletionTimestamp != nil {
		return reconcile.Result{}, err
	}
	nss := make([]string, 10)
	scope := instance.Spec.Scope
	switch scope {
	case ScopeCluster:
		nsl := &corev1.NamespaceList{}
		er1 := r.client.List(context.Background(), &client.ListOptions{}, nsl)
		if er1 != nil {
			reqLogger.Error(err, "拉取namespaces失败")
			return reconcile.Result{Requeue: true}, er1
		}
		for _, ns := range nsl.Items {
			nss = append(nss, ns.Name)
		}
	case ScopeNamespaced:
		nss = instance.Spec.Namespaces
	default:
		reqLogger.Info("scope not in [\"Cluster\",\"Namespaced\"]")
		return reconcile.Result{Requeue: false}, nil
	}
	reqLogger.Info(fmt.Sprintf("管理的命名空间为:%s", nss))
	tls := instance.Spec.Tls

	for _, ns := range nss {
		if ns == "" {
			continue
		}
		for _, t := range tls {
			time.Sleep(time.Second)
			// 获取secret
			secret := &corev1.Secret{}
			secret.SetGroupVersionKind(schema.GroupVersionKind{
				Group:   appv1alpha1.SchemeGroupVersion.Group,
				Version: appv1alpha1.SchemeGroupVersion.Version,
				Kind:    "CertSecret",
			})
			err2 := r.client.Get(context.Background(), client.ObjectKey{
				Namespace: ns,
				Name:      t.Name,
			}, secret)
			if err2 != nil && errors.IsNotFound(err2) {
				//reqLogger.Info(fmt.Sprintf("1|查询secret %s/%s 报错，不存在", ns, t.Name))
				// 不存在创建
				reqLogger.Info(fmt.Sprintf("2|创建secret:%s/%s", ns, t.Name))
				newSecret := resources.NewSecret(instance, t, ns)
				if errn := r.client.Create(context.Background(), newSecret); errn != nil {
					reqLogger.Error(errn, "创建secret失败")
					//return reconcile.Result{}, nil
				}
			} else {
				reqLogger.Info(fmt.Sprintf("4| secret:%s/%s 已存在", ns, t.Name))

			}
		}

	}
	// 关联Annotations 保留crd 上次配置，用此字段值与当前配置对比，来判断是否需要更新
	data, _ := json.Marshal(instance.Spec)
	if instance.Annotations != nil {
		oldInstanceSpec := &appv1alpha1.CertSecretSpec{}
		if err := json.Unmarshal([]byte(instance.Annotations["spec"]), oldInstanceSpec); err != nil {
			reqLogger.Info("Annotations[spec]不存在或格式错误")
			reqLogger.Info("7| 第一次创建，保存spec内容到当前实例")
			instance.Annotations["spec"] = string(data)
			if err7 := r.client.Update(context.Background(), instance); err7 != nil {
				reqLogger.Error(err7, "更新Annotations失败")
				//return reconcile.Result{}, nil
			}
			return reconcile.Result{Requeue: false}, nil
		} else if !reflect.DeepEqual(instance.Spec, *oldInstanceSpec) {
			reqLogger.Info("Crd资源更新")
			// 更新secret
			for _, ns := range nss {
				if ns == "" {
					continue
				}
				reqLogger.Info("操作命名空间：" + ns)
				for _, t := range tls {
					// 获取secret
					secret := &corev1.Secret{}
					secret.SetGroupVersionKind(schema.GroupVersionKind{
						Group:   appv1alpha1.SchemeGroupVersion.Group,
						Version: appv1alpha1.SchemeGroupVersion.Version,
						Kind:    "CertSecret",
					})
					err2 := r.client.Get(context.Background(), client.ObjectKey{
						Namespace: ns,
						Name:      t.Name,
					}, secret)
					if err2 != nil {
						reqLogger.Error(err2, fmt.Sprintf("9| 查询Secret %s/%s失败", ns, t.Name))
						return reconcile.Result{Requeue: false}, nil
					}
					reqLogger.Info(fmt.Sprintf("10| 更新Secret %s/%s", ns, t.Name))
					// todo: 判断证书内容是否变化
					if erru := r.client.Update(context.Background(), resources.UpdateSecret(t, secret)); erru != nil {
						reqLogger.Error(erru, fmt.Sprintf("11| 更新Secret %s/%s失败", ns, t.Name))
						return reconcile.Result{Requeue: true}, nil
					}
					time.Sleep(time.Second)
				}
			}
		}
	} else {
		reqLogger.Info("8| 第一次创建，保存spec内容到当前实例")
		instance.Annotations = map[string]string{"spec": string(data)}
		if err := r.client.Update(context.TODO(), instance); err != nil {
			return reconcile.Result{}, nil
		}
		return reconcile.Result{Requeue: false}, nil
	}

	return reconcile.Result{}, nil
}
