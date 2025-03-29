// Package composer 提供用于解析和操作PHP Composer JSON文件的功能。
//
// 此包允许开发者：
// - 将composer.json文件解析为Go结构体
// - 创建新的composer.json配置
// - 验证包名
// - 管理依赖项（添加、删除、检查）
// - 操作PSR-0/PSR-4自动加载配置
// - 将结构体转换回JSON
// - 保存到文件
package composer

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/scagogogo/php-composer-json-parser/pkg/composer/archive"
	"github.com/scagogogo/php-composer-json-parser/pkg/composer/autoload"
	"github.com/scagogogo/php-composer-json-parser/pkg/composer/config"
	"github.com/scagogogo/php-composer-json-parser/pkg/composer/dependency"
	"github.com/scagogogo/php-composer-json-parser/pkg/composer/parser"
	"github.com/scagogogo/php-composer-json-parser/pkg/composer/repository"
	"github.com/scagogogo/php-composer-json-parser/pkg/composer/serializer"
)

// 包版本
const Version = "1.0.0"

// 错误定义
var (
	// ErrInvalidJSON 表示JSON格式无效
	ErrInvalidJSON = parser.ErrInvalidJSON

	// ErrFileNotFound 表示composer.json文件未找到
	ErrFileNotFound = parser.ErrFileNotFound

	// ErrReadingFile 表示读取文件时出错
	ErrReadingFile = parser.ErrReadingFile

	// ErrUnmarshallingJSON 表示JSON反序列化时出错
	ErrUnmarshallingJSON = parser.ErrUnmarshallingJSON
)

// ParseFile 从文件路径解析composer.json文件
//
// 参数:
//   - filePath: composer.json文件路径
//
// 返回:
//   - *ComposerJSON: 解析后的结构体
//   - error: 如果解析失败，返回错误
//
// 示例:
//
//	composer, err := composer.ParseFile("./composer.json")
//	if err != nil {
//		log.Fatal(err)
//	}
//	fmt.Println("包名:", composer.Name)
//	fmt.Println("版本:", composer.Version)
//	fmt.Println("PHP依赖版本:", composer.Require["php"])
func ParseFile(filePath string) (*ComposerJSON, error) {
	rawData, err := parser.ParseFile(filePath)
	if err != nil {
		return nil, err
	}

	return convertToComposerJSON(rawData)
}

// ParseDir 在指定目录中查找并解析composer.json文件
//
// 参数:
//   - dir: 要查找composer.json的目录路径
//
// 返回:
//   - *ComposerJSON: 解析后的结构体
//   - error: 如果解析失败，返回错误
//
// 示例:
//
//	// 解析当前项目的composer.json
//	composer, err := composer.ParseDir(".")
//	if err != nil {
//		log.Fatal(err)
//	}
//	// 解析指定PHP项目的composer.json
//	composer, err = composer.ParseDir("/path/to/php/project")
//	if err != nil {
//		log.Fatal(err)
//	}
func ParseDir(dir string) (*ComposerJSON, error) {
	filePath := filepath.Join(dir, "composer.json")
	return ParseFile(filePath)
}

// Parse 从io.Reader读取JSON并解析为ComposerJSON结构体
//
// 参数:
//   - r: io.Reader接口，可以是文件、字符串等
//
// 返回:
//   - *ComposerJSON: 解析后的结构体
//   - error: 如果解析失败，返回错误
//
// 示例:
//
//	file, err := os.Open("composer.json")
//	if err != nil {
//		log.Fatal(err)
//	}
//	defer file.Close()
//
//	composer, err := composer.Parse(file)
//	if err != nil {
//		log.Fatal(err)
//	}
//
//	// 也可以从HTTP响应体解析
//	resp, err := http.Get("https://example.com/composer.json")
//	if err != nil {
//		log.Fatal(err)
//	}
//	defer resp.Body.Close()
//
//	composer, err = composer.Parse(resp.Body)
func Parse(r io.Reader) (*ComposerJSON, error) {
	rawData, err := parser.Parse(r)
	if err != nil {
		return nil, err
	}

	return convertToComposerJSON(rawData)
}

