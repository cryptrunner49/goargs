package parser

import (
	"bytes"
	"reflect"
	"strings"
	"testing"
)

func TestParserLongFlags(t *testing.T) {
	args := []string{"--name=John", "--age=30", "install", "package"}
	out := &bytes.Buffer{}
	p := NewParser(out)
	p.SetProgramName("testprog")

	var name string
	var age int
	var verbose bool
	p.StringVar(&name, "n", "name", "default", "The name")
	p.IntVar(&age, "a", "age", 0, "The age")
	p.BoolVar(&verbose, "v", "verbose", false, "Verbose output")

	err := p.Parse(args)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	if name != "John" {
		t.Errorf("Expected name 'John', got %s", name)
	}
	if age != 30 {
		t.Errorf("Expected age 30, got %d", age)
	}
	if verbose != false {
		t.Errorf("Expected verbose false, got %t", verbose)
	}
	if !reflect.DeepEqual(p.Args(), []string{"install", "package"}) {
		t.Errorf("Expected args [install package], got %v", p.Args())
	}
}

func TestParserShortFlags(t *testing.T) {
	args := []string{"-n", "Jane", "-a=25", "-v", "push"}
	out := &bytes.Buffer{}
	p := NewParser(out)

	var name string
	var age int
	var verbose bool
	p.StringVar(&name, "n", "name", "default", "The name")
	p.IntVar(&age, "a", "age", 0, "The age")
	p.BoolVar(&verbose, "v", "verbose", false, "Verbose output")

	err := p.Parse(args)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	if name != "Jane" {
		t.Errorf("Expected name 'Jane', got %s", name)
	}
	if age != 25 {
		t.Errorf("Expected age 25, got %d", age)
	}
	if verbose != true {
		t.Errorf("Expected verbose true, got %t", verbose)
	}
	if !reflect.DeepEqual(p.Args(), []string{"push"}) {
		t.Errorf("Expected args [push], got %v", p.Args())
	}
}

func TestParserHelp(t *testing.T) {
	args := []string{"--help"}
	out := &bytes.Buffer{}
	p := NewParser(out)
	p.SetProgramName("testprog")

	var name string
	p.StringVar(&name, "n", "name", "default", "The name")

	err := p.Parse(args)
	if err != ErrHelpRequested {
		t.Errorf("Expected ErrHelpRequested, got %v", err)
	}
	if out.Len() == 0 {
		t.Error("Expected usage output, got none")
	}
}

func TestParserInvalidFlag(t *testing.T) {
	args := []string{"--unknown"}
	out := &bytes.Buffer{}
	p := NewParser(out)

	var name string
	p.StringVar(&name, "n", "name", "default", "The name")

	err := p.Parse(args)
	if err == nil || err.Error() != "unknown flag: unknown" {
		t.Errorf("Expected unknown flag error, got %v", err)
	}
}

func TestParserInvalidInt(t *testing.T) {
	args := []string{"--age=invalid"}
	out := &bytes.Buffer{}
	p := NewParser(out)

	var age int
	p.IntVar(&age, "a", "age", 0, "The age")

	err := p.Parse(args)
	if err == nil || !strings.Contains(err.Error(), "invalid value") {
		t.Errorf("Expected invalid value error, got %v", err)
	}
}

func TestParserEmptyArgs(t *testing.T) {
	args := []string{}
	out := &bytes.Buffer{}
	p := NewParser(out)

	var name string
	p.StringVar(&name, "n", "name", "default", "The name")

	err := p.Parse(args)
	if err != nil {
		t.Errorf("Unexpected error with empty args: %v", err)
	}
	if name != "default" {
		t.Errorf("Expected default name, got %s", name)
	}
	if len(p.Args()) != 0 {
		t.Errorf("Expected no positional args, got %v", p.Args())
	}
}
