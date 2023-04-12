/*
Copyright 2023 apulis.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package controllers

import (
	"context"
	"fmt"
	slurmoperatorv1beta1 "github.com/aistudio/slurm-operator/api/v1beta1"
	"github.com/go-logr/logr"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/tools/record"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

// SlurmApplicationReconciler reconciles a SlurmApplication object
type SlurmApplicationReconciler struct {
	client.Client
	Log      logr.Logger
	Scheme   *runtime.Scheme
	Recorder record.EventRecorder
}

// +kubebuilder:rbac:groups=slurmoperator.apulis.cn,resources=slurmapplications,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=slurmoperator.apulis.cn,resources=slurmapplications/status,verbs=get;update;patch
// +kubebuilder:rbac:groups=core,resources=services,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=core,resources=pods,verbs=get;list;watch;create;update;patch;delete

func (r *SlurmApplicationReconciler) Reconcile(req ctrl.Request) (ctrl.Result, error) {
	ctx := context.Background()
	log := r.Log.WithValues("Slurmapplication Reconcile", req.NamespacedName)

	resource, err := r.GetSlurmApplicationResource(ctx, req)
	if err != nil {
		return ctrl.Result{}, err
	}
	var action Action
	switch {
	case resource == nil || resource.app == nil: // 被删除了
		log.Info("slurm app not found. Ignoring.")
	case !resource.app.DeletionTimestamp.IsZero(): // 标记了删除
		log.Info("slurm app has been deleted. Ignoring.")
	case resource.app.Status == "":
		log.Info("slurm app staring. Updating status.")
		newBackup := resource.app.DeepCopy()
		newBackup.Status = slurmoperatorv1beta1.SlurmClusterStatusRunning
		action = &PatchStatus{client: r.Client, original: resource.app, new: newBackup}
	case resource.app.Status == slurmoperatorv1beta1.SlurmClusterStatusRunning ||
		resource.app.Status == slurmoperatorv1beta1.SlurmClusterStatusError:
		log.Info("slurm app start create or update resources...")
		originals := map[string]runtime.Object{}
		news := map[string]runtime.Object{}
		for key, item := range resource.actual.pods {
			originals[key] = item
		}
		for key, item := range resource.actual.svcs {
			originals[fmt.Sprintf("%ssvc", key)] = item
		}
		for key, item := range resource.desired.pods {
			news[key] = item
		}
		for key, item := range resource.desired.svcs {
			news[fmt.Sprintf("%ssvc", key)] = item
		}
		action = &BatchCreateOrUpdateObject{client: r.Client, news: news, originals: originals}
	}
	var ignore bool
	if action != nil {
		if ignore, err = action.Execute(ctx); err != nil {
			return ctrl.Result{}, fmt.Errorf("executing action error: %s", err)
		}
	}
	//未有变动说明监听到子资源的变化
	if ignore {
		isHealthy := r.IsSlurmApplicationHealthy(resource)
		newBackup := resource.app.DeepCopy()
		//如果子资源处于健康状态并且自身是不正常的则更新自身状态为正常状态
		if isHealthy && resource.app.Status == slurmoperatorv1beta1.SlurmClusterStatusError {
			newBackup.Status = slurmoperatorv1beta1.SlurmClusterStatusRunning
		}
		//如果子资源处于非健康状态并且自身是正常的则更新自身状态为错误状态
		if !isHealthy && resource.app.Status == slurmoperatorv1beta1.SlurmClusterStatusRunning {
			newBackup.Status = slurmoperatorv1beta1.SlurmClusterStatusError
			log.Info("slurm resources has error. please check event to fix")
		}
		action = &PatchStatus{client: r.Client, original: resource.app, new: newBackup}
		if _, err = action.Execute(ctx); err != nil {
			return ctrl.Result{}, fmt.Errorf("executing PatchStatus error: %s", err)
		}
	}
	return ctrl.Result{}, nil
}

func (r *SlurmApplicationReconciler) GetSlurmApplicationResource(ctx context.Context, req ctrl.Request) (*SlurmApplicationResource, error) {
	var (
		saResource *SlurmApplicationResource
		slurmApp   slurmoperatorv1beta1.SlurmApplication
	)
	if err := r.Get(ctx, req.NamespacedName, &slurmApp); err != nil {
		return saResource, client.IgnoreNotFound(err)
	}
	saResource = NewSlurmApplicationResource(&slurmApp)
	if err := saResource.SetStateActual(ctx, r.Scheme, r.Client); err != nil {
		return nil, fmt.Errorf("setting actual state error: %s", err)
	}
	if err := saResource.SetStateDesired(ctx, r.Scheme, r.Client); err != nil {
		return nil, fmt.Errorf("setting desired state error: %s", err)
	}
	return saResource, nil
}

func (r *SlurmApplicationReconciler) IsSlurmApplicationHealthy(resource *SlurmApplicationResource) bool {
	for _, pod := range resource.actual.pods {
		if pod.Status.Phase != corev1.PodRunning {
			return false
		}
	}
	return true
}

func (r *SlurmApplicationReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&slurmoperatorv1beta1.SlurmApplication{}).
		Owns(&corev1.Pod{}).
		Owns(&corev1.Service{}).
		Complete(r)
}
