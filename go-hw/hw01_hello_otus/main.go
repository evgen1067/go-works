package main

import (
	"fmt"
	"github.com/vigo/stringutils-demo"
)

func main() {
	reversed, _ := stringutils.Reverse("Hello, OTUS!")
	fmt.Println(reversed)
}
