package maintenancewatcher

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/go-logr/logr"
	machinemaintenancev1alpha1 "github.com/openshift/machine-maintenance-operator/pkg/apis/machinemaintenance/v1alpha1"
	"github.com/openshift/machine-maintenance-operator/pkg/awsclient"
	"github.com/openshift/machine-maintenance-operator/pkg/controller/utils"

	machinev1 "github.com/openshift/machine-api-operator/pkg/apis/machine/v1beta1"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

const (
	// Name of secret that holds AWS creds for IAM user.
	AWSSecretName     = "machine-maintenance-operator-credentials"
	OperatorNamespace = "machine-maintenance-operator"

	// AWS Event Codes
	eventCode          = "event.code"
	instanceReboot     = "instance-reboot"
	systemReboot       = "system-reboot"
	systemMaintenance  = "system-maintenance"
	instanceRetirement = "instance-retirement"
	instanceStop       = "instance-stop"
)

// Filter for DescribeInstanceStatusInput so that only instances with
// events are returned.
var maintenanceFilter = []*ec2.Filter{
	{
		Name: aws.String(eventCode),
		Values: []*string{aws.String(instanceStop), aws.String(instanceReboot),
			aws.String(systemReboot), aws.String(instanceRetirement), aws.String(systemMaintenance)},
	},
}

// SecretWatcher global var for SecretWatcher.
var MaintenanceWatcher *maintenanceWatcher

// maintenanceWatcher struct that the below methods are built around.
type maintenanceWatcher struct {
	watchInterval time.Duration
	client        client.Client
}

// Initialize creates a global instance of the SecretWatcher.
func Initialize(client client.Client, watchInterval time.Duration) {
	MaintenanceWatcher = NewMaintenanceWatcher(client, watchInterval)
}

// NewMaintenanceWatcher returns a new instance of the maintenanceWatcher struct.
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
	log.Info("Initial scanning of Machines for maintenances")

	// Scan machines for maintenances and publish them to a machinemaintenance CR.
	err := m.scanAndPublishMachineMaintences(log)
	if err != nil {
		log.Error(err, "Machine maintenance scanner failed initial run")
	}

	// While true, perform a scan ever m.watchInterval and process any instances
	// that have returned with maintenance events.
	for {
		select {
		case <-time.After(m.watchInterval):
			log.Info("Scanning Machines for maintenances")
			// Scan machines for maintenances and publish them to a machinemaintenance CR.
			err := m.scanAndPublishMachineMaintences(log)

			if err != nil {
				log.Error(err, "Machine maintenance scanner failed to run")
			}

			// Listens for SIGTERM or SIGINT.
		case <-stopCh:
			log.Info("Stopping the maintenanceWatcher")
			break
		}
	}
}

// scanAndPublishMachineMaintences will retrieve each machine resources (Name and ID) and query for
// scheduled maintenances from the cloud provider. If a maintenance is present, a
// machinemaintenance CR will be created for the machinemaintenance controller to reconcile.
func (m *maintenanceWatcher) scanAndPublishMachineMaintences(log logr.Logger) error {
	// Retrieve machine resources.
	machineResources, err := m.getMachineResources(log)
	if err != nil {
		log.Error(err, "failed to return machine resources IDs")
		return err
	}

	// Get cluster region.
	region, err := utils.GetClusterRegion(m.client)
	if err != nil {
		return err
	}

	// Get awsclient using operator secret and cluster region.
	awsClient, err := awsclient.GetAWSClient(m.client, awsclient.NewAwsClientInput{
		SecretName: AWSSecretName,
		NameSpace:  OperatorNamespace,
		AwsRegion:  region,
	})

	if err != nil {
		log.Error(err, "Failed to get AWS client")
		return err
	}

	// Query the AWS api for maintenance for the given ID.
	// TODO : AWS can actually accept a list.
	for _, machineResource := range machineResources {
		log.Info(fmt.Sprintf("Checking maintenance for %s", machineResource["name"]))
		err := m.checkMachineMaintenance(log, machineResource, awsClient)
		if err != nil {
			return err
		}
	}

	return nil
}

// Create machineMap type so getMachineResources can return it
type machineMap map[string]string

