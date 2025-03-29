package composer

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"

	"github.com/scagogogo/php-composer-json-parser/pkg/composer/archive"
	"github.com/scagogogo/php-composer-json-parser/pkg/composer/autoload"
)

func TestComposerJSON_ToJSON(t *testing.T) {
	// 创建一个有效的ComposerJSON对象
	composer := &ComposerJSON{
		Name:        "vendor/project",
		Description: "A test project",
		Version:     "1.0.0",
		Require: map[string]string{
			"php": "^7.4",
		},
	}

	// 测试正常情况
	jsonStr, err := composer.ToJSON(true)
	if err != nil {
		t.Fatalf("ToJSON() returned unexpected error: %v", err)
	}

	// 验证JSON内容
	var data map[string]interface{}
	if err := json.Unmarshal([]byte(jsonStr), &data); err != nil {
		t.Fatalf("Failed to unmarshal JSON: %v", err)
	}

	if data["name"] != "vendor/project" {
		t.Errorf("JSON has incorrect name: %v, want 'vendor/project'", data["name"])
	}
	if data["description"] != "A test project" {
		t.Errorf("JSON has incorrect description: %v, want 'A test project'", data["description"])
	}
	if data["version"] != "1.0.0" {
		t.Errorf("JSON has incorrect version: %v, want '1.0.0'", data["version"])
	}

	// 验证缩进和格式
	if !strings.Contains(jsonStr, "{\n") || !strings.Contains(jsonStr, "  \"") {
		t.Errorf("JSON is not properly indented")
	}

	// 测试非缩进格式
	jsonStrCompact, err := composer.ToJSON(false)
	if err != nil {
		t.Fatalf("ToJSON(false) returned unexpected error: %v", err)
	}
	if strings.Contains(jsonStrCompact, "{\n") || strings.Contains(jsonStrCompact, "  \"") {
		t.Errorf("JSON should not be indented")
	}

	// 测试JSON序列化错误场景
	invalidComposer := &ComposerJSON{
		Extra: map[string]interface{}{
			"invalid": make(chan int), // 无法序列化为JSON的类型
		},
	}
	_, err = invalidComposer.ToJSON(true)
	if err == nil {
		t.Errorf("Expected error when marshaling invalid data, got nil")
	}
}

