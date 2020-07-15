# golang-study
Go语言学习笔记

## 基础
- 模块被导入时会执行该模块下所有代码文件的init函数，main包下也可以有init函数，init函数先于main函数执行。
- 如果用`var xxx type`的形式声明一个变量，则变量会被初始化为零值；如果有一个确切的非零值用于初始化变量，或者想要使用函数返回值初始化变量，
则应该使用简化变量声明运算符，即 `xxx := something`。
- channel、map、切片、指针、函数变量和接口变量都是引用类型，初始化的零值是nil：
```go
func main() {
	var a *int
	var b []int
	var c map[string]int
	var d chan int
	var e func(string) int
	var f error // error是接口，接口类型只有函数，所以不应该作为值类型

	fmt.Println(a == nil) // 输出：true
	fmt.Println(b == nil) // 输出：true
	fmt.Println(c == nil) // 输出：true
	fmt.Println(d == nil) // 输出：true
	fmt.Println(e == nil) // 输出：true
	fmt.Println(f == nil) // 输出：true
}
```
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
- 如果数组元素类型是指针，则赋值的时候复制的就是指针的值，而不是指针所指向的值。
- 上面的性质也意味着在函数间传递数组是一个开销很大的操作。在函数之间传递变量时，总是以值的方式传递的，如果这个变量是一个数组，意味着整个
数组，不管有多长，都会完整复制，并传递给函数。可以将函数参数由数组改为数组指针（如`func foo(array*[1e6]int)`）从而避免不必要的复制，但是
也要注意函数内对数组的修改会反映在
传入的数组指针对应的数组上。
- 使用切片能更好的解决上面的问题，在函数间传递切片只会复制切片的属性，不会复制切片的底层数组（无论底层数组多大，切片大小都是24个字节）。切片有3个属性，分别是：指向底层数组的指针、切片访问的元素的个数（即长度）和切片允许增长到的元素个数（即容
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
- 尽管如上所述，nil切片和空切片很相似，但是也还是有两个需要注意的不同点，一个是json序列化，一个是`reflect.DeepEqual`的结果：
```go
type Res struct {
	Data []string
}

func main() {
	var nilSlice []string
	emptySlice := make([]string, 0)

	// 使用json序列化
	res, _ := json.Marshal(Res{Data: nilSlice})
	res2, _ := json.Marshal(Res{Data: emptySlice})

	fmt.Println(string(res))  // 输出：{"Data":null}
	fmt.Println(string(res2)) // 输出：{"Data":[]}

	fmt.Println(reflect.DeepEqual(nilSlice, emptySlice)) // 输出：false
	fmt.Printf("Got: %+v, Want: %+v\n", nilSlice, emptySlice) // 输出：Got: [], Want: []，DeepEqual为false，但是打印时又看不出差别，出问题时可能影响问题的定位
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
- append方法可以向切片追加元素，并返回一个新的切片。新切片的长度总是大于原来的切片，但是容量可能大于也可能等于原来的切片，取决于原来的切片是否有可用容量。如果容量不够，则执行append时golang会创建一个两倍原切片容量的新数组并复制数组元素（不一定都是成倍增加，当容量很大时，如
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
- 切片、map、函数具有引用语义，不能用于==比较，而map的key要求能够进行==比较，所以这些类型不能用作map的key。另外需要注意的是，如果数组的
元素类型为切片、map或函数，或者某个结构化类型包含切片、map或函数的属性，则也不能用于==比较，所以也不能作为map的key
- 接口类型是可以比较的。可以将接口类型作为map的key或者执行==比较，接口类型执行==操作的返回值取决于两边接口类型是否都是nil或者他们的动态类
型相同并且动态值进行==操作也相同。接口类型的比较不一定是安全的，其它类型要么是安全的可比较类型（如基本类型和指针）要么是完全不可比较的类
型（如切片，映射类型，和函数），在比较接口值或者包含了接口值的聚合类型时，如果其动态类型是不可比较的，则执行==操作是会引发panic：
```go
type Demo struct {
	name string
}

