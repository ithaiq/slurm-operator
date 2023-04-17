package controllers

import (
	"context"
	slurmoperatorv1beta1 "github.com/aistudio/slurm-operator/api/v1beta1"
	"k8s.io/apimachinery/pkg/runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"strconv"
)

type SlurmApplicationManager struct {
	app *slurmoperatorv1beta1.SlurmApplication
	// 真实的资源
	actual *SlurmApplicationResources
	// 期望的资源
	desired *SlurmApplicationResources
	// 处理器
	handler *SlurmApplicationHandler
}

func NewSlurmApplicationManager(app *slurmoperatorv1beta1.SlurmApplication, scheme *runtime.Scheme, client client.Client) *SlurmApplicationManager {
	return &SlurmApplicationManager{
		app:     app,
		actual:  NewSlurmApplicationResources(app),
		desired: NewSlurmApplicationResources(app),
		handler: NewSlurmApplicationHandler(scheme, client),
	}
}

func (r *SlurmApplicationManager) GetApp() *slurmoperatorv1beta1.SlurmApplication {
	return r.app
}

func (r *SlurmApplicationManager) GetActual() *SlurmApplicationResources {
	return r.actual
}

func (r *SlurmApplicationManager) GetDesired() *SlurmApplicationResources {
	return r.desired
}

func (r *SlurmApplicationManager) GetActualSlurmResource(ctx context.Context) (err error) {
	var names []string
	names = append(names, GetSlurmResourceName(r.app.Name, SlurmJupyter))
	names = append(names, GetSlurmResourceName(r.app.Name, SlurmMaster))
	for i := 1; i <= r.app.Spec.Node.Instance; i++ {
		names = append(names, GetSlurmResourceName(r.app.Name, SlurmNode+strconv.Itoa(i)))
	}
	return r.getActualSlurmResources(ctx, names...)
}

func (r *SlurmApplicationManager) GetDesiredSlurmResource(ctx context.Context) (err error) {
	if err = r.generateDesiredJupyterResource(r.app, GetSlurmResourceName(r.app.Name, SlurmJupyter)); err != nil {
		return err
	}
	if err = r.generateDesiredSubSlurmResource(r.app, &r.app.Spec.Master.CommonSpec, GetSlurmResourceName(r.app.Name, SlurmMaster)); err != nil {
		return err
	}
	for i := 1; i <= r.app.Spec.Node.Instance; i++ {
		if err = r.generateDesiredSubSlurmResource(r.app, &r.app.Spec.Node.CommonSpec, GetSlurmResourceName(r.app.Name, SlurmNode+strconv.Itoa(i))); err != nil {
			return err
		}
	}
	return nil
}

func (r *SlurmApplicationManager) getActualSlurmResources(ctx context.Context, names ...string) error {
	for _, name := range names {
		pod, err := r.handler.GetPodByNamespacedName(ctx, name, r.app.Namespace)
		if err != nil {
			return err
		}
		svc, err := r.handler.GetSvcByNamespacedName(ctx, name, r.app.Namespace)
		if err != nil {
			return err
		}
		r.actual.SetPods(pod)
		r.actual.SetSvcs(svc)
	}
	return nil
}

func (r *SlurmApplicationManager) generateDesiredJupyterResource(app *slurmoperatorv1beta1.SlurmApplication, name string) error {
	pod, err := r.handler.GenerateJupyterPod(app, name)
	if err != nil {
		return err
	}
	svc, err := r.handler.GenerateJupyterSvc(app, name)
	if err != nil {
		return err
	}
	r.desired.SetPods(pod)
	r.desired.SetSvcs(svc)
	return nil
}

func (r *SlurmApplicationManager) generateDesiredSubSlurmResource(app *slurmoperatorv1beta1.SlurmApplication, cs *slurmoperatorv1beta1.CommonSpec, name string) error {
	pod, err := r.handler.GenerateSubSlurmPod(app, cs, name)
	if err != nil {
		return err
	}
	svc, err := r.handler.GenerateSubSlurmSvc(app, name)
	if err != nil {
		return err
	}
	r.desired.SetPods(pod)
	r.desired.SetSvcs(svc)
	return nil
}
