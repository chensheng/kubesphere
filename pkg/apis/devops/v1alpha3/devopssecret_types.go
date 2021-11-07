/*
Copyright 2020 The KubeSphere Authors.

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

package v1alpha3

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

// DevOpsSecretSpec defines the desired state of DevOpsSecret
type DevOpsSecretSpec struct {
	Type       string `json:"type,omitempty" description:"Credential type, e.g. kubeconfig, basic-auth, secret-text, ssh-auth"`
	Content    string `json:"content,omitempty" description:"Content of kubeconfig"`
	Username   string `json:"username,omitempty" description:"Username of basic-auth or ssh-auth"`
	Password   string `json:"password,omitempty" description:"Password of basic-auth or ssh-auth"`
	Secret     string `json:"secret,omitempty" description:"Secret of secret-text"`
	PrivateKey string `json:"privateKey,omitempty" description:"Private key of ssh-auth"`
}

// DevOpsSecretStatus defines the observed state of DevOpsSecret
type DevOpsSecretStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "make" to regenerate code after modifying this file
}

// +genclient
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// DevOpsSecret is the Schema for the devopssecrets API
// +k8s:openapi-gen=true
type DevOpsSecret struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   DevOpsSecretSpec   `json:"spec,omitempty"`
	Status DevOpsSecretStatus `json:"status,omitempty"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// DevOpsSecretList contains a list of DevOpsSecret
type DevOpsSecretList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []DevOpsSecret `json:"items"`
}

func init() {
	SchemeBuilder.Register(&DevOpsSecret{}, &DevOpsSecretList{})
}
