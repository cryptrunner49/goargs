package parser

import (
	"os"
	"testing"
)

func TestParser(t *testing.T) {
	// Test long flags with positional arguments
	os.Args = []string{"program", "--name=John", "--age=30", "install", "package"}
	p := NewParser()
	var name string
	var age int
	var verbose bool
	p.StringVar(&name, "n", "name", "default", "The name")
	p.IntVar(&age, "a", "age", 0, "The age")
	p.BoolVar(&verbose, "v", "verbose", false, "Verbose output")
	p.Parse()

	if name != "John" {
		t.Errorf("Expected name 'John', got %s", name)
	}
	if age != 30 {
		t.Errorf("Expected age 30, got %d", age)
	}
	if verbose != false {
		t.Errorf("Expected verbose false, got %t", verbose)
	}
	args := p.Args()
	if len(args) != 2 || args[0] != "install" || args[1] != "package" {
		t.Errorf("Expected args [install package], got %v", args)
	}

	// Test short flags
	os.Args = []string{"program", "-n", "Jane", "-a=25", "-v", "push"}
	p = NewParser()
	p.StringVar(&name, "n", "name", "default", "The name")
	p.IntVar(&age, "a", "age", 0, "The age")
	p.BoolVar(&verbose, "v", "verbose", false, "Verbose output")
	p.Parse()
	if name != "Jane" {
		t.Errorf("Expected name 'Jane', got %s", name)
	}
	if age != 25 {
		t.Errorf("Expected age 25, got %d", age)
	}
	if verbose != true {
		t.Errorf("Expected verbose true, got %t", verbose)
	}
	args = p.Args()
	if len(args) != 1 || args[0] != "push" {
		t.Errorf("Expected args [push], got %v", args)
	}

	// Test bool flag with long name
	os.Args = []string{"program", "--verbose"}
	p = NewParser()
	p.BoolVar(&verbose, "v", "verbose", false, "Verbose output")
	p.Parse()
	if verbose != true {
		t.Errorf("Expected verbose true, got %t", verbose)
	}

	// Test default values with bare command
	os.Args = []string{"program", "install"}
	p = NewParser()
	p.StringVar(&name, "n", "name", "default", "The name")
	p.IntVar(&age, "a", "age", 0, "The age")
	p.Parse()
	if name != "default" {
		t.Errorf("Expected name 'default', got %s", name)
	}
	if age != 0 {
		t.Errorf("Expected age 0, got %d", age)
	}
	args = p.Args()
	if len(args) != 1 || args[0] != "install" {
		t.Errorf("Expected args [install], got %v", args)
	}

	// Test mixed short and long flags
	os.Args = []string{"program", "-n=Bob", "--age", "40", "build"}
	p = NewParser()
	p.StringVar(&name, "n", "name", "default", "The name")
	p.IntVar(&age, "a", "age", 0, "The age")
	p.BoolVar(&verbose, "v", "verbose", false, "Verbose output")
	p.Parse()
	if name != "Bob" {
		t.Errorf("Expected name 'Bob', got %s", name)
	}
	if age != 40 {
		t.Errorf("Expected age 40, got %d", age)
	}
	if verbose != false {
		t.Errorf("Expected verbose false, got %t", verbose)
	}
	args = p.Args()
	if len(args) != 1 || args[0] != "build" {
		t.Errorf("Expected args [build], got %v", args)
	}
}
