# Argonaut Documentation

## Argonaut Parser Library

The `Argonaut Parser` is a simple, lightweight command-line argument parsing library for Go. It allows developers to define and parse flags (string, int, and bool) with optional short (e.g., `-v`) and long (e.g., `--verbose`) forms, and handle positional arguments, including bare commands like `install` or `push`. This library is ideal for small to medium-sized CLI applications where you need flexible, test-friendly flag parsing without additional dependencies.

This document provides a comprehensive guide to using the library, including installation instructions, usage examples, and a full API reference.

---

## Table of Contents

1. [Installation](#installation)
2. [Quick Start](#quick-start)
3. [Usage](#usage)
   - [Defining Flags](#defining-flags)
     - [Using Short Flags](#using-short-flags)
     - [Using Long Flags](#using-long-flags)
   - [Parsing Arguments](#parsing-arguments)
   - [Accessing Positional Arguments and Bare Commands](#accessing-positional-arguments-and-bare-commands)
   - [Displaying Usage Information](#displaying-usage-information)
   - [Error Handling](#error-handling)
4. [Examples](#examples)
   - [Bare Commands Only](#bare-commands-only)
   - [Short Flags Only](#short-flags-only)
   - [Long Flags Only](#long-flags-only)
   - [Mixed Short and Long Flags with Commands](#mixed-short-and-long-flags-with-commands)
5. [API Reference](#api-reference)
   - [Parser Struct](#parser-struct)
   - [NewParser](#newparser)
   - [SetProgramName](#setprogramname)
   - [StringVar](#stringvar)
   - [IntVar](#intvar)
   - [BoolVar](#boolvar)
   - [Parse](#parse)
   - [Args](#args)
   - [Usage](#usage)
6. [Testing](#testing)
7. [Contributing](#contributing)
8. [License](#license)

---

## Installation

To use the `Argonaut Parser` library in your Go project, follow these steps:

1. Ensure you have Go installed (version 1.11 or later recommended).
2. Add the library to your project using Go modules:

```bash
go get github.com/cryptrunner49/argonaut/parser
```

3. Import the library in your Go code:

```go
import "github.com/cryptrunner49/argonaut/parser"
```

The library has no external dependencies beyond the Go standard library.

---

## Quick Start

Here’s a minimal example using bare commands and a flag, with manual usage display:

```go
package main

import (
 "fmt"
 "os"
 "github.com/cryptrunner49/argonaut/parser"
)

func main() {
 p := parser.NewParser(os.Stdout)
 p.SetProgramName("myapp")

 var verbose bool
 p.BoolVar(&verbose, "v", "verbose", false, "Enable verbose output")

 err := p.Parse(os.Args[1:])
 if err != nil {
  if err == parser.ErrHelpRequested {
   return
  }
  fmt.Fprintf(os.Stderr, "Error: %v\n", err)
  os.Exit(1)
 }

 // Show usage if no arguments provided
 if len(os.Args) == 1 {
  p.Usage()
  return
 }

 args := p.Args()
 if len(args) > 0 {
  fmt.Printf("Command: %s, Verbose: %t\n", args[0], verbose)
 }
}
```

Run it:

```bash
go run main.go start -v
# Output:
# Command: start, Verbose: true

go run main.go
# Output:
# Usage of myapp:
#   -v, --verbose
#         Enable verbose output (default false)
```

---

## Usage

### Defining Flags

Flags are defined using `StringVar`, `IntVar`, and `BoolVar`. Each method registers a flag with optional short and long names, a default value, and a description. Set unused names to `""` to use only short or long flags.

#### Using Short Flags

- **Syntax**: `-flag=value` or `-flag value` (for string and int), `-flag` (for bool, sets to `true`).
- **Setup**: Provide a short name (e.g., `"v"`) and set `longName` to `""`.
- **Note**: Long flags are ignored if `longName` is empty.

#### Using Long Flags

- **Syntax**: `--flag=value` or `--flag value` (for string and int), `--flag` (for bool, sets to `true`).
- **Setup**: Provide a long name (e.g., `"verbose"`) and set `shortName` to `""`.
- **Note**: Short flags are ignored if `shortName` is empty.

- **Default Values**: Flags retain their default value if not provided.
- **Duplicates**: If multiple flags share the same name, the last registered flag takes precedence.

### Parsing Arguments

Call `Parse(args []string)` to process command-line arguments. It sets flag values, collects positional arguments, and returns an error for invalid inputs or help requests.

### Accessing Positional Arguments and Bare Commands

Use `Args()` to retrieve a slice of positional arguments, which include bare commands (non-flag arguments like `install` or `push`). Handle these manually as needed.

- **Bare Commands**: Non-flag arguments are treated as positional. Use without flags for bare command focus or combine with flags.

### Displaying Usage Information

The library automatically handles `--help` and `-h` flags, printing usage to the `io.Writer` and returning `ErrHelpRequested`. Use `SetProgramName` to customize the program name. Call `Usage()` manually to display usage at any time (e.g., when no arguments are provided).

### Error Handling

`Parse` returns an error for:

- `ErrHelpRequested`: When `--help` or `-h` is provided.
- Invalid flag values (e.g., non-integer for int flags).
- Unknown flags.

Handle errors appropriately, exiting gracefully for help requests.

---

## Examples

### Bare Commands Only

A program using only bare commands:

```go
package main

import (
 "fmt"
 "os"
 "github.com/cryptrunner49/argonaut/parser"
)

func main() {
 p := parser.NewParser(os.Stdout)
 p.SetProgramName("myapp")

 err := p.Parse(os.Args[1:])
 if err != nil {
  if err == parser.ErrHelpRequested {
   return
  }
  fmt.Fprintf(os.Stderr, "Error: %v\n", err)
  os.Exit(1)
 }

 if len(os.Args) == 1 {
  p.Usage()
  return
 }

 fmt.Println("Commands:", p.Args())
}
```

Run it:

```bash
go run main.go push origin
# Output:
# Commands: [push origin]
```

### Short Flags Only

A program using only short flags:

```go
package main

import (
 "fmt"
 "os"
 "github.com/cryptrunner49/argonaut/parser"
)

func main() {
 p := parser.NewParser(os.Stdout)
 p.SetProgramName("myapp")

 var name string
 var verbose bool
 p.StringVar(&name, "n", "", "default", "The name of the user")
 p.BoolVar(&verbose, "v", "", false, "Enable verbose output")

 err := p.Parse(os.Args[1:])
 if err != nil {
  if err == parser.ErrHelpRequested {
   return
  }
  fmt.Fprintf(os.Stderr, "Error: %v\n", err)
  os.Exit(1)
 }

 fmt.Println("Name:", name)
 if verbose {
  fmt.Println("Verbose mode enabled")
 }
}
```

Run it:

```bash
go run main.go -n Bob -v
# Output:
# Name: Bob
# Verbose mode enabled
```

### Long Flags Only

A program using only long flags:

```go
package main

import (
 "fmt"
 "os"
 "github.com/cryptrunner49/argonaut/parser"
)

func main() {
 p := parser.NewParser(os.Stdout)
 p.SetProgramName("myapp")

 var name string
 var verbose bool
 p.StringVar(&name, "", "name", "default", "The name of the user")
 p.BoolVar(&verbose, "", "verbose", false, "Enable verbose output")

 err := p.Parse(os.Args[1:])
 if err != nil {
  if err == parser.ErrHelpRequested {
   return
  }
  fmt.Fprintf(os.Stderr, "Error: %v\n", err)
  os.Exit(1)
 }

 fmt.Println("Name:", name)
 if verbose {
  fmt.Println("Verbose mode enabled")
 }
}
```

Run it:

```bash
go run main.go --name=Bob --verbose
# Output:
# Name: Bob
# Verbose mode enabled
```

### Mixed Short and Long Flags with Commands

A program mixing short and long flags with bare commands:

```go
package main

import (
 "fmt"
 "os"
 "github.com/cryptrunner49/argonaut/parser"
)

func main() {
 p := parser.NewParser(os.Stdout)
 p.SetProgramName("myapp")

 var name string
 var age int
 var debug bool
 p.StringVar(&name, "n", "name", "unknown", "User's name")
 p.IntVar(&age, "", "age", 18, "User's age")
 p.BoolVar(&debug, "d", "", false, "Enable debug mode")

 err := p.Parse(os.Args[1:])
 if err != nil {
  if err == parser.ErrHelpRequested {
   return
  }
  fmt.Fprintf(os.Stderr, "Error: %v\n", err)
  os.Exit(1)
 }

 if len(os.Args) == 1 {
  p.Usage()
  return
 }

 args := p.Args()
 if len(args) > 0 {
  switch args[0] {
  case "start":
   fmt.Printf("Starting %s (age %d), Debug: %t\n", name, age, debug)
  case "stop":
   fmt.Printf("Stopping %s (age %d), Debug: %t\n", name, age, debug)
  default:
   fmt.Fprintf(os.Stderr, "Unknown command: %s\n", args[0])
   os.Exit(1)
  }
 }
}
```

Run it with short flags and a command:

```bash
go run main.go -n Alice -d start
# Output:
# Starting Alice (age 18), Debug: true
```

Run it with long flags and a command:

```bash
go run main.go --name=Bob --age=30 stop
# Output:
# Stopping Bob (age 30), Debug: false
```

Run it with mixed flags and a command:

```bash
go run main.go -n Charlie --age=35 -d start
# Output:
# Starting Charlie (age 35), Debug: true
```

Run with `--help`:

```bash
go run main.go --help
# Output:
# Usage of myapp:
#   -n, --name string
#         User's name (default "unknown")
#   --age int
#         User's age (default 18)
#   -d
#         Enable debug mode (default false)
```

---

## API Reference

### Parser Struct

```go
type Parser struct {
    flags      []flagInfo // Stores registered flags
    positional []string   // Stores positional arguments
    output     io.Writer  // Where to write usage and errors
    program    string     // Program name for usage
}
```

Represents the command-line parser instance.

### NewParser

```go
func NewParser(out io.Writer) *Parser
```

Creates and returns a new `Parser` instance with an output writer.

- **Parameters**:
  - `out io.Writer`: Where usage and error messages are written (e.g., `os.Stdout` or `bytes.Buffer`).
- **Returns**: `*Parser` – A pointer to an initialized parser.

### SetProgramName

```go
func (p *Parser) SetProgramName(name string)
```

Sets the program name used in usage output.

- **Parameters**:
  - `name string`: The name to display in usage (e.g., `"myapp"`).

### StringVar

```go
func (p *Parser) StringVar(ptr *string, shortName, longName string, defaultValue string, description string)
```

Registers a string flag with optional short and long names.

- **Parameters**:
  - `ptr *string`: Pointer to the variable to store the flag value.
  - `shortName string`: Short flag name (e.g., `"n"`), or `""` if unused.
  - `longName string`: Long flag name (e.g., `"name"`), or `""` if unused.
  - `defaultValue string`: Default value if the flag isn’t provided.
  - `description string`: Description for usage output.

### IntVar

```go
func (p *Parser) IntVar(ptr *int, shortName, longName string, defaultValue int, description string)
```

Registers an integer flag with optional short and long names.

- **Parameters**:
  - `ptr *int`: Pointer to the variable to store the flag value.
  - `shortName string`: Short flag name, or `""` if unused.
  - `longName string`: Long flag name, or `""` if unused.
  - `defaultValue int`: Default value.
  - `description string`: Description.

### BoolVar

```go
func (p *Parser) BoolVar(ptr *bool, shortName, longName string, defaultValue bool, description string)
```

Registers a boolean flag with optional short and long names (set to `true` when present, unless explicitly `=false`).

- **Parameters**:
  - `ptr *bool`: Pointer to the variable to store the flag value.
  - `shortName string`: Short flag name, or `""` if unused.
  - `longName string`: Long flag name, or `""` if unused.
  - `defaultValue bool`: Default value.
  - `description string`: Description.

### Parse

```go
func (p *Parser) Parse(args []string) error
```

Parses the provided arguments and sets flag values. Returns an error for invalid inputs or help requests.

- **Parameters**:
  - `args []string`: Arguments to parse (e.g., `os.Args[1:]`).
- **Returns**: `error` – `nil` on success, `ErrHelpRequested` for help flags, or an error for invalid/unknown flags.

### Args

```go
func (p *Parser) Args() []string
```

Returns the list of positional arguments, including bare commands.

- **Returns**: `[]string` – Slice of non-flag arguments (e.g., `[install package]`).

### Usage

```go
func (p *Parser) Usage()
```

Prints usage information for all registered flags to the parser’s `output` writer, showing only defined short and/or long names. Can be called manually or triggered by `Parse` for `--help` or `-h`.

---

## Testing

The library is designed to be test-friendly. Use a `bytes.Buffer` for output and pass arguments directly to `Parse`. Example:

```go
package parser

import (
 "bytes"
 "testing"
)

func TestParser(t *testing.T) {
 out := &bytes.Buffer{}
 p := NewParser(out)
 p.SetProgramName("testapp")
 var name string
 p.StringVar(&name, "n", "name", "default", "User name")

 // Test parsing
 err := p.Parse([]string{"-n", "test"})
 if err != nil {
  t.Fatalf("Unexpected error: %v", err)
 }
 if name != "test" {
  t.Errorf("Expected 'test', got %s", name)
 }

 // Test manual usage
 out.Reset()
 p.Usage()
 expected := `Usage of testapp:
  -n, --name string
        User name (default "default")
`
 if out.String() != expected {
  t.Errorf("Expected usage:\n%sGot:\n%s", expected, out.String())
 }

 // Test help request
 out.Reset()
 err = p.Parse([]string{"--help"})
 if err != ErrHelpRequested {
  t.Errorf("Expected ErrHelpRequested, got %v", err)
 }
 if out.String() != expected {
  t.Errorf("Expected usage on help:\n%sGot:\n%s", expected, out.String())
 }
}
```

Run the tests:

```bash
cd parser
go test
```

The included test suite covers flag parsing, error cases, usage output, and edge cases like empty flag names and duplicates.

---

## Contributing

Contributions are welcome! To contribute:

1. Fork the repository.
2. Create a feature branch (`git checkout -b feature/new-feature`).
3. Commit your changes (`git commit -m "Add new feature"`).
4. Push to the branch (`git push origin feature/new-feature`).
5. Open a pull request.

Please include tests for new features or bug fixes.

---

## License

This library is licensed under the MIT License. See the `LICENSE` file for details.
