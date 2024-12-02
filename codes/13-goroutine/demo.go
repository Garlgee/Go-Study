package main

import (
	"fmt"
	"time"
)

func main() {
	f1()
	f3()
}

// goroutine泄露
func f1() {
	ch := make(chan int)

	// Goroutine 会一直阻塞
	go func() {
		ch <- 1
	}()
}

// 使用 sync.WaitGroup 或 context.Context 管理 Goroutine 生命周期。
func f2() {

}

// goroutine 调度问题
func f3() {
	go func() {
		for {
			// 无抢占点，可能导致调度延迟
			fmt.Println("Hello world")
		}
	}()

	time.Sleep(1 * time.Second)
	fmt.Println("Main goroutine finished")
}