// ParseString 解析composer.json字符串
//
// 参数:
//   - jsonStr: 要解析的JSON字符串
//
// 返回:
//   - *ComposerJSON: 解析后的结构体
//   - error: 如果解析失败，返回错误
//
// 示例:
//
//	jsonStr := `{
//		"name": "vendor/project",
//		"description": "My composer package",
//		"require": {
//			"php": ">=7.4"
//		}
//	}`
//
//	composer, err := composer.ParseString(jsonStr)
//	if err != nil {
//		log.Fatal(err)
//	}
//	fmt.Println("包名:", composer.Name)
func ParseString(jsonStr string) (*ComposerJSON, error) {
	rawData, err := parser.ParseString(jsonStr)
	if err != nil {
		return nil, err
	}

	return convertToComposerJSON(rawData)
}

// convertToComposerJSON 将原始map转换为ComposerJSON结构体
//
// 参数:
//   - data: 从解析JSON获得的原始map
//
// 返回:
//   - *ComposerJSON: 转换后的结构体
//   - error: 如果转换失败，返回错误
//
// 注意: 这是内部函数，用于将解析器返回的通用map转换为强类型的ComposerJSON结构体
func convertToComposerJSON(data map[string]interface{}) (*ComposerJSON, error) {
	// 重新序列化为JSON然后反序列化为ComposerJSON结构体
	// 这是处理嵌套结构的有效方法
	jsonData, err := json.Marshal(data)
	if err != nil {
		return nil, fmt.Errorf("error marshalling raw data: %v", err)
	}

	var composer ComposerJSON
	if err := json.Unmarshal(jsonData, &composer); err != nil {
		return nil, fmt.Errorf("error converting to ComposerJSON: %v", err)
	}

	return &composer, nil
}

// ToJSON 将ComposerJSON结构体转换为JSON字符串
//
// 参数:
//   - indent: 是否缩进格式化JSON（true为美化输出，false为紧凑输出）
//
// 返回:
//   - string: 转换后的JSON字符串
//   - error: 如果转换失败，返回错误
//
// 示例:
//
//	composer, _ := composer.ParseFile("composer.json")
//
//	// 美化输出（带缩进）
//	prettyJSON, err := composer.ToJSON(true)
//	if err != nil {
//		log.Fatal(err)
//	}
//	fmt.Println(prettyJSON)
//
//	// 紧凑输出（不带缩进）
//	compactJSON, err := composer.ToJSON(false)
//	if err != nil {
//		log.Fatal(err)
//	}
func (c *ComposerJSON) ToJSON(indent bool) (string, error) {
	// 将结构体转换为map
	jsonData, err := json.Marshal(c)
	if err != nil {
		return "", fmt.Errorf("error marshalling to JSON: %v", err)
	}

	var rawData map[string]interface{}
	if err := json.Unmarshal(jsonData, &rawData); err != nil {
		return "", fmt.Errorf("error unmarshalling to map: %v", err)
	}

	return serializer.ToJSON(rawData, indent)
}

// Save 将ComposerJSON结构体保存为JSON文件
//
// 参数:
//   - filePath: 保存的文件路径
//   - indent: 是否缩进格式化JSON（true为美化输出，false为紧凑输出）
//
// 返回:
//   - error: 如果保存失败，返回错误
//
// 示例:
//
//	// 创建新的composer.json
//	composer, _ := composer.CreateNew("vendor/project", "我的PHP项目")
//
//	// 添加依赖
//	composer.AddDependency("php", ">=7.4")
//	composer.AddDependency("symfony/console", "^5.4")
//
//	// 保存到文件（带缩进美化格式）
//	err := composer.Save("./composer.json", true)
//	if err != nil {
//		log.Fatal(err)
//	}
func (c *ComposerJSON) Save(filePath string, indent bool) error {
	jsonData, err := c.ToJSON(indent)
	if err != nil {
		return err
	}

	return os.WriteFile(filePath, []byte(jsonData), 0644)
}

