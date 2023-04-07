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

package v1beta1

import (
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

// SlurmApplicationSpec defines the desired state of SlurmApplication
type SlurmApplicationSpec struct {
	// INSERT ADDITIONAL SPEC FIELDS - desired state of cluster
	// Important: Run "make" to regenerate code after modifying this file

	// Foo is an example field of SlurmApplication. Edit SlurmApplication_types.go to remove/update
	Jupyter *SlurmJupyterSpec `json:"jupyter"`
	Master  *SlurmMasterSpec  `json:"master"`
	Node    *SlurmNodeSpec    `json:"node"`
}

type SlurmJupyterSpec struct {
	CommonSpec `json:".,inline"`
}

type SlurmMasterSpec struct {
	CommonSpec `json:".,inline"`
}

type SlurmNodeSpec struct {
	CommonSpec `json:".,inline"`
}

type CommonSpec struct {
	// +kubebuilder:validation:MinLength=1
	Image string `json:"image"`

	// +kubebuilder:validation:Minimum=1
	// +kubebuilder:default=1
	Instance int `json:"instance,omitempty"`

	Resources    *corev1.ResourceRequirements `json:"resources,omitempty"`
	Labels       *metav1.LabelSelector        `json:"labels,omitempty"`
	NodeSelector *corev1.NodeSelector         `json:"nodeSelector,omitempty"`
}

// SlurmApplicationStatus defines the observed state of SlurmApplication
type SlurmApplicationStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "make" to regenerate code after modifying this file
}

// +kubebuilder:object:root=true

// SlurmApplication is the Schema for the slurmapplications API
type SlurmApplication struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   SlurmApplicationSpec   `json:"spec,omitempty"`
	Status SlurmApplicationStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true

// SlurmApplicationList contains a list of SlurmApplication
type SlurmApplicationList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []SlurmApplication `json:"items"`
}

func init() {
	SchemeBuilder.Register(&SlurmApplication{}, &SlurmApplicationList{})
}
