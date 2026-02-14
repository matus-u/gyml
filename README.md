# ⚙️ gyml: Generic YAML Manipulator

![GitHub go.mod Go version](https://img.shields.io/github/go-mod/go-version/matus-u/gyml)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)

`gyml` is a lightweight Go library designed for intuitive and flexible manipulation of YAML files using path-based lookups. Say goodbye to cumbersome struct definitions for every YAML configuration – `gyml` allows you to access, set, and delete values within any YAML document without needing to know its full structure beforehand.

## ✨ Features

- **Path-Based Access:** Easily target specific YAML nodes using a simple key-based path (e.g., `["database", "connections", "[0]", "port"]`).
- **Dynamic Type Handling:** Automatically handles different YAML node types (maps, sequences, scalars) during read and write operations.
- **Value Setting & Appending:** Set new values or append items to sequences, with automatic creation of intermediate nodes if they don't exist.
- **Value Deletion:** Remove specific nodes or entire sub-trees from your YAML document.
- **Generic Data Types:** Work with any Go data type that can be marshaled/unmarshaled by `gopkg.in/yaml.v3`.
- **Error Handling:** Comprehensive error reporting for invalid paths, out-of-bound indices, and unsupported operations.

## 🚀 Installation

To integrate `gyml` into your Go project, simply run the following command:

```bash
go get github.com/matus-u/gyml
```

## 📚 Usage

Here's how you can use `gyml` to manipulate your YAML data:

### Initializing a YAML Document

You can start with an empty YAML node or parse an existing YAML string:

```go
package main

import (
	"fmt"
	"log"

	"gopkg.in/yaml.v3"
	"github.com/matus-u/gyml"
)

func main() {
	// Start with an empty root node
	var root yaml.Node
	fmt.Println("Initial YAML:", toYAMLString(&root))

	// Or parse an existing YAML string
	yamlString := `
application:
  name: MyWebApp
  version: 1.0.0
database:
  host: localhost
  port: 5432
  users:
    - name: admin
      role: superuser
    - name: guest
      role: readonly
`
	var parsedRoot yaml.Node
	err := yaml.Unmarshal([]byte(yamlString), &parsedRoot)
	if err != nil {
		log.Fatalf("Error unmarshaling YAML: %v", err)
	}
	fmt.Println("Parsed YAML:", toYAMLString(&parsedRoot))
}

// Helper to convert yaml.Node back to string for printing
func toYAMLString(node *yaml.Node) string {
	if len(node.Content) == 0 {
		return "{}" // Represent empty YAML as an empty map
	}
	data, err := yaml.Marshal(node)
	if err != nil {
		return fmt.Sprintf("Error marshaling: %v", err)
	}
	return string(data)
}
```

### Setting Values

Use `gyml.SetValue` to add or modify values. Intermediate nodes will be created as needed.

```go
package main

import (
	"fmt"
	"log"

	"gopkg.in/yaml.v3"
	"github.com/matus-u/gyml"
)

func main() {
	var root yaml.Node

	// Set a simple scalar value
	err := gyml.SetValue(&root, "My Application", "application", "name")
	if err != nil {
		log.Fatalf("SetValue error: %v", err)
	}
	fmt.Println("YAML after setting application name:", toYAMLString(&root))

	// Set a nested struct
	type DBConfig struct {
		Host string `yaml:"host"`
		Port int    `yaml:"port"`
	}
	err = gyml.SetValue(&root, DBConfig{Host: "db.example.com", Port: 8080}, "application", "database")
	if err != nil {
		log.Fatalf("SetValue error: %v", err)
	}
	fmt.Println("YAML after setting database config:", toYAMLString(&root))

	// Append to a list
	err = gyml.SetValue(&root, "user1", "application", "users", "[]")
	if err != nil {
		log.Fatalf("SetValue error: %v", err)
	}
	err = gyml.SetValue(&root, "user2", "application", "users", "[]")
	if err != nil {
		log.Fatalf("SetValue error: %v", err)
	}
	fmt.Println("YAML after appending users:", toYAMLString(&root))

	// Set a value at a specific index in a list
	err = gyml.SetValue(&root, "user_one_updated", "application", "users", "[0]")
	if err != nil {
		log.Fatalf("SetValue error: %v", err)
	}
	fmt.Println("YAML after updating user at index 0:", toYAMLString(&root))

	// Set a boolean value
	err = gyml.SetValue(&root, true, "application", "settings", "notifications")
	if err != nil {
		log.Fatalf("SetValue error: %v", err)
	}
	fmt.Println("YAML after setting notifications:", toYAMLString(&root))
}

// Helper to convert yaml.Node back to string for printing
func toYAMLString(node *yaml.Node) string {
	if len(node.Content) == 0 {
		return "{}" // Represent empty YAML as an empty map
	}
	data, err := yaml.Marshal(node)
	if err != nil {
		return fmt.Sprintf("Error marshaling: %v", err)
	}
	return string(data)
}
```

### Getting Values

Use `gyml.GetValue` to retrieve values from your YAML document and deserialize them into Go types.

```go
package main

import (
	"fmt"
	"log"

	"gopkg.in/yaml.v3"
	"github.com/matus-u/gyml"
)

func main() {
	yamlString := `
application:
  name: MyWebApp
  version: 1.0.0
database:
  host: localhost
  port: 5432
  users:
    - name: admin
      role: superuser
    - name: guest
      role: readonly
settings:
  notifications: true
`
	var root yaml.Node
	err := yaml.Unmarshal([]byte(yamlString), &root)
	if err != nil {
		log.Fatalf("Error unmarshaling YAML: %v", err)
	}

	// Get a string value
	appName, err := gyml.GetValue[string](&root, "application", "name")
	if err != nil {
		log.Fatalf("GetValue error for app name: %v", err)
	}
	fmt.Printf("Application Name: %s
", *appName)

	// Get an int value
	dbPort, err := gyml.GetValue[int](&root, "database", "port")
	if err != nil {
		log.Fatalf("GetValue error for database port: %v", err)
	}
	fmt.Printf("Database Port: %d
", *dbPort)

	// Get a boolean value
	notificationsEnabled, err := gyml.GetValue[bool](&root, "settings", "notifications")
	if err != nil {
		log.Fatalf("GetValue error for notifications: %v", err)
	}
	fmt.Printf("Notifications Enabled: %t
", *notificationsEnabled)

	// Get a struct
	type User struct {
		Name string `yaml:"name"`
		Role string `yaml:"role"`
	}
	firstUser, err := gyml.GetValue[User](&root, "database", "users", "[0]")
	if err != nil {
		log.Fatalf("GetValue error for first user: %v", err)
	}
	fmt.Printf("First User: %+v
", *firstUser)

	// Get a slice of strings
	type Application struct {
		Users []string `yaml:"users"`
	}
	// For example, if you set users as just strings.
	// You might need to adjust based on how you set it.
	// usersList, err := gyml.GetValue[[]string](&root, "database", "users")
	// if err != nil {
	// 	log.Fatalf("GetValue error for users list: %v", err)
	// }
	// fmt.Printf("Users List: %v
", *usersList)
}
```

### Deleting Values

Use `gyml.DeleteValue` to remove nodes from your YAML document.

```go
package main

import (
	"fmt"
	"log"

	"gopkg.in/yaml.v3"
	"github.com/matus-u/gyml"
)

func main() {
	yamlString := `
application:
  name: MyWebApp
  version: 1.0.0
database:
  host: localhost
  port: 5432
  users:
    - name: admin
      role: superuser
    - name: guest
      role: readonly
`
	var root yaml.Node
	err := yaml.Unmarshal([]byte(yamlString), &root)
	if err != nil {
		log.Fatalf("Error unmarshaling YAML: %v", err)
	}
	fmt.Println("Original YAML:", toYAMLString(&root))

	// Delete a specific key
	err = gyml.DeleteValue(&root, "application", "version")
	if err != nil {
		log.Fatalf("DeleteValue error: %v", err)
	}
	fmt.Println("YAML after deleting application version:", toYAMLString(&root))

	// Delete an item from a sequence by index
	err = gyml.DeleteValue(&root, "database", "users", "[0]")
	if err != nil {
		log.Fatalf("DeleteValue error: %v", err)
	}
	fmt.Println("YAML after deleting first user:", toYAMLString(&root))

	// Delete an entire sub-map
	err = gyml.DeleteValue(&root, "database", "host") // Delete host only
	if err != nil {
		log.Fatalf("DeleteValue error: %v", err)
	}
	fmt.Println("YAML after deleting database host:", toYAMLString(&root))

	err = gyml.DeleteValue(&root, "database", "port") // Delete port, now database map is empty
	if err != nil {
		log.Fatalf("DeleteValue error: %v", err)
	}
	fmt.Println("YAML after deleting database port (database should be gone if empty):", toYAMLString(&root))
}

// Helper to convert yaml.Node back to string for printing
func toYAMLString(node *yaml.Node) string {
	if len(node.Content) == 0 {
		return "{}" // Represent empty YAML as an empty map
	}
	data, err := yaml.Marshal(node)
	if err != nil {
		return fmt.Sprintf("Error marshaling: %v", err)
	}
	return string(data)
}
```

## 🤝 Contributing

Contributions are welcome! Please feel free to open issues or submit pull requests.

## 📄 License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.
