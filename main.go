package main

import (
	"fmt"
	"sha256/hash"
)

func main() {
	fmt.Println(hash.Compute([]byte("Hello World!")))
}
