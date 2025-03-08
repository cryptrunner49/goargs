package parser

import (
	"os"
	"testing"
)

func TestParser(t *testing.T) {
	// Test string and int flags with positional arguments
	os.Args = []string{"program", "--name=John", "--age=30", "file1", "file2"}
	p := NewParser()
	var name string
	var age int
	var verbose bool
	p.StringVar(&name, "name", "default", "The name")
	p.IntVar(&age, "age", 0, "The age")
	p.BoolVar(&verbose, "verbose", false, "Verbose output")
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
	if len(args) != 2 || args[0] != "file1" || args[1] != "file2" {
		t.Errorf("Expected args [file1 file2], got %v", args)
	}

	// Test bool flag
	os.Args = []string{"program", "--verbose"}
	p = NewParser()
	p.BoolVar(&verbose, "verbose", false, "Verbose output")
	p.Parse()
	if verbose != true {
		t.Errorf("Expected verbose true, got %t", verbose)
	}

	// Test default values
	os.Args = []string{"program"}
	p = NewParser()
	p.StringVar(&name, "name", "default", "The name")
	p.IntVar(&age, "age", 0, "The age")
	p.Parse()
	if name != "default" {
		t.Errorf("Expected name 'default', got %s", name)
	}
	if age != 0 {
		t.Errorf("Expected age 0, got %d", age)
	}
}
