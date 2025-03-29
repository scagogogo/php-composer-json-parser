package validation

import (
	"strings"
	"testing"
)

func TestValidateComposerJSON(t *testing.T) {
	tests := []struct {
		name          string
		packageName   string
		description   string
		stability     string
		wantErr       bool
		errorContains string
	}{
		{
			name:        "All valid",
			packageName: "vendor/package",
			description: "This is a valid description",
			stability:   "stable",
			wantErr:     false,
		},
		{
			name:        "Empty values",
			packageName: "",
			description: "",
			stability:   "",
			wantErr:     false,
		},
		{
			name:          "Invalid package name",
			packageName:   "invalid-package-name",
			description:   "Valid description",
			stability:     "stable",
			wantErr:       true,
			errorContains: "包名必须符合",
		},
		{
			name:          "Description too short",
			packageName:   "vendor/package",
			description:   "Too short",
			stability:     "stable",
			wantErr:       true,
			errorContains: "description is too short",
		},
		{
			name:          "Invalid stability",
			packageName:   "vendor/package",
			description:   "This is a valid description",
			stability:     "invalid",
			wantErr:       true,
			errorContains: "invalid minimum-stability",
		},
		{
			name:        "Valid stability - dev",
			packageName: "vendor/package",
			description: "This is a valid description",
			stability:   "dev",
			wantErr:     false,
		},
		{
			name:        "Valid stability - alpha",
			packageName: "vendor/package",
			description: "This is a valid description",
			stability:   "alpha",
			wantErr:     false,
		},
		{
			name:        "Valid stability - beta",
			packageName: "vendor/package",
			description: "This is a valid description",
			stability:   "beta",
			wantErr:     false,
		},
		{
			name:        "Valid stability - RC",
			packageName: "vendor/package",
			description: "This is a valid description",
			stability:   "RC",
			wantErr:     false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateComposerJSON(tt.packageName, tt.description, tt.stability)

			// 检查错误情况
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateComposerJSON() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if err != nil && tt.errorContains != "" && !strings.Contains(err.Error(), tt.errorContains) {
				t.Errorf("ValidateComposerJSON() error = %v, should contain %v", err, tt.errorContains)
			}
		})
	}
}

func TestValidateVersion(t *testing.T) {
	tests := []struct {
		name          string
		version       string
		wantErr       bool
		errorContains string
	}{
		{
			name:    "Empty version",
			version: "",
			wantErr: false,
		},
		{
			name:    "Wildcard",
			version: "*",
			wantErr: false,
		},
		{
			name:    "Dev branch",
			version: "dev-master",
			wantErr: false,
		},
		{
			name:    "Dev feature branch",
			version: "dev-feature-branch",
			wantErr: false,
		},
		{
			name:    "Simple version",
			version: "1.0.0",
			wantErr: false,
		},
		{
			name:    "Version with caret",
			version: "^1.0.0",
			wantErr: false,
		},
		{
			name:    "Version with tilde",
			version: "~1.0.0",
			wantErr: false,
		},
		{
			name:    "Version with greater than",
			version: ">1.0.0",
			wantErr: false,
		},
		{
			name:    "Version with less than",
			version: "<1.0.0",
			wantErr: false,
		},
		{
			name:    "Version with greater than or equal",
			version: ">=1.0.0",
			wantErr: false,
		},
		{
			name:    "Version with less than or equal",
			version: "<=1.0.0",
			wantErr: false,
		},
		{
			name:    "Version with not equal",
			version: "!=1.0.0",
			wantErr: false,
		},
		{
			name:    "Version with equal",
			version: "==1.0.0",
			wantErr: false,
		},
		{
			name:    "Version with pre-release",
			version: "1.0.0-alpha",
			wantErr: false,
		},
		{
			name:    "Version with build metadata",
			version: "1.0.0+build.1",
			wantErr: false,
		},
		{
			name:    "Version with pre-release and build metadata",
			version: "1.0.0-alpha+build.1",
			wantErr: false,
		},
		{
			name:    "Range with space",
			version: ">1.0.0 <2.0.0",
			wantErr: false,
		},
		{
			name:          "Invalid version format",
			version:       "invalid",
			wantErr:       true,
			errorContains: "invalid version format",
		},
		{
			name:          "Invalid characters",
			version:       "1.0.0$",
			wantErr:       true,
			errorContains: "invalid version format",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateVersion(tt.version)

			// 检查错误情况
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateVersion() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if err != nil && tt.errorContains != "" && !strings.Contains(err.Error(), tt.errorContains) {
				t.Errorf("ValidateVersion() error = %v, should contain %v", err, tt.errorContains)
			}
		})
	}
}
