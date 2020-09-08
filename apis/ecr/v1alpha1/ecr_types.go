/*
Copyright 2020 The Crossplane Authors.

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

	runtimev1alpha1 "github.com/crossplane/crossplane-runtime/apis/core/v1alpha1"
)

// RegistryParameters defines the desired state of an AWS Elastic Container Registry.
type RegistryParameters struct {
	ImmutableTags bool `json:"immutableTags"`
	ScanOnPush    bool `json:"scanOnPush"`
}

// An RegistrySpec defines the desired state of an Elastic Container Registry.
type RegistrySpec struct {
	runtimev1alpha1.ResourceSpec `json:",inline"`
	RegistryParameters           `json:",inline"`
}

// RegistryStatus defines the observed state of the Elastic Container Registry.
type RegistryStatus struct {
	runtimev1alpha1.ResourceStatus `json:",inline"`
	RegistryID                     string `json:"registryID,omitempty"`
	RepositoryName                 string `json:"repositoryName,omitempty"`
	RepositoryURI                  string `json:"repositoryURI,omitempty"`
}

// +kubebuilder:object:root=true

// Registry is a managed resource that represents an AWS Elastic Container Registry.
// +kubebuilder:printcolumn:name="NAME",type="string",JSONPath=".statuc.repositoryName"
// +kubebuilder:printcolumn:name="URI",type="string",JSONPath=".status.repositoryURI"
// +kubebuilder:printcolumn:name="READY",type="string",JSONPath=".status.conditions[?(@.type=='Ready')].status"
// +kubebuilder:printcolumn:name="SYNCED",type="string",JSONPath=".status.conditions[?(@.type=='Synced')].status"
// +kubebuilder:printcolumn:name="AGE",type="date",JSONPath=".metadata.creationTimestamp"
// +kubebuilder:resource:scope=Cluster,categories={crossplane,managed,aws}
type Registry struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   RegistrySpec   `json:"spec"`
	Status RegistryStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true

// RegistryList contains a list of Registries
type RegistryList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []Registry `json:"items"`
}
