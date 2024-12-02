package main

import (
	"context"
	"fmt"
	"time"
)

func main() {
	f0()
	f1() // 带超时的context
	f2() // 带deadline的context
	f3()
}

func f0() {
	ctx, cancel := context.WithCancel(context.Background())

	go func(ctx context.Context) {
		for {
			select {
			case <-ctx.Done():
				fmt.Println("Go routine stopped	")
				return
			default:
				fmt.Println("Go routine running")
				time.Sleep(time.Millisecond * 500)
			}
		}
	}(ctx)

	time.Sleep(2 * time.Second)
	cancel() //通知go routine停止
	time.Sleep(time.Second)
}

// 带超时的context
func f1() {
	fmt.Println("f1 start.")

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*2)
	defer cancel()

	go func(ctx context.Context) {
		for {
			select {
			case <-ctx.Done():
				fmt.Println("f1 Operation timed out.")
				return
			default:
				fmt.Println("f1 go routine running")
				time.Sleep(time.Millisecond * 500)
			}
		}
	}(ctx)

	time.Sleep(3 * time.Second) // 模拟耗时操作
	fmt.Println("f1 finished.")
}

// 带deadline的context
func f2() {
	fmt.Println("f2 start.")
	deadline := time.Now().Add(2 * time.Second)
	ctx, cancel := context.WithDeadline(context.Background(), deadline)
	defer cancel()
	go func(ctx context.Context) {
		for {
			select {
			case <-ctx.Done():
				fmt.Println("f2 Operation timed out.")
				return
			default:
				fmt.Println("f2 go routine running")
				time.Sleep(time.Millisecond * 500)
			}
		}
	}(ctx)

	time.Sleep(3 * time.Second) // 模拟耗时操作
	fmt.Println("f2 finished.")
}

// 携带值的context
func f3() {
	fmt.Println("f3 start.")

	// 返回一个携带值的上下文，用于在 API 边界之间传递元数据。
	// 不推荐将 context 用作数据存储，应尽量避免频繁使用。
	ctx := context.WithValue(context.Background(), "key", "value")
	go func() {
		fmt.Println("Go routine running with value:", ctx.Value("key"))
	}()

	time.Sleep(3 * time.Second)
	fmt.Println("f3 finished.")
}
