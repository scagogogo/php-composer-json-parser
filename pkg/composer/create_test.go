package composer

import (
	"reflect"
	"testing"
)

func TestCreateNew(t *testing.T) {
	tests := []struct {
		name          string
		packageName   string
		description   string
		wantName      string
		wantDesc      string
		wantType      string
		wantRequire   map[string]string
		wantNamespace string
		wantErr       bool
	}{
		{
			name:          "Create with valid name and description",
			packageName:   "vendor/project",
			description:   "A test project",
			wantName:      "vendor/project",
			wantDesc:      "A test project",
			wantType:      "library",
			wantRequire:   map[string]string{},
			wantNamespace: "Vendor\\Project\\",
			wantErr:       false,
		},
		{
			name:          "Create with empty name and description",
			packageName:   "",
			description:   "",
			wantName:      "",
			wantDesc:      "",
			wantType:      "library",
			wantRequire:   map[string]string{},
			wantNamespace: "",
			wantErr:       false,
		},
		{
			name:          "Create with invalid name",
			packageName:   "invalid-name",
			description:   "A test project",
			wantName:      "",
			wantDesc:      "",
			wantType:      "",
			wantRequire:   nil,
			wantNamespace: "",
			wantErr:       true,
		},
		{
			name:          "Create with short description",
			packageName:   "vendor/project",
			description:   "Short",
			wantName:      "",
			wantDesc:      "",
			wantType:      "",
			wantRequire:   nil,
			wantNamespace: "",
			wantErr:       true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			composer, err := CreateNew(tt.packageName, tt.description)

			// 检查错误
			if (err != nil) != tt.wantErr {
				t.Errorf("CreateNew() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			// 如果期望错误，则不检查数据
			if tt.wantErr {
				return
			}

			// 检查基本字段
			if composer.Name != tt.wantName {
				t.Errorf("composer.Name = %v, want %v", composer.Name, tt.wantName)
			}

			if composer.Description != tt.wantDesc {
				t.Errorf("composer.Description = %v, want %v", composer.Description, tt.wantDesc)
			}

			if composer.Type != tt.wantType {
				t.Errorf("composer.Type = %v, want %v", composer.Type, tt.wantType)
			}

			// 检查require映射是否存在且为空
			if !reflect.DeepEqual(composer.Require, tt.wantRequire) {
				t.Errorf("composer.Require = %v, want %v", composer.Require, tt.wantRequire)
			}

			if tt.packageName != "" {
				// 检查PSR-4命名空间
				psr4Map, ok := composer.GetPSR4Map()
				if !ok {
					t.Errorf("GetPSR4Map() ok = false, want true")
					return
				}

				// 检查命名空间映射
				if psr4Map[tt.wantNamespace] != "src/" {
					t.Errorf("PSR4 namespace mapping: got %v, want %v -> src/", psr4Map, tt.wantNamespace)
				}
			}
		})
	}
}

func TestCreateProject(t *testing.T) {
	tests := []struct {
		name           string
		packageName    string
		description    string
		phpVersion     string
		wantName       string
		wantDesc       string
		wantType       string
		wantPhpVer     string
		wantPHPUnitVer string
		wantErr        bool
	}{
		{
			name:           "Create project with valid parameters",
			packageName:    "vendor/project",
			description:    "A test project",
			phpVersion:     "^8.0",
			wantName:       "vendor/project",
			wantDesc:       "A test project",
			wantType:       "project",
			wantPhpVer:     "^8.0",
			wantPHPUnitVer: "^9.0",
			wantErr:        false,
		},
		{
			name:           "Create project with empty PHP version",
			packageName:    "vendor/project",
			description:    "A test project",
			phpVersion:     "",
			wantName:       "vendor/project",
			wantDesc:       "A test project",
			wantType:       "project",
			wantPhpVer:     "^7.4",
			wantPHPUnitVer: "^9.0",
			wantErr:        false,
		},
		{
			name:           "Create project with invalid package name",
			packageName:    "invalid-name",
			description:    "A test project",
			phpVersion:     "^8.0",
			wantName:       "",
			wantDesc:       "",
			wantType:       "",
			wantPhpVer:     "",
			wantPHPUnitVer: "",
			wantErr:        true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			composer, err := CreateProject(tt.packageName, tt.description, tt.phpVersion)

			// 检查错误
			if (err != nil) != tt.wantErr {
				t.Errorf("CreateProject() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			// 如果期望错误，则不检查数据
			if tt.wantErr {
				return
			}

			// 检查基本字段
			if composer.Name != tt.wantName {
				t.Errorf("composer.Name = %v, want %v", composer.Name, tt.wantName)
			}

			if composer.Description != tt.wantDesc {
				t.Errorf("composer.Description = %v, want %v", composer.Description, tt.wantDesc)
			}

			if composer.Type != tt.wantType {
				t.Errorf("composer.Type = %v, want %v", composer.Type, tt.wantType)
			}

			// 检查PHP版本
			if composer.Require["php"] != tt.wantPhpVer {
				t.Errorf("composer.Require[\"php\"] = %v, want %v", composer.Require["php"], tt.wantPhpVer)
			}

			// 检查PHPUnit开发依赖
			if composer.RequireDev["phpunit/phpunit"] != tt.wantPHPUnitVer {
				t.Errorf("composer.RequireDev[\"phpunit/phpunit\"] = %v, want %v", composer.RequireDev["phpunit/phpunit"], tt.wantPHPUnitVer)
			}
		})
	}
}

func TestCreateLibrary(t *testing.T) {
	tests := []struct {
		name        string
		packageName string
		description string
		phpVersion  string
		wantName    string
		wantDesc    string
		wantType    string
		wantPhpVer  string
		wantErr     bool
	}{
		{
			name:        "Create library with valid parameters",
			packageName: "vendor/library",
			description: "A test library",
			phpVersion:  "^8.1",
			wantName:    "vendor/library",
			wantDesc:    "A test library",
			wantType:    "library",
			wantPhpVer:  "^8.1",
			wantErr:     false,
		},
		{
			name:        "Create library with empty PHP version",
			packageName: "vendor/library",
			description: "A test library",
			phpVersion:  "",
			wantName:    "vendor/library",
			wantDesc:    "A test library",
			wantType:    "library",
			wantPhpVer:  "^7.4",
			wantErr:     false,
		},
		{
			name:        "Create library with invalid name",
			packageName: "invalid-name",
			description: "A test library",
			phpVersion:  "^8.1",
			wantName:    "",
			wantDesc:    "",
			wantType:    "",
			wantPhpVer:  "",
			wantErr:     true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			composer, err := CreateLibrary(tt.packageName, tt.description, tt.phpVersion)

			// 检查错误
			if (err != nil) != tt.wantErr {
				t.Errorf("CreateLibrary() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			// 如果期望错误，则不检查数据
			if tt.wantErr {
				return
			}

			// 检查基本字段
			if composer.Name != tt.wantName {
				t.Errorf("composer.Name = %v, want %v", composer.Name, tt.wantName)
			}

			if composer.Description != tt.wantDesc {
				t.Errorf("composer.Description = %v, want %v", composer.Description, tt.wantDesc)
			}

			if composer.Type != tt.wantType {
				t.Errorf("composer.Type = %v, want %v", composer.Type, tt.wantType)
			}

			// 检查PHP版本
			if composer.Require["php"] != tt.wantPhpVer {
				t.Errorf("composer.Require[\"php\"] = %v, want %v", composer.Require["php"], tt.wantPhpVer)
			}

			// 确保RequireDev为空或空Map
			if len(composer.RequireDev) > 0 {
				t.Errorf("composer.RequireDev should be empty, got %v", composer.RequireDev)
			}
		})
	}
}

func TestHelperFunctions(t *testing.T) {
	// 测试toNamespace函数
	tests := []struct {
		name       string
		vendor     string
		project    string
		wantResult string
	}{
		{
			name:       "Basic conversion",
			vendor:     "vendor",
			project:    "project",
			wantResult: "Vendor\\Project",
		},
		{
			name:       "Already capitalized",
			vendor:     "Vendor",
			project:    "Project",
			wantResult: "Vendor\\Project",
		},
		{
			name:       "Empty strings",
			vendor:     "",
			project:    "",
			wantResult: "\\",
		},
		{
			name:       "Mixed case",
			vendor:     "myVendor",
			project:    "myProject",
			wantResult: "MyVendor\\MyProject",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := toNamespace(tt.vendor, tt.project)
			if got != tt.wantResult {
				t.Errorf("toNamespace() = %v, want %v", got, tt.wantResult)
			}
		})
	}

	// 测试ucfirst函数
	ucfirstTests := []struct {
		name       string
		input      string
		wantResult string
	}{
		{
			name:       "Lowercase first letter",
			input:      "test",
			wantResult: "Test",
		},
		{
			name:       "Already capitalized",
			input:      "Test",
			wantResult: "Test",
		},
		{
			name:       "Empty string",
			input:      "",
			wantResult: "",
		},
		{
			name:       "Single character",
			input:      "a",
			wantResult: "A",
		},
		{
			name:       "Non-letter first character",
			input:      "123test",
			wantResult: "123test",
		},
	}

	for _, tt := range ucfirstTests {
		t.Run(tt.name, func(t *testing.T) {
			got := ucfirst(tt.input)
			if got != tt.wantResult {
				t.Errorf("ucfirst() = %v, want %v", got, tt.wantResult)
			}
		})
	}
}

func TestValidationFunctions(t *testing.T) {
	// 测试ValidateComposerJSON函数
	validateTests := []struct {
		name        string
		packageName string
		description string
		stability   string
		wantErr     bool
	}{
		{
			name:        "Valid data",
			packageName: "vendor/project",
			description: "This is a sufficient description",
			stability:   "stable",
			wantErr:     false,
		},
		{
			name:        "Invalid package name",
			packageName: "invalidname",
			description: "This is a sufficient description",
			stability:   "stable",
			wantErr:     true,
		},
		{
			name:        "Short description",
			packageName: "vendor/project",
			description: "Short",
			stability:   "stable",
			wantErr:     true,
		},
		{
			name:        "Invalid stability",
			packageName: "vendor/project",
			description: "This is a sufficient description",
			stability:   "invalid",
			wantErr:     true,
		},
		{
			name:        "Empty values",
			packageName: "",
			description: "",
			stability:   "",
			wantErr:     false,
		},
	}

	for _, tt := range validateTests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateComposerJSON(tt.packageName, tt.description, tt.stability)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateComposerJSON() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}

	// 测试ValidateVersion函数
	versionTests := []struct {
		name    string
		version string
		wantErr bool
	}{
		{
			name:    "Valid version",
			version: "^7.4",
			wantErr: false,
		},
		{
			name:    "Another valid version",
			version: "~8.0.0",
			wantErr: false,
		},
		{
			name:    "Empty version",
			version: "",
			wantErr: false,
		},
		{
			name:    "Invalid version",
			version: "invalid$",
			wantErr: true,
		},
	}

	for _, tt := range versionTests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateVersion(tt.version)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateVersion() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
