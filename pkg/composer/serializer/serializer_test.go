package serializer

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestToJSON(t *testing.T) {
	tests := []struct {
		name    string
		data    map[string]interface{}
		indent  bool
		want    string
		wantErr bool
	}{
		{
			name:    "Empty map without indentation",
			data:    map[string]interface{}{},
			indent:  false,
			want:    "{}",
			wantErr: false,
		},
		{
			name:    "Empty map with indentation",
			data:    map[string]interface{}{},
			indent:  true,
			want:    "{}",
			wantErr: false,
		},
		{
			name: "Simple data without indentation",
			data: map[string]interface{}{
				"name": "vendor/project",
				"type": "library",
			},
			indent:  false,
			want:    `{"name":"vendor/project","type":"library"}`,
			wantErr: false,
		},
		{
			name: "Simple data with indentation",
			data: map[string]interface{}{
				"name": "vendor/project",
				"type": "library",
			},
			indent: true,
			// 使用contains检查以避免空格/换行符的问题
			want:    `"name": "vendor/project"`,
			wantErr: false,
		},
		{
			name: "Complex data with nested structures",
			data: map[string]interface{}{
				"name": "vendor/project",
				"require": map[string]interface{}{
					"php": "^7.4",
				},
			},
			indent:  false,
			want:    `{"name":"vendor/project","require":{"php":"^7.4"}}`,
			wantErr: false,
		},
		{
			name: "Data with non-JSON serializable types",
			data: map[string]interface{}{
				"invalid": make(chan int), // 通道无法序列化为JSON
			},
			indent:  false,
			want:    "",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ToJSON(tt.data, tt.indent)

			if (err != nil) != tt.wantErr {
				t.Errorf("ToJSON() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if tt.wantErr {
				return // 如果期望错误，不继续检查输出
			}

			if tt.indent {
				// 对于缩进的情况，检查结果是否包含预期内容
				if !strings.Contains(got, tt.want) {
					t.Errorf("ToJSON() = %v, should contain %v", got, tt.want)
				}
				// 验证缩进格式
				if !strings.Contains(got, "\n") && len(got) > 2 { // 非空JSON应该有换行符
					t.Errorf("ToJSON() with indent=true should contain newlines: %s", got)
				}
			} else {
				// 对于非缩进的情况，精确匹配
				if got != tt.want {
					t.Errorf("ToJSON() = %v, want %v", got, tt.want)
				}
			}
		})
	}
}

