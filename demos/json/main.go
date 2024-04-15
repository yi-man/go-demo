package main

import (
	// "encoding/json"
	"fmt"
)

// func main() {
// 	request := `{"mobile":"18601762076","template":"sms-tmpl-CvsFUg07712","contentVar":"{\"jobType\": \"jupyter\"}"}`

// 	// 解析 request 字符串为一个 map
// 	var requestMap map[string]interface{}
// 	err := json.Unmarshal([]byte(request), &requestMap)
// 	if err != nil {
// 		fmt.Println("解析 JSON 失败:", err)
// 		return
// 	}

// 	// 获取 contentVar 字段的值
// 	contentVarStr, ok := requestMap["contentVar"].(string)
// 	if !ok {
// 		fmt.Println("contentVar 字段不存在或不是字符串类型")
// 		return
// 	}

// 	// 解析 contentVar 字符串为一个 map
// 	var contentVarMap map[string]interface{}
// 	err = json.Unmarshal([]byte(contentVarStr), &contentVarMap)
// 	if err != nil {
// 		fmt.Println("解析 contentVar 字符串失败:", err)
// 		return
// 	}

// 	// 获取 jobType 的值
// 	jobType, ok := contentVarMap["jobType"].(string)
// 	if !ok {
// 		fmt.Println("jobType 字段不存在或不是字符串类型")
// 		return
// 	}

// 	fmt.Println("jobType 的值为:", jobType)
// }

// func main() {
// 	template := `["jobType"]`

// 	var contentVar []string
// 	json.Unmarshal([]byte(template), &contentVar)

// 	fmt.Println("contentVar ", contentVar[0])

// }

func main() {
	arr := [5]int{1, 2, 3, 4, 5}
	for index, value := range arr {
		fmt.Println(index, value) // 输出将是：0 1 1 2 3 4
	}
}
