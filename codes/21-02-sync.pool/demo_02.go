package main

import (
	"fmt"
	"sync"
)

// sync.Pool 对象池 顺序性：不能保证顺序

func main() {
	// 创建一个对象池
	pool := &sync.Pool{
		New: func() any {
			return "New Object"
		},
	}

	// 向池中放入三个对象
	pool.Put("Object 1")
	pool.Put("Object 2")
	pool.Put("Object 3")

	// 并发获取对象，观察顺序
	var wg sync.WaitGroup
	for i := 0; i < 5; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			obj := pool.Get()
			fmt.Printf("Goroutine %d got: %v\n", id, obj)
		}(i)
	}

	wg.Wait()
}

/*
输出1：
Goroutine 4 got: Object 1
Goroutine 2 got: New Object
Goroutine 0 got: Object 2
Goroutine 3 got: Object 3
Goroutine 1 got: New Object

输出2：
Goroutine 0 got: Object 2
Goroutine 1 got: Object 3
Goroutine 4 got: Object 1
Goroutine 3 got: New Object
Goroutine 2 got: New Object
*/
