// Package dependency 提供与PHP Composer依赖项相关的功能
//
// 本包处理Composer依赖项的各种操作，包括：
// - 依赖项格式验证
// - 检查依赖项是否存在
// - 添加和删除依赖项
// - 合并依赖映射
package dependency

import (
	"fmt"
	"regexp"
	"strings"
)

// GetPackageNameParts 将包名分割为供应商和项目部分
//
// 参数:
//   - packageName: 要分割的包名，格式应为"vendor/project"
//
// 返回:
//   - string: 供应商名称（如"symfony"）
//   - string: 项目名称（如"console"）
//   - error: 如果包名格式无效则返回错误
//
// 示例:
//
//	vendor, project, err := dependency.GetPackageNameParts("symfony/console")
//	if err != nil {
//		log.Fatal(err)
//	}
//	fmt.Printf("供应商: %s, 项目: %s\n", vendor, project)
//	// 输出: 供应商: symfony, 项目: console
//
//	// 无效格式示例
//	_, _, err = dependency.GetPackageNameParts("invalid-name")
//	if err != nil {
//		fmt.Println(err)
//		// 输出: invalid package name format: invalid-name, expected vendor/project
//	}
func GetPackageNameParts(packageName string) (string, string, error) {
	parts := strings.Split(packageName, "/")
	if len(parts) != 2 {
		return "", "", fmt.Errorf("invalid package name format: %s, expected vendor/project", packageName)
	}
	return parts[0], parts[1], nil
}

// ValidatePackageName 根据Composer规范检查包名是否有效
//
// 参数:
//   - packageName: 要验证的包名
//
// 返回:
//   - error: 如果包名无效则返回描述性错误，有效则返回nil
//
// 规则:
//   - 包名不能为空
//   - 必须符合"vendor/project"格式
//   - 供应商和项目名只能包含小写字母数字字符、下划线、短横线和点号
//   - 必须以字母数字字符开头
//
// 示例:
//
//	// 有效包名
//	err := dependency.ValidatePackageName("symfony/console")
//	if err != nil {
//		log.Fatal(err)
//	}
//	fmt.Println("包名有效")
//
//	// 无效包名
//	err = dependency.ValidatePackageName("Invalid/Name")
//	if err != nil {
//		fmt.Println(err)
//		// 输出: invalid vendor name 'Invalid': must contain only lowercase alphanumeric characters, '_', '-', and '.'
//	}
func ValidatePackageName(packageName string) error {
	if packageName == "" {
		return fmt.Errorf("包名不能为空")
	}

	// 检查格式vendor/project
	parts := strings.Split(packageName, "/")
	if len(parts) != 2 {
		return fmt.Errorf("包名必须符合'vendor/project'格式")
	}

	// 验证供应商和项目名
	validNameRegex := regexp.MustCompile(`^[a-z0-9]([_.-]?[a-z0-9]+)*$`)
	if !validNameRegex.MatchString(parts[0]) {
		return fmt.Errorf("无效的供应商名'%s': 只能包含小写字母数字字符、'_'、'-'和'.'", parts[0])
	}
	if !validNameRegex.MatchString(parts[1]) {
		return fmt.Errorf("无效的项目名'%s': 只能包含小写字母数字字符、'_'、'-'和'.'", parts[1])
	}

	return nil
}

// DependencyExists 检查依赖项是否存在于require部分
//
// 参数:
//   - require: 依赖映射，key为包名，value为版本约束
//   - packageName: 要检查的包名
//
// 返回:
//   - bool: 如果依赖项存在返回true，否则返回false
//
// 示例:
//
//	// 检查PHP依赖是否存在
//	require := map[string]string{
//		"php": ">=7.4",
//		"symfony/console": "^5.4",
//	}
//
//	if dependency.DependencyExists(require, "php") {
//		fmt.Println("PHP依赖存在，版本要求:", require["php"])
//	} else {
//		fmt.Println("PHP依赖不存在")
//	}
//
//	// 检查不存在的依赖
//	if !dependency.DependencyExists(require, "not-exists/package") {
//		fmt.Println("依赖项不存在")
//	}
func DependencyExists(require map[string]string, packageName string) bool {
	if require == nil {
		return false
	}
	_, exists := require[packageName]
	return exists
}

