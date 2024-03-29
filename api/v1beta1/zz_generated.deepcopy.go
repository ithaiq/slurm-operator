//go:build !ignore_autogenerated
// +build !ignore_autogenerated

/*
Copyright 2023 xxx.

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

package v1beta1

import (
	"k8s.io/api/core/v1"
	runtime "k8s.io/apimachinery/pkg/runtime"
)

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *CommonSpec) DeepCopyInto(out *CommonSpec) {
	*out = *in
	in.Resources.DeepCopyInto(&out.Resources)
	if in.Labels != nil {
		in, out := &in.Labels, &out.Labels
		*out = make(map[string]string, len(*in))
		for key, val := range *in {
			(*out)[key] = val
		}
	}
	if in.NodeSelector != nil {
		in, out := &in.NodeSelector, &out.NodeSelector
		*out = make(map[string]string, len(*in))
		for key, val := range *in {
			(*out)[key] = val
		}
	}
	if in.VolumeMounts != nil {
		in, out := &in.VolumeMounts, &out.VolumeMounts
		*out = make([]v1.VolumeMount, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new CommonSpec.
func (in *CommonSpec) DeepCopy() *CommonSpec {
	if in == nil {
		return nil
	}
	out := new(CommonSpec)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *SlurmApplication) DeepCopyInto(out *SlurmApplication) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ObjectMeta.DeepCopyInto(&out.ObjectMeta)
	in.Spec.DeepCopyInto(&out.Spec)
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new SlurmApplication.
func (in *SlurmApplication) DeepCopy() *SlurmApplication {
	if in == nil {
		return nil
	}
	out := new(SlurmApplication)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *SlurmApplication) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *SlurmApplicationList) DeepCopyInto(out *SlurmApplicationList) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ListMeta.DeepCopyInto(&out.ListMeta)
	if in.Items != nil {
		in, out := &in.Items, &out.Items
		*out = make([]SlurmApplication, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new SlurmApplicationList.
func (in *SlurmApplicationList) DeepCopy() *SlurmApplicationList {
	if in == nil {
		return nil
	}
	out := new(SlurmApplicationList)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *SlurmApplicationList) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *SlurmApplicationSpec) DeepCopyInto(out *SlurmApplicationSpec) {
	*out = *in
	if in.Volumes != nil {
		in, out := &in.Volumes, &out.Volumes
		*out = make([]v1.Volume, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
	in.Jupyter.DeepCopyInto(&out.Jupyter)
	in.Master.DeepCopyInto(&out.Master)
	in.Node.DeepCopyInto(&out.Node)
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new SlurmApplicationSpec.
func (in *SlurmApplicationSpec) DeepCopy() *SlurmApplicationSpec {
	if in == nil {
		return nil
	}
	out := new(SlurmApplicationSpec)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *SlurmJupyterSpec) DeepCopyInto(out *SlurmJupyterSpec) {
	*out = *in
	in.CommonSpec.DeepCopyInto(&out.CommonSpec)
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new SlurmJupyterSpec.
func (in *SlurmJupyterSpec) DeepCopy() *SlurmJupyterSpec {
	if in == nil {
		return nil
	}
	out := new(SlurmJupyterSpec)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *SlurmMasterSpec) DeepCopyInto(out *SlurmMasterSpec) {
	*out = *in
	in.CommonSpec.DeepCopyInto(&out.CommonSpec)
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new SlurmMasterSpec.
func (in *SlurmMasterSpec) DeepCopy() *SlurmMasterSpec {
	if in == nil {
		return nil
	}
	out := new(SlurmMasterSpec)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *SlurmNodeSpec) DeepCopyInto(out *SlurmNodeSpec) {
	*out = *in
	in.CommonSpec.DeepCopyInto(&out.CommonSpec)
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new SlurmNodeSpec.
func (in *SlurmNodeSpec) DeepCopy() *SlurmNodeSpec {
	if in == nil {
		return nil
	}
	out := new(SlurmNodeSpec)
	in.DeepCopyInto(out)
	return out
}
