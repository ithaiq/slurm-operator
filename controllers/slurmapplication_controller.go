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

	"github.com/go-logr/logr"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"

	slurmoperatorv1beta1 "github.com/aistudio/slurm-operator/api/v1beta1"
)

// SlurmApplicationReconciler reconciles a SlurmApplication object
type SlurmApplicationReconciler struct {
	client.Client
	Log    logr.Logger
	Scheme *runtime.Scheme
}

// +kubebuilder:rbac:groups=slurmoperator.apulis.cn,resources=slurmapplications,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=slurmoperator.apulis.cn,resources=slurmapplications/status,verbs=get;update;patch

func (r *SlurmApplicationReconciler) Reconcile(req ctrl.Request) (ctrl.Result, error) {
	_ = context.Background()
	_ = r.Log.WithValues("slurmapplication", req.NamespacedName)

	// your logic here

	return ctrl.Result{}, nil
}

func (r *SlurmApplicationReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&slurmoperatorv1beta1.SlurmApplication{}).
		Complete(r)
}
