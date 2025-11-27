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

package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

// LangfuseProjectSpec defines the desired state of LangfuseProject
type LangfuseProjectSpec struct {
	// Name is the name of the project in Langfuse.
	// +required
	Name string `json:"name"`
}

// LangfuseProjectStatus defines the observed state of LangfuseProject.
type LangfuseProjectStatus struct {
	// ID is the unique identifier of the project in Langfuse.
	// +optional
	ID string `json:"id,omitempty"`

	// State represents the current state of the project (e.g., Ready, Error).
	// +optional
	State string `json:"state,omitempty"`

	// Conditions represent the latest available observations of an object's state
	// +optional
	Conditions []metav1.Condition `json:"conditions,omitempty"`
}

// +kubebuilder:object:root=true
// +kubebuilder:subresource:status

// LangfuseProject is the Schema for the langfuseprojects API
type LangfuseProject struct {
	metav1.TypeMeta `json:",inline"`

	// metadata is a standard object metadata
	// +optional
	metav1.ObjectMeta `json:"metadata,omitzero"`

	// spec defines the desired state of LangfuseProject
	// +required
	Spec LangfuseProjectSpec `json:"spec"`

	// status defines the observed state of LangfuseProject
	// +optional
	Status LangfuseProjectStatus `json:"status,omitzero"`
}

// +kubebuilder:object:root=true

// LangfuseProjectList contains a list of LangfuseProject
type LangfuseProjectList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitzero"`
	Items           []LangfuseProject `json:"items"`
}

func init() {
	SchemeBuilder.Register(&LangfuseProject{}, &LangfuseProjectList{})
}
