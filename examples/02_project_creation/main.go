package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/scagogogo/php-composer-json-parser/pkg/composer"
)

func main() {
	fmt.Println("PHP Composer JSON Parser - 项目创建示例")
	fmt.Println("========================================")

	// 创建临时目录保存示例文件
	tmpDir, err := os.MkdirTemp("", "composer-project-example-*")
	if err != nil {
		log.Fatalf("创建临时目录失败: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	// 1. 创建PHP Web应用项目
	fmt.Println("\n1. 创建PHP Web应用项目")
	fmt.Println("----------------------")

	// 使用CreateProject创建一个Web应用
	webApp, err := composer.CreateProject(
		"acme/web-app", // 包名
		"Acme公司Web应用",  // 描述
		"^8.1",         // PHP版本要求
	)
	if err != nil {
		log.Fatalf("创建Web应用项目失败: %v", err)
	}

	// 设置Web应用特有的配置
	webApp.Type = "project" // 项目类型

	// 添加Web开发常用的依赖
	webApp.AddDependency("symfony/framework-bundle", "^6.0")
	webApp.AddDependency("symfony/yaml", "^6.0")
	webApp.AddDependency("doctrine/orm", "^2.11")
	webApp.AddDependency("twig/twig", "^3.3")

	// 添加开发依赖
	webApp.AddDevDependency("symfony/debug-bundle", "^6.0")
	webApp.AddDevDependency("symfony/maker-bundle", "^1.36")
	webApp.AddDevDependency("phpstan/phpstan", "^1.4")

	// 保存Web应用的composer.json
	webAppPath := filepath.Join(tmpDir, "web-app-composer.json")
	if err := webApp.Save(webAppPath, true); err != nil {
		log.Fatalf("保存Web应用composer.json失败: %v", err)
	}

	fmt.Printf("Web应用composer.json已保存到: %s\n", webAppPath)
	printFileContent(webAppPath)

	// 2. 创建PHP库/组件
	fmt.Println("\n2. 创建PHP库/组件")
	fmt.Println("----------------")

	// 使用CreateLibrary创建一个可复用库
	library, err := composer.CreateLibrary(
		"acme/data-validator", // 包名
		"Acme数据验证库",           // 描述
		"^7.4",                // PHP版本要求
	)
	if err != nil {
		log.Fatalf("创建库项目失败: %v", err)
	}

	// 设置库特有的配置
	library.License = "MIT"                                    // 许可证
	library.Keywords = []string{"validation", "forms", "data"} // 关键词

	// 添加库所需依赖
	library.AddDependency("symfony/validator", "^5.4")

	// 添加开发依赖
	library.AddDevDependency("phpunit/phpunit", "^9.5")
	library.AddDevDependency("phpstan/phpstan", "^1.4")

	// 设置PSR-4自动加载
	library.SetPSR4("Acme\\DataValidator\\", "src/")
	library.SetPSR4("Acme\\DataValidator\\Tests\\", "tests/")

	// 保存库的composer.json
	libraryPath := filepath.Join(tmpDir, "library-composer.json")
	if err := library.Save(libraryPath, true); err != nil {
		log.Fatalf("保存库composer.json失败: %v", err)
	}

	fmt.Printf("库composer.json已保存到: %s\n", libraryPath)
	printFileContent(libraryPath)

	// 3. 创建API工具项目
	fmt.Println("\n3. 创建API工具项目")
	fmt.Println("------------------")

	// 使用CreateNew创建一个API工具项目
	apiTool, err := composer.CreateNew(
		"acme/api-tools", // 包名
		"Acme API工具集",    // 描述
	)
	if err != nil {
		log.Fatalf("创建API工具项目失败: %v", err)
	}

	// 自定义配置
	apiTool.Type = "project"
	apiTool.License = "proprietary" // 私有许可

	// 添加基本依赖
	apiTool.AddDependency("php", "^8.0")
	apiTool.AddDependency("guzzlehttp/guzzle", "^7.4")
	apiTool.AddDependency("symfony/http-client", "^6.0")
	apiTool.AddDependency("monolog/monolog", "^2.3")

	// 创建作者信息
	apiTool.Authors = []composer.Author{
		{
			Name:     "Acme开发团队",
			Email:    "dev@acme.example.com",
			Homepage: "https://acme.example.com",
			Role:     "Developer",
		},
	}

	// 设置支持信息
	apiTool.Support = composer.Support{
		Email:  "support@acme.example.com",
		Issues: "https://github.com/acme/api-tools/issues",
		Docs:   "https://docs.acme.example.com",
	}

	// 保存API工具的composer.json
	apiToolPath := filepath.Join(tmpDir, "api-tool-composer.json")
	if err := apiTool.Save(apiToolPath, true); err != nil {
		log.Fatalf("保存API工具composer.json失败: %v", err)
	}

	fmt.Printf("API工具composer.json已保存到: %s\n", apiToolPath)
	printFileContent(apiToolPath)

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
