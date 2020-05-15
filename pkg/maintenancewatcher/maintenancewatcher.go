package maintenancewatcher

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/dofinn/machine-maintenance-operator/pkg/awsclient"
	"github.com/go-logr/logr"
	machinev1 "github.com/openshift/machine-api-operator/pkg/apis/machine/v1beta1"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

const (
	AWSSecretName     = "machine-maintenance-operator-credentials"
	OperatorNamespace = "machine-maintenance-operator"

	// AWS Event Codes
	instanceReboot     = "instance-reboot"
	systemReboot       = "system-reboot"
	systemMaintenance  = "system-maintenance"
	instanceRetirement = "instance-retirement"
	instanceStop       = "instance-stop"
)

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
		log.Error(err, "machine maintenance scanner failed initial run")
	}

	for {
		select {
		case <-time.After(m.watchInterval):
			log.Info("Scanning Machines for maintenances")
			// Scan machines for maintenances and publish them to a machinemaintenance CR.
			err := m.scanAndPublishMachineMaintences(log)

			if err != nil {
				log.Error(err, "machine maintenance scanner failed to run")
			}

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

	// Get region
	//region, err := utils.GetClusterRegion(m.client)
	//if err != nil {
	//	return err
	//}
	region := "us-east-1"

	awsClient, err := awsclient.GetAWSClient(m.client, awsclient.NewAwsClientInput{
		SecretName: AWSSecretName,
		NameSpace:  OperatorNamespace,
		AwsRegion:  region,
	})

	if err != nil {
		log.Error(err, "Failed to get AWS client")
		return err
	}

	machineResourceIDs = make([]string, 0)
	machineResourceIDs = append(machineResourceIDs, "i-00b75cd359cf95e26")

	// Query the AWS api for maintenance for the given ID.
	// TODO : AWS can actually accept a list.
	for _, machineResourceID := range machineResourceIDs {
		// pass aws client here
		err := checkMachineMaintenance(log, machineResourceID, awsClient)
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

// checkMachineMaintenance uses each machines resource ID to query the API and check
// if that machine has a scheduled maintenance.
func checkMachineMaintenance(log logr.Logger, mri string, awsclient *awsclient.AwsClient) error {

	// Prepare input for request.
	input := &ec2.DescribeInstanceStatusInput{
		InstanceIds: []*string{
			aws.String(mri),
		},
		Filters: []*ec2.Filter{
			{
				Name:   aws.String("event.code"),
				Values: []*string{aws.String(instanceStop), aws.String(instanceReboot)},
			},
		}}

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

	fmt.Println(result)
	// receive aws clinet
	// query given machineID for maintenance
	return nil
}

/*
[dofinn@bastion-nasa-1 ~]$ aws --region us-east-1 ec2 describe-instance-status --instance-id i-00b75cd359cf95e26 --query "InstanceStatuses[].Events"
[
    [
        {
            "Code": "instance-stop",
            "Description": "The instance is running on degraded hardware",
            "NotBefore": "2020-05-25T06:00:00.000Z"
        }
    ]
]
*/
