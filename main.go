package main

import "fmt"

type Demo struct {
	name string
}

func main() {
	var x interface{} = []int{1, 2, 3}
	var y interface{}
	// panic: runtime error: comparing uncomparable type []int，切片不可比较
	//fmt.Println(x == x)

	x = nil
	fmt.Println(x == y) // 输出：true

	x = Demo{"x"}
	y = Demo{"x"}
	fmt.Println(x == y) // 输出：true

	x = Demo{"x"}
	y = Demo{"y"}
	fmt.Println(x == y) // 输出：false
}
