package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/scagogogo/php-composer-json-parser/pkg/composer"
)

func main() {
	fmt.Println("PHP Composer JSON Parser - 高级功能示例")
	fmt.Println("========================================")

	// 创建临时目录保存示例文件
	tmpDir, err := os.MkdirTemp("", "composer-advanced-example-*")
	if err != nil {
		log.Fatalf("创建临时目录失败: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	// 创建基础Composer项目
	composerJsonPath := filepath.Join(tmpDir, "composer.json")

	// 1. PSR-4自动加载配置
	fmt.Println("\n1. PSR-4自动加载配置")
	fmt.Println("------------------")

	// 创建一个新的库
	c, err := composer.CreateNew(
		"acme/advanced-example",
		"展示高级Composer功能的示例",
	)
	if err != nil {
		log.Fatalf("创建Composer配置失败: %v", err)
	}

	// 基本配置
	c.Type = "project"
	c.License = "MIT"
	c.AddDependency("php", "^8.0")

	// 查看默认的PSR-4配置
	psr4Map, ok := c.GetPSR4Map()
	if ok {
		fmt.Println("默认PSR-4配置:")
		for namespace, path := range psr4Map {
			fmt.Printf("- %s => %s\n", namespace, path)
		}
	}

	// 设置多个PSR-4命名空间映射
	fmt.Println("\n添加更多PSR-4配置:")

	// 添加核心命名空间
	c.SetPSR4("Acme\\Advanced\\", "src/")

	// 添加测试命名空间
	c.SetPSR4("Acme\\Advanced\\Tests\\", "tests/")

	// 添加特定模块命名空间
	c.SetPSR4("Acme\\Advanced\\Api\\", "src/Api/")
	c.SetPSR4("Acme\\Advanced\\Utils\\", "src/Utils/")

	// 检查更新后的PSR-4映射
	psr4Map, _ = c.GetPSR4Map()
	for namespace, path := range psr4Map {
		fmt.Printf("- %s => %s\n", namespace, path)
	}

	// 移除一个命名空间映射
	if c.RemovePSR4("Acme\\Advanced\\Utils\\") {
		fmt.Printf("\n已移除 Acme\\Advanced\\Utils\\ 命名空间\n")
	}

	// 2. 仓库管理
	fmt.Println("\n2. 仓库管理")
	fmt.Println("----------")

	// 创建并添加不同类型的仓库

	// 添加Packagist仓库
	packagist := composer.NewRepository("composer", "https://packagist.org")
	c.AddRepository(*packagist)
	fmt.Printf("添加了Packagist仓库: %s\n", packagist.URL)

	// 添加Composer私有仓库
	privateRepo := composer.NewRepository("composer", "https://composer.example.com")
	c.AddRepository(*privateRepo)
	fmt.Printf("添加了私有Composer仓库: %s\n", privateRepo.URL)

	// 添加版本控制仓库
	vcsRepo := composer.NewRepository("vcs", "https://github.com/acme/private-package")
	c.AddRepository(*vcsRepo)
	fmt.Printf("添加了VCS仓库: %s\n", vcsRepo.URL)

	// 添加本地路径仓库
	pathRepo := composer.NewRepository("path", "../local-package")
	c.AddRepository(*pathRepo)
	fmt.Printf("添加了本地路径仓库: %s\n", pathRepo.URL)

	// 3. 归档排除配置
	fmt.Println("\n3. 归档排除配置")
	fmt.Println("-------------")

	// 添加要从归档中排除的文件和目录
	excludePatterns := []string{
		"/tests",
		"/docs",
		"/.github",
		"/phpunit.xml.dist",
		"/phpstan.neon",
		"/.gitignore",
		"/.editorconfig",
	}

	fmt.Println("添加归档排除模式:")
	for _, pattern := range excludePatterns {
		c.AddExclusion(pattern)
		fmt.Printf("- %s\n", pattern)
	}

	// 移除一个排除模式
	c.RemoveExclusion("/docs")
	fmt.Println("\n移除了 /docs 排除模式")

	// 4. 配置选项
	fmt.Println("\n4. 配置选项")
	fmt.Println("-----------")

	// 获取默认配置
	defaultConfig := composer.DefaultConfig()
	c.Config = *defaultConfig

	fmt.Println("默认配置值:")
	fmt.Printf("- process-timeout: %d秒\n", c.Config.ProcessTimeout)
	fmt.Printf("- vendor-dir: %s\n", c.Config.VendorDir)
	fmt.Printf("- bin-dir: %s\n", c.Config.BinDir)

	// 修改配置
	c.Config.ProcessTimeout = 600 // 增加到10分钟
	c.Config.OptimizeAutoloader = true
	c.Config.PreferredInstall = "dist"

	fmt.Println("\n修改后的配置值:")
	fmt.Printf("- process-timeout: %d秒\n", c.Config.ProcessTimeout)
	fmt.Printf("- optimize-autoloader: %v\n", c.Config.OptimizeAutoloader)
	fmt.Printf("- preferred-install: %s\n", c.Config.PreferredInstall)

	// 5. 作者和支持信息
	fmt.Println("\n5. 作者和支持信息")
	fmt.Println("-----------------")

	// 添加作者信息
	c.Authors = []composer.Author{
		{
			Name:     "张三",
			Email:    "zhangsan@example.com",
			Homepage: "https://zhangsan.example.com",
			Role:     "Developer",
		},
		{
			Name:  "李四",
			Email: "lisi@example.com",
			Role:  "Project Manager",
		},
	}

	// 添加支持信息
	c.Support = composer.Support{
		Email:  "support@acme.example.com",
		Issues: "https://github.com/acme/advanced-example/issues",
		Forum:  "https://forum.acme.example.com",
		Wiki:   "https://wiki.acme.example.com",
		Source: "https://github.com/acme/advanced-example",
		Docs:   "https://docs.acme.example.com",
	}

	// 保存最终的composer.json
	if err := c.Save(composerJsonPath, true); err != nil {
		log.Fatalf("保存composer.json失败: %v", err)
	}

	fmt.Printf("\n高级示例的composer.json已保存到: %s\n", composerJsonPath)
	printFileContent(composerJsonPath)

	// 6. 解析保存的文件（确认配置完整性）
	fmt.Println("\n6. 解析保存的文件（确认配置完整性）")
	fmt.Println("-------------------------------")

	parsed, err := composer.ParseFile(composerJsonPath)
	if err != nil {
		log.Fatalf("解析保存的composer.json失败: %v", err)
	}

	fmt.Println("成功解析配置文件!")
	fmt.Printf("- 包名: %s\n", parsed.Name)
	fmt.Printf("- 作者数量: %d\n", len(parsed.Authors))
	fmt.Printf("- 仓库数量: %d\n", len(parsed.Repositories))

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