// getMachineResources gets each machines resource ID and returns them in a string array
// to the caller.
func (m *maintenanceWatcher) getMachineResources(log logr.Logger) ([]machineMap, error) {
	// machineMaps holds a slice of machineMap
	var machineMaps []machineMap

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
	// to the machineResources string array.
	for _, machineObj := range machineList.Items {
		machineName := *&machineObj.ObjectMeta.Name
		r := strings.LastIndex(*machineObj.Spec.ProviderID, "/")
		if r != -1 {
			n := *machineObj.Spec.ProviderID
			mMap := machineMap{
				"name": machineName,
				"id":   n[r+1:],
			}
			machineMaps = append(machineMaps, mMap)
		}
	}
	// Return machineMaps to caller.
	return machineMaps, nil
}

// checkMachineMaintenance uses each machines resource ID to query the API and check
// if that machine has a scheduled maintenance.
func (m *maintenanceWatcher) checkMachineMaintenance(log logr.Logger, mMap machineMap, awsclient *awsclient.AwsClient) error {
	// Split map into ID and Name for convenience.
	mID := mMap["id"]
	mName := mMap["name"]

	// Prepare input for request.
	input := &ec2.DescribeInstanceStatusInput{
		InstanceIds: []*string{
			aws.String(mID),
		},
		Filters: maintenanceFilter,
	}

	// Call to AWS API for instance status.
	result, err := awsclient.DescribeInstanceStatus(input)
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			switch aerr.Code() {
			default:
				fmt.Println(aerr.Error())
			}
		} else {
			// Print the error, cast err to awserr.Error to get the Code and
			// Message from an error.
			fmt.Println(err.Error())
		}
		return err
	}

	// Based on the maintenanceFilter, no results will be returned for machines
	// without a maintenance scheduled.
	if len(result.InstanceStatuses) != 0 {
		log.Info(fmt.Sprintf("Found maintenance for %s", mName))
		m.createMaintenanceCR(log, result, mMap)
		log.Info(fmt.Sprintf("Creating maintenance for %s", mName))
	} else {
		log.Info(fmt.Sprintf("Currently no maintenance for %s", mName))
	}

	return nil
}

// createMaintenanceCR populates a MachineMaintenance CR with required inputs and creates it
// so it can be reconiled by the machinemaintenance controller.
func (m *maintenanceWatcher) createMaintenanceCR(log logr.Logger, maintenance *ec2.DescribeInstanceStatusOutput, mMap machineMap) error {
	// Get a MachineMaintenance CR for population.
	machineMaintenance := machinemaintenancev1alpha1.MachineMaintenance{}

	// Retrieve event details from the instance status.
	// TODO: Handle multiple events for the one instance.
	mCurrent := maintenance.InstanceStatuses[0]

	// Construct MachineMaintenance CR
	machineMaintenance.Spec.MaintenanceScheduled = false
	machineMaintenance.Spec.EventCode = *mCurrent.Events[0].Code
	machineMaintenance.Spec.EventID = *mCurrent.Events[0].InstanceEventId
	machineMaintenance.Spec.MachineID = *mCurrent.InstanceId
	machineMaintenance.ObjectMeta.Namespace = OperatorNamespace
	machineMaintenance.ObjectMeta.Name = "mm-" + mMap["name"]
	machineMaintenance.Spec.MachineLink = mMap["name"]
	// TODO: handle JSON Time.
	//	machineMaintenance.Spec.NotBefore = *mCurrent.Events[0].NotBefore

	err := m.client.Create(context.TODO(), &machineMaintenance)
	if err != nil {
		log.Error(err, "Failed to create machineMaintenance CR")
	}

	// Return nil if successful creation of MachineMaintenance CR.
	return nil
}

// Example Event for an actual scheduled machine.
/*
{
  InstanceStatuses: [{
      AvailabilityZone: "us-east-1a",
      Events: [{
          Code: "instance-stop",
          Description: "The instance is running on degraded hardware",
          InstanceEventId: "instance-event-01d0903276a5d038c",
          NotBefore: 2020-05-25 06:00:00 +0000 UTC
        }],
      InstanceId: "i-00b75cd359cf95e26",
      InstanceState: {
        Code: 16,
        Name: "running"
      },
      InstanceStatus: {
        Details: [{
            Name: "reachability",
            Status: "passed"
          }],
        Status: "ok"
      },
      SystemStatus: {
        Details: [{
            Name: "reachability",
            Status: "passed"
          }],
        Status: "ok"
      }
    }]
}
*/
