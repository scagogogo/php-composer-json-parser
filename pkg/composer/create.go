package composer

import (
	"github.com/scagogogo/php-composer-json-parser/pkg/composer/autoload"
	"github.com/scagogogo/php-composer-json-parser/pkg/composer/dependency"
	"github.com/scagogogo/php-composer-json-parser/pkg/composer/validation"
)

// CreateNew 创建一个新的空composer.json结构体，并设置基本字段
//
// 参数:
//   - name: 包名，格式为"vendor/project"
//   - description: 包描述
//
// 返回:
//   - *ComposerJSON: 创建的结构体
//   - error: 如果创建失败，返回错误
//
// 默认值:
//   - Type: "library"
//   - License: "MIT"
//   - Autoload.PSR4: 基于包名自动生成，如"Vendor\\Project\\"对应"src/"
//
// 示例:
//
//	// 创建一个基本的Composer包
//	composer, err := composer.CreateNew("vendor/project", "一个示例PHP包")
//	if err != nil {
//		log.Fatal(err)
//	}
//
//	// 查看生成的命名空间
//	psr4, ok := composer.GetPSR4Map()
//	if ok {
//		fmt.Println("PSR-4命名空间:", psr4) // 输出: PSR-4命名空间: map[Vendor\Project\:src/]
//	}
//
//	// 保存到文件
//	composer.Save("composer.json", true)
func CreateNew(name, description string) (*ComposerJSON, error) {
	// 验证包名（如果提供）
	if name != "" {
		if err := dependency.ValidatePackageName(name); err != nil {
			return nil, err
		}
	}

	// 验证基本信息
	if err := validation.ValidateComposerJSON(name, description, ""); err != nil {
		return nil, err
	}

	// 创建PSR-4自动加载映射
	psr4Map := make(map[string]interface{})

	// 如果有包名，为其生成默认的命名空间映射
	if name != "" {
		vendorName, projectName, _ := dependency.GetPackageNameParts(name)
		namespace := toNamespace(vendorName, projectName)
		psr4Map[namespace+"\\"] = "src/"
	}

	// 创建新结构体
	return &ComposerJSON{
		Name:        name,
		Description: description,
		Type:        "library",
		License:     "MIT",
		Require:     make(map[string]string),
		RequireDev:  make(map[string]string),
		Autoload: autoload.Autoload{
			PSR4: psr4Map,
		},
	}, nil
}

// CreateProject 创建一个新的PHP项目composer.json结构体
//
// 参数:
//   - name: 包名，格式为"vendor/project"
//   - description: 包描述
//   - phpVersion: PHP版本要求，例如"^7.4"（为空时默认使用"^7.4"）
//
// 返回:
//   - *ComposerJSON: 创建的结构体
//   - error: 如果创建失败，返回错误
//
// 特点:
//   - Type设置为"project"
//   - 自动添加phpunit/phpunit作为开发依赖
//   - 基于包名自动生成PSR-4命名空间映射
//
// 示例:
//
//	// 创建一个PHP项目
//	composer, err := composer.CreateProject("acme/blog", "Acme博客应用", "^8.0")
//	if err != nil {
//		log.Fatal(err)
//	}
//
//	// 添加框架依赖
//	composer.AddDependency("symfony/framework-bundle", "^5.4")
//	composer.AddDependency("doctrine/orm", "^2.10")
//
//	// 添加更多开发依赖
//	composer.AddDevDependency("symfony/phpunit-bridge", "^5.4")
//
//	// 保存到文件
//	composer.Save("composer.json", true)
func CreateProject(name, description, phpVersion string) (*ComposerJSON, error) {
	composer, err := CreateNew(name, description)
	if err != nil {
		return nil, err
	}

	// 设置项目类型
	composer.Type = "project"

	// 添加PHP版本要求
	if phpVersion != "" {
		composer.Require["php"] = phpVersion
	} else {
		composer.Require["php"] = "^7.4"
	}

	// 添加常用开发依赖
	composer.RequireDev["phpunit/phpunit"] = "^9.0"

	return composer, nil
}

// CreateLibrary 创建一个新的PHP库composer.json结构体
//
// 参数:
//   - name: 包名，格式为"vendor/project"
//   - description: 包描述
//   - phpVersion: PHP版本要求，例如"^7.4"（为空时默认使用"^7.4"）
//
// 返回:
//   - *ComposerJSON: 创建的结构体
//   - error: 如果创建失败，返回错误
//
// 特点:
//   - Type设置为"library"
//   - 适合创建可复用的PHP库/组件
//   - 基于包名自动生成PSR-4命名空间映射
//
// 示例:
//
//	// 创建一个PHP库
//	composer, err := composer.CreateLibrary("acme/utils", "Acme实用工具集", "^7.4")
//	if err != nil {
//		log.Fatal(err)
//	}
//
//	// 添加库依赖
//	composer.AddDependency("symfony/http-foundation", "^5.4")
//
//	// 设置更多自动加载规则
//	composer.SetPSR4("Acme\\Utils\\Tests\\", "tests/")
//
//	// 保存到文件
//	composer.Save("composer.json", true)
//
//	// 库的composer.json会包含类似以下内容:
//	// {
//	//   "name": "acme/utils",
//	//   "description": "Acme实用工具集",
//	//   "type": "library",
//	//   "require": {
//	//     "php": "^7.4",
//	//     "symfony/http-foundation": "^5.4"
//	//   },
//	//   "autoload": {
//	//     "psr-4": {
//	//       "Acme\\Utils\\": "src/",
//	//       "Acme\\Utils\\Tests\\": "tests/"
//	//     }
//	//   }
//	// }
func CreateLibrary(name, description, phpVersion string) (*ComposerJSON, error) {
	composer, err := CreateNew(name, description)
	if err != nil {
		return nil, err
	}

	// 设置库类型
	composer.Type = "library"

	// 添加PHP版本要求
	if phpVersion != "" {
		composer.Require["php"] = phpVersion
	} else {
		composer.Require["php"] = "^7.4"
	}

	return composer, nil
}

