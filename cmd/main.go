package main

import (
	"fmt"
	"os"

	"github.com/cryptrunner49/argonaut/parser"
)

func main() {
	// Create a new parser with output to os.Stdout
	p := parser.NewParser(os.Stdout)

	// Define variables and register flags
	var name string
	var age int
	var verbose bool
	p.StringVar(&name, "n", "name", "default", "The name of the user")
	p.IntVar(&age, "a", "age", 0, "The age of the user")
	p.BoolVar(&verbose, "v", "verbose", false, "Enable verbose output")

	// Parse command-line arguments and handle errors
	err := p.Parse(os.Args[1:])
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}

	// Print parsed values
	fmt.Println("Name:", name)
	fmt.Println("Age:", age)
	fmt.Println("Verbose:", verbose)
	fmt.Println("Positional arguments:", p.Args())
}
