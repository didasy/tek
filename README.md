# tek
### A golang package to get tags of an article
----------------------------------------------
[![GoDoc](https://godoc.org/github.com/didasy/tek?status.svg)](https://godoc.org/github.com/didasy/tek)
[![Coverage Status](https://coveralls.io/repos/github/didasy/tek/badge.svg?branch=feat/perf%0Amaster)](https://coveralls.io/github/didasy/tek?branch=feat/perf%0Amaster)

### Installation
`go get github.com/didasy/tek`

### Dependencies
None

### Usage
```
package main

import (
	"fmt"
	"io/ioutil"
	"github.com/didasy/tek"
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
```
### Testing
To test, just run `go test`, but you need to have [gomega](http://github.com/onsi/gomega) and [ginkgo](http://github.com/onsi/ginkgo) installed.

### Benchmark
Using i3-3217U @1.8GHz with 370 total words from the sample.txt provided and command `go test -bench . -benchtime=5s -cpu 4`:
```
BenchmarkGetTagEn-4         1000          11097404 ns/op
BenchmarkGetTagId-4         1000           6295201 ns/op
```

### License
See LICENSE file, it is MIT