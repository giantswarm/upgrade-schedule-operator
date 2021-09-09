package controllers

import (
	"time"

	"github.com/blang/semver"
	"github.com/giantswarm/apiextensions/v3/pkg/annotation"
	"github.com/giantswarm/apiextensions/v3/pkg/label"
	clusterv1 "sigs.k8s.io/cluster-api/api/v1alpha3"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
)

func defaultRequeue() reconcile.Result {
	return ctrl.Result{
		Requeue:      true,
		RequeueAfter: time.Minute * 5,
	}
}

func timedRequeue(upgradeTime time.Time) reconcile.Result {
	if upgradeTime.Sub(time.Now().In(upgradeTime.Location())) > 5*time.Minute {
		return defaultRequeue()
	}
	return ctrl.Result{
		Requeue:      true,
		RequeueAfter: upgradeTime.Sub(time.Now().In(upgradeTime.Location())) + time.Second,
	}
}

func getClusterReleaseVersionLabel(cluster *clusterv1.Cluster) string {
	labels := cluster.GetLabels()
	return labels[label.ReleaseVersion]
}

func getClusterUpgradeTimeAnnotation(cluster *clusterv1.Cluster) string {
	annotations := cluster.GetAnnotations()
	return annotations[annotation.UpdateScheduleTargetTime]
}

func getClusterUpgradeVersionAnnotation(cluster *clusterv1.Cluster) string {
	annotations := cluster.GetAnnotations()
	return annotations[annotation.UpdateScheduleTargetRelease]
}

func upgradeApplied(targetVersion semver.Version, currentVersion semver.Version) bool {
	return currentVersion.GE(targetVersion)
}

func upgradeTimeReached(upgradeTime time.Time) bool {
	return upgradeTime.Before(time.Now().In(upgradeTime.Location()))
}
