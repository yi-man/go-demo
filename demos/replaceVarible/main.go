package main

import (
	"fmt"
	"regexp"
	"strings"
)

// ReplaceVariables 替换字符串中的变量
func ReplaceVariables(input string, variables map[string]string) string {
	// 匹配形如 ${variable} 的模式
	re := regexp.MustCompile(`\${(\w+)}`)

	// 用于替换的函数
	replacer := func(match string) string {
		// 提取变量名
		variableName := strings.Trim(match, "${}")
		// 从映射中获取变量的值
		value, ok := variables[variableName]
		if !ok {
			// 如果映射中没有这个变量，则保持原样
			return match
		}
		return value
	}

	// 替换字符串中的变量
	output := re.ReplaceAllStringFunc(input, replacer)

	return output
}

func main() {
	// 示例字符串
	str := "The job type is ${jobType}, and the job name is ${jobName}."

	// 示例变量映射
	variables := map[string]string{
		"jobType": "engineer",
		"jobName": "software developer",
	}

	// 替换变量并打印结果
	result := ReplaceVariables(str, variables)
	fmt.Println(result)
}
