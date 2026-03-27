package utils

import (
	"strings"
)

// GetBrowser 获取浏览器类型
func GetBrowser(ua string) string {
	if strings.Contains(ua, "Chrome") {
		return "Chrome"
	} else if strings.Contains(ua, "Firefox") {
		return "Firefox"
	} else if strings.Contains(ua, "Safari") {
		return "Safari"
	} else if strings.Contains(ua, "Edge") {
		return "Edge"
	} else if strings.Contains(ua, "MSIE") || strings.Contains(ua, "Trident") {
		return "IE"
	}
	return "Other"
}

// GetOS 获取操作系统类型
func GetOS(ua string) string {
	if strings.Contains(ua, "Windows") {
		return "Windows"
	} else if strings.Contains(ua, "Mac") {
		return "Mac OS"
	} else if strings.Contains(ua, "Linux") {
		return "Linux"
	} else if strings.Contains(ua, "Android") {
		return "Android"
	} else if strings.Contains(ua, "iPhone") || strings.Contains(ua, "iPad") {
		return "iOS"
	}
	return "Other"
}
