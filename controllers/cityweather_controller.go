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

package controllers

import (
	"context"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
	//"k8s.io/apimachinery/pkg/types"

	"sigs.k8s.io/controller-runtime/pkg/reconcile"
	//"strconv"

	"fmt"
	"github.com/go-logr/logr"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"

	weatherv1alpha1 "github.com/aneeshkp/weather-reprot/api/v1alpha1"
)

// CityweatherReconciler reconciles a Cityweather object
type CityweatherReconciler struct {
	client.Client
	Log    logr.Logger
	Scheme *runtime.Scheme
}

// +kubebuilder:rbac:groups=weather.aputtur.com,resources=cityweathers,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=weather.aputtur.com,resources=cityweathers/status,verbs=get;update;patch

func (r *CityweatherReconciler) Reconcile(req ctrl.Request) (ctrl.Result, error) {
	_ = context.Background()
	_ = r.Log.WithValues("cityweather", req.NamespacedName)

	// your logic here
	// take the cr valus and check if pod exists wioth that value if not create
	// adn delete the pods that are not with that
	//newpod, _ := c.kubeClient.CoreV1().Pods("default").Create(pod)

	weather := &weatherv1alpha1.Cityweather{}

	err := r.Client.Get(context.TODO(), req.NamespacedName, weather)
	if err != nil {
		if errors.IsNotFound(err) {
			// Request object not found, could have been deleted after reconcile request.
			// Owned objects are automatically garbage collected. For additional cleanup logic use finalizers.
			// Return and don't requeue
			return reconcile.Result{}, nil
		}
		// Error reading the object - requeue the request.
		return reconcile.Result{}, err
	}
	//get all pods with labesl
	// loop thought and see if names are mathcinhg
	podList := &corev1.PodList{}
	lbs := map[string]string{
		"app":     "weather-report",
		"version": "v0.1",
	}

	labelSelector := labels.SelectorFromSet(lbs)
	lisOpts := &client.ListOptions{Namespace: req.Namespace, LabelSelector: labelSelector}
	if err = r.Client.List(context.TODO(), podList, lisOpts); err != nil {
		r.Log.Error(err, "Unable to retrieve pod list %v")
		return reconcile.Result{}, err
	}
	r.Log.Info(req.Namespace)
	r.Log.Info("pod list")
	for _, pod := range podList.Items {
		r.Log.Info(pod.Name)
	}
	// these are available pods right now
	var available map[string]corev1.Pod
	available = make(map[string]corev1.Pod)

	var availablecity map[string]string
	availablecity = make(map[string]string)
	for _, pod := range podList.Items {
		if pod.ObjectMeta.DeletionTimestamp != nil {
			continue
		}
		if pod.Status.Phase == corev1.PodRunning || pod.Status.Phase == corev1.PodPending {
			available[pod.ObjectMeta.Name] = pod
			availablecity[pod.Labels["city"]] = pod.ObjectMeta.Name
		}
	}

	// status doesnt matter now
	//Clean the pods that are not in the spec
	podToDelete := &corev1.Pod{}
	for city, podname := range availablecity {
		if !itemExists(weather.Spec.City, city) {
			key := client.ObjectKey{Namespace: req.Namespace, Name: podname}
			err := r.Client.Get(context.TODO(), key, podToDelete)
			if err != nil && !errors.IsNotFound(err) {
				r.Log.Error(err, "Unable to retrieve pod %v")
			} else {
				r.Client.Delete(context.TODO(), podToDelete)
			}
			return reconcile.Result{Requeue: true}, nil
		}
	}

	for _, city := range weather.Spec.City {
		r.Log.Info(city)
		if _, ok := availablecity[city]; !ok {
			r.Client.Create(context.TODO(), CreatePod(weather, city))
			r.Log.Info("pod creating")
			return reconcile.Result{Requeue: true}, nil
		}
	}

	//update status
	weather.Status.City = nil
	var statusCity map[string]string
	statusCity = make(map[string]string)
	if err = r.Client.List(context.TODO(), podList, lisOpts); err != nil {
		return reconcile.Result{}, err
	}
	for _, pod := range podList.Items {
		fmt.Printf("Pod labels %#v", pod.Labels)
		if pod.ObjectMeta.DeletionTimestamp != nil {
			continue
		}
		if pod.Status.Phase == corev1.PodRunning || pod.Status.Phase == corev1.PodPending {
			statusCity[pod.Labels["city"]] = pod.ObjectMeta.Name
		}
	}
	weather.Status.City = statusCity
	r.Client.Status().Update(context.Background(), weather)

	return ctrl.Result{}, nil
}

func (r *CityweatherReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&weatherv1alpha1.Cityweather{}).
		Complete(r)
}

func CreatePod(instance *weatherv1alpha1.Cityweather, city string) (pod *corev1.Pod) {
	url := fmt.Sprintf("http://wttr.in/%s?%d", city, 1)
	labels := map[string]string{
		"app":     "weather-report",
		"city":    city,
		"days":    "1",
		"version": "v0.1",
	}
	pod = &corev1.Pod{
		TypeMeta: metav1.TypeMeta{
			Kind:       "Pod",
			APIVersion: "v1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      instance.ObjectMeta.Name + city, //strconv.Itoa(time.Now().Nanosecond()),
			Namespace: "default",
			Labels:    labels,
		},
		Spec: corev1.PodSpec{
			Containers: []corev1.Container{
				{
					Name:    instance.ObjectMeta.Name + city,
					Image:   "tutum/curl",
					Command: []string{"sh", "-c", "curl -s " + url + " && sleep 3600"},
				},
			},
		},
	}
	return
}

func itemExists(arr []string, item string) bool {
	for _, city := range arr {
		if city == item {
			return true
		}
	}
	return false
}
