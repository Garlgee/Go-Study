// Thread Pool in Go

package main

import (
	"context"
	"fmt"
	"sync"
	"time"
)

type Task struct {
	Input    any
	Execute  func(input any) (output any, err error) // 任务执行函数
	Callback func(result any, err error)             // 结果回调
}

type ThreadPool struct {
	taskChan    chan Task          // Tasks to be executed
	workerCount int                // Number of workers
	wg          sync.WaitGroup     // Wait group to wait for all workers to finish
	ctx         context.Context    // Context for all workers
	cancel      context.CancelFunc // Cancel context to stop all workers
}

func NewThreadPool(workerCount int, taskQueueSize int) *ThreadPool {
	ctx, cancel := context.WithCancel(context.Background())
	tp := &ThreadPool{
		taskChan:    make(chan Task, taskQueueSize),
		workerCount: workerCount,
		ctx:         ctx,
		cancel:      cancel,
	}

	return tp
}

func (tp *ThreadPool) Start() {
	for i := 0; i < tp.workerCount; i++ {
		tp.wg.Add(1)
		go tp.worker(i)
	}
}
func (tp *ThreadPool) Submit(t Task) error {
	select {
	case tp.taskChan <- t:
		return nil
	case <-tp.ctx.Done(): // 收到关闭信号
		return fmt.Errorf("Pool is stopped")
	default: // 任务队列已满，返回错误
		return fmt.Errorf("Task queue is full")
	}
}
func (tp *ThreadPool) Stop() {
	// 取消 context
	tp.cancel()
	// 关闭任务队列
	close(tp.taskChan)
	// 等待所有 worker 结束
	tp.wg.Wait()
	fmt.Println("Pool stopped")
}

func (tp *ThreadPool) worker(id int) {
	defer tp.wg.Done()
	fmt.Println("Worker", id, "started")
	for {
		select {
		case <-tp.ctx.Done():
			// 收到关闭信号
			fmt.Println("Worker", id, "stopped")
			return
		case task, ok := <-tp.taskChan:
			if !ok {
				// 任务队列已关闭
				fmt.Printf("Worker %d stopping\n", id)
				return
			}
			output, err := task.Execute(task.Input)
			if task.Callback != nil {
				task.Callback(output, err)
				fmt.Println("Worker", id, "finished task")
			}

		}
	}
}

func main() {
	tp := NewThreadPool(3, 10)
	tp.Start()

	var wg sync.WaitGroup
	for i := 0; i < 20; i++ {
		err := tp.Submit(Task{
			Input: i,
			Execute: func(input any) (output any, err error) {
				id := input.(int)
				fmt.Printf("Processing task %d\n", id)
				time.Sleep(2 * time.Second) // 模拟任务耗时
				return fmt.Sprintf("Result of task %d", id), nil
			},
			Callback: func(result any, err error) {
				defer wg.Done()
				if err != nil {
					fmt.Printf("Task failed: %v\n", err)
				} else {
					fmt.Printf("Task completed: %v\n", result)
				}
			}})
		if err != nil {
			fmt.Println(err)
		} else {
			wg.Add(1)
		}
	}

	wg.Wait()
	tp.Stop()
}
