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
	shortName    string      // Short flag name (e.g., "v")
	longName     string      // Long flag name (e.g., "verbose")
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

// StringVar registers a string flag with short and long names, default value, and description.
func (p *Parser) StringVar(ptr *string, shortName, longName string, defaultValue string, description string) {
	*ptr = defaultValue // Set the variable to its default value
	p.flags = append(p.flags, flagInfo{shortName, longName, ptr, defaultValue, description, "string"})
}

// IntVar registers an integer flag with short and long names, default value, and description.
func (p *Parser) IntVar(ptr *int, shortName, longName string, defaultValue int, description string) {
	*ptr = defaultValue
	p.flags = append(p.flags, flagInfo{shortName, longName, ptr, defaultValue, description, "int"})
}

// BoolVar registers a boolean flag with short and long names, default value, and description.
func (p *Parser) BoolVar(ptr *bool, shortName, longName string, defaultValue bool, description string) {
	*ptr = defaultValue
	p.flags = append(p.flags, flagInfo{shortName, longName, ptr, defaultValue, description, "bool"})
}

// Parse processes command-line arguments and sets flag values.
// It also handles the --help flag to display usage information.
func (p *Parser) Parse() {
	args := os.Args[1:] // Skip the program name

	// Check for --help or -h flag
	if contains(args, "--help") || contains(args, "-h") {
		p.Usage()
		os.Exit(0)
	}

	// Map to store flags and their values from the command line
	setFlags := make(map[string]string)
	p.positional = []string{}

	// Parse each argument
	for i := 0; i < len(args); i++ {
		arg := args[i]
		if strings.HasPrefix(arg, "--") {
			// Long flag
			if strings.Contains(arg, "=") {
				parts := strings.SplitN(arg[2:], "=", 2)
				flag := parts[0]
				value := parts[1]
				setFlags[flag] = value
			} else {
				flag := arg[2:]
				// Check if next arg is a value (not a flag or positional)
				if i+1 < len(args) && !strings.HasPrefix(args[i+1], "-") && p.needsValue(flag) {
					setFlags[flag] = args[i+1]
					i++
				} else {
					setFlags[flag] = "true" // Default for bool flags
				}
			}
		} else if strings.HasPrefix(arg, "-") && arg != "-" {
			// Short flag
			flag := arg[1:]
			if strings.Contains(arg, "=") {
				parts := strings.SplitN(arg[1:], "=", 2)
				flag = parts[0]
				value := parts[1]
				setFlags[flag] = value
			} else if i+1 < len(args) && !strings.HasPrefix(args[i+1], "-") && p.needsValue(flag) {
				setFlags[flag] = args[i+1]
				i++
			} else {
				setFlags[flag] = "true" // Default for bool flags
			}
		} else {
			// Positional argument (no dashes)
			p.positional = append(p.positional, arg)
		}
	}

	// Set registered flag values
	for _, fi := range p.flags {
		// Check both short and long names
		for _, name := range []string{fi.shortName, fi.longName} {
			if val, ok := setFlags[name]; ok {
				switch fi.kind {
				case "string":
					*fi.ptr.(*string) = val
				case "int":
					i, err := strconv.Atoi(val)
					if err != nil {
						fmt.Printf("Invalid value for flag -%s/--%s: %s\n", fi.shortName, fi.longName, val)
						os.Exit(1)
					}
					*fi.ptr.(*int) = i
				case "bool":
					if val == "true" {
						*fi.ptr.(*bool) = true
					} else if val == "false" {
						*fi.ptr.(*bool) = false
					} else {
						fmt.Printf("Invalid value for bool flag -%s/--%s: %s\n", fi.shortName, fi.longName, val)
						os.Exit(1)
					}
				}
			}
		}
		// If flag not provided, it retains its default value set during registration
	}

	// Check for unknown flags
	for flag := range setFlags {
		if !p.hasFlag(flag) {
			fmt.Printf("Unknown flag: %s\n", flag)
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
			fmt.Printf("  -%s, --%s string\n        %s (default %q)\n", fi.shortName, fi.longName, fi.description, fi.defaultValue)
		case "int":
			fmt.Printf("  -%s, --%s int\n        %s (default %d)\n", fi.shortName, fi.longName, fi.description, fi.defaultValue)
		case "bool":
			fmt.Printf("  -%s, --%s\n        %s (default %t)\n", fi.shortName, fi.longName, fi.description, fi.defaultValue)
		}
	}
}

// hasFlag checks if a flag (short or long) is registered.
func (p *Parser) hasFlag(name string) bool {
	for _, fi := range p.flags {
		if fi.shortName == name || fi.longName == name {
			return true
		}
	}
	return false
}

// needsValue checks if a flag requires a value (i.e., not a bool flag).
func (p *Parser) needsValue(name string) bool {
	for _, fi := range p.flags {
		if (fi.shortName == name || fi.longName == name) && fi.kind != "bool" {
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
