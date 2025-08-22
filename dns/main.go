package main

import "fmt"

func main() {

	i, j := 12, 40

	// Print the Values
	fmt.Println("Print The Values")
	fmt.Println("i value --> ", i)
	fmt.Println("j Value --> ", j)

	// Print the address to the values
	fmt.Println("Print The address to the values")
	fmt.Println(" address of I --> ", &i)
	fmt.Println(" address of J  --> ", &j)

	pi := &i
	pj := &j

	fmt.Println("Print The address to the values")
	fmt.Println(" address of I --> ", pi)
	fmt.Println(" address of J  --> ", pj)

	// What is a the * shit then ?

	fmt.Println("Print the value of the addresses")
	fmt.Println("Value at address of pi :", *pi)
	fmt.Println("Value at address of pj :", *pj)

	// *pj --> pi --> i  This is Dereferencing.

	fmt.Println("Change the dereferenced value", *pi)
	*pi = 800
	fmt.Println("Change the dereferenced value", *pi)

}
