package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

// MachineMaintenanceSpec defines the desired state of MachineMaintenance
type MachineMaintenanceSpec struct {
	// INSERT ADDITIONAL SPEC FIELDS - desired state of cluster
	// Important: Run "operator-sdk generate k8s" to regenerate code after modifying this file
	// Add custom validation using kubebuilder tags: https://book-v1.book.kubebuilder.io/beyond_basics/generating_crd.html
	Maintenance bool   `json:"maintenance,required"`
	EventCode   string `json:"eventcode,omitempty"`
	EventID     string `json:"eventid,omitempty"`
	//	NotBefore   time.Time `json:"notbefore,omitempty"`
	MachineID string `json:"machineid,required"`
}

// MachineMaintenanceStatus defines the observed state of MachineMaintenance
type MachineMaintenanceStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "operator-sdk generate k8s" to regenerate code after modifying this file
	// Add custom validation using kubebuilder tags: https://book-v1.book.kubebuilder.io/beyond_basics/generating_crd.html
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// MachineMaintenance is the Schema for the machinemaintenances API
// +kubebuilder:subresource:status
// +kubebuilder:resource:path=machinemaintenances,scope=Namespaced
type MachineMaintenance struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   MachineMaintenanceSpec   `json:"spec,omitempty"`
	Status MachineMaintenanceStatus `json:"status,omitempty"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// MachineMaintenanceList contains a list of MachineMaintenance
type MachineMaintenanceList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []MachineMaintenance `json:"items"`
}

func init() {
	SchemeBuilder.Register(&MachineMaintenance{}, &MachineMaintenanceList{})
}