// toNamespace 将vendor和project名称转换为符合PSR-4的命名空间
//
// 参数:
//   - vendor: 供应商名称
//   - project: 项目名称
//
// 返回:
//   - string: 转换后的命名空间（不含结尾的反斜杠）
//
// 说明:
//   - 将供应商和项目名称首字母大写，并用反斜杠连接
//   - 例如："symfony"和"console"转换为"Symfony\\Console"
//
// 示例:
//
//	namespace := toNamespace("acme", "blog")
//	fmt.Println(namespace) // 输出: Acme\Blog
//
//	// 在Composer文件中使用时会添加结尾反斜杠
//	namespaceKey := namespace + "\\"
//	// namespaceKey == "Acme\Blog\"
func toNamespace(vendor, project string) string {
	return ucfirst(vendor) + "\\" + ucfirst(project)
}

// ucfirst 将字符串的首字母转换为大写
//
// 参数:
//   - s: 要处理的字符串
//
// 返回:
//   - string: 首字母大写后的字符串
//
// 说明:
//   - 如果字符串为空，返回空字符串
//   - 只转换第一个字符，其余字符保持不变
//   - 仅支持ASCII字母a-z的转换
//
// 示例:
//
//	fmt.Println(ucfirst("symfony")) // 输出: Symfony
//	fmt.Println(ucfirst("acme"))    // 输出: Acme
//	fmt.Println(ucfirst(""))        // 输出: 空字符串
//	fmt.Println(ucfirst("123abc"))  // 输出: 123abc (数字不变)
func ucfirst(s string) string {
	if s == "" {
		return ""
	}
	r := []rune(s)
	if r[0] >= 'a' && r[0] <= 'z' {
		r[0] = r[0] - 'a' + 'A'
	}
	return string(r)
}

// ValidateComposerJSON 验证ComposerJSON结构体
//
// 参数:
//   - name: 包名，格式为"vendor/project"
//   - description: 包描述
//   - stability: 稳定性标识，如"stable"、"beta"、"alpha"、"dev"（可为空）
//
// 返回:
//   - error: 如果验证失败，返回错误；验证通过返回nil
//
// 验证规则:
//   - 包名必须符合"vendor/project"格式
//   - 描述不能为空
//   - 如果提供了稳定性标识，必须是有效值
//
// 示例:
//
//	// 验证基本信息
//	err := composer.ValidateComposerJSON("acme/blog", "Acme博客应用", "")
//	if err != nil {
//		fmt.Println("验证失败:", err)
//	} else {
//		fmt.Println("验证通过")
//	}
//
//	// 验证带稳定性标识
//	err = composer.ValidateComposerJSON("acme/blog", "Acme博客应用", "beta")
//	if err != nil {
//		fmt.Println("验证失败:", err)
//	} else {
//		fmt.Println("验证通过")
//	}
func ValidateComposerJSON(name, description, stability string) error {
	return validation.ValidateComposerJSON(name, description, stability)
}

// ValidateVersion 验证版本字符串是否符合语义化版本规范
//
// 参数:
//   - version: 要验证的版本字符串，如"1.0.0"、"^2.1"、">=7.4"
//
// 返回:
//   - error: 如果验证失败，返回错误；验证通过返回nil
//
// 支持的格式:
//   - 精确版本: "1.0.0"
//   - 范围版本: ">=1.0.0"、"<=2.0.0"、">1.0.0 <2.0.0"
//   - 通配符: "1.0.*"
//   - 赋值符: "^1.0.0"（兼容1.x.x）、"~1.0.0"（兼容1.0.x）
//   - 稳定性标识: "1.0.0-beta"、"1.0.0-RC1"
//
// 示例:
//
//	// 验证不同类型的版本字符串
//	versions := []string{
//		"1.0.0",       // 精确版本
//		"^1.0.0",      // 赋值符版本
//		">=7.4",       // 范围版本
//		"1.0.*",       // 通配符版本
//		"1.0.0-beta1", // 带稳定性标识的版本
//		"invalid",     // 无效版本
//	}
//
//	for _, v := range versions {
//		err := composer.ValidateVersion(v)
//		if err != nil {
//			fmt.Printf("版本'%s'无效: %v\n", v, err)
//		} else {
//			fmt.Printf("版本'%s'有效\n", v)
//		}
//	}
func ValidateVersion(version string) error {
	return validation.ValidateVersion(version)
}
