package action

import (
	"context"
	"fmt"
	corev1 "k8s.io/api/core/v1"
	"reflect"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

// BatchHandlePod 批量处理Pod资源
type BatchHandlePod struct {
	client    client.Client
	originals map[string]*corev1.Pod
	news      map[string]*corev1.Pod
}

func NewBatchHandlePod(client client.Client) *BatchHandlePod {
	return &BatchHandlePod{
		client: client,
	}
}

func (o *BatchHandlePod) SetOriginal(key string, original *corev1.Pod) {
	if o.originals == nil {
		o.originals = map[string]*corev1.Pod{}
	}
	o.originals[key] = original
}

func (o *BatchHandlePod) SetNew(key string, new *corev1.Pod) {
	if o.news == nil {
		o.news = map[string]*corev1.Pod{}
	}
	o.news[key] = new
}

func (o *BatchHandlePod) Execute(ctx context.Context) (bool, error) {
	noChange := true
	for key, newPod := range o.news {
		originalPod := o.originals[key]
		if originalPod == nil {
			noChange = false
			if err := o.client.Create(ctx, newPod); err != nil {
				return noChange, fmt.Errorf("BatchHandlePod create error %q", err)
			}
			continue
		}
		needUpdate, needDelete := false, false
		if originalPod.Spec.Containers[0].Image != newPod.Spec.Containers[0].Image {
			originalPod.Spec.Containers[0].Image = newPod.Spec.Containers[0].Image
			needUpdate = true
		}
		if !reflect.DeepEqual(originalPod.Labels, newPod.Labels) {
			originalPod.Labels = newPod.Labels
			needUpdate = true
		}
		if !reflect.DeepEqual(originalPod.Spec.NodeSelector, newPod.Spec.NodeSelector) {
			needDelete = true
		}
		if !reflect.DeepEqual(originalPod.Spec.Containers[0].Resources, newPod.Spec.Containers[0].Resources) {
			needDelete = true
		}
		for index, volume := range newPod.Spec.Volumes {
			if !reflect.DeepEqual(volume, originalPod.Spec.Volumes[index]) {
				needDelete = true
				break
			}
		}
		for index, volume := range newPod.Spec.Containers[0].VolumeMounts {
			if !reflect.DeepEqual(volume, originalPod.Spec.Containers[0].VolumeMounts[index]) {
				needDelete = true
				break
			}
		}
		if needDelete {
			noChange = false
			if err := o.client.Delete(ctx, originalPod); err != nil {
				return noChange, fmt.Errorf("BatchHandlePod delete error %q", err)
			}
		} else if needUpdate {
			noChange = false
			if err := o.client.Update(ctx, originalPod); err != nil {
				return noChange, fmt.Errorf("BatchHandlePod update error %q", err)
			}
		}
	}
	return noChange, nil
}
