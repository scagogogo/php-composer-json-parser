// Package validation provides functionality for validating PHP Composer JSON structures
package validation

import (
	"fmt"
	"regexp"

	"github.com/scagogogo/php-composer-json-parser/pkg/composer/dependency"
)

// ValidateComposerJSON validates a ComposerJSON structure
func ValidateComposerJSON(name, description, stability string) error {
	// Validate package name if provided
	if name != "" {
		if err := dependency.ValidatePackageName(name); err != nil {
			return err
		}
	}

	// Validate description if provided
	if description != "" && len(description) < 10 {
		return fmt.Errorf("description is too short, should be at least 10 characters")
	}

	// Validate minimum-stability if provided
	if stability != "" {
		validStabilities := []string{"dev", "alpha", "beta", "RC", "stable"}
		valid := false
		for _, s := range validStabilities {
			if stability == s {
				valid = true
				break
			}
		}
		if !valid {
			return fmt.Errorf("invalid minimum-stability '%s', should be one of: dev, alpha, beta, RC, stable", stability)
		}
	}

	return nil
}

// ValidateVersion validates a version string
func ValidateVersion(version string) error {
	if version == "" {
		return nil // Empty version is valid (omitted)
	}

	// Allow "*" as a wildcard
	if version == "*" {
		return nil
	}

	// Allow "dev-master", "dev-*", or specific branch
	if regexp.MustCompile(`^dev-`).MatchString(version) {
		return nil
	}

	// Check for valid semantic version with constraints
	// Example patterns: "1.0.0", "^1.0", "~1.2.3", ">=1.0", ">1.0 <2.0", etc.
	versionRegex := regexp.MustCompile(`^(\^|~|>=|<=|>|<|!=|==)?[0-9]+(\.[0-9]+)?(\.[0-9]+)?(\-[0-9A-Za-z-]+(\.[0-9A-Za-z-]+)*)?(\+[0-9A-Za-z-]+(\.[0-9A-Za-z-]+)*)?$`)
	if versionRegex.MatchString(version) {
		return nil
	}

	// Check for version ranges
	rangeRegex := regexp.MustCompile(`^([><=!]?[0-9]+(\.[0-9]+)?(\.[0-9]+)?)(\s+[><=!]?[0-9]+(\.[0-9]+)?(\.[0-9]+)?)$`)
	if rangeRegex.MatchString(version) {
		return nil
	}

	return fmt.Errorf("invalid version format: %s", version)
}
