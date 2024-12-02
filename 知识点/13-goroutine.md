# **Go Goroutine 知识全面整理**

## **1. Goroutine 的基础知识**

### **1.1 什么是 Goroutine?**

- Goroutine 是 Go 语言中的轻量级线程，由 Go 运行时调度。
- 每个 Goroutine 有独立的栈空间，初始栈仅 2KB，随着需要动态增长（最大 1GB）。
- Goroutine 是通过 `go` 关键字创建的。

### **1.2 Goroutine 的优势**

- **轻量级**: Goroutine 比系统线程更轻量，能同时支持数十万甚至更多 Goroutine。
- **独立栈**: 栈大小动态调整，内存使用更加高效。
- **调度模型**: Go 使用 M:N 调度模型（N 个 Goroutine 运行在 M 个 OS 线程上），避免频繁的线程切换和创建成本。

---

## **2. Goroutine 的底层实现**

### **2.1 调度模型（GMP 模型）**

- **G (Goroutine):**
  - Goroutine 是 Go 的任务执行单元。
- **M (Machine):**
  - M 代表操作系统线程，承载 G 的执行。
- **P (Processor):**
  - P 代表逻辑处理器，负责执行 G 的调度。
  - P 的数量由 `GOMAXPROCS` 设置决定。

### **2.2 Goroutine 调度过程**

1. **创建 Goroutine (G):**
   - 创建 Goroutine 时，Go 运行时会分配一个 G 对象。
   - G 被放入一个队列，等待分配给 P 执行。
2. **运行时调度:**
   - P 会从其本地队列中获取 G 并将其分配给 M 执行。
   - 如果本地队列为空，P 会尝试从全局队列或其他 P 的队列窃取任务。
3. **抢占机制:**
   - Goroutine 是协作式调度，Go 1.14 引入了异步抢占机制，防止长时间运行的 Goroutine 阻塞。

### **2.3 Goroutine 栈**

- 初始栈只有 2KB，动态增长。
- **栈复制:** 当 Goroutine 栈空间不足时，Go 会复制栈到更大的内存区域。
- **陷阱:** Goroutine 的栈是连续的，过大的栈增长可能导致内存不足。

---

## **3. Goroutine 相关考点与陷阱**

### **3.1 Goroutine 泄漏**

- **问题描述:**
  - Goroutine 没有正确退出，可能会导致 Goroutine 泄漏，进而占用系统资源。例如，当一个 Goroutine 在等待通道数据时，如果通道没有关闭，Goroutine 会一直阻塞，造成泄漏。
- **示例:**

  ```go
  func main() {
      ch := make(chan int)

      // Goroutine 会一直阻塞，因为没有关闭通道
      go func() {
          <-ch
      }()
  }
  ```

- **解决方法:**
  - 使用 `sync.WaitGroup` 或 `context.Context` 管理 Goroutine 生命周期。
  - 确保通道关闭或显式终止 Goroutine。

- **示例:**

```go

// 示例 1 ：
// 使用 sync.WaitGroup 管理 Goroutine 生命周期
// 适合单纯的管理
func main() {
 var wg sync.WaitGroup // 创建 WaitGroup 实例
 ch := make(chan int)

 // 启动 Goroutine
 wg.Add(1) // 增加计数器
 go func() {
  defer wg.Done() // Goroutine 执行完后调用 Done()
  <-ch
  fmt.Println("Goroutine finished.")
 }()

 // 模拟一些处理，关闭通道
 close(ch)

 // 等待 Goroutine 执行完毕
 wg.Wait()
 fmt.Println("Main function exits.")
}



// 示例 2 ：
// 使用 context.Context 管理 Goroutine 生命周期
// 适合协作与取消的场景
func main() {
  ctx,cancel := context.WithCancel(context.Background())

  go func() {
    select {
      case <-ctx.Done():
        fmt.Println("Goroutine cancelled.")
        return
    }
  }()

  // 模拟一些处理
 time.Sleep(2 * time.Second)

 // 取消 Goroutine
 cancel()

 // 等待 Goroutine 完成
 time.Sleep(1 * time.Second)

 fmt.Println("Main function exits.")
}
```

