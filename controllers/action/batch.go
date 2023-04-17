package action

import (
	"context"
	corev1 "k8s.io/api/core/v1"

	"k8s.io/apimachinery/pkg/runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

// BatchHandleObject 批量处理对象资源
type BatchHandleObject struct {
	client    client.Client
	originals map[string]runtime.Object
	news      map[string]runtime.Object
}

func NewBatchHandleObject(client client.Client, originals map[string]runtime.Object, news map[string]runtime.Object) *BatchHandleObject {
	return &BatchHandleObject{client: client, originals: originals, news: news}
}

func (o *BatchHandleObject) Execute(ctx context.Context) (bool, error) {
	var isNoChange = true
	podHandle := NewBatchHandlePod(o.client)
	svcHandle := NewBatchHandleSvc(o.client)
	actions := []Action{podHandle, svcHandle}
	for key, newObj := range o.news {
		switch obj := newObj.(type) {
		case *corev1.Pod:
			podHandle.SetNew(key, obj)
			if o.originals[key] != nil {
				podHandle.SetOriginal(key, o.originals[key].(*corev1.Pod))
			}
		case *corev1.Service:
			svcHandle.SetNew(key, obj)
			if o.originals[key] != nil {
				svcHandle.SetOriginal(key, o.originals[key].(*corev1.Service))
			}
		}
	}
	for _, action := range actions {
		noChange, err := action.Execute(ctx)
		if err != nil {
			return isNoChange, err
		}
		if !noChange {
			isNoChange = false
		}
	}
	return isNoChange, nil
}
