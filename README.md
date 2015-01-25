# tek
### A golang package to get tags of an article
----------------------------------------------

### Installation
`go get github.com/JesusIslam/tek`

### Dependencies
None

### Usage
```
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
```

Or just open `test/sample.go`

### Benchmark
Using i3-3217U @1.8GHz with 370 total words from the sample.txt provided and command `go test -bench . -benchtime=5s -cpu 4`:
```
BenchmarkGetTag-4            300          27111330 ns/op
```

### License
See LICENSE file, it is MIT