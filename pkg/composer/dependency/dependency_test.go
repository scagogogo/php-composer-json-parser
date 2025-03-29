package dependency

import (
	"strings"
	"testing"
)

func TestGetPackageNameParts(t *testing.T) {
	tests := []struct {
		name          string
		packageName   string
		wantVendor    string
		wantProject   string
		wantErr       bool
		errorContains string
	}{
		{
			name:        "Valid package name",
			packageName: "vendor/project",
			wantVendor:  "vendor",
			wantProject: "project",
			wantErr:     false,
		},
		{
			name:          "Empty package name",
			packageName:   "",
			wantVendor:    "",
			wantProject:   "",
			wantErr:       true,
			errorContains: "invalid package name format",
		},
		{
			name:          "No slash",
			packageName:   "invalidname",
			wantVendor:    "",
			wantProject:   "",
			wantErr:       true,
			errorContains: "invalid package name format",
		},
		{
			name:          "Too many slashes",
			packageName:   "vendor/sub/project",
			wantVendor:    "",
			wantProject:   "",
			wantErr:       true,
			errorContains: "invalid package name format",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			vendor, project, err := GetPackageNameParts(tt.packageName)

			// 检查错误情况
			if (err != nil) != tt.wantErr {
				t.Errorf("GetPackageNameParts() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if err != nil && tt.errorContains != "" && !contains(err.Error(), tt.errorContains) {
				t.Errorf("GetPackageNameParts() error = %v, should contain %v", err, tt.errorContains)
				return
			}

			// 检查返回值
			if vendor != tt.wantVendor {
				t.Errorf("GetPackageNameParts() vendor = %v, want %v", vendor, tt.wantVendor)
			}

			if project != tt.wantProject {
				t.Errorf("GetPackageNameParts() project = %v, want %v", project, tt.wantProject)
			}
		})
	}
}

func TestValidatePackageName(t *testing.T) {
	tests := []struct {
		name          string
		packageName   string
		wantErr       bool
		errorContains string
	}{
		{
			name:        "Valid package name",
			packageName: "vendor/project",
			wantErr:     false,
		},
		{
			name:          "Empty package name",
			packageName:   "",
			wantErr:       true,
			errorContains: "cannot be empty",
		},
		{
			name:          "No slash",
			packageName:   "invalidname",
			wantErr:       true,
			errorContains: "format 'vendor/project'",
		},
		{
			name:          "Invalid vendor name with uppercase",
			packageName:   "Vendor/project",
			wantErr:       true,
			errorContains: "invalid vendor name",
		},
		{
			name:          "Invalid project name with uppercase",
			packageName:   "vendor/Project",
			wantErr:       true,
			errorContains: "invalid project name",
		},
		{
			name:          "Invalid characters in vendor",
			packageName:   "vendor$/project",
			wantErr:       true,
			errorContains: "invalid vendor name",
		},
		{
			name:          "Invalid characters in project",
			packageName:   "vendor/project$",
			wantErr:       true,
			errorContains: "invalid project name",
		},
		{
			name:        "Valid with hyphen",
			packageName: "my-vendor/my-project",
			wantErr:     false,
		},
		{
			name:        "Valid with underscore",
			packageName: "my_vendor/my_project",
			wantErr:     false,
		},
		{
			name:        "Valid with dot",
			packageName: "my.vendor/my.project",
			wantErr:     false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidatePackageName(tt.packageName)

			// 检查错误情况
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidatePackageName() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if err != nil && tt.errorContains != "" && !contains(err.Error(), tt.errorContains) {
				t.Errorf("ValidatePackageName() error = %v, should contain %v", err, tt.errorContains)
			}
		})
	}
}

