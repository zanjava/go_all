package main

import (
	"fmt"
)

type Person struct {
	Name string
	Age  int
}

func (p Person) Greet() string {
	p.Name = "Bob" // This modification won't affect the original struct
	return fmt.Sprintf("Hello, my name is %s and I am %d years old.", p.Name, p.Age)
}

func (p *Person) Greet1() string {
	return fmt.Sprintf("Hello, my name is %s and I am %d years old.", p.Name, p.Age)
}

func main() {
	person := Person{Name: "Alice", Age: 30}
	fmt.Println(person.Greet())
	fmt.Println(person.Name) // Outputs "Alice", not "Bob"
	fmt.Println(person.Greet1())
}