func main() {
	var x interface{} = []int{1, 2, 3}
	var y interface{}
	// panic: runtime error: comparing uncomparable type []int，x的动态类型为切片，不可比较
	//fmt.Println(x == x)

	x = nil
	fmt.Println(x == y) // 输出：true，x和y都是nil

	x = Demo{"x"}
	y = Demo{"x"}
	fmt.Println(x == y) // 输出：true，x和y的动态类型相同，动态类型的值==结果相同

	x = Demo{"x"}
	y = Demo{"y"}
	fmt.Println(x == y) // 输出：false
}
```
- 一个包含nil指针的接口不是nil接口，下面的代码在debug为true时可以正常执行，但是如果debug为false时会panic，函数f已经对参数out作为nil判
断，但是并没有起作用，原因是debug为false时，buf为nil，`*bytes.Buffer`实现了`io.Writer`接口，所以可以被传入f函数，此时参数out被赋值为
一个`*bytes.Buffer`的空指针，即out接口的动态类型为`*bytes.Buffer`，动态值为nil，这个时候out是一个非空接口：
```go
const debug = false

func main() {
	var buf *bytes.Buffer
	if debug {
		buf = new(bytes.Buffer)
	}
	f(buf)
}

func f(out io.Writer) {
	if out != nil {
		out.Write([]byte("done!\n")) // panic: runtime error: invalid memory address or nil pointer dereference
	}
}
```

上面的问题的解决方案是在main函数中声明buf为io.Writer，这样在debug为false时，buf为nil接口，传入函数f时out参数被赋值为nil接口：
```go
func main() {
	var buf io.Writer
	if debug {
		buf = new(bytes.Buffer)
	}
	f(buf) // OK
}
```
- map和切片一样，在函数间传递时不会创建一个map的副本，而是创建一个map的引用
- 如果使用类型的值作为方法的接收者，则在调用方法时，方法会接收到一个类型值的副本，如果使用类型的指针作为方法的接收者，则方法接收到的是指针
的副本，所以可以通过指针修改类型的值，如果类型的值复制代价很大，则应该避免使用类型的值作为方法的接收者：
```go
type User struct {
	username string
	password string
}

func (u User) test1() { // 类型的值作为接收者
	u.username = "test1"
}

func (u *User) test2() { // 类型的指针作为接收者
	u.username = "test2"
}

func main() {
	u := User{"dhf", "pwd"}
	fmt.Println(u) // 输出：{dhf pwd}
	
	u.test1()
	fmt.Println(u) // 输出：{dhf pwd}，没有被test1方法修改
	
	u.test2()
	fmt.Println(u) // 输出：{dhf pwd}，被test2方法修改了
}
```
- 如果类型的基础类型是引用类型，则使用类型的值作为方法接收者或者传递类型的值到函数，复制的也是引用：
```go
type Demo []string // 基础类型是切片

func (d Demo) test1() { // 使用类型的值作为方法接收者
	d[0] = "test1"
}

func test2(d Demo) { // 使用类型的值作为函数参数
	d[1] = "test2"
}

func main() {
	d := Demo{"1", "2"}

	fmt.Println(d)
	d.test1()
	test2(d)
	fmt.Println(d) // 值被改变了
}
```
- 一般情况下使用类型的值还是类型的指针作为方法接收者或者函数参数，取决于是否想要方法或函数能够对类型的值直接做修改。如果不想要修改，则应该使用
类型的值作为方法接收者或者函数参数，这样方法或函数操作的是类型的值的副本。如果想要被修改，则应该使用类型的指针。

	如果单纯为了效率考虑，可以使用类型的指针避免类型的值被不必要的复制，如果类型的值可能包含一个大数组。

	有些类型不能被安全的复制，同时也不允许修改，则可以使用指针并且不公开属性，如标准库中的File：
```go
type File struct {
	*file // 内嵌类型
}

