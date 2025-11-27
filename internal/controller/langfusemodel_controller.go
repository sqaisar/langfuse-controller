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
	"strconv"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	logf "sigs.k8s.io/controller-runtime/pkg/log"

	langfusev1alpha1 "github.com/sqaisar/langfuse-controller/api/v1alpha1"
	"github.com/sqaisar/langfuse-controller/internal/langfuse"
)

// LangfuseModelReconciler reconciles a LangfuseModel object
type LangfuseModelReconciler struct {
	client.Client
	Scheme         *runtime.Scheme
	LangfuseClient *langfuse.Client
}

// +kubebuilder:rbac:groups=langfuse.io,resources=langfusemodels,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=langfuse.io,resources=langfusemodels/status,verbs=get;update;patch
// +kubebuilder:rbac:groups=langfuse.io,resources=langfusemodels/finalizers,verbs=update

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the LangfuseModel object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.22.4/pkg/reconcile
func (r *LangfuseModelReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	log := logf.FromContext(ctx)

	var model langfusev1alpha1.LangfuseModel
	if err := r.Get(ctx, req.NamespacedName, &model); err != nil {
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}

	if len(model.Status.Conditions) > 0 && model.Status.Conditions[0].Status == metav1.ConditionTrue {
		return ctrl.Result{}, nil
	}

	inputPrice, _ := strconv.ParseFloat(model.Spec.InputPrice, 64)
	outputPrice, _ := strconv.ParseFloat(model.Spec.OutputPrice, 64)
	totalPrice, _ := strconv.ParseFloat(model.Spec.TotalPrice, 64)

	lfModel := langfuse.Model{
		ModelName:       model.Spec.ModelName,
		MatchPattern:    model.Spec.MatchPattern,
		StartDate:       model.Spec.StartDate,
		Unit:            model.Spec.Unit,
		InputPrice:      inputPrice,
		OutputPrice:     outputPrice,
		TotalPrice:      totalPrice,
		TokenizerId:     model.Spec.TokenizerId,
		TokenizerConfig: model.Spec.TokenizerConfig,
	}

	log.Info("Creating Langfuse Model", "name", model.Spec.ModelName)
	_, err := r.LangfuseClient.CreateModel(lfModel)
	if err != nil {
		log.Error(err, "Failed to create Model")
		return ctrl.Result{}, err
	}

	// Update Status
	model.Status.Conditions = []metav1.Condition{
		{
			Type:               "Available",
			Status:             metav1.ConditionTrue,
			Reason:             "Created",
			Message:            "Model created successfully",
			LastTransitionTime: metav1.Now(),
		},
	}
	if err := r.Status().Update(ctx, &model); err != nil {
		return ctrl.Result{}, err
	}

	return ctrl.Result{}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *LangfuseModelReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&langfusev1alpha1.LangfuseModel{}).
		Named("langfusemodel").
		Complete(r)
}
