package machinemaintenance

import (
	"context"
	"fmt"
	"time"

	machinev1 "github.com/openshift/machine-api-operator/pkg/apis/machine/v1beta1"
	machinemaintenancev1alpha1 "github.com/openshift/machine-maintenance-operator/pkg/apis/machinemaintenance/v1alpha1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
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

	// Fetch the MachineMaintenance mm
	mm := &machinemaintenancev1alpha1.MachineMaintenance{}
	err := r.client.Get(context.TODO(), request.NamespacedName, mm)
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

	/* If maintenance is false:
	1. Query the:
		* CustomerFreezeWindow
		* AdminFreezeWindow
		* PreferredUpgradeTime -> could use this as "maintenance" time
	2. Calculate appropriate maintenance time based on above
	3. set the above time as a deletion timestamp on the machines machineconfig
	4. Update status of CR to "scheduled"
	5. Set the owner of CR to that of the machinepool object? IE when machine is deleted,
	the machine maintenance is deleted?
	*/
	if mm.Spec.MaintenanceScheduled != true {
		reqLogger.Info(fmt.Sprintf("Scheduling maintenance for %s", mm.Spec.MachineLink))
		// Query freeze windows here
		// Calculate appropriate maintenanance time here

		// fetch machine CR from .MachineLink
		targetMachine := &machinev1.Machine{}

		machineLink := mm.Spec.MachineLink

		err := r.client.Get(context.TODO(),
			types.NamespacedName{
				Name:      machineLink,
				Namespace: "openshift-machine-api",
			},
			targetMachine)

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

		// Set deletion timestamp on targetMachine
		//		targetMachine.ObjectMeta.DeletionTimestamp = &metav1.Time{Time: time.Now()}
		targetMachine.ObjectMeta.SetDeletionTimestamp(&metav1.Time{Time: time.Now()})

		err = r.client.Update(context.TODO(), targetMachine)
		if err != nil {
			log.Error(err, "unable to update machine CR with deletiontimestamp")
			return reconcile.Result{}, err
		}

		fmt.Println("This is retrived from machine CR")
		fmt.Println(targetMachine.ObjectMeta.Name)
		// set deletiontimestamp on machine CR
		// set machinemaintenance as finalizer on machine CR
		mm.Spec.MaintenanceScheduled = true
		err = r.client.Update(context.TODO(), mm)
		if err != nil {
			reqLogger.Info((fmt.Sprintf("An error occurred: %s", err)))
			return reconcile.Result{}, err
		}
		return reconcile.Result{}, nil
	} else {
		reqLogger.Info(fmt.Sprintf("Maintenance already scheduled for: %s", mm.Spec.MachineLink))
		return reconcile.Result{}, nil
	}
}
