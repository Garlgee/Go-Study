package main

import (
	"fmt"
	"sync"
	"time"
)

/*
当需要启动大量 Goroutine 时，
如果直接用 sync.WaitGroup 管理，每个 Goroutine 都需要消耗内存和资源（如栈空间）。
这种情况可能导致资源耗尽，尤其在任务数量远大于系统支持的并发 Goroutine 时。
*/
func main() {
	f1()
}

func worker(id int, wg *sync.WaitGroup, tasks <-chan int) {
	defer wg.Done()

	for task := range tasks {
		fmt.Println("worker", id, "started  task", task)
		time.Sleep(time.Millisecond * 10)
		fmt.Println("worker", id, "finished task", task)
	}
}

// 方法1：使用固定 Goroutine 池，限制同时运行的 Goroutine 数量，通过通道分发任务。
func f1() {
	var wg sync.WaitGroup
	taskChan := make(chan int, 10)

	// 使用工作池模式限制并发go协程数量
	// 创建固定的并发数
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go worker(i, &wg, taskChan)

	}

	// 发送任务
	for i := 0; i < 100; i++ {
		taskChan <- i
	}
	// 发送完关闭通道，这样worker goroutine就会结束
	close(taskChan)

	wg.Wait()
	fmt.Println("All tasks completed")
}

type Task struct {
	id int
	//params interface{}
}

func (t *Task) Reset() {
	t.id = -1 // 标记可复用
}

func worker1(id int, taskPool *sync.Pool, wg *sync.WaitGroup) {
	defer wg.Done()

	for {
		task := taskPool.Get().(*Task)
		if task == nil || task.id < 0 {
			break
		}
		fmt.Println("worker[", id, "] started deal task[", task.id, "]")
		time.Sleep(time.Millisecond * 10)
		fmt.Println("worker[", id, "] finished task[", task.id, "]")
		task.Reset()
		taskPool.Put(task)
	}
}