// CreateBackup 在修改前创建composer.json的备份
//
// 参数:
//   - filePath: 原始文件路径
//   - backupSuffix: 备份文件后缀（为空时默认使用.bak）
//
// 返回:
//   - string: 备份文件路径
//   - error: 如果备份失败，返回错误
//
// 示例:
//
//	// 在修改前备份composer.json
//	backupPath, err := composer.CreateBackup("./composer.json", "")
//	if err != nil {
//		log.Fatal(err)
//	}
//	fmt.Printf("已创建备份: %s\n", backupPath) // 输出：已创建备份: ./composer.json.bak
//
//	// 使用自定义后缀
//	backupPath, err = composer.CreateBackup("./composer.json", ".backup")
//	if err != nil {
//		log.Fatal(err)
//	}
//	fmt.Printf("已创建备份: %s\n", backupPath) // 输出：已创建备份: ./composer.json.backup
func CreateBackup(filePath string, backupSuffix string) (string, error) {
	return serializer.CreateBackup(filePath, backupSuffix)
}

// DependencyExists 检查依赖项是否存在于require部分
//
// 参数:
//   - packageName: 要检查的包名，如"symfony/console"
//
// 返回:
//   - bool: 如果依赖项存在返回true，否则返回false
//
// 示例:
//
//	composer, _ := composer.ParseFile("./composer.json")
//
//	if composer.DependencyExists("php") {
//		fmt.Println("PHP依赖存在，版本为:", composer.Require["php"])
//	} else {
//		fmt.Println("PHP依赖不存在")
//	}
func (c *ComposerJSON) DependencyExists(packageName string) bool {
	return dependency.DependencyExists(c.Require, packageName)
}

// DevDependencyExists 检查依赖项是否存在于require-dev部分
//
// 参数:
//   - packageName: 要检查的包名，如"phpunit/phpunit"
//
// 返回:
//   - bool: 如果开发依赖项存在返回true，否则返回false
//
// 示例:
//
//	composer, _ := composer.ParseFile("./composer.json")
//
//	if composer.DevDependencyExists("phpunit/phpunit") {
//		fmt.Println("PHPUnit依赖存在，版本为:", composer.RequireDev["phpunit/phpunit"])
//	} else {
//		fmt.Println("PHPUnit依赖不存在")
//	}
func (c *ComposerJSON) DevDependencyExists(packageName string) bool {
	return dependency.DependencyExists(c.RequireDev, packageName)
}

// AddDependency 向require部分添加包
//
// 参数:
//   - packageName: 要添加的包名，格式为"vendor/package"
//   - version: 依赖版本，如"^5.4"、">=7.4"
//
// 返回:
//   - error: 如果添加失败，返回错误
//
// 示例:
//
//	composer, _ := composer.ParseFile("./composer.json")
//
//	// 添加PHP版本需求
//	err := composer.AddDependency("php", ">=7.4")
//	if err != nil {
//		log.Fatal(err)
//	}
//
//	// 添加Symfony组件
//	err = composer.AddDependency("symfony/console", "^5.4")
//	if err != nil {
//		log.Fatal(err)
//	}
//
//	// 保存修改
//	composer.Save("./composer.json", true)
func (c *ComposerJSON) AddDependency(packageName, version string) error {
	return dependency.AddDependency(c.Require, packageName, version)
}

// AddDevDependency 向require-dev部分添加包
//
// 参数:
//   - packageName: 要添加的包名，格式为"vendor/package"
//   - version: 依赖版本，如"^9.0"
//
// 返回:
//   - error: 如果添加失败，返回错误
//
// 示例:
//
//	composer, _ := composer.ParseFile("./composer.json")
//
//	// 添加测试依赖
//	err := composer.AddDevDependency("phpunit/phpunit", "^9.0")
//	if err != nil {
//		log.Fatal(err)
//	}
//
//	// 添加代码风格检查工具
//	err = composer.AddDevDependency("squizlabs/php_codesniffer", "^3.6")
//	if err != nil {
//		log.Fatal(err)
//	}
//
//	// 保存修改
//	composer.Save("./composer.json", true)
func (c *ComposerJSON) AddDevDependency(packageName, version string) error {
	return dependency.AddDependency(c.RequireDev, packageName, version)
}

