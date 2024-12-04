### **1. 什么是 `sync.WaitGroup`？**

`sync.WaitGroup` 是 Go 中一个简单的并发同步工具，用于等待一组 Goroutine 执行完成。它通过计数器机制，协调主 Goroutine 和多个子 Goroutine 的协作。

---

### **2. `sync.WaitGroup` 的基本用法**

#### 示例：等待多个 Goroutine 执行完成

```go
package main

import (
 "fmt"
 "sync"
)

func worker(id int, wg *sync.WaitGroup) {
 defer wg.Done() // 确保计数器减一
 fmt.Printf("Worker %d starting\n", id)
 // 模拟工作
 fmt.Printf("Worker %d done\n", id)
}

func main() {
 var wg sync.WaitGroup

 for i := 1; i <= 3; i++ {
  wg.Add(1) // 增加计数器
  go worker(i, &wg)
 }

 wg.Wait() // 阻塞直到计数器变为零
 fmt.Println("All workers finished")
}
```

---

### **3. 工作原理**

1. **计数器**：
   - `Add(delta int)`：增加或减少计数器。
   - 每个 Goroutine 开始前调用 `Add(1)`，完成后调用 `Done()`（相当于 `Add(-1)`）。

2. **等待**：
   - 主 Goroutine 调用 `Wait()`，阻塞直到计数器变为 0。

3. **线程安全**：
   - 内部实现使用原子操作和条件变量，确保线程安全。

---

### **4. 考点**

#### **4.1. 确保 `Add` 在 `Wait` 之前调用**

如果在调用 `Wait` 后再调用 `Add`，可能导致死锁。

**示例问题**：

```go
package main

import (
 "fmt"
 "sync"
)

func main() {
 var wg sync.WaitGroup

 go func() {
  wg.Add(1) // 在 Wait() 之后调用
  defer wg.Done()
  fmt.Println("Worker")
 }()

 wg.Wait() // 死锁
 fmt.Println("Finished")
}
```

**原因**：

- `Wait` 一旦调用，Goroutine 阻塞等待计数器变为 0，而此时 Goroutine 还未开始。

**解决方案**：

- 确保 `Add` 总是在 `Wait` 之前调用。

---

#### **4.2. Goroutine 中忘记 `Done`**

如果忘记调用 `Done`，计数器无法归零，会导致 `Wait` 永远阻塞。

**示例问题**：

```go
package main

import (
 "fmt"
 "sync"
)

func main() {
 var wg sync.WaitGroup

 wg.Add(1) // 增加计数器
 go func() {
  fmt.Println("Worker starting")
  // 忘记调用 Done
 }()

 wg.Wait() // 永远阻塞
 fmt.Println("Finished")
}
```

**解决方案**：

- 使用 `defer wg.Done()` 确保计数器减一，即使函数中发生错误或提前退出。

---

#### **4.3. `WaitGroup` 的拷贝问题**

`WaitGroup` 是不可复制的，拷贝会导致逻辑混乱。

**示例问题**：

```go
package main

import (
 "fmt"
 "sync"
)

func main() {
 var wg sync.WaitGroup
 wg.Add(1)

 wgCopy := wg // 拷贝 WaitGroup
 go func() {
  defer wg.Done()
  fmt.Println("Worker done")
 }()

 wgCopy.Wait() // 无法等待 Goroutine 完成
}
```

**原因**：

- `WaitGroup` 内部使用计数器和条件变量，拷贝后这些状态无法共享。

**解决方案**：

- 始终通过指针传递 `WaitGroup`。

---

### **5. 常见陷阱**

#### **5.1. 过度使用 `WaitGroup`**

当 Goroutine 数量过多时，`WaitGroup` 可能消耗大量内存和资源。此时可以使用以下方法优化：

- **使用固定 Goroutine 池**：限制并发 Goroutine 的数量。
- **使用 `sync.Pool` 或 `chan` 进行任务分发**。

---

#### **5.2. 动态 `Add` 的陷阱**

**场景**：

- 动态生成任务时，`Add` 的数量和 `Done` 的调用可能不匹配。

**解决方案**：

- 确保任务生成与计数器更新在同一逻辑中，避免漏掉。

```go
package main

import (
 "fmt"
 "sync"
)

func main() {
 var wg sync.WaitGroup
 tasks := []int{1, 2, 3, 4, 5}

 for _, task := range tasks {
  wg.Add(1) // 每生成一个任务增加计数器
  go func(t int) {
   defer wg.Done()
   fmt.Printf("Processing task %d\n", t)
  }(task)
 }

 wg.Wait()
 fmt.Println("All tasks completed")
}
```

---

#### **5.3. 并发量过大导致资源竞争**

当 `WaitGroup` 管理大量 Goroutine 时，调度可能出现资源竞争。

**优化建议**：

- 使用工作池模式限制并发 Goroutine 数量：

  ```go
  package main

  import (
      "fmt"
      "sync"
  )

  func worker(id int, wg *sync.WaitGroup, tasks <-chan int) {
      defer wg.Done()
      for task := range tasks {
          fmt.Printf("Worker %d processing task %d\n", id, task)
      }
  }

  func main() {
      var wg sync.WaitGroup
      taskChannel := make(chan int, 10)

      // 创建固定数量的 Goroutine
      for i := 1; i <= 3; i++ {
          wg.Add(1)
          go worker(i, &wg, taskChannel)
      }

      // 发送任务
      for i := 1; i <= 10; i++ {
          taskChannel <- i
      }
      close(taskChannel) // 关闭任务通道

      wg.Wait()
      fmt.Println("All tasks completed")
  }
  ```

---

### **6. 总结**

- **考点**：
  - `WaitGroup` 的生命周期管理。
  - 确保 `Add` 和 `Done` 调用的配合。
  - 避免对 `WaitGroup` 的错误拷贝。

- **陷阱**：
  - 忘记调用 `Done`。
  - 动态 `Add` 任务时漏掉 Goroutine。
  - 大量 Goroutine 导致资源竞争。

通过合理地使用 `sync.WaitGroup`，可以有效管理 Goroutine 的并发性，但需要小心应对上述常见问题和陷阱，以确保程序的稳定性和性能。