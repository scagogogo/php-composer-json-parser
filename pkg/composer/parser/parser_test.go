package parser

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"reflect"
	"runtime"
	"strings"
	"testing"
)

func TestParseString(t *testing.T) {
	// 测试 null 值处理，了解预期行为
	testNullHandling(t)

	tests := []struct {
		name     string
		jsonStr  string
		wantData map[string]interface{}
		wantErr  error
	}{
		{
			name:    "Valid simple JSON",
			jsonStr: `{"name": "vendor/project"}`,
			wantData: map[string]interface{}{
				"name": "vendor/project",
			},
			wantErr: nil,
		},
		{
			name: "Valid complex JSON",
			jsonStr: `{
				"name": "vendor/project",
				"description": "Test project",
				"require": {
					"php": "^7.4"
				}
			}`,
			wantData: map[string]interface{}{
				"name":        "vendor/project",
				"description": "Test project",
				"require": map[string]interface{}{
					"php": "^7.4",
				},
			},
			wantErr: nil,
		},
		{
			name:     "Invalid JSON",
			jsonStr:  `{"name": "vendor/project"`,
			wantData: nil,
			wantErr:  ErrInvalidJSON,
		},
		{
			name:     "Empty JSON",
			jsonStr:  `{}`,
			wantData: map[string]interface{}{},
			wantErr:  nil,
		},
		{
			name:     "Empty string",
			jsonStr:  "",
			wantData: nil,
			wantErr:  ErrInvalidJSON,
		},
		{
			name:     "Malformed JSON that doesn't pass validation",
			jsonStr:  `{"name": "\u0G00"}`, // 无效的 Unicode 转义序列
			wantData: nil,
			wantErr:  ErrInvalidJSON,
		},
		// JSON null 值能通过 json.Valid，但 Unmarshal 到 map 时会失败
		{
			name:     "JSON null value",
			jsonStr:  `null`,
			wantData: nil,
			wantErr:  nil, // 不期望错误
		},
		// JSON array 值能通过 json.Valid，但 Unmarshal 到 map 时会失败
		{
			name:     "JSON array instead of object",
			jsonStr:  `[1, 2, 3]`,
			wantData: nil,
			wantErr:  ErrUnmarshallingJSON,
		},
		// JSON boolean 值能通过 json.Valid，但 Unmarshal 到 map 时会失败
		{
			name:     "JSON boolean value",
			jsonStr:  `true`,
			wantData: nil,
			wantErr:  ErrUnmarshallingJSON,
		},
		// JSON 数值能通过 json.Valid，但 Unmarshal 到 map 时会失败
		{
			name:     "JSON numeric value",
			jsonStr:  `42`,
			wantData: nil,
			wantErr:  ErrUnmarshallingJSON,
		},
		{
			name:     "Very large JSON object",
			jsonStr:  generateLargeJSON(),
			wantData: nil, // 不检查具体数据，只关心解析是否成功
			wantErr:  nil,
		},
		{
			name:    "Special characters in JSON string",
			jsonStr: `{"special": "特殊字符 çñáéíóú 😀🔥👍"}`,
			wantData: map[string]interface{}{
				"special": "特殊字符 çñáéíóú 😀🔥👍",
			},
			wantErr: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotData, gotErr := ParseString(tt.jsonStr)

			// 检查错误情况
			if (gotErr != nil) != (tt.wantErr != nil) {
				t.Errorf("ParseString() error = %v, wantErr %v", gotErr, tt.wantErr)
				return
			}

			// 更精确地检查错误类型
			if gotErr != nil && tt.wantErr != nil {
				if !errors.Is(gotErr, tt.wantErr) && !strings.Contains(gotErr.Error(), tt.wantErr.Error()) {
					t.Errorf("ParseString() error = %v, want error containing %v", gotErr, tt.wantErr)
				}
				return
			}

			// 如果期望错误，则不检查数据
			if tt.wantErr != nil {
				return
			}

			// 对于超大JSON，只检查解析是否成功
			if tt.name == "Very large JSON object" {
				if gotData == nil {
					t.Errorf("ParseString() failed to parse large JSON")
				}
				return
			}

			// 检查数据
			if !reflect.DeepEqual(gotData, tt.wantData) {
				t.Errorf("ParseString() gotData = %v, want %v", gotData, tt.wantData)
			}
		})
	}
}

