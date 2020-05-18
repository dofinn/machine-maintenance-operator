# Machine Maintenance Operator

## What problem does it solve?
Currenty the `machine-config-api` does not handle machine termination via the scheduled events from AWS. This creates alerts that require manual intervention by deleting the missing machines `machine` CR. 

This will not scale. 

## What it does
* Queries AWS API for each instance ID in the cluster for scheduled events every 60 minutes. 
* Creates a `machinemaintenance` CR that the machinemaintenance controller reconciles
* Current reconciliation is deleting the target machines `machine` CR in the `openshift-machine-api` namespace
* machine-config-api gracefully terminates node

## What could do
* Query Xchangewindows CRs that are being developed in the managed-upgrade-operator.
* If its not a suitable time, exit reconcile loop with Result.Requeue = false
* Recncile will check again in 15 mins as per SyncPeriod. 
* If it is a suitable time, delete `machine` CR for target machine. 
* This should also handle different terminations for worker/infra vs master nodes. 

## Open Questions
From testing, it takes approximately 7 minutes for a node to be available in the "Running" state. If this is not acceptable, the reconciler could have the logic to increase machine pool by count+1, then terminate then reduce machine pool by count-1.
