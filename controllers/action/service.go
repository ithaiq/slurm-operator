package action

import (
	"context"
	"fmt"
	corev1 "k8s.io/api/core/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

// BatchHandleSvc 批量处理Service资源
type BatchHandleSvc struct {
	client    client.Client
	originals map[string]*corev1.Service
	news      map[string]*corev1.Service
}

func NewBatchHandleSvc(client client.Client) *BatchHandleSvc {
	return &BatchHandleSvc{
		client: client,
	}
}

func (o *BatchHandleSvc) SetOriginal(key string, original *corev1.Service) {
	if o.originals == nil {
		o.originals = map[string]*corev1.Service{}
	}
	o.originals[key] = original
}

func (o *BatchHandleSvc) SetNew(key string, new *corev1.Service) {
	if o.news == nil {
		o.news = map[string]*corev1.Service{}
	}
	o.news[key] = new
}

func (o *BatchHandleSvc) Execute(ctx context.Context) (bool, error) {
	noChange := true
	for key, newPod := range o.news {
		originalPod := o.originals[key]
		if originalPod != nil { //暂未考虑service更新的问题
			continue
		}
		noChange = false
		if err := o.client.Create(ctx, newPod); err != nil {
			return noChange, fmt.Errorf("BatchHandleSvc create error %q", err)
		}
	}
	return noChange, nil
}
