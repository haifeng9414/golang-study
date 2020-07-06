package main

import "fmt"

func main() {
	slice := []int{10, 20, 30, 40}

	for index, value := range slice {
		fmt.Printf("Value:%d　ValueAddr:%X　ElemAddr:%X\n", value, &value, &slice[index])
	}

	fmt.Println()

	array := [...]int{10, 20, 30, 40}
	for index, value := range array {
		fmt.Printf("Value:%d　ValueAddr:%X　ElemAddr:%X\n", value, &value, &array[index])
	}
}