type file struct {
	pfd         poll.FD
	name        string
	dirinfo     *dirInfo // nil unless directory being read
	nonblock    bool     // whether we set nonblocking mode
	stdoutOrErr bool     // whether this is stdout or stderr
	appendMode  bool     // whether file is opened for appending
}
```

File类型的实际类型是file，使用指针可以使得File作为函数参数时复制的只是指针的值，同时不公开file属性，使得客户端无法修改file类型的属性。

实际上使用值接收者还是指针接收者，不应该只由方法或函数是否修改了接收到的值来决定，应该基于类型的本质。如果类型的值可以被安全的复制，如时间
`time.Time`、数字`int64`等，则应该使用值接收者；如果类型的值不能被安全的复制，如上面的`os.File`，则即使方法或者函数没有修改类型的值，
也应该使用指针。
- 类型的指针可以使用以类型的值作为接收者的方法和以类型的指针作为接收者的方法，类型的值只能使用以类型的值作为接收者的方法：
```go
type notifier interface {
	notify() string
}

type Demo struct {
}

func (d *Demo) notify() string {
	return "demo"
}

func test(n notifier)  {
	fmt.Println(n.notify())
}

func main() {
	d := Demo{}

	//test(d) // cannot use d (type Demo) as type notifier in argument to test
	test(&d) // 只有类型的指针才能作为notifier接口的实现类

	d.notify() // 类型的值不能直接调用以类型的指针作为接收者的方法，这里是因为golang帮忙转成了(&d).notify()
}
```
- 方法的接受者只能是类型的值或者类型的指针，为了避免歧义，如果一个类型本身是一个指针的话，是不允许作为方法的接收者的：
```go
type P *int // 类型本身是个指针

// invalid receiver type P (P is a pointer type)
//func (p P) test() {
//
//}

type Q int

func (p Q) test1() {}
func (p *Q) test2() {}
```
- 类型的指针如果是nil，也可以进行方法调用，这一点需要注意：
```go
type Demo struct {
}

func (d *Demo) test() {
	fmt.Println(d == nil)
}

func main() {
	var d *Demo
	d.test() // 输出：true
}
```
- 可以通过嵌入类型实现类型的复用，已有的类型可以被嵌入到新的类型，已有类型称为内部类型，新的类型称为外部类型。内部类型的标识符会提升到外部
类型上。这些被提升的标识符就像直接声明在外部类型里的标识符一样，也是外部类型的一部分。外部类型也可以通过声明与内部类型标识符同名的标识符来
覆盖内部标识符的字段或者方法。

另外需要注意的是，如果内部类型的某个方法使用指针作为接收者，则该方法只能被外部类型的指针访问：
```go
type notifier interface {
	notify()
}

type user struct {
	name  string
	email string
}

func (u *user) notify() { // user用指针作为接收者实现了notifier接口
	fmt.Printf("Sending user email to %s<%s>\n",
		u.name,
		u.email)
}

func (u *user) print() {
	fmt.Printf("name: %s, email: %s\n", u.name, u.email)
}

type admin struct {
	user  // 嵌入的内部类型
	level string
}

func (u *admin) print() { // 可以覆盖内部类型的方法，这里用u *admin或者u admin都可以
	u.user.print()
	fmt.Printf("level: %s\n", u.level)
}

func test(n notifier)  {
	n.notify()
}

func main() {
	ad := admin{
		user: user{
			name:  "dhf",
			email: "dhf@yahoo.com",
		},
		level: "super",
	}

	ad.user.notify() // 可以访问内部类型的方法

	ad.notify() // 直接调用内部类型的方法也可以

	fmt.Println(ad.user.name) // 还可以访问属性

	ad.print() // 可以覆盖内部类型的方法

	test(&ad) // 由于user的指针类型实现了notifier接口，所以必须使用admin的指针才能作为notifier接口的实现类
}
```
- 未公开的类型可以通过被公开的方法暴露出去，并被短变量声明所引用，通过短变量可以访问未公开类型的公开方法，不过这样不是一个好的编码习惯，ide
也会给出警告：
```go
// demo包下
package demo

type user struct { // user类型未公开
	Name  string // 未公开类型的公开属性
	email string // 未公开属性
}

func User(name, email string) user { // 公开的方法返回未公开的类型，ide会给出警告：Exported function with unexported return type
	return user{name, email}
}

// main包下
func main() {
	user := demo.User("dhf", "dhf@demo.com") // 通过短变量声明引用返回的未公开方法
	fmt.Println(user.Name) // 只能够访问被公开的方法
}
```
- 和上面类似，未公开的内部类型的公开属性可以被公开的外部类型所公开：
```go
// demo包下
package demo

