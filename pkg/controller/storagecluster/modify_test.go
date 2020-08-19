package storagecluster

import (
	"context"
	"testing"

	api "github.com/openshift/ocs-operator/pkg/apis/ocs/v1"
	v1 "github.com/openshift/ocs-operator/pkg/apis/ocs/v1"
	"github.com/stretchr/testify/assert"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
)

func TestUpdate(t *testing.T) {
	request := reconcile.Request{
		NamespacedName: types.NamespacedName{
			Name:      "storage-test",
			Namespace: "storage-test-ns",
		},
	}
	cases := []struct {
		label     string
		sc        *v1.StorageCluster
		condition string
	}{
		{
			label:     "case 1", // Update event
			condition: "update",
		},
		{
			label:     "case 2", // No Update Event
			condition: "noUpdate",
		},
		{
			label:     "case 3", // No Event At All
			condition: "",
		},
	}
	for _, c := range cases {
		c.sc = &v1.StorageCluster{}
		mockStorageCluster.DeepCopyInto(c.sc)
		fakeReconciler := createFakeStorageClusterReconciler(t, c.sc)
		orig := &api.StorageCluster{}
		fakeReconciler.client.Get(context.TODO(), request.NamespacedName, orig)
		if c.condition == "update" {
			request.NamespacedName.Name = "storage-test-modify"
			fakeReconciler.Reconcile(request)
			fakeReconciler.StatusUpdate(orig)
			assert.Equal(t, c.sc.ObjectMeta.Name, orig.ObjectMeta.Name)
		} else if c.condition == "noUpdate" {
			fakeReconciler.Reconcile(request)
			fakeReconciler.StatusUpdate(orig)
			assert.Equal(t, c.sc.ObjectMeta.Name, orig.ObjectMeta.Name)
			// TODO: This is Failing saying actual value is " "
			// and expecting value is "storage-test"
		} else {
			// TODO: What param am I supposed to pass to StatusUpdate()
			// when there is "no Event"
		}

	}

}
