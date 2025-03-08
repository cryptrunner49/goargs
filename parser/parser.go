package parser

import (
	"fmt"
	"os"
	"strconv"
	"strings"
)

// Parser represents a command-line argument parser.
type Parser struct {
	flags      []flagInfo // Stores registered flags
	positional []string   // Stores positional arguments
}

// flagInfo holds metadata about a registered flag.
type flagInfo struct {
	name         string      // Flag name (e.g., "verbose")
	ptr          interface{} // Pointer to the variable to store the value
	defaultValue interface{} // Default value of the flag
	description  string      // Description for documentation
	kind         string      // Type of the flag: "string", "int", or "bool"
}

// NewParser creates and initializes a new Parser instance.
func NewParser() *Parser {
	return &Parser{
		flags:      []flagInfo{},
		positional: []string{},
	}
}

// StringVar registers a string flag with a name, default value, and description.
func (p *Parser) StringVar(ptr *string, name string, defaultValue string, description string) {
	*ptr = defaultValue // Set the variable to its default value
	p.flags = append(p.flags, flagInfo{name, ptr, defaultValue, description, "string"})
}

// IntVar registers an integer flag with a name, default value, and description.
func (p *Parser) IntVar(ptr *int, name string, defaultValue int, description string) {
	*ptr = defaultValue
	p.flags = append(p.flags, flagInfo{name, ptr, defaultValue, description, "int"})
}

// BoolVar registers a boolean flag with a name, default value, and description.
func (p *Parser) BoolVar(ptr *bool, name string, defaultValue bool, description string) {
	*ptr = defaultValue
	p.flags = append(p.flags, flagInfo{name, ptr, defaultValue, description, "bool"})
}

// Parse processes command-line arguments and sets flag values.
// It also handles the --help flag to display usage information.
func (p *Parser) Parse() {
	args := os.Args[1:] // Skip the program name

	// Check for --help flag
	if contains(args, "--help") {
		p.Usage()
		os.Exit(0)
	}

	// Map to store flags and their values from the command line
	setFlags := make(map[string]string)
	p.positional = []string{}

	// Parse each argument
	for _, arg := range args {
		if strings.HasPrefix(arg, "--") {
			if strings.Contains(arg, "=") {
				// Flag with value, e.g., --name=John
				parts := strings.SplitN(arg[2:], "=", 2)
				flag := parts[0]
				value := parts[1]
				setFlags[flag] = value
			} else {
				// Bool flag without value, e.g., --verbose, defaults to true
				flag := arg[2:]
				setFlags[flag] = "true"
			}
		} else {
			// Positional argument
			p.positional = append(p.positional, arg)
		}
	}

	// Set registered flag values
	for _, fi := range p.flags {
		if val, ok := setFlags[fi.name]; ok {
			switch fi.kind {
			case "string":
				*fi.ptr.(*string) = val
			case "int":
				i, err := strconv.Atoi(val)
				if err != nil {
					fmt.Printf("Invalid value for flag --%s: %s\n", fi.name, val)
					os.Exit(1)
				}
				*fi.ptr.(*int) = i
			case "bool":
				if val == "true" {
					*fi.ptr.(*bool) = true
				} else if val == "false" {
					*fi.ptr.(*bool) = false
				} else {
					fmt.Printf("Invalid value for bool flag --%s: %s\n", fi.name, val)
					os.Exit(1)
				}
			}
		}
		// If flag not provided, it retains its default value set during registration
	}

	// Check for unknown flags
	for flag := range setFlags {
		if !p.hasFlag(flag) {
			fmt.Printf("Unknown flag: --%s\n", flag)
			os.Exit(1)
		}
	}
}

// Args returns the list of positional arguments.
func (p *Parser) Args() []string {
	return p.positional
}

// Usage generates and prints the documentation for all registered flags.
func (p *Parser) Usage() {
	fmt.Printf("Usage of %s:\n", os.Args[0])
	for _, fi := range p.flags {
		switch fi.kind {
		case "string":
			fmt.Printf("  --%s string\n        %s (default %q)\n", fi.name, fi.description, fi.defaultValue)
		case "int":
			fmt.Printf("  --%s int\n        %s (default %d)\n", fi.name, fi.description, fi.defaultValue)
		case "bool":
			fmt.Printf("  --%s\n        %s (default %t)\n", fi.name, fi.description, fi.defaultValue)
		}
	}
}

// hasFlag checks if a flag is registered.
func (p *Parser) hasFlag(name string) bool {
	for _, fi := range p.flags {
		if fi.name == name {
			return true
		}
	}
	return false
}

// contains checks if a string slice includes a specific item.
func contains(slice []string, item string) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
}
