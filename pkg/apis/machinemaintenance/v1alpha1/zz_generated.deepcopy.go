// +build !ignore_autogenerated

// Code generated by operator-sdk. DO NOT EDIT.

package v1alpha1

import (
	runtime "k8s.io/apimachinery/pkg/runtime"
)

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *MachineMaintenance) DeepCopyInto(out *MachineMaintenance) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ObjectMeta.DeepCopyInto(&out.ObjectMeta)
	out.Spec = in.Spec
	out.Status = in.Status
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new MachineMaintenance.
func (in *MachineMaintenance) DeepCopy() *MachineMaintenance {
	if in == nil {
		return nil
	}
	out := new(MachineMaintenance)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *MachineMaintenance) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *MachineMaintenanceList) DeepCopyInto(out *MachineMaintenanceList) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ListMeta.DeepCopyInto(&out.ListMeta)
	if in.Items != nil {
		in, out := &in.Items, &out.Items
		*out = make([]MachineMaintenance, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new MachineMaintenanceList.
func (in *MachineMaintenanceList) DeepCopy() *MachineMaintenanceList {
	if in == nil {
		return nil
	}
	out := new(MachineMaintenanceList)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *MachineMaintenanceList) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *MachineMaintenanceSpec) DeepCopyInto(out *MachineMaintenanceSpec) {
	*out = *in
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new MachineMaintenanceSpec.
func (in *MachineMaintenanceSpec) DeepCopy() *MachineMaintenanceSpec {
	if in == nil {
		return nil
	}
	out := new(MachineMaintenanceSpec)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *MachineMaintenanceStatus) DeepCopyInto(out *MachineMaintenanceStatus) {
	*out = *in
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new MachineMaintenanceStatus.
func (in *MachineMaintenanceStatus) DeepCopy() *MachineMaintenanceStatus {
	if in == nil {
		return nil
	}
	out := new(MachineMaintenanceStatus)
	in.DeepCopyInto(out)
	return out
}