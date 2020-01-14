package v1

import (
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

// RedisSpec defines the desired state of Redis
type RedisSpec struct {
	// INSERT ADDITIONAL SPEC FIELDS - desired state of cluster
	// Important: Run "operator-sdk generate k8s" to regenerate code after modifying this file
	// Add custom validation using kubebuilder tags: https://book-v1.book.kubebuilder.io/beyond_basics/generating_crd.html
	Replicas          *int32                      `json:"replicas"`
	RedisImage        string                      `json:"redisImage"`
	ClusterType       string                      `json:"clusterType"`
	RedisResources    corev1.ResourceRequirements `json:"redisResources,omitempty"`
	SentinelResources corev1.ResourceRequirements `json:"sentinelResources,omitempty"`
	MasterRouter      RedisRouter                 `json:"masterRouter,omitempty"`
	SlaveRouter       RedisRouter                 `json:"slaveRouter,omitempty"`
	SentinelRouter    RedisRouter                 `json:"sentinelRouter,omitempty"`
}

// cluster network router
type RedisRouter struct {
	Key  string `json:"key"`
	Port int32  `json:"port"`
}

// RedisStatus defines the observed state of Redis
type RedisStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "operator-sdk generate k8s" to regenerate code after modifying this file
	// Add custom validation using kubebuilder tags: https://book-v1.book.kubebuilder.io/beyond_basics/generating_crd.html
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// Redis is the Schema for the redis API
// +kubebuilder:subresource:status
// +kubebuilder:resource:path=redis,scope=Namespaced
type Redis struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   RedisSpec   `json:"spec,omitempty"`
	Status RedisStatus `json:"status,omitempty"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// RedisList contains a list of Redis
type RedisList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []Redis `json:"items"`
}

func init() {
	SchemeBuilder.Register(&Redis{}, &RedisList{})
}
