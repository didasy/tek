package main

import (
	"fmt"
	"io/ioutil"
	"github.com/JesusIslam/tek"
)

func main() {
	Tb, _ := ioutil.ReadFile("../sample.txt")
	text := string(Tb)
	tek.SetLang("en")
	tags := tek.GetTags(text, 10)
	for _, tag := range tags {
		fmt.Println(tag.Term, tag.Tfidf)
	}
}