func TestDependencyExists(t *testing.T) {
	// 准备测试数据
	require := map[string]string{
		"vendor/package1": "^1.0",
		"vendor/package2": "^2.0",
	}

	tests := []struct {
		name        string
		require     map[string]string
		packageName string
		want        bool
	}{
		{
			name:        "Existing package",
			require:     require,
			packageName: "vendor/package1",
			want:        true,
		},
		{
			name:        "Non-existing package",
			require:     require,
			packageName: "vendor/nonexistent",
			want:        false,
		},
		{
			name:        "Nil require map",
			require:     nil,
			packageName: "vendor/package1",
			want:        false,
		},
		{
			name:        "Empty require map",
			require:     map[string]string{},
			packageName: "vendor/package1",
			want:        false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := DependencyExists(tt.require, tt.packageName); got != tt.want {
				t.Errorf("DependencyExists() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestAddDependency(t *testing.T) {
	tests := []struct {
		name        string
		require     map[string]string
		packageName string
		version     string
		wantErr     bool
		wantRequire map[string]string
	}{
		{
			name:        "Add new dependency",
			require:     map[string]string{"existing/pkg": "^1.0"},
			packageName: "vendor/package",
			version:     "^2.0",
			wantErr:     false,
			wantRequire: map[string]string{
				"existing/pkg":   "^1.0",
				"vendor/package": "^2.0",
			},
		},
		{
			name:        "Update existing dependency",
			require:     map[string]string{"vendor/package": "^1.0"},
			packageName: "vendor/package",
			version:     "^2.0",
			wantErr:     false,
			wantRequire: map[string]string{"vendor/package": "^2.0"},
		},
		{
			name:        "Invalid package name",
			require:     map[string]string{},
			packageName: "invalid-name", // 没有slash
			version:     "^1.0",
			wantErr:     true,
			wantRequire: map[string]string{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// 复制原始require，确保不影响其他测试
			require := make(map[string]string)
			for k, v := range tt.require {
				require[k] = v
			}

			err := AddDependency(require, tt.packageName, tt.version)

			// 检查错误情况
			if (err != nil) != tt.wantErr {
				t.Errorf("AddDependency() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			// 如果没有错误，检查是否添加成功
			if !tt.wantErr {
				if len(require) != len(tt.wantRequire) {
					t.Errorf("AddDependency() require length = %v, want %v", len(require), len(tt.wantRequire))
				}

				for pkg, version := range tt.wantRequire {
					if require[pkg] != version {
						t.Errorf("AddDependency() require[%s] = %v, want %v", pkg, require[pkg], version)
					}
				}
			}
		})
	}
}

func TestRemoveDependency(t *testing.T) {
	tests := []struct {
		name        string
		require     map[string]string
		packageName string
		want        bool
		wantRequire map[string]string
	}{
		{
			name:        "Remove existing dependency",
			require:     map[string]string{"vendor/package": "^1.0", "other/pkg": "^2.0"},
			packageName: "vendor/package",
			want:        true,
			wantRequire: map[string]string{"other/pkg": "^2.0"},
		},
		{
			name:        "Remove non-existing dependency",
			require:     map[string]string{"other/pkg": "^2.0"},
			packageName: "vendor/package",
			want:        false,
			wantRequire: map[string]string{"other/pkg": "^2.0"},
		},
		{
			name:        "Nil require map",
			require:     nil,
			packageName: "vendor/package",
			want:        false,
			wantRequire: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// 复制原始require，确保不影响其他测试
			var require map[string]string
			if tt.require != nil {
				require = make(map[string]string)
				for k, v := range tt.require {
					require[k] = v
				}
			}

			if got := RemoveDependency(require, tt.packageName); got != tt.want {
				t.Errorf("RemoveDependency() = %v, want %v", got, tt.want)
			}

			// 检查移除后的结果
			if require != nil {
				if len(require) != len(tt.wantRequire) {
					t.Errorf("RemoveDependency() require length = %v, want %v", len(require), len(tt.wantRequire))
				}

				for pkg, version := range tt.wantRequire {
					if require[pkg] != version {
						t.Errorf("RemoveDependency() require[%s] = %v, want %v", pkg, require[pkg], version)
					}
				}
			}
		})
	}
}

func TestMergeDependencies(t *testing.T) {
	tests := []struct {
		name       string
		require    map[string]string
		requireDev map[string]string
		want       map[string]string
	}{
		{
			name:       "Both maps have values",
			require:    map[string]string{"pkg1": "^1.0", "pkg2": "^2.0"},
			requireDev: map[string]string{"pkg3": "^3.0", "pkg4": "^4.0"},
			want:       map[string]string{"pkg1": "^1.0", "pkg2": "^2.0", "pkg3": "^3.0", "pkg4": "^4.0"},
		},
		{
			name:       "Require is empty",
			require:    map[string]string{},
			requireDev: map[string]string{"pkg3": "^3.0", "pkg4": "^4.0"},
			want:       map[string]string{"pkg3": "^3.0", "pkg4": "^4.0"},
		},
		{
			name:       "RequireDev is empty",
			require:    map[string]string{"pkg1": "^1.0", "pkg2": "^2.0"},
			requireDev: map[string]string{},
			want:       map[string]string{"pkg1": "^1.0", "pkg2": "^2.0"},
		},
		{
			name:       "Both maps are empty",
			require:    map[string]string{},
			requireDev: map[string]string{},
			want:       map[string]string{},
		},
		{
			name:       "Overlapping packages",
			require:    map[string]string{"pkg1": "^1.0", "common": "^1.0"},
			requireDev: map[string]string{"pkg2": "^2.0", "common": "^2.0"},
			want:       map[string]string{"pkg1": "^1.0", "pkg2": "^2.0", "common": "^2.0"},
		},
		{
			name:       "Nil maps",
			require:    nil,
			requireDev: nil,
			want:       map[string]string{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := MergeDependencies(tt.require, tt.requireDev)

			if len(got) != len(tt.want) {
				t.Errorf("MergeDependencies() length = %v, want %v", len(got), len(tt.want))
			}

			for pkg, version := range tt.want {
				if got[pkg] != version {
					t.Errorf("MergeDependencies() got[%s] = %v, want %v", pkg, got[pkg], version)
				}
			}
		})
	}
}

// 辅助函数，检查字符串是否包含指定子串
func contains(s, substring string) bool {
	return strings.Contains(s, substring)
}