type user struct { // user类型未公开
	Name  string // 未公开的公开属性
	Email string // 未公开的公开属性
}

type Admin struct {
	user  // 未公开的内部类型
}

// main包下
func main() {
	ad := demo.Admin{}
	ad.Name = "dhf" // 可以正常访问
	ad.Email = "dhf@demo.com" // 可以正常访问
	fmt.Println(ad.Name, ad.Email)
}
```
- 通道在不使用后一定要记得close
- `reflect.TypeOf`方法返回`reflect.Type`，表示传入的参数的动态类型的接口值，并且总是表示具体的类型：
```go
func main() {
	t := reflect.TypeOf(3) // TypeOf方法返回reflect.Type，表示传入的参数的动态类型的接口值
	fmt.Println(t.String()) // "int"
	fmt.Println(t) // "int"

    var w io.Writer = os.Stdout
	fmt.Println(reflect.TypeOf(w)) // 输出："*os.File"，TypeOf方法总是返回具体的类型
}
```
- 当打印值时，可以使用`%v`及其变形形式打印：
```go
func main() {
	t := &T{7, -2.35, "abc\tdef"}
	fmt.Println(t) // 输出：&{7 -2.35 abc   def}
	fmt.Printf("%v\n", t) // 输出：&{7 -2.35 abc   def}，和直接执行Println的结果一样
	fmt.Printf("%+v\n", t) // 输出：&{a:7 b:-2.35 c:abc     def}，会带上字段名
	fmt.Printf("%#v\n", t) // 输出：&main.T{a:7, b:-2.35, c:"abc\tdef"}，go语法表示形式
	fmt.Printf("%#v\n", time.UTC) // 输出：&time.Location{name:"UTC", zone:[]time.zone(nil), tx:[]time.zoneTrans(nil), cacheStart:0, cacheEnd:0, cacheZone:(*time.zone)(nil)}
}
```
-  `reflect.ValueOf`方法返回`reflect.Value`，`reflect.Value`类似`interface{}`，可以持有任意类型的值。`reflect.ValueOf`返回的结
果也是针对具体的类型，除非`reflect.Value`持有的是字符串值，否则其`String()`方法返回具体的类型的字符串表示形式。`%v`标志参数，可以输出
`reflect.Value`持有的值：
```go
func main() {
	var x interface{}

	x = 3
	v := reflect.ValueOf(x) // a reflect.Value
	fmt.Println(v) // 输出：3
	fmt.Printf("%v\n", v) // 输出：3
	fmt.Println(v.String()) // 输出：<int Value>

	x = []int{1}
	v = reflect.ValueOf(x)
	fmt.Println(v) // 输出：[1]
	fmt.Printf("%v\n", v) // 输出：[1]
	fmt.Println(v.String()) // 输出：<[]int Value>

	x = time.Now()
	v = reflect.ValueOf(x)
	fmt.Println(v) // 输出：2020-07-14 21:00:10.911226 +0800 CST m=+0.000141229
	fmt.Printf("%v\n", v) // 输出：2020-07-14 21:00:10.911226 +0800 CST m=+0.000141229
	fmt.Println(v.String()) // 输出：<time.Time Value>
}
```
- `reflect.Value`还有一些其他有用的方法：
    - `Kind()`：获取值的种类，`reflect.Value`的值的类型有无限多种，但是种类是有限的，可以利用`Kind()`方法实现自己的`fmt.Println()`方法：
        - Bool、String和所有数字类型的基础类型; 
        - Array和Struct对应的聚合类型; 
        - Chan、Func、Ptr、Slice和Map对应的引用类似;
        - 接口类型;
        - 表示空值的无效类型(空的`reflect.Value`对应`Invalid`);
    - `Interface()`：返回一个`interface{}`，其指向`reflect.Value`表示的值
    - `Type()`：返回一个`reflect.Type`
    - 其他方法就不一一列举了，可以看下面的例子
```go
type Movie struct {
	Title, Subtitle string
	Year            int
	Color           bool
	Actor           map[string]string
	Oscars          []string
	Sequel          *string
}

