package main

import "fmt"

type Person struct {
	Name string
	Age  int
}

type Node struct {
	Value string
	Prev  *Node
	Next  *Node
}

func main() {
	str := []string{"a", "b"}
	fmt.Println(str)

	var str2 []string
	str2 = append(str2, "c", "d")
	fmt.Println(str2)

	var str3 [2]string
	str3[0] = "a"
	str3[1] = "b"
	fmt.Println(str3)

	var slice []string = make([]string, 5)
	slice[0] = "10"
	slice[1] = "20"
	fmt.Println(slice)

	var int1 []int = make([]int, 5)
	int1[0] = 10
	int1[1] = 20
	fmt.Println(int1)

	int2 := int1[:len(int1)-2]
	fmt.Println(int2)

	obj := make(map[string]string)
	obj["key"] = "value"
	fmt.Println(obj)

	originalSlice := []int{0, 1, 2, 3, 4, 5}
	subSlice := originalSlice[1:4]
	fmt.Println("Sub slice:", subSlice)

	person := make(map[string]Person)
	person["wade"] = Person{Name: "xuxin", Age: 18}
	fmt.Println(person)

}