func TestParse(t *testing.T) {
	tests := []struct {
		name     string
		input    io.Reader
		wantData map[string]interface{}
		wantErr  bool
		errType  error
	}{
		{
			name:  "Valid JSON",
			input: strings.NewReader(`{"name": "vendor/project"}`),
			wantData: map[string]interface{}{
				"name": "vendor/project",
			},
			wantErr: false,
		},
		{
			name:     "Invalid JSON",
			input:    strings.NewReader(`{"name": "vendor/project"`),
			wantData: nil,
			wantErr:  true,
			errType:  ErrInvalidJSON,
		},
		{
			name:     "Read error",
			input:    &errorReader{},
			wantData: nil,
			wantErr:  true,
			errType:  ErrReadingFile,
		},
		{
			name:     "Empty JSON",
			input:    strings.NewReader(`{}`),
			wantData: map[string]interface{}{},
			wantErr:  false,
		},
		{
			name:     "Malformed JSON that doesn't pass validation",
			input:    strings.NewReader(`{"name": "\u0G00"}`), // 无效的 Unicode 转义序列
			wantData: nil,
			wantErr:  true,
			errType:  ErrInvalidJSON,
		},
		{
			name:     "Large JSON payload",
			input:    strings.NewReader(generateLargeJSON()),
			wantData: nil, // 我们不关心实际数据，只关心解析是否成功
			wantErr:  false,
		},
		{
			name:     "JSON null value",
			input:    strings.NewReader(`null`),
			wantData: nil,
			wantErr:  false, // 不期望错误
			errType:  nil,
		},
		{
			name:     "JSON array instead of object",
			input:    strings.NewReader(`[1, 2, 3]`),
			wantData: nil,
			wantErr:  true,
			errType:  ErrUnmarshallingJSON,
		},
		{
			name:     "Empty input",
			input:    strings.NewReader(``),
			wantData: nil,
			wantErr:  true,
			errType:  ErrInvalidJSON,
		},
		{
			name:     "Extremely large input (simulated)",
			input:    &limitedReader{content: []byte(`{"name": "test"}`)},
			wantData: map[string]interface{}{"name": "test"},
			wantErr:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotData, gotErr := Parse(tt.input)

			// 检查错误
			if (gotErr != nil) != tt.wantErr {
				t.Errorf("Parse() error = %v, wantErr %v", gotErr, tt.wantErr)
				return
			}

			if tt.wantErr && tt.errType != nil {
				if !errors.Is(gotErr, tt.errType) && !strings.Contains(gotErr.Error(), tt.errType.Error()) {
					t.Errorf("Parse() error type = %v, want error containing %v", gotErr, tt.errType)
				}
				return
			}

			// 对于大型JSON，我们只检查解析是否成功，不比较具体数据
			if tt.name == "Large JSON payload" {
				if gotData == nil {
					t.Errorf("Parse() failed to parse large JSON")
				}
				return
			}

			// 对于模拟极大输入，只检查基本结构
			if tt.name == "Extremely large input (simulated)" {
				if gotData == nil || gotData["name"] != "test" {
					t.Errorf("Parse() failed with simulated large input: %v", gotData)
				}
				return
			}

			// 检查数据
			if !reflect.DeepEqual(gotData, tt.wantData) {
				t.Errorf("Parse() gotData = %v, want %v", gotData, tt.wantData)
			}
		})
	}
}