func main() {
	now := time.Now()
	v := reflect.ValueOf(now) // 获取reflect.Value
	fmt.Printf("%T\n", v)     // 输出：reflect.Value

	i := v.Interface()    // 返回一个interface{}，其指向reflect.Value表示的值
	a := i.(time.Time)    // 强转
	fmt.Printf("%T\n", a) // 输出：time.Time
	fmt.Println(a)        // 输出：2020-07-14 21:05:19.764024 +0800 CST m=+0.000107937

	b := v.Type()                         // 返回一个reflect.Type
	fmt.Printf("%T\n", b)                 // 输出：*reflect.rtype
	fmt.Println(b)                        // 输出：time.Time
	fmt.Println(b == reflect.TypeOf(now)) // 输出：true

	c := v.Kind()                          // 获取值的类型
	fmt.Printf("%T\n", c)                  // 输出：reflect.Kind
	fmt.Println(c)                         // 输出：struct
	fmt.Println(reflect.ValueOf(1).Kind()) // int

	var x int64 = 1
	d := 1 * time.Nanosecond
	fmt.Println(printAny(x))                  // 输出：1
	fmt.Println(printAny(d))                  // 输出：1
	fmt.Println(printAny([]int64{x}))         // 输出：[]int64 0xc0000b4090
	fmt.Println(printAny([]time.Duration{d})) // 输出：[]time.Duration 0xc0000b4098

	strangelove := Movie{
		Title:    "Dr. Strangelove",
		Subtitle: "How I Learned to Stop Worrying and Love the Bomb", Year: 1964,
		Color: false,
		Actor: map[string]string{
			"Dr. Strangelove":            "Peter Sellers",
			"Grp. Capt. Lionel Mandrake": "Peter Sellers",
			"Pres. Merkin Muffley":       "Peter Sellers",
			"Gen. Buck Turgidson":        "George C. Scott",
			"Brig. Gen. Jack D. Ripper":  "Sterling Hayden",
			`Maj. T.J. "King" Kong`:      "Slim Pickens",
		},
		Oscars: []string{
			"Best Actor (Nomin.)",
			"Best Adapted Screenplay (Nomin.)", "Best Director (Nomin.)",
			"Best Picture (Nomin.)",
		},
	}

	/*
	输出：
	Display strangelove (main.Movie):
	strangelove.Title = "Dr. Strangelove"
	strangelove.Subtitle = "How I Learned to Stop Worrying and Love the Bomb"
	strangelove.Year = 1964
	strangelove.Color = false
	strangelove.Actor["Gen. Buck Turgidson"] = "George C. Scott"
	strangelove.Actor["Brig. Gen. Jack D. Ripper"] = "Sterling Hayden"
	strangelove.Actor["Maj. T.J. \"King\" Kong"] = "Slim Pickens"
	strangelove.Actor["Dr. Strangelove"] = "Peter Sellers"
	strangelove.Actor["Grp. Capt. Lionel Mandrake"] = "Peter Sellers"
	strangelove.Actor["Pres. Merkin Muffley"] = "Peter Sellers"
	strangelove.Oscars[0] = "Best Actor (Nomin.)"
	strangelove.Oscars[1] = "Best Adapted Screenplay (Nomin.)"
	strangelove.Oscars[2] = "Best Director (Nomin.)"
	strangelove.Oscars[3] = "Best Picture (Nomin.)"
	strangelove.Sequel = nil
	*/
	Display("strangelove", strangelove)

	/*
	输出：
	Display os.Stderr (*os.File):
	(*(*os.Stderr).file).pfd.fdmu.state = 0
	(*(*os.Stderr).file).pfd.fdmu.rsema = 0
	(*(*os.Stderr).file).pfd.fdmu.wsema = 0
	(*(*os.Stderr).file).pfd.Sysfd = 2
	(*(*os.Stderr).file).pfd.pd.runtimeCtx = 0
	(*(*os.Stderr).file).pfd.iovecs = nil
	(*(*os.Stderr).file).pfd.csema = 0
	(*(*os.Stderr).file).pfd.isBlocking = 1
	(*(*os.Stderr).file).pfd.IsStream = true
	(*(*os.Stderr).file).pfd.ZeroReadIsEOF = true
	(*(*os.Stderr).file).pfd.isFile = true
	(*(*os.Stderr).file).name = "/dev/stderr"
	(*(*os.Stderr).file).dirinfo = nil
	(*(*os.Stderr).file).nonblock = false
	(*(*os.Stderr).file).stdoutOrErr = true
	(*(*os.Stderr).file).appendMode = false
	 */
	Display("os.Stderr", os.Stderr)
}