---

### **3.2 Goroutine 调度问题**

- **问题描述:**
  - Goroutine 的抢占式调度不等于实时调度，可能会出现任务延迟。
- **陷阱示例:**

  ```go
  func main() {
      go func() {
          for {
              // 无抢占点，可能导致调度延迟
          }
      }()

      time.Sleep(1 * time.Second)
      fmt.Println("Main goroutine finished")
  }
  ```

- **解决方法:**
  - 在循环中添加**显式让出调度**的操作，例如 `runtime.Gosched()`。

---

### **3.3 通道死锁**

- **问题描述:**
  - Goroutine 通过通道通信时，若通道操作不完整，可能导致死锁。例如，发送数据到通道时，如果没有接收方，发送操作会阻塞。
- **陷阱示例:**

  ```go
  func main() {
      ch := make(chan int)

      go func() {
          ch <- 1
      }()
      // 没有读取通道数据，导致 Goroutine 被阻塞
      // 主程序没有读取通道数据，造成死锁

      //// 确保主程序从通道接收数据，避免死锁
      // data := <-ch
  }
  ```

- **解决方法:**
  - 确保每个通道的操作都有对应的接收方或发送方。
  - 使用 `select` 避免通道操作阻塞。

---

### **3.4 过多 Goroutine 带来的调度压力**

- **问题描述:**
  - 启动过多 Goroutine 会增加调度器的负担，影响性能。
- **陷阱示例:**

  ```go
  func main() {
      for i := 0; i < 1_000_000; i++ {
          go func() {
              time.Sleep(1 * time.Second)
          }()
      }
      time.Sleep(2 * time.Second)
  }
  ```

- **解决方法:**
  - 使用 Goroutine 池限制并发数。

---

### **3.5 数据竞争 (Race Condition)**

- **问题描述:**
  - 多个 Goroutine 操作共享数据时，未使用同步机制，导致结果不一致。
- **陷阱示例:**

  ```go
  var count int

  func main() {
      for i := 0; i < 10; i++ {
          go func() {
              count++
          }()
      }
      time.Sleep(1 * time.Second)
      fmt.Println(count) // 输出可能不为 10
  }
  ```

- **解决方法:**
  - 使用 `sync.Mutex` 或 `sync/atomic` 保护共享数据。

---

### **3.6 Goroutine 的 panic 传播**

- **问题描述:**
  - Goroutine 的 `panic` 不会影响其他 Goroutine，但可能导致程序状态不一致。
- **陷阱示例:**

  ```go
  func main() {
      go func() {
          panic("Goroutine panic")
      }()
      time.Sleep(1 * time.Second)
      fmt.Println("Main goroutine finished")
  }
  ```

- **解决方法:**
  - 使用 `recover` 捕获 `panic`。

---

## **4. Goroutine 的面试常见问题**

1. **Goroutine 与线程的区别是什么？**
   - Goroutine 是用户态线程，由 Go 运行时调度。
   - 系统线程由操作系统管理，创建和切换的成本更高。

2. **如何避免 Goroutine 泄漏？**
   - 使用 `context.Context` 管理 Goroutine 生命周期。
   - 确保通道操作匹配，避免阻塞。

3. **为什么 Goroutine 的栈大小是动态的？**
   - 初始栈较小以减少内存占用；动态调整栈以适应不同任务的内存需求。

4. **Goroutine 如何调度？**
   - Go 使用 GMP 模型，P 决定了 G 的调度策略。

5. **如何限制 Goroutine 并发数？**
   - 使用 `sync.WaitGroup` 或第三方库（如 `errgroup`）限制 Goroutine 的数量。

---

## **5. Goroutine 使用建议**

- **控制 Goroutine 数量**: 避免无限制创建 Goroutine。
- **显式管理生命周期**: 使用 `sync.WaitGroup` 或 `context.Context`。
- **正确处理共享数据**: 使用同步工具（如 `sync.Mutex`）。
- **关注资源释放**: 确保 Goroutine 中分配的资源正确释放。
- **避免阻塞操作**: 使用非阻塞通道操作和超时机制。