func TestParseFile(t *testing.T) {
	// 创建测试文件
	tempDir := t.TempDir()
	validFilePath := filepath.Join(tempDir, "valid.json")
	invalidFilePath := filepath.Join(tempDir, "invalid.json")
	nonExistentPath := filepath.Join(tempDir, "nonexistent.json")
	invalidUnicodeFilePath := filepath.Join(tempDir, "invalid_unicode.json")
	largeFilePath := filepath.Join(tempDir, "large.json")
	specialCharsFilePath := filepath.Join(tempDir, "special_chars.json")
	emptyFilePath := filepath.Join(tempDir, "empty.json")
	nonJSONExtFilePath := filepath.Join(tempDir, "composer.txt")
	nullValuePath := filepath.Join(tempDir, "null_value.json")
	arrayValuePath := filepath.Join(tempDir, "array_value.json")
	noPermissionFilePath := filepath.Join(tempDir, "no_permission.json")

	// 写测试文件内容
	err := os.WriteFile(validFilePath, []byte(`{"name": "test/package"}`), 0644)
	if err != nil {
		t.Fatal(err)
	}

	err = os.WriteFile(invalidFilePath, []byte(`{"name": "test/package", "version": "1.0.0"`), 0644)
	if err != nil {
		t.Fatal(err)
	}

	err = os.WriteFile(invalidUnicodeFilePath, []byte(`{"name": "test/package", "description": "\u00g5"}`), 0644)
	if err != nil {
		t.Fatal(err)
	}

	// 创建一个较大的JSON文件
	largeJSON := []byte(`{"name": "test/package", "description": "`)
	largeJSON = append(largeJSON, bytes.Repeat([]byte("a"), 10000)...)
	largeJSON = append(largeJSON, []byte(`"}`)...)
	err = os.WriteFile(largeFilePath, largeJSON, 0644)
	if err != nil {
		t.Fatal(err)
	}

	// 包含特殊字符的JSON
	err = os.WriteFile(specialCharsFilePath, []byte(`{"name": "test/package", "description": "Special chars: éèêë àâäãå"}`), 0644)
	if err != nil {
		t.Fatal(err)
	}

	// 空文件
	err = os.WriteFile(emptyFilePath, []byte{}, 0644)
	if err != nil {
		t.Fatal(err)
	}

	// 没有.json扩展名的文件
	err = os.WriteFile(nonJSONExtFilePath, []byte(`{"name": "test/package"}`), 0644)
	if err != nil {
		t.Fatal(err)
	}

	// null值JSON
	err = os.WriteFile(nullValuePath, []byte(`null`), 0644)
	if err != nil {
		t.Fatal(err)
	}

	// 数组值JSON
	err = os.WriteFile(arrayValuePath, []byte(`[1, 2, 3]`), 0644)
	if err != nil {
		t.Fatal(err)
	}

	// 创建无权限文件
	err = os.WriteFile(noPermissionFilePath, []byte(`{"name": "test/package"}`), 0644)
	if err != nil {
		t.Fatal(err)
	}
	// 在非Windows系统上设置无读取权限
	if runtime.GOOS != "windows" {
		err = os.Chmod(noPermissionFilePath, 0000)
		if err != nil {
			t.Fatal(err)
		}
	}

	tests := []struct {
		name     string
		filePath string
		wantData map[string]interface{}
		wantErr  bool
		errType  error
	}{
		{
			name:     "Valid file",
			filePath: validFilePath,
			wantData: map[string]interface{}{
				"name": "test/package",
			},
			wantErr: false,
		},
		{
			name:     "Invalid JSON file",
			filePath: invalidFilePath,
			wantData: nil,
			wantErr:  true,
			errType:  ErrInvalidJSON,
		},
		{
			name:     "File not found",
			filePath: nonExistentPath,
			wantData: nil,
			wantErr:  true,
			errType:  ErrFileNotFound,
		},
		{
			name:     "Invalid Unicode escape sequence",
			filePath: invalidUnicodeFilePath,
			wantData: nil,
			wantErr:  true,
			errType:  ErrInvalidJSON,
		},
		{
			name:     "Large JSON file",
			filePath: largeFilePath,
			wantData: nil, // 不检查具体数据
			wantErr:  false,
		},
		{
			name:     "Special characters in JSON file",
			filePath: specialCharsFilePath,
			wantData: map[string]interface{}{
				"name":        "test/package",
				"description": "Special chars: éèêë àâäãå",
			},
			wantErr: false,
		},
		{
			name:     "Empty file",
			filePath: emptyFilePath,
			wantData: nil,
			wantErr:  true,
			errType:  ErrInvalidJSON,
		},
		{
			name:     "File without .json extension",
			filePath: nonJSONExtFilePath,
			wantData: map[string]interface{}{
				"name": "test/package",
			},
			wantErr: false,
		},
		{
			name:     "File with JSON null value",
			filePath: nullValuePath,
			wantData: nil,
			wantErr:  false, // 不期望错误
			errType:  nil,
		},
		{
			name:     "File with JSON array value",
			filePath: arrayValuePath,
			wantData: nil,
			wantErr:  true,
			errType:  ErrUnmarshallingJSON,
		},
		{
			name:     "File with no read permission",
			filePath: noPermissionFilePath,
			wantData: nil,
			wantErr:  true,
			errType:  ErrReadingFile,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotData, gotErr := ParseFile(tt.filePath)

			// 检查错误
			if (gotErr != nil) != tt.wantErr {
				t.Errorf("ParseFile() error = %v, wantErr %v", gotErr, tt.wantErr)
				return
			}

			if tt.wantErr && tt.errType != nil {
				if !errors.Is(gotErr, tt.errType) && !strings.Contains(gotErr.Error(), tt.errType.Error()) {
					t.Errorf("ParseFile() error type = %v, want error containing %v", gotErr, tt.errType)
				}
				return
			}

			// 如果是大型JSON，只检查解析是否成功
			if tt.name == "Large JSON file" {
				if gotData == nil {
					t.Errorf("ParseFile() failed to parse large JSON file")
				}
				return
			}

			// 检查数据
			if !reflect.DeepEqual(gotData, tt.wantData) {
				t.Errorf("ParseFile() gotData = %v, want %v", gotData, tt.wantData)
			}
		})
	}
}

