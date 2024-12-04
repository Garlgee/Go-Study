### **Go 语言 Struct 基础**

在 Go 中，`struct` 是一种用户自定义的数据类型，用于将不同类型的数据组合在一起。

### **Struct 定义与使用**

#### 定义 Struct

```go
type Person struct {
    Name string
    Age  int
}
```

#### 创建 Struct 实例

1. **基本方式**：
   ```go
   p := Person{Name: "Alice", Age: 25}
   ```

2. **省略字段名（字段按顺序匹配）**：
   ```go
   p := Person{"Alice", 25}
   ```

3. **创建指针**：
   ```go
   p := &Person{Name: "Alice", Age: 25}
   ```

4. **匿名字段（嵌套结构体）**：
   ```go
   type Address struct {
       City    string
       ZipCode string
   }

   type Employee struct {
       Name    string
       Address // 匿名字段
   }

   e := Employee{Name: "Bob", Address: Address{City: "NY", ZipCode: "10001"}}
   fmt.Println(e.City) // 可直接访问匿名字段的成员
   ```

---

### **JSON 转换**

Go 提供 `encoding/json` 包用于将结构体与 JSON 相互转换。

#### **结构体转 JSON**

```go
package main

import (
    "encoding/json"
    "fmt"
)

type Person struct {
    Name string `json:"name"` // 自定义 JSON 字段名
    Age  int    `json:"age"`
}

func main() {
    p := Person{Name: "Alice", Age: 25}
    jsonData, _ := json.Marshal(p) // 转换为 JSON 字符串
    fmt.Println(string(jsonData))  // 输出: {"name":"Alice","age":25}
}
```

#### **JSON 转结构体**

```go
func main() {
    jsonStr := `{"name":"Alice","age":25}`
    var p Person
    _ = json.Unmarshal([]byte(jsonStr), &p) // 将 JSON 字符串转换为结构体
    fmt.Println(p)                         // 输出: {Alice 25}
}
```

#### **常见 JSON 标签**

- `omitempty`: 忽略空值字段
- `-`: 忽略字段
- 示例：
  ```go
  type Person struct {
      Name string `json:"name"`
      Age  int    `json:"age,omitempty"` // 如果 Age 为 0，则不会出现在 JSON 中
      Temp string `json:"-"`             // 忽略此字段
  }
  ```

---

### **Struct 常见面试题与陷阱**

#### 1. **Struct 值拷贝与指针**
   
   Struct 是值类型，赋值时会**拷贝所有字段**：

   ```go
   type Person struct {
       Name string
       Age  int
   }

   func main() {
       p1 := Person{Name: "Alice", Age: 25}
       p2 := p1 // 值拷贝
       p2.Name = "Bob"
       fmt.Println(p1.Name) // 输出: Alice，不受影响
   }
   ```

   **陷阱**：使用指针可以避免拷贝，减少内存消耗，但要注意数据共享问题。

---

#### 2. **JSON 反序列化忽略未导出字段**

未导出的字段（首字母小写）无法参与 JSON 转换：

```go
type Person struct {
    Name string `json:"name"`
    age  int    // 未导出字段
}

func main() {
    jsonStr := `{"name":"Alice","age":25}`
    var p Person
    _ = json.Unmarshal([]byte(jsonStr), &p)
    fmt.Println(p) // 输出: {Alice 0}，age 未赋值
}
```

---

#### 3. **字段冲突**

```go
type Address struct {
    City string
}

type Employee struct {
    Name string
    Address
    City string // 与 Address 中的 City 字段冲突
}

func main() {
    e := Employee{Name: "Bob", City: "NY", Address: Address{City: "SF"}}
    fmt.Println(e.City)        // 输出: NY
    fmt.Println(e.Address.City) // 输出: SF
}
```

---

#### 4. **嵌套结构体与 JSON**

```go
type Address struct {
    City string `json:"city"`
}

type Employee struct {
    Name    string  `json:"name"`
    Address Address `json:"address"`
}

func main() {
    e := Employee{Name: "Bob", Address: Address{City: "NY"}}
    jsonData, _ := json.Marshal(e)
    fmt.Println(string(jsonData)) // 输出: {"name":"Bob","address":{"city":"NY"}}
}
```

---

#### 5. **`omitempty` 的误区**

`omitempty` 只会忽略字段的“零值”：

```go
type Person struct {
    Name  string `json:"name,omitempty"`
    Age   int    `json:"age,omitempty"`
    Alive bool   `json:"alive,omitempty"`
}

func main() {
    p := Person{Name: "", Age: 0, Alive: false}
    jsonData, _ := json.Marshal(p)
    fmt.Println(string(jsonData)) // 输出: {}，所有字段都被忽略
}
```

---

#### 6. **深拷贝与浅拷贝**

Go 结构体默认是浅拷贝（拷贝字段值），但嵌套指针或切片时需要注意：

```go
type Person struct {
    Name    string
    Friends []string
}

func main() {
    p1 := Person{Name: "Alice", Friends: []string{"Bob", "Charlie"}}
    p2 := p1
    p2.Friends[0] = "David" // 修改切片
    fmt.Println(p1.Friends) // 输出: [David Charlie]，p1 也受影响
}
```

**解决方案**：手动深拷贝。

---

### **总结**

- **基础**：
  - 结构体是值类型，默认按值传递。
  - 通过 JSON 标签控制序列化和反序列化行为。

- **常见面试考点**：
  - 结构体值拷贝与指针引用。
  - JSON 转换及未导出字段的处理。
  - 嵌套结构体、字段冲突、`omitempty` 使用。

- **陷阱**：
  - 共享底层数组/切片时的数据变动。
  - JSON 处理中的忽略未导出字段和零值逻辑。