// AddDependency 向require部分添加包
//
// 参数:
//   - require: 要修改的依赖映射
//   - packageName: 要添加的包名，格式为"vendor/project"
//   - version: 依赖版本约束，如"^5.4"、">=7.4"
//
// 返回:
//   - error: 如果包名格式无效则返回错误，成功则返回nil
//
// 注意:
//   - 如果包已存在，将更新其版本约束
//   - require映射会被直接修改，无需重新赋值
//
// 示例:
//
//	// 初始化依赖映射
//	require := make(map[string]string)
//
//	// 添加PHP版本要求
//	err := dependency.AddDependency(require, "php", ">=7.4")
//	if err != nil {
//		log.Fatal(err)
//	}
//
//	// 添加Symfony依赖
//	err = dependency.AddDependency(require, "symfony/console", "^5.4")
//	if err != nil {
//		log.Fatal(err)
//	}
//
//	// 更新现有依赖的版本
//	err = dependency.AddDependency(require, "php", ">=8.0")
//	if err != nil {
//		log.Fatal(err)
//	}
//
//	fmt.Println(require)
//	// 输出: map[php:>=8.0 symfony/console:^5.4]
func AddDependency(require map[string]string, packageName, version string) error {
	if err := ValidatePackageName(packageName); err != nil {
		return err
	}

	require[packageName] = version
	return nil
}

// RemoveDependency 从require部分移除包
//
// 参数:
//   - require: 要修改的依赖映射
//   - packageName: 要移除的包名
//
// 返回:
//   - bool: 如果包被成功移除返回true，如果包不存在或映射为nil则返回false
//
// 注意:
//   - require映射会被直接修改，无需重新赋值
//
// 示例:
//
//	// 已有依赖映射
//	require := map[string]string{
//		"php": ">=7.4",
//		"symfony/console": "^5.4",
//	}
//
//	// 移除Symfony依赖
//	if dependency.RemoveDependency(require, "symfony/console") {
//		fmt.Println("symfony/console已移除")
//	} else {
//		fmt.Println("symfony/console不存在")
//	}
//
//	// 尝试移除不存在的依赖
//	if !dependency.RemoveDependency(require, "not-exists/package") {
//		fmt.Println("not-exists/package不存在，无需移除")
//	}
//
//	fmt.Println(require)
//	// 输出: map[php:>=7.4]
func RemoveDependency(require map[string]string, packageName string) bool {
	if require == nil {
		return false
	}

	if _, exists := require[packageName]; !exists {
		return false
	}

	delete(require, packageName)
	return true
}

// MergeDependencies 返回两个依赖映射的合并结果
//
// 参数:
//   - require: 运行时依赖映射
//   - requireDev: 开发时依赖映射
//
// 返回:
//   - map[string]string: 包含所有依赖的新映射，key为包名，value为版本约束
//
// 注意:
//   - 如果两个映射中存在相同的包，requireDev中的版本会覆盖require中的版本
//   - 返回的是新映射，不会修改输入参数
//
// 示例:
//
//	// 运行时依赖
//	require := map[string]string{
//		"php": ">=7.4",
//		"symfony/console": "^5.4",
//	}
//
//	// 开发时依赖
//	requireDev := map[string]string{
//		"phpunit/phpunit": "^9.5",
//		"php": ">=8.0", // 会覆盖require中的版本
//	}
//
//	// 合并依赖
//	allDeps := dependency.MergeDependencies(require, requireDev)
//
//	fmt.Println(allDeps)
//	// 输出: map[php:>=8.0 symfony/console:^5.4 phpunit/phpunit:^9.5]
func MergeDependencies(require, requireDev map[string]string) map[string]string {
	result := make(map[string]string)

	// 添加运行时依赖
	for pkg, version := range require {
		result[pkg] = version
	}

	// 添加开发时依赖
	for pkg, version := range requireDev {
		result[pkg] = version
	}

	return result
}
