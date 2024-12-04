package main

import (
	"fmt"
	"sync"
	"time"
)

func main() {
	f1()
	f2()
	//f3()
	f4()
}

func worker(id int, wg *sync.WaitGroup) {
	defer wg.Done()
	fmt.Printf("Worker %d starting\n", id)
	time.Sleep(time.Millisecond * 10) // 模拟工作
	fmt.Printf("Worker %d done\n", id)
}
func f1() {
	fmt.Println("------------ basic ------------")
	var wg sync.WaitGroup
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go worker(i, &wg)
	}

	wg.Wait()
	fmt.Println("all goroutines are done")
}

// 确保 Add 在 Wait 之前调用
func f2() {
	fmt.Println("------------ dead lock ------------")
	var wg sync.WaitGroup

	go func() {
		wg.Add(1) //should call wg.Add(1) before starting the goroutine to avoid a race (SA2000)
		defer wg.Done()
		fmt.Println("Worker")
	}()

	wg.Wait() // 死锁
	fmt.Println("Finished")
}

// fatal error: all goroutines are asleep - deadlock!
func f3() {
	fmt.Println("------------ miss Done() ------------")
	var wg sync.WaitGroup

	wg.Add(1) // 增加计数器
	go func() {
		fmt.Println("Worker starting")
		// 忘记调用 Done
	}()

	wg.Wait() // 永远阻塞
	fmt.Println("Finished")
}

// waitgroup 不可以拷贝,始终通过指针传递 WaitGroup。
func f4() {
	fmt.Println("------------ copy ------------")
	var wg sync.WaitGroup
	wg1 := wg // 错误

	wg.Add(1)
	go func() {
		defer wg.Done()
		fmt.Println("Worker done")
	}()

	//wg.Wait()
	wg1.Wait() // 无法等待 Goroutine 完成
	fmt.Println("Finished")

}