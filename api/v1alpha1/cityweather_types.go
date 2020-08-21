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

// CityweatherSpec defines the desired state of Cityweather
type CityweatherSpec struct {
	// INSERT ADDITIONAL SPEC FIELDS - desired state of cluster
	// Important: Run "make" to regenerate code after modifying this file

	// Foo is an example field of Cityweather. Edit Cityweather_types.go to remove/update
	City []string `json:"city,omitempty"`
	Days int      `json:"days"`
}

// CityweatherStatus defines the observed state of Cityweather
type CityweatherStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "make" to regenerate code after modifying this file
	//State string `json:"state"`
	City map[string]string `json:"city,omitempty"`
}

// +kubebuilder:object:root=true
// +kubebuilder:subresource:status

// Cityweather is the Schema for the cityweathers API
type Cityweather struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   CityweatherSpec   `json:"spec,omitempty"`
	Status CityweatherStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true

// CityweatherList contains a list of Cityweather
type CityweatherList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []Cityweather `json:"items"`
}

func init() {
	SchemeBuilder.Register(&Cityweather{}, &CityweatherList{})
}
