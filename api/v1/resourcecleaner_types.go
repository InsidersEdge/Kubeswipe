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

package v1

import (
	"time"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

type MyTime struct {
	time.Time
}

// ResourceCleanerSpec defines the desired state of ResourceCleaner
type ResourceCleanerSpec struct {
	// For example, "* * * * *" represents a schedule that runs every minute.
	Schedule      string        `json:"schedule"`
	Resources     ResourcesSpec `json:"resources,omitempty"`
	CloudProvider CloudName     `json:"cloudProvider,omitempty"`
	Expire        metav1.Time   `json:"expire,omitempty"`
}

type CloudName string

const (
	AWS   CloudName = "AWS"
	GCP   CloudName = "GCP"
	Azure CloudName = "Azure"
)

type ResourcesSpec struct {
	Include []Resource `json:"include,omitempty"`
	Exclude []Resource `json:"exclude,omitempty"`
}

type Resource struct {
	Name      string `json:"name"`
	Namespace string `json:"namespace,omitempty"`
	Backup    bool   `json:"backup,omitempty"`
}

type ResourceNames string

// ResourceCleanerStatus defines the observed state of ResourceCleaner
type ResourceCleanerStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "make" to regenerate code after modifying this file
}

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status

// ResourceCleaner is the Schema for the resourcecleaners API
type ResourceCleaner struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   ResourceCleanerSpec   `json:"spec,omitempty"`
	Status ResourceCleanerStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

// ResourceCleanerList contains a list of ResourceCleaner
type ResourceCleanerList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []ResourceCleaner `json:"items"`
}

func init() {
	SchemeBuilder.Register(&ResourceCleaner{}, &ResourceCleanerList{})
}
