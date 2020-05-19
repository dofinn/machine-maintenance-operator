# Machine Maintenance Operator Investigation

The follow info is a justification for a machine-maintenance operator that will resolve what is not currently handled by the `openshift-machine-api`. 

The below examines how the `openshift-machine-api` handles different states that are introduced via the cloud provider, in this case AWS. 

## Set instance state to Stop via the aws console.

### Summary

`openshift-machine-api` detects changes however the machine CR indicates this instance is still "Running"
This fails to see that the node is unavailable leaving the cluster in a degraded state that is not reconciled by the `openshift-machine-api`.

```
oc get machines -n openshift-machine-api | grep dofinn-20201705-blz22-worker-ap-southeast-2a-xj2rc
dofinn-20201705-blz22-worker-ap-southeast-2a-xj2rc   Running   m5.xlarge   ap-southeast-2   ap-southeast-2a   10h
```

The node becomes NotReady

```
oc get nodes | grep ip-10-0-140-28.ap-southeast-2.compute.internal                                                          130 ↵
ip-10-0-140-28.ap-southeast-2.compute.internal    NotReady   worker         10h   v1.16.2
```

`openshift-machine-api` still sees the instance as running. 

```
I0518 23:26:43.230082       1 controller.go:161] Reconciling Machine "dofinn-20201705-blz22-worker-ap-southeast-2a-xj2rc"
I0518 23:26:43.230109       1 actuator.go:500] dofinn-20201705-blz22-worker-ap-southeast-2a-xj2rc: Checking if machine exists
I0518 23:26:48.323587       1 actuator.go:549] dofinn-20201705-blz22-worker-ap-southeast-2a-xj2rc: Found instance by id: i-07b30d9dbde6351dc
I0518 23:26:48.323617       1 controller.go:258] Reconciling machine "dofinn-20201705-blz22-worker-ap-southeast-2a-xj2rc" triggers idempotent update
I0518 23:26:48.323628       1 actuator.go:406] dofinn-20201705-blz22-worker-ap-southeast-2a-xj2rc: updating machine
I0518 23:26:48.323715       1 actuator.go:414] dofinn-20201705-blz22-worker-ap-southeast-2a-xj2rc: obtaining EC2 client for region
I0518 23:26:48.381203       1 actuator.go:549] dofinn-20201705-blz22-worker-ap-southeast-2a-xj2rc: Found instance by id: i-07b30d9dbde6351dc
I0518 23:26:48.381221       1 actuator.go:430] dofinn-20201705-blz22-worker-ap-southeast-2a-xj2rc: found 1 existing instances for machine
I0518 23:26:48.391870       1 actuator.go:185] dofinn-20201705-blz22-worker-ap-southeast-2a-xj2rc: ProviderID already set in the machine Spec with value:aws:///ap-southeast-2a/i-07b30d9dbde6351dc
I0518 23:26:48.391893       1 actuator.go:599] dofinn-20201705-blz22-worker-ap-southeast-2a-xj2rc: Updating status
I0518 23:26:48.391940       1 actuator.go:628] dofinn-20201705-blz22-worker-ap-southeast-2a-xj2rc: finished calculating AWS status
I0518 23:26:48.392007       1 actuator.go:224] dofinn-20201705-blz22-worker-ap-southeast-2a-xj2rc: machine status has changed, updating
I0518 23:26:48.403094       1 controller.go:161] Reconciling Machine "dofinn-20201705-blz22-worker-ap-southeast-2a-xj2rc"
I0518 23:26:48.403116       1 actuator.go:500] dofinn-20201705-blz22-worker-ap-southeast-2a-xj2rc: Checking if machine exists
I0518 23:26:48.468628       1 actuator.go:549] dofinn-20201705-blz22-worker-ap-southeast-2a-xj2rc: Found instance by id: i-07b30d9dbde6351dc
I0518 23:26:48.468720       1 controller.go:258] Reconciling machine "dofinn-20201705-blz22-worker-ap-southeast-2a-xj2rc" triggers idempotent update
I0518 23:26:48.468764       1 actuator.go:406] dofinn-20201705-blz22-worker-ap-southeast-2a-xj2rc: updating machine
I0518 23:26:48.468859       1 actuator.go:414] dofinn-20201705-blz22-worker-ap-southeast-2a-xj2rc: obtaining EC2 client for region
I0518 23:26:48.531517       1 actuator.go:549] dofinn-20201705-blz22-worker-ap-southeast-2a-xj2rc: Found instance by id: i-07b30d9dbde6351dc
I0518 23:26:48.531541       1 actuator.go:430] dofinn-20201705-blz22-worker-ap-southeast-2a-xj2rc: found 1 existing instances for machine
E0518 23:26:48.535740       1 controller.go:260] Error updating machine "openshift-machine-api/dofinn-20201705-blz22-worker-ap-southeast-2a-xj2rc": dofinn-20201705-blz22-worker-ap-southeast-2a-xj2rc: failed to set machine cloud provider specifics: dofinn-20201705-blz22-worker-ap-southeast-2a-xj2rc: error updating machine spec: Operation cannot be fulfilled on machines.machine.openshift.io "dofinn-20201705-blz22-worker-ap-southeast-2a-xj2rc": the object has been modified; please apply your changes to the latest version and try again
I0518 23:26:49.535960       1 controller.go:161] Reconciling Machine "dofinn-20201705-blz22-worker-ap-southeast-2a-xj2rc"
I0518 23:26:49.535989       1 actuator.go:500] dofinn-20201705-blz22-worker-ap-southeast-2a-xj2rc: Checking if machine exists
I0518 23:26:49.599332       1 actuator.go:549] dofinn-20201705-blz22-worker-ap-southeast-2a-xj2rc: Found instance by id: i-07b30d9dbde6351dc
I0518 23:26:49.599357       1 controller.go:258] Reconciling machine "dofinn-20201705-blz22-worker-ap-southeast-2a-xj2rc" triggers idempotent update
I0518 23:26:49.599367       1 actuator.go:406] dofinn-20201705-blz22-worker-ap-southeast-2a-xj2rc: updating machine
I0518 23:26:49.599448       1 actuator.go:414] dofinn-20201705-blz22-worker-ap-southeast-2a-xj2rc: obtaining EC2 client for region
I0518 23:26:49.667946       1 actuator.go:549] dofinn-20201705-blz22-worker-ap-southeast-2a-xj2rc: Found instance by id: i-07b30d9dbde6351dc
I0518 23:26:49.667969       1 actuator.go:430] dofinn-20201705-blz22-worker-ap-southeast-2a-xj2rc: found 1 existing instances for machine
I0518 23:26:49.679773       1 actuator.go:185] dofinn-20201705-blz22-worker-ap-southeast-2a-xj2rc: ProviderID already set in the machine Spec with value:aws:///ap-southeast-2a/i-07b30d9dbde6351dc
I0518 23:26:49.679792       1 actuator.go:599] dofinn-20201705-blz22-worker-ap-southeast-2a-xj2rc: Updating status
I0518 23:26:49.679831       1 actuator.go:628] dofinn-20201705-blz22-worker-ap-southeast-2a-xj2rc: finished calculating AWS status
I0518 23:26:49.679909       1 actuator.go:233] dofinn-20201705-blz22-worker-ap-southeast-2a-xj2rc: status unchanged
I0518 23:27:55.909080       1 controller.go:161] Reconciling Machine "dofinn-20201705-blz22-worker-ap-southeast-2a-7pd2r"
I0518 23:27:55.909108       1 actuator.go:500] dofinn-20201705-blz22-worker-ap-southeast-2a-7pd2r: Checking if machine exists
I0518 23:27:56.005066       1 actuator.go:549] dofinn-20201705-blz22-worker-ap-southeast-2a-7pd2r: Found instance by id: i-0488bee4f8178449b
I0518 23:27:56.005089       1 controller.go:258] Reconciling machine "dofinn-20201705-blz22-worker-ap-southeast-2a-7pd2r" triggers idempotent update
I0518 23:27:56.005098       1 actuator.go:406] dofinn-20201705-blz22-worker-ap-southeast-2a-7pd2r: updating machine
I0518 23:27:56.005182       1 actuator.go:414] dofinn-20201705-blz22-worker-ap-southeast-2a-7pd2r: obtaining EC2 client for region
I0518 23:27:56.069229       1 actuator.go:549] dofinn-20201705-blz22-worker-ap-southeast-2a-7pd2r: Found instance by id: i-0488bee4f8178449b
I0518 23:27:56.069248       1 actuator.go:430] dofinn-20201705-blz22-worker-ap-southeast-2a-7pd2r: found 1 existing instances for machine
I0518 23:27:56.069258       1 actuator.go:459] dofinn-20201705-blz22-worker-ap-southeast-2a-7pd2r: found 1 running instances for machine
I0518 23:27:56.080843       1 actuator.go:185] dofinn-20201705-blz22-worker-ap-southeast-2a-7pd2r: ProviderID already set in the machine Spec with value:aws:///ap-southeast-2a/i-0488bee4f8178449b
I0518 23:27:56.080859       1 actuator.go:599] dofinn-20201705-blz22-worker-ap-southeast-2a-7pd2r: Updating status
I0518 23:27:56.080897       1 actuator.go:628] dofinn-20201705-blz22-worker-ap-southeast-2a-7pd2r: finished calculating AWS status
I0518 23:27:56.080978       1 actuator.go:233] dofinn-20201705-blz22-worker-ap-southeast-2a-7pd2r: status unchanged
I0518 23:27:56.081057       1 controller.go:161] Reconciling Machine "dofinn-20201705-blz22-infra-ap-southeast-2a-w82t2"
I0518 23:27:56.081099       1 actuator.go:500] dofinn-20201705-blz22-infra-ap-southeast-2a-w82t2: Checking if machine exists
I0518 23:27:56.175772       1 actuator.go:549] dofinn-20201705-blz22-infra-ap-southeast-2a-w82t2: Found instance by id: i-0e240367c7b5aef32
I0518 23:27:56.175799       1 controller.go:258] Reconciling machine "dofinn-20201705-blz22-infra-ap-southeast-2a-w82t2" triggers idempotent update
I0518 23:27:56.175809       1 actuator.go:406] dofinn-20201705-blz22-infra-ap-southeast-2a-w82t2: updating machine
I0518 23:27:56.175893       1 actuator.go:414] dofinn-20201705-blz22-infra-ap-southeast-2a-w82t2: obtaining EC2 client for region
I0518 23:27:56.237648       1 actuator.go:549] dofinn-20201705-blz22-infra-ap-southeast-2a-w82t2: Found instance by id: i-0e240367c7b5aef32
I0518 23:27:56.237669       1 actuator.go:430] dofinn-20201705-blz22-infra-ap-southeast-2a-w82t2: found 1 existing instances for machine
I0518 23:27:56.237681       1 actuator.go:459] dofinn-20201705-blz22-infra-ap-southeast-2a-w82t2: found 1 running instances for machine
I0518 23:27:56.250016       1 actuator.go:185] dofinn-20201705-blz22-infra-ap-southeast-2a-w82t2: ProviderID already set in the machine Spec with value:aws:///ap-southeast-2a/i-0e240367c7b5aef32
I0518 23:27:56.250032       1 actuator.go:599] dofinn-20201705-blz22-infra-ap-southeast-2a-w82t2: Updating status
I0518 23:27:56.250072       1 actuator.go:628] dofinn-20201705-blz22-infra-ap-southeast-2a-w82t2: finished calculating AWS status
I0518 23:27:56.250143       1 actuator.go:233] dofinn-20201705-blz22-infra-ap-southeast-2a-w82t2: status unchanged
I0518 23:27:56.250190       1 controller.go:161] Reconciling Machine "dofinn-20201705-blz22-infra-ap-southeast-2a-4pdbs"
I0518 23:27:56.250202       1 actuator.go:500] dofinn-20201705-blz22-infra-ap-southeast-2a-4pdbs: Checking if machine exists
I0518 23:27:56.358535       1 actuator.go:549] dofinn-20201705-blz22-infra-ap-southeast-2a-4pdbs: Found instance by id: i-031d2ef084bfc6894
I0518 23:27:56.358586       1 controller.go:258] Reconciling machine "dofinn-20201705-blz22-infra-ap-southeast-2a-4pdbs" triggers idempotent update
I0518 23:27:56.358595       1 actuator.go:406] dofinn-20201705-blz22-infra-ap-southeast-2a-4pdbs: updating machine
I0518 23:27:56.358668       1 actuator.go:414] dofinn-20201705-blz22-infra-ap-southeast-2a-4pdbs: obtaining EC2 client for region
I0518 23:27:56.416100       1 actuator.go:549] dofinn-20201705-blz22-infra-ap-southeast-2a-4pdbs: Found instance by id: i-031d2ef084bfc6894
I0518 23:27:56.416121       1 actuator.go:430] dofinn-20201705-blz22-infra-ap-southeast-2a-4pdbs: found 1 existing instances for machine
I0518 23:27:56.416130       1 actuator.go:459] dofinn-20201705-blz22-infra-ap-southeast-2a-4pdbs: found 1 running instances for machine
I0518 23:27:56.428609       1 actuator.go:185] dofinn-20201705-blz22-infra-ap-southeast-2a-4pdbs: ProviderID already set in the machine Spec with value:aws:///ap-southeast-2a/i-031d2ef084bfc6894
I0518 23:27:56.428628       1 actuator.go:599] dofinn-20201705-blz22-infra-ap-southeast-2a-4pdbs: Updating status
I0518 23:27:56.428660       1 actuator.go:628] dofinn-20201705-blz22-infra-ap-southeast-2a-4pdbs: finished calculating AWS status
I0518 23:27:56.428732       1 actuator.go:233] dofinn-20201705-blz22-infra-ap-southeast-2a-4pdbs: status unchanged
I0518 23:27:56.428785       1 controller.go:161] Reconciling Machine "dofinn-20201705-blz22-infra-ap-southeast-2a-ph26m"
I0518 23:27:56.428800       1 actuator.go:500] dofinn-20201705-blz22-infra-ap-southeast-2a-ph26m: Checking if machine exists
I0518 23:27:56.532138       1 actuator.go:549] dofinn-20201705-blz22-infra-ap-southeast-2a-ph26m: Found instance by id: i-030bf61b1e820df6a
I0518 23:27:56.532160       1 controller.go:258] Reconciling machine "dofinn-20201705-blz22-infra-ap-southeast-2a-ph26m" triggers idempotent update
I0518 23:27:56.532168       1 actuator.go:406] dofinn-20201705-blz22-infra-ap-southeast-2a-ph26m: updating machine
I0518 23:27:56.532244       1 actuator.go:414] dofinn-20201705-blz22-infra-ap-southeast-2a-ph26m: obtaining EC2 client for region
I0518 23:27:56.587025       1 actuator.go:549] dofinn-20201705-blz22-infra-ap-southeast-2a-ph26m: Found instance by id: i-030bf61b1e820df6a
I0518 23:27:56.587043       1 actuator.go:430] dofinn-20201705-blz22-infra-ap-southeast-2a-ph26m: found 1 existing instances for machine
I0518 23:27:56.587053       1 actuator.go:459] dofinn-20201705-blz22-infra-ap-southeast-2a-ph26m: found 1 running instances for machine
I0518 23:27:56.597858       1 actuator.go:185] dofinn-20201705-blz22-infra-ap-southeast-2a-ph26m: ProviderID already set in the machine Spec with value:aws:///ap-southeast-2a/i-030bf61b1e820df6a
I0518 23:27:56.597875       1 actuator.go:599] dofinn-20201705-blz22-infra-ap-southeast-2a-ph26m: Updating status
I0518 23:27:56.597923       1 actuator.go:628] dofinn-20201705-blz22-infra-ap-southeast-2a-ph26m: finished calculating AWS status
I0518 23:27:56.598006       1 actuator.go:233] dofinn-20201705-blz22-infra-ap-southeast-2a-ph26m: status unchanged
I0518 23:27:56.598063       1 controller.go:161] Reconciling Machine "dofinn-20201705-blz22-worker-ap-southeast-2a-xj2rc"
I0518 23:27:56.598076       1 actuator.go:500] dofinn-20201705-blz22-worker-ap-southeast-2a-xj2rc: Checking if machine exists
I0518 23:27:56.665630       1 actuator.go:549] dofinn-20201705-blz22-worker-ap-southeast-2a-xj2rc: Found instance by id: i-07b30d9dbde6351dc
I0518 23:27:56.665655       1 controller.go:258] Reconciling machine "dofinn-20201705-blz22-worker-ap-southeast-2a-xj2rc" triggers idempotent update
I0518 23:27:56.665663       1 actuator.go:406] dofinn-20201705-blz22-worker-ap-southeast-2a-xj2rc: updating machine
I0518 23:27:56.665738       1 actuator.go:414] dofinn-20201705-blz22-worker-ap-southeast-2a-xj2rc: obtaining EC2 client for region
I0518 23:27:56.733963       1 actuator.go:549] dofinn-20201705-blz22-worker-ap-southeast-2a-xj2rc: Found instance by id: i-07b30d9dbde6351dc
I0518 23:27:56.733983       1 actuator.go:430] dofinn-20201705-blz22-worker-ap-southeast-2a-xj2rc: found 1 existing instances for machine
I0518 23:27:56.745100       1 actuator.go:185] dofinn-20201705-blz22-worker-ap-southeast-2a-xj2rc: ProviderID already set in the machine Spec with value:aws:///ap-southeast-2a/i-07b30d9dbde6351dc
I0518 23:27:56.745134       1 actuator.go:599] dofinn-20201705-blz22-worker-ap-southeast-2a-xj2rc: Updating status
I0518 23:27:56.745177       1 actuator.go:628] dofinn-20201705-blz22-worker-ap-southeast-2a-xj2rc: finished calculating AWS status
I0518 23:27:56.745251       1 actuator.go:233] dofinn-20201705-blz22-worker-ap-southeast-2a-xj2rc: status unchanged
I0518 23:27:56.745299       1 controller.go:161] Reconciling Machine "dofinn-20201705-blz22-master-1"
I0518 23:27:56.745316       1 actuator.go:500] dofinn-20201705-blz22-master-1: Checking if machine exists
I0518 23:27:56.870845       1 actuator.go:549] dofinn-20201705-blz22-master-1: Found instance by id: i-07aedb09e9f984af0
I0518 23:27:56.870868       1 controller.go:258] Reconciling machine "dofinn-20201705-blz22-master-1" triggers idempotent update
I0518 23:27:56.870875       1 actuator.go:406] dofinn-20201705-blz22-master-1: updating machine
I0518 23:27:56.870950       1 actuator.go:414] dofinn-20201705-blz22-master-1: obtaining EC2 client for region
I0518 23:27:56.928814       1 actuator.go:549] dofinn-20201705-blz22-master-1: Found instance by id: i-07aedb09e9f984af0
I0518 23:27:56.928833       1 actuator.go:430] dofinn-20201705-blz22-master-1: found 1 existing instances for machine
I0518 23:27:56.928846       1 actuator.go:459] dofinn-20201705-blz22-master-1: found 1 running instances for machine
I0518 23:27:58.024956       1 actuator.go:185] dofinn-20201705-blz22-master-1: ProviderID already set in the machine Spec with value:aws:///ap-southeast-2a/i-07aedb09e9f984af0
I0518 23:27:58.024975       1 actuator.go:599] dofinn-20201705-blz22-master-1: Updating status
I0518 23:27:58.025009       1 actuator.go:628] dofinn-20201705-blz22-master-1: finished calculating AWS status
I0518 23:27:58.025078       1 actuator.go:233] dofinn-20201705-blz22-master-1: status unchanged
I0518 23:27:58.025126       1 controller.go:161] Reconciling Machine "dofinn-20201705-blz22-master-2"
I0518 23:27:58.025139       1 actuator.go:500] dofinn-20201705-blz22-master-2: Checking if machine exists
I0518 23:27:58.133375       1 actuator.go:549] dofinn-20201705-blz22-master-2: Found instance by id: i-0f04ac6744f0d8a03
I0518 23:27:58.133399       1 controller.go:258] Reconciling machine "dofinn-20201705-blz22-master-2" triggers idempotent update
I0518 23:27:58.133409       1 actuator.go:406] dofinn-20201705-blz22-master-2: updating machine
I0518 23:27:58.133499       1 actuator.go:414] dofinn-20201705-blz22-master-2: obtaining EC2 client for region
I0518 23:27:58.197014       1 actuator.go:549] dofinn-20201705-blz22-master-2: Found instance by id: i-0f04ac6744f0d8a03
I0518 23:27:58.197036       1 actuator.go:430] dofinn-20201705-blz22-master-2: found 1 existing instances for machine
I0518 23:27:58.197046       1 actuator.go:459] dofinn-20201705-blz22-master-2: found 1 running instances for machine
I0518 23:27:59.104148       1 actuator.go:185] dofinn-20201705-blz22-master-2: ProviderID already set in the machine Spec with value:aws:///ap-southeast-2a/i-0f04ac6744f0d8a03
I0518 23:27:59.104167       1 actuator.go:599] dofinn-20201705-blz22-master-2: Updating status
I0518 23:27:59.104206       1 actuator.go:628] dofinn-20201705-blz22-master-2: finished calculating AWS status
I0518 23:27:59.104278       1 actuator.go:233] dofinn-20201705-blz22-master-2: status unchanged
I0518 23:27:59.104325       1 controller.go:161] Reconciling Machine "dofinn-20201705-blz22-worker-ap-southeast-2a-rh8xv"
I0518 23:27:59.104339       1 actuator.go:500] dofinn-20201705-blz22-worker-ap-southeast-2a-rh8xv: Checking if machine exists
I0518 23:27:59.211246       1 actuator.go:549] dofinn-20201705-blz22-worker-ap-southeast-2a-rh8xv: Found instance by id: i-0725991df42113c47
I0518 23:27:59.211274       1 controller.go:258] Reconciling machine "dofinn-20201705-blz22-worker-ap-southeast-2a-rh8xv" triggers idempotent update
I0518 23:27:59.211284       1 actuator.go:406] dofinn-20201705-blz22-worker-ap-southeast-2a-rh8xv: updating machine
I0518 23:27:59.211369       1 actuator.go:414] dofinn-20201705-blz22-worker-ap-southeast-2a-rh8xv: obtaining EC2 client for region
I0518 23:27:59.272581       1 actuator.go:549] dofinn-20201705-blz22-worker-ap-southeast-2a-rh8xv: Found instance by id: i-0725991df42113c47
I0518 23:27:59.272602       1 actuator.go:430] dofinn-20201705-blz22-worker-ap-southeast-2a-rh8xv: found 1 existing instances for machine
I0518 23:27:59.272613       1 actuator.go:459] dofinn-20201705-blz22-worker-ap-southeast-2a-rh8xv: found 1 running instances for machine
I0518 23:27:59.283835       1 actuator.go:185] dofinn-20201705-blz22-worker-ap-southeast-2a-rh8xv: ProviderID already set in the machine Spec with value:aws:///ap-southeast-2a/i-0725991df42113c47
I0518 23:27:59.283853       1 actuator.go:599] dofinn-20201705-blz22-worker-ap-southeast-2a-rh8xv: Updating status
I0518 23:27:59.283892       1 actuator.go:628] dofinn-20201705-blz22-worker-ap-southeast-2a-rh8xv: finished calculating AWS status
I0518 23:27:59.283969       1 actuator.go:233] dofinn-20201705-blz22-worker-ap-southeast-2a-rh8xv: status unchanged
I0518 23:27:59.284039       1 controller.go:161] Reconciling Machine "dofinn-20201705-blz22-master-0"
I0518 23:27:59.284053       1 actuator.go:500] dofinn-20201705-blz22-master-0: Checking if machine exists
I0518 23:27:59.381608       1 actuator.go:549] dofinn-20201705-blz22-master-0: Found instance by id: i-070a0828b9fd1d0e7
I0518 23:27:59.381635       1 controller.go:258] Reconciling machine "dofinn-20201705-blz22-master-0" triggers idempotent update
I0518 23:27:59.381645       1 actuator.go:406] dofinn-20201705-blz22-master-0: updating machine
I0518 23:27:59.381733       1 actuator.go:414] dofinn-20201705-blz22-master-0: obtaining EC2 client for region
I0518 23:27:59.446986       1 actuator.go:549] dofinn-20201705-blz22-master-0: Found instance by id: i-070a0828b9fd1d0e7
I0518 23:27:59.447005       1 actuator.go:430] dofinn-20201705-blz22-master-0: found 1 existing instances for machine
I0518 23:27:59.447017       1 actuator.go:459] dofinn-20201705-blz22-master-0: found 1 running instances for machine
I0518 23:28:00.387301       1 actuator.go:185] dofinn-20201705-blz22-master-0: ProviderID already set in the machine Spec with value:aws:///ap-southeast-2a/i-070a0828b9fd1d0e7
I0518 23:28:00.387320       1 actuator.go:599] dofinn-20201705-blz22-master-0: Updating status
I0518 23:28:00.387359       1 actuator.go:628] dofinn-20201705-blz22-master-0: finished calculating AWS status
I0518 23:28:00.387428       1 actuator.go:233] dofinn-20201705-blz22-master-0: status unchanged
I0518 23:28:00.387478       1 controller.go:161] Reconciling Machine "dofinn-20201705-blz22-worker-ap-southeast-2a-7kvd4"
I0518 23:28:00.387487       1 actuator.go:500] dofinn-20201705-blz22-worker-ap-southeast-2a-7kvd4: Checking if machine exists
I0518 23:28:00.495594       1 actuator.go:549] dofinn-20201705-blz22-worker-ap-southeast-2a-7kvd4: Found instance by id: i-040be59194cd8ed64
I0518 23:28:00.495620       1 controller.go:258] Reconciling machine "dofinn-20201705-blz22-worker-ap-southeast-2a-7kvd4" triggers idempotent update
I0518 23:28:00.495628       1 actuator.go:406] dofinn-20201705-blz22-worker-ap-southeast-2a-7kvd4: updating machine
I0518 23:28:00.495706       1 actuator.go:414] dofinn-20201705-blz22-worker-ap-southeast-2a-7kvd4: obtaining EC2 client for region
I0518 23:28:00.549265       1 actuator.go:549] dofinn-20201705-blz22-worker-ap-southeast-2a-7kvd4: Found instance by id: i-040be59194cd8ed64
I0518 23:28:00.549287       1 actuator.go:430] dofinn-20201705-blz22-worker-ap-southeast-2a-7kvd4: found 1 existing instances for machine
I0518 23:28:00.549297       1 actuator.go:459] dofinn-20201705-blz22-worker-ap-southeast-2a-7kvd4: found 1 running instances for machine
I0518 23:28:00.559910       1 actuator.go:185] dofinn-20201705-blz22-worker-ap-southeast-2a-7kvd4: ProviderID already set in the machine Spec with value:aws:///ap-southeast-2a/i-040be59194cd8ed64
I0518 23:28:00.559928       1 actuator.go:599] dofinn-20201705-blz22-worker-ap-southeast-2a-7kvd4: Updating status
I0518 23:28:00.559960       1 actuator.go:628] dofinn-20201705-blz22-worker-ap-southeast-2a-7kvd4: finished calculating AWS status
I0518 23:28:00.560035       1 actuator.go:233] dofinn-20201705-blz22-worker-ap-southeast-2a-7kvd4: status unchanged
```

