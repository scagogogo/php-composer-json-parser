# PHP Composer JSON Parser 使用示例

本目录包含了一系列示例，展示了如何使用 `php-composer-json-parser` 库来创建、解析、修改和保存 PHP Composer 的 `composer.json` 文件。

## 示例目录

### 1. 基本用法 (01_basic_usage)

展示库的基本功能，包括：
- 创建新的 composer.json 文件
- 解析现有的 composer.json 文件
- 修改配置并保存

**运行方法：**
```bash
cd 01_basic_usage
go run main.go
```

### 2. 项目创建 (02_project_creation)

展示如何创建不同类型的 PHP 项目：
- Web 应用项目
- PHP 库/组件
- API 工具项目

**运行方法：**
```bash
cd 02_project_creation
go run main.go
```

### 3. 依赖管理 (03_dependencies)

专注于展示依赖项管理的功能：
- 解析和查看依赖
- 检查依赖是否存在
- 添加新依赖
- 更新依赖版本
- 移除依赖
- 合并运行时和开发依赖

**运行方法：**
```bash
cd 03_dependencies
go run main.go
```

### 4. 高级功能 (04_advanced)

展示更高级的功能，包括：
- PSR-4 自动加载配置
- 仓库管理
- 归档排除配置
- 配置选项设置
- 作者和支持信息

**运行方法：**
```bash
cd 04_advanced
go run main.go
```

## 注意事项

- 这些示例都是独立运行的，每个示例都会在临时目录中创建和操作文件，不会影响您的实际项目
- 示例代码包含详细的注释，解释每个操作的目的和效果
- 运行结果会直接输出到控制台，包括创建的 composer.json 文件内容

## 自定义示例

您可以通过修改这些示例代码，尝试不同的参数和配置，以便更好地了解库的功能。

例如，可以尝试：
- 在 `01_basic_usage` 示例中，修改添加的依赖版本
- 在 `02_project_creation` 示例中，创建其他类型的项目
- 在 `03_dependencies` 示例中，添加和移除自定义依赖
- 在 `04_advanced` 示例中，尝试其他 PSR-4 命名空间映射 