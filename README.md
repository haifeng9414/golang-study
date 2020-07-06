# golang-study
Go语言学习笔记

## 基础
- 模块被导入时会执行该模块下所有代码文件的init函数，main包下也可以有init函数，init函数先于main函数执行。
- 如果用`var xxx type`的形式声明一个变量，则变量会被初始化为零值；如果有一个确切的非零值用于初始化变量，或者想要使用函数返回值初始化变量，
则应该使用简化变量声明运算符，即 `xxx := something`。
- channel、map、切片都是引用类型，初始化的零值是nil
- Go中所有变量都是值传递，指针变量的值是其所指向的内存地址，传递指针变量时传递的就是这个内存地址
- 下面的代码会输出10个11：
```go
func main() {
	waitGroup := sync.WaitGroup{}
	waitGroup.Add(10)

	for i := 1; i <= 10; i++ {
		go func() {
			fmt.Println(i)
			waitGroup.Done()
		}()
	}

	waitGroup.Wait()
}
```
上面用闭包的方式访问i的值，随着i值的改变，内层的匿名函数也会感知到这些改变，所有的goroutine都会因为闭包共享同样的变量，导致可能所有的
goroutine都输出i的最后一个值，也就是11（也可能是其他值，取决于goroutine什么时候运行的）。有两种方式避免这个问题，思路是一样的，避免
以闭包的方式访问i：
```go
func main() {
	waitGroup := sync.WaitGroup{}
	waitGroup.Add(10)

	for i := 1; i <= 10; i++ {
		j := i
		go func() {
			fmt.Println(j)
			waitGroup.Done()
		}()
	}

	waitGroup.Wait()
}

func main() {
	waitGroup := sync.WaitGroup{}
	waitGroup.Add(10)

	for i := 1; i <= 10; i++ {
		go func(i int) {
			fmt.Println(i)
			waitGroup.Done()
		}(i)
	}

	waitGroup.Wait()
}
```
- golang中命名接口的时候，如果接口只包含一个方法，那么这个的名字以方法名 + er结尾，这是golang的命名惯例。
- 空结构（如`var matcher defaultMatcher`）在创建实例时，不会分配任何内存，这种结构很适合创建没有任何状态的类型。
- 如果一个类型的方法会修改该类型实例的状态，则应该使用指针作为接收者声明方法，而不是类型的值，即`func (m *defaultMatcher) Search`，而
不是`func (m defaultMatcher) Search`。
- 使用指针作为接收者声明的方法，只能在接口类型的值是一个指针的时候被调用。使用值作为接收者声明的方法，在接口类型的值为值或者指针时，都可以被
调用。
- 类型相同（数组长度和元素类型相同）的数组是可以赋值的，赋值时会复制数组的所有元素，如：
```go
func main() {
	var array1 [3]int
	array2 := [3]int{1, 2, 3} // 等价于array2 := [...]int{1, 2, 3}，也可以用array2 := []int{1, 2, 3}声明

	// 如果用array2 := []int{1, 2, 3}声明，则这里会报error：Cannot use 'array2' (type []int) as type [3]int in assignment，这样声明出来的array2实际上是切片
	// 必须用array2 := [3]int{1, 2, 3}或者array2 := [...]int{1, 2, 3}
	array1 = array2 
	array2[0] = -1
	fmt.Println(array1) // 输出：[1 2 3]
	fmt.Println(array2) // 输出：[-1 2 3]
}

// 二维数组也会进行复制
func main() {
	var array1 [4][2]int
	array2 := [4][2]int{{10, 11}, {20, 21}, {30, 31}, {40, 41}}

	array1 = array2
	array2[0][0] = -1
	fmt.Println(array1) // 输出：[[10 11] [20 21] [30 31] [40 41]]
	fmt.Println(array2) // 输出：[[-1 11] [20 21] [30 31] [40 41]]
}
```
如果数组元素类型是指针，则赋值的时候复制的就是指针的值，而不是指针所指向的值。
- 上面的性质也意味着在函数间传递数组是一个开销很大的操作。在函数之间传递变量时，总是以值的方式传递的，如果这个变量是一个数组，意味着整个
数组，不管有多长，都会完整复制，并传递给函数。可以将函数参数由数组改为数组指针（如`func foo(array*[1e6]int)`）从而避免不必要的复制，但是
也要注意函数内对数组的修改会反映在
传入的数组指针对应的数组上。
- 使用切片能更好的解决上面的问题，切片有3个属性，分别是：指向底层数组的指针、切片访问的元素的个数（即长度）和切片允许增长到的元素个数（即容
量）。和数组不一样的地方就是，切片在赋值时会共享底层数组：
```go
func main() {
	array1 := [1]int{1}
	array2 := make([]int, 1) // 声明一个切片，也可以写成[]int{1}，也可以写成make([]int, 1, 1)，分别声明切片的长度和容量
	array3 := array1
	array4 := array2

	array1[0] = -1
	array2[0] = -2

	fmt.Println(array1) // 输出：[-1]
	fmt.Println(array2) // 输出：[-2]
	fmt.Println(array3) // 输出：[1]
	fmt.Println(array4) // 输出：[-2]
}
```
- 可以声明一个nil切片或者空切片，对nil切片和空切片指向append、len、cap的结果是一样：
```go
func main() {
	var slice1 []int // nil切片，即底层的数组指针为nil
	slice2 := make([]int, 0) // 空切片，底层的数组指针指向长度为0的数组
	slice3 := []int{} // 空切片，同上

	fmt.Println(slice1) // 输出：[]
	fmt.Println(slice2) // 输出：[]
	fmt.Println(slice3) // 输出：[]

	fmt.Println(len(slice1), cap(slice1)) // 输出：0 0
	fmt.Println(len(slice2), cap(slice2)) // 输出：0 0
	fmt.Println(len(slice3), cap(slice3)) // 输出：0 0

	slice1 = append(slice1, 1) // 可以安全的append
	fmt.Println(slice1) // 输出：[1]
	slice2 = append(slice2, 1) // 可以安全的append
	fmt.Println(slice2) // 输出：[1]
	slice3 = append(slice3, 1) // 可以安全的append
	fmt.Println(slice3) // 输出：[1]
}
```
- 切片和数组最大的不同在于，切片可以通过创建一个新的切片将底层的数组元素共享出去，而数组的赋值操作是复制整个数组的元素，不过也可以通过一个
已存在的数组创建切片：
```go
func main() {
	slice1 := []int{1, 2, 3, 4, 5} // 创建一个长度为5，容量为5的切片，指针指向拥有5个元素的数组的0号元素
	slice2 := slice1[1:3] // 创建一个长度为2，容量为4的切片，指针指向拥有5个元素的数组的1号元素，所以容量为4

	fmt.Println(slice1) // 输出：[1 2 3 4 5]
	fmt.Println(slice2) // 输出：[2 3]

	slice1[1] = -1
	fmt.Println(slice1) // 输出：[1 -1 3 4 5]
	fmt.Println(slice2) // 输出：[-1 3]

	array1 := [...]int{1, 2, 3, 4, 5}
	array2 := array1[1:3] // 在数组上创建切片

	fmt.Println(array1) // 输出：[1 2 3 4 5]
	fmt.Println(array2) // 输出：[2 3]

	array1[1] = -1
	fmt.Println(array1) // 输出：[1 -1 3 4 5]
	fmt.Println(array2) // 输出：[-1 3]，array2跟着一块变了
}
```
- append方法可以向切片追加元素，并返回一个新的切片。新切片的长度总是大于原来的切片，但是容量可能大于也可能等于原来的切片，取决于原来的切片
是否有可用容量。如果容量不够，则执行append时golang会创建一个两倍原切片容量的新数组并复制数组元素（不一定都是成倍增加，当容量很大时，如
2000，则增加的倍数可能只有1.25，具体的倍数取决的golang的版本是如何实现的），如：
```go
func main() {
	slice1 := []int{1, 2, 3, 4, 5}
	slice2 := slice1[1:4]
	slice3 := append(slice2, -1)

	fmt.Println("after slice3 := append(slice2, -1)")
	fmt.Println(slice1) // 输出：[1 2 3 4 -1]
	fmt.Println(slice2) // 输出：[2 3 4]
	fmt.Println(slice3) // 输出：[2 3 4 -1] // slice2的容量还足够，所以append(slice2, -1)后共享数组的元素被改变了

	slice3[0] = -2
	fmt.Println("after slice3[0] = -2")
	fmt.Println(slice1) // 输出：[1 -2 3 4 -1]
	fmt.Println(slice2) // 输出：[-2 3 4]
	fmt.Println(slice3) // 输出：[-2 3 4 -1] // 3个切片的数组指针指向的还是同一个数组

	slice3 = append(slice3, -3)
	fmt.Println("after slice3 = append(slice3, -3)")
	fmt.Println(slice1) // 输出：[1 -2 3 4 -1]
	fmt.Println(slice2) // 输出：[-2 3 4]
	fmt.Println(slice3) // 输出：[-2 3 4 -1 -3] // slice3的容量不够了，执行append(slice3, -3)时golang会创建一个两边容量的新数组并复制原数组的元素

	slice3[0] = -4
	fmt.Println("after slice3[0] = -4")
	fmt.Println(slice1) // 输出：[1 -2 3 4 -1]
	fmt.Println(slice2) // 输出：[-2 3 4]
	fmt.Println(slice3) // 输出：[-4 3 4 -1 -3] // 此时修改slice3的元素不会影响slice1和slice2的数组元素
	fmt.Println(cap(slice1)) // 输出：5
	fmt.Println(cap(slice2)) // 输出：4
	fmt.Println(len(slice3)) // 输出：5
	fmt.Println(cap(slice3)) // 输出：8，原容量的两倍
}
```
- 创建切片时还可以设置切片的容量，防止append时修改了原底层数组的元素，不过要注意容量不能超过限制：
```go
func main() {
	source := []int{0, 1, 2, 3, 4}

	slice1 := source[2:3:3] // 创建一个新切片，指针指向数组2号元素，希望长度为1（3 - 2 = 1，即从2号元素开始，希望包括原数组的1个元素），容量为1（3 - 2 = 1，即希望能容纳1个元素）
	fmt.Println(len(slice1), cap(slice1)) // 输出：1 1
	slice1 = append(slice1, -1)
	fmt.Println(source) // 输出：[0 1 2 3 4]
	fmt.Println(slice1) // 输出：[2 -1]，由于容量限制，append时创建了一个新数组

	slice2 := source[2:3:4] // 创建一个新切片，指针指向数组2号元素，希望长度为1（3 - 2 = 1，即从2号元素开始，希望包括原数组的1个元素），容量为2（4 - 2 = 2，即希望能容纳2个元素）
	fmt.Println(len(slice2), cap(slice2)) // 输出：1 2
	slice2 = append(slice2, -1)
	fmt.Println(source) // 输出：[0 1 2 -1 4]
	fmt.Println(slice2) // 输出：[2 -1]
    
    //slice3 := source[2:3:6] error: panic: runtime error: slice bounds out of range [::6] with capacity 5
}
```
- append时可以用过`...`将后一个切片的数组元素复制到前一个切片，容量不够会自动创建一个新数组：
```go
func main() {
	s1 := []int{1, 2}
	s2 := []int{3, 4}

	s3 := append(s1, s2...)
	fmt.Println(s1) // 输出：[1 2]
	fmt.Println(s2) // 输出：[3 4]
	fmt.Println(s3) // 输出：[1 2 3 4]

	s3[0] = -1
	fmt.Println(s1) // 输出：[1 2] // s3的修改不会影响s1
	fmt.Println(s2) // 输出：[3 4]
	fmt.Println(s3) // 输出：[-1 2 3 4]，append时容量不够了，所以append返回的切片使用的是新数组

	s4 := make([]int, 1, 10)
	s5 := append(s4, s2...)
	fmt.Println(s5) // 输出：[0 3 4]
	s5[0] = -1
	s5[2] = -2
	s6 := s4[:cap(s4)] // 用于输出底层数组的元素
	fmt.Println(s2) // 输出：[3 4] // s5的修改不会影响s2，因为s2的数组元素在append时是被复制到s4的
	fmt.Println(s4) // 输出：[-1] // s5的修改会影响s4、s6的底层数组
	fmt.Println(s5) // 输出：[-1 3 -2]
	fmt.Println(s6) // 输出：[-1 3 -2 0 0 0 0 0 0 0]，append时容量还够了，所以append返回的切片使用的是原数组
}
```
- 使用range关键字遍历切片和数组时，会创建每个遍历到的元素的副本：
```go
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
    
    /*
    输出：
    Value:10　ValueAddr:C00001C088　ElemAddr:C00001A140
    Value:20　ValueAddr:C00001C088　ElemAddr:C00001A148
    Value:30　ValueAddr:C00001C088　ElemAddr:C00001A150
    Value:40　ValueAddr:C00001C088　ElemAddr:C00001A158
    
    Value:10　ValueAddr:C00001C0B8　ElemAddr:C00001A180
    Value:20　ValueAddr:C00001C0B8　ElemAddr:C00001A188
    Value:30　ValueAddr:C00001C0B8　ElemAddr:C00001A190
    Value:40　ValueAddr:C00001C0B8　ElemAddr:C00001A198
    
    可以发现每次迭代ValueAddr和ElemAddr的值都不一样，ValueAddr地址不变是因为迭代过程使用的是同一个value变量，只不过值在变而已
    */
}
```