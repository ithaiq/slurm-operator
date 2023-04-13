package controllers

import (
	"context"
	slurmoperatorv1beta1 "github.com/aistudio/slurm-operator/api/v1beta1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"strconv"
)

type SlurmApplicationResource struct {
	app *slurmoperatorv1beta1.SlurmApplication
	// 真实的资源
	actual *SlurmApplicationContainer
	// 期望的资源
	desired *SlurmApplicationContainer
}

type SlurmApplicationContainer struct {
	pods map[string]*corev1.Pod
	svcs map[string]*corev1.Service
}

func NewSlurmApplicationResource(app *slurmoperatorv1beta1.SlurmApplication) *SlurmApplicationResource {
	return &SlurmApplicationResource{
		app: app,
		actual: &SlurmApplicationContainer{
			pods: make(map[string]*corev1.Pod),
			svcs: make(map[string]*corev1.Service),
		},
		desired: &SlurmApplicationContainer{
			pods: make(map[string]*corev1.Pod),
			svcs: make(map[string]*corev1.Service),
		},
	}
}

func (r *SlurmApplicationResource) SetActualSlurmResource(ctx context.Context, scheme *runtime.Scheme, client client.Client) (err error) {
	ns := r.app.Namespace
	saResource := NewSlurmApplicationHandler(scheme, client)
	if r.actual.pods[PodSlurmJupyter], r.actual.svcs[SvcSlurmJupyter], err = saResource.GetSlurmApplicationResource(ctx, SlurmJupyter, ns); err != nil {
		return err
	}
	if r.actual.pods[PodSlurmMaster], r.actual.svcs[SvcSlurmMaster], err = saResource.GetSlurmApplicationResource(ctx, SlurmMaster, ns); err != nil {
		return err
	}
	for i := 1; i <= r.app.Spec.Node.Instance; i++ {
		name := SlurmNode + strconv.Itoa(i)
		podName := PodSlurmNode + strconv.Itoa(i)
		svcName := SvcSlurmNode + strconv.Itoa(i)
		if r.actual.pods[podName], r.actual.svcs[svcName], err = saResource.GetSlurmApplicationResource(ctx, name, ns); err != nil {
			return err
		}
	}
	return nil
}

func (r *SlurmApplicationResource) SetDesiredSlurmResource(ctx context.Context, scheme *runtime.Scheme, client client.Client) (err error) {
	saResource := NewSlurmApplicationHandler(scheme, client)
	if r.desired.pods[PodSlurmJupyter], r.desired.svcs[SvcSlurmJupyter], err = saResource.GenerateJupyterResource(
		r.app,
		r.actual.pods[PodSlurmJupyter],
		r.actual.svcs[SvcSlurmJupyter],
		SlurmJupyter,
	); err != nil {
		return err
	}
	if r.desired.pods[PodSlurmMaster], r.desired.svcs[SvcSlurmMaster], err = saResource.GenerateSlurmResource(
		r.app,
		&r.app.Spec.Master.CommonSpec,
		r.actual.pods[PodSlurmMaster],
		r.actual.svcs[SvcSlurmMaster],
		SlurmMaster,
	); err != nil {
		return err
	}
	for i := 1; i <= r.app.Spec.Node.Instance; i++ {
		name := SlurmNode + strconv.Itoa(i)
		podName := PodSlurmNode + strconv.Itoa(i)
		svcName := SvcSlurmNode + strconv.Itoa(i)
		if r.desired.pods[podName], r.desired.svcs[svcName], err = saResource.GenerateSlurmResource(
			r.app,
			&r.app.Spec.Node.CommonSpec,
			r.actual.pods[podName],
			r.actual.svcs[svcName],
			name,
		); err != nil {
			return err
		}
	}
	return nil
}
