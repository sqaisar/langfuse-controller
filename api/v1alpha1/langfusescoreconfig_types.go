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

// LangfuseScoreConfigSpec defines the desired state of LangfuseScoreConfig
type LangfuseScoreConfigSpec struct {
	// ProjectRef is the name of the LangfuseProject CR this score config belongs to.
	// +required
	ProjectRef string `json:"projectRef"`

	// Name is the name of the score config.
	// +required
	Name string `json:"name"`

	// DataType is the type of the score (e.g., NUMERIC, BOOLEAN).
	// +required
	DataType string `json:"dataType"`

	// MinValue is the minimum value for numeric scores.
	// +optional
	MinValue *int `json:"minValue,omitempty"`

	// MaxValue is the maximum value for numeric scores.
	// +optional
	MaxValue *int `json:"maxValue,omitempty"`

	// Categories are the allowed categories for categorical scores.
	// +optional
	Categories []string `json:"categories,omitempty"`
}

// LangfuseScoreConfigStatus defines the observed state of LangfuseScoreConfig.
type LangfuseScoreConfigStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "make" to regenerate code after modifying this file

	// For Kubernetes API conventions, see:
	// https://github.com/kubernetes/community/blob/master/contributors/devel/sig-architecture/api-conventions.md#typical-status-properties

	// conditions represent the current state of the LangfuseScoreConfig resource.
	// Each condition has a unique type and reflects the status of a specific aspect of the resource.
	//
	// Standard condition types include:
	// - "Available": the resource is fully functional
	// - "Progressing": the resource is being created or updated
	// - "Degraded": the resource failed to reach or maintain its desired state
	//
	// The status of each condition is one of True, False, or Unknown.
	// +listType=map
	// +listMapKey=type
	// +optional
	Conditions []metav1.Condition `json:"conditions,omitempty"`
}

// +kubebuilder:object:root=true
// +kubebuilder:subresource:status

// LangfuseScoreConfig is the Schema for the langfusescoreconfigs API
type LangfuseScoreConfig struct {
	metav1.TypeMeta `json:",inline"`

	// metadata is a standard object metadata
	// +optional
	metav1.ObjectMeta `json:"metadata,omitzero"`

	// spec defines the desired state of LangfuseScoreConfig
	// +required
	Spec LangfuseScoreConfigSpec `json:"spec"`

	// status defines the observed state of LangfuseScoreConfig
	// +optional
	Status LangfuseScoreConfigStatus `json:"status,omitzero"`
}

// +kubebuilder:object:root=true

// LangfuseScoreConfigList contains a list of LangfuseScoreConfig
type LangfuseScoreConfigList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitzero"`
	Items           []LangfuseScoreConfig `json:"items"`
}

func init() {
	SchemeBuilder.Register(&LangfuseScoreConfig{}, &LangfuseScoreConfigList{})
}
