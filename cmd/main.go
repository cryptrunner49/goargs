package main

import (
	"fmt"

	"github.com/cryptrunner49/argonaut/parser"
)

func main() {
	// Create a new parser
	p := parser.NewParser()

	// Define variables and register flags with short and long names
	var name string
	var age int
	var verbose bool
	p.StringVar(&name, "n", "name", "default", "The name of the user")
	p.IntVar(&age, "a", "age", 0, "The age of the user")
	p.BoolVar(&verbose, "v", "verbose", false, "Enable verbose output")

	// Parse command-line arguments
	p.Parse()

	// Print parsed values and positional arguments
	fmt.Println("Name:", name)
	fmt.Println("Age:", age)
	fmt.Println("Verbose:", verbose)
	fmt.Println("Positional arguments:", p.Args())
}
