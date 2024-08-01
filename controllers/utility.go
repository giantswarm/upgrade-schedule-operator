package controllers

import (
	"time"

	"github.com/blang/semver"
	"github.com/giantswarm/k8smetadata/pkg/annotation"
	"github.com/giantswarm/k8smetadata/pkg/label"
	clusterv1 "sigs.k8s.io/cluster-api/api/v1beta1"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
)

const (
	ClusterUpgradeAnnouncement = "alpha.giantswarm.io/update-schedule-upgrade-announcement"
	OutOfHoursContact          = "kaascloud@giantswarm.io"
)

func defaultRequeue() reconcile.Result {
	return ctrl.Result{
		Requeue:      true,
		RequeueAfter: time.Minute * 5,
	}
}

func timedRequeue(upgradeTime time.Time) reconcile.Result {
	if upgradeTime.Sub(time.Now().UTC()) > 5*time.Minute {
		return defaultRequeue()
	}
	return ctrl.Result{
		Requeue:      true,
		RequeueAfter: upgradeTime.Sub(time.Now().UTC()) + time.Second,
	}
}

func getClusterReleaseVersionLabel(cluster *clusterv1.Cluster) string {
	labels := cluster.GetLabels()
	return labels[label.ReleaseVersion]
}

func isCAPIProvider(cluster *clusterv1.Cluster) bool {
	labels := cluster.GetLabels()
	name, ok := labels["cluster.x-k8s.io/watch-filter"]
	if !ok {
		return false
	}
	if name == "capi" {
		return true
	}
	return false
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
	return upgradeTime.Before(time.Now().UTC())
}

func upgradeAnnouncementTimeReached(upgradeTime time.Time) bool {
	return upgradeTime.Add(-15 * time.Minute).Before(time.Now().UTC())
}

func outOfOffice(upgradeTime time.Time) bool {
	if upgradeTime.Weekday() == time.Saturday || upgradeTime.Weekday() == time.Sunday {
		return true
	}
	if upgradeTime.UTC().Hour() <= 7 || upgradeTime.UTC().Hour() >= 16 {
		return true
	}
	return false
}
