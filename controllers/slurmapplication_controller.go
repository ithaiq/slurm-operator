/*
Copyright 2023 xxx.

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
	. "github.com/aistudio/slurm-operator/controllers/action"
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

// +kubebuilder:rbac:groups=slurmoperator.xxx.cn,resources=slurmapplications,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=slurmoperator.xxx.cn,resources=slurmapplications/status,verbs=get;update;patch
// +kubebuilder:rbac:groups=core,resources=services,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=core,resources=pods,verbs=get;list;watch;create;update;patch;delete

func (r *SlurmApplicationReconciler) Reconcile(req ctrl.Request) (ctrl.Result, error) {

	ctx := context.Background()
	log := r.Log.WithValues("Slurmapplication Reconcile", req.NamespacedName)

	manager, err := r.GetSlurmApplicationManager(ctx, req)
	if err != nil {
		r.Recorder.Event(manager.app, corev1.EventTypeWarning, "GetSlurmApplicationManager-Failed", err.Error())
		return ctrl.Result{}, err
	}
	if manager == nil {
		return ctrl.Result{}, nil
	}
	app := manager.GetApp()
	var action Action

	switch {
	case app == nil: // 被删除了
		log.Info("slurm app not found. Ignoring.")
	case !app.DeletionTimestamp.IsZero(): // 标记了删除
		log.Info("slurm app has been deleted. Ignoring.")
	case app.Status == "":
		log.Info("slurm app staring. Updating status.")
		newBackup := app.DeepCopy()
		newBackup.Status = slurmoperatorv1beta1.SlurmClusterStatusPending
		action = NewPatchStatus(r.Client, app, newBackup)
	case app.Status == slurmoperatorv1beta1.SlurmClusterStatusPending ||
		app.Status == slurmoperatorv1beta1.SlurmClusterStatusRunning ||
		app.Status == slurmoperatorv1beta1.SlurmClusterStatusError:
		log.Info("slurm app start create or update resources...")
		originals, news := r.convertToObject(manager.GetActual(), manager.GetDesired())
		action = NewBatchHandleObject(r.Client, originals, news)
	}
	var noChange bool
	if action != nil {
		if noChange, err = action.Execute(ctx); err != nil {
			r.Recorder.Event(app, corev1.EventTypeWarning, "ActionExecute-Failed", err.Error())
			return ctrl.Result{}, fmt.Errorf("executing action error: %s", err)
		}
	}
	//未有变动说明监听到子资源的变化
	if noChange {
		isHealthy := r.isSlurmApplicationHealthy(manager.GetActual())
		newBackup := app.DeepCopy()
		//如果子资源处于健康状态并且自身不是Running状态则更新Running状态
		if isHealthy && app.Status != slurmoperatorv1beta1.SlurmClusterStatusRunning {
			newBackup.Status = slurmoperatorv1beta1.SlurmClusterStatusRunning
		}
		//如果子资源处于非健康状态并且自身是Running状态则更新Error状态（兼容pending状态）
		if !isHealthy && app.Status == slurmoperatorv1beta1.SlurmClusterStatusRunning {
			newBackup.Status = slurmoperatorv1beta1.SlurmClusterStatusError
			log.Info("slurm resources has error. please check event to fix")
		}
		action = NewPatchStatus(r.Client, app, newBackup)
		if _, err = action.Execute(ctx); err != nil {
			r.Recorder.Event(app, corev1.EventTypeWarning, "PatchStatus-Failed", err.Error())
			return ctrl.Result{}, fmt.Errorf("executing PatchStatus error: %s", err)
		}
	}
	return ctrl.Result{}, nil
}

func (r *SlurmApplicationReconciler) convertToObject(actual *SlurmApplicationResources, desired *SlurmApplicationResources) (map[string]runtime.Object, map[string]runtime.Object) {
	originals := map[string]runtime.Object{}
	news := map[string]runtime.Object{}
	for key, item := range actual.GetPods() {
		originals[CombineSlurmResourceType(key, SlurmPod)] = item
	}
	for key, item := range actual.GetSvcs() {
		originals[CombineSlurmResourceType(key, SlurmService)] = item
	}
	for key, item := range desired.GetPods() {
		news[CombineSlurmResourceType(key, SlurmPod)] = item
	}
	for key, item := range desired.GetSvcs() {
		news[CombineSlurmResourceType(key, SlurmService)] = item
	}
	return originals, news
}

func (r *SlurmApplicationReconciler) GetSlurmApplicationManager(ctx context.Context, req ctrl.Request) (*SlurmApplicationManager, error) {
	var (
		manager  *SlurmApplicationManager
		slurmApp slurmoperatorv1beta1.SlurmApplication
	)
	if err := r.Get(ctx, req.NamespacedName, &slurmApp); err != nil {
		return manager, client.IgnoreNotFound(err)
	}
	manager = NewSlurmApplicationManager(&slurmApp, r.Scheme, r.Client)
	if err := manager.GetActualSlurmResource(ctx); err != nil {
		return nil, fmt.Errorf("get actual slurm resource error: %s", err)
	}
	if err := manager.GetDesiredSlurmResource(ctx); err != nil {
		return nil, fmt.Errorf("get desired slurm resource error: %s", err)
	}
	return manager, nil
}

func (r *SlurmApplicationReconciler) isSlurmApplicationHealthy(resource *SlurmApplicationResources) bool {
	for _, pod := range resource.GetPods() {
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
