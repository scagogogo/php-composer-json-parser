# PHP Composer JSON 解析器

[![Go](https://github.com/scagogogo/php-composer-json-parser/actions/workflows/go.yml/badge.svg)](https://github.com/scagogogo/php-composer-json-parser/actions/workflows/go.yml)
[![GoDoc](https://godoc.org/github.com/scagogogo/php-composer-json-parser?status.svg)](https://pkg.go.dev/github.com/scagogogo/php-composer-json-parser)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)
[![Go Report Card](https://goreportcard.com/badge/github.com/scagogogo/php-composer-json-parser)](https://goreportcard.com/report/github.com/scagogogo/php-composer-json-parser)
[![Coverage](https://img.shields.io/badge/coverage-97%25-brightgreen.svg)](https://github.com/scagogogo/php-composer-json-parser/actions)

这是一个Go语言实现的PHP Composer配置文件（`composer.json`）解析和操作工具库。通过简洁直观的API，开发者可以轻松创建、解析、修改和保存PHP项目的`composer.json`文件，适用于自动化工具、CI/CD流程和PHP项目管理工具。

## 🌟 功能特点

- **全面的文件操作**
  - 解析`composer.json`文件并转换为结构化的Go对象
  - 创建新的`composer.json`配置（支持项目、库等不同类型）
  - 将修改后的配置保存回文件，支持格式化和备份
  
- **严格的验证功能**
  - 校验包名格式（`vendor/project`）
  - 验证版本约束格式（语义化版本）
  - 确保生成的配置符合Composer规范
  
- **依赖项管理**
  - 添加、更新和删除运行时及开发依赖
  - 查询依赖是否存在及其版本
  - 合并和过滤依赖列表
  
- **高级配置支持**
  - PSR-4/PSR-0 自动加载配置
  - 自定义仓库管理
  - 配置选项设置（如缓存目录、超时等）
  - 归档排除规则
  - 作者和支持信息

- **优秀的开发体验**
  - 详细的中文注释
  - 完整的测试覆盖（97%+）
  - 类型安全的API设计
  - 丰富的示例代码

## 📦 安装

```bash
go get github.com/scagogogo/php-composer-json-parser
```

需要 Go 1.18 或更高版本。

## 🚀 快速开始

### 解析 composer.json 文件

```go
package main

import (
	"fmt"
	"log"
	
	"github.com/scagogogo/php-composer-json-parser/pkg/composer"
)

func main() {
	// 从文件解析
	c, err := composer.ParseFile("./composer.json")
	if err != nil {
		log.Fatalf("解析失败: %v", err)
	}
	
	fmt.Printf("包名: %s\n", c.Name)
	fmt.Printf("描述: %s\n", c.Description)
	fmt.Printf("类型: %s\n", c.Type)
	
	// 查看依赖
	fmt.Println("依赖项:")
	for pkg, version := range c.Require {
		fmt.Printf("- %s: %s\n", pkg, version)
	}
}
```

### 创建新的 composer.json

```go
package main

import (
	"fmt"
	"log"
	
	"github.com/scagogogo/php-composer-json-parser/pkg/composer"
)

func main() {
	// 创建一个PHP库项目
	c, err := composer.CreateLibrary(
		"vendor/my-library",
		"一个实用的PHP库",
		"^8.0", // PHP版本要求
	)
	if err != nil {
		log.Fatalf("创建失败: %v", err)
	}
	
	// 添加依赖
	c.AddDependency("symfony/console", "^6.0")
	c.AddDependency("monolog/monolog", "^2.3")
	
	// 添加开发依赖
	c.AddDevDependency("phpunit/phpunit", "^9.5")
	
	// 保存到文件（美化格式）
	err = c.Save("./composer.json", true)
	if err != nil {
		log.Fatalf("保存失败: %v", err)
	}
	
	fmt.Println("composer.json文件已创建")
}
```

## 📚 详细文档

### 解析API

从多种来源解析 composer.json：

```go
// 从文件解析
composer, err := composer.ParseFile("./composer.json")

// 从目录解析（自动寻找目录下的composer.json文件）
composer, err := composer.ParseDir("./project-dir")

// 从字符串解析
jsonStr := `{"name": "vendor/project", "require": {"php": "^7.4"}}`
composer, err := composer.ParseString(jsonStr)

// 从io.Reader解析
file, _ := os.Open("composer.json")
composer, err := composer.Parse(file)
```

### 创建API

创建不同类型的 composer.json：

```go
// 创建基本的composer.json
composer, err := composer.CreateNew("vendor/package", "描述")

// 创建PHP项目
composer, err := composer.CreateProject("vendor/project", "项目描述", "^8.0")

// 创建PHP库
composer, err := composer.CreateLibrary("vendor/library", "库描述", "^7.4")
```

### 依赖管理

管理项目依赖项：

```go
// 检查依赖是否存在
if composer.DependencyExists("symfony/console") {
    fmt.Println("已安装Symfony Console")
}

if composer.DevDependencyExists("phpunit/phpunit") {
    fmt.Println("已安装PHPUnit作为开发依赖")
}

// 添加依赖
composer.AddDependency("php", "^8.0")
composer.AddDependency("monolog/monolog", "^2.3")

// 添加开发依赖
composer.AddDevDependency("phpunit/phpunit", "^9.5")
composer.AddDevDependency("phpstan/phpstan", "^1.4")

// 移除依赖
removed := composer.RemoveDependency("old/package")
removed = composer.RemoveDevDependency("outdated/tool")

// 获取所有依赖（合并require和require-dev）
allDeps := composer.GetAllDependencies()
```

### PSR-4 自动加载

配置PSR-4自动加载：

```go
// 获取当前PSR-4配置
psr4Map, ok := composer.GetPSR4Map()
if ok {
    for namespace, path := range psr4Map {
        fmt.Printf("%s => %s\n", namespace, path)
    }
}

// 设置PSR-4命名空间
composer.SetPSR4("App\\", "src/")
composer.SetPSR4("App\\Tests\\", "tests/")

// 移除PSR-4命名空间
removed := composer.RemovePSR4("App\\Utils\\")
```

### 归档排除

配置打包时要排除的文件和目录：

```go
// 添加排除模式
composer.AddExclusion("/tests")
composer.AddExclusion("/docs")
composer.AddExclusion("/.github")

// 移除排除模式
composer.RemoveExclusion("/docs")
```

### 仓库管理

添加不同类型的仓库：

```go
// 创建仓库
packagist := composer.NewRepository("composer", "https://packagist.org")
vcsRepo := composer.NewRepository("vcs", "https://github.com/vendor/package")
pathRepo := composer.NewRepository("path", "../local-package")

// 添加仓库
composer.AddRepository(*packagist)
composer.AddRepository(*vcsRepo)
```

### 输出和保存

```go
// 转换为JSON字符串（美化格式）
jsonStr, err := composer.ToJSON(true)
fmt.Println(jsonStr)

// 转换为JSON字符串（紧凑格式）
jsonStr, err := composer.ToJSON(false)

// 保存到文件（美化格式）
err = composer.Save("./composer.json", true)

// 在修改前创建备份
backupPath, err := composer.CreateBackup("./composer.json", ".bak")
fmt.Printf("备份文件创建在: %s\n", backupPath)
```

### 配置管理

使用和修改Composer配置选项：

```go
// 获取默认配置
config := composer.DefaultConfig()

// 修改配置
composer.Config.ProcessTimeout = 600 // 10分钟超时
composer.Config.OptimizeAutoloader = true
composer.Config.PreferredInstall = "dist"
composer.Config.VendorDir = "vendors" // 自定义vendor目录
```

### 错误处理

库使用特定错误类型帮助识别问题：

```go
c, err := composer.ParseFile("./composer.json")
if err != nil {
    switch {
    case errors.Is(err, composer.ErrFileNotFound):
        fmt.Println("文件不存在，将创建新文件")
        // 创建新文件的代码
    case errors.Is(err, composer.ErrInvalidJSON):
        fmt.Println("JSON格式无效")
    case errors.Is(err, composer.ErrUnmarshallingJSON):
        fmt.Println("JSON结构不匹配")
    default:
        fmt.Printf("未知错误: %v\n", err)
    }
}
```

## 🔍 包结构

该项目采用模块化设计，将不同功能分解到子包中：

- `pkg/composer`: 主包，提供高级API
  - `pkg/composer/archive`: 存档相关功能
  - `pkg/composer/autoload`: 自动加载配置
  - `pkg/composer/config`: 配置相关功能
  - `pkg/composer/dependency`: 依赖项管理
  - `pkg/composer/parser`: JSON解析功能
  - `pkg/composer/repository`: 仓库管理
  - `pkg/composer/serializer`: JSON序列化
  - `pkg/composer/validation`: 数据验证

## 📋 示例代码

完整的示例代码可以在 [examples](examples/) 目录中找到：

1. [基本用法](examples/01_basic_usage/main.go) - 演示基本的解析和创建功能
2. [项目创建](examples/02_project_creation/main.go) - 演示如何创建不同类型的PHP项目
3. [依赖管理](examples/03_dependencies/main.go) - 演示依赖项管理功能
4. [高级功能](examples/04_advanced/main.go) - 演示PSR-4、仓库管理等高级功能

查看 [examples/README.md](examples/README.md) 获取更多关于示例的信息。

## 🔄 集成测试

项目使用GitHub Actions自动运行单元测试和示例代码。每次提交代码时，会执行以下操作：

1. 运行所有单元测试并检查覆盖率
2. 执行所有示例代码确保功能正常
3. 进行代码格式检查和静态分析

## 🛠️ 用法场景

本库适用于多种场景：

- **PHP项目生成工具** - 自动创建新的PHP项目结构
- **依赖分析工具** - 分析PHP项目依赖关系
- **CI/CD流程** - 自动化PHP项目的构建和部署
- **项目迁移工具** - 帮助从一个框架迁移到另一个框架
- **代码质量工具** - 检查项目的依赖项是否过时或存在漏洞

## 👥 贡献指南

我们欢迎各种形式的贡献！如果您想参与贡献，请遵循以下步骤：

1. Fork本仓库
2. 创建您的特性分支 (`git checkout -b feature/amazing-feature`)
3. 提交您的更改 (`git commit -m 'Add some amazing feature'`)
4. 推送到分支 (`git push origin feature/amazing-feature`)
5. 创建新的Pull Request

请确保您的代码通过测试并且遵循项目的代码风格。

## 📝 许可证

本项目使用MIT许可证 - 详情请参阅[LICENSE](LICENSE)文件。 