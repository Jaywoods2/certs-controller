package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.
type TlsData struct {
	Name string `json:"name,omitempty"`
	Key  string `json:"key,omitempty"`
	Crt  string `json:"crt,omitempty"`
}

// CertSecretSpec defines the desired state of CertSecret
// +k8s:openapi-gen=true
type CertSecretSpec struct {
	// INSERT ADDITIONAL SPEC FIELDS - desired state of cluster
	// Important: Run "operator-sdk generate k8s" to regenerate code after modifying this file
	// Add custom validation using kubebuilder tags: https://book-v1.book.kubebuilder.io/beyond_basics/generating_crd.html
	Tls []TlsData `json:"tls,omitempty"`
}

// CertSecretStatus defines the observed state of CertSecret
// +k8s:openapi-gen=true
type CertSecretStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "operator-sdk generate k8s" to regenerate code after modifying this file
	// Add custom validation using kubebuilder tags: https://book-v1.book.kubebuilder.io/beyond_basics/generating_crd.html
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// CertSecret is the Schema for the certsecrets API
// +k8s:openapi-gen=true
// +kubebuilder:subresource:status
type CertSecret struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   CertSecretSpec   `json:"spec,omitempty"`
	Status CertSecretStatus `json:"status,omitempty"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// CertSecretList contains a list of CertSecret
type CertSecretList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []CertSecret `json:"items"`
}

func init() {
	SchemeBuilder.Register(&CertSecret{}, &CertSecretList{})
}
