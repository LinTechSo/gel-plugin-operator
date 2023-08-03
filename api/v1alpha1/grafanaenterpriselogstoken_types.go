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

type AccessPolicyRef struct {
	Name string `json:"name,omitempty"`
}

// GrafanaEnterpriseLogsTokenSpec defines the desired state of GrafanaEnterpriseLogsToken
type GrafanaEnterpriseLogsTokenSpec struct {
	// INSERT ADDITIONAL SPEC FIELDS - desired state of cluster
	// Important: Run "make" to regenerate code after modifying this file

	AccessPolicyRef AccessPolicyRef `json:"accessPolicyRef,omitempty"`
	ExpirationTime  string          `json:"expirationTime,omitempty"`
}

// GrafanaEnterpriseLogsTokenStatus defines the observed state of GrafanaEnterpriseLogsToken
type GrafanaEnterpriseLogsTokenStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "make" to regenerate code after modifying this file
	Status string `json:"status,omitempty"`
}

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status

// GrafanaEnterpriseLogsToken is the Schema for the grafanaenterpriselogstokens API
type GrafanaEnterpriseLogsToken struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`
	Finalizers        []string `json:"finalizers,omitempty"`

	Spec   GrafanaEnterpriseLogsTokenSpec   `json:"spec,omitempty"`
	Status GrafanaEnterpriseLogsTokenStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

// GrafanaEnterpriseLogsTokenList contains a list of GrafanaEnterpriseLogsToken
type GrafanaEnterpriseLogsTokenList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []GrafanaEnterpriseLogsToken `json:"items"`
}

func init() {
	SchemeBuilder.Register(&GrafanaEnterpriseLogsToken{}, &GrafanaEnterpriseLogsTokenList{})
}
