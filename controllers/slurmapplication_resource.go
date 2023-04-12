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

func (r *SlurmApplicationResource) SetStateActual(ctx context.Context, scheme *runtime.Scheme, client client.Client) error {
	var err error
	ns := r.app.Namespace
	saResource := NewSlurmApplicationHandler(scheme, client)
	r.actual.pods[SlurmJupyter], r.actual.svcs[SlurmJupyter], err = saResource.GetSlurmApplicationResource(ctx, SlurmJupyter, ns)
	if err != nil {
		return err
	}
	r.actual.pods[SlurmMaster], r.actual.svcs[SlurmMaster], err = saResource.GetSlurmApplicationResource(ctx, SlurmMaster, ns)
	if err != nil {
		return err
	}
	for i := 1; i <= r.app.Spec.Node.Instance; i++ {
		name := SlurmNode + strconv.Itoa(i)
		r.actual.pods[name], r.actual.svcs[name], err = saResource.GetSlurmApplicationResource(ctx, name, ns)
		if err != nil {
			return err
		}
	}
	return nil
}

func (r *SlurmApplicationResource) SetStateDesired(ctx context.Context, scheme *runtime.Scheme, client client.Client) error {
	var err error
	saResource := NewSlurmApplicationHandler(scheme, client)
	r.desired.pods[SlurmJupyter], r.desired.svcs[SlurmJupyter], err = saResource.GenerateJupyterResource(
		r.app,
		r.actual.pods[SlurmJupyter],
		r.actual.svcs[SlurmJupyter],
		SlurmJupyter,
	)
	if err != nil {
		return err
	}
	r.desired.pods[SlurmMaster], r.desired.svcs[SlurmMaster], err = saResource.GenerateSlurmResource(
		r.app,
		&r.app.Spec.Master.CommonSpec,
		r.actual.pods[SlurmMaster],
		r.actual.svcs[SlurmMaster],
		SlurmMaster,
	)
	if err != nil {
		return err
	}
	for i := 1; i <= r.app.Spec.Node.Instance; i++ {
		name := SlurmNode + strconv.Itoa(i)
		r.desired.pods[name], r.desired.svcs[name], err = saResource.GenerateSlurmResource(r.app, &r.app.Spec.Node.CommonSpec, r.actual.pods[name], r.actual.svcs[name], name)
		if err != nil {
			return err
		}
	}
	return nil
}
