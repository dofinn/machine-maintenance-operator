package machinemaintenance

import (
	"context"
	"fmt"

	machinemaintenancev1alpha1 "github.com/openshift/machine-maintenance-operator/pkg/apis/machinemaintenance/v1alpha1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/manager"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
	"sigs.k8s.io/controller-runtime/pkg/source"
)

var log = logf.Log.WithName("controller_machinemaintenance")

/**
* USER ACTION REQUIRED: This is a scaffold file intended for the user to modify with their own Controller
* business logic.  Delete these comments after modifying this file.*
 */

// Add creates a new MachineMaintenance Controller and adds it to the Manager. The Manager will set fields on the Controller
// and Start it when the Manager is Started.
func Add(mgr manager.Manager) error {
	return add(mgr, newReconciler(mgr))
}

// newReconciler returns a new reconcile.Reconciler
func newReconciler(mgr manager.Manager) reconcile.Reconciler {
	return &ReconcileMachineMaintenance{client: mgr.GetClient(), scheme: mgr.GetScheme()}
}

// add adds a new Controller to mgr with r as the reconcile.Reconciler
func add(mgr manager.Manager, r reconcile.Reconciler) error {
	// Create a new controller
	c, err := controller.New("machinemaintenance-controller", mgr, controller.Options{Reconciler: r})
	if err != nil {
		return err
	}

	// Watch for changes to primary resource MachineMaintenance
	err = c.Watch(&source.Kind{Type: &machinemaintenancev1alpha1.MachineMaintenance{}}, &handler.EnqueueRequestForObject{})
	if err != nil {
		return err
	}

	// TODO(user): Modify this to be the types you create that are owned by the primary resource
	// Watch for changes to secondary resource Pods and requeue the owner MachineMaintenance
	err = c.Watch(&source.Kind{Type: &corev1.Pod{}}, &handler.EnqueueRequestForOwner{
		IsController: true,
		OwnerType:    &machinemaintenancev1alpha1.MachineMaintenance{},
	})
	if err != nil {
		return err
	}

	return nil
}

// blank assignment to verify that ReconcileMachineMaintenance implements reconcile.Reconciler
var _ reconcile.Reconciler = &ReconcileMachineMaintenance{}

// ReconcileMachineMaintenance reconciles a MachineMaintenance object
type ReconcileMachineMaintenance struct {
	// This client, initialized using mgr.Client() above, is a split client
	// that reads objects from the cache and writes to the apiserver
	client client.Client
	scheme *runtime.Scheme
}

// Reconcile reads that state of the cluster for a MachineMaintenance object and makes changes based on the state read
// and what is in the MachineMaintenance.Spec
// TODO(user): Modify this Reconcile function to implement your Controller logic.  This example creates
// a Pod as an example
// Note:
// The Controller will requeue the Request to be processed again if the returned error is non-nil or
// Result.Requeue is true, otherwise upon completion it will remove the work from the queue.
func (r *ReconcileMachineMaintenance) Reconcile(request reconcile.Request) (reconcile.Result, error) {
	reqLogger := log.WithValues("Request.Namespace", request.Namespace, "Request.Name", request.Name)
	reqLogger.Info("Reconciling MachineMaintenance")

	// Fetch the MachineMaintenance instance
	instance := &machinemaintenancev1alpha1.MachineMaintenance{}
	err := r.client.Get(context.TODO(), request.NamespacedName, instance)
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
	if instance.Spec.Maintenance != true {
		reqLogger.Info("Maintenance is NOT scheduled :(")
		reqLogger.Info("Setting maintenance to true")
		instance.Spec.Maintenance = true
		err := r.client.Update(context.TODO(), instance)
		if err != nil {
			reqLogger.Info((fmt.Sprintf("An error occurred: %s", err)))
			return reconcile.Result{}, err
		}
		return reconcile.Result{}, nil
	} else {
		reqLogger.Info("Maintenance is currently scheduled :)")
		log.Info(fmt.Sprintf("machine maintenance CR\n%+v\n", instance.Spec))
		return reconcile.Result{}, nil
	}
}

//	// Define a new Pod object
//	pod := newPodForCR(instance)
//
//	// Set MachineMaintenance instance as the owner and controller
//	if err := controllerutil.SetControllerReference(instance, pod, r.scheme); err != nil {
//		return reconcile.Result{}, err
//	}
//
//	// Check if this Pod already exists
//	found := &corev1.Pod{}
//	err = r.client.Get(context.TODO(), types.NamespacedName{Name: pod.Name, Namespace: pod.Namespace}, found)
//	if err != nil && errors.IsNotFound(err) {
//		reqLogger.Info("Creating a new Pod", "Pod.Namespace", pod.Namespace, "Pod.Name", pod.Name)
//		err = r.client.Create(context.TODO(), pod)
//		if err != nil {
//			return reconcile.Result{}, err
//		}
//
//		// Pod created successfully - don't requeue
//		return reconcile.Result{}, nil
//	} else if err != nil {
//		return reconcile.Result{}, err
//	}
//
//	// Pod already exists - don't requeue
//	reqLogger.Info("Skip reconcile: Pod already exists", "Pod.Namespace", found.Namespace, "Pod.Name", found.Name)
//	return reconcile.Result{}, nil
//}

// newPodForCR returns a busybox pod with the same name/namespace as the cr
func newPodForCR(cr *machinemaintenancev1alpha1.MachineMaintenance) *corev1.Pod {
	labels := map[string]string{
		"app": cr.Name,
	}
	return &corev1.Pod{
		ObjectMeta: metav1.ObjectMeta{
			Name:      cr.Name + "-pod",
			Namespace: cr.Namespace,
			Labels:    labels,
		},
		Spec: corev1.PodSpec{
			Containers: []corev1.Container{
				{
					Name:    "busybox",
					Image:   "busybox",
					Command: []string{"sleep", "3600"},
				},
			},
		},
	}
}
