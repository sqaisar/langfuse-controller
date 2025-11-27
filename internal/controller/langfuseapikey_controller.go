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
	"sigs.k8s.io/controller-runtime/pkg/controller"
	logf "sigs.k8s.io/controller-runtime/pkg/log"

	langfusev1alpha1 "github.com/sqaisar/langfuse-controller/api/v1alpha1"
	"github.com/sqaisar/langfuse-controller/internal/langfuse"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
)

// LangfuseAPIKeyReconciler reconciles a LangfuseAPIKey object
type LangfuseAPIKeyReconciler struct {
	client.Client
	Scheme         *runtime.Scheme
	LangfuseClient *langfuse.Client
}

// +kubebuilder:rbac:groups=langfuse.io,resources=langfuseapikeys,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=langfuse.io,resources=langfuseapikeys/status,verbs=get;update;patch
// +kubebuilder:rbac:groups=langfuse.io,resources=langfuseapikeys/finalizers,verbs=update

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the LangfuseAPIKey object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.22.4/pkg/reconcile
func (r *LangfuseAPIKeyReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	log := logf.FromContext(ctx)

	var apiKey langfusev1alpha1.LangfuseAPIKey
	if err := r.Get(ctx, req.NamespacedName, &apiKey); err != nil {
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}

	if len(apiKey.Status.Conditions) > 0 && apiKey.Status.Conditions[0].Status == metav1.ConditionTrue {
		return ctrl.Result{}, nil
	}

	// Fetch Project
	var project langfusev1alpha1.LangfuseProject
	if err := r.Get(ctx, types.NamespacedName{Name: apiKey.Spec.ProjectRef, Namespace: req.Namespace}, &project); err != nil {
		log.Error(err, "Failed to get Project", "project", apiKey.Spec.ProjectRef)
		return ctrl.Result{}, err
	}

	if project.Status.ID == "" {
		log.Info("Project not ready yet", "project", project.Name)
		return ctrl.Result{Requeue: true}, nil
	}

	log.Info("Creating Langfuse API Key", "name", apiKey.Spec.Name, "projectID", project.Status.ID)
	lfAPIKey, err := r.LangfuseClient.CreateAPIKey(project.Status.ID, apiKey.Spec.Name)
	if err != nil {
		log.Error(err, "Failed to create API Key")
		return ctrl.Result{}, err
	}

	// Create Secret
	secret := &corev1.Secret{
		ObjectMeta: metav1.ObjectMeta{
			Name:      apiKey.Spec.SecretName,
			Namespace: apiKey.Namespace,
		},
		StringData: map[string]string{
			"LANGFUSE_PUBLIC_KEY": lfAPIKey.PublicKey,
			"LANGFUSE_SECRET_KEY": lfAPIKey.SecretKey,
			"LANGFUSE_HOST":       r.LangfuseClient.BaseURL,
		},
	}
	if err := ctrl.SetControllerReference(&apiKey, secret, r.Scheme); err != nil {
		return ctrl.Result{}, err
	}

	if err := r.Create(ctx, secret); err != nil {
		if client.IgnoreAlreadyExists(err) != nil {
			return ctrl.Result{}, err
		}
		// If exists, update? For now assume immutable or manual fix.
	}

	// Update Status
	apiKey.Status.Conditions = []metav1.Condition{
		{
			Type:               "Available",
			Status:             metav1.ConditionTrue,
			Reason:             "Created",
			Message:            "API Key created successfully",
			LastTransitionTime: metav1.Now(),
		},
	}
	if err := r.Status().Update(ctx, &apiKey); err != nil {
		return ctrl.Result{}, err
	}

	return ctrl.Result{}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *LangfuseAPIKeyReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&langfusev1alpha1.LangfuseAPIKey{}).
		WithOptions(controller.Options{
			MaxConcurrentReconciles: 1, // avoids event storms
		}).
		Named("langfuseapikey").
		Complete(r)
}