## Starting the same instance manually via the AWS Console. 

### Summary

The node becomes ready.

```
oc get machines -n openshift-machine-api | grep dofinn-20201705-blz22-worker-ap-southeast-2a-xj2rc                          130 ↵
dofinn-20201705-blz22-worker-ap-southeast-2a-xj2rc   Running   m5.xlarge   ap-southeast-2   ap-southeast-2a   10h
```

`openshift-machine-api` detects status change but fails to update state.

```
I0518 23:32:21.279779       1 controller.go:161] Reconciling Machine "dofinn-20201705-blz22-worker-ap-southeast-2a-xj2rc"
I0518 23:32:21.279893       1 actuator.go:500] dofinn-20201705-blz22-worker-ap-southeast-2a-xj2rc: Checking if machine exists
I0518 23:32:21.383933       1 actuator.go:549] dofinn-20201705-blz22-worker-ap-southeast-2a-xj2rc: Found instance by id: i-07b30d9dbde6351dc
I0518 23:32:21.383956       1 controller.go:258] Reconciling machine "dofinn-20201705-blz22-worker-ap-southeast-2a-xj2rc" triggers idempotent update
I0518 23:32:21.383967       1 actuator.go:406] dofinn-20201705-blz22-worker-ap-southeast-2a-xj2rc: updating machine
I0518 23:32:21.384036       1 actuator.go:414] dofinn-20201705-blz22-worker-ap-southeast-2a-xj2rc: obtaining EC2 client for region
I0518 23:32:21.459398       1 actuator.go:549] dofinn-20201705-blz22-worker-ap-southeast-2a-xj2rc: Found instance by id: i-07b30d9dbde6351dc
I0518 23:32:21.459421       1 actuator.go:430] dofinn-20201705-blz22-worker-ap-southeast-2a-xj2rc: found 1 existing instances for machine
I0518 23:32:21.459434       1 actuator.go:459] dofinn-20201705-blz22-worker-ap-southeast-2a-xj2rc: found 1 running instances for machine
I0518 23:32:21.471855       1 actuator.go:185] dofinn-20201705-blz22-worker-ap-southeast-2a-xj2rc: ProviderID already set in the machine Spec with value:aws:///ap-southeast-2a/i-07b30d9dbde6351dc
I0518 23:32:21.471877       1 actuator.go:599] dofinn-20201705-blz22-worker-ap-southeast-2a-xj2rc: Updating status
I0518 23:32:21.471918       1 actuator.go:628] dofinn-20201705-blz22-worker-ap-southeast-2a-xj2rc: finished calculating AWS status
I0518 23:32:21.471982       1 actuator.go:224] dofinn-20201705-blz22-worker-ap-southeast-2a-xj2rc: machine status has changed, updating
I0518 23:32:21.483980       1 controller.go:161] Reconciling Machine "dofinn-20201705-blz22-worker-ap-southeast-2a-xj2rc"
I0518 23:32:21.484069       1 actuator.go:500] dofinn-20201705-blz22-worker-ap-southeast-2a-xj2rc: Checking if machine exists
I0518 23:32:21.550936       1 actuator.go:549] dofinn-20201705-blz22-worker-ap-southeast-2a-xj2rc: Found instance by id: i-07b30d9dbde6351dc
I0518 23:32:21.550962       1 controller.go:258] Reconciling machine "dofinn-20201705-blz22-worker-ap-southeast-2a-xj2rc" triggers idempotent update
I0518 23:32:21.550972       1 actuator.go:406] dofinn-20201705-blz22-worker-ap-southeast-2a-xj2rc: updating machine
I0518 23:32:21.551063       1 actuator.go:414] dofinn-20201705-blz22-worker-ap-southeast-2a-xj2rc: obtaining EC2 client for region
I0518 23:32:21.620140       1 actuator.go:549] dofinn-20201705-blz22-worker-ap-southeast-2a-xj2rc: Found instance by id: i-07b30d9dbde6351dc
I0518 23:32:21.620162       1 actuator.go:430] dofinn-20201705-blz22-worker-ap-southeast-2a-xj2rc: found 1 existing instances for machine
I0518 23:32:21.620171       1 actuator.go:459] dofinn-20201705-blz22-worker-ap-southeast-2a-xj2rc: found 1 running instances for machine
E0518 23:32:21.624177       1 controller.go:260] Error updating machine "openshift-machine-api/dofinn-20201705-blz22-worker-ap-southeast-2a-xj2rc": dofinn-20201705-blz22-worker-ap-southeast-2a-xj2rc: failed to set machine cloud provider specifics: dofinn-20201705-blz22-worker-ap-southeast-2a-xj2rc: error updating machine spec: Operation cannot be fulfilled on machines.machine.openshift.io "dofinn-20201705-blz22-worker-ap-southeast-2a-xj2rc": the object has been modified; please apply your changes to the latest version and try again
I0518 23:32:22.624388       1 controller.go:161] Reconciling Machine "dofinn-20201705-blz22-worker-ap-southeast-2a-xj2rc"
I0518 23:32:22.624420       1 actuator.go:500] dofinn-20201705-blz22-worker-ap-southeast-2a-xj2rc: Checking if machine exists
I0518 23:32:22.693730       1 actuator.go:549] dofinn-20201705-blz22-worker-ap-southeast-2a-xj2rc: Found instance by id: i-07b30d9dbde6351dc
I0518 23:32:22.693755       1 controller.go:258] Reconciling machine "dofinn-20201705-blz22-worker-ap-southeast-2a-xj2rc" triggers idempotent update
I0518 23:32:22.693763       1 actuator.go:406] dofinn-20201705-blz22-worker-ap-southeast-2a-xj2rc: updating machine
I0518 23:32:22.693831       1 actuator.go:414] dofinn-20201705-blz22-worker-ap-southeast-2a-xj2rc: obtaining EC2 client for region
I0518 23:32:22.760081       1 actuator.go:549] dofinn-20201705-blz22-worker-ap-southeast-2a-xj2rc: Found instance by id: i-07b30d9dbde6351dc
I0518 23:32:22.760101       1 actuator.go:430] dofinn-20201705-blz22-worker-ap-southeast-2a-xj2rc: found 1 existing instances for machine
I0518 23:32:22.760113       1 actuator.go:459] dofinn-20201705-blz22-worker-ap-southeast-2a-xj2rc: found 1 running instances for machine
I0518 23:32:22.770807       1 actuator.go:185] dofinn-20201705-blz22-worker-ap-southeast-2a-xj2rc: ProviderID already set in the machine Spec with value:aws:///ap-southeast-2a/i-07b30d9dbde6351dc
I0518 23:32:22.770824       1 actuator.go:599] dofinn-20201705-blz22-worker-ap-southeast-2a-xj2rc: Updating status
I0518 23:32:22.770860       1 actuator.go:628] dofinn-20201705-blz22-worker-ap-southeast-2a-xj2rc: finished calculating AWS status
I0518 23:32:22.770931       1 actuator.go:233] dofinn-20201705-blz22-worker-ap-southeast-2a-xj2rc: status unchanged
```

