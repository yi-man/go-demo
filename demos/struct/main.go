package main

import (
	"fmt"
	"reflect"
)

type Greeter struct {
	Hello string
}

func main() {
	greeter := &Greeter{
		Hello: "hello",
	}

	fmt.Println(*greeter, &greeter, reflect.TypeOf((greeter)), reflect.TypeOf((*greeter)))
}
