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
	if !strings.Contains(out.String(), "Usage of testprog:") {
		t.Errorf("Expected usage output with program name, got %s", out.String())
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
	if age != 0 {
		t.Errorf("Expected age to remain default 0, got %d", age)
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

func TestUsageOutput(t *testing.T) {
	out := &bytes.Buffer{}
	p := NewParser(out)
	p.SetProgramName("testapp")

	var name string
	var age int
	var verbose bool
	p.StringVar(&name, "n", "name", "default", "User name")
	p.IntVar(&age, "a", "age", 0, "User age")
	p.BoolVar(&verbose, "v", "verbose", false, "Verbose mode")

	p.Usage()

	expected := `Usage of testapp:
  -n, --name string
        User name (default "default")
  -a, --age int
        User age (default 0)
  -v, --verbose
        Verbose mode (default false)
`
	if out.String() != expected {
		t.Errorf("Expected usage:\n%sGot:\n%s", expected, out.String())
	}
}

func TestSetProgramName(t *testing.T) {
	tests := []struct {
		name       string
		program    string
		wantPrefix string
	}{
		{"CustomName", "customapp", "Usage of customapp:"},
		{"DefaultName", "", "Usage of program:"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			out := &bytes.Buffer{}
			p := NewParser(out)
			if tt.program != "" {
				p.SetProgramName(tt.program)
			}
			var flag string
			p.StringVar(&flag, "f", "flag", "", "Test flag")

			p.Usage()
			if !strings.HasPrefix(out.String(), tt.wantPrefix) {
				t.Errorf("Expected usage to start with %q, got %s", tt.wantPrefix, out.String())
			}
		})
	}
}

func TestEmptyFlagNames(t *testing.T) {
	t.Run("OnlyShortName", func(t *testing.T) {
		out := &bytes.Buffer{}
		p := NewParser(out)
		p.SetProgramName("testapp")
		var name string
		p.StringVar(&name, "n", "", "default", "User name")

		// Test parsing
		err := p.Parse([]string{"-n", "test"})
		if err != nil {
			t.Fatalf("Unexpected error: %v", err)
		}
		if name != "test" {
			t.Errorf("Expected name 'test', got %s", name)
		}

		// Long name should fail
		err = p.Parse([]string{"--name=test"})
		if err == nil || err.Error() != "unknown flag: name" {
			t.Errorf("Expected unknown flag error, got %v", err)
		}

		// Test usage output
		out.Reset()
		p.Usage()
		expected := `Usage of testapp:
  -n string
        User name (default "default")
`
		if out.String() != expected {
			t.Errorf("Expected usage:\n%sGot:\n%s", expected, out.String())
		}
	})

	t.Run("OnlyLongName", func(t *testing.T) {
		out := &bytes.Buffer{}
		p := NewParser(out)
		p.SetProgramName("testapp")
		var name string
		p.StringVar(&name, "", "name", "default", "User name")

		// Test parsing
		err := p.Parse([]string{"--name=test"})
		if err != nil {
			t.Fatalf("Unexpected error: %v", err)
		}
		if name != "test" {
			t.Errorf("Expected name 'test', got %s", name)
		}

		// Short name should fail
		err = p.Parse([]string{"-n", "test"})
		if err == nil || err.Error() != "unknown flag: n" {
			t.Errorf("Expected unknown flag error, got %v", err)
		}

		// Test usage output
		out.Reset()
		p.Usage()
		expected := `Usage of testapp:
  --name string
        User name (default "default")
`
		if out.String() != expected {
			t.Errorf("Expected usage:\n%sGot:\n%s", expected, out.String())
		}
	})
}

func TestDuplicateFlagNames(t *testing.T) {
	out := &bytes.Buffer{}
	p := NewParser(out)
	p.SetProgramName("testapp")

	var name1, name2 string
	p.StringVar(&name1, "n", "name", "default1", "First name")
	p.StringVar(&name2, "n", "name", "default2", "Second name") // Duplicate flags

	err := p.Parse([]string{"-n", "test"})
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}
	if name2 != "test" {
		t.Errorf("Expected name2 'test', got %s", name2)
	}
	if name1 != "default1" {
		t.Errorf("Expected name1 to retain default 'default1', got %s", name1)
	}

	out.Reset()
	p.Usage()
	expected := `Usage of testapp:
  -n, --name string
        First name (default "default1")
  -n, --name string
        Second name (default "default2")
`
	if out.String() != expected {
		t.Errorf("Expected usage:\n%sGot:\n%s", expected, out.String())
	}
}

func TestMixedFlagUsage(t *testing.T) {
	out := &bytes.Buffer{}
	p := NewParser(out)
	p.SetProgramName("testapp")

	var shortOnly string
	var longOnly string
	var both string
	p.StringVar(&shortOnly, "s", "", "short", "Short only flag")
	p.StringVar(&longOnly, "", "long", "long", "Long only flag")
	p.StringVar(&both, "b", "both", "both", "Both flags")

	err := p.Parse([]string{"-s", "sval", "--long=longval", "-b", "bval"})
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}
	if shortOnly != "sval" {
		t.Errorf("Expected shortOnly 'sval', got %s", shortOnly)
	}
	if longOnly != "longval" {
		t.Errorf("Expected longOnly 'longval', got %s", longOnly)
	}
	if both != "bval" {
		t.Errorf("Expected both 'bval', got %s", both)
	}

	out.Reset()
	p.Usage()
	expected := `Usage of testapp:
  -s string
        Short only flag (default "short")
  --long string
        Long only flag (default "long")
  -b, --both string
        Both flags (default "both")
`
	if out.String() != expected {
		t.Errorf("Expected usage:\n%sGot:\n%s", expected, out.String())
	}
}

func TestBoolFlagEdgeCases(t *testing.T) {
	tests := []struct {
		name    string
		args    []string
		wantVal bool
		wantErr bool
		errMsg  string
	}{
		{"FlagPresent", []string{"-v"}, true, false, ""},
		{"ExplicitTrue", []string{"-v=true"}, true, false, ""},
		{"ExplicitFalse", []string{"-v=false"}, false, false, ""},
		{"InvalidValue", []string{"-v=invalid"}, false, true, "invalid value for bool flag -v/--verbose: invalid"},
		{"NotPresent", []string{}, false, false, ""},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			out := &bytes.Buffer{}
			p := NewParser(out)
			var verbose bool
			p.BoolVar(&verbose, "v", "verbose", false, "Verbose mode")

			err := p.Parse(tt.args)
			if tt.wantErr {
				if err == nil || err.Error() != tt.errMsg {
					t.Errorf("Expected error %q, got %v", tt.errMsg, err)
				}
			} else if err != nil {
				t.Fatalf("Unexpected error: %v", err)
			}

			if verbose != tt.wantVal {
				t.Errorf("Expected verbose %t, got %t", tt.wantVal, verbose)
			}
		})
	}
}