func TestParseDir(t *testing.T) {
	// 创建临时目录
	tmpDir, err := os.MkdirTemp("", "composer-test")
	if err != nil {
		t.Fatalf("Failed to create temp directory: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	// 创建composer.json文件
	validJSON := `{"name": "vendor/project"}`
	composerPath := filepath.Join(tmpDir, "composer.json")
	if err := os.WriteFile(composerPath, []byte(validJSON), 0644); err != nil {
		t.Fatalf("Failed to create composer.json file: %v", err)
	}

	// 创建不包含composer.json的空目录
	emptyDir, err := os.MkdirTemp("", "composer-empty")
	if err != nil {
		t.Fatalf("Failed to create empty directory: %v", err)
	}
	defer os.RemoveAll(emptyDir)

	tests := []struct {
		name     string
		dir      string
		wantData map[string]interface{}
		wantErr  bool
		errType  error
	}{
		{
			name: "Directory with composer.json",
			dir:  tmpDir,
			wantData: map[string]interface{}{
				"name": "vendor/project",
			},
			wantErr: false,
		},
		{
			name:     "Directory without composer.json",
			dir:      emptyDir,
			wantData: nil,
			wantErr:  true,
			errType:  ErrFileNotFound,
		},
		{
			name:     "Non-existent directory",
			dir:      filepath.Join(tmpDir, "nonexistent"),
			wantData: nil,
			wantErr:  true,
			errType:  ErrFileNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotData, gotErr := ParseDir(tt.dir)

			// 检查错误
			if (gotErr != nil) != tt.wantErr {
				t.Errorf("ParseDir() error = %v, wantErr %v", gotErr, tt.wantErr)
				return
			}

			if tt.wantErr && tt.errType != nil {
				if !errors.Is(gotErr, tt.errType) && !strings.Contains(gotErr.Error(), tt.errType.Error()) {
					t.Errorf("ParseDir() error type = %v, want error containing %v", gotErr, tt.errType)
				}
				return
			}

			// 检查数据
			if !reflect.DeepEqual(gotData, tt.wantData) {
				t.Errorf("ParseDir() gotData = %v, want %v", gotData, tt.wantData)
			}
		})
	}
}

// errorReader 是一个始终返回错误的io.Reader实现
type errorReader struct{}

func (e *errorReader) Read(p []byte) (n int, err error) {
	return 0, io.ErrUnexpectedEOF
}

// generateLargeJSON 生成一个大型JSON对象用于测试
func generateLargeJSON() string {
	var sb strings.Builder
	sb.WriteString(`{"name":"vendor/large-project","description":"A very large project","require":{`)

	// 添加许多依赖项
	for i := 0; i < 100; i++ {
		if i > 0 {
			sb.WriteString(",")
		}
		sb.WriteString(fmt.Sprintf(`"package%d":"^1.0.%d"`, i, i))
	}

	sb.WriteString(`},"repositories":[`)

	// 添加多个仓库
	for i := 0; i < 20; i++ {
		if i > 0 {
			sb.WriteString(",")
		}
		sb.WriteString(fmt.Sprintf(`{"type":"vcs","url":"https://github.com/vendor/repo%d"}`, i))
	}

	sb.WriteString(`]}`)
	return sb.String()
}

// limitedReader 是一个模拟读取大文件的 io.Reader 实现
type limitedReader struct {
	content []byte
	pos     int
}

func (r *limitedReader) Read(p []byte) (n int, err error) {
	if r.pos >= len(r.content) {
		return 0, io.EOF
	}

	n = copy(p, r.content[r.pos:])
	r.pos += n
	return n, nil
}

// 扩展testNullHandling以测试更多JSON基本类型
func testNullHandling(t *testing.T) {
	t.Helper()

	// 测试各种基本类型的JSON值
	testCases := []struct {
		name    string
		jsonStr string
	}{
		{"null", "null"},
		{"array", "[1,2,3]"},
		{"boolean", "true"},
		{"number", "42"},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			data := []byte(tc.jsonStr)

			// 检查json.Valid
			if !json.Valid(data) {
				t.Errorf("Expected %s to be valid JSON, but json.Valid returned false", tc.jsonStr)
			}

			// 尝试解析为map
			var mapResult map[string]interface{}
			mapErr := json.Unmarshal(data, &mapResult)

			// 尝试解析为interface{}
			var anyResult interface{}
			anyErr := json.Unmarshal(data, &anyResult)

			t.Logf("%s - To map: result=%v, err=%v", tc.name, mapResult, mapErr)
			t.Logf("%s - To interface{}: result=%v, err=%v", tc.name, anyResult, anyErr)
		})
	}
}

// TestParseFile_CorruptContent 测试损坏的文件内容
func TestParseFile_CorruptContent(t *testing.T) {
	// 创建临时目录和文件
	tempDir := t.TempDir()
	corruptFilePath := filepath.Join(tempDir, "corrupt.json")

	// 写入损坏的内容（不是有效的UTF-8编码）
	err := os.WriteFile(corruptFilePath, []byte{0xFF, 0xFE, 0xFD}, 0644)
	if err != nil {
		t.Fatal(err)
	}

	// 测试解析损坏的文件
	_, err = ParseFile(corruptFilePath)
	if err == nil {
		t.Error("Expected error when parsing corrupt file, but got nil")
	}
}

// TestParseFile_DirectoryAsFile 测试尝试将目录作为文件解析
func TestParseFile_DirectoryAsFile(t *testing.T) {
	// 创建临时目录
	tempDir := t.TempDir()

	// 尝试将目录作为文件解析
	_, err := ParseFile(tempDir)
	if err == nil {
		t.Error("Expected error when parsing directory as file, but got nil")
	}
}
