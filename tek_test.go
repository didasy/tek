package tek_test

import (
	. "github.com/didasy/tek"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"io/ioutil"
	"testing"
)

var (
	err error
	indonesian []byte
	sample []byte
)

func init() {
	indonesian, err = ioutil.ReadFile("./indonesian.txt")
	if err != nil {
		Fail("Failed to read indonesian.txt")
	}
	sample, err = ioutil.ReadFile("./sample.txt")
	if err != nil {
		Fail("Failed to read sample.txt")
	}
}

var _ = Describe("Tek", func() {

	Describe("Testing the tagger", func() {
		Context("Set language", func() {
			It("Should not return error if language is not `en` or `id`", func() {
				err := SetLang("us")
				Expect(err).To(BeNil())
			})
			It("Should not return error if language is set to `id`", func() {
				err := SetLang("id")
				Expect(err).To(BeNil())
			})
			It("Should not return error if language is set to `en`", func() {
				err := SetLang("en")
				Expect(err).To(BeNil())
			})
		})
		Context("Get tags of indonesian.txt", func() {
			It("Should return a slice of *Info with length of 5", func() {
				var term string
				var tf, idf, tfidf float64

				SetLang("id")
				tags := GetTags(string(indonesian), 5)
				Expect(tags).To(HaveLen(5))
				Expect(tags[0].Term).ToNot(BeEmpty())
				Expect(tags[0].Term).To(BeAssignableToTypeOf(term))
				Expect(tags[0].Idf).ToNot(BeZero())
				Expect(tags[0].Idf).To(BeAssignableToTypeOf(idf))
				Expect(tags[0].Idf).ToNot(BeZero())
				Expect(tags[0].Tf).To(BeAssignableToTypeOf(tf))
				Expect(tags[0].Tf).ToNot(BeZero())
				Expect(tags[0].Tfidf).To(BeAssignableToTypeOf(tfidf))
			})
		})
		Context("Get tags of sample.txt", func() {
			It("Should return a slice of *Info with length of 5", func() {
				var term string
				var tf, idf, tfidf float64

				SetLang("id")
				tags := GetTags(string(sample), 5)
				Expect(tags).To(HaveLen(5))
				Expect(tags[0].Term).ToNot(BeEmpty())
				Expect(tags[0].Term).To(BeAssignableToTypeOf(term))
				Expect(tags[0].Idf).ToNot(BeZero())
				Expect(tags[0].Idf).To(BeAssignableToTypeOf(idf))
				Expect(tags[0].Idf).ToNot(BeZero())
				Expect(tags[0].Tf).To(BeAssignableToTypeOf(tf))
				Expect(tags[0].Tf).ToNot(BeZero())
				Expect(tags[0].Tfidf).To(BeAssignableToTypeOf(tfidf))
			})
		})
	})

})

func BenchmarkGetTagEn(b *testing.B) {
	SetLang("en")
	for i := 0; i < b.N; i++ {
		GetTags(string(sample), 10)
	}
}

func BenchmarkGetTagId(b *testing.B) {
	SetLang("id")
	for i := 0; i < b.N; i++ {
		GetTags(string(indonesian), 10)
	}
}
