package main

import (
	"fmt"

	"github.com/cryptrunner49/argonaut/parser"
)

func main() {
	// Create a new parser
	p := parser.NewParser()

	// Define variables and register flags
	var name string
	var age int
	var verbose bool
	p.StringVar(&name, "name", "default", "The name of the user")
	p.IntVar(&age, "age", 0, "The age of the user")
	p.BoolVar(&verbose, "verbose", false, "Enable verbose output")

	// Parse command-line arguments
	p.Parse()

	// Print parsed values and positional arguments
	fmt.Println("Name:", name)
	fmt.Println("Age:", age)
	fmt.Println("Verbose:", verbose)
	fmt.Println("Positional arguments:", p.Args())
}
