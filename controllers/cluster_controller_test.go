package controllers

import (
	"context"
	"fmt"
	"strconv"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	utilruntime "k8s.io/apimachinery/pkg/util/runtime"
	clientgoscheme "k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/tools/record"

	capi "sigs.k8s.io/cluster-api/api/v1beta1"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client/fake"
)

var (
	fakeScheme = runtime.NewScheme()
)

func init() {
	utilruntime.Must(clientgoscheme.AddToScheme(fakeScheme))
	_ = capi.AddToScheme(fakeScheme)
}

func TestClusterController(t *testing.T) {
	testCases := []struct {
		name                   string
		expectedReleaseVersion string
		expectedEventTriggered bool
		annotationsKept        bool

		cluster   *capi.Cluster
		configMap *corev1.ConfigMap
	}{
		// event triggered, within office time
		{
			name:                   "case 0",
			expectedReleaseVersion: "15.2.1",
			annotationsKept:        false,
			expectedEventTriggered: true,
			cluster: &capi.Cluster{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "test",
					Namespace: "default",
					Labels: map[string]string{
						"cluster-operator.giantswarm.io/version": "3.7.0",
						"cluster.x-k8s.io/cluster-name":          "dh82p",
						"giantswarm.io/cluster":                  "dh82p",
						"giantswarm.io/organization":             "giantswarm",
						"release.giantswarm.io/version":          "14.2.2",
					},
					Annotations: map[string]string{
						"alpha.giantswarm.io/update-schedule-target-release": "15.2.1",
						"alpha.giantswarm.io/update-schedule-target-time":    "10 Sep 21 12:00 UTC",
					},
				},
			},
		},
		// event triggered, weekend
		{
			name:                   "case 1",
			expectedReleaseVersion: "15.2.1",
			expectedEventTriggered: true,
			annotationsKept:        false,
			cluster: &capi.Cluster{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "test",
					Namespace: "default",
					Labels: map[string]string{
						"cluster-operator.giantswarm.io/version": "3.7.0",
						"cluster.x-k8s.io/cluster-name":          "ga83x",
						"giantswarm.io/cluster":                  "ga83x",
						"giantswarm.io/organization":             "giantswarm",
						"release.giantswarm.io/version":          "14.2.2",
					},
					Annotations: map[string]string{
						"alpha.giantswarm.io/update-schedule-target-release": "15.2.1",
						"alpha.giantswarm.io/update-schedule-target-time":    "11 Sep 21 12:00 UTC",
					},
				},
			},
		},
		// event triggered, out of office time
		{
			name:                   "case 2",
			expectedReleaseVersion: "15.2.1",
			expectedEventTriggered: true,
			annotationsKept:        false,
			cluster: &capi.Cluster{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "test",
					Namespace: "default",
					Labels: map[string]string{
						"cluster-operator.giantswarm.io/version": "3.7.0",
						"cluster.x-k8s.io/cluster-name":          "ga83x",
						"giantswarm.io/cluster":                  "ga83x",
						"giantswarm.io/organization":             "giantswarm",
						"release.giantswarm.io/version":          "14.2.2",
					},
					Annotations: map[string]string{
						"alpha.giantswarm.io/update-schedule-target-release": "15.2.1",
						"alpha.giantswarm.io/update-schedule-target-time":    "13 Sep 21 19:00 UTC",
					},
				},
			},
		},
		// future event, no change
		{
			name:                   "case 3",
			expectedReleaseVersion: "14.2.2",
			expectedEventTriggered: false,
			annotationsKept:        true,
			cluster: &capi.Cluster{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "test2",
					Namespace: "default",
					Labels: map[string]string{
						"cluster-operator.giantswarm.io/version": "3.7.0",
						"cluster.x-k8s.io/cluster-name":          "dh82p",
						"giantswarm.io/cluster":                  "dh82p",
						"giantswarm.io/organization":             "giantswarm",
						"release.giantswarm.io/version":          "14.2.2",
					},
					Annotations: map[string]string{
						"alpha.giantswarm.io/update-schedule-target-release": "15.2.1",
						"alpha.giantswarm.io/update-schedule-target-time":    "31 Dec 50 20:00 UTC",
					},
				},
			},
		},
		{
			name:                   "case 4",
			expectedReleaseVersion: "26.0.0",
			expectedEventTriggered: true,
			annotationsKept:        false,
			cluster: &capi.Cluster{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "test3",
					Namespace: "org-giantswarm",
					Labels: map[string]string{
						"cluster.x-k8s.io/cluster-name": "dh82p",
						"giantswarm.io/cluster":         "dh82p",
						"giantswarm.io/organization":    "giantswarm",
						"release.giantswarm.io/version": "25.0.0",
						"cluster.x-k8s.io/watch-filter": "capi",
					},
					Annotations: map[string]string{
						"alpha.giantswarm.io/update-schedule-target-release": "26.0.0",
						"alpha.giantswarm.io/update-schedule-target-time":    "31 Jul 24 14:00 UTC",
					},
				},
			},
			configMap: &corev1.ConfigMap{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "test3-userconfig",
					Namespace: "org-giantswarm",
				},
				Data: map[string]string{
					"values": `
					global:
					  connectivity:
						availabilityZoneUsageLimit: 3
						network: {}
						topology: {}
					  controlPlane: {}
					  metadata:
						description: Franco tests things from the bundle
						name: franco055
						organization: giantswarm
						preventDeletion: false
					  nodePools:
						nodepool0:
						  instanceType: m5.xlarge
						  maxSize: 10
						  minSize: 3
						  rootVolumeSizeGB: 8
					  providerSpecific: {}
					  release:
					  	version: 25.0.0
					`,
				},
			},
		},
	}

	for i, tc := range testCases {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			t.Log(tc.name)
			fakeClient := fake.NewClientBuilder().WithScheme(fakeScheme).WithObjects(tc.cluster).Build()
			fakeRecorder := record.NewFakeRecorder(1)
			r := &ClusterReconciler{
				Client:   fakeClient,
				Scheme:   fakeScheme,
				Log:      ctrl.Log.WithName("fake"),
				recorder: fakeRecorder,
			}
			ctx := context.TODO()

			if isCAPIProvider(tc.cluster) {
				err := fakeClient.Create(ctx, tc.configMap)
				if err != nil {
					t.Error(err)
				}
			}

			_, err := r.Reconcile(ctx, ctrl.Request{NamespacedName: types.NamespacedName{Name: tc.cluster.GetName(), Namespace: tc.cluster.GetNamespace()}})
			if err != nil {
				t.Error(err)
			}

			obj := &capi.Cluster{}
			err = fakeClient.Get(ctx, types.NamespacedName{Name: tc.cluster.GetName(), Namespace: tc.cluster.GetNamespace()}, obj)
			if err != nil {
				t.Error(err)
			}

			if isCAPIProvider(tc.cluster) {
				cm := &corev1.ConfigMap{}
				err = fakeClient.Get(ctx, types.NamespacedName{Name: tc.configMap.GetName(), Namespace: tc.configMap.GetNamespace()}, cm)
				if err != nil {
					t.Error(err)
				}
				if !strings.Contains(cm.Data["values"], fmt.Sprintf("version: %s", tc.expectedReleaseVersion)) {
					t.Fatalf("expected release to be %v, got %s", tc.expectedReleaseVersion, cm.Data["values"])
				}
			} else {
				if obj.Labels["release.giantswarm.io/version"] != tc.expectedReleaseVersion {
					t.Fatalf("expected release.giantswarm.io/version to be %v, got %s", tc.expectedReleaseVersion, obj.Labels["release.giantswarm.io/version"])
				}
			}

			if _, exists := obj.Annotations["alpha.giantswarm.io/update-schedule-target-release"]; exists != tc.annotationsKept {
				t.Fatalf("update schedule target release annotation expected to be %v, got %v", tc.annotationsKept, obj.Annotations["alpha.giantswarm.io/update-schedule-target-release"])
			}

			if _, exists := obj.Annotations["alpha.giantswarm.io/update-schedule-target-time"]; exists != tc.annotationsKept {
				t.Fatalf("update schedule target time annotation expected to be %v, got %v", tc.annotationsKept, obj.Annotations["alpha.giantswarm.io/update-schedule-target-release"])
			}

			triggered := false
			for eventsLeft := true; eventsLeft; {
				select {
				case event := <-fakeRecorder.Events:
					if strings.Contains(event, "ClusterUpgradeAnnouncement") {
						t.Log(event)
						triggered = true
					} else {
						t.Fatalf("test case %v failed. unexpected event %v", tc.name, event)
					}
				default:
					eventsLeft = false
				}
			}
			assert.Equal(t, tc.expectedEventTriggered, triggered, "test case %v failed.", tc.name)
		})
	}
}

func StringPtr(s string) *string {
	return &s
}
