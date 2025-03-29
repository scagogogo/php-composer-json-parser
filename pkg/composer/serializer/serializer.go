// Package serializer 提供将PHP Composer结构体序列化为JSON的功能
package serializer

import (
	"encoding/json"
	"fmt"
	"os"
)

// ToJSON 将map数据结构转换为JSON字符串
//
// 参数:
//   - data: 要转换的数据
//   - indent: 是否缩进格式化JSON（true为美化输出，false为紧凑输出）
//
// 返回:
//   - string: 转换后的JSON字符串
//   - error: 如果转换失败，返回错误
//
// 示例:
//
//	jsonStr, err := serializer.ToJSON(composerData, true)
//	if err != nil {
//		log.Fatal(err)
//	}
//	fmt.Println(jsonStr)
func ToJSON(data map[string]interface{}, indent bool) (string, error) {
	var (
		jsonData []byte
		err      error
	)

	if indent {
		jsonData, err = json.MarshalIndent(data, "", "    ")
	} else {
		jsonData, err = json.Marshal(data)
	}

	if err != nil {
		return "", fmt.Errorf("error marshalling to JSON: %v", err)
	}

	return string(jsonData), nil
}

// SaveToFile 将数据保存为JSON文件
//
// 参数:
//   - data: 要保存的数据
//   - filePath: 保存的文件路径
//   - indent: 是否缩进格式化JSON（true为美化输出，false为紧凑输出）
//
// 返回:
//   - error: 如果保存失败，返回错误
//
// 示例:
//
//	err := serializer.SaveToFile(composerData, "./composer.json", true)
//	if err != nil {
//		log.Fatal(err)
//	}
//	fmt.Println("文件已保存")
func SaveToFile(data map[string]interface{}, filePath string, indent bool) error {
	jsonData, err := ToJSON(data, indent)
	if err != nil {
		return err
	}

	return os.WriteFile(filePath, []byte(jsonData), 0644)
}

// CreateBackup 在修改前创建composer.json的备份
//
// 参数:
//   - filePath: 原始文件路径
//   - backupSuffix: 备份文件后缀（默认为.bak）
//
// 返回:
//   - string: 备份文件路径
//   - error: 如果备份失败，返回错误
func CreateBackup(filePath string, backupSuffix string) (string, error) {
	if backupSuffix == "" {
		backupSuffix = ".bak"
	}

	backupPath := filePath + backupSuffix

	// 读取原始文件
	data, err := os.ReadFile(filePath)
	if err != nil {
		return "", fmt.Errorf("failed to read original file: %v", err)
	}

	// 写入备份文件
	err = os.WriteFile(backupPath, data, 0644)
	if err != nil {
		return "", fmt.Errorf("failed to create backup file: %v", err)
	}

	return backupPath, nil
}
