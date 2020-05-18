# Machine Maintenance Operator

## What problem does it solve?
Currenty the `machine-api` does not handle machine termination via the scheduled events from AWS. This creates alerts that require manual intervention by deleting the missing machines `machine` CR. 

This will not scale. 

Further developement to take into consideration the use of freeze windows will also enable customers to be assurred that these events are occurring during their designated maintenance windows. 

## What it does
* Queries AWS API for each instance ID in the cluster for scheduled events every 60 minutes. 
* Creates a `machinemaintenance` CR that the machinemaintenance controller reconciles
* Current reconciliation is deleting the target machines `machine` CR in the `openshift-machine-api` namespace
* machine-api gracefully terminates node

## What could do
* Query Xchangewindows CRs that are being developed in the managed-upgrade-operator.
	* PreferredUpgradeStartTime -> this could be used as maintenance time
	* AdminFreezeWindow
	* CustomerFreezeWindow
* If its not a suitable time, exit reconcile loop with Result.Requeue = false
* Recncile will check again in 15 mins as per SyncPeriod. 
* If it is a suitable time, delete `machine` CR for target machine. 
* This should also handle different terminations for worker/infra vs master nodes. 

## Open Questions
From testing, it takes approximately 7 minutes for a node to be available in the "Running" state. If this is not acceptable, the reconciler could have the logic to increase machine pool by count+1, then terminate then reduce machine pool by count-1.

## Example Logs

