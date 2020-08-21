package controllers

import (
	"context"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
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

var _ = Describe("Cwcontroller", func() {
	var (
		//name      string
		//namespace string
		//request   reconcile.Request
		ctx     = context.Background()
		fetched = &weatherv1alpha1.Cityweather{}
	)
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
	Describe("ResourceAccess CRD", func() {
		var (
			instance *weatherv1alpha1.Cityweather
		)
		Context("With one Resource ", func() {
			It("should update the status of the CR", func() {
				instance = createInstanceInCluster()
				Expect(k8sClient.Create(ctx, instance)).Should(Succeed())
				//time.Sleep(time.Second * 5)

				//check CRD was created okay
				By("Expecting submitted")
				Eventually(func() error {
					err := k8sClient.Get(ctx, types.NamespacedName{Name: instance.Name, Namespace: "default"}, fetched)
					return err
				}, timeout, interval).ShouldNot(HaveOccurred())

				/*pod1 := CreatePod(instance, instance.Spec.City[0])
				By("Expecting submitted")
				Eventually(func() error {
					err := k8sClient.Create(ctx, pod1)
					return err
				}, timeout, interval).ShouldNot(HaveOccurred())*/

				//check pod is created
				pod2 := &corev1.Pod{}
				Eventually(func() error {
					err := k8sClient.Get(ctx, types.NamespacedName{Name: instance.Name + "london", Namespace: "default"}, pod2)
					return err
				}, timeout, interval).ShouldNot(HaveOccurred())

				Expect(fetched.Name).To(Equal("weather-report"))
				Expect(len(fetched.Status.City)).To(Equal(2))
				//Expect(fetched.Spec.City).To(Equal(""))
				//Expect(fetched.Status.City).To(Equal(""))

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
