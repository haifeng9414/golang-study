# golang-study
Go语言学习笔记

## 基础
- 模块被导入时会执行该模块下所有代码文件的init函数，main包下也可以有init函数，init函数先于main函数执行。
- 如果用`var xxx type`的形式声明一个变量，则变量会被初始化为零值；如果有一个确切的非零值用于初始化变量，或者想要使用函数返回值初始化变量，则应该使用简化变量声明运算符，即 `xxx := something`。
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