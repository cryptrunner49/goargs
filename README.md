# ArGonaut 🏹🚀  

**ArGonaut** is an easy-to-use CLI argument parser for Go that automatically generates documentation for every command and flag, ensuring your docs stay up-to-date with zero effort! 🚀📚

## Features ✨

- **Auto-Documentation:** Generates docs for commands & flags automatically.
- **Simple & Intuitive:** Easily define CLI arguments with minimal setup.
- **Flexible & Powerful:** Supports nested commands and custom flag types.
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
    "github.com/cryptrunner49/argonaut"
)

func main() {
    // Create a new parser instance
    parser := argonaut.NewParser()

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

## Contributing 🤝

We welcome contributions! Check out our [CONTRIBUTING.md](CONTRIBUTING.md) for guidelines on how to help improve ArGonaut.

## License 📜

This project is licensed under the MIT License. See the [LICENSE](LICENSE) file for details.
