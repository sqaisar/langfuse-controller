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
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
)

// LangfuseScoreConfigReconciler reconciles a LangfuseScoreConfig object
type LangfuseScoreConfigReconciler struct {
	client.Client
	Scheme         *runtime.Scheme
	LangfuseClient *langfuse.Client
}

// +kubebuilder:rbac:groups=langfuse.io,resources=langfusescoreconfigs,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=langfuse.io,resources=langfusescoreconfigs/status,verbs=get;update;patch
// +kubebuilder:rbac:groups=langfuse.io,resources=langfusescoreconfigs/finalizers,verbs=update

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the LangfuseScoreConfig object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.22.4/pkg/reconcile
func (r *LangfuseScoreConfigReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	log := logf.FromContext(ctx)

	var config langfusev1alpha1.LangfuseScoreConfig
	if err := r.Get(ctx, req.NamespacedName, &config); err != nil {
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}

	if len(config.Status.Conditions) > 0 && config.Status.Conditions[0].Status == metav1.ConditionTrue {
		return ctrl.Result{}, nil
	}

	var project langfusev1alpha1.LangfuseProject
	if err := r.Get(ctx, types.NamespacedName{Name: config.Spec.ProjectRef, Namespace: req.Namespace}, &project); err != nil {
		log.Error(err, "Failed to get Project")
		return ctrl.Result{}, err
	}

	if project.Status.ID == "" {
		return ctrl.Result{Requeue: true}, nil
	}

	log.Info("Creating Score Config", "name", config.Spec.Name)
	if err := r.LangfuseClient.CreateScoreConfig(project.Status.ID, map[string]interface{}{
		"name":       config.Spec.Name,
		"dataType":   config.Spec.DataType,
		"minValue":   config.Spec.MinValue,
		"maxValue":   config.Spec.MaxValue,
		"categories": config.Spec.Categories,
	}); err != nil {
		log.Error(err, "Failed to create Score Config")
		return ctrl.Result{}, err
	}

	config.Status.Conditions = []metav1.Condition{{
		Type:               "Available",
		Status:             metav1.ConditionTrue,
		Reason:             "Created",
		Message:            "Score Config created",
		LastTransitionTime: metav1.Now(),
	}}
	if err := r.Status().Update(ctx, &config); err != nil {
		return ctrl.Result{}, err
	}

	return ctrl.Result{}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *LangfuseScoreConfigReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&langfusev1alpha1.LangfuseScoreConfig{}).
		Named("langfusescoreconfig").
		Complete(r)
}
