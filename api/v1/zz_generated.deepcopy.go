//go:build !ignore_autogenerated

/*
Copyright 2024.

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

package v1

import (
	runtime "k8s.io/apimachinery/pkg/runtime"
)

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *Resource) DeepCopyInto(out *Resource) {
	*out = *in
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new Resource.
func (in *Resource) DeepCopy() *Resource {
	if in == nil {
		return nil
	}
	out := new(Resource)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *ResourceCleaner) DeepCopyInto(out *ResourceCleaner) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ObjectMeta.DeepCopyInto(&out.ObjectMeta)
	in.Spec.DeepCopyInto(&out.Spec)
	out.Status = in.Status
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new ResourceCleaner.
func (in *ResourceCleaner) DeepCopy() *ResourceCleaner {
	if in == nil {
		return nil
	}
	out := new(ResourceCleaner)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *ResourceCleaner) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *ResourceCleanerList) DeepCopyInto(out *ResourceCleanerList) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ListMeta.DeepCopyInto(&out.ListMeta)
	if in.Items != nil {
		in, out := &in.Items, &out.Items
		*out = make([]ResourceCleaner, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new ResourceCleanerList.
func (in *ResourceCleanerList) DeepCopy() *ResourceCleanerList {
	if in == nil {
		return nil
	}
	out := new(ResourceCleanerList)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *ResourceCleanerList) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *ResourceCleanerSpec) DeepCopyInto(out *ResourceCleanerSpec) {
	*out = *in
	in.Resources.DeepCopyInto(&out.Resources)
	in.Expire.DeepCopyInto(&out.Expire)
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new ResourceCleanerSpec.
func (in *ResourceCleanerSpec) DeepCopy() *ResourceCleanerSpec {
	if in == nil {
		return nil
	}
	out := new(ResourceCleanerSpec)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *ResourceCleanerStatus) DeepCopyInto(out *ResourceCleanerStatus) {
	*out = *in
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new ResourceCleanerStatus.
func (in *ResourceCleanerStatus) DeepCopy() *ResourceCleanerStatus {
	if in == nil {
		return nil
	}
	out := new(ResourceCleanerStatus)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *ResourcesSpec) DeepCopyInto(out *ResourcesSpec) {
	*out = *in
	if in.Include != nil {
		in, out := &in.Include, &out.Include
		*out = make([]Resource, len(*in))
		copy(*out, *in)
	}
	if in.Exclude != nil {
		in, out := &in.Exclude, &out.Exclude
		*out = make([]Resource, len(*in))
		copy(*out, *in)
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new ResourcesSpec.
func (in *ResourcesSpec) DeepCopy() *ResourcesSpec {
	if in == nil {
		return nil
	}
	out := new(ResourcesSpec)
	in.DeepCopyInto(out)
	return out
}
