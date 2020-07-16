/*


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

package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

// DeployassistSpec defines the desired state of Deployassist
type DeployassistSpec struct {
	// INSERT ADDITIONAL SPEC FIELDS - desired state of cluster
	// Important: Run "make" to regenerate code after modifying this file

	// Foo is an example field of Deployassist. Edit Deployassist_types.go to remove/update
	Files []string `json:"files,omitempty"`
}

// DeployassistStatus defines the observed state of Deployassist
type DeployassistStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "make" to regenerate code after modifying this file
	AppliedFiles []string `json:"files,omitempty"`
}

// +kubebuilder:object:root=true

// Deployassist is the Schema for the deployassists API
type Deployassist struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   DeployassistSpec   `json:"spec,omitempty"`
	Status DeployassistStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true

// DeployassistList contains a list of Deployassist
type DeployassistList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []Deployassist `json:"items"`
}

func init() {
	SchemeBuilder.Register(&Deployassist{}, &DeployassistList{})
}
