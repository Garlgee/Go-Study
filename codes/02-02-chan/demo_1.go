package main

import (
	"fmt"
	"time"
)

func main() {
	fmt.Println("--------------------1--------------------")
	// ch1 := make(chan string)     //不带缓冲的通道,进和出都会阻塞
	// ch3 := make(<-chan string)   //只读通道
	// ch4 := make(chan<- string)   //只写通道
	ch2 := make(chan string, 10) //带10个缓冲的通道,进一次长度 +1，出一次长度 -1，如果长度等于缓冲长度时，再进就会阻塞。
	ch2 <- "ch2"
	if val, ok := <-ch2; ok {
		fmt.Println(val)
	}

	if false {
		fmt.Println("--------------------2--------------------")
		ch5 := make(chan string, 1)
		close(ch5)
		ch5 <- "ch5" //panic : send on closed channel
	}

	if false {
		fmt.Println("--------------------3--------------------")
		ch6 := make(chan string, 1)
		close(ch6)
		close(ch6) //panic : close of closed channel
	}

	// fmt.Println("--------------------4--------------------")
	// ch7 := make(<-chan string, 1)
	// close(ch7) // 编译: invalid operation: cannot close receive-only channel ch7 (variable of type <-chan string)

	// close之后还可以读数据
	fmt.Println("--------------------5--------------------")
	ch8 := make(chan string, 1)
	ch8 <- "ch8"
	close(ch8)
	if val, ok := <-ch8; ok {
		fmt.Println("after close, read chan:", val)
	}

	f6()
	f7()
	f8()
	f9()
}

// 无缓存的chan,读写数据都会阻塞当前goroutine
func f6() {
	fmt.Println("f6 start")
	ch9 := make(chan string)
	// ch9 <- "ch9" //fatal error: all goroutines are asleep - deadlock!
	// <-ch9 //fatal error: all goroutines are asleep - deadlock!
	close(ch9)

}

// 无缓存chan,在不同goroutine中读写数据
func f7() {
	fmt.Println("f7 start")
	ch10 := make(chan string)
	go func() {
		ch10 <- "a"
	}()
	go func() {
		val := <-ch10
		fmt.Println(val)
	}()
	time.Sleep(1 * time.Second)
	fmt.Println("f7 end")
}

/*
	带缓冲的通道，如果长度等于缓冲长度时，再进就会阻塞。

f8 start
producer start
f8 end
*/
func f8() {
	fmt.Println("f8 start")
	ch := make(chan string, 3)
	go producer(ch)
	time.Sleep(1 * time.Second)
	fmt.Println("f8 end")
}
func producer(ch chan string) {
	fmt.Println("producer start")
	ch <- "a" // 写入通道，不阻塞
	ch <- "b" // 写入通道，不阻塞
	ch <- "c" // 写入通道，不阻塞
	ch <- "d" // 缓冲已满，阻塞，等待消费者取出数据
	fmt.Println("producer end")
}

func customer(ch chan string) {
	for {
		//msg := <-ch
		fmt.Println(<-ch)
	}
}
func f9() {
	fmt.Println("f9 start")
	ch := make(chan string, 3)
	go producer(ch)
	go customer(ch)

	time.Sleep(1 * time.Second)
	fmt.Println("f9 end")
}
