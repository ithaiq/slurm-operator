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
	Jupyter SlurmJupyterSpec `json:"jupyter"`
	Master  SlurmMasterSpec  `json:"master"`
	Node    SlurmNodeSpec    `json:"node"`
}

type SlurmJupyterSpec struct {
	CommonSpec `json:",inline"`
}

type SlurmMasterSpec struct {
	CommonSpec `json:",inline"`
}

type SlurmNodeSpec struct {
	CommonSpec `json:",inline"`
}

type CommonSpec struct {
	// +kubebuilder:validation:MinLength=1
	Image string `json:"image"`

	// +kubebuilder:validation:Minimum=1
	// +kubebuilder:default=1
	Instance int `json:"instance,omitempty"`

	Resources    corev1.ResourceRequirements `json:"resources,omitempty"`
	Labels       map[string]string           `json:"labels,omitempty"`
	NodeSelector map[string]string           `json:"nodeSelector,omitempty"`
}

type SlurmClusterPhase string

var (
	SlurmClusterStatusRunning SlurmClusterPhase = "Running"
	SlurmClusterStatusError   SlurmClusterPhase = "Error"
)

// +kubebuilder:object:root=true
// +kubebuilder:printcolumn:name="Status",type="string",JSONPath=".status",description="Slurm app status"
// +kubebuilder:printcolumn:name="NodeInstance",type="integer",JSONPath=".spec.node.instance",description="The Num of Slurm node"
// +kubebuilder:printcolumn:name="JupyterImage",type="string",JSONPath=".spec.jupyter.image",description="The Docker Image of jupyter"
// +kubebuilder:printcolumn:name="SlurmMasterImage",type="string",JSONPath=".spec.master.image",description="The Docker Image of Slurm master"
// +kubebuilder:printcolumn:name="SlurmNodeImage",type="string",JSONPath=".spec.node.image",description="The Docker Image of Slurm node"
// +kubebuilder:printcolumn:name="Age",type="date",JSONPath=".metadata.creationTimestamp"
// +kubebuilder:subresource:status

// SlurmApplication is the Schema for the slurmapplications API
type SlurmApplication struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   SlurmApplicationSpec `json:"spec,omitempty"`
	Status SlurmClusterPhase    `json:"status,omitempty"`
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
