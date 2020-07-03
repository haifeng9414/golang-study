package main

import (
	"fmt"
	"sync"
)

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
