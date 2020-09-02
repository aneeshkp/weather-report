package controllers

import (
	"context"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"math/rand"
	"strings"
	"time"

	weatherv1alpha1 "github.com/aneeshkp/weather-reprot/api/v1alpha1"
	"k8s.io/apimachinery/pkg/types"
)

const timeout = time.Second * 60
const interval = time.Second * 1

var (
	//name      string
	//namespace string
	//request   reconcile.Request
	ctx     = context.Background()
	fetched = &weatherv1alpha1.Cityweather{}
)

/*var _ = Describe("deployment Resource ", func() {

	dep := GetDeployment()
	Context("Test Deployment Resource", func() {
		It("Should create  deployment ", func() {
			Expect(k8sClient.Create(ctx, dep)).Should(Succeed())
		})
		It("Should have readyReplicas=1", func() {
			depInstance := &appsv1.Deployment{}
			By("Expecting deployment  created")
			Eventually(func() error {
				err := k8sClient.Get(ctx, types.NamespacedName{Name: dep.Name, Namespace: "default"}, depInstance)
				return err
			}, timeout, interval).ShouldNot(HaveOccurred())
			Expect(depInstance.Status.ReadyReplicas).To(Equal(1))

		})
	})

})*/
var _ = Describe("Cwcontroller", func() {

	BeforeEach(func() {
		/*name = "test-resource"
		namespace = "default"
		request = reconcile.Request{
			NamespacedName: types.NamespacedName{
				Name:      name,
				Namespace: namespace,
			},
		}*/
	})
	//var _ = Describe("")
	Describe("Weather Report Operator", func() {
		var (
			instance *weatherv1alpha1.Cityweather
		)
		Context("Initially  ", func() {
			It("Should create CR Successfully ", func() {
				instance = createInstanceInCluster()
				Expect(k8sClient.Create(ctx, instance)).Should(Succeed())
			})
			It("Should created CR exist", func() {
				By("Expecting kind created")
				Eventually(func() error {
					err := k8sClient.Get(ctx, types.NamespacedName{Name: instance.Name, Namespace: "default"}, fetched)
					return err
				}, timeout, interval).ShouldNot(HaveOccurred())
			})

			Context("When CR is Exists    ", func() {
				It("Should create Pod(s) ", func() {
					cityPod := &corev1.Pod{}
					Eventually(func() error {
						err := k8sClient.Get(ctx, types.NamespacedName{Name: instance.Name + "london", Namespace: "default"}, cityPod)
						return err
					}, timeout, interval).ShouldNot(HaveOccurred())
					Expect(cityPod.Name).To(Equal(instance.Name + "london"))
				})
			})
			Context("Once Pod exists", func() {
				It("Should update CR  status ", func() {
					Eventually(func() error {
						err := k8sClient.Get(ctx, types.NamespacedName{Name: instance.Name, Namespace: "default"}, fetched)
						return err
					}, timeout, interval).ShouldNot(HaveOccurred())
					//CHECK CR STATUS
					Expect(len(fetched.Status.City)).To(Equal(2))
				})
			})
			Context("Update CR to remove london and add chicago and austin", func() {
				It("Should update new cr", func() {
					fetched.Spec.City = []string{"newyork", "chicago", "austin"}
					Eventually(func() error {
						err := k8sClient.Update(ctx, fetched)
						return err
					}, timeout, interval).ShouldNot(HaveOccurred())
				})
				It("Should have new city in the spec  ", func() {
					Eventually(func() error {
						err := k8sClient.Get(ctx, types.NamespacedName{Name: instance.Name, Namespace: "default"}, fetched)
						return err
					}, timeout, interval).ShouldNot(HaveOccurred())
					Expect(len(fetched.Spec.City)).To(Equal(3))
				})
				//CHECK CR STATUS
			})
			Context("Once CR is updated", func() {
				It("should delete london pod", func() {
					londonPod := &corev1.Pod{}
					Eventually(func() error {
						err := k8sClient.Get(ctx, types.NamespacedName{Name: instance.Name + "london", Namespace: "default"}, londonPod)
						return err
					}, timeout, interval).ShouldNot(HaveOccurred())
					Expect(londonPod).ToNot(BeNil())
				})
				It("should create new austin pod", func() {
					austinPod := &corev1.Pod{}
					Eventually(func() error {
						err := k8sClient.Get(ctx, types.NamespacedName{Name: instance.Name + "austin", Namespace: "default"}, austinPod)
						return err
					}, timeout, interval).ShouldNot(HaveOccurred())
					Expect(austinPod).ToNot(BeNil())
				})
				It("Should update CR status with new pods ", func() {
					Eventually(func() error {
						err := k8sClient.Get(ctx, types.NamespacedName{Name: instance.Name, Namespace: "default"}, fetched)
						return err
					}, timeout, interval).ShouldNot(HaveOccurred())
					//CHECK CR STATUS
					Expect(len(fetched.Status.City)).To(Equal(3))
					Expect(fetched.Status.City).Should((HaveKey("austin")))
					Expect(fetched.Status.City).Should((HaveKey("chicago")))
					Expect(fetched.Status.City).Should((HaveKey("newyork")))
				})
			})

			/*BeforeEach(func() {
				// Create a new resource using k8sClient.Create()
				// I'm just going to assume you've done this in
				// a method called createInstanceInCluster()
				instance = createInstanceInCluster()
				Expect(k8sClient.Create(ctx,instance)).Should(Succeed())

				Eventually(func() bool{
					err:=k8sClient.Get(ctx,types.NamespacedName{Name:instance.Name,Namespace:request.Namespace},fetched)
					return err==nil
				},timeout,interval).Should(BeTrue())
			})
			AfterEach(func() {
				// Remember to clean up the cluster after each test
				//deleteInstanceInCluster(instance)

			})
			It("should update the status of the CR", func() {
				// Some method where you've implemented Gomega Expect()'s
				//assertResourceUpdated(instance)
				Expect(fetched.Name).To(Equal("weather-report"))
				//Expect(len(instance.Status.City)).To(Equal(2))
			}) */
		})
	})
})

