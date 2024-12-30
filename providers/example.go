package providers_example

import "fmt"

func SayHello(name string) string {
	const defaultName = "World"
	if name == "" {
		name = defaultName
	}
	return fmt.Sprintf("Hello, %s!", name)
}
