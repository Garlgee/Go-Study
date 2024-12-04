### **`sync.WaitGroup` 底层实现**

`sync.WaitGroup` 是 Go 标准库中用于并发同步的工具，底层实现基于原子操作和条件变量（`sync.Cond`），保证了高效且线程安全的操作。

---

### **1. `sync.WaitGroup` 的内部结构**

以下是 `sync.WaitGroup` 的简化定义（源于 `runtime/sync` 包）：

```go
type WaitGroup struct {
    state1 [3]uint32 // 包含计数器和一些内部状态信息
    sema   uint32    // 信号量，用于阻塞和唤醒
}
```

- `state1[3]`：
  - `counter`（计数器）：记录需要等待的 Goroutine 数量。
  - 其他字段：用于存储状态信息（如是否有正在等待的 Goroutine）。
- `sema`：
  - 用于实现 Goroutine 的阻塞和唤醒机制。

---

### **2. 核心方法的底层逻辑**

#### **2.1. `Add(delta int)`**

- `Add` 用于增加或减少计数器值。
- 底层通过原子操作 `atomic.AddInt32` 修改计数器。
- 当计数器变为负数时会触发 panic（非法操作）。

**实现逻辑（简化版）**：

```go
func (wg *WaitGroup) Add(delta int) {
    // 使用原子操作修改计数器
    new := atomic.AddUint32(&wg.state1[0], uint32(delta))
    if new < 0 {
        panic("sync: negative WaitGroup counter")
    }

    // 如果减少后计数器为 0，唤醒所有等待 Goroutine
    if new == 0 {
        runtime_Semrelease(&wg.sema)
    }
}
```

---

#### **2.2. `Done()`**

- `Done` 是 `Add(-1)` 的快捷方式，用于减少计数器。

**实现逻辑**：

```go
func (wg *WaitGroup) Done() {
    wg.Add(-1)
}
```

---

#### **2.3. `Wait()`**

- `Wait` 会阻塞调用的 Goroutine，直到计数器变为 0。
- 底层通过检查计数器状态和信号量配合实现阻塞和唤醒。

**实现逻辑（简化版）**：

```go
func (wg *WaitGroup) Wait() {
    // 如果计数器已经为 0，直接返回
    if atomic.LoadUint32(&wg.state1[0]) == 0 {
        return
    }

    // 使用信号量阻塞
    runtime_Semacquire(&wg.sema)
}
```

- **`runtime_Semacquire`**：
  - 将当前 Goroutine 阻塞，直到有其他 Goroutine 调用 `runtime_Semrelease` 唤醒。

- **`runtime_Semrelease`**：
  - 唤醒阻塞在信号量上的 Goroutine。

---

### **3. 工作机制分析**

#### **3.1. 状态更新**

1. **`Add` 增加计数器**：
   - 主 Goroutine 调用 `Add` 增加任务计数。
   - 如果计数器大于 0，则表示有任务未完成。

2. **子 Goroutine 完成任务后调用 `Done`**：
   - 计数器减少，当计数器变为 0 时：
     - 唤醒所有在 `Wait` 上阻塞的 Goroutine。

3. **主 Goroutine 调用 `Wait`**：
   - 如果计数器大于 0，则阻塞。
   - 当计数器为 0 时被唤醒，继续执行。

---

### **4. 常见问题**

#### **4.1. 计数器为负**

当 `Add` 被错误使用（减少次数多于增加）时，计数器可能变为负数，触发 panic。

**示例问题**：

```go
var wg sync.WaitGroup
wg.Add(1)
wg.Done()
wg.Done() // panic: sync: negative WaitGroup counter
```

**解决方法**：

- 确保 `Add` 和 `Done` 的调用严格匹配。

---

#### **4.2. 提前调用 `Wait`**

在 `Add` 调用完成前调用 `Wait` 可能导致主 Goroutine 提前返回。

**示例问题**：

```go
var wg sync.WaitGroup
go func() {
    wg.Add(1) // Wait 调用后增加计数器
    defer wg.Done()
}()
wg.Wait() // 死锁，计数器没有更新
```

**解决方法**：

- 确保 `Add` 在所有 Goroutine 启动之前调用。

---

### **5. 优化与陷阱**

#### **5.1. 优化**

- **批量启动 Goroutine**：
  - 避免频繁调用 `Add` 和 `Done`。
- **任务队列**：
  - 使用任务队列分发工作，避免每个任务单独创建 Goroutine。

#### **5.2. 陷阱**

- **过度依赖 `WaitGroup`**：
  - 当任务复杂时，推荐使用 `context.Context` 管理任务的生命周期。
- **未配合通道使用**：
  - `WaitGroup` 不负责任务的取消和超时控制，需配合通道或 `context` 使用。

---

### **6. 总结**

- **`sync.WaitGroup` 的特点**：
  - 使用简单，适用于等待固定数量的任务完成。
  - 依赖原子操作和信号量，性能高效。
  
- **核心考点**：
  - 确保计数器增加和减少严格匹配。
  - 避免拷贝 `WaitGroup`，始终通过指针传递。