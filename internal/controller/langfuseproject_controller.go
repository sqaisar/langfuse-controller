/*
Copyright 2025.

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

package controller

import (
	"context"

	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	logf "sigs.k8s.io/controller-runtime/pkg/log"

	langfusev1alpha1 "github.com/sqaisar/langfuse-controller/api/v1alpha1"
	"github.com/sqaisar/langfuse-controller/internal/langfuse"
)

// LangfuseProjectReconciler reconciles a LangfuseProject object
type LangfuseProjectReconciler struct {
	client.Client
	Scheme         *runtime.Scheme
	LangfuseClient *langfuse.Client
}

// +kubebuilder:rbac:groups=langfuse.io,resources=langfuseprojects,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=langfuse.io,resources=langfuseprojects/status,verbs=get;update;patch
// +kubebuilder:rbac:groups=langfuse.io,resources=langfuseprojects/finalizers,verbs=update

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the LangfuseProject object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.22.4/pkg/reconcile
func (r *LangfuseProjectReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	log := logf.FromContext(ctx)

	var project langfusev1alpha1.LangfuseProject
	if err := r.Get(ctx, req.NamespacedName, &project); err != nil {
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}

	if project.Status.ID != "" {
		// Project already exists, maybe check if it still exists?
		// For now, assume it's synced.
		return ctrl.Result{}, nil
	}

	log.Info("Creating Langfuse Project", "name", project.Spec.Name)
	lfProject, err := r.LangfuseClient.CreateProject(project.Spec.Name)
	if err != nil {
		log.Error(err, "Failed to create Langfuse Project")
		project.Status.State = "Error"
		if err := r.Status().Update(ctx, &project); err != nil {
			return ctrl.Result{}, err
		}
		return ctrl.Result{}, err
	}

	project.Status.ID = lfProject.ID
	project.Status.State = "Ready"
	if err := r.Status().Update(ctx, &project); err != nil {
		log.Error(err, "Failed to update LangfuseProject status")
		return ctrl.Result{}, err
	}

	log.Info("Langfuse Project created successfully", "id", lfProject.ID)
	return ctrl.Result{}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *LangfuseProjectReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&langfusev1alpha1.LangfuseProject{}).
		Named("langfuseproject").
		Complete(r)
}