func printAny(value interface{}) string {
	return printValue(reflect.ValueOf(value))
}

func Display(name string, x interface{}) {
	fmt.Printf("Display %s (%T):\n", name, x)
	display(name, reflect.ValueOf(x))
}

func display(path string, v reflect.Value) {
	switch v.Kind() {
	case reflect.Invalid:
		fmt.Printf("%s = invalid\n", path)
	case reflect.Slice, reflect.Array:
		for i := 0; i < v.Len(); i++ {
			display(fmt.Sprintf("%s[%d]", path, i), v.Index(i)) // 通过Index方法获取切片和数组的指定元素
		}
	case reflect.Struct:
		for i := 0; i < v.NumField(); i++ {
			fieldPath := fmt.Sprintf("%s.%s", path, v.Type().Field(i).Name) // 结构体通过reflect.Type的Field方法获取属性
			display(fieldPath, v.Field(i))
		}
	case reflect.Map:
		for _, key := range v.MapKeys() {
			display(fmt.Sprintf("%s[%s]", path, printValue(key)), v.MapIndex(key)) // map通过MapKeys方法获取所有的key，并通过MapIndex获取指定key的value
		}
	case reflect.Ptr:
		if v.IsNil() {
			fmt.Printf("%s = nil\n", path)
		} else {
			display(fmt.Sprintf("(*%s)", path), v.Elem()) // Elem方法返回一个reflect.Value变量，其持有指针指向的变量
		}
	case reflect.Interface:
		if v.IsNil() {
			fmt.Printf("%s = nil\n", path)
		} else {
			fmt.Printf("%s.type = %s\n", path, v.Elem().Type())
			display(path+".value", v.Elem()) // Elem方法返回一个reflect.Value变量，其持有接口指向的变量
		}
	default: // basic types, channels, funcs
		fmt.Printf("%s = %s\n", path, printValue(v))
	}
}

func printValue(v reflect.Value) string {
	switch v.Kind() {
	case reflect.Invalid:
		return "invalid"
	case reflect.Int, reflect.Int8, reflect.Int16,
		reflect.Int32, reflect.Int64:
		return strconv.FormatInt(v.Int(), 10)
	case reflect.Uint, reflect.Uint8, reflect.Uint16,
		reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		return strconv.FormatUint(v.Uint(), 10)
	// 简单起见，省略了一些浮点类型的判断
	case reflect.Bool:
		return strconv.FormatBool(v.Bool())
	case reflect.String:
		return strconv.Quote(v.String())
	case reflect.Chan, reflect.Func, reflect.Ptr, reflect.Slice, reflect.Map:
		return v.Type().String() + " 0x" + strconv.FormatUint(uint64(v.Pointer()), 16)
	default: // reflect.Array, reflect.Struct, reflect.Interface...
		return v.Type().String() + " value"
	}
}
```
- go中有些值是不可寻址的，不可寻址的值不能获取其指针，不可寻址的情况有：
    - 常量的值；
    - 基本类型值的字面量；
    - 算术操作的结果值；
    - 对各种字面量的索引表达式和切片表达式的结果值。不过有一个例外，对切片字面量的索引结果值是可寻址的；
    - 对字符串变量的索引表达式和切片表达式的结果值；
    - 对字典变量的索引表达式的结果值；函数字面量和方法字面量，以及对它们的调用表达式的结果值；
    - 结构体字面量的字段值，也就是对结构体字面量的选择表达式的结果值；
    - 类型转换表达式的结果值；
    - 类型断言表达式的结果值；
    - 接收表达式的结果值；
    - 指向函数或者方法的值；
上面的情况其实可以总结为：
    - 不可变的值；
    - 临时结果（对切片字面量的索引结果值是可寻址的。因为不论怎样，每个切片值都会持有一个底层数组，而这个底层数组中的每个元素值都是有一个确切的内存地址的。）；
    - 不安全的，若拿到某值的指针可能会破坏程序的一致性，那么就是不安全的，该值就不可寻址；

	例子：
```go
type Named interface {
	Name() string
}