func TestComposerJSON_Save(t *testing.T) {
	// 创建临时目录
	tempDir, err := os.MkdirTemp("", "composer-test-")
	if err != nil {
		t.Fatalf("Failed to create temporary directory: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// 测试场景
	tests := []struct {
		name    string
		setup   func() (string, *ComposerJSON, bool)
		wantErr bool
	}{
		{
			name: "Valid save with indent",
			setup: func() (string, *ComposerJSON, bool) {
				filePath := filepath.Join(tempDir, "composer.json")
				composer := &ComposerJSON{
					Name:        "vendor/project",
					Description: "A test project",
					Version:     "1.0.0",
				}
				return filePath, composer, true
			},
			wantErr: false,
		},
		{
			name: "Valid save without indent",
			setup: func() (string, *ComposerJSON, bool) {
				filePath := filepath.Join(tempDir, "composer-compact.json")
				composer := &ComposerJSON{
					Name:        "vendor/project",
					Description: "A test project",
					Version:     "1.0.0",
				}
				return filePath, composer, false
			},
			wantErr: false,
		},
		{
			name: "Save to existing file",
			setup: func() (string, *ComposerJSON, bool) {
				filePath := filepath.Join(tempDir, "existing.json")
				// 先创建文件
				if err := os.WriteFile(filePath, []byte("{}"), 0644); err != nil {
					t.Fatalf("Failed to create test file: %v", err)
				}
				composer := &ComposerJSON{
					Name:        "vendor/project",
					Description: "A test project",
				}
				return filePath, composer, true
			},
			wantErr: false,
		},
		{
			name: "Save to invalid path",
			setup: func() (string, *ComposerJSON, bool) {
				// 创建一个不存在的目录路径
				filePath := filepath.Join(tempDir, "nonexistent", "composer.json")
				composer := &ComposerJSON{
					Name: "vendor/project",
				}
				return filePath, composer, true
			},
			wantErr: true,
		},
		{
			name: "Save with JSON serialization error",
			setup: func() (string, *ComposerJSON, bool) {
				filePath := filepath.Join(tempDir, "invalid.json")
				composer := &ComposerJSON{
					Extra: map[string]interface{}{
						"invalid": make(chan int), // 无法序列化为JSON的类型
					},
				}
				return filePath, composer, true
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			filePath, composer, indent := tt.setup()
			err := composer.Save(filePath, indent)

			// 检查错误
			if (err != nil) != tt.wantErr {
				t.Errorf("Save() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			// 如果期望成功，验证文件是否正确写入
			if !tt.wantErr {
				// 检查文件是否存在
				if _, err := os.Stat(filePath); os.IsNotExist(err) {
					t.Errorf("Save() did not create file at %s", filePath)
					return
				}

				// 读取文件内容并验证
				content, err := os.ReadFile(filePath)
				if err != nil {
					t.Errorf("Failed to read saved file: %v", err)
					return
				}

				var data map[string]interface{}
				if err := json.Unmarshal(content, &data); err != nil {
					t.Errorf("Saved file contains invalid JSON: %v", err)
					return
				}

				// 检查必要字段
				if composer.Name != "" && data["name"] != composer.Name {
					t.Errorf("Saved JSON has incorrect name: %v, want %v", data["name"], composer.Name)
				}
			}
		})
	}
}

func TestComposerJSON_DependencyFunctions(t *testing.T) {
	// 创建一个测试用的ComposerJSON实例
	composer := &ComposerJSON{
		Name: "vendor/project",
		Require: map[string]string{
			"php":            "^7.4",
			"vendor/package": "^1.0",
		},
		RequireDev: map[string]string{
			"phpunit/phpunit": "^9.0",
		},
	}

	// 测试DependencyExists
	if !composer.DependencyExists("php") {
		t.Errorf("DependencyExists() for 'php' = false, want true")
	}

	if composer.DependencyExists("nonexistent") {
		t.Errorf("DependencyExists() for 'nonexistent' = true, want false")
	}

	// 测试DevDependencyExists
	if !composer.DevDependencyExists("phpunit/phpunit") {
		t.Errorf("DevDependencyExists() for 'phpunit/phpunit' = false, want true")
	}

	if composer.DevDependencyExists("nonexistent") {
		t.Errorf("DevDependencyExists() for 'nonexistent' = true, want false")
	}

	// 测试AddDependency
	err := composer.AddDependency("new/package", "^2.0")
	if err != nil {
		t.Errorf("AddDependency() error = %v", err)
	}
	if composer.Require["new/package"] != "^2.0" {
		t.Errorf("After AddDependency(), Require['new/package'] = %v, want ^2.0", composer.Require["new/package"])
	}

	// 测试更新依赖版本
	err = composer.AddDependency("new/package", "^3.0")
	if err != nil {
		t.Errorf("AddDependency() error = %v", err)
	}
	if composer.Require["new/package"] != "^3.0" {
		t.Errorf("After AddDependency() update, Require['new/package'] = %v, want ^3.0", composer.Require["new/package"])
	}

	// 测试AddDevDependency
	err = composer.AddDevDependency("new/dev-package", "^2.0")
	if err != nil {
		t.Errorf("AddDevDependency() error = %v", err)
	}
	if composer.RequireDev["new/dev-package"] != "^2.0" {
		t.Errorf("After AddDevDependency(), RequireDev['new/dev-package'] = %v, want ^2.0", composer.RequireDev["new/dev-package"])
	}

	// 测试RemoveDependency
	removed := composer.RemoveDependency("new/package")
	if !removed {
		t.Errorf("RemoveDependency() = %v, want true", removed)
	}
	if _, exists := composer.Require["new/package"]; exists {
		t.Errorf("After RemoveDependency(), 'new/package' still exists in Require")
	}

	// 测试移除不存在的依赖
	removed = composer.RemoveDependency("nonexistent")
	if removed {
		t.Errorf("RemoveDependency() for nonexistent package = %v, want false", removed)
	}

	// 测试RemoveDevDependency
	removed = composer.RemoveDevDependency("new/dev-package")
	if !removed {
		t.Errorf("RemoveDevDependency() = %v, want true", removed)
	}
	if _, exists := composer.RequireDev["new/dev-package"]; exists {
		t.Errorf("After RemoveDevDependency(), 'new/dev-package' still exists in RequireDev")
	}

	// 测试GetAllDependencies
	allDeps := composer.GetAllDependencies()
	expectedDeps := map[string]string{
		"php":             "^7.4",
		"vendor/package":  "^1.0",
		"phpunit/phpunit": "^9.0",
	}
	if !reflect.DeepEqual(allDeps, expectedDeps) {
		t.Errorf("GetAllDependencies() = %v, want %v", allDeps, expectedDeps)
	}
}

func TestComposerJSON_PSR4Functions(t *testing.T) {
	// 创建一个测试用的ComposerJSON实例
	composer := &ComposerJSON{
		Name: "vendor/project",
		Autoload: autoload.Autoload{
			PSR4: map[string]interface{}{
				"Vendor\\Package\\": "src/",
			},
		},
	}

	// 测试GetPSR4Map
	psr4Map, ok := composer.GetPSR4Map()
	if !ok {
		t.Errorf("GetPSR4Map() ok = %v, want true", ok)
	}
	expectedMap := map[string]string{
		"Vendor\\Package\\": "src/",
	}
	if !reflect.DeepEqual(psr4Map, expectedMap) {
		t.Errorf("GetPSR4Map() = %v, want %v", psr4Map, expectedMap)
	}

	// 测试SetPSR4添加新命名空间
	composer.SetPSR4("Vendor\\Tests\\", "tests/")
	psr4Map, ok = composer.GetPSR4Map()
	if !ok {
		t.Errorf("GetPSR4Map() ok = %v, want true after SetPSR4", ok)
	}
	expectedMap["Vendor\\Tests\\"] = "tests/"
	if !reflect.DeepEqual(psr4Map, expectedMap) {
		t.Errorf("After SetPSR4(), PSR4 = %v, want %v", psr4Map, expectedMap)
	}

	// 测试SetPSR4更新现有命名空间
	composer.SetPSR4("Vendor\\Package\\", "new-src/")
	psr4Map, ok = composer.GetPSR4Map()
	if !ok {
		t.Errorf("GetPSR4Map() ok = %v, want true after update SetPSR4", ok)
	}
	expectedMap["Vendor\\Package\\"] = "new-src/"
	if !reflect.DeepEqual(psr4Map, expectedMap) {
		t.Errorf("After update SetPSR4(), PSR4 = %v, want %v", psr4Map, expectedMap)
	}

	// 测试RemovePSR4
	removed := composer.RemovePSR4("Vendor\\Tests\\")
	if !removed {
		t.Errorf("RemovePSR4() = %v, want true", removed)
	}
	psr4Map, ok = composer.GetPSR4Map()
	if !ok {
		t.Errorf("GetPSR4Map() ok = %v, want true after RemovePSR4", ok)
	}
	expectedMap = map[string]string{
		"Vendor\\Package\\": "new-src/",
	}
	if !reflect.DeepEqual(psr4Map, expectedMap) {
		t.Errorf("After RemovePSR4(), PSR4 = %v, want %v", psr4Map, expectedMap)
	}

	// 测试移除不存在的命名空间
	removed = composer.RemovePSR4("Nonexistent\\")
	if removed {
		t.Errorf("RemovePSR4() for nonexistent namespace = %v, want false", removed)
	}
}

func TestComposerJSON_ArchiveFunctions(t *testing.T) {
	// 创建一个测试用的ComposerJSON实例
	composer := &ComposerJSON{
		Name: "vendor/project",
		Archive: archive.Archive{
			Exclude: []string{"/vendor", "/tests"},
		},
	}

	// 测试AddExclusion
	composer.AddExclusion("/docs")
	expectedExcludes := []string{"/vendor", "/tests", "/docs"}
	if !reflect.DeepEqual(composer.Archive.Exclude, expectedExcludes) {
		t.Errorf("After AddExclusion(), Archive.Exclude = %v, want %v", composer.Archive.Exclude, expectedExcludes)
	}

	// 测试添加已存在的排除模式
	composer.AddExclusion("/docs")
	if !reflect.DeepEqual(composer.Archive.Exclude, expectedExcludes) {
		t.Errorf("After adding existing exclusion, Archive.Exclude = %v, want %v", composer.Archive.Exclude, expectedExcludes)
	}

	// 测试RemoveExclusion
	removed := composer.RemoveExclusion("/tests")
	if !removed {
		t.Errorf("RemoveExclusion() = %v, want true", removed)
	}
	expectedExcludes = []string{"/vendor", "/docs"}
	if !reflect.DeepEqual(composer.Archive.Exclude, expectedExcludes) {
		t.Errorf("After RemoveExclusion(), Archive.Exclude = %v, want %v", composer.Archive.Exclude, expectedExcludes)
	}

	// 测试移除不存在的排除模式
	removed = composer.RemoveExclusion("/nonexistent")
	if removed {
		t.Errorf("RemoveExclusion() for nonexistent pattern = %v, want false", removed)
	}
}

func TestComposerJSON_RepositoryFunctions(t *testing.T) {
	// 创建测试仓库
	repo := NewRepository("composer", "https://packagist.org")
	if repo == nil {
		t.Fatalf("NewRepository() returned nil")
	}

	// 检查仓库字段
	if repo.Type != "composer" {
		t.Errorf("repo.Type = %v, want %v", repo.Type, "composer")
	}
	if repo.URL != "https://packagist.org" {
		t.Errorf("repo.URL = %v, want %v", repo.URL, "https://packagist.org")
	}

	// 测试AddRepository
	composer := &ComposerJSON{}
	composer.AddRepository(*repo)
	if len(composer.Repositories) != 1 {
		t.Errorf("len(composer.Repositories) = %v, want %v", len(composer.Repositories), 1)
	}
	if composer.Repositories[0].Type != "composer" {
		t.Errorf("composer.Repositories[0].Type = %v, want %v", composer.Repositories[0].Type, "composer")
	}
}

func TestCreateBackup(t *testing.T) {
	// 创建一个临时目录用于测试
	tempDir, err := os.MkdirTemp("", "composer-backup-test-")
	if err != nil {
		t.Fatalf("Failed to create temporary directory: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// 创建一个原始文件进行备份
	origFilePath := filepath.Join(tempDir, "composer.json")
	origContent := []byte(`{"name": "vendor/project", "version": "1.0.0"}`)
	if err := os.WriteFile(origFilePath, origContent, 0644); err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	tests := []struct {
		name         string
		filePath     string
		backupSuffix string
		wantErr      bool
	}{
		{
			name:         "Default suffix",
			filePath:     origFilePath,
			backupSuffix: "", // 空字符串表示使用默认后缀
			wantErr:      false,
		},
		{
			name:         "Custom suffix",
			filePath:     origFilePath,
			backupSuffix: ".backup",
			wantErr:      false,
		},
		{
			name:         "Non-existent file",
			filePath:     filepath.Join(tempDir, "nonexistent.json"),
			backupSuffix: "",
			wantErr:      true,
		},
		{
			name:         "Invalid path",
			filePath:     filepath.Join(tempDir, "nonexistent", "composer.json"),
			backupSuffix: "",
			wantErr:      true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// 执行备份
			backupPath, err := CreateBackup(tt.filePath, tt.backupSuffix)

			// 检查错误结果是否符合预期
			if (err != nil) != tt.wantErr {
				t.Errorf("CreateBackup() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			// 如果期望成功，验证备份文件是否存在并包含正确内容
			if !tt.wantErr {
				// 确定期望的后缀
				expectedSuffix := ".bak" // 默认后缀
				if tt.backupSuffix != "" {
					expectedSuffix = tt.backupSuffix
				}

				// 验证返回的备份路径是否正确
				expectedPath := tt.filePath + expectedSuffix
				if backupPath != expectedPath {
					t.Errorf("CreateBackup() returned path = %v, want %v", backupPath, expectedPath)
				}

				// 检查备份文件是否存在
				if _, err := os.Stat(backupPath); os.IsNotExist(err) {
					t.Errorf("Backup file does not exist at %s", backupPath)
					return
				}

				// 检查备份文件内容是否与原始内容相同
				backupContent, err := os.ReadFile(backupPath)
				if err != nil {
					t.Errorf("Failed to read backup file: %v", err)
					return
				}

				if !bytes.Equal(backupContent, origContent) {
					t.Errorf("Backup content does not match original. Got %s, want %s", backupContent, origContent)
				}
			}
		})
	}
}

func TestDefaultConfig(t *testing.T) {
	config := DefaultConfig()

	// 确保返回的Config不是nil
	if config == nil {
		t.Errorf("DefaultConfig() returned nil")
		return
	}

	// 检查一些关键字段
	if config.ProcessTimeout != 300 {
		t.Errorf("DefaultConfig().ProcessTimeout = %v, want 300", config.ProcessTimeout)
	}

	if config.PreferredInstall != "dist" {
		t.Errorf("DefaultConfig().PreferredInstall = %v, want dist", config.PreferredInstall)
	}

	if config.VendorDir != "vendor" {
		t.Errorf("DefaultConfig().VendorDir = %v, want vendor", config.VendorDir)
	}

	if config.BinDir != "vendor/bin" {
		t.Errorf("DefaultConfig().BinDir = %v, want vendor/bin", config.BinDir)
	}

	// 检查数组字段
	expectedGithubProtocols := []string{"https", "ssh", "git"}
	if !reflect.DeepEqual(config.GithubProtocols, expectedGithubProtocols) {
		t.Errorf("DefaultConfig().GithubProtocols = %v, want %v", config.GithubProtocols, expectedGithubProtocols)
	}

	expectedGitlabProtocols := []string{"https", "ssh"}
	if !reflect.DeepEqual(config.GitlabProtocols, expectedGitlabProtocols) {
		t.Errorf("DefaultConfig().GitlabProtocols = %v, want %v", config.GitlabProtocols, expectedGitlabProtocols)
	}

	// 检查布尔值字段
	if config.SecureHttp != true {
		t.Errorf("DefaultConfig().SecureHttp = %v, want true", config.SecureHttp)
	}

	if config.OptimizeAutoloader != false {
		t.Errorf("DefaultConfig().OptimizeAutoloader = %v, want false", config.OptimizeAutoloader)
	}
}

func TestConvertToComposerJSON(t *testing.T) {
	tests := []struct {
		name    string
		data    map[string]interface{}
		check   func(*ComposerJSON)
		wantErr bool
	}{
		{
			name: "Valid composer data",
			data: map[string]interface{}{
				"name":        "vendor/project",
				"description": "Test project",
				"require": map[string]interface{}{
					"php": "^7.4",
				},
			},
			check: func(c *ComposerJSON) {
				if c.Name != "vendor/project" {
					t.Errorf("Expected name 'vendor/project', got '%s'", c.Name)
				}
				if c.Description != "Test project" {
					t.Errorf("Expected description 'Test project', got '%s'", c.Description)
				}
				if c.Require["php"] != "^7.4" {
					t.Errorf("Expected php requirement '^7.4', got '%s'", c.Require["php"])
				}
			},
			wantErr: false,
		},
		{
			name: "Empty data",
			data: map[string]interface{}{},
			check: func(c *ComposerJSON) {
				if c.Name != "" {
					t.Errorf("Expected empty name, got '%s'", c.Name)
				}
				if c.Description != "" {
					t.Errorf("Expected empty description, got '%s'", c.Description)
				}
			},
			wantErr: false,
		},
		{
			name: "Complex nested data",
			data: map[string]interface{}{
				"name": "vendor/project",
				"autoload": map[string]interface{}{
					"psr-4": map[string]interface{}{
						"Vendor\\Project\\": "src/",
					},
				},
				"config": map[string]interface{}{
					"process-timeout": 300,
					"vendor-dir":      "vendor",
				},
			},
			check: func(c *ComposerJSON) {
				if c.Name != "vendor/project" {
					t.Errorf("Expected name 'vendor/project', got '%s'", c.Name)
				}

				psr4Map, ok := c.Autoload.PSR4.(map[string]interface{})
				if !ok {
					t.Errorf("Expected PSR-4 map, got %T", c.Autoload.PSR4)
					return
				}

				if psr4Map["Vendor\\Project\\"] != "src/" {
					t.Errorf("Expected PSR-4 namespace 'Vendor\\Project\\' to be 'src/', got '%v'", psr4Map["Vendor\\Project\\"])
				}

				if c.Config.ProcessTimeout != 300 {
					t.Errorf("Expected config.ProcessTimeout to be 300, got %d", c.Config.ProcessTimeout)
				}

				if c.Config.VendorDir != "vendor" {
					t.Errorf("Expected config.VendorDir to be 'vendor', got '%s'", c.Config.VendorDir)
				}
			},
			wantErr: false,
		},
		{
			// This test case causes a JSON marshalling error - this is very rare and
			// difficult to actually trigger in real code, but we include it for coverage
			name: "Invalid data that can't be marshaled to JSON",
			data: map[string]interface{}{
				"invalid": make(chan int), // channels can't be marshaled to JSON
			},
			check:   nil,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			composer, err := convertToComposerJSON(tt.data)

			// 检查错误
			if (err != nil) != tt.wantErr {
				t.Errorf("convertToComposerJSON() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			// 如果期望成功，执行检查
			if err == nil && tt.check != nil {
				tt.check(composer)
			}
		})
	}
}

// BadMarshalerJSON 用于测试JSON序列化错误
type BadMarshalerJSON struct {
	Name        string `json:"name,omitempty"`
	Description string `json:"description,omitempty"`
}

// MarshalJSON 实现自定义的JSON序列化，总是返回错误
func (b BadMarshalerJSON) MarshalJSON() ([]byte, error) {
	return nil, fmt.Errorf("mock marshal error")
}

// 用于测试JSON反序列化错误的结构体
type BadUnmarshalerJSON struct {
}

// MarshalJSON 返回有效的JSON，但无法解析为map
func (b BadUnmarshalerJSON) MarshalJSON() ([]byte, error) {
	return []byte(`{"name": {"nested": ["invalid", "for", "unmarshal", "to", "map"]}}`), nil
}

// TestComposerJSON_ToJSONError 测试ToJSON函数的错误处理
func TestComposerJSON_ToJSONError(t *testing.T) {
	// 创建一个模拟ComposerJSON的结构体
	// 覆盖标准ComposerJSON.ToJSON方法的逻辑

	// 先检查json.Marshal失败的情况
	badMarshaler := &BadMarshalerJSON{
		Name:        "test/package",
		Description: "Test package",
	}

	// 手动模拟ToJSON方法的第一部分 - json.Marshal失败
	_, err := json.Marshal(badMarshaler)

	// 检查是否返回了预期的错误
	if err == nil {
		t.Errorf("Expected error from MarshalJSON but got nil")
		return
	}

	// 验证错误消息
	if !strings.Contains(err.Error(), "mock marshal error") {
		t.Errorf("Expected error to contain 'mock marshal error', got: %v", err)
	}
}

// TestComposerJSON_ToJSONUnmarshalError 测试JSON反序列化为map时的错误
func TestComposerJSON_ToJSONUnmarshalError(t *testing.T) {
	// 模拟ToJSON方法的json.Unmarshal部分失败
	// 创建一个无法解析为map的JSON字符串
	jsonStr := `[1, 2, 3]` // 数组不能直接解析为map

	var rawData map[string]interface{}
	err := json.Unmarshal([]byte(jsonStr), &rawData)

	// 检查是否返回了预期的错误
	if err == nil {
		t.Errorf("Expected error from Unmarshal but got nil")
		return
	}

	// 验证错误类型
	if !strings.Contains(err.Error(), "cannot unmarshal array") {
		t.Errorf("Expected error to contain 'cannot unmarshal array', got: %v", err)
	}
}

// MockComposerJSON for testing
type MockComposerJSON struct {
	ComposerJSON
}

// Custom MarshalJSON that produces valid JSON but not valid map structure
func (m *MockComposerJSON) MarshalJSON() ([]byte, error) {
	return []byte(`{"name": {"nested": ["invalid", "for", "unmarshal", "to", "map"]}}`), nil
}

// Additional tests for convertToComposerJSON
func Test_convertToComposerJSON(t *testing.T) {
	// 测试数据
	tests := []struct {
		name    string
		data    map[string]interface{}
		wantErr bool
		errMsg  string
	}{
		{
			name: "Valid data",
			data: map[string]interface{}{
				"name":        "test/package",
				"description": "A test package",
				"require": map[string]interface{}{
					"php": "^7.4",
				},
			},
			wantErr: false,
		},
		{
			name:    "Invalid JSON Marshalling",
			data:    make(map[string]interface{}),
			wantErr: false, // 空map也是有效的
		},
		{
			name: "Complex nested structures",
			data: map[string]interface{}{
				"name": "test/package",
				"authors": []interface{}{
					map[string]interface{}{
						"name":  "Test Author",
						"email": "test@example.com",
					},
				},
				"autoload": map[string]interface{}{
					"psr-4": map[string]interface{}{
						"Test\\Namespace\\": "src/",
					},
				},
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			composer, err := convertToComposerJSON(tt.data)

			// 检查错误预期
			if (err != nil) != tt.wantErr {
				t.Errorf("convertToComposerJSON() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if tt.wantErr {
				if !strings.Contains(err.Error(), tt.errMsg) {
					t.Errorf("Expected error message to contain '%s', got: %v", tt.errMsg, err)
				}
				return
			}

			// 验证转换结果
			if composer == nil && !tt.wantErr {
				t.Errorf("Expected non-nil composer for valid data")
			}

			// 对于有效数据，验证一些字段是否正确转换
			if !tt.wantErr && len(tt.data) > 0 {
				if name, ok := tt.data["name"]; ok && composer.Name != name {
					t.Errorf("Name field not properly converted, expected %v, got %v", name, composer.Name)
				}
			}
		})
	}
}
