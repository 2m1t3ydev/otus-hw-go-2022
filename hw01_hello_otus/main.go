package main

import (
	"fmt"

	"golang.org/x/example/stringutil"
)

func main() {
	sourceString := "Hello, OTUS!"
	reversedString := stringutil.Reverse(sourceString)
	fmt.Println(reversedString)
}