type Dog struct {
	name string
}

func (dog *Dog) SetName(name string) {
	dog.name = name
}

func (dog Dog) Name() string {
	return dog.name
}

func main() {
	const num = 123
	//_ = &num // 常量属于不可变的值，不可寻址。
	//_ = &(123) // 基本类型的字面量属于不可变的值，不可寻址。

	var str = "abc"
	_ = str // 变量是可寻址的
	//_ = &(str[0]) // 对字符串变量的索引结果值不可寻址，因为字符串不可变，寻址了也没意义
	//_ = &(str[0:2]) // 对字符串变量的切片结果值不可寻址，因为字符串不可变，寻址了也没意义
	str2 := str[0]
	_ = &str2 // 但这样的寻址就是合法的，因为str2 := str[0]实际上是复制一份str[0]并赋值给str2

	//_ = &(123 + 456) // 算术操作的结果值属于临时变量，不可寻址。
	num2 := 456
	_ = &num2 // 变量是可寻址的
	//_ = &(num + num2) // 算术操作的结果值属于临时变量，不可寻址。

	//_ = &([3]int{1, 2, 3}[0]) // 数组字面量的索引结果值属于临时变量，不可寻址。
	//_ = &([3]int{1, 2, 3}[0:2]) // 数组字面量的切片结果值属于临时变量，不可寻址。
	_ = &([]int{1, 2, 3}[0]) // 需要注意的是，对切片字面量的索引结果值却是可寻址的，因为不论怎样，每个切片值都会持有一个底层数组，而这个底层数组中的每个元素值都是有一个确切的内存地址。
	//_ = &([]int{1, 2, 3}[0:2]) // 切片字面量的切片结果值属于临时变量，不可寻址。
	//_ = &(map[int]string{1: "a"}[0]) // 字典字面量的索引结果值属于临时变量，不可寻址。

	var map1 = map[int]string{1: "a", 2: "b", 3: "c"}
	_ = &map1 // 变量是可寻址的
	//_ = &(map1[2]) // 字典变量的索引结果值不可寻址，因为字典中的每个键值对的存储位置都可能会变化，而且这种变化外界是无法感知的。

	//_ = &(func(x, y int) int { // 字面量代表的函数不可寻址。
	//	return x + y
	//})
	//_ = &(fmt.Sprintf) // 标识符代表的函数不可寻址。
	//_ = &(fmt.Sprintln("abc")) // 对函数的调用结果值属于临时变量，不可寻址。

	dog := Dog{"little pig"}
	_ = &dog // 变量是可寻址的
	//_ = &(dog.Name) // 标识符代表的函数不可寻址。
	//_ = &(dog.Name()) // 对方法的调用结果值不可寻址。

	//_ = &(Dog{"little pig"}.name) // 结构体字面量的字段属于临时变量，不可寻址。

	//_ = &(interface{}(dog)) // 类型转换表达式的结果值属于临时变量，不可寻址。
	dogI := interface{}(dog)
	_ = &dogI // 变量是可寻址的
	//_ = &(dogI.(Named)) // 类型断言表达式的结果值属于临时变量，不可寻址。
	named := dogI.(Named)
	_ = &named // 变量是可寻址的
	//_ = &(named.(Dog)) // 类型断言表达式的结果值属于临时变量，不可寻址。

	var chan1 = make(chan int, 1)
	chan1 <- 1
	//_ = &(<-chan1) // 接收表达式的结果值属于临时变量，不可寻址。
}
```
- 下面看看如果通过`reflect.Value`修改值。可寻址的值表示可以通过内存地址来更新的值，一个变量就是一个可寻址的内存空间，里面存储了一个值，
并且存储的值可以通过内存地址来更新：
```go
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

    // 观察c和d输出的不同，就能理解为什么直接通过ValueOf方法获取的Value是不可寻址的，c是个指向指针变量的reflect.Value，不能直接
    // 修改一个指针变量的值，需要通过Elem方法解引用
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
	// reflect.ValueOf(demo)不能调用Elem方法，因为引用变量是不能解引用的
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
```