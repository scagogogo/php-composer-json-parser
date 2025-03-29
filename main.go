package main

import (
	"fmt"
	"log"
	"os"

	"github.com/scagogogo/php-composer-json-parser/pkg/composer"
)

func main() {
	// Check command-line arguments
	if len(os.Args) > 1 {
		// If a file path is provided, parse that file
		composerFile := os.Args[1]
		fmt.Printf("Parsing composer.json file: %s\n", composerFile)

		c, err := composer.ParseFile(composerFile)
		if err != nil {
			log.Fatalf("Error parsing composer.json: %v", err)
		}

		displayComposerInfo(c)
		return
	}

	// Otherwise, create a new composer.json example
	fmt.Println("Creating a new composer.json example")

	c, err := composer.CreateNew("example/project", "A PHP Composer Parser Example")
	if err != nil {
		log.Fatalf("Error creating composer.json: %v", err)
	}

	// Add some dependencies
	c.AddDependency("php", "^8.0")
	c.AddDependency("symfony/console", "^6.0")
	c.AddDependency("monolog/monolog", "^2.3")
	c.AddDevDependency("phpunit/phpunit", "^9.5")

	// Set up PSR-4 autoloading
	psr4Map := make(map[string]interface{})
	psr4Map["Example\\Project\\"] = "src/"
	psr4Map["Example\\Tests\\"] = "tests/"
	c.Autoload.PSR4 = psr4Map

	// Display the created composer.json
	displayComposerInfo(c)

	// Convert to JSON and print
	jsonStr, err := c.ToJSON(true)
	if err != nil {
		log.Fatalf("Error converting to JSON: %v", err)
	}

	fmt.Println("\nGenerated composer.json:")
	fmt.Println("------------------------")
	fmt.Println(jsonStr)
}

func displayComposerInfo(c *composer.ComposerJSON) {
	fmt.Println("\nComposer.json Information:")
	fmt.Println("-------------------------")
	fmt.Printf("Package Name: %s\n", c.Name)
	fmt.Printf("Description: %s\n", c.Description)
	fmt.Printf("Type: %s\n", c.Type)
	fmt.Printf("License: %v\n", c.License)

	fmt.Println("\nDependencies:")
	for pkg, version := range c.Require {
		fmt.Printf("  - %s: %s\n", pkg, version)
	}

	fmt.Println("\nDev Dependencies:")
	for pkg, version := range c.RequireDev {
		fmt.Printf("  - %s: %s\n", pkg, version)
	}

	// Display autoloading info if present
	if psr4, ok := c.Autoload.PSR4.(map[string]interface{}); ok && len(psr4) > 0 {
		fmt.Println("\nPSR-4 Autoloading:")
		for namespace, path := range psr4 {
			fmt.Printf("  - %s => %s\n", namespace, path)
		}
	}
}
