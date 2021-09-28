/*
Copyright 2021.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package controllers

import (
	"context"
	"fmt"
	"time"

	"github.com/blang/semver"
	"github.com/giantswarm/apiextensions/v3/pkg/annotation"
	"github.com/giantswarm/apiextensions/v3/pkg/label"

	"github.com/go-logr/logr"
	"github.com/pkg/errors"
	corev1 "k8s.io/api/core/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/tools/record"
	clusterv1 "sigs.k8s.io/cluster-api/api/v1alpha3"
	"sigs.k8s.io/cluster-api/util/annotations"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

// ClusterReconciler reconciles a Cluster object
type ClusterReconciler struct {
	client.Client
	Log    logr.Logger
	Scheme *runtime.Scheme

	recorder record.EventRecorder
}

//+kubebuilder:rbac:groups=cluster.x-k8s.io,resources=clusters,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=cluster.x-k8s.io,resources=clusters/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=cluster.x-k8s.io,resources=clusters/finalizers,verbs=update

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the Cluster object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.8.3/pkg/reconcile
func (r *ClusterReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	log := r.Log.WithValues("cluster", req.NamespacedName)

	// Fetch the Cluster instance.
	cluster := &clusterv1.Cluster{}
	if err := r.Client.Get(ctx, req.NamespacedName, cluster); err != nil {
		if apierrors.IsNotFound(err) {
			// Object not found, return.  Created objects are automatically garbage collected.
			// For additional cleanup logic use finalizers.
			UpgradesInfo.WithLabelValues(cluster.Name, cluster.Namespace).Set(0)
			return ctrl.Result{}, nil
		}

		// Error reading the object - requeue the request.
		UpgradesInfo.WithLabelValues(cluster.Name, cluster.Namespace).Set(-1)
		return ctrl.Result{}, err
	}

	// Return if the Cluster is paused.
	if annotations.IsPaused(cluster, cluster) {
		log.Info("The cluster is paused.")
		UpgradesInfo.WithLabelValues(cluster.Name, cluster.Namespace).Set(0)
		return defaultRequeue(), nil
	}

	// Return if the Cluster is deleted.
	if !cluster.DeletionTimestamp.IsZero() {
		log.Info("The cluster is deleted.")
		UpgradesInfo.WithLabelValues(cluster.Name, cluster.Namespace).Set(0)
		return ctrl.Result{}, nil
	}

	// Return if there is no upgrade time scheduled.
	if getClusterUpgradeTimeAnnotation(cluster) == "" {
		log.Info("The cluster has no upgrade scheduled.")
		UpgradesInfo.WithLabelValues(cluster.Name, cluster.Namespace).Set(0)
		return defaultRequeue(), nil
	}

	// Return if the upgrade release version is not specified.
	if getClusterUpgradeVersionAnnotation(cluster) == "" {
		log.Info(fmt.Sprintf("The scheduled update at %v can not proceed because no target release version has been set via annotation %v.", getClusterUpgradeTimeAnnotation(cluster), annotation.UpdateScheduleTargetRelease))
		UpgradesInfo.WithLabelValues(cluster.Name, cluster.Namespace).Set(0)
		return defaultRequeue(), nil
	}
	return r.ReconcileUpgrade(ctx, cluster, log)
}

func (r *ClusterReconciler) ReconcileUpgrade(ctx context.Context, cluster *clusterv1.Cluster, log logr.Logger) (ctrl.Result, error) {
	upgradeTime, err := time.Parse(time.RFC822, getClusterUpgradeTimeAnnotation(cluster))
	if err != nil {
		log.Error(err, fmt.Sprintf("Failed to parse cluster upgrade time annotation %v. The value has to be in RFC822 Format and UTC time zone. e.g. 30 Jan 21 15:04 UTC", getClusterUpgradeTimeAnnotation(cluster)))
		UpgradesInfo.WithLabelValues(cluster.Name, cluster.Namespace).Set(-1)
		return ctrl.Result{}, err
	}

	// Send scheduled cluster upgrade announcement.
	if _, exists := cluster.Annotations[ClusterUpgradeAnnouncement]; !exists {
		if upgradeAnnouncementTimeReached(upgradeTime) {
			cluster.Annotations[ClusterUpgradeAnnouncement] = "true"
			err = r.Client.Update(ctx, cluster)
			if err != nil {
				log.Error(err, "Failed to set upgrade announcement annotation.")
				UpgradesInfo.WithLabelValues(cluster.Name, cluster.Namespace).Set(-1)
				return ctrl.Result{}, err
			}
			log.Info("Sending cluster upgrade announcement event.")

			msg := fmt.Sprintf("The cluster %s/%s upgrade from release version %v to %v is scheduled to start in %v.",
				cluster.Namespace,
				cluster.Name,
				getClusterReleaseVersionLabel(cluster),
				getClusterUpgradeVersionAnnotation(cluster),
				upgradeTime.Sub(time.Now().UTC()).Round(time.Minute),
			)
			if outOfOffice(upgradeTime) {
				msg += fmt.Sprintf(" Please contact us via %s in case of anomalies.", OutOfHoursContact)
			}
			r.sendClusterUpgradeEvent(cluster, msg)
		}
	}

	// Return if the scheduled upgrade time is not reached yet.
	if !upgradeTimeReached(upgradeTime) {
		log.Info(fmt.Sprintf("The scheduled update time is not reached yet. Cluster will be upgraded in %v at %v.", upgradeTime.Sub(time.Now().UTC()).Round(time.Minute), upgradeTime))
		UpgradesInfo.WithLabelValues(cluster.Name, cluster.Namespace).Set(float64(upgradeTime.Unix()))
		return timedRequeue(upgradeTime), nil
	}

	currentVersion, err := semver.New(getClusterReleaseVersionLabel(cluster))
	if err != nil {
		log.Error(err, "Failed to parse current cluster release version label.")
		UpgradesInfo.WithLabelValues(cluster.Name, cluster.Namespace).Set(-1)
		return ctrl.Result{}, err
	}
	targetVersion, err := semver.New(getClusterUpgradeVersionAnnotation(cluster))
	if err != nil {
		log.Error(err, fmt.Sprintf("Failed to parse cluster upgrade target version annotation %v. The value has to be only the desired release version, e.g 15.2.1.", getClusterUpgradeVersionAnnotation(cluster)))
		UpgradesInfo.WithLabelValues(cluster.Name, cluster.Namespace).Set(-1)
		return ctrl.Result{}, err
	}

	// Return if the upgrade to the target release has already been performed.
	if upgradeApplied(*targetVersion, *currentVersion) {
		log.Info(fmt.Sprintf("The upgrade to target version %v has already been applied. The current release version is %v.", targetVersion, currentVersion))
		UpgradesInfo.WithLabelValues(cluster.Name, cluster.Namespace).Set(0)
		return defaultRequeue(), nil
	}

	// Apply the upgrade and remove annotations
	log.Info(fmt.Sprintf("The cluster will be upgraded from version %v to %v.", currentVersion, targetVersion))
	cluster.Labels[label.ReleaseVersion] = getClusterUpgradeVersionAnnotation(cluster)
	delete(cluster.Annotations, annotation.UpdateScheduleTargetTime)
	delete(cluster.Annotations, annotation.UpdateScheduleTargetRelease)
	delete(cluster.Annotations, ClusterUpgradeAnnouncement)
	UpgradesTotal.WithLabelValues(cluster.Name, cluster.Namespace, currentVersion.String(), targetVersion.String()).Inc()
	err = r.Client.Update(ctx, cluster)
	if err != nil {
		log.Error(err, "Failed to update Release version tag and remove scheduled upgrade annotations.")
		FailuresTotal.WithLabelValues(cluster.Name, cluster.Namespace, currentVersion.String(), targetVersion.String()).Inc()
		UpgradesInfo.WithLabelValues(cluster.Name, cluster.Namespace).Set(-1)
		return ctrl.Result{}, err
	}
	log.Info(fmt.Sprintf("The cluster CR was upgraded from version %v to %v.", currentVersion, targetVersion))
	SuccessTotal.WithLabelValues(cluster.Name, cluster.Namespace, currentVersion.String(), targetVersion.String()).Inc()
	UpgradesInfo.WithLabelValues(cluster.Name, cluster.Namespace).Set(0)
	return defaultRequeue(), nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *ClusterReconciler) SetupWithManager(mgr ctrl.Manager) error {
	err := ctrl.NewControllerManagedBy(mgr).
		For(&clusterv1.Cluster{}).
		Complete(r)
	if err != nil {
		return errors.Wrap(err, "failed setting up with a controller manager")
	}

	r.recorder = mgr.GetEventRecorderFor("cluster-controller")
	return nil
}

func (r *ClusterReconciler) sendClusterUpgradeEvent(cluster *clusterv1.Cluster, message string) {
	r.recorder.Eventf(cluster, corev1.EventTypeNormal, "ClusterUpgradeAnnouncement", message)
}
