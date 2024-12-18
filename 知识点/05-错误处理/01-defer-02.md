 **`defer` 捕获变量的方式** 和 **执行时变量的作用域** 对比分析。

 闭包获取变量相当于引用传递，而非值传递。

---

### **1. 函数 `f0` 的分析**
```go
func f0() {
	fmt.Println("---------------- f0 -------------------")
	var i int // 定义一个变量 i

	for i = 0; i < 3; i++ {
		defer func() {
			fmt.Println("Deferred:", i) // 捕获的是变量 i 的引用
		}()
	}
}
```

#### **输出结果：**
```
---------------- f0 -------------------
Deferred: 3
Deferred: 3
Deferred: 3
```

#### **原因：**
- **变量 `i` 是一个循环外部定义的变量**，`defer` 中的匿名函数是一个闭包，它捕获了变量 `i` 的引用（地址）。
- 在 `for` 循环中，`i` 的值不断被更新。
- 当函数执行到 `defer` 时，延迟的匿名函数已经记录了 `i` 的引用，但直到函数 `f0` 返回时，`defer` 才执行，此时 `i` 的值已变为 `3`（循环结束后的值）。

#### **总结：**
- **闭包捕获的是变量 `i` 的引用**，所以 `defer` 执行时看到的是最终的值 `3`。

---

### **2. 函数 `f1` 的分析**
```go
func f1() {
	fmt.Println("---------------- f1 -------------------")
	for i := 0; i < 3; i++ {
		defer fmt.Println("Deferred: ", i)
	}
}
```

#### **输出结果：**
```
---------------- f1 -------------------
Deferred:  2
Deferred:  1
Deferred:  0
```

#### **原因：**
- **`i` 是 `for` 循环内部的变量**，且 `defer` 直接调用的是 `fmt.Println`。
- 每次执行 `defer` 时，当前的 `i` 值会被捕获并绑定到 `defer` 的调用中。
- 由于 `defer` 是后进先出（LIFO），输出顺序为 `2 -> 1 -> 0`。

#### **总结：**
- **`defer` 捕获的是当前 `i` 的值**（不是引用），所以每次输出的是循环中 `i` 的当时值，顺序遵循后进先出。

---

### **3. 函数 `f11` 的分析**
```go
func f11() {
	fmt.Println("---------------- f11 -------------------")
	for i := 0; i < 3; i++ {
		defer func(i int) {
			fmt.Println("Deferred: ", i)
		}(i)
	}
}
```

#### **输出结果：**
```
---------------- f11 -------------------
Deferred:  2
Deferred:  1
Deferred:  0
```

#### **原因：**
- 与 `f1` 不同，这里 `defer` 中使用了一个匿名函数，且通过参数显式传递了 `i` 的值。
- 每次循环中，当前的 `i` 值通过闭包参数传递给匿名函数，确保了每个 `defer` 调用中的 `i` 值是独立的。
- **后进先出（LIFO）** 的顺序导致输出顺序为 `2 -> 1 -> 0`。

#### **总结：**
- 通过显式传参的方式解决了闭包捕获引用的问题，每次 `defer` 捕获的是独立的 `i` 值。

---

### **核心区别总结**

| 函数名 | `defer` 捕获变量的方式 | 捕获的值 | 输出顺序 | 关键点                                   |
|--------|------------------------|----------|----------|----------------------------------------|
| `f0`   | 捕获 `i` 的引用         | 最终值 3 | 3, 3, 3  | 闭包捕获引用，变量更新导致值变化         |
| `f1`   | 捕获 `i` 的值           | 2, 1, 0 | 2, 1, 0  | `defer` 捕获当前值，后进先出             |
| `f11`  | 捕获显式传参的值         | 2, 1, 0 | 2, 1, 0  | 参数传递值，解决闭包捕获引用的问题       |

---

### **深入思考**

- **为何 `f11` 可以解决 `f0` 的问题？**
  - `f11` 使用了显式参数传递，这使得每次 `defer` 都能获取独立的值，而不是捕获循环变量的引用。
  - `f0` 中，闭包捕获的变量是共享的，而不是独立的。

- **`defer` 捕获行为的核心规则：**
  1. 如果直接调用函数（如 `f1`），捕获的是当前变量值。
  2. 如果是闭包（如 `f0`），捕获的是变量的引用。
  3. 显式参数传递（如 `f11`）可以避免闭包引用问题。

如果还有疑问，欢迎继续探讨！