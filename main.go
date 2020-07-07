package main

import (
	"fmt"
	"golang-study/demo"
)

func main() {
	ad := demo.Admin{}
	ad.Name = "dhf"
	ad.Email = "dhf@demo.com"
	fmt.Println(ad.Name, ad.Email)
}
