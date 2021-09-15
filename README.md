[![CircleCI](https://circleci.com/gh/giantswarm/upgrade-schedule-operator.svg?style=shield)](https://circleci.com/gh/giantswarm/upgrade-schedule-operator)

# upgrade-schedule-operator

This operator is intended to automate scheduling, triggering and notification for upgrades of giant swarm workload clusters.

## prerequisites

Running this operator alone on a management cluster will enable the automatic scheduling and triggering of upgrades of workload clusters in any provider.
However, on AWS, we use the following components together for observability and avoiding errors:

- [`event-exporter-app`](https://github.com/giantswarm/event-exporter-app) for exporting the upgrade events and notifying stakeholders in their slack channel.
  This requires adding a token for the slack channels (using `Giant Swarm Cluster Upgrade` app in slack) to the MC config if it is not there yet.
- [`aws-admission-controller`](https://github.com/giantswarm/aws-admission-controller) to validate the format of the upgrade annotations.
  This could be added to another validating webhook.
- [`cluster-operator`](https://github.com/giantswarm/cluster-operator) that carries out the actual upgrading process and emits upgrade events on the cluster.
  This could be added to another controller.

## how to schedule the upgrade

To schedule your upgrade, simply add the following annotations to the `Cluster` CR to specify the time and desired version.
```
annotations:
  alpha.giantswarm.io/update-schedule-target-release: 15.2.1
  alpha.giantswarm.io/update-schedule-target-time: "15 Sep 21 08:00 UTC"
```
Please note that the release version has to be an existing release higher than the current release version.
The time has to be given in RFC822 format and UTC.
Furthermore, only times that are at least 16 minutes in the future but not more than 6 months are accepted.
(16 minutes to ensure that a notification about the upgrade can be sent in advance)

## what happens next

For your scheduled upgrades you should be able to see the remaining time in the logs of the `upgrade-schedule-operator`

```
2021-09-14T17:23:29.605Z	INFO	controllers.Cluster	The scheduled update time is not reached yet. Cluster will be upgraded in 14h37m0s at 2021-09-15 08:00:00 +0000 UTC.	{"cluster": "default/xyz01"}
```
Around 10-15 minutes before the scheduled upgrade, a slack message should appear in the specified slack channel.
```
Giant Swarm Cluster Upgrade (APP)  9:48 AM
The cluster default/xyz01 upgrade from release version 15.1.0 to 15.2.1 is scheduled to start in 12m0s.
```
At the scheduled upgrade time, you should see the release version label change on the `Cluster` CR and a second slack message.
```
Workload cluster upgrade triggered for default/xyz01 on gauss.
```

## debugging

Generally take the same precautions/actions you would as when you trigger the upgrade manually. Some additional advice:

- If for any reason you need to cancel the scheduled upgrade, just remove (one of) the annotations.
  Or change them to reschedule.

- If something does not go as expected, always check the `upgrade-schedule-operator` logs for clues.

- To observe the process before/during the upgrade we recommend to look at the `aws cluster status` dashboard on the MC occasionally.
  ```
  opsctl open -a grafana -i alpaca
  ```
  Search for `aws cluster status` and don't forget to select your cluster from the drop down.

## observability

The operator exposes a couple of prometheus metrics.

- `scheduled_upgrades_applied_total`: the total number of times an upgrade was attempted to apply.
  this is counted by cluster as well as target and origin version.
- `scheduled_upgrades_failed_total`: the total number of times an upgrade was attempted to apply but failed.
  this is counted by cluster as well as target and origin version.
- `scheduled_upgrades_succeeded_total`: the total number of times an upgrade was attempted to apply and succeeded.
  this is counted by cluster as well as target and origin version.
- `scheduled_upgrades_time`: the scheduled upgrade time for each cluster in unix format.
  In case a cluster has no scheduled upgrade it will be 0.
  In case there is some sort of error with the upgrade it will be -1.