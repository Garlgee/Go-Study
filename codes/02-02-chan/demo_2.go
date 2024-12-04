// 工作池实现

package main

import (
	"fmt"
	"sync"
	"time"

	"math/rand"
)

type Task struct {
	ID      int
	Payload int
}

type Result struct {
	TaskID int
	Output int
}

func main() {
	rand.Seed(time.Now().UnixNano())

	numWorkers := 3
	numTasks := 10
	tasks := make(chan Task, numWorkers)
	results := make(chan Result, numTasks)

	// woker的等待组
	var wg sync.WaitGroup
	// 启动worker
	for i := 0; i < numWorkers; i++ {
		wg.Add(1)
		go worker(i, tasks, results, &wg)
	}

	go func() {
		// 生成任务
		for i := 0; i < numTasks; i++ {
			fmt.Printf("Main: Generated Task ID: %d\n", i)
			tasks <- Task{ID: i, Payload: rand.Intn(100)}
		}

		close(tasks) // 关闭任务通道，通知workers没有更多任务
	}()

	go func() {
		wg.Wait()      // 等待所有worker完成
		close(results) // 关闭结果通道，通知主程序已经没有结果
	}()

	// 打印结果
	// 会从通道中接收数据，直到通道关闭为止。
	// 当results通道关闭时，range会自动结束，退出循环。
	// 在通道关闭之前，主goroutine会一直阻塞在for循环中等待通道中的新数据。
	for result := range results {
		fmt.Printf("Main: Result{ Task ID: %d, Output: %d }\n", result.TaskID, result.Output)
	}

	fmt.Println("Main: All workers completed.")

}

func worker(id int, tasks <-chan Task, results chan<- Result, wg *sync.WaitGroup) {
	defer wg.Done()
	for task := range tasks {
		// 模拟任务处理
		time.Sleep(time.Millisecond * time.Duration(rand.Intn(100)))

		fmt.Printf("Worker %d processing Task ID: %d\n", id, task.ID)
		// 将结果发送到结果通道
		results <- Result{TaskID: task.ID, Output: task.Payload * 2}
	}

}
