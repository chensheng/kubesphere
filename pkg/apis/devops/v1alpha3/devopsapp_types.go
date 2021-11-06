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
	"k8s.io/apimachinery/pkg/util/intstr"
)

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

// DevOpsAppSpec defines the desired state of DevOpsApp
type DevOpsAppSpec struct {
	// INSERT ADDITIONAL SPEC FIELDS - desired state of cluster
	// Important: Run "make" to regenerate code after modifying this file
	AutoUpdate    bool          `json:"autoUpdate,omitempty" description:"Whether to update related resources when DevOpsApp CRD changed"`
	Desc          string        `json:"desc,omitempty" description:"Description of the DevOpsApp"`
	Mode          string        `json:"mode,omitempty" description:"DevOps application mode: java, nodejs"`
	Registry      *Registry     `json:"registry,omitempty" description:"Docker registry for this DevOpsApp"`
	Git           *Git          `json:"git,omitempty" description:"Git information for this DevOpsApp"`
	ConfigCenter  *ConfigCenter `json:"configCenter,omitempty" description:"Configuration center information"`
	Environments  []Environment `json:"environments,omitempty" description:"Environments holds all environment for  this DevOpsApp"`
	DevOpsProject string        `json:"devopsproject,omitemty" description:"Name of related DevOpsProject"`
}

type Registry struct {
	Url      string `json:"url,omitempty" description:"URL of this docker registry"`
	Username string `json:"username,omitempty" description:"Username of this docker registry"`
	Password string `json:"password,omitempty" description:"Password of this docker registry"`
	Email    string `json:"email,omitempty" description:"Email of this docker registry"`
}

type Git struct {
	Repo       string `json:"repo,omitempty" description:"Git repository for this DevOpsApp"`
	Username   string `json:"username,omitempty" description:"Username of git repository"`
	Password   string `json:"password,omitempty" description:"Password of git repository"`
	TargetPath string `json:"targetPath,omitempty" description:"Build target path of git repository"`
}

type ConfigCenter struct {
	Url        string `json:"url,omitempty" description:"URL for configuration center"`
	ClusterUrl string `json:"clusterUrl,omitempty" description:"Cluster URL for the configuration center"`
	Username   string `json:"username,omitempty" description:"Username of configuration center"`
	Password   string `json:"password,omitempty" description:"Password of configuration center"`
}

type Environment struct {
	Namespace    string             `json:"namespace,omitempty" description:"Namespace of the environment"`
	Name         string             `json:"name,omitempty" description:"Name of the environment"`
	Desc         string             `json:"desc,omitempty" description:"Description for the environment"`
	Cluster      string             `json:"cluster,omitempty" description:"Cluster of the environment"`
	Gateway      string             `json:"gateway,omitempty" description:"Gateway of the environment"`
	ConfigCenter *ConfigCenter      `json:"configCenter,omitempty" description:"Configuration center information for the environment"`
	Registry     *Registry          `json:"registry,omitempty" description:"Docker registry for the environment"`
	Resource     *EnvResource       `json:"resource,omitempty" description:"Resource of the environment"`
	Dependencies []Dependency       `json:"dependencies,omitempty" description:"Dependencies of the environment"`
	Credentials  []Credential       `json:"credentials,omitempty" description:"Credentials of the environment"`
	Pipeline     PipelineDefinition `json:"pipeline,omitempty" description:"Pipeline definition of the environment"`
}

type EnvResource struct {
	Replicas int                 `json:"replicas,omitempty" description:"Replicas of the environment"`
	Memory   string              `json:"memory,omitempty" description:"Max memory of the environment"`
	CPU      *intstr.IntOrString `json:"cpu,omitempty" description:"'Max cpu of the environment"`
}

type Dependency struct {
	Category   string `json:"category,omitempty" description:"Category of the dependency, e.g. mysql, mongodb, redis, seata, rocketmq, kafka"`
	Desc       string `json:"desc,omitempty" description:"Description of the dependency"`
	Url        string `json:"url,omitempty" description:"URL of the dependency"`
	ClusterUrl string `json:"clusterUrl,omitempty" description:"Cluster URL of the dependency"`
	Username   string `json:"username,omitempty" description:"Username of the dependency"`
	Password   string `json:"password,omitempty" description:"Password of the dependency"`
	Database   string `json:"database,omitempty" description:"Database of the dependency"`
}

type Credential struct {
	Name       string `json:"name,omitempty" description:"Credentail name"`
	Type       string `json:"type,omitempty" description:"Credential type, e.g. kubeconfig, basic-auth, secret-text, ssh-auth"`
	Content    string `json:"content,omitempty" description:"Content of kubeconfig"`
	Username   string `json:"username,omitempty" description:"Username of basic-auth or ssh-auth"`
	Password   string `json:"password,omitempty" description:"Password of basic-auth or ssh-auth"`
	Secret     string `json:"secret,omitempty" description:"Secret of secret-text"`
	PrivateKey string `json:"privateKey,omitempty" description:"Private key of ssh-auth"`
}

type PipelineDefinition struct {
	Name         string            `json:"name" description:"Pipeline name"`
	TemplateName string            `json:"templateName,omitempty" description:"Pipeline template name"`
	Parameters   map[string]string `json:"parameters,omitempty" description:"Pipeline parameters"`
}

// DevOpsAppStatus defines the observed state of DevOpsApp
type DevOpsAppStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "make" to regenerate code after modifying this file
}

// +genclient
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// DevOpsApp is the Schema for the devopsapps API
// +k8s:openapi-gen=true
type DevOpsApp struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   DevOpsAppSpec   `json:"spec,omitempty"`
	Status DevOpsAppStatus `json:"status,omitempty"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// DevOpsAppList contains a list of DevOpsApp
type DevOpsAppList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []DevOpsApp `json:"items"`
}

func init() {
	SchemeBuilder.Register(&DevOpsApp{}, &DevOpsAppList{})
}
