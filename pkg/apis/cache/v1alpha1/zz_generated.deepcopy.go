// +build !ignore_autogenerated

// Code generated by operator-sdk. DO NOT EDIT.

package v1alpha1

import (
	runtime "k8s.io/apimachinery/pkg/runtime"
)

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *Ikshvaku) DeepCopyInto(out *Ikshvaku) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ObjectMeta.DeepCopyInto(&out.ObjectMeta)
	out.Spec = in.Spec
	in.Status.DeepCopyInto(&out.Status)
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new Ikshvaku.
func (in *Ikshvaku) DeepCopy() *Ikshvaku {
	if in == nil {
		return nil
	}
	out := new(Ikshvaku)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *Ikshvaku) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *IkshvakuList) DeepCopyInto(out *IkshvakuList) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ListMeta.DeepCopyInto(&out.ListMeta)
	if in.Items != nil {
		in, out := &in.Items, &out.Items
		*out = make([]Ikshvaku, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new IkshvakuList.
func (in *IkshvakuList) DeepCopy() *IkshvakuList {
	if in == nil {
		return nil
	}
	out := new(IkshvakuList)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *IkshvakuList) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *IkshvakuSpec) DeepCopyInto(out *IkshvakuSpec) {
	*out = *in
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new IkshvakuSpec.
func (in *IkshvakuSpec) DeepCopy() *IkshvakuSpec {
	if in == nil {
		return nil
	}
	out := new(IkshvakuSpec)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *IkshvakuStatus) DeepCopyInto(out *IkshvakuStatus) {
	*out = *in
	if in.Masters != nil {
		in, out := &in.Masters, &out.Masters
		*out = make([]string, len(*in))
		copy(*out, *in)
	}
	if in.Slaves != nil {
		in, out := &in.Slaves, &out.Slaves
		*out = make([]string, len(*in))
		copy(*out, *in)
	}
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new IkshvakuStatus.
func (in *IkshvakuStatus) DeepCopy() *IkshvakuStatus {
	if in == nil {
		return nil
	}
	out := new(IkshvakuStatus)
	in.DeepCopyInto(out)
	return out
}
