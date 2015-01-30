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

func TestGetTagEn(t *testing.T) {
	textB, _ := ioutil.ReadFile("../sample.txt")
	text := string(textB)
	tek.SetLang("en")
	result = tek.GetTags(text, num)
}

func TestGetTagId(t *testing.T) {
	textB, _ := ioutil.ReadFile("../indonesian.txt")
	text := string(textB)
	tek.SetLang("id")
	result = tek.GetTags(text, num)
}

func BenchmarkGetTagEn(b *testing.B) {
	textB, _ := ioutil.ReadFile("../sample.txt")
	text := string(textB)
	var r []*tek.Info
	tek.SetLang("en")
	for n := 0; n < b.N; n++ {
		r = tek.GetTags(text, num)
	}
	result = r
}

func BenchmarkGetTagId(b *testing.B) {
	textB, _ := ioutil.ReadFile("../indonesian.txt")
	text := string(textB)
	var r []*tek.Info
	tek.SetLang("id")
	for n := 0; n < b.N; n++ {
		r = tek.GetTags(text, num)
	}
	result = r
}