## Testing Reboot via the AWS Console

### Summary

`openshift-machine-api` detects a change, but not status are updated. 


machine-api logs

```
I0518 23:36:13.424011       1 controller.go:161] Reconciling Machine "dofinn-20201705-blz22-worker-ap-southeast-2a-xj2rc"
I0518 23:36:13.426587       1 actuator.go:500] dofinn-20201705-blz22-worker-ap-southeast-2a-xj2rc: Checking if machine exists
I0518 23:36:18.522620       1 actuator.go:549] dofinn-20201705-blz22-worker-ap-southeast-2a-xj2rc: Found instance by id: i-07b30d9dbde6351dc
I0518 23:36:18.522644       1 controller.go:258] Reconciling machine "dofinn-20201705-blz22-worker-ap-southeast-2a-xj2rc" triggers idempotent update
I0518 23:36:18.522654       1 actuator.go:406] dofinn-20201705-blz22-worker-ap-southeast-2a-xj2rc: updating machine
I0518 23:36:18.522733       1 actuator.go:414] dofinn-20201705-blz22-worker-ap-southeast-2a-xj2rc: obtaining EC2 client for region
I0518 23:36:18.589743       1 actuator.go:549] dofinn-20201705-blz22-worker-ap-southeast-2a-xj2rc: Found instance by id: i-07b30d9dbde6351dc
I0518 23:36:18.589764       1 actuator.go:430] dofinn-20201705-blz22-worker-ap-southeast-2a-xj2rc: found 1 existing instances for machine
I0518 23:36:18.589777       1 actuator.go:459] dofinn-20201705-blz22-worker-ap-southeast-2a-xj2rc: found 1 running instances for machine
I0518 23:36:18.601155       1 actuator.go:185] dofinn-20201705-blz22-worker-ap-southeast-2a-xj2rc: ProviderID already set in the machine Spec with value:aws:///ap-southeast-2a/i-07b30d9dbde6351dc
I0518 23:36:18.601172       1 actuator.go:599] dofinn-20201705-blz22-worker-ap-southeast-2a-xj2rc: Updating status
I0518 23:36:18.601211       1 actuator.go:628] dofinn-20201705-blz22-worker-ap-southeast-2a-xj2rc: finished calculating AWS status
I0518 23:36:18.601285       1 actuator.go:233] dofinn-20201705-blz22-worker-ap-southeast-2a-xj2rc: status unchanged
```

