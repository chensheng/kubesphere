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

// PipelineTemplateSpec defines the desired state of PipelineTemplate
type PipelineTemplateSpec struct {
	// INSERT ADDITIONAL SPEC FIELDS - desired state of cluster
	// Important: Run "make" to regenerate code after modifying this file
	Desc                  string                 `json:"desc,omitempty" description:"Description of the pipeline template"`
	Mode                  string                 `json:"mode,moitempty" description:"Pipeline mode, e.g. java,nodejs"`
	Jenkinsfile           string                 `json:"jenkinsfile,omitempty" description:"Jenkinsfile content"`
	JenkinsfileParameters []JenkinsfileParameter `json:"jenkinsfileParameters,omitempty" description:"Pipeline jenkinsfile parameters"`
	BuildParameters       []Parameter            `json:"buildParameters,omitempty" description:"Pipeline build parameters"`
	Discarder             *DiscarderProperty     `json:"discarder,omitempty" description:"Pipeline build discarder"`
	DisableConcurrent     bool                   `json:"disableConcurrent,omitempty" description:"Whether to prohibit the pipeline from running in parallel"`
	TimerTrigger          *TimerTrigger          `json:"timerTrigger,omitempty" description:"Timer to trigger pipeline run"`
	RemoteTrigger         *RemoteTrigger         `json:"remoteTrigger,omitempty" description:"Remote api define to trigger pipeline run"`
}

type JenkinsfileParameter struct {
	Name         string `json:"name,omitempty" description:"Name of the template parameter"`
	Description  string `json:"description,omitempty" desription:"Description of the template parameter"`
	DefaultValue string `json:"defaultValue,omitempty" description:"Default value of the template parameter"`
}

// PipelineTemplateStatus defines the observed state of PipelineTemplate
type PipelineTemplateStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "make" to regenerate code after modifying this file
}

// +genclient
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// PipelineTemplate is the Schema for the pipelinetemplates API
// +k8s:openapi-gen=true
type PipelineTemplate struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   PipelineTemplateSpec   `json:"spec,omitempty"`
	Status PipelineTemplateStatus `json:"status,omitempty"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// PipelineTemplateList contains a list of PipelineTemplate
type PipelineTemplateList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []PipelineTemplate `json:"items"`
}

func init() {
	SchemeBuilder.Register(&PipelineTemplate{}, &PipelineTemplateList{})
}
