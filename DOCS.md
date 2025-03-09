# Argonaut Documentation

## Argonaut Parser Library

The `Argonaut Parser` is a simple, lightweight command-line argument parsing library for Go. It allows developers to define and parse flags (string, int, and bool) and handle positional arguments with ease. This library is ideal for small to medium-sized CLI applications where you need straightforward flag parsing without additional dependencies.

This document provides a comprehensive guide to using the library, including installation instructions, usage examples, and a full API reference.

---

## Table of Contents

1. [Installation](#installation)
2. [Quick Start](#quick-start)
3. [Usage](#usage)
   - [Defining Flags](#defining-flags)
   - [Parsing Arguments](#parsing-arguments)
   - [Accessing Positional Arguments](#accessing-positional-arguments)
   - [Displaying Usage Information](#displaying-usage-information)
4. [Examples](#examples)
   - [Basic Example](#basic-example)
   - [Advanced Example](#advanced-example)
5. [API Reference](#api-reference)
   - [Parser Struct](#parser-struct)
   - [NewParser](#newparser)
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

Here’s a minimal example to get you started:

```go
package main

import (
 "fmt"
 "github.com/cryptrunner49/argonaut/parser"
)

func main() {
 p := parser.NewParser()

 var name string
 p.StringVar(&name, "name", "guest", "The name of the user")

 p.Parse()

 fmt.Println("Hello,", name)
}
```

Run it:

```bash
go run main.go --name=Alice
# Output: Hello, Alice
```

---

## Usage

### Defining Flags

Flags are defined using methods like `StringVar`, `IntVar`, and `BoolVar`. Each method registers a flag with a name, default value, and description.

- **Syntax**: `--flag=value` (for string and int) or `--flag` (for bool, sets to `true`).
- **Default Values**: If a flag isn’t provided, it retains its default value.

### Parsing Arguments

Call `Parse()` to process command-line arguments (`os.Args`). It sets registered flag values and collects positional arguments.

### Accessing Positional Arguments

Use `Args()` to retrieve a slice of positional arguments (non-flag arguments).

### Displaying Usage Information

The library automatically handles the `--help` flag, displaying usage information for all registered flags.

---

## Examples

### Basic Example

A simple program with a string and bool flag:

```go
package main

import (
 "fmt"
 "github.com/cryptrunner49/argonaut/parser"
)

func main() {
 p := parser.NewParser()

 var name string
 var verbose bool
 p.StringVar(&name, "name", "default", "The name of the user")
 p.BoolVar(&verbose, "verbose", false, "Enable verbose output")

 p.Parse()

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

### Advanced Example

A program with multiple flag types and positional arguments:

```go
package main

import (
 "fmt"
 "github.com/cryptrunner49/argonaut/parser"
)

func main() {
 p := parser.NewParser()

 var name string
 var age int
 var debug bool
 p.StringVar(&name, "name", "unknown", "User's name")
 p.IntVar(&age, "age", 18, "User's age")
 p.BoolVar(&debug, "debug", false, "Enable debug mode")

 p.Parse()

 fmt.Printf("Name: %s, Age: %d, Debug: %t\n", name, age, debug)
 fmt.Println("Files:", p.Args())
}
```

Run it:

```bash
go run main.go --name=Alice --age=25 --debug file1.txt file2.txt
# Output:
# Name: Alice, Age: 25, Debug: true
# Files: [file1.txt file2.txt]
```

Run with `--help`:

```bash
go run main.go --help
# Output:
# Usage of ./main:
#   --name string
#         User's name (default "unknown")
#   --age int
#         User's age (default 18)
#   --debug
#         Enable debug mode (default false)
```

---

## API Reference

### Parser Struct

```go
type Parser struct {
    flags      []flagInfo // Stores registered flags
    positional []string   // Stores positional arguments
}
```

Represents the command-line parser instance.

### NewParser

```go
func NewParser() *Parser
```

Creates and returns a new `Parser` instance.

- **Returns**: `*Parser` – A pointer to an initialized parser.

### StringVar

```go
func (p *Parser) StringVar(ptr *string, name string, defaultValue string, description string)
```

Registers a string flag.

- **Parameters**:
  - `ptr *string`: Pointer to the variable to store the flag value.
  - `name string`: Flag name (e.g., `"name"` for `--name`).
  - `defaultValue string`: Default value if the flag isn’t provided.
  - `description string`: Description for `--help` output.

### IntVar

```go
func (p *Parser) IntVar(ptr *int, name string, defaultValue int, description string)
```

Registers an integer flag.

- **Parameters**:
  - `ptr *int`: Pointer to the variable to store the flag value.
  - `name string`: Flag name.
  - `defaultValue int`: Default value.
  - `description string`: Description.

### BoolVar

```go
func (p *Parser) BoolVar(ptr *bool, name string, defaultValue bool, description string)
```

Registers a boolean flag (set to `true` when present, unless explicitly `=false`).

- **Parameters**:
  - `ptr *bool`: Pointer to the variable to store the flag value.
  - `name string`: Flag name.
  - `defaultValue bool`: Default value.
  - `description string`: Description.

### Parse

```go
func (p *Parser) Parse()
```

Parses command-line arguments from `os.Args` and sets flag values. Exits with usage info if `--help` is provided or with an error for invalid/unknown flags.

### Args

```go
func (p *Parser) Args() []string
```

Returns the list of positional arguments.

- **Returns**: `[]string` – Slice of non-flag arguments.

### Usage

```go
func (p *Parser) Usage()
```

Prints usage information for all registered flags. Called automatically by `Parse()` when `--help` is detected.

---

## Testing

The library includes a test suite in `parser/parser_test.go`. To run the tests:

```bash
cd parser
go test
```

The tests cover:

- Parsing string, int, and bool flags.
- Handling positional arguments.
- Default value behavior.

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
