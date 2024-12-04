### **Golang 基础知识点整理**

#### **1. 基础语法**

- **数据类型**
  - 基本类型：`int`, `float32`, `string`, `bool`。
  - 复合类型：数组、切片（`slice`）、映射（`map`）、结构体（`struct`）。
  - 指针：`*T` 表示指向 `T` 类型的指针。
  - 接口：`interface` 是一种抽象类型，用于定义方法集合。
- **变量与常量**
  - 声明：`var`、短变量声明 `:=`。
  - 常量：`const` 声明不可改变的值。
- **控制结构**
  - 条件语句：`if-else`，支持简短的声明。
  - 循环：`for`，Go 中唯一的循环结构。
  - `switch`：支持多种匹配规则。
- **函数**
  - 多返回值：函数可以返回多个值。
  - 可变参数：`func f(args ...int)`。
  - 匿名函数与闭包。

---

#### **2. 数据结构**

- **数组**：固定大小的序列。
- **切片（Slice）**
  - 动态数组，底层依赖数组。
  - 支持追加（`append`）和拷贝（`copy`）。
- **Map**
  - 无序的键值对集合。
  - 使用 `make` 初始化。
- **结构体（Struct）**
  - 自定义类型，用于组织复杂数据。
  - 支持嵌套与匿名字段。
- **指针**
  - 支持直接操作地址。
  - 值传递与指针传递。

---

#### **3. 并发模型**

- **Goroutine**
  - 轻量级线程，使用 `go` 启动。
  - 每个 Goroutine 都有独立的栈空间。
- **Channel**
  - Goroutine 间通信的管道。
  - 无缓冲和带缓冲的 Channel。
  - 关闭 Channel：通过 `close(chan)` 关闭。
- **`select`**
  - 用于监听多个 Channel 的通信。
  - 支持超时（`time.After`）和阻塞操作。

---

#### **4. 包和模块**

- **包（Package）**
  - 每个 Go 文件必须声明包名。
  - `import` 用于导入其他包。
- **模块（Module）**
  - Go 1.11 引入的模块化管理工具。
  - 使用 `go mod init` 初始化模块，`go.mod` 文件用于管理依赖。

---

#### **5. 错误处理**

- **内置 `error` 接口**
  - `func (e *MyError) Error() string`。
- **`defer`、`panic`、`recover`**
  - `defer`：延迟执行，用于资源清理。
  - `panic`：引发运行时错误。
  - `recover`：捕获 `panic`，避免程序崩溃。

---

#### **6. 文件与 IO**

- **文件操作**
  - 使用 `os` 包进行文件读写。
  - `os.Open`、`os.Create`、`os.ReadFile`。
- **标准输入输出**
  - `fmt` 包提供 `Println`、`Printf`。
  - `bufio` 提供缓冲读写。
- **JSON 处理**
  - 使用 `encoding/json` 包解析与生成 JSON。

---

#### **7. 测试**

- **单元测试**
  - 使用 `testing` 包。
  - 测试文件以 `_test.go` 结尾。
  - 测试方法以 `Test` 开头。
- **性能测试**
  - 使用 `Benchmark` 函数。
  - 使用 `b.N` 控制循环次数。

---

### **Golang 进阶知识点整理**

#### **1. 高级数据结构**

- **Context**
  - 用于 Goroutine 间传递上下文信息。
  - 常见方法：`WithCancel`、`WithTimeout`、`WithValue`。
- **Sync 包**
  - `sync.Mutex`：互斥锁。
  - `sync.WaitGroup`：等待多个 Goroutine 结束。
  - `sync.Once`：确保只执行一次。
  - `sync.Map`：并发安全的 Map。
- **Time 包**
  - 定时器：`time.Timer` 和 `time.Ticker`。
  - 时间格式化：`time.Parse` 和 `time.Format`。

---

#### **2. Goroutine 调度**

- **Go 调度器（GMP 模型）**
  - G：Goroutine。
  - M：OS 线程。
  - P：逻辑处理器。
- **`runtime` 包**
  - 控制调度器行为，如 `GOMAXPROCS`。
  - 获取调用堆栈信息。

---

#### **3. 性能优化**

- **内存管理**
  - 内存分配：`make` 和 `new`。
  - 垃圾回收（GC）：自动回收不再使用的内存。
- **工具**
  - `pprof`：性能剖析工具。
  - `trace`：跟踪 Goroutine 行为。
- **逃逸分析**
  - 使用 `-gcflags="-m"` 检查变量是否逃逸到堆。

---

#### **4. 编译和链接**

- **`go build`**
  - 编译选项：`-race` 检测竞态条件，`-gcflags` 和 `-ldflags` 自定义编译和链接。
- **交叉编译**
  - 设置环境变量：`GOOS` 和 `GOARCH`。
  - 示例：`GOOS=linux GOARCH=amd64 go build main.go`。

---

#### **5. 反射**

- **`reflect` 包**
  - 动态获取类型和值。
  - 常用方法：`TypeOf` 和 `ValueOf`。
- **注意点**
  - 性能开销较高。
  - 不建议在性能关键代码中使用。

---

#### **6. 网络与通信**

- **标准库支持**
  - `net/http`：构建 HTTP 服务器和客户端。
  - `net/rpc`：远程过程调用。
- **第三方框架**
  - HTTP 服务：Gin、Echo。
  - RPC：gRPC。

---

#### **7. 并发控制**

- **高级并发模式**
  - 工作池模式：限制并发 Goroutine 数量。
  - `fan-in` 和 `fan-out` 模式。
- **注意事项**
  - 避免死锁。
  - 避免资源竞争：使用锁或 Channel。

---

#### **8. 安全与最佳实践**

- **避免内存泄漏**
  - 合理使用 `defer` 和 `Context`。
- **错误处理**
  - 错误信息应包含上下文。
- **代码规范**
  - 使用 `gofmt` 格式化代码。
  - 使用 `golint` 检查代码风格。

---

#### **工具链与生态**

- **构建与依赖管理**
  - `go mod`：模块化依赖管理。
- **调试工具**
  - `delve`：Go 语言调试器。
- **监控与日志**
  - `logrus`、`zap`：高性能日志库。
  - `prometheus`：监控指标采集与展示。

---
