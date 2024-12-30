package main

import (
	providers_example "layered-stack/cli/providers"
	"os"
)

func main() {
	name := "World"
	if len(os.Args) > 1 {
		name = os.Args[1]
	}
	println(providers_example.SayHello(name))
}
