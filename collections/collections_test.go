package collections_test

import (
	"errors"
	"fmt"
	"strings"
	"testing"

	"github.com/paketo-buildpacks/libdependency/collections"
	"github.com/sclevine/spec"

	. "github.com/onsi/gomega"
)

func testCollections(t *testing.T, context spec.G, it spec.S) {
	Expect := NewWithT(t).Expect

	context("FilterFunc", func() {
		it("will filter", func() {
			filtered := collections.FilterFunc([]string{"a", "aa", "b", "bb", "bab"}, func(s string) bool {
				return strings.Contains(s, "a")
			})
			Expect(filtered).To(ConsistOf("a", "aa", "bab"))
		})

		it("will gracefully handle an empty input", func() {
			filtered := collections.FilterFunc([]float64{}, func(_ float64) bool {
				return true
			})
			Expect(filtered).To(BeEmpty())
		})

		it("will gracefully handle a nil input", func() {
			filtered := collections.FilterFunc([]spec.Spec(nil), func(s spec.Spec) bool {
				return true
			})
			Expect(filtered).To(BeEmpty())
		})

		it("will gracefully handle a nil func", func() {
			filtered := collections.FilterFunc([]string{"a", "aa", "b", "bb", "bab"}, nil)
			Expect(filtered).To(BeEmpty())
		})
	})

	context("TransformFunc", func() {
		it("will transform", func() {
			transformed := collections.TransformFunc([]int{1, 2, 3, 4, 5}, func(i int) int {
				return i * 2
			})
			Expect(transformed).To(ConsistOf(2, 4, 6, 8, 10))
		})

		it("will transform types", func() {
			transformed := collections.TransformFunc([]int{1, 2, 3, 4, 5}, func(i int) string {
				return fmt.Sprintf("%d", i*2)
			})
			Expect(transformed).To(ConsistOf("2", "4", "6", "8", "10"))
		})

		it("will gracefully handle an empty input", func() {
			transformed := collections.TransformFunc([]float64{}, func(_ float64) bool {
				return true
			})
			Expect(transformed).To(BeEmpty())
		})

		it("will gracefully handle a nil input", func() {
			transformed := collections.TransformFunc(nil, func(s spec.Spec) bool {
				return true
			})
			Expect(transformed).To(BeEmpty())
		})

		it("will gracefully handle a nil func", func() {
			transformed := collections.TransformFunc[int, int]([]int{1, 2, 3, 4, 5}, nil)
			Expect(transformed).To(BeEmpty())
		})
	})

	context("TransformFuncWithError", func() {
		it("will transform", func() {
			transformed, err := collections.TransformFuncWithError([]int{1, 2, 3, 4, 5}, func(i int) (string, error) {
				return fmt.Sprintf("%d", i), nil
			})
			Expect(err).NotTo(HaveOccurred())
			Expect(transformed).To(ConsistOf("1", "2", "3", "4", "5"))
		})

		it("will gracefully handle an empty input", func() {
			transformed, err := collections.TransformFuncWithError([]float64{}, func(_ float64) (string, error) {
				return "", nil
			})
			Expect(err).NotTo(HaveOccurred())
			Expect(transformed).To(BeEmpty())
		})

		it("will gracefully handle a nil input", func() {
			transformed, err := collections.TransformFuncWithError(nil, func(s spec.Spec) (string, error) {
				return "", nil
			})
			Expect(err).NotTo(HaveOccurred())
			Expect(transformed).To(BeEmpty())
		})

		it("will gracefully handle a nil func", func() {
			transformed, err := collections.TransformFuncWithError[int, int]([]int{1, 2, 3, 4, 5}, nil)
			Expect(err).NotTo(HaveOccurred())
			Expect(transformed).To(BeEmpty())
		})

		context("failure cases", func() {
			it("will return the error", func() {
				_, err := collections.TransformFuncWithError([]int{1}, func(i int) (string, error) {
					return "", errors.New("error from inside transform with error")
				})

				Expect(err).To(MatchError("error from inside transform with error"))
			})
		})
	})
}