Start up the operator
```
INFO[0000] Running the operator locally in namespace machine-maintenance-operator. 
{"level":"info","ts":1589806376.6754005,"logger":"cmd","msg":"Operator Version: 0.0.1"}
{"level":"info","ts":1589806376.6754684,"logger":"cmd","msg":"Go Version: go1.14.2"}
{"level":"info","ts":1589806376.6754832,"logger":"cmd","msg":"Go OS/Arch: linux/amd64"}
{"level":"info","ts":1589806376.675494,"logger":"cmd","msg":"Version of operator-sdk: v0.15.2"}
{"level":"info","ts":1589806376.693169,"logger":"leader","msg":"Trying to become the leader."}
{"level":"info","ts":1589806376.6931949,"logger":"leader","msg":"Skipping leader election; not running in a cluster."}
{"level":"info","ts":1589806379.315759,"logger":"controller-runtime.metrics","msg":"metrics server is starting to listen","addr":":8080"}
{"level":"info","ts":1589806379.3171518,"logger":"cmd","msg":"Registering Components."}
{"level":"info","ts":1589806379.3180697,"logger":"cmd","msg":"Initializing maintenanceWatcher"}
{"level":"info","ts":1589806379.3181596,"logger":"cmd","msg":"Starting the Cmd."}
{"level":"info","ts":1589806379.3182573,"logger":"cmd","msg":"Starting the maintenanceWatcher"}
{"level":"info","ts":1589806379.3183444,"logger":"cmd","msg":"Initial scanning of Machines for maintenances"}
{"level":"info","ts":1589806379.3184586,"logger":"controller-runtime.manager","msg":"starting metrics server","path":"/metrics"}
{"level":"info","ts":1589806379.3185785,"logger":"controller-runtime.controller","msg":"Starting EventSource","controller":"machinemaintenance-controller","source":"kind source: /, Kind="}
{"level":"info","ts":1589806379.4195812,"logger":"controller-runtime.controller","msg":"Starting EventSource","controller":"machinemaintenance-controller","source":"kind source: /, Kind="}
{"level":"info","ts":1589806379.9218109,"logger":"controller-runtime.controller","msg":"Starting Controller","controller":"machinemaintenance-controller"}
{"level":"info","ts":1589806381.0211668,"logger":"cmd","msg":"Checking maintenance for dofinn-20201705-blz22-worker-ap-southeast-2a-7kvd4"}
{"level":"info","ts":1589806381.0221722,"logger":"controller-runtime.controller","msg":"Starting workers","controller":"machinemaintenance-controller","worker count":1}
{"level":"info","ts":1589806381.2525415,"logger":"cmd","msg":"Currently no maintenance for dofinn-20201705-blz22-worker-ap-southeast-2a-7kvd4"}
{"level":"info","ts":1589806381.2526298,"logger":"cmd","msg":"Checking maintenance for dofinn-20201705-blz22-worker-ap-southeast-2a-984h4"}
{"level":"info","ts":1589806381.3415031,"logger":"cmd","msg":"Currently no maintenance for dofinn-20201705-blz22-worker-ap-southeast-2a-984h4"}
{"level":"info","ts":1589806381.3415604,"logger":"cmd","msg":"Checking maintenance for dofinn-20201705-blz22-worker-ap-southeast-2a-rh8xv"}
{"level":"info","ts":1589806381.4398959,"logger":"cmd","msg":"Currently no maintenance for dofinn-20201705-blz22-worker-ap-southeast-2a-rh8xv"}
{"level":"info","ts":1589806381.4399655,"logger":"cmd","msg":"Checking maintenance for dofinn-20201705-blz22-infra-ap-southeast-2a-4pdbs"}
{"level":"info","ts":1589806381.5298276,"logger":"cmd","msg":"Currently no maintenance for dofinn-20201705-blz22-infra-ap-southeast-2a-4pdbs"}
{"level":"info","ts":1589806381.5298805,"logger":"cmd","msg":"Checking maintenance for dofinn-20201705-blz22-master-1"}
{"level":"info","ts":1589806381.6231802,"logger":"cmd","msg":"Currently no maintenance for dofinn-20201705-blz22-master-1"}
{"level":"info","ts":1589806381.6232882,"logger":"cmd","msg":"Checking maintenance for dofinn-20201705-blz22-worker-ap-southeast-2a-7pd2r"}
{"level":"info","ts":1589806381.7127972,"logger":"cmd","msg":"Currently no maintenance for dofinn-20201705-blz22-worker-ap-southeast-2a-7pd2r"}
{"level":"info","ts":1589806381.712888,"logger":"cmd","msg":"Checking maintenance for dofinn-20201705-blz22-master-0"}
{"level":"info","ts":1589806381.8086264,"logger":"cmd","msg":"Currently no maintenance for dofinn-20201705-blz22-master-0"}
{"level":"info","ts":1589806381.8086927,"logger":"cmd","msg":"Checking maintenance for dofinn-20201705-blz22-infra-ap-southeast-2a-ph26m"}
{"level":"info","ts":1589806381.8795257,"logger":"cmd","msg":"Currently no maintenance for dofinn-20201705-blz22-infra-ap-southeast-2a-ph26m"}
{"level":"info","ts":1589806381.8795803,"logger":"cmd","msg":"Checking maintenance for dofinn-20201705-blz22-master-2"}
{"level":"info","ts":1589806381.9503958,"logger":"cmd","msg":"Currently no maintenance for dofinn-20201705-blz22-master-2"}
{"level":"info","ts":1589806381.9504538,"logger":"cmd","msg":"Checking maintenance for dofinn-20201705-blz22-infra-ap-southeast-2a-w82t2"}
{"level":"info","ts":1589806382.050584,"logger":"cmd","msg":"Currently no maintenance for dofinn-20201705-blz22-infra-ap-southeast-2a-w82t2"}
```

Inject a CR that would trigger the `maintenanceWatcher` to create a `machinemaintenance` CR to be reconciled.

```
oc create -f deploy/crds/machinemaintenance.managed.openshift.io_v1alpha1_machinemaintenance_cr.yaml
machinemaintenance.machinemaintenance.managed.openshift.io/mm-dofinn-20201705-blz22-worker-ap-southeast-2a-984h4 created
```

Reconcile `machinemaintenance` CR

