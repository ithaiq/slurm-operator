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
		r.Recorder.Event(resource.app, corev1.EventTypeWarning, "GetSlurmApplicationResource-Failed", err.Error())
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
		newBackup.Status = slurmoperatorv1beta1.SlurmClusterStatusPending
		action = &PatchStatus{client: r.Client, original: resource.app, new: newBackup}
	case resource.app.Status == slurmoperatorv1beta1.SlurmClusterStatusPending ||
		resource.app.Status == slurmoperatorv1beta1.SlurmClusterStatusRunning ||
		resource.app.Status == slurmoperatorv1beta1.SlurmClusterStatusError:
		log.Info("slurm app start create or update resources...")
		originals, news := r.convertToObject(resource)
		action = &BatchCreateOrUpdateObject{client: r.Client, news: news, originals: originals}
	}
	var noChange bool
	if action != nil {
		if noChange, err = action.Execute(ctx); err != nil {
			r.Recorder.Event(resource.app, corev1.EventTypeWarning, "ActionExecute-Failed", err.Error())
			return ctrl.Result{}, fmt.Errorf("executing action error: %s", err)
		}
	}
	//未有变动说明监听到子资源的变化
	if noChange {
		isHealthy := r.isSlurmApplicationHealthy(resource)
		newBackup := resource.app.DeepCopy()
		//如果子资源处于健康状态并且自身不是Running状态则更新Running状态
		if isHealthy && resource.app.Status != slurmoperatorv1beta1.SlurmClusterStatusRunning {
			newBackup.Status = slurmoperatorv1beta1.SlurmClusterStatusRunning
		}
		//如果子资源处于非健康状态并且自身是Running状态则更新Error状态（兼容pending状态）
		if !isHealthy && resource.app.Status == slurmoperatorv1beta1.SlurmClusterStatusRunning {
			newBackup.Status = slurmoperatorv1beta1.SlurmClusterStatusError
			log.Info("slurm resources has error. please check event to fix")
		}
		action = &PatchStatus{client: r.Client, original: resource.app, new: newBackup}
		if _, err = action.Execute(ctx); err != nil {
			r.Recorder.Event(resource.app, corev1.EventTypeWarning, "PatchStatus-Failed", err.Error())
			return ctrl.Result{}, fmt.Errorf("executing PatchStatus error: %s", err)
		}
	}
	return ctrl.Result{}, nil
}

func (r *SlurmApplicationReconciler) convertToObject(resource *SlurmApplicationResource) (map[string]runtime.Object, map[string]runtime.Object) {
	originals := map[string]runtime.Object{}
	news := map[string]runtime.Object{}
	for key, item := range resource.actual.pods {
		originals[key] = item
	}
	for key, item := range resource.actual.svcs {
		originals[key] = item
	}
	for key, item := range resource.desired.pods {
		news[key] = item
	}
	for key, item := range resource.desired.svcs {
		news[key] = item
	}
	return originals, news
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
	if err := saResource.SetActualSlurmResource(ctx, r.Scheme, r.Client); err != nil {
		return nil, fmt.Errorf("setting actual slurm resource error: %s", err)
	}
	if err := saResource.SetDesiredSlurmResource(ctx, r.Scheme, r.Client); err != nil {
		return nil, fmt.Errorf("setting desired slurm resource error: %s", err)
	}
	return saResource, nil
}

func (r *SlurmApplicationReconciler) isSlurmApplicationHealthy(resource *SlurmApplicationResource) bool {
	for _, pod := range resource.actual.pods {
		if pod.Status.Phase != corev1.PodRunning {
			return false
		}
		ready := false
		for _, condition := range pod.Status.Conditions {
			if condition.Type == corev1.PodReady && condition.Status == corev1.ConditionTrue {
				ready = true
				break
			}
		}
		if !ready {
			return false
		}
	}
	//TODO 是否需要检查service是否健康
	return true
}

func (r *SlurmApplicationReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&slurmoperatorv1beta1.SlurmApplication{}).
		Owns(&corev1.Pod{}).
		Owns(&corev1.Service{}).
		Complete(r)
}
