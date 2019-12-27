package nginx

import (
	"context"
	"reflect"

	daasv1 "nginx-operator/pkg/apis/daas/v1"
	"nginx-operator/pkg/common"

	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/manager"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
	"sigs.k8s.io/controller-runtime/pkg/source"
)

var log = logf.Log.WithName("controller_nginx")

/**
* USER ACTION REQUIRED: This is a scaffold file intended for the user to modify with their own Controller
* business logic.  Delete these comments after modifying this file.*
 */

// Add creates a new Nginx Controller and adds it to the Manager. The Manager will set fields on the Controller
// and Start it when the Manager is Started.
func Add(mgr manager.Manager) error {
	return add(mgr, newReconciler(mgr))
}

// newReconciler returns a new reconcile.Reconciler
func newReconciler(mgr manager.Manager) reconcile.Reconciler {
	return &ReconcileNginx{client: mgr.GetClient(), scheme: mgr.GetScheme()}
}

// add adds a new Controller to mgr with r as the reconcile.Reconciler
func add(mgr manager.Manager, r reconcile.Reconciler) error {
	// Create a new controller
	c, err := controller.New("nginx-controller", mgr, controller.Options{Reconciler: r})
	if err != nil {
		return err
	}

	// Watch for changes to primary resource Nginx
	if err = c.Watch(&source.Kind{Type: &daasv1.Nginx{}}, &handler.EnqueueRequestForObject{}); err != nil {
		return err
	}

	// Watch for changes to secondary resource Pods and requeue the owner Nginx
	if err = c.Watch(&source.Kind{Type: &appsv1.Deployment{}}, &handler.EnqueueRequestForOwner{IsController: true, OwnerType: &daasv1.Nginx{}}); err != nil {
		return err
	}

	// Watch for changes to secondary resource Pods and requeue the owner Nginx
	if err = c.Watch(&source.Kind{Type: &corev1.Service{}}, &handler.EnqueueRequestForOwner{IsController: true, OwnerType: &daasv1.Nginx{}}); err != nil {
		return err
	}

	return nil
}

// blank assignment to verify that ReconcileNginx implements reconcile.Reconciler
var _ reconcile.Reconciler = &ReconcileNginx{}

// ReconcileNginx reconciles a Nginx object
type ReconcileNginx struct {
	// This client, initialized using mgr.Client() above, is a split client
	// that reads objects from the cache and writes to the apiserver
	client client.Client
	scheme *runtime.Scheme
}

// Reconcile reads that state of the cluster for a Nginx object and makes changes based on the state read
// and what is in the Nginx.Spec
// Note:
// The Controller will requeue the Request to be processed again if the returned error is non-nil or
// Result.Requeue is true, otherwise upon completion it will remove the work from the queue.
func (r *ReconcileNginx) Reconcile(request reconcile.Request) (reconcile.Result, error) {
	reqLogger := log.WithValues("Request.Namespace", request.Namespace, "Request.Name", request.Name)
	reqLogger.Info("Reconciling Nginx")

	// Fetch the Nginx instance
	instance := &daasv1.Nginx{}
	if err := r.client.Get(context.TODO(), request.NamespacedName, instance); err != nil || instance.DeletionTimestamp != nil {
		if errors.IsNotFound(err) {
			// Request object not found, could have been deleted after reconcile request.
			// Owned objects are automatically garbage collected. For additional cleanup logic use finalizers.
			// Return and don't requeue
			return reconcile.Result{}, nil
		}
		// Error reading the object - requeue the request.
		return reconcile.Result{}, err
	}

	// Install Deployment
	foundDeploy := &appsv1.Deployment{}
	currentDeploy := common.NewDeploymentForCR(instance)
	if err := r.client.Get(context.TODO(), request.NamespacedName, foundDeploy); err != nil {
		if errors.IsNotFound(err) {
			if err := controllerutil.SetControllerReference(instance, currentDeploy, r.scheme); err != nil {
				return reconcile.Result{}, err
			}
			if err := r.client.Create(context.TODO(), currentDeploy); err != nil {
				return reconcile.Result{}, err
			}
		} else {
			return reconcile.Result{}, err
		}
	} else {
		if !reflect.DeepEqual(foundDeploy.Spec, currentDeploy.Spec) {
			foundDeploy.Spec = *currentDeploy.Spec.DeepCopy()
			if err := r.client.Update(context.TODO(), foundDeploy); err != nil {
				return reconcile.Result{}, err
			}
		}
	}

	// Install Service
	foundService := &corev1.Service{}
	currentService := common.NewServiceForCR(instance)
	if err := r.client.Get(context.TODO(), request.NamespacedName, foundService); err != nil {
		if err := controllerutil.SetControllerReference(instance, currentService, r.scheme); err != nil {
			return reconcile.Result{}, err
		}
		if err := r.client.Create(context.TODO(), currentService); err != nil {
			return reconcile.Result{}, err
		}
	} else {
		if !reflect.DeepEqual(foundService.Spec, currentService.Spec) {
			foundService.Spec.Ports = currentService.Spec.Ports
			if err := r.client.Update(context.TODO(), foundService); err != nil {
				return reconcile.Result{}, err
			}
		}
	}

	return reconcile.Result{}, nil
}
