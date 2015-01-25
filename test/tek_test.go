package main

import (
	"testing"
	"io/ioutil"
	"github.com/JesusIslam/tek"
)

const (
	num = 10
	lang = "en"
)

var result []*tek.Info

func TestGetTag(t *testing.T) {
	textB, _ := ioutil.ReadFile("../sample.txt")
	text := string(textB)
	tek.SetLang("en")
	result = tek.GetTags(text, num)
}

func BenchmarkGetTag(b *testing.B) {
	textB, _ := ioutil.ReadFile("../sample.txt")
	text := string(textB)
	var r []*tek.Info
	tek.SetLang("en")
	for n := 0; n < b.N; n++ {
		r = tek.GetTags(text, num)
	}
	result = r
}