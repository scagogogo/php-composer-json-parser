// Package parser 提供解析PHP Composer JSON文件的功能
package parser

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
)

// ComposerJSON结构将在parser.Parse等函数返回时被父包转换为composer.ComposerJSON
type ComposerJSON map[string]interface{}

// 解析器错误定义
var (
	// ErrInvalidJSON 表示JSON格式无效
	ErrInvalidJSON = fmt.Errorf("invalid JSON format")

	// ErrFileNotFound 表示composer.json文件未找到
	ErrFileNotFound = fmt.Errorf("composer.json file not found")

	// ErrReadingFile 表示读取文件时出错
	ErrReadingFile = fmt.Errorf("error reading file")

	// ErrUnmarshallingJSON 表示JSON反序列化时出错
	ErrUnmarshallingJSON = fmt.Errorf("error unmarshalling JSON")
)

// ParseFile 从文件路径解析composer.json文件
//
// 参数:
//   - filePath: composer.json文件路径
//
// 返回:
//   - map[string]interface{}: 解析后的原始JSON数据
//   - error: 如果解析失败，返回错误
//
// 示例:
//
//	rawData, err := parser.ParseFile("./composer.json")
//	if err != nil {
//		log.Fatal(err)
//	}
func ParseFile(filePath string) (map[string]interface{}, error) {
	// 检查文件是否存在
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		return nil, fmt.Errorf("%w: %s", ErrFileNotFound, filePath)
	}

	// 打开文件
	file, err := os.Open(filePath)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrReadingFile, err)
	}
	defer file.Close()

	return Parse(file)
}

// ParseDir 在指定目录中查找并解析composer.json文件
//
// 参数:
//   - dir: 要查找composer.json的目录路径
//
// 返回:
//   - map[string]interface{}: 解析后的原始JSON数据
//   - error: 如果解析失败，返回错误
func ParseDir(dir string) (map[string]interface{}, error) {
	filePath := filepath.Join(dir, "composer.json")
	return ParseFile(filePath)
}

// Parse 从io.Reader读取JSON并解析为原始map结构
//
// 参数:
//   - r: io.Reader接口，可以是文件、字符串等
//
// 返回:
//   - map[string]interface{}: 解析后的原始JSON数据
//   - error: 如果解析失败，返回错误
func Parse(r io.Reader) (map[string]interface{}, error) {
	// 读取所有数据
	data, err := io.ReadAll(r)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrReadingFile, err)
	}

	// 验证JSON
	if !json.Valid(data) {
		return nil, ErrInvalidJSON
	}

	var result map[string]interface{}
	if err := json.Unmarshal(data, &result); err != nil {
		return nil, fmt.Errorf("%w: %v", ErrUnmarshallingJSON, err)
	}

	return result, nil
}

// ParseString 解析composer.json字符串
//
// 参数:
//   - jsonStr: 要解析的JSON字符串
//
// 返回:
//   - map[string]interface{}: 解析后的原始JSON数据
//   - error: 如果解析失败，返回错误
func ParseString(jsonStr string) (map[string]interface{}, error) {
	// 验证JSON
	if !json.Valid([]byte(jsonStr)) {
		return nil, ErrInvalidJSON
	}

	var result map[string]interface{}
	if err := json.Unmarshal([]byte(jsonStr), &result); err != nil {
		return nil, fmt.Errorf("%w: %v", ErrUnmarshallingJSON, err)
	}

	return result, nil
}
