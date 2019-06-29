# Patchwerk

![Patchwerk logo](https://github.com/herkyl/patchwerk/blob/master/patchwerk.jpg)

Use Patchwerk to create [RFC6902 JSON patches](https://tools.ietf.org/html/rfc6902).

At the moment of writing this is the only working Go library for creating JSON patches. If you wish to apply the patches I recommend using [evanphx/json-patch](https://github.com/evanphx/json-patch) (it only allows for applying patches, not generating them).


The project was originally cloned from [mattbaird/jsonpatch](https://github.com/mattbaird/jsonpatch).

## Installation

```bash
go get github.com/herkyl/patchwerk
```

## Usage

```go
package main

import (
	"fmt"
	"github.com/herkyl/patchwerk"
)

func main() {
	a := `{"a":100, "b":200}`
	b := `{"a":100, "b":200, "c":300}`
	patch, err := patchwerk.DiffBytes([]byte(a), []byte(b))
	if err != nil {
		fmt.Printf("Error creating JSON patch: %v", err)
		return
	}
	fmt.Println(string(patch)) // [{"op": "add", "path": "/c", "value": 300}]
}

```