package action

import (
	"context"
	"fmt"
	"k8s.io/apimachinery/pkg/runtime"
	"reflect"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

// PatchStatus 更新对象status状态
type PatchStatus struct {
	client   client.Client
	original runtime.Object
	new      runtime.Object
}

func NewPatchStatus(client client.Client, original runtime.Object, new runtime.Object) *PatchStatus {
	return &PatchStatus{client: client, original: original, new: new}
}

func (o *PatchStatus) Execute(ctx context.Context) (noChange bool, err error) {
	if reflect.DeepEqual(o.original, o.new) {
		return true, nil
	}
	if err = o.client.Status().Patch(ctx, o.new, client.MergeFrom(o.original)); err != nil {
		return true, fmt.Errorf("PatchStatus error %q", err)
	}
	return false, nil
}
