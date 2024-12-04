# select 通道泄露

## 示例代码

```go
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
  //  fmt.Println("received from ch2",msg2)
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
```

在上述示例中，**通道泄露的具体资源**，主要体现在以下几个方面：

---

## **1. 阻塞的 Goroutine**

- **泄漏资源：**
  - **Goroutine** 是轻量级线程，由 Go 运行时管理。
  - 如果一个 Goroutine 永远等待某个通道操作（发送或接收），就会进入**阻塞状态**，占用系统资源（例如内存堆栈、调度器中的上下文等）。
  - 在我们的例子中：
    - 向 `ch2` 发送数据的 Goroutine 被阻塞，因为没有对应的接收者。
    - 阻塞的 Goroutine 永远不会被调度运行，也无法被显式回收。
  - **泄露表现：**
    - Goroutine 的资源（例如运行时栈、元数据）会一直占用内存。

---

## **2. 未被处理的通道缓冲**

- **泄漏资源：**
  - 在带缓冲的通道（`make(chan int, N)`）中，缓冲区占用内存。
  - 如果通道中的数据未被及时读取，或者通道未被正确关闭，缓冲区中的数据无法释放，从而导致**内存泄漏**。
  - 在无缓冲通道的情况下，虽然没有缓冲区，但会因为阻塞操作占用 Goroutine。
- **在本例中：**
  - `ch2` 是无缓冲通道，阻塞的是发送 Goroutine，没有缓冲泄漏。

---

## **3. 调度器的负担**

- **泄漏资源：**
  - 阻塞 Goroutine 会对 Go 调度器增加负担，因为调度器需要跟踪所有活跃的 Goroutine 状态。
  - 即使阻塞的 Goroutine 不消耗 CPU，调度器仍会维护其元数据（例如堆栈指针、状态位等）。
  - 在长时间运行的程序中，随着 Goroutine 堆积，调度器的效率会降低。

---

## **4. 未释放的系统资源**

- **泄漏资源：**
  - 如果阻塞 Goroutine 内部分配了其他资源（例如文件句柄、网络连接、锁等），这些资源也可能无法正确释放。
  - 在我们的例子中，如果 `ch2` 发送 Goroutine 包含文件操作、数据库连接等，它会因为阻塞而无法释放这些系统资源。

---

## **泄漏总结**

在代码中，`ch2` 的发送 Goroutine 被阻塞，导致的泄漏情况包括：

1. 阻塞 Goroutine 永远无法退出，占用运行时内存和调度器资源。
2. 如果阻塞 Goroutine 申请了额外资源（如文件、网络、锁），这些资源也可能无法释放。

---

## **改进方法小结**

- **清晰的 Goroutine 生命周期管理：** 使用 `sync.WaitGroup` 或 `context.Context` 控制 Goroutine 退出。
- **确保通道操作的完整性：** 在 `select` 中处理所有通道，避免遗漏。
- **及时释放资源：** 关闭不再使用的通道或显式释放资源。
