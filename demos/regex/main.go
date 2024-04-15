package main

import (
	"fmt"
	"regexp"
)

func main() {
	re := regexp.MustCompile(`^[\p{Han}A-Za-z][\p{Han}A-Za-z0-9_]*[^\d_]$`)
	fmt.Println("测试字符串: ", re.MatchString("测试字符串"))
	fmt.Println("_测试字符串:", re.MatchString("_测试字符串"))
	fmt.Println("测试字符串_:", re.MatchString("测试字符串_"))
	fmt.Println("11测试字符串: ", re.MatchString("11测试字符串"))
	fmt.Println("aaa_测试字符串:", re.MatchString("aaa_测试字符串"))
	fmt.Println("bbb:", re.MatchString("bbb"))
}
