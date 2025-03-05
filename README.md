# GoArgonaut 🏹🚀  

**GoArgonaut** is an easy-to-use CLI argument parser for Go that automatically generates documentation for every command and flag—keeping your docs up-to-date effortlessly!  

## Features ✨

- **Auto-Documentation:** Generates docs for commands & flags automatically.
- **Simple & Intuitive:** Easily define CLI arguments with minimal setup.
- **Flexible & Powerful:** Supports nested commands and custom flag types.
- **Go-Powered:** Lightweight, efficient, and designed for performance.

## Installation 🛠️  

Install GoArgonaut via `go get`:  

```bash
go get github.com/cryptrunner49/goargonaut
```

## Quick Start 🚀  

Here's a simple example:  

```go
package main

import (
    "fmt"
    "github.com/cryptrunner49/goargonaut"
)

func main() {
    // Create a new parser instance
    parser := goargonaut.NewParser()

    // Define flags and commands
    parser.AddFlag("verbose", "Enable verbose output", false)
    parser.AddCommand("start", "Start the service", func(args []string) {
        fmt.Println("Service started!")
    })

    // Parse arguments and auto-generate documentation
    if err := parser.Parse(); err != nil {
        fmt.Println("Error:", err)
    }
}
```

## Documentation 📖

After running your application, find the auto-generated docs in the `docs/` folder. They update automatically as you modify your commands and flags!  

## Contributing 🤝

We welcome contributions! Check out our [CONTRIBUTING.md](CONTRIBUTING.md) for guidelines on how to help improve GoArgonaut.

## License 📄

This project is licensed under the MIT License. See the [LICENSE](LICENSE) file for details.
