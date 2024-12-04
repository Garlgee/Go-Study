package main

import (
	"fmt"
	"time"
)

func main() {
	f0()
	f1()
}

// nil通道阻塞
func f0() {
	fmt.Println("------------------------------ f0")
	var ch1 chan int // 未初始化，为nil通道
	ch2 := make(chan int, 1)

	go func() {
		ch2 <- 1
	}()

	select {
	case <-make(chan int): // 使用一个永远不被关闭的通道来模拟阻塞
		fmt.Println("This will block forever")
	case <-ch1: // 通道ch1未初始化，为nil通道，会阻塞
		fmt.Println("received from ch1")
	case <-ch2:
		fmt.Println("received from ch2")

	}

	time.Sleep(1 * time.Second)
	close(ch2)
}

// 通道泄露
func f1() {
	fmt.Println("------------------------------ f1")
	ch1 := make(chan int)
	ch2 := make(chan int)

	go func() {
		fmt.Println("ch1 send start")
		ch1 <- 1
		fmt.Println("ch1 send end")
	}()

	// select中忽略了 case msg2:=<-ch2:
	// ch2 中的数据无法被接收，这意味着 ch2 上的发送操作会被阻塞，进而导致发送方（go func()）永远等待，形成死锁。
	// 未处理的通道，通道泄露
	go func() {
		fmt.Println("ch2 send start")
		// 如果这里有资源申请，忽略ch2的select时，资源会泄露
		ch2 <- 2
		fmt.Println("ch2 send end")
	}()

	select {
	case msg1 := <-ch1:
		fmt.Println("received from ch1", msg1)
		// 这里只读取了 ch1，漏掉了 ch2
		// case msg2:=<-ch2:
		//	 fmt.Println("received from ch2",msg2)
	}

	// 这将导致 ch2 的发送方（goroutine）在此阻塞，造成资源泄漏
	time.Sleep(3 * time.Second)
	close(ch1)
	fmt.Println("ch1 closed")
	close(ch2)
	fmt.Println("ch2 closed")
	time.Sleep(3 * time.Second)
	// panic: send on closed channel
}
