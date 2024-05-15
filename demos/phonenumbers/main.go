package main

import (
	"fmt"
	"regexp"

	"github.com/nyaruka/phonenumbers"
)

func isValidPhoneNumber2(mobile string) bool {
	// 定义手机号正则表达式
	phoneRegex := `^\+\d{1,3}\d{5,14}$`

	phoneNumber := mobile
	if len(phoneNumber) > 0 && phoneNumber[0] != '+' {
		phoneNumber = "+86" + phoneNumber
	}

	// 编译正则表达式
	regexp := regexp.MustCompile(phoneRegex)

	// 使用正则表达式匹配手机号
	return regexp.MatchString(phoneNumber)
}

func isValidPhoneNumber(mobile string) bool {

	phonenumbers.PhoneNumber
	// 定义手机号正则表达式
	phoneRegex := `^\+\d{1,3}\d{5,14}$`

	phoneNumber := mobile
	if len(phoneNumber) > 0 && phoneNumber[0] != '+' {
		phoneNumber = "+86" + phoneNumber
	}

	// 编译正则表达式
	regexp := regexp.MustCompile(phoneRegex)

	// 使用正则表达式匹配手机号
	return regexp.MatchString(phoneNumber)
}

func main() {
	phoneNumber := "+8867751875413"
	if isValidPhoneNumber(phoneNumber) {
		fmt.Println("Valid phone number")
	} else {
		fmt.Println("Invalid phone number")
	}

	phoneNumber = "+8618601762076"
	if isValidPhoneNumber(phoneNumber) {
		fmt.Println("Valid phone number")
	} else {
		fmt.Println("Invalid phone number")
	}

	phoneNumber = "18601762076"
	if isValidPhoneNumber(phoneNumber) {
		fmt.Println("Valid phone number")
	} else {
		fmt.Println("Invalid phone number")
	}
}
