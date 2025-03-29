package composer

import (
	"github.com/scagogogo/php-composer-json-parser/pkg/composer/archive"
	"github.com/scagogogo/php-composer-json-parser/pkg/composer/autoload"
	"github.com/scagogogo/php-composer-json-parser/pkg/composer/config"
	"github.com/scagogogo/php-composer-json-parser/pkg/composer/repository"
)

// ComposerJSON 表示composer.json文件的根结构
//
// 该结构体映射了PHP Composer项目中composer.json文件的所有可能字段，
// 允许通过Go代码读取和修改PHP项目的依赖管理配置。
//
// 示例:
//
//	{
//	  "name": "vendor/project",
//	  "description": "项目描述",
//	  "type": "library",
//	  "require": {
//	    "php": ">=7.4",
//	    "symfony/console": "^5.4"
//	  },
//	  "autoload": {
//	    "psr-4": {
//	      "App\\": "src/"
//	    }
//	  }
//	}
type ComposerJSON struct {
	// Name 包名，格式为"vendor/project"，如"symfony/console"
	Name string `json:"name,omitempty"`

	// Description 项目描述
	Description string `json:"description,omitempty"`

	// Type 包类型，如"library"、"project"、"metapackage"、"composer-plugin"
	Type string `json:"type,omitempty"`

	// Keywords 关键词数组，用于Packagist搜索
	Keywords []string `json:"keywords,omitempty"`

	// Homepage 项目主页URL
	Homepage string `json:"homepage,omitempty"`

	// Version 版本号（通常由VCS标签自动确定，很少在composer.json中手动指定）
	Version string `json:"version,omitempty"`

	// License 许可证，可以是单个字符串或字符串数组
	License interface{} `json:"license,omitempty"` // 可以是字符串或字符串数组

	// Authors 作者信息数组
	Authors []Author `json:"authors,omitempty"`

	// Support 支持信息
	Support Support `json:"support,omitempty"`

	// Require 运行时依赖，key为包名，value为版本约束
	Require map[string]string `json:"require,omitempty"`

	// RequireDev 开发时依赖，key为包名，value为版本约束
	RequireDev map[string]string `json:"require-dev,omitempty"`

	// Conflict 冲突依赖，指定与当前包冲突的包和版本
	Conflict map[string]string `json:"conflict,omitempty"`

	// Replace 替换依赖，指定当前包可以替代的其他包和版本
	Replace map[string]string `json:"replace,omitempty"`

	// Provide 提供依赖，指定当前包可以满足的其他包的需求
	Provide map[string]string `json:"provide,omitempty"`

	// Suggest 建议依赖，推荐但非必需的包
	Suggest map[string]string `json:"suggest,omitempty"`

	// Autoload 自动加载配置，包含PSR-4、PSR-0、classmap等
	Autoload autoload.Autoload `json:"autoload,omitempty"`

	// AutoloadDev 开发时自动加载配置，通常用于测试代码
	AutoloadDev autoload.Autoload `json:"autoload-dev,omitempty"`

	// Repositories 自定义包仓库配置
	Repositories []repository.Repository `json:"repositories,omitempty"`

	// Config Composer配置选项
	Config config.Config `json:"config,omitempty"`

	// Scripts Composer脚本定义，可以是字符串或字符串数组
	Scripts map[string]interface{} `json:"scripts,omitempty"`

	// ScriptsDescriptions 脚本的说明文本
	ScriptsDescriptions map[string]string `json:"scripts-descriptions,omitempty"`

	// Extra 附加元数据，供第三方工具使用
	Extra map[string]interface{} `json:"extra,omitempty"`

	// Bin 可执行文件列表，会被软链接到vendor/bin目录
	Bin []string `json:"bin,omitempty"`

	// Archive 打包时的配置，如排除文件
	Archive archive.Archive `json:"archive,omitempty"`

	// Abandoned 标记包已被废弃，可以是布尔值或推荐替代包的字符串
	Abandoned interface{} `json:"abandoned,omitempty"` // 可以是布尔值或字符串

	// NonFeatureBranches 指定哪些分支不应被视为功能分支
	NonFeatureBranches []string `json:"non-feature-branches,omitempty"`

	// MinimumStability 最低稳定性要求，如"stable"、"RC"、"beta"、"alpha"、"dev"
	MinimumStability string `json:"minimum-stability,omitempty"`

	// PreferStable 优先使用稳定版本
	PreferStable bool `json:"prefer-stable,omitempty"`
}

// Author 表示Composer包的作者信息
//
// 示例:
//
//	{
//	  "name": "张三",
//	  "email": "zhangsan@example.com",
//	  "homepage": "https://example.com",
//	  "role": "Developer"
//	}
type Author struct {
	// Name 作者姓名（必填）
	Name string `json:"name,omitempty"`

	// Email 作者邮箱
	Email string `json:"email,omitempty"`

	// Homepage 作者个人主页
	Homepage string `json:"homepage,omitempty"`

	// Role 在项目中的角色，如"Developer"、"Maintainer"
	Role string `json:"role,omitempty"`
}

// Support 包含包的支持信息
//
// 示例:
//
//	{
//	  "email": "support@example.com",
//	  "issues": "https://github.com/vendor/project/issues",
//	  "docs": "https://docs.example.com"
//	}
type Support struct {
	// Email 支持邮箱
	Email string `json:"email,omitempty"`

	// Issues 问题跟踪URL
	Issues string `json:"issues,omitempty"`

	// Forum 论坛URL
	Forum string `json:"forum,omitempty"`

	// Wiki Wiki页面URL
	Wiki string `json:"wiki,omitempty"`

	// IRC IRC频道
	IRC string `json:"irc,omitempty"`

	// Source 源代码URL
	Source string `json:"source,omitempty"`

	// Docs 文档URL
	Docs string `json:"docs,omitempty"`

	// Chat 聊天渠道URL
	Chat string `json:"chat,omitempty"`
}
