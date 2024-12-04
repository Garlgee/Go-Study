### Go 中 `context` 相关知识整理

`context` 是 Go 标准库中专门为 **跨 API 和 Goroutine 边界传递请求范围的上下文信息** 而设计的包，通常用于：

- **控制 Goroutine 生命周期**、
- **传递元数据**（例如请求的用户认证信息、请求 ID）
  - 与 请求作用域 紧密相关的小型值（例如用户 ID、权限）。
  - 轻量级信息。

---

## 1. **`context` 的基本概念**

`context` 是 Go 的标准库，常用于：

1. **控制 Goroutine 生命周期**：
   - 通过传递上下文，可以在特定事件（如超时、取消等）发生时，优雅地停止 Goroutine。
2. **传递请求范围的数据**：
   - 使用 `context.WithValue` 可以在不同 API 和 Goroutine 之间共享元数据（注意性能问题）。

---

## 2. **主要函数与使用场景**

### **2.1 `context.Background()`**

**描述**: 返回一个空的上下文，常作为根上下文使用。

**使用场景**:

- 程序的主函数、初始化、测试等根级上下文。

**示例**:

```go
ctx := context.Background()
```

---

### **2.2 `context.WithCancel(parent)`**

**描述**: 创建一个可取消的上下文，调用返回的 `cancel` 函数会通知相关 Goroutine 退出。

**使用场景**:

- 控制 Goroutine 的生命周期。

**示例**:

```go
package main

import (
 "context"
 "fmt"
 "time"
)

func main() {
 ctx, cancel := context.WithCancel(context.Background())

 go func(ctx context.Context) {
  for {
   select {
   case <-ctx.Done():
    fmt.Println("Goroutine stopped.")
    return
   default:
    fmt.Println("Goroutine running...")
    time.Sleep(500 * time.Millisecond)
   }
  }
 }(ctx)

 time.Sleep(2 * time.Second)
 cancel() // 通知 Goroutine 停止
 time.Sleep(1 * time.Second)
}
```

---

### **2.3 `context.WithTimeout(parent, timeout)`**

**描述**: 创建一个带超时的上下文，在超时时会自动调用 `cancel`。

**使用场景**:

- 设置某个操作的最大执行时间。

**示例**:

```go
package main

import (
 "context"
 "fmt"
 "time"
)

func main() {
 ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
 defer cancel() // 避免资源泄漏

 go func(ctx context.Context) {
  select {
  case <-ctx.Done():
   fmt.Println("Operation timed out.")
   return
  }
 }(ctx)

 time.Sleep(3 * time.Second) // 模拟耗时操作
}
```

---

### **2.4 `context.WithDeadline(parent, deadline)`**

**描述**: 创建一个带截止时间的上下文，功能类似 `WithTimeout`。

**使用场景**:

- 当操作需要在特定时间点之前完成时使用。

**示例**:

```go
deadline := time.Now().Add(2 * time.Second)
ctx, cancel := context.WithDeadline(context.Background(), deadline)
defer cancel()

select {
case <-ctx.Done():
 fmt.Println("Operation exceeded deadline.")
}
```

---

### **2.5 `context.WithValue(parent, key, value)`**

**描述**: 返回一个携带值的上下文，用于在 API 边界之间传递元数据。

**使用场景**:

- 传递请求范围的额外信息。

**注意事项**:

- 不推荐将 `context` 用作数据存储，应尽量避免频繁使用。

**示例**:

```go
package main

import (
 "context"
 "fmt"
)

func main() {
 // ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
 // defer cancel() // 确保资源释放
 ctx := context.WithValue(context.Background(), "key", "value")

 value := ctx.Value("key")
 if value != nil {
  fmt.Println("Value:", value)
 }
}
```

---

## 3. **`context` 的考点与陷阱**

### **3.1 考点**

1. **Goroutine 的生命周期管理**:
   - 能否正确使用 `WithCancel` 或 `WithTimeout` 停止 Goroutine。
2. **上下文传递**:
   - 在函数调用链中能否正确传递 `context`。
3. **避免 `context.WithValue` 滥用**:
   - 不推荐使用 `context` 来存储大量或频繁访问的数据。
4. **`cancel` 函数的调用**:
   - 如果没有调用 `cancel`，会导致资源泄漏。

---

### **3.2 常见陷阱**

#### **陷阱 1：未调用 `cancel` 导致资源泄漏**

**描述**: 未调用 `cancel` 函数会导致 `context` 未被释放，从而浪费资源。

- 当上下文通过 context.WithCancel、context.WithTimeout 或 context.WithDeadline 创建时，背景中可能启动了一个辅助 Goroutine 来监听上下文的取消或超时事件。如果没有显式调用 cancel，这些 Goroutine 会一直存在，造成资源泄漏。
- 内存泄露：未释放的 Goroutine 会一直保留上下文中的引用数据，造成 对象无法被垃圾回收，从而占用内存。

**示例**:

```go
ctx, cancel := context.WithTimeout(context.Background(), time.Second)
// forget to call cancel()
```

**解决方法**:

- 使用 `defer cancel()`。

#### **陷阱 2：传递空上下文**

**描述**: 有些开发者可能会直接传递 `nil` 而不是 `context.Background()` 或 `context.TODO()`。

**问题**: 如果上下文为 `nil`，在某些函数中会导致 panic。

**解决方法**:

- 始终传递有效的上下文。

#### **陷阱 3：过度使用 `context.WithValue`**

**描述**: 滥用 `WithValue` 存储大量数据。

**问题**: 会导致上下文成为隐式的数据传递工具，增加复杂性和性能开销。

**解决方法**:

- 仅用于携带与上下文有关的元数据。
- 传递的数据应该是不可变的、小型的。

```go
package main

import (
 "context"
 "fmt"
)

// 请求中传递用户 ID
func main() {
 ctx := context.WithValue(context.Background(), "userID", 12345)
 processRequest(ctx)
}

func processRequest(ctx context.Context) {
 userID := ctx.Value("userID")
 fmt.Println("Processing request for user:", userID)
}

```

#### **陷阱 4：超时和取消未正确处理**

**描述**: 忽略 `ctx.Done()` 信号，导致 Goroutine 无法优雅退出。

**示例**:

```go
go func(ctx context.Context) {
 for {
  // 忽略 ctx.Done()
 }
}(ctx)
```

**解决方法**:

- 在 Goroutine 中始终监听 `ctx.Done()` 信号。

---

## 4. **总结**

`context` 是 Go 中进行并发编程的核心工具之一。正确理解其作用和用法，对于编写高效、安全的代码至关重要。

| 功能          | 场景                                                   | API 示例                                    |
|---------------|--------------------------------------------------------|--------------------------------------------|
| 根上下文       | 程序主入口、测试                                       | `context.Background()`                     |
| 可取消上下文   | 控制 Goroutine 的生命周期                              | `context.WithCancel(parent)`               |
| 超时上下文     | 限制操作最大时间                                       | `context.WithTimeout(parent, timeout)`     |
| 截止时间上下文 | 限制操作截止时间                                       | `context.WithDeadline(parent, deadline)`   |
| 值传递上下文   | 传递请求范围的元数据（尽量避免滥用）                    | `context.WithValue(parent, key, value)`    |

通过正确使用 `context`，可以有效地管理 Goroutine 的生命周期，减少资源泄漏和并发问题，提高代码的可读性和健壮性。