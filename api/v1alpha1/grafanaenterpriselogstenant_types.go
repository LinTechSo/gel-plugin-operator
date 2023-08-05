/*
Copyright 2023.

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
type TenantInfo struct {
	Name        string `json:"name,omitempty"`
	DisplayName string `json:"displayName,omitempty"`
	ClusterName string `json:"clusterName,omitempty"`
}

// GrafanaEnterpriseLogsTenantSpec defines the desired state of GrafanaEnterpriseLogsTenant
type GrafanaEnterpriseLogsTenantSpec struct {
	// INSERT ADDITIONAL SPEC FIELDS - desired state of cluster
	// Important: Run "make" to regenerate code after modifying this file

	TenantInfo TenantInfo `json:"tenantInfo,omitempty"`
}

// GrafanaEnterpriseLogsTenantStatus defines the observed state of GrafanaEnterpriseLogsTenant
type GrafanaEnterpriseLogsTenantStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "make" to regenerate code after modifying this file
	Status string `json:"status,omitempty"`
}

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status

// GrafanaEnterpriseLogsTenant is the Schema for the grafanaenterpriselogstenants API
type GrafanaEnterpriseLogsTenant struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`
	Finalizers        []string `json:"finalizers,omitempty"`

	Spec   GrafanaEnterpriseLogsTenantSpec   `json:"spec,omitempty"`
	Status GrafanaEnterpriseLogsTenantStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

// GrafanaEnterpriseLogsTenantList contains a list of GrafanaEnterpriseLogsTenant
type GrafanaEnterpriseLogsTenantList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []GrafanaEnterpriseLogsTenant `json:"items"`
}

func init() {
	SchemeBuilder.Register(
		&GrafanaEnterpriseLogsTenant{},
		&GrafanaEnterpriseLogsTenantList{},
	)
}
