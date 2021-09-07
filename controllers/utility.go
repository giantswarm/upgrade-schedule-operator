package controllers

import (
	"time"

	"github.com/blang/semver"
	"github.com/giantswarm/apiextensions/v3/pkg/annotation"
	"github.com/giantswarm/apiextensions/v3/pkg/label"
	clusterv1 "sigs.k8s.io/cluster-api/api/v1alpha3"
)

func getClusterReleaseVersionLabel(cluster *clusterv1.Cluster) string {
	labels := cluster.GetLabels()
	return labels[label.ReleaseVersion]
}

func getClusterUpgradeTimeAnnotation(cluster *clusterv1.Cluster) string {
	annotations := cluster.GetAnnotations()
	return annotations[annotation.AWSUpdateScheduleTargetTime]
}

func getClusterUpgradeVersionAnnotation(cluster *clusterv1.Cluster) string {
	annotations := cluster.GetAnnotations()
	return annotations[annotation.AWSUpdateScheduleTargetRelease]
}

func upgradeApplied(targetVersion semver.Version, currentVersion semver.Version) bool {
	return currentVersion.GE(targetVersion)
}

func upgradeTimeReached(upgradeTime time.Time) bool {
	return upgradeTime.After(time.Now())
}
