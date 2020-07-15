package main

import (
	"fmt"
	"reflect"
)

type Demo struct {
	A string
	b string
}

func main() {
	x := 2
	a := reflect.ValueOf(2) // 直接通过ValueOf方法获取的所有Value都是不可寻址的
	b := reflect.ValueOf(x) // 直接通过ValueOf方法获取的所有Value都是不可寻址的
	c := reflect.ValueOf(&x) // 直接通过ValueOf方法获取的所有Value都是不可寻址的
	d := c.Elem() // d是通过c解引用获取到的，是可寻址的
	d.SetInt(1) // 通过d修改x的值
	fmt.Println(x) // 输出：1

	fmt.Printf("value: %v, type: %v, can addr: %v\n", a, a.Type(), a.CanAddr()) // 输出：value: 2, type: int, can addr: false
	fmt.Printf("value: %v, type: %v, can addr: %v\n", b, b.Type(), b.CanAddr()) // 输出：value: 2, type: int, can addr: false
	fmt.Printf("value: %v, type: %v, can addr: %v\n", c, c.Type(), c.CanAddr()) // 输出：value: 0xc0000b4008, type: *int, can addr: false
	fmt.Printf("value: %v, type: %v, can addr: %v\n", d, d.Type(), d.CanAddr()) // 输出：value: 1, type: int, can addr: true

	// 如果已知变量的类型，也可以强转得到指针再更新
	// 先通过Addr()方法获取一个Value，里面保存了指向变量的指针
	addr := d.Addr()
	// 在Value上调用Interface()方法，获得一个interface{}，里面包含指向变量的指针
	inter := addr.Interface()
	// 强转
	px := inter.(*int)
	// 上面的3补等价于px := &x
	*px = 3 // x = 3
	fmt.Println(x) // 输出：3

	// 也可以直接通过set方法更新，但是如果更新的值的类型和变量的类型不匹配，会panic
	d.Set(reflect.ValueOf(4))
	fmt.Println(x) // 输出：4
	//d.Set(reflect.ValueOf(int64(4))) panic: reflect.Set: value of type int64 is not assignable to type int

	// 如果在不可寻址的变量上调用set方法也会panic
	//b.Set(reflect.ValueOf(3)) panic: reflect: reflect.flag.mustBeAssignable using unaddressable value

	// reflect.Value会记录一个结构体成员是否是未导出成员，如果是的话则拒绝修改操作，所以通过CanAddr方法判断变量是否可寻址来判断是否
	// 能够执行修改操作这是不对的，不过可以直接通过CanSet方法来判断是否可以修改reflect.Value对应的值
	demo := Demo{"testA", "testB"}
	// reflect.ValueOf(demo)不能调用Elem方法，因为demo是不能解引用的
	//_ = reflect.ValueOf(demo).Elem() panic: reflect: call of reflect.Value.Elem on struct Value
	elem := reflect.ValueOf(&demo).Elem()
	aFiled := elem.FieldByName("A")
	bFiled := elem.FieldByName("b")
	fmt.Println(aFiled) // 输出：testA
	fmt.Println(bFiled) // 输出：testB
	fmt.Println(aFiled.CanSet()) // 输出：true
	fmt.Println(bFiled.CanSet()) // 输出：false

	aFiled.SetString("test")
	fmt.Println(aFiled) // 输出：test
	//bFiled.SetString("test") panic: reflect: reflect.flag.mustBeAssignable using value obtained using unexported field
}
