package parser

import (
	"errors"
	"fmt"
	"io"
	"strconv"
	"strings"
)

// Parser represents a command-line argument parser.
type Parser struct {
	flags      []flagInfo // Stores registered flags
	positional []string   // Stores positional arguments
	output     io.Writer  // Where to write usage and errors
	program    string     // Program name for usage (optional)
}

// flagInfo holds metadata about a registered flag.
type flagInfo struct {
	shortName    string      // Short flag name (e.g., "v")
	longName     string      // Long flag name (e.g., "verbose")
	ptr          interface{} // Pointer to the variable to store the value
	defaultValue interface{} // Default value of the flag
	description  string      // Description for documentation
	kind         string      // Type of the flag: "string", "int", or "bool"
}

// NewParser creates and initializes a new Parser instance with an output writer.
func NewParser(out io.Writer) *Parser {
	return &Parser{
		flags:      []flagInfo{},
		positional: []string{},
		output:     out,
		program:    "program", // Default program name
	}
}

// SetProgramName sets the program name for usage output.
func (p *Parser) SetProgramName(name string) {
	p.program = name
}

// StringVar registers a string flag.
func (p *Parser) StringVar(ptr *string, shortName, longName string, defaultValue string, description string) {
	*ptr = defaultValue
	p.flags = append(p.flags, flagInfo{shortName, longName, ptr, defaultValue, description, "string"})
}

// IntVar registers an integer flag.
func (p *Parser) IntVar(ptr *int, shortName, longName string, defaultValue int, description string) {
	*ptr = defaultValue
	p.flags = append(p.flags, flagInfo{shortName, longName, ptr, defaultValue, description, "int"})
}

// BoolVar registers a boolean flag.
func (p *Parser) BoolVar(ptr *bool, shortName, longName string, defaultValue bool, description string) {
	*ptr = defaultValue
	p.flags = append(p.flags, flagInfo{shortName, longName, ptr, defaultValue, description, "bool"})
}

// Parse processes command-line arguments and returns an error if parsing fails.
// If --help or -h is present, it prints usage and returns a special error.
func (p *Parser) Parse(args []string) error {
	// Check for help flags
	if contains(args, "--help") || contains(args, "-h") {
		p.Usage()
		return ErrHelpRequested
	}

	// Reset positional arguments
	p.positional = []string{}
	setFlags := make(map[string]string)

	// Parse arguments
	for i := 0; i < len(args); i++ {
		arg := args[i]
		if strings.HasPrefix(arg, "--") {
			if strings.Contains(arg, "=") {
				parts := strings.SplitN(arg[2:], "=", 2)
				setFlags[parts[0]] = parts[1]
			} else {
				flag := arg[2:]
				if i+1 < len(args) && !strings.HasPrefix(args[i+1], "-") && p.needsValue(flag) {
					setFlags[flag] = args[i+1]
					i++
				} else {
					setFlags[flag] = "true"
				}
			}
		} else if strings.HasPrefix(arg, "-") && arg != "-" {
			flag := arg[1:]
			if strings.Contains(arg, "=") {
				parts := strings.SplitN(arg[1:], "=", 2)
				setFlags[parts[0]] = parts[1]
			} else if i+1 < len(args) && !strings.HasPrefix(args[i+1], "-") && p.needsValue(flag) {
				setFlags[flag] = args[i+1]
				i++
			} else {
				setFlags[flag] = "true"
			}
		} else {
			p.positional = append(p.positional, arg)
		}
	}

	// Set flag values, processing in reverse to ensure last registration wins
	for i := len(p.flags) - 1; i >= 0; i-- {
		fi := p.flags[i]
		for _, name := range []string{fi.shortName, fi.longName} {
			if val, ok := setFlags[name]; ok {
				switch fi.kind {
				case "string":
					*fi.ptr.(*string) = val
				case "int":
					i, err := strconv.Atoi(val)
					if err != nil {
						return fmt.Errorf("invalid value for flag -%s/--%s: %s", fi.shortName, fi.longName, val)
					}
					*fi.ptr.(*int) = i
				case "bool":
					if val == "true" {
						*fi.ptr.(*bool) = true
					} else if val == "false" {
						*fi.ptr.(*bool) = false
					} else {
						return fmt.Errorf("invalid value for bool flag -%s/--%s: %s", fi.shortName, fi.longName, val)
					}
				}
				// Remove the flag from setFlags to prevent earlier duplicates from being set
				delete(setFlags, name)
			}
		}
	}

	// Check for unknown flags
	for flag := range setFlags {
		if !p.hasFlag(flag) {
			return fmt.Errorf("unknown flag: %s", flag)
		}
	}

	return nil
}

// Args returns the list of positional arguments.
func (p *Parser) Args() []string {
	return p.positional
}

// Usage generates and writes the usage documentation.
func (p *Parser) Usage() {
	fmt.Fprintf(p.output, "Usage of %s:\n", p.program)
	for _, fi := range p.flags {
		var flagNames string
		switch {
		case fi.shortName != "" && fi.longName != "":
			flagNames = fmt.Sprintf("-%s, --%s", fi.shortName, fi.longName)
		case fi.shortName != "":
			flagNames = fmt.Sprintf("-%s", fi.shortName)
		case fi.longName != "":
			flagNames = fmt.Sprintf("--%s", fi.longName)
		}

		switch fi.kind {
		case "string":
			fmt.Fprintf(p.output, "  %s string\n        %s (default %q)\n", flagNames, fi.description, fi.defaultValue)
		case "int":
			fmt.Fprintf(p.output, "  %s int\n        %s (default %d)\n", flagNames, fi.description, fi.defaultValue)
		case "bool":
			fmt.Fprintf(p.output, "  %s\n        %s (default %t)\n", flagNames, fi.description, fi.defaultValue)
		}
	}
}

// hasFlag checks if a flag is registered.
func (p *Parser) hasFlag(name string) bool {
	for _, fi := range p.flags {
		if fi.shortName == name || fi.longName == name {
			return true
		}
	}
	return false
}

// needsValue checks if a flag requires a value.
func (p *Parser) needsValue(name string) bool {
	for _, fi := range p.flags {
		if (fi.shortName == name || fi.longName == name) && fi.kind != "bool" {
			return true
		}
	}
	return false
}

// contains checks if a string slice includes an item.
func contains(slice []string, item string) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
}

// ErrHelpRequested is returned when --help or -h is encountered.
var ErrHelpRequested = errors.New("help requested")
