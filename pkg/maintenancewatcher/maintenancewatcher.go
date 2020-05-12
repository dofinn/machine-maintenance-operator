package maintenancewatcher

import (
	"context"
	"fmt"
	"time"

	machinev1 "github.com/openshift/machine-api-operator/pkg/apis/machine/v1beta1"
	"github.com/prometheus/common/log"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

const ()

// SecretWatcher global var for SecretWatcher
var MaintenanceWatcher *maintenanceWatcher

type maintenanceWatcher struct {
	watchInterval time.Duration
	client        client.Client
}

// Initialize creates a global instance of the SecretWatcher
func Initialize(client client.Client, watchInterval time.Duration) {
	MaintenanceWatcher = NewMaintenanceWatcher(client, watchInterval)
}

// NewMaintenanceWatcher returns a new instance of the maintenanceWatcher struct
func NewMaintenanceWatcher(client client.Client, watchInterval time.Duration) *maintenanceWatcher {
	return &maintenanceWatcher{
		watchInterval: watchInterval,
		client:        client,
	}
}

// MaintenanceWatcher will trigger CredentialsRotator every `scanInternal` and only stop if the operator is killed or a
// message is sent on the stopCh
//func (m *maintenanceWatcher) Start(log logr.Logger, stopCh <-chan struct{}) {
func (m *maintenanceWatcher) Start() {
	log.Info("Starting the maintenanceWatcher")
	log.Info("MaintenanceWatcher initial run")
	// Get list of nodes by ID
	// Query if any events scheduled for node
	// if maintenance, update PlatformMaintenanceSpec.Maintenance = true

	//	machineList := &machinev1.MachineList{}
	machineList := &machinev1.Machine{}
	machineNamespacedName := types.NamespacedName{Namespace: "openshift-machine-api-operator"}

	if err := m.client.Get(context.TODO(), machineNamespacedName, machineList); err != nil {
		log.Error(err, fmt.Sprintf("Unable to list machines in namespace %s", "openshift-machine-api-operator"))
	}

	fmt.Printf("%v", machineList)

	//	for {
	//		select {
	//		case <-time.After(m.watchInterval):
	log.Info("maintenanceWatcher: scanning for maintenances")
	time.Sleep(time.Duration(5) * time.Minute)
	fmt.Printf("%v", machineList)

	//		case <-stopCh:
	//			log.Info("Stopping the maintenance scanner")
	//			break
	//		}
}
