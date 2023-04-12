// controllers/action.go

package controllers

import (
	"context"
	"fmt"
	"reflect"

	"k8s.io/apimachinery/pkg/runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

type Action interface {
	Execute(context.Context) (ignore bool, err error)
}

// PatchStatus 更新对象status状态
type PatchStatus struct {
	client   client.Client
	original runtime.Object
	new      runtime.Object
}

func (o *PatchStatus) Execute(ctx context.Context) (ignore bool, err error) {
	if reflect.DeepEqual(o.original, o.new) {
		return true, nil
	}
	if err := o.client.Status().Patch(ctx, o.new, client.MergeFrom(o.original)); err != nil {
		return false, fmt.Errorf("PatchStatus error %q", err)
	}
	return false, nil
}

// PatchObject 更新对象资源
type PatchObject struct {
	client   client.Client
	original runtime.Object
	new      runtime.Object
}

func (o *PatchObject) Execute(ctx context.Context) (ignore bool, err error) {
	if reflect.DeepEqual(o.original, o.new) {
		return true, nil
	}
	if err := o.client.Patch(ctx, o.new, client.MergeFrom(o.original)); err != nil {
		return false, fmt.Errorf("PatchObject error %q", err)
	}
	return false, nil
}

// CreateObject 创建对象资源
type CreateObject struct {
	client client.Client
	obj    runtime.Object
}

func (o *CreateObject) Execute(ctx context.Context) (ignore bool, err error) {
	if err = o.client.Create(ctx, o.obj); err != nil {
		return false, fmt.Errorf("CreateObject error %q", err)
	}
	return false, nil
}

// BatchCreateObject 批量创建对象资源
type BatchCreateObject struct {
	client client.Client
	objs   []runtime.Object
}

func (o *BatchCreateObject) Execute(ctx context.Context) (ignore bool, err error) {
	for _, item := range o.objs {
		if err = o.client.Create(ctx, item); err != nil {
			return false, fmt.Errorf("BatchCreateObject error %q", err)
		}
	}
	return false, nil
}

// BatchPatchObject 批量更新对象资源
type BatchPatchObject struct {
	client    client.Client
	originals map[string]runtime.Object
	news      map[string]runtime.Object
}

func (o *BatchPatchObject) Execute(ctx context.Context) (ignore bool, err error) {
	ignore = true
	for key, original := range o.originals {
		if reflect.DeepEqual(original, o.news[key]) {
			continue
		}
		ignore = false
		if err = o.client.Patch(ctx, o.news[key], client.MergeFrom(original)); err != nil {
			return ignore, fmt.Errorf("BatchPatchObject error %q", err)
		}
	}
	return ignore, nil
}

// BatchCreateOrUpdateObject 批量创建或更新对象资源
type BatchCreateOrUpdateObject struct {
	client    client.Client
	originals map[string]runtime.Object
	news      map[string]runtime.Object
}

func (o *BatchCreateOrUpdateObject) Execute(ctx context.Context) (ignore bool, err error) {
	ignore = true
	for key, original := range o.originals {
		if reflect.DeepEqual(original, o.news[key]) {
			continue
		}
		ignore = false
		if o.news[key] == nil {
			return ignore, fmt.Errorf("BatchCreateOrUpdateObject news key %q is empty", key)
		}
		if reflect.ValueOf(original).IsNil() {
			fmt.Println("=============3", key)
			if err = o.client.Create(ctx, o.news[key]); err != nil {
				return ignore, fmt.Errorf("BatchCreateOrUpdateObject create error %q", err)
			}
		} else {
			fmt.Println("=============4", key)
			if err = o.client.Update(ctx, o.news[key]); err != nil {
				return ignore, fmt.Errorf("BatchCreateOrUpdateObject update error %q", err)
			}
		}
	}
	return ignore, nil
}
