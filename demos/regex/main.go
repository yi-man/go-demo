package main

import (
	"fmt"
	"regexp"
)

func main() {
	// 定义正则表达式
	re := regexp.MustCompile(`^[A-Z][A-Z0-9_]*[A-Z0-9]$`)

	// 测试字符串
	strs := []string{
		"ABC123_",   // 不匹配，以下划线结尾
		"_ABC123",   // 不匹配，以下划线开头
		"123ABC_",   // 不匹配，以数字开头
		"ABC_123",   // 匹配
		"AB_1C23",   // 匹配
		"ABCD",      // 匹配
		"AB_CD_123", // 匹配
	}

	// 遍历测试字符串
	for _, str := range strs {
		if re.MatchString(str) {
			fmt.Printf("%s 匹配\n", str)
		} else {
			fmt.Printf("%s 不匹配\n", str)
		}
	}
}