## Testing Terminate via the AWS Console

### Summary
Node becomes `NotReady` prior being deleted from `oc get nodes`
`machine-api` detects a change, but still sees the machine and keeps it in a running state. 
machineset remains unchanged illustrating a false indication of machine status. 
Required manual intervention of deleting the target machines `machine` CR. 
The machine-api then created a new machine to enable the cluster to return to the correct state.

```
oc get machines -n openshift-machine-api  | grep dofinn-20201705-blz22-worker-ap-southeast-2a-xj2rc     
dofinn-20201705-blz22-worker-ap-southeast-2a-xj2rc   Running   m5.xlarge   ap-southeast-2   ap-southeast-2a   10h
```

machine sets are also not updated. 

```
oc get machinesets -n openshift-machine-api
NAME                                           DESIRED   CURRENT   READY   AVAILABLE   AGE
dofinn-20201705-blz22-infra-ap-southeast-2a    3         3         3       3           46h
dofinn-20201705-blz22-worker-ap-southeast-2a   4         4         3       3           46h
```


machine-api logs

```
I0518 23:39:43.660206       1 controller.go:161] Reconciling Machine "dofinn-20201705-blz22-worker-ap-southeast-2a-xj2rc"
I0518 23:39:43.660492       1 actuator.go:500] dofinn-20201705-blz22-worker-ap-southeast-2a-xj2rc: Checking if machine exists
I0518 23:39:43.759147       1 actuator.go:549] dofinn-20201705-blz22-worker-ap-southeast-2a-xj2rc: Found instance by id: i-07b30d9dbde6351dc
I0518 23:39:43.759246       1 controller.go:258] Reconciling machine "dofinn-20201705-blz22-worker-ap-southeast-2a-xj2rc" triggers idempotent update
I0518 23:39:43.759281       1 actuator.go:406] dofinn-20201705-blz22-worker-ap-southeast-2a-xj2rc: updating machine
I0518 23:39:43.759397       1 actuator.go:414] dofinn-20201705-blz22-worker-ap-southeast-2a-xj2rc: obtaining EC2 client for region
I0518 23:39:43.828414       1 actuator.go:549] dofinn-20201705-blz22-worker-ap-southeast-2a-xj2rc: Found instance by id: i-07b30d9dbde6351dc
I0518 23:39:43.828442       1 actuator.go:430] dofinn-20201705-blz22-worker-ap-southeast-2a-xj2rc: found 1 existing instances for machine
I0518 23:39:43.846922       1 actuator.go:185] dofinn-20201705-blz22-worker-ap-southeast-2a-xj2rc: ProviderID already set in the machine Spec with value:aws:///ap-southeast-2a/i-07b30d9dbde6351dc
I0518 23:39:43.846942       1 actuator.go:599] dofinn-20201705-blz22-worker-ap-southeast-2a-xj2rc: Updating status
I0518 23:39:43.846987       1 actuator.go:628] dofinn-20201705-blz22-worker-ap-southeast-2a-xj2rc: finished calculating AWS status
I0518 23:39:43.847052       1 actuator.go:224] dofinn-20201705-blz22-worker-ap-southeast-2a-xj2rc: machine status has changed, updating
I0518 23:39:43.864032       1 controller.go:161] Reconciling Machine "dofinn-20201705-blz22-worker-ap-southeast-2a-xj2rc"
I0518 23:39:43.864134       1 actuator.go:500] dofinn-20201705-blz22-worker-ap-southeast-2a-xj2rc: Checking if machine exists
I0518 23:39:43.930875       1 actuator.go:549] dofinn-20201705-blz22-worker-ap-southeast-2a-xj2rc: Found instance by id: i-07b30d9dbde6351dc
I0518 23:39:43.930900       1 controller.go:258] Reconciling machine "dofinn-20201705-blz22-worker-ap-southeast-2a-xj2rc" triggers idempotent update
I0518 23:39:43.930908       1 actuator.go:406] dofinn-20201705-blz22-worker-ap-southeast-2a-xj2rc: updating machine
I0518 23:39:43.930979       1 actuator.go:414] dofinn-20201705-blz22-worker-ap-southeast-2a-xj2rc: obtaining EC2 client for region
I0518 23:39:43.994345       1 actuator.go:549] dofinn-20201705-blz22-worker-ap-southeast-2a-xj2rc: Found instance by id: i-07b30d9dbde6351dc
I0518 23:39:43.994364       1 actuator.go:430] dofinn-20201705-blz22-worker-ap-southeast-2a-xj2rc: found 1 existing instances for machine
I0518 23:39:44.005316       1 actuator.go:185] dofinn-20201705-blz22-worker-ap-southeast-2a-xj2rc: ProviderID already set in the machine Spec with value:aws:///ap-southeast-2a/i-07b30d9dbde6351dc
I0518 23:39:44.005337       1 actuator.go:599] dofinn-20201705-blz22-worker-ap-southeast-2a-xj2rc: Updating status
I0518 23:39:44.005383       1 actuator.go:628] dofinn-20201705-blz22-worker-ap-southeast-2a-xj2rc: finished calculating AWS status
I0518 23:39:44.005469       1 actuator.go:233] dofinn-20201705-blz22-worker-ap-southeast-2a-xj2rc: status unchanged
```