```
{"level":"info","ts":1589806388.3073452,"logger":"controller_machinemaintenance","msg":"Reconciling MachineMaintenance","Request.Namespace":"machine-maintenance-operator","Request.Name":"mm-dofinn-20201705-blz22-worker-ap-southeast-2a-984h4"}
{"level":"info","ts":1589806388.307381,"logger":"controller_machinemaintenance","msg":"Scheduling maintenance for dofinn-20201705-blz22-worker-ap-southeast-2a-984h4","Request.Namespace":"machine-maintenance-operator","Request.Name":"mm-dofinn-20201705-blz22-worker-ap-southeast-2a-984h4"}
{"level":"info","ts":1589806388.3412569,"logger":"controller_machinemaintenance","msg":"Deleting machine dofinn-20201705-blz22-worker-ap-southeast-2a-984h4","Request.Namespace":"machine-maintenance-operator","Request.Name":"mm-dofinn-20201705-blz22-worker-ap-southeast-2a-984h4"}
{"level":"info","ts":1589806388.36995,"logger":"controller_machinemaintenance","msg":"Reconciling MachineMaintenance","Request.Namespace":"machine-maintenance-operator","Request.Name":"mm-dofinn-20201705-blz22-worker-ap-southeast-2a-984h4"}
{"level":"info","ts":1589806388.3700304,"logger":"controller_machinemaintenance","msg":"Maintenance already scheduled for: dofinn-20201705-blz22-worker-ap-southeast-2a-984h4","Request.Namespace":"machine-maintenance-operator","Request.Name":"mm-dofinn-20201705-blz22-worker-ap-southeast-2a-984h4"}
```

The `machine-api` then gracefully terminates the target node

```
I0518 12:53:08.335599       1 controller.go:161] Reconciling Machine "dofinn-20201705-blz22-worker-ap-southeast-2a-984h4"
I0518 12:53:08.335626       1 controller.go:383] Machine "dofinn-20201705-blz22-worker-ap-southeast-2a-984h4" going into phase "Deleting"
I0518 12:53:08.348692       1 controller.go:201] Reconciling machine "dofinn-20201705-blz22-worker-ap-southeast-2a-984h4" triggers delete
I0518 12:53:08.365994       1 info.go:20] cordoned node "ip-10-0-138-180.ap-southeast-2.compute.internal"
I0518 12:53:08.407791       1 info.go:16] ignoring DaemonSet-managed pods: tuned-kdxgq, dns-default-57pg9, node-ca-7b6cr, machine-config-daemon-66xg5, node-exporter-64v68, sre-dns-latency-exporter-w99c6, multus-glfvw, ovs-ghc9n, sdn-l284k, splunkforwarder-ds-rgj2m; deleting pods not managed by ReplicationController, ReplicaSet, Job, DaemonSet or StatefulSet: splunk-forwarder-operator-catalog-96tfs
I0518 12:53:08.445884       1 info.go:20] pod "deployments-pruner-1589796000-89ww4" removed (evicted)
I0518 12:53:08.449840       1 info.go:20] pod "builds-pruner-1589799600-vf5n2" removed (evicted)
I0518 12:53:08.452340       1 info.go:20] pod "builds-pruner-1589803200-62xfd" removed (evicted)
I0518 12:53:08.524154       1 info.go:20] pod "builds-pruner-1589796000-5rflg" removed (evicted)
I0518 12:53:08.602280       1 info.go:20] pod "alert-pruner-1589799600-72rfn" removed (evicted)
.....TRUNCATED.....
I0518 12:53:44.077230       1 actuator.go:628] dofinn-20201705-blz22-worker-ap-southeast-2a-xj2rc: finished calculating AWS status
I0518 12:53:44.077339       1 actuator.go:224] dofinn-20201705-blz22-worker-ap-southeast-2a-xj2rc: machine status has changed, updating
I0518 12:53:44.090960       1 controller.go:383] Machine "dofinn-20201705-blz22-worker-ap-southeast-2a-xj2rc" going into phase "Provisioned"
```