// RemoveDependency 从require部分移除包
//
// 参数:
//   - packageName: 要移除的包名
//
// 返回:
//   - bool: 如果成功移除返回true，如果包不存在返回false
//
// 示例:
//
//	composer, _ := composer.ParseFile("./composer.json")
//
//	if composer.RemoveDependency("symfony/console") {
//		fmt.Println("成功移除symfony/console")
//	} else {
//		fmt.Println("symfony/console不存在于依赖中")
//	}
//
//	// 保存修改
//	composer.Save("./composer.json", true)
func (c *ComposerJSON) RemoveDependency(packageName string) bool {
	return dependency.RemoveDependency(c.Require, packageName)
}

// RemoveDevDependency 从require-dev部分移除包
//
// 参数:
//   - packageName: 要移除的包名
//
// 返回:
//   - bool: 如果成功移除返回true，如果包不存在返回false
//
// 示例:
//
//	composer, _ := composer.ParseFile("./composer.json")
//
//	if composer.RemoveDevDependency("phpunit/phpunit") {
//		fmt.Println("成功移除phpunit/phpunit")
//	} else {
//		fmt.Println("phpunit/phpunit不存在于开发依赖中")
//	}
//
//	// 保存修改
//	composer.Save("./composer.json", true)
func (c *ComposerJSON) RemoveDevDependency(packageName string) bool {
	return dependency.RemoveDependency(c.RequireDev, packageName)
}

// GetAllDependencies 返回所有依赖项（require和require-dev合并）
//
// 返回:
//   - map[string]string: 包含所有依赖的映射，key为包名，value为版本需求
//
// 示例:
//
//	composer, _ := composer.ParseFile("./composer.json")
//
//	allDeps := composer.GetAllDependencies()
//	fmt.Println("项目共有", len(allDeps), "个依赖项:")
//
//	for pkg, version := range allDeps {
//		fmt.Printf("- %s: %s\n", pkg, version)
//	}
func (c *ComposerJSON) GetAllDependencies() map[string]string {
	return dependency.MergeDependencies(c.Require, c.RequireDev)
}

// GetPSR4Map 获取PSR-4自动加载命名空间映射
//
// 返回:
//   - map[string]string: 命名空间到目录的映射，key为命名空间，value为目录路径
//   - bool: 是否成功获取映射，如果PSR-4配置不存在或无效则返回false
//
// 示例:
//
//	composer, _ := composer.ParseFile("./composer.json")
//
//	psr4Map, ok := composer.GetPSR4Map()
//	if !ok {
//		fmt.Println("PSR-4配置不存在或无效")
//		return
//	}
//
//	fmt.Println("PSR-4自动加载映射:")
//	for namespace, path := range psr4Map {
//		fmt.Printf("- %s => %s\n", namespace, path)
//	}
func (c *ComposerJSON) GetPSR4Map() (map[string]string, bool) {
	return autoload.GetPSR4Map(&c.Autoload)
}

// SetPSR4 设置PSR-4命名空间映射
//
// 参数:
//   - namespace: 命名空间，必须以\\结尾，如"App\\"
//   - path: 目录路径，如"src/"
//
// 示例:
//
//	composer, _ := composer.ParseFile("./composer.json")
//
//	// 设置命名空间映射
//	composer.SetPSR4("App\\", "src/")
//	composer.SetPSR4("App\\Tests\\", "tests/")
//
//	// 保存修改
//	composer.Save("./composer.json", true)
//
//	// composer.json中会包含以下内容:
//	// "autoload": {
//	//   "psr-4": {
//	//     "App\\": "src/",
//	//     "App\\Tests\\": "tests/"
//	//   }
//	// }
func (c *ComposerJSON) SetPSR4(namespace, path string) {
	autoload.SetPSR4(&c.Autoload, namespace, path)
}

// RemovePSR4 移除PSR-4命名空间映射
//
// 参数:
//   - namespace: 要移除的命名空间，必须以\\结尾，如"App\\"
//
// 返回:
//   - bool: 如果成功移除返回true，如果命名空间不存在返回false
//
// 示例:
//
//	composer, _ := composer.ParseFile("./composer.json")
//
//	if composer.RemovePSR4("App\\Tests\\") {
//		fmt.Println("成功移除App\\Tests\\命名空间映射")
//	} else {
//		fmt.Println("App\\Tests\\命名空间映射不存在")
//	}
//
//	// 保存修改
//	composer.Save("./composer.json", true)
func (c *ComposerJSON) RemovePSR4(namespace string) bool {
	return autoload.RemovePSR4(&c.Autoload, namespace)
}

