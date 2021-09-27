package controllers

import (
	"context"
	"strconv"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	utilruntime "k8s.io/apimachinery/pkg/util/runtime"
	clientgoscheme "k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/tools/record"

	"sigs.k8s.io/cluster-api/api/v1alpha3"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client/fake"
)

var (
	fakeScheme = runtime.NewScheme()
)

func init() {
	utilruntime.Must(clientgoscheme.AddToScheme(fakeScheme))
	_ = v1alpha3.AddToScheme(fakeScheme)
}

func TestClusterController(t *testing.T) {
	testCases := []struct {
		name                        string
		expectedReleaseVersion      string
		expectedEventTriggeredCount int
		annotationsKept             bool

		cluster *v1alpha3.Cluster
	}{
		// event triggered, within office time
		{
			name:                        "case 0",
			expectedReleaseVersion:      "15.2.1",
			annotationsKept:             false,
			expectedEventTriggeredCount: 2,
			cluster: &v1alpha3.Cluster{
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
			name:                        "case 1",
			expectedReleaseVersion:      "15.2.1",
			expectedEventTriggeredCount: 2,
			annotationsKept:             false,
			cluster: &v1alpha3.Cluster{
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
			name:                        "case 2",
			expectedReleaseVersion:      "15.2.1",
			expectedEventTriggeredCount: 2,
			annotationsKept:             false,
			cluster: &v1alpha3.Cluster{
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
			name:                        "case 3",
			expectedReleaseVersion:      "14.2.2",
			expectedEventTriggeredCount: 0,
			annotationsKept:             true,
			cluster: &v1alpha3.Cluster{
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
		// invalid "upgrade"
		{
			name:                        "case 3",
			expectedReleaseVersion:      "14.2.2",
			expectedEventTriggeredCount: 1,
			annotationsKept:             true,
			cluster: &v1alpha3.Cluster{
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
						"alpha.giantswarm.io/update-schedule-target-release": "14.2.0",
						"alpha.giantswarm.io/update-schedule-target-time":    "13 Sep 21 19:00 UTC",
					},
				},
			},
		},
	}

	for i, tc := range testCases {
		t.Run(strconv.Itoa(i), func(t *testing.T) {

			fakeClient := fake.NewClientBuilder().WithScheme(fakeScheme).WithObjects(tc.cluster).Build()
			fakeRecorder := record.NewFakeRecorder(2)
			r := &ClusterReconciler{
				Client:   fakeClient,
				Scheme:   fakeScheme,
				Log:      ctrl.Log.WithName("fake"),
				recorder: fakeRecorder,
			}
			ctx := context.TODO()
			_, err := r.Reconcile(ctx, ctrl.Request{NamespacedName: types.NamespacedName{Name: tc.cluster.GetName(), Namespace: tc.cluster.GetNamespace()}})
			if err != nil {
				t.Error(err)
			}

			obj := &v1alpha3.Cluster{}
			err = fakeClient.Get(ctx, types.NamespacedName{Name: tc.cluster.GetName(), Namespace: tc.cluster.GetNamespace()}, obj)
			if err != nil {
				t.Error(err)
			}

			if obj.Labels["release.giantswarm.io/version"] != tc.expectedReleaseVersion {
				t.Fatalf("expected release.giantswarm.io/version to be %v, got %s", tc.expectedReleaseVersion, obj.Labels["release.giantswarm.io/version"])
			}

			if _, exists := obj.Annotations["alpha.giantswarm.io/update-schedule-target-release"]; exists != tc.annotationsKept {
				t.Fatalf("update schedule target release annotation expected to be %v, got %v", tc.annotationsKept, obj.Annotations["alpha.giantswarm.io/update-schedule-target-release"])
			}

			if _, exists := obj.Annotations["alpha.giantswarm.io/update-schedule-target-time"]; exists != tc.annotationsKept {
				t.Fatalf("update schedule target time annotation expected to be %v, got %v", tc.annotationsKept, obj.Annotations["alpha.giantswarm.io/update-schedule-target-release"])
			}

			triggeredCount := 0
			for eventsLeft := true; eventsLeft; {
				select {
				case event := <-fakeRecorder.Events:
					if strings.Contains(event, "ClusterUpgradeAnnouncement") {
						t.Log(event)
						triggeredCount++
					} else {
						t.Fatalf("test case %v failed. unexpected event %v", tc.name, event)
					}
				default:
					eventsLeft = false
				}
			}
			assert.Equal(t, tc.expectedEventTriggeredCount, triggeredCount, "test case %v failed.", tc.name)
		})
	}
}

func StringPtr(s string) *string {
	return &s
}
