package maintenancewatcher

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/go-logr/logr"
	machinev1 "github.com/openshift/machine-api-operator/pkg/apis/machine/v1beta1"
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

// MaintenanceWatcher will scan each machine in the cluster and publish maintenances for them.
// on the frequency of m.watchInterval. Process wil cease when message is sent on the stopCh.
func (m *maintenanceWatcher) Start(log logr.Logger, stopCh <-chan struct{}) {
	log.Info("Starting the maintenanceWatcher")
	log.Info("Scanning Machines for maintenances")
	// Scan machines for maintenances and publish them to a machinemaintenance CR.
	err := m.scanAndPublishMachineMaintences(log)
	if err != nil {
		log.Error(err, "machine maintenance scanner failed to run")
	}

	for {
		select {
		case <-time.After(m.watchInterval):
			log.Info("This is printing in a go routine")
		case <-stopCh:
			log.Info("Stopping the maintenanceWatcher")
			break
		}
	}
}

// scanAndPublishMachineMaintences will retrieve each machine resource ID and query for
// scheduled maintenances from the cloud provider. If a maintenance is present, a
// machinemaintenance CR will be created for the machinemaintenance controller to reconcile.
func (m *maintenanceWatcher) scanAndPublishMachineMaintences(log logr.Logger) error {
	// Get the machine resource IDs.
	machineResourceIDs := make([]string, 0)
	machineResourceIDs, err := m.getMachineResourceIDs(log)
	if err != nil {
		log.Error(err, "failed to return machine resources IDs")
		return err
	}

	// get aws client

	// query the AWS api for maintenance for the given ID.
	for _, machineResourceID := range machineResourceIDs {
		// pass aws client here
		err := checkMachineMaintenance(log, machineResourceID)
		if err != nil {
			return err
		}
	}

	return nil
}

// getMachineResourceIDs gets each machines resource ID and returns them in a string array
// to the caller.
func (m *maintenanceWatcher) getMachineResourceIDs(log logr.Logger) ([]string, error) {
	// Create array to hold machine IDs.
	machineResourceIDs := make([]string, 0)

	// Create machineList type to use as runtime object.
	machineList := &machinev1.MachineList{}

	// Set list options for client.List.
	listOptions := []client.ListOption{
		client.InNamespace("openshift-machine-api"),
	}

	// List based on above configuration.
	err := m.client.List(context.TODO(), machineList, listOptions...)
	if err != nil {
		log.Error(err, "Failed to get machine list")
		return nil, err
	}

	// Iterate through each machine item and append the resource ID returned
	// to the machineResourceIDs string array.
	for _, machineObj := range machineList.Items {
		r := strings.LastIndex(*machineObj.Spec.ProviderID, "/")
		if r != -1 {
			n := *machineObj.Spec.ProviderID
			machineResourceIDs = append(machineResourceIDs, n[r+1:])
		}
	}
	return machineResourceIDs, nil
}

func checkMachineMaintenance(log logr.Logger, mri string) error {
	fmt.Println(mri)
	// receive aws clinet
	// query given machineID for maintenance
	return nil
}
