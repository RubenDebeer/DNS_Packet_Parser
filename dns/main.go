package main

import "fmt"

func main() {
	a := "Hello World"

	fmt.Println("Before modifyString:", a)
	modifyString(&a)
	fmt.Println("After modifyString:", a)
}

// Try and Figure out what is a Pointer and how it works

func modifyString(a *string) {
	*a = "Modified String"
}