func createInstanceInCluster() *weatherv1alpha1.Cityweather {
	spec := weatherv1alpha1.CityweatherSpec{
		City: []string{"newyork", "london"},
		Days: 1,
	}

	toCreate := weatherv1alpha1.Cityweather{
		TypeMeta: v1.TypeMeta{},
		ObjectMeta: v1.ObjectMeta{
			Name:      "weather-report", //-" + randString(),
			Namespace: "default",
		},
		Spec: spec,
	}
	return &toCreate
}

func randString() string {
	rand.Seed(time.Now().UnixNano())
	chars := []rune("abcdefghijklmnopqrstuvwxyz")
	length := 6
	var b strings.Builder
	for i := 0; i < length; i++ {
		b.WriteRune(chars[rand.Intn(len(chars))])
	}
	return b.String()
}

// Define the desired Deployment object1
func GetDeployment() *appsv1.Deployment {
	deploy1 := &appsv1.Deployment{
		ObjectMeta: v1.ObjectMeta{
			Name:      "deployment1",
			Namespace: "default",
		},
		Spec: appsv1.DeploymentSpec{
			Selector: &v1.LabelSelector{
				MatchLabels: map[string]string{"deployment": "deployment1-deployment"},
			},
			Template: corev1.PodTemplateSpec{
				ObjectMeta: v1.ObjectMeta{Labels: map[string]string{"deployment": "deployment1-deployment"}},
				Spec: corev1.PodSpec{
					Containers: []corev1.Container{
						{
							Name:  "nginx",
							Image: "nginx",
						},
					},
				},
			},
		},
	}
	return deploy1
}
