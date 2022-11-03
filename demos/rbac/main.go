package main

import (
	"fmt"
	"log"

	"github.com/casbin/casbin/v2"
)

// type Sub struct {
// 	Age uint32
// }

// func main() {
// 	e, err := casbin.NewEnforcer("./rules/abac_rule_model.conf", "./rules/abac_rule_policy.csv")

// 	if err != nil {
// 		log.Fatalf("error: enforcer: %s", err)
// 	}

// 	sub := &Sub{Age: 20} // 想要访问资源的用户。
// 	obj := "/data1"      // 将被访问的资源。
// 	act := "read"        // 用户对资源执行的操作。

// 	ok, err := e.Enforce(sub, obj, act)

// 	if err != nil {
// 		// 处理err
// 		fmt.Println(err)
// 	}

// 	if ok {
// 		// 允许alice读取data1
// 		fmt.Println("can read")
// 	} else {
// 		// 拒绝请求，抛出异常
// 		fmt.Println("can not read")
// 	}
// }

// func main() {
// 	e, err := casbin.NewEnforcer("./rules/rbac_model.conf", "./rules/rbac_policy.csv")

// 	if err != nil {
// 		log.Fatalf("error: enforcer: %s", err)
// 	}

// 	sub := "data2_admin" // 想要访问资源的用户。
// 	obj := "data2"       // 将被访问的资源。
// 	act := "read"        // 用户对资源执行的操作。

// 	ok, err := e.Enforce(sub, obj, act)

// 	if err != nil {
// 		// 处理err
// 		fmt.Println(err)
// 	}

// 	if ok {
// 		// 允许alice读取data1
// 		fmt.Println("can read")
// 	} else {
// 		// 拒绝请求，抛出异常
// 		fmt.Println("can not read")
// 	}
// }

func main() {
	e, err := casbin.NewEnforcer("./rules/rbac_with_resource_roles_model.conf", "./rules/rbac_with_resource_roles_model.csv")

	if err != nil {
		log.Fatalf("error: enforcer: %s", err)
	}

	sub := "data2_admin" // 想要访问资源的用户。
	obj := "data2"       // 将被访问的资源。
	act := "read"        // 用户对资源执行的操作。

	ok, err := e.Enforce(sub, obj, act)

	if err != nil {
		// 处理err
		fmt.Println(err)
	}

	if ok {
		// 允许alice读取data1
		fmt.Println("can read")
	} else {
		// 拒绝请求，抛出异常
		fmt.Println("can not read")
	}
}