// AddExclusion 向归档排除列表添加路径模式
//
// 参数:
//   - pattern: 要排除的路径模式，如"/tests"、"/.github"
//
// 示例:
//
//	composer, _ := composer.ParseFile("./composer.json")
//
//	// 排除测试目录和文档
//	composer.AddExclusion("/tests")
//	composer.AddExclusion("/docs")
//	composer.AddExclusion("/.github")
//
//	// 保存修改
//	composer.Save("./composer.json", true)
//
//	// composer.json中会包含以下内容:
//	// "archive": {
//	//   "exclude": [
//	//     "/tests",
//	//     "/docs",
//	//     "/.github"
//	//   ]
//	// }
func (c *ComposerJSON) AddExclusion(pattern string) {
	archive.AddExclusion(&c.Archive, pattern)
}

// RemoveExclusion 从归档排除列表中删除路径模式
//
// 参数:
//   - pattern: 要移除的路径模式，如"/tests"
//
// 返回:
//   - bool: 如果成功移除返回true，如果模式不存在返回false
//
// 示例:
//
//	composer, _ := composer.ParseFile("./composer.json")
//
//	if composer.RemoveExclusion("/tests") {
//		fmt.Println("成功移除/tests排除模式")
//	} else {
//		fmt.Println("/tests排除模式不存在")
//	}
//
//	// 保存修改
//	composer.Save("./composer.json", true)
func (c *ComposerJSON) RemoveExclusion(pattern string) bool {
	return archive.RemoveExclusion(&c.Archive, pattern)
}

// AddRepository 添加一个仓库到仓库列表
//
// 参数:
//   - repo: 要添加的仓库对象
//
// 示例:
//
//	composer, _ := composer.ParseFile("./composer.json")
//
//	// 添加自定义Composer仓库
//	repo := composer.NewRepository("composer", "https://asset-packagist.org")
//	composer.AddRepository(*repo)
//
//	// 添加版本控制系统仓库
//	vcsRepo := composer.NewRepository("vcs", "https://github.com/myvendor/private-package")
//	composer.AddRepository(*vcsRepo)
//
//	// 保存修改
//	composer.Save("./composer.json", true)
func (c *ComposerJSON) AddRepository(repo repository.Repository) {
	c.Repositories = append(c.Repositories, repo)
}

// NewRepository 创建一个新的仓库
//
// 参数:
//   - repoType: 仓库类型，如"composer"、"vcs"、"path"
//   - url: 仓库URL或路径
//
// 返回:
//   - *repository.Repository: 新创建的仓库对象
//
// 示例:
//
//	// 创建Packagist仓库
//	packagist := composer.NewRepository("composer", "https://packagist.org")
//
//	// 创建GitHub仓库
//	githubRepo := composer.NewRepository("vcs", "https://github.com/symfony/console")
//
//	// 创建本地路径仓库
//	localRepo := composer.NewRepository("path", "../my-local-package")
//
//	// 向Composer配置添加仓库
//	composer, _ := composer.ParseFile("./composer.json")
//	composer.AddRepository(*packagist)
//	composer.Save("./composer.json", true)
func NewRepository(repoType, url string) *repository.Repository {
	return repository.NewRepository(repoType, url)
}

// DefaultConfig 返回带有合理默认值的Config
//
// 返回:
//   - *config.Config: 一个配置好标准默认值的Config对象
//
// 示例:
//
//	// 创建一个新的composer.json文件
//	composer, _ := composer.CreateNew("vendor/project", "我的PHP项目")
//
//	// 使用默认配置
//	composer.Config = *composer.DefaultConfig()
//
//	// 修改特定配置
//	composer.Config.ProcessTimeout = 600 // 设置为10分钟
//	composer.Config.VendorDir = "lib"    // 更改vendor目录
//
//	// 保存修改
//	composer.Save("./composer.json", true)
func DefaultConfig() *config.Config {
	return config.DefaultConfig()
}