func TestSaveToFile(t *testing.T) {
	// 创建测试数据
	testData := map[string]interface{}{
		"name":        "vendor/project",
		"description": "Test project",
		"require": map[string]interface{}{
			"php": "^7.4",
		},
	}

	// 带有无法序列化类型的数据
	invalidData := map[string]interface{}{
		"invalid": make(chan int),
	}

	// 创建临时目录
	tmpDir, err := os.MkdirTemp("", "composer-serializer-test")
	if err != nil {
		t.Fatalf("Failed to create temp directory: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	// 创建一个只读目录用于测试写入失败
	readOnlyDir := filepath.Join(tmpDir, "readonly")
	if err := os.Mkdir(readOnlyDir, 0755); err != nil {
		t.Fatalf("Failed to create readonly directory: %v", err)
	}
	// 尝试将目录设为只读
	if err := os.Chmod(readOnlyDir, 0500); err != nil {
		t.Logf("Warning: Could not set directory to read-only: %v", err)
	}

	// 创建一个文件但无写入权限
	nonWritablePath := filepath.Join(tmpDir, "non-writable.json")
	if err := os.WriteFile(nonWritablePath, []byte("{}"), 0644); err != nil {
		t.Fatalf("Failed to create non-writable file: %v", err)
	}
	// 设置为只读
	if err := os.Chmod(nonWritablePath, 0400); err != nil {
		t.Logf("Warning: Could not set file to read-only: %v", err)
	}

	tests := []struct {
		name     string
		data     map[string]interface{}
		filePath string
		indent   bool
		wantErr  bool
	}{
		{
			name:     "Save to valid path without indentation",
			data:     testData,
			filePath: filepath.Join(tmpDir, "composer1.json"),
			indent:   false,
			wantErr:  false,
		},
		{
			name:     "Save to valid path with indentation",
			data:     testData,
			filePath: filepath.Join(tmpDir, "composer2.json"),
			indent:   true,
			wantErr:  false,
		},
		{
			name:     "Save to invalid path",
			data:     testData,
			filePath: filepath.Join(tmpDir, "nonexistent-dir", "composer.json"),
			indent:   true,
			wantErr:  true,
		},
		{
			name:     "Save invalid data that can't be marshaled",
			data:     invalidData,
			filePath: filepath.Join(tmpDir, "invalid.json"),
			indent:   false,
			wantErr:  true,
		},
		{
			name:     "Save to read-only directory",
			data:     testData,
			filePath: filepath.Join(readOnlyDir, "composer.json"),
			indent:   false,
			wantErr:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := SaveToFile(tt.data, tt.filePath, tt.indent)

			if (err != nil) != tt.wantErr {
				t.Errorf("SaveToFile() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr {
				// 验证文件是否已创建且包含正确的内容
				fileContent, err := os.ReadFile(tt.filePath)
				if err != nil {
					t.Errorf("Failed to read saved file: %v", err)
					return
				}

				jsonStr := string(fileContent)

				// 尝试解析JSON
				var parsed map[string]interface{}
				if err := json.Unmarshal(fileContent, &parsed); err != nil {
					t.Errorf("Saved file contains invalid JSON: %v", err)
					return
				}

				// 验证文件内容是否包含必要的信息
				if parsed["name"] != "vendor/project" {
					t.Errorf("Saved file has incorrect name: %v", parsed["name"])
				}

				// 验证缩进
				if tt.indent && !strings.Contains(jsonStr, "\n") {
					t.Errorf("File should be indented but doesn't contain newlines: %s", jsonStr)
				}
			}
		})
	}
}

func TestCreateBackup(t *testing.T) {
	// 创建临时目录
	tmpDir, err := os.MkdirTemp("", "composer-backup-test")
	if err != nil {
		t.Fatalf("Failed to create temp directory: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	// 创建原始文件
	originalContent := `{"name": "vendor/project"}`
	originalPath := filepath.Join(tmpDir, "composer.json")
	if err := os.WriteFile(originalPath, []byte(originalContent), 0644); err != nil {
		t.Fatalf("Failed to create original file: %v", err)
	}

	// 创建一个文件在只读目录中
	readOnlyDir := filepath.Join(tmpDir, "readonly")
	if err := os.Mkdir(readOnlyDir, 0755); err != nil {
		t.Fatalf("Failed to create readonly directory: %v", err)
	}
	readOnlyFile := filepath.Join(readOnlyDir, "composer.json")
	if err := os.WriteFile(readOnlyFile, []byte(originalContent), 0644); err != nil {
		t.Fatalf("Failed to create file in readonly directory: %v", err)
	}
	// 设置目录为只读
	if err := os.Chmod(readOnlyDir, 0500); err != nil {
		t.Logf("Warning: Could not set directory to read-only: %v", err)
	}

	// 创建不存在的目录路径
	invalidPath := filepath.Join(tmpDir, "nonexistent", "composer.json")

	tests := []struct {
		name         string
		filePath     string
		backupSuffix string
		wantErr      bool
	}{
		{
			name:         "Create backup with default suffix",
			filePath:     originalPath,
			backupSuffix: "",
			wantErr:      false,
		},
		{
			name:         "Create backup with custom suffix",
			filePath:     originalPath,
			backupSuffix: ".backup",
			wantErr:      false,
		},
		{
			name:         "Create backup of non-existent file",
			filePath:     filepath.Join(tmpDir, "nonexistent.json"),
			backupSuffix: "",
			wantErr:      true,
		},
		{
			name:         "Create backup to invalid path",
			filePath:     invalidPath,
			backupSuffix: "",
			wantErr:      true,
		},
		{
			name:         "Create backup to read-only directory",
			filePath:     readOnlyFile,
			backupSuffix: ".bak",
			wantErr:      true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			suffix := tt.backupSuffix
			if suffix == "" {
				suffix = ".bak"
			}

			backupPath, err := CreateBackup(tt.filePath, tt.backupSuffix)

			if (err != nil) != tt.wantErr {
				t.Errorf("CreateBackup() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr {
				// 确认返回的备份路径是正确的
				expectedPath := tt.filePath + suffix
				if backupPath != expectedPath {
					t.Errorf("CreateBackup() backupPath = %v, want %v", backupPath, expectedPath)
				}

				// 确认备份文件存在
				_, err := os.Stat(backupPath)
				if err != nil {
					t.Errorf("Backup file doesn't exist: %v", err)
					return
				}

				// 确认备份文件内容与原始文件相同
				backupContent, err := os.ReadFile(backupPath)
				if err != nil {
					t.Errorf("Failed to read backup file: %v", err)
					return
				}

				if string(backupContent) != originalContent {
					t.Errorf("Backup content = %v, want %v", string(backupContent), originalContent)
				}
			}
		})
	}
}
