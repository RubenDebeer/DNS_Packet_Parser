package main

import "fmt"

type Employee struct {
	firstName string
	lastName  string
	age       int
	isRemote  bool
	wage      int
	workHours int
}

func (e Employee) sallary() int {
	return e.workHours * e.wage
}

func main() {
	// ###################### Struct ######################
	employeeOne := Employee{
		firstName: "Rubas",
		lastName:  "De Beer",
		age:       25,
		isRemote:  true,
		wage:      200,
		workHours: 40,
	}

	fmt.Println("First Name", employeeOne.firstName)
	fmt.Println("Last Name", employeeOne.lastName)
	fmt.Println("Age", employeeOne.age)

	fmt.Println("Sallary", employeeOne.sallary())

	// ###################### Struct ######################

}
