package composer

import (
	"io"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestParseString(t *testing.T) {
	validJSON := `{
		"name": "vendor/project",
		"description": "Test project",
		"type": "library",
		"license": "MIT",
		"require": {
			"php": "^7.4",
			"ext-json": "*"
		},
		"require-dev": {
			"phpunit/phpunit": "^9.0"
		},
		"autoload": {
			"psr-4": {
				"Vendor\\Project\\": "src/"
			}
		}
	}`

	composer, err := ParseString(validJSON)
	if err != nil {
		t.Fatalf("Expected successful parsing but got error: %v", err)
	}

	// Check basic properties
	if composer.Name != "vendor/project" {
		t.Errorf("Expected name 'vendor/project', got '%s'", composer.Name)
	}
	if composer.Description != "Test project" {
		t.Errorf("Expected description 'Test project', got '%s'", composer.Description)
	}
	if composer.Type != "library" {
		t.Errorf("Expected type 'library', got '%s'", composer.Type)
	}
	if composer.License != "MIT" {
		t.Errorf("Expected license 'MIT', got '%v'", composer.License)
	}

	// Check dependencies
	if composer.Require["php"] != "^7.4" {
		t.Errorf("Expected PHP requirement '^7.4', got '%s'", composer.Require["php"])
	}
	if composer.Require["ext-json"] != "*" {
		t.Errorf("Expected ext-json requirement '*', got '%s'", composer.Require["ext-json"])
	}
	if composer.RequireDev["phpunit/phpunit"] != "^9.0" {
		t.Errorf("Expected phpunit requirement '^9.0', got '%s'", composer.RequireDev["phpunit/phpunit"])
	}

	// Check PSR-4 autoloading
	psr4Map, ok := composer.Autoload.PSR4.(map[string]interface{})
	if !ok {
		t.Fatalf("Expected PSR-4 autoload to be a map, got %T", composer.Autoload.PSR4)
	}
	if psr4Map["Vendor\\Project\\"] != "src/" {
		t.Errorf("Expected PSR-4 namespace 'Vendor\\Project\\' to map to 'src/', got '%v'", psr4Map["Vendor\\Project\\"])
	}

	// Test invalid JSON
	invalidJSON := `{ "name": "invalid" `
	_, err = ParseString(invalidJSON)
	if err == nil {
		t.Errorf("Expected error for invalid JSON, but got none")
	}
}

func TestParse(t *testing.T) {
	validJSON := `{"name": "vendor/project"}`
	reader := strings.NewReader(validJSON)

	composer, err := Parse(reader)
	if err != nil {
		t.Fatalf("Expected successful parsing but got error: %v", err)
	}
	if composer.Name != "vendor/project" {
		t.Errorf("Expected name 'vendor/project', got '%s'", composer.Name)
	}

	// Test error case
	errorReader := &errorReader{}
	_, err = Parse(errorReader)
	if err == nil {
		t.Errorf("Expected error for reader failure, but got none")
	}
}

func TestParseFile(t *testing.T) {
	// Create a temporary directory
	tmpDir, err := os.MkdirTemp("", "composer-test")
	if err != nil {
		t.Fatalf("Failed to create temp directory: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	// Create a test composer.json file
	testFilePath := filepath.Join(tmpDir, "composer.json")
	testJSON := `{"name": "vendor/project"}`
	if err := os.WriteFile(testFilePath, []byte(testJSON), 0644); err != nil {
		t.Fatalf("Failed to write test file: %v", err)
	}

	// Test parsing the file
	composer, err := ParseFile(testFilePath)
	if err != nil {
		t.Fatalf("Expected successful file parsing but got error: %v", err)
	}
	if composer.Name != "vendor/project" {
		t.Errorf("Expected name 'vendor/project', got '%s'", composer.Name)
	}

	// Test parsing non-existent file
	nonExistentPath := filepath.Join(tmpDir, "nonexistent.json")
	_, err = ParseFile(nonExistentPath)
	if err == nil {
		t.Errorf("Expected error for non-existent file, but got none")
	}
}

func TestParseDir(t *testing.T) {
	// Create a temporary directory
	tmpDir, err := os.MkdirTemp("", "composer-test")
	if err != nil {
		t.Fatalf("Failed to create temp directory: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	// Create a test composer.json file
	testFilePath := filepath.Join(tmpDir, "composer.json")
	testJSON := `{"name": "vendor/project"}`
	if err := os.WriteFile(testFilePath, []byte(testJSON), 0644); err != nil {
		t.Fatalf("Failed to write test file: %v", err)
	}

	// Test parsing from directory
	composer, err := ParseDir(tmpDir)
	if err != nil {
		t.Fatalf("Expected successful directory parsing but got error: %v", err)
	}
	if composer.Name != "vendor/project" {
		t.Errorf("Expected name 'vendor/project', got '%s'", composer.Name)
	}

	// Test parsing from directory without composer.json
	emptyDir, err := os.MkdirTemp("", "composer-empty")
	if err != nil {
		t.Fatalf("Failed to create empty temp directory: %v", err)
	}
	defer os.RemoveAll(emptyDir)

	_, err = ParseDir(emptyDir)
	if err == nil {
		t.Errorf("Expected error for directory without composer.json, but got none")
	}
}

// Helper error reader for testing error cases
type errorReader struct{}

func (e *errorReader) Read(p []byte) (n int, err error) {
	return 0, io.ErrUnexpectedEOF
}
