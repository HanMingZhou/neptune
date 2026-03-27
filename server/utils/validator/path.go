package validator

import (
	"github.com/pkg/errors"
	"strings"
)

var forbiddenPaths = map[string]bool{
	"/":      true,
	"/bin":   true,
	"/sbin":  true,
	"/usr":   true,
	"/etc":   true,
	"/root":  true,
	"/var":   true,
	"/lib":   true,
	"/lib64": true,
	"/proc":  true,
	"/sys":   true,
	"/dev":   true,
	"/boot":  true,
	"/data":  true,
}

// ValidateMountPath 验证挂载路径的安全性
func ValidateMountPath(path string) error {
	if path == "" {
		return errors.New("挂载路径不能为空")
	}

	// 清理路径中的空格
	path = strings.TrimSpace(path)

	// 必须是绝对路径
	if !strings.HasPrefix(path, "/") {
		return errors.New("挂载路径必须以 / 开头")
	}

	// 禁止路径逃逸
	if strings.Contains(path, "..") {
		return errors.New("挂载路径不能包含 '..'")
	}

	// 禁止特殊字符
	invalidChars := []string{";", "&", "|", "`", "$", "(", ")", "{", "}", "[", "]", "<", ">", "\\", "\"", "'"}
	for _, char := range invalidChars {
		if strings.Contains(path, char) {
			return errors.New("挂载路径包含非法字符")
		}
	}

	// 禁止以 ~ 开头（用户目录）
	if strings.HasPrefix(path, "~") {
		return errors.New("挂载路径不能使用 ~ 符号")
	}

	// 路径长度限制
	if len(path) > 200 {
		return errors.New("挂载路径过长（最多200个字符）")
	}

	// 验证路径只包含安全字符（字母、数字、-、_、/、.）
	for _, char := range path {
		if !((char >= 'a' && char <= 'z') ||
			(char >= 'A' && char <= 'Z') ||
			(char >= '0' && char <= '9') ||
			char == '-' || char == '_' || char == '/' || char == '.') {
			return errors.New("挂载路径只能包含字母、数字、-、_、/ 和 .")
		}
	}

	// 禁止挂载到系统关键目录
	for forbidden := range forbiddenPaths {
		if path == forbidden || strings.HasPrefix(path, forbidden+"/") {
			return errors.New("禁止挂载到系统关键目录: " + forbidden)
		}
	}

	return nil
}

// ValidateSubPath 验证子路径的安全性（相对路径，用于 TensorBoard 日志路径等）
func ValidateSubPath(path string) error {
	// 允许空路径（使用默认值）
	if path == "" {
		return nil
	}

	// 清理路径中的空格
	path = strings.TrimSpace(path)

	// 移除开头的 / 以便统一处理
	path = strings.TrimPrefix(path, "/")

	// 禁止路径逃逸
	if strings.Contains(path, "..") {
		return errors.New("路径不能包含 '..'")
	}

	// 禁止特殊字符
	invalidChars := []string{";", "&", "|", "`", "$", "(", ")", "{", "}", "[", "]", "<", ">", "\\", "\"", "'"}
	for _, char := range invalidChars {
		if strings.Contains(path, char) {
			return errors.New("路径包含非法字符")
		}
	}

	// 禁止以 ~ 开头（用户目录）
	if strings.HasPrefix(path, "~") {
		return errors.New("路径不能使用 ~ 符号")
	}

	// 路径长度限制
	if len(path) > 200 {
		return errors.New("路径过长（最多200个字符）")
	}

	// 验证路径只包含安全字符（字母、数字、-、_、/、.）
	for _, char := range path {
		if !((char >= 'a' && char <= 'z') ||
			(char >= 'A' && char <= 'Z') ||
			(char >= '0' && char <= '9') ||
			char == '-' || char == '_' || char == '/' || char == '.') {
			return errors.New("路径只能包含字母、数字、-、_、/ 和 .")
		}
	}

	return nil
}
