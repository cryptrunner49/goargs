# 🏹 ArGonaut 🚀  

**ArGonaut** is an easy-to-use CLI argument parser for Go that automatically generates documentation for every command and flag, ensuring your docs stay up-to-date with zero effort! 🚀📚

## Features ✨

- **Auto-Documentation:** Generates docs for flags automatically via usage output.
- **Simple & Intuitive:** Easily define CLI arguments with minimal setup.
- **Test-Friendly:** No process exits, capturable output, and error-based API.
- **Flexible & Powerful:** Supports short/long flags (optional), positional arguments, and manual usage display.
- **Go-Powered:** Lightweight, efficient, and designed for performance.

## Quick Start 🚀

1. **Install via Go Modules:**

   ```bash
   go get github.com/cryptrunner49/argonaut
   ```

2. **Basic Usage Example:**

```go
package main

import (
    "fmt"
    "os"
    "github.com/cryptrunner49/argonaut/parser"
)

func main() {
    // Create a new parser instance with output writer
    p := parser.NewParser(os.Stdout)
    p.SetProgramName("myapp") // Optional: customize program name for usage

    // Define variables for flags
    var name string
    var age int
    var verbose bool

    // Register flags with short/long names (either can be "")
    p.StringVar(&name, "n", "name", "default", "The name of the user")
    p.IntVar(&age, "a", "age", 0, "The age of the user")
    p.BoolVar(&verbose, "v", "verbose", false, "Enable verbose output")

    // Parse arguments and handle errors
    if err := p.Parse(os.Args[1:]); err != nil {
        if err == parser.ErrHelpRequested {
            return // Help was printed, exit gracefully
        }
        fmt.Fprintf(os.Stderr, "Error: %v\n", err)
        os.Exit(1)
    }

    // Manually show usage if no args provided (optional)
    if len(os.Args) == 1 {
        p.Usage()
        return
    }

    // Use parsed values
    fmt.Println("Name:", name)
    fmt.Println("Age:", age)
    fmt.Println("Verbose:", verbose)
    fmt.Println("Args:", p.Args())
}
```

## API Overview 📖

### Creating a Parser

```go
p := parser.NewParser(out io.Writer) // out can be os.Stdout or a bytes.Buffer for testing
p.SetProgramName("myapp")           // Optional: sets program name for usage
```

### Registering Flags

```go
var s string
var i int
var b bool
p.StringVar(&s, "s", "string", "def", "A string flag") // Both short and long
p.IntVar(&i, "i", "", 0, "An integer flag")            // Only short
p.BoolVar(&b, "", "bool", false, "A boolean flag")     // Only long
```

### Parsing Arguments

```go
err := p.Parse(args []string) // Returns error instead of exiting
if err == parser.ErrHelpRequested {
    // Handle help request (usage was printed)
}
```

### Showing Usage Manually

```go
p.Usage() // Prints usage to the configured io.Writer
```

### Getting Positional Arguments

```go
args := p.Args() // Returns []string of positional arguments
```

## Testing 🧪

The parser is designed to be test-friendly:

- Use a `bytes.Buffer` to capture usage output.
- Pass argument slices directly to `Parse`.
- Check returned errors instead of handling process exits.
- No global state dependencies.

Example test:

```go
func TestExample(t *testing.T) {
    out := &bytes.Buffer{}
    p := parser.NewParser(out)
    var flag string
    p.StringVar(&flag, "f", "flag", "", "Test flag")

    // Test parsing
    err := p.Parse([]string{"--flag=test"})
    if err != nil {
        t.Fatal(err)
    }
    if flag != "test" {
        t.Errorf("expected 'test', got %q", flag)
    }

    // Test usage output
    out.Reset()
    p.Usage()
    expected := `Usage of program:
  -f, --flag string
        Test flag (default "")
`
    if out.String() != expected {
        t.Errorf("expected usage:\n%s\ngot:\n%s", expected, out.String())
    }
}
```

## Notes 📝

- **Flag Names**: Either short (`-f`) or long (`--flag`) names can be omitted by setting them to `""`. Duplicate flag names are allowed, with the last registered flag taking precedence.
- **Bare Commands**: Positional arguments can be used as bare commands (e.g., `myapp start`), handled manually via `p.Args()`.

## Contributing 🤝

We welcome contributions! Check out our [CONTRIBUTING.md](CONTRIBUTING.md) for guidelines on how to help improve ArGonaut.

## License 📜

This project is licensed under the MIT License. See the [LICENSE](LICENSE) file for details.
