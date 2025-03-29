package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/scagogogo/php-composer-json-parser/pkg/composer"
)

func main() {
	fmt.Println("PHP Composer JSON Parser - 基本用法示例")
	fmt.Println("========================================")

	// 创建临时目录保存示例文件
	tmpDir, err := os.MkdirTemp("", "composer-example-*")
	if err != nil {
		log.Fatalf("创建临时目录失败: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	composerJsonPath := filepath.Join(tmpDir, "composer.json")

	// 1. 创建新的composer.json
	fmt.Println("\n1. 创建新的composer.json文件")
	fmt.Println("--------------------------")

	// 创建一个新的库
	myLibrary, err := composer.CreateLibrary(
		"acme/example-library", // 包名
		"PHP Composer解析器示例库",   // 描述
		"^7.4",                 // PHP版本要求
	)
	if err != nil {
		log.Fatalf("创建库失败: %v", err)
	}

	// 添加一些依赖
	myLibrary.AddDependency("symfony/console", "^5.4")
	myLibrary.AddDependency("monolog/monolog", "^2.3")

	// 添加开发依赖
	myLibrary.AddDevDependency("phpunit/phpunit", "^9.5")

	// 保存到文件（美化格式）
	if err := myLibrary.Save(composerJsonPath, true); err != nil {
		log.Fatalf("保存composer.json失败: %v", err)
	}

	fmt.Printf("composer.json文件已保存到: %s\n", composerJsonPath)
	printFileContent(composerJsonPath)

	// 2. 解析现有的composer.json
	fmt.Println("\n2. 解析现有的composer.json文件")
	fmt.Println("-----------------------------")

	parsedComposer, err := composer.ParseFile(composerJsonPath)
	if err != nil {
		log.Fatalf("解析composer.json失败: %v", err)
	}

	fmt.Printf("包名: %s\n", parsedComposer.Name)
	fmt.Printf("描述: %s\n", parsedComposer.Description)
	fmt.Printf("类型: %s\n", parsedComposer.Type)
	fmt.Printf("PHP版本要求: %s\n", parsedComposer.Require["php"])

	fmt.Println("\n依赖项:")
	for pkg, version := range parsedComposer.Require {
		fmt.Printf("- %s: %s\n", pkg, version)
	}

	fmt.Println("\n开发依赖项:")
	for pkg, version := range parsedComposer.RequireDev {
		fmt.Printf("- %s: %s\n", pkg, version)
	}

	// 3. 修改composer.json
	fmt.Println("\n3. 修改composer.json")
	fmt.Println("-------------------")

	// 更新PHP版本要求
	parsedComposer.AddDependency("php", "^8.0")

	// 添加新依赖
	parsedComposer.AddDependency("guzzlehttp/guzzle", "^7.4")

	// 移除依赖
	removed := parsedComposer.RemoveDependency("monolog/monolog")
	fmt.Printf("移除monolog/monolog: %v\n", removed)

	// 打印PSR-4自动加载信息
	psr4Map, ok := parsedComposer.GetPSR4Map()
	if ok {
		fmt.Println("\nPSR-4自动加载配置:")
		for namespace, path := range psr4Map {
			fmt.Printf("- %s => %s\n", namespace, path)
		}
	}

	// 添加新的命名空间映射
	parsedComposer.SetPSR4("Acme\\Example\\Tests\\", "tests/")

	// 保存修改后的文件
	modifiedPath := filepath.Join(tmpDir, "composer.modified.json")
	if err := parsedComposer.Save(modifiedPath, true); err != nil {
		log.Fatalf("保存修改后的composer.json失败: %v", err)
	}

	fmt.Printf("\n修改后的composer.json文件已保存到: %s\n", modifiedPath)
	printFileContent(modifiedPath)

	fmt.Println("\n示例执行完成!")
}

// 辅助函数：打印文件内容
func printFileContent(filePath string) {
	content, err := os.ReadFile(filePath)
	if err != nil {
		fmt.Printf("无法读取文件内容: %v\n", err)
		return
	}

	fmt.Println("\n文件内容:")
	fmt.Println("----------")
	fmt.Println(string(content))
	fmt.Println("----------")
}
