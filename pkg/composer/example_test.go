package composer

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
)

func Example() {
	// Create a new empty composer.json
	composer, err := CreateNew("vendor/project", "Example PHP project")
	if err != nil {
		log.Fatalf("Failed to create composer.json: %v", err)
	}

	// Add dependencies
	composer.AddDependency("php", "^8.0")
	composer.AddDependency("monolog/monolog", "^2.0")
	composer.AddDevDependency("phpunit/phpunit", "^9.0")

	// Set up PSR-4 autoloading
	psr4Map := make(map[string]interface{})
	psr4Map["Vendor\\Project\\"] = "src/"
	composer.Autoload.PSR4 = psr4Map

	// Convert to JSON and print
	jsonStr, _ := composer.ToJSON(true)
	fmt.Println(jsonStr)

	// Output contains:
	// "name": "vendor/project",
	// "description": "Example PHP project",
}

func Example_parseExisting() {
	// Sample composer.json content
	jsonStr := `{
		"name": "vendor/existing-project",
		"require": {
			"php": "^7.4"
		}
	}`

	// Parse from string
	composer, err := ParseString(jsonStr)
	if err != nil {
		log.Fatalf("Failed to parse composer.json: %v", err)
	}

	// Check if a dependency exists
	if composer.DependencyExists("php") {
		fmt.Println("PHP dependency exists with version:", composer.Require["php"])
	}

	// Add more dependencies
	composer.AddDependency("symfony/console", "^5.0")

	// Output:
	// PHP dependency exists with version: ^7.4
}

func Example_saveAndLoad() {
	// Create a temporary directory for the example
	tmpDir, err := os.MkdirTemp("", "composer-example")
	if err != nil {
		log.Fatalf("Failed to create temp directory: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	// Create a new composer.json
	composer, _ := CreateNew("vendor/save-example", "Example for saving and loading")
	composer.AddDependency("php", "^8.0")

	// Save to file
	filePath := filepath.Join(tmpDir, "composer.json")
	if err := composer.Save(filePath, true); err != nil {
		log.Fatalf("Failed to save composer.json: %v", err)
	}
	fmt.Println("Saved composer.json")

	// Load from file
	loadedComposer, err := ParseFile(filePath)
	if err != nil {
		log.Fatalf("Failed to load composer.json: %v", err)
	}
	fmt.Println("Loaded composer.json for:", loadedComposer.Name)

	// Output:
	// Saved composer.json
	// Loaded composer.json for: vendor/save-example
}
