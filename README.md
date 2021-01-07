# go-dpkg-parser
fetch debian repository's metadata and parse it without dpkg/apt


```
package main

import (
	"fmt"

	"github.com/snyh/go-dpkg-parser"
)

func main() {
	r := dpkg.NewRepository("cache") // where to save the local cache data

	//         source url and the suite name
	r.AddSuite("http://mirrors.163.com/debian/", "stable", "")

	a, err := r.Archive("amd64") //architecture name
	if err != nil {
		panic(err)
	}
	for name, p := range a.Packages {
		fmt.Println(name, "---------", p.Get("filename"))
	}
}

```
