package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/scagogogo/php-composer-json-parser/pkg/composer"
)

func main() {
	fmt.Println("PHP Composer JSON Parser - 依赖管理示例")
	fmt.Println("========================================")

	// 创建临时目录保存示例文件
	tmpDir, err := os.MkdirTemp("", "composer-deps-example-*")
	if err != nil {
		log.Fatalf("创建临时目录失败: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	// 创建示例composer.json文件
	composerJsonPath := filepath.Join(tmpDir, "composer.json")

	// 使用字符串创建一个现有的composer.json
	composerJsonContent := `{
		"name": "acme/dependency-example",
		"description": "依赖管理示例",
		"type": "project",
		"require": {
			"php": "^7.4",
			"symfony/console": "^5.4",
			"monolog/monolog": "^2.3",
			"doctrine/orm": "^2.10"
		},
		"require-dev": {
			"phpunit/phpunit": "^9.5"
		},
		"license": "MIT"
	}`

	// 保存示例文件
	if err := os.WriteFile(composerJsonPath, []byte(composerJsonContent), 0644); err != nil {
		log.Fatalf("无法写入示例composer.json文件: %v", err)
	}

	// 1. 解析composer.json文件
	fmt.Println("\n1. 解析composer.json文件中的依赖")
	fmt.Println("-------------------------------")

	c, err := composer.ParseFile(composerJsonPath)
	if err != nil {
		log.Fatalf("解析composer.json失败: %v", err)
	}

	// 显示原始依赖
	fmt.Println("当前运行时依赖:")
	for pkg, version := range c.Require {
		fmt.Printf("- %s: %s\n", pkg, version)
	}

	fmt.Println("\n当前开发依赖:")
	for pkg, version := range c.RequireDev {
		fmt.Printf("- %s: %s\n", pkg, version)
	}

	// 2. 检查特定依赖是否存在
	fmt.Println("\n2. 检查依赖是否存在")
	fmt.Println("------------------")

	pkgsToCheck := []string{
		"php",
		"symfony/console",
		"symfony/framework-bundle", // 不存在
		"phpunit/phpunit",          // 开发依赖
	}

	for _, pkg := range pkgsToCheck {
		if c.DependencyExists(pkg) {
			fmt.Printf("依赖 '%s' 存在于运行时依赖中，版本: %s\n", pkg, c.Require[pkg])
		} else if c.DevDependencyExists(pkg) {
			fmt.Printf("依赖 '%s' 存在于开发依赖中，版本: %s\n", pkg, c.RequireDev[pkg])
		} else {
			fmt.Printf("依赖 '%s' 不存在\n", pkg)
		}
	}

	// 3. 添加依赖
	fmt.Println("\n3. 添加新依赖")
	fmt.Println("-------------")

	// 添加运行时依赖
	depsToAdd := map[string]string{
		"symfony/http-client": "^5.4",
		"guzzlehttp/guzzle":   "^7.4",
		"twig/twig":           "^3.3",
	}

	for pkg, version := range depsToAdd {
		if err := c.AddDependency(pkg, version); err != nil {
			fmt.Printf("添加 %s 失败: %v\n", pkg, err)
		} else {
			fmt.Printf("添加 %s: %s 成功\n", pkg, version)
		}
	}

	// 添加开发依赖
	devDepsToAdd := map[string]string{
		"phpstan/phpstan":           "^1.4",
		"squizlabs/php_codesniffer": "^3.6",
	}

	for pkg, version := range devDepsToAdd {
		if err := c.AddDevDependency(pkg, version); err != nil {
			fmt.Printf("添加 %s 失败: %v\n", pkg, err)
		} else {
			fmt.Printf("添加开发依赖 %s: %s 成功\n", pkg, version)
		}
	}

	// 4. 更新依赖版本
	fmt.Println("\n4. 更新依赖版本")
	fmt.Println("---------------")

	// 更新PHP版本要求
	oldVersion := c.Require["php"]
	if err := c.AddDependency("php", "^8.0"); err != nil {
		fmt.Printf("更新PHP版本失败: %v\n", err)
	} else {
		fmt.Printf("PHP版本从 %s 更新到 %s\n", oldVersion, c.Require["php"])
	}

	// 更新Symfony Console版本
	oldVersion = c.Require["symfony/console"]
	if err := c.AddDependency("symfony/console", "^6.0"); err != nil {
		fmt.Printf("更新symfony/console版本失败: %v\n", err)
	} else {
		fmt.Printf("symfony/console版本从 %s 更新到 %s\n", oldVersion, c.Require["symfony/console"])
	}

	// 5. 移除依赖
	fmt.Println("\n5. 移除依赖")
	fmt.Println("------------")

	depsToRemove := []string{"monolog/monolog", "doctrine/orm", "non-existent/package"}

	for _, pkg := range depsToRemove {
		if c.RemoveDependency(pkg) {
			fmt.Printf("依赖 '%s' 已成功移除\n", pkg)
		} else {
			fmt.Printf("依赖 '%s' 移除失败（可能不存在）\n", pkg)
		}
	}

	// 移除开发依赖
	if c.RemoveDevDependency("phpunit/phpunit") {
		fmt.Printf("开发依赖 'phpunit/phpunit' 已成功移除\n")
	}

	// 6. 获取所有依赖（合并运行时和开发依赖）
	fmt.Println("\n6. 获取所有依赖（合并运行时和开发依赖）")
	fmt.Println("---------------------------------------")

	allDeps := c.GetAllDependencies()
	fmt.Printf("总计 %d 个依赖:\n", len(allDeps))

	for pkg, version := range allDeps {
		fmt.Printf("- %s: %s\n", pkg, version)
	}

	// 7. 保存修改后的composer.json
	fmt.Println("\n7. 保存修改后的composer.json")
	fmt.Println("---------------------------")

	modifiedPath := filepath.Join(tmpDir, "composer.modified.json")
	if err := c.Save(modifiedPath, true); err != nil {
		log.Fatalf("保存修改后的composer.json失败: %v", err)
	}

	fmt.Printf("修改后的composer.json已保存到: %s\n", modifiedPath)
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
