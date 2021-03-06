package utils_test

import (
	"os"
	"regexp"

	"github.com/Azure/application-gateway-kubernetes-ingress/pkg/utils"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Utils", func() {

	Describe("Testing `UnorderedSets`", func() {
		var testSet utils.UnorderedSet
		BeforeEach(func() {
			testSet = utils.NewUnorderedSet()
			testSet.Insert(1)
			testSet.Insert(2)
			testSet.Insert(3)
			testSet.Insert(4)
			testSet.Insert(4)
		})

		Context("Inserting non-unique elements", func() {
			It("Should only store unique elements", func() {
				Expect(testSet.Size()).To(Equal(4))
			})
		})

		Context("Erasing an element", func() {
			It("Should remove the element", func() {
				testSet.Erase(4)
				Expect(testSet.Contains(4)).To(BeFalse())
				Expect(testSet.Size()).To(Equal(3))
			})
		})

		Context("Clearing the unordered set", func() {
			It("Should erase all elements", func() {
				testSet.Clear()
				Expect(testSet.Size()).To(Equal(0))
				Expect(testSet.IsEmpty()).To(BeTrue())
			})
		})

		Context("Testing union of sets", func() {
			set := utils.NewUnorderedSet()
			set.Insert(5)
			set.Insert(6)
			set.Insert(7)

			It("Should contain union of both sets", func() {
				uSet := testSet.Union(set)
				Expect(uSet.Size()).To(Equal(7))
				Expect(uSet.Contains(5)).To(BeTrue())
				Expect(uSet.Contains(6)).To(BeTrue())
				Expect(uSet.Contains(7)).To(BeTrue())
			})
		})

		Context("Testing intersection of sets", func() {
			set := utils.NewUnorderedSet()
			set.Insert(3)
			set.Insert(4)

			It("Should only contain intersection of both sets", func() {
				iSet := testSet.Intersect(set)
				Expect(iSet.Size()).To(Equal(2))
				Expect(iSet.Contains(3)).To(BeTrue())
				Expect(iSet.Contains(4)).To(BeTrue())
				Expect(iSet.Contains(1)).To(BeFalse())
				Expect(iSet.Contains(2)).To(BeFalse())
			})
		})

	})

	Describe("Testing `utils` helpers", func() {
		Context("Testing integer comparators", func() {
			It("Should return maximum of two 64-bit integers", func() {
				Expect(utils.MaxInt64(int64(101), int64(100))).To(Equal(int64(101)))
				Expect(utils.MaxInt64(int64(100), int64(101))).To(Equal(int64(101)))
			})

			It("Should return maximum of two 32-bit integers", func() {
				Expect(utils.MaxInt32(int32(101), int32(100))).To(Equal(int32(101)))
				Expect(utils.MaxInt32(int32(100), int32(101))).To(Equal(int32(101)))
			})
		})

		Context("Testing string helpers", func() {
			It("Should return a string, which is a formatted list of integers", func() {
				Expect(utils.IntsToString([]int{1, 2, 3, 4, 5, 6}, ";")).To(Equal("1;2;3;4;5;6"))
			})
		})

		Context("Testing the Kubernetes namespace generator", func() {
			It("Given a namespace and resource it should return the Kubernetes resource identifier.", func() {
				Expect(utils.GetResourceKey("default", "pod")).To(Equal("default/pod"))
			})
		})

		Context("Testing the GetEnv helper", func() {
			const (
				expectedEnvVarValue = "expected-value"
				envVar              = "---some--environment--variable--with-low-likelihood-that-will-collide---"
			)
			BeforeEach(func() {
				// Make sure the environment variable we are using for this test does not already exist in the OS.
				_, exists := os.LookupEnv(envVar)
				Expect(exists).To(BeFalse())
				// Set it
				_ = os.Setenv(envVar, expectedEnvVarValue)
				_, exists = os.LookupEnv(envVar)
				Expect(exists).To(BeTrue())
			})
			AfterEach(func() {
				// Clean up the env var after the tests are done
				_ = os.Unsetenv(envVar)
			})
			It("returns default value in absence of an env var", func() {
				Expect(utils.GetEnv("-non-existent-key-we-hope", "expected-value", nil)).To(Equal("expected-value"))
			})
			It("returns expected value", func() {
				defaultValue := "--default--value--"
				passingValidator := regexp.MustCompile(`^[a-zA-Z\-]+$`)
				Expect(utils.GetEnv(envVar, defaultValue, passingValidator)).To(Equal("expected-value"))
			})
			It("returns default value in absence of an env var", func() {
				defaultValue := "--default--value--"
				// without a validator we get the environment variable's value
				Expect(utils.GetEnv(envVar, defaultValue, nil)).To(Equal(expectedEnvVarValue))

				// with a non-passing validator we get the default value
				failingValidator := regexp.MustCompile(`^[0-9]+$`)
				Expect(utils.GetEnv(envVar, defaultValue, failingValidator)).To(Equal(defaultValue))
			})
		})
	})
})
