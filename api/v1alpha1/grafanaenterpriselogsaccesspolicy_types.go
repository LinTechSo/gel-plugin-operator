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
// Desired AccessPolicy configurations
type TenantInfoRef struct {
	TenantName     string   `json:"tenantName,omitempty"`
	ClusterName    string   `json:"clusterName,omitempty"`
	AccessPolicies []string `json:"accessPolicies,omitempty"`
	LabelSelectors []string `json:"labelSelectors,omitempty"`
}

// MatchLabels Selector for finding the token secret
type MatchLabels struct {
	App string `json:"app,omitempty"`
}

// GrafanaEnterpriseLogsAccessPolicySpec defines the desired state of GrafanaEnterpriseLogsAccessPolicy
type GrafanaEnterpriseLogsAccessPolicySpec struct {
	// INSERT ADDITIONAL SPEC FIELDS - desired state of cluster
	// Important: Run "make" to regenerate code after modifying this file

	TenantInfoRef TenantInfoRef `json:"tenantInfoRef,omitempty"`
}

// GrafanaEnterpriseLogsAccessPolicyStatus defines the observed state of GrafanaEnterpriseLogsAccessPolicy
type GrafanaEnterpriseLogsAccessPolicyStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "make" to regenerate code after modifying this file
	Status string `json:"status,omitempty"`
}

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status

// GrafanaEnterpriseLogsAccessPolicy is the Schema for the grafanaenterpriselogsaccesspolicies API
type GrafanaEnterpriseLogsAccessPolicy struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`
	Finalizers        []string `json:"finalizers,omitempty"`

	Spec   GrafanaEnterpriseLogsAccessPolicySpec   `json:"spec,omitempty"`
	Status GrafanaEnterpriseLogsAccessPolicyStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

// GrafanaEnterpriseLogsAccessPolicyList contains a list of GrafanaEnterpriseLogsAccessPolicy
type GrafanaEnterpriseLogsAccessPolicyList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []GrafanaEnterpriseLogsAccessPolicy `json:"items"`
}

func init() {
	SchemeBuilder.Register(
		&GrafanaEnterpriseLogsAccessPolicy{},
		&GrafanaEnterpriseLogsAccessPolicyList{},
	)
}
