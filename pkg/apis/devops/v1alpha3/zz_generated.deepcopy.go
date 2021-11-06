// +build !ignore_autogenerated

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

// Code generated by controller-gen. DO NOT EDIT.

package v1alpha3

import (
	runtime "k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/util/intstr"
)

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *BitbucketServerSource) DeepCopyInto(out *BitbucketServerSource) {
	*out = *in
	if in.DiscoverPRFromForks != nil {
		in, out := &in.DiscoverPRFromForks, &out.DiscoverPRFromForks
		*out = new(DiscoverPRFromForks)
		**out = **in
	}
	if in.CloneOption != nil {
		in, out := &in.CloneOption, &out.CloneOption
		*out = new(GitCloneOption)
		**out = **in
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new BitbucketServerSource.
func (in *BitbucketServerSource) DeepCopy() *BitbucketServerSource {
	if in == nil {
		return nil
	}
	out := new(BitbucketServerSource)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *ConfigCenter) DeepCopyInto(out *ConfigCenter) {
	*out = *in
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new ConfigCenter.
func (in *ConfigCenter) DeepCopy() *ConfigCenter {
	if in == nil {
		return nil
	}
	out := new(ConfigCenter)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *Credential) DeepCopyInto(out *Credential) {
	*out = *in
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new Credential.
func (in *Credential) DeepCopy() *Credential {
	if in == nil {
		return nil
	}
	out := new(Credential)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *Dependency) DeepCopyInto(out *Dependency) {
	*out = *in
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new Dependency.
func (in *Dependency) DeepCopy() *Dependency {
	if in == nil {
		return nil
	}
	out := new(Dependency)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *DevOpsApp) DeepCopyInto(out *DevOpsApp) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ObjectMeta.DeepCopyInto(&out.ObjectMeta)
	in.Spec.DeepCopyInto(&out.Spec)
	out.Status = in.Status
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new DevOpsApp.
func (in *DevOpsApp) DeepCopy() *DevOpsApp {
	if in == nil {
		return nil
	}
	out := new(DevOpsApp)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *DevOpsApp) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *DevOpsAppList) DeepCopyInto(out *DevOpsAppList) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ListMeta.DeepCopyInto(&out.ListMeta)
	if in.Items != nil {
		in, out := &in.Items, &out.Items
		*out = make([]DevOpsApp, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new DevOpsAppList.
func (in *DevOpsAppList) DeepCopy() *DevOpsAppList {
	if in == nil {
		return nil
	}
	out := new(DevOpsAppList)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *DevOpsAppList) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *DevOpsAppSpec) DeepCopyInto(out *DevOpsAppSpec) {
	*out = *in
	if in.Registry != nil {
		in, out := &in.Registry, &out.Registry
		*out = new(Registry)
		**out = **in
	}
	if in.Git != nil {
		in, out := &in.Git, &out.Git
		*out = new(Git)
		**out = **in
	}
	if in.ConfigCenter != nil {
		in, out := &in.ConfigCenter, &out.ConfigCenter
		*out = new(ConfigCenter)
		**out = **in
	}
	if in.Environments != nil {
		in, out := &in.Environments, &out.Environments
		*out = make([]Environment, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new DevOpsAppSpec.
func (in *DevOpsAppSpec) DeepCopy() *DevOpsAppSpec {
	if in == nil {
		return nil
	}
	out := new(DevOpsAppSpec)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *DevOpsAppStatus) DeepCopyInto(out *DevOpsAppStatus) {
	*out = *in
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new DevOpsAppStatus.
func (in *DevOpsAppStatus) DeepCopy() *DevOpsAppStatus {
	if in == nil {
		return nil
	}
	out := new(DevOpsAppStatus)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *DevOpsProject) DeepCopyInto(out *DevOpsProject) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ObjectMeta.DeepCopyInto(&out.ObjectMeta)
	out.Spec = in.Spec
	out.Status = in.Status
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new DevOpsProject.
func (in *DevOpsProject) DeepCopy() *DevOpsProject {
	if in == nil {
		return nil
	}
	out := new(DevOpsProject)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *DevOpsProject) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *DevOpsProjectList) DeepCopyInto(out *DevOpsProjectList) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ListMeta.DeepCopyInto(&out.ListMeta)
	if in.Items != nil {
		in, out := &in.Items, &out.Items
		*out = make([]DevOpsProject, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new DevOpsProjectList.
func (in *DevOpsProjectList) DeepCopy() *DevOpsProjectList {
	if in == nil {
		return nil
	}
	out := new(DevOpsProjectList)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *DevOpsProjectList) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *DevOpsProjectSpec) DeepCopyInto(out *DevOpsProjectSpec) {
	*out = *in
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new DevOpsProjectSpec.
func (in *DevOpsProjectSpec) DeepCopy() *DevOpsProjectSpec {
	if in == nil {
		return nil
	}
	out := new(DevOpsProjectSpec)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *DevOpsProjectStatus) DeepCopyInto(out *DevOpsProjectStatus) {
	*out = *in
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new DevOpsProjectStatus.
func (in *DevOpsProjectStatus) DeepCopy() *DevOpsProjectStatus {
	if in == nil {
		return nil
	}
	out := new(DevOpsProjectStatus)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *DiscarderProperty) DeepCopyInto(out *DiscarderProperty) {
	*out = *in
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new DiscarderProperty.
func (in *DiscarderProperty) DeepCopy() *DiscarderProperty {
	if in == nil {
		return nil
	}
	out := new(DiscarderProperty)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *DiscoverPRFromForks) DeepCopyInto(out *DiscoverPRFromForks) {
	*out = *in
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new DiscoverPRFromForks.
func (in *DiscoverPRFromForks) DeepCopy() *DiscoverPRFromForks {
	if in == nil {
		return nil
	}
	out := new(DiscoverPRFromForks)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *EnvResource) DeepCopyInto(out *EnvResource) {
	*out = *in
	if in.CPU != nil {
		in, out := &in.CPU, &out.CPU
		*out = new(intstr.IntOrString)
		**out = **in
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new EnvResource.
func (in *EnvResource) DeepCopy() *EnvResource {
	if in == nil {
		return nil
	}
	out := new(EnvResource)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *Environment) DeepCopyInto(out *Environment) {
	*out = *in
	if in.ConfigCenter != nil {
		in, out := &in.ConfigCenter, &out.ConfigCenter
		*out = new(ConfigCenter)
		**out = **in
	}
	if in.Registry != nil {
		in, out := &in.Registry, &out.Registry
		*out = new(Registry)
		**out = **in
	}
	if in.Resource != nil {
		in, out := &in.Resource, &out.Resource
		*out = new(EnvResource)
		(*in).DeepCopyInto(*out)
	}
	if in.Dependencies != nil {
		in, out := &in.Dependencies, &out.Dependencies
		*out = make([]Dependency, len(*in))
		copy(*out, *in)
	}
	if in.Credentials != nil {
		in, out := &in.Credentials, &out.Credentials
		*out = make([]Credential, len(*in))
		copy(*out, *in)
	}
	in.Pipeline.DeepCopyInto(&out.Pipeline)
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new Environment.
func (in *Environment) DeepCopy() *Environment {
	if in == nil {
		return nil
	}
	out := new(Environment)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *Git) DeepCopyInto(out *Git) {
	*out = *in
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new Git.
func (in *Git) DeepCopy() *Git {
	if in == nil {
		return nil
	}
	out := new(Git)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *GitCloneOption) DeepCopyInto(out *GitCloneOption) {
	*out = *in
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new GitCloneOption.
func (in *GitCloneOption) DeepCopy() *GitCloneOption {
	if in == nil {
		return nil
	}
	out := new(GitCloneOption)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *GitSource) DeepCopyInto(out *GitSource) {
	*out = *in
	if in.CloneOption != nil {
		in, out := &in.CloneOption, &out.CloneOption
		*out = new(GitCloneOption)
		**out = **in
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new GitSource.
func (in *GitSource) DeepCopy() *GitSource {
	if in == nil {
		return nil
	}
	out := new(GitSource)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *GithubSource) DeepCopyInto(out *GithubSource) {
	*out = *in
	if in.DiscoverPRFromForks != nil {
		in, out := &in.DiscoverPRFromForks, &out.DiscoverPRFromForks
		*out = new(DiscoverPRFromForks)
		**out = **in
	}
	if in.CloneOption != nil {
		in, out := &in.CloneOption, &out.CloneOption
		*out = new(GitCloneOption)
		**out = **in
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new GithubSource.
func (in *GithubSource) DeepCopy() *GithubSource {
	if in == nil {
		return nil
	}
	out := new(GithubSource)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *GitlabSource) DeepCopyInto(out *GitlabSource) {
	*out = *in
	if in.DiscoverPRFromForks != nil {
		in, out := &in.DiscoverPRFromForks, &out.DiscoverPRFromForks
		*out = new(DiscoverPRFromForks)
		**out = **in
	}
	if in.CloneOption != nil {
		in, out := &in.CloneOption, &out.CloneOption
		*out = new(GitCloneOption)
		**out = **in
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new GitlabSource.
func (in *GitlabSource) DeepCopy() *GitlabSource {
	if in == nil {
		return nil
	}
	out := new(GitlabSource)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *JenkinsfileParameter) DeepCopyInto(out *JenkinsfileParameter) {
	*out = *in
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new JenkinsfileParameter.
func (in *JenkinsfileParameter) DeepCopy() *JenkinsfileParameter {
	if in == nil {
		return nil
	}
	out := new(JenkinsfileParameter)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *MultiBranchJobTrigger) DeepCopyInto(out *MultiBranchJobTrigger) {
	*out = *in
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new MultiBranchJobTrigger.
func (in *MultiBranchJobTrigger) DeepCopy() *MultiBranchJobTrigger {
	if in == nil {
		return nil
	}
	out := new(MultiBranchJobTrigger)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *MultiBranchPipeline) DeepCopyInto(out *MultiBranchPipeline) {
	*out = *in
	if in.Discarder != nil {
		in, out := &in.Discarder, &out.Discarder
		*out = new(DiscarderProperty)
		**out = **in
	}
	if in.TimerTrigger != nil {
		in, out := &in.TimerTrigger, &out.TimerTrigger
		*out = new(TimerTrigger)
		**out = **in
	}
	if in.GitSource != nil {
		in, out := &in.GitSource, &out.GitSource
		*out = new(GitSource)
		(*in).DeepCopyInto(*out)
	}
	if in.GitHubSource != nil {
		in, out := &in.GitHubSource, &out.GitHubSource
		*out = new(GithubSource)
		(*in).DeepCopyInto(*out)
	}
	if in.GitlabSource != nil {
		in, out := &in.GitlabSource, &out.GitlabSource
		*out = new(GitlabSource)
		(*in).DeepCopyInto(*out)
	}
	if in.SvnSource != nil {
		in, out := &in.SvnSource, &out.SvnSource
		*out = new(SvnSource)
		**out = **in
	}
	if in.SingleSvnSource != nil {
		in, out := &in.SingleSvnSource, &out.SingleSvnSource
		*out = new(SingleSvnSource)
		**out = **in
	}
	if in.BitbucketServerSource != nil {
		in, out := &in.BitbucketServerSource, &out.BitbucketServerSource
		*out = new(BitbucketServerSource)
		(*in).DeepCopyInto(*out)
	}
	if in.MultiBranchJobTrigger != nil {
		in, out := &in.MultiBranchJobTrigger, &out.MultiBranchJobTrigger
		*out = new(MultiBranchJobTrigger)
		**out = **in
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new MultiBranchPipeline.
func (in *MultiBranchPipeline) DeepCopy() *MultiBranchPipeline {
	if in == nil {
		return nil
	}
	out := new(MultiBranchPipeline)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *NoScmPipeline) DeepCopyInto(out *NoScmPipeline) {
	*out = *in
	if in.Discarder != nil {
		in, out := &in.Discarder, &out.Discarder
		*out = new(DiscarderProperty)
		**out = **in
	}
	if in.Parameters != nil {
		in, out := &in.Parameters, &out.Parameters
		*out = make([]Parameter, len(*in))
		copy(*out, *in)
	}
	if in.TimerTrigger != nil {
		in, out := &in.TimerTrigger, &out.TimerTrigger
		*out = new(TimerTrigger)
		**out = **in
	}
	if in.RemoteTrigger != nil {
		in, out := &in.RemoteTrigger, &out.RemoteTrigger
		*out = new(RemoteTrigger)
		**out = **in
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new NoScmPipeline.
func (in *NoScmPipeline) DeepCopy() *NoScmPipeline {
	if in == nil {
		return nil
	}
	out := new(NoScmPipeline)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *Parameter) DeepCopyInto(out *Parameter) {
	*out = *in
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new Parameter.
func (in *Parameter) DeepCopy() *Parameter {
	if in == nil {
		return nil
	}
	out := new(Parameter)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *Pipeline) DeepCopyInto(out *Pipeline) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ObjectMeta.DeepCopyInto(&out.ObjectMeta)
	in.Spec.DeepCopyInto(&out.Spec)
	out.Status = in.Status
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new Pipeline.
func (in *Pipeline) DeepCopy() *Pipeline {
	if in == nil {
		return nil
	}
	out := new(Pipeline)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *Pipeline) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *PipelineDefinition) DeepCopyInto(out *PipelineDefinition) {
	*out = *in
	if in.Parameters != nil {
		in, out := &in.Parameters, &out.Parameters
		*out = make(map[string]string, len(*in))
		for key, val := range *in {
			(*out)[key] = val
		}
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new PipelineDefinition.
func (in *PipelineDefinition) DeepCopy() *PipelineDefinition {
	if in == nil {
		return nil
	}
	out := new(PipelineDefinition)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *PipelineList) DeepCopyInto(out *PipelineList) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ListMeta.DeepCopyInto(&out.ListMeta)
	if in.Items != nil {
		in, out := &in.Items, &out.Items
		*out = make([]Pipeline, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new PipelineList.
func (in *PipelineList) DeepCopy() *PipelineList {
	if in == nil {
		return nil
	}
	out := new(PipelineList)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *PipelineList) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *PipelineSpec) DeepCopyInto(out *PipelineSpec) {
	*out = *in
	if in.Pipeline != nil {
		in, out := &in.Pipeline, &out.Pipeline
		*out = new(NoScmPipeline)
		(*in).DeepCopyInto(*out)
	}
	if in.MultiBranchPipeline != nil {
		in, out := &in.MultiBranchPipeline, &out.MultiBranchPipeline
		*out = new(MultiBranchPipeline)
		(*in).DeepCopyInto(*out)
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new PipelineSpec.
func (in *PipelineSpec) DeepCopy() *PipelineSpec {
	if in == nil {
		return nil
	}
	out := new(PipelineSpec)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *PipelineStatus) DeepCopyInto(out *PipelineStatus) {
	*out = *in
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new PipelineStatus.
func (in *PipelineStatus) DeepCopy() *PipelineStatus {
	if in == nil {
		return nil
	}
	out := new(PipelineStatus)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *PipelineTemplate) DeepCopyInto(out *PipelineTemplate) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ObjectMeta.DeepCopyInto(&out.ObjectMeta)
	in.Spec.DeepCopyInto(&out.Spec)
	out.Status = in.Status
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new PipelineTemplate.
func (in *PipelineTemplate) DeepCopy() *PipelineTemplate {
	if in == nil {
		return nil
	}
	out := new(PipelineTemplate)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *PipelineTemplate) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *PipelineTemplateList) DeepCopyInto(out *PipelineTemplateList) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ListMeta.DeepCopyInto(&out.ListMeta)
	if in.Items != nil {
		in, out := &in.Items, &out.Items
		*out = make([]PipelineTemplate, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new PipelineTemplateList.
func (in *PipelineTemplateList) DeepCopy() *PipelineTemplateList {
	if in == nil {
		return nil
	}
	out := new(PipelineTemplateList)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *PipelineTemplateList) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *PipelineTemplateSpec) DeepCopyInto(out *PipelineTemplateSpec) {
	*out = *in
	if in.JenkinsfileParameters != nil {
		in, out := &in.JenkinsfileParameters, &out.JenkinsfileParameters
		*out = make([]JenkinsfileParameter, len(*in))
		copy(*out, *in)
	}
	if in.BuildParameters != nil {
		in, out := &in.BuildParameters, &out.BuildParameters
		*out = make([]Parameter, len(*in))
		copy(*out, *in)
	}
	if in.Discarder != nil {
		in, out := &in.Discarder, &out.Discarder
		*out = new(DiscarderProperty)
		**out = **in
	}
	if in.TimerTrigger != nil {
		in, out := &in.TimerTrigger, &out.TimerTrigger
		*out = new(TimerTrigger)
		**out = **in
	}
	if in.RemoteTrigger != nil {
		in, out := &in.RemoteTrigger, &out.RemoteTrigger
		*out = new(RemoteTrigger)
		**out = **in
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new PipelineTemplateSpec.
func (in *PipelineTemplateSpec) DeepCopy() *PipelineTemplateSpec {
	if in == nil {
		return nil
	}
	out := new(PipelineTemplateSpec)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *PipelineTemplateStatus) DeepCopyInto(out *PipelineTemplateStatus) {
	*out = *in
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new PipelineTemplateStatus.
func (in *PipelineTemplateStatus) DeepCopy() *PipelineTemplateStatus {
	if in == nil {
		return nil
	}
	out := new(PipelineTemplateStatus)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *Registry) DeepCopyInto(out *Registry) {
	*out = *in
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new Registry.
func (in *Registry) DeepCopy() *Registry {
	if in == nil {
		return nil
	}
	out := new(Registry)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *RemoteTrigger) DeepCopyInto(out *RemoteTrigger) {
	*out = *in
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new RemoteTrigger.
func (in *RemoteTrigger) DeepCopy() *RemoteTrigger {
	if in == nil {
		return nil
	}
	out := new(RemoteTrigger)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *SingleSvnSource) DeepCopyInto(out *SingleSvnSource) {
	*out = *in
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new SingleSvnSource.
func (in *SingleSvnSource) DeepCopy() *SingleSvnSource {
	if in == nil {
		return nil
	}
	out := new(SingleSvnSource)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *SvnSource) DeepCopyInto(out *SvnSource) {
	*out = *in
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new SvnSource.
func (in *SvnSource) DeepCopy() *SvnSource {
	if in == nil {
		return nil
	}
	out := new(SvnSource)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *TimerTrigger) DeepCopyInto(out *TimerTrigger) {
	*out = *in
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new TimerTrigger.
func (in *TimerTrigger) DeepCopy() *TimerTrigger {
	if in == nil {
		return nil
	}
	out := new(TimerTrigger)
	in.DeepCopyInto(out)
	return out
}
