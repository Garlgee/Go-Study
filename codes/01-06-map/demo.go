package main

import (
	"encoding/json"
	"fmt"
	"strings"
	"sync"
)

func main() {
	// 声明但未初始化（nil map）
	var m map[string]int
	fmt.Println(m == nil) // prints "true"

	//使用 make 初始化：
	m1 := make(map[string]int)
	fmt.Println(m1 == nil) // prints "false"
	fmt.Println(len(m1))   // prints "0"

	//使用字面量初始化
	m2 := map[string]int{"one": 1, "two": 2}
	fmt.Printf("m2=%v\n", m2) // prints "{one 1 two 2}"

	// m必须初始化后才可以赋值
	m = make(map[string]int)
	// 添加或更新
	m["Alice"] = 26
	// 获取键的值
	fmt.Println("age=", m["Alice"])

	// 检查键是否存在
	value, ok := m["Bob"]
	if !ok {
		fmt.Println("Bob not found")
	} else {
		fmt.Printf("Bob's age=%d", value)
	}

	// 删除键
	delete(m, "Alice")

	m["Tom"] = 26
	// 遍历
	for key, value := range m {
		fmt.Printf("%s=%d\n", key, value)
	}

	map_json()
	f1()
	f2()
	//f3() //panic
	f4()
	f5()
	f6()

	f7() //并发异常
	f8()

}

func map_json() {
	m := map[string]int{"Alice": 25, "Bob": 30}
	jsonData, _ := json.Marshal(m)
	fmt.Println("json data :", string(jsonData))

	jsonStr := `{"Alice":25,"Bob":30}`
	var m2 map[string]int
	json.Unmarshal([]byte(jsonStr), &m2)
	fmt.Println("map data :", m2["Alice"])
}

// map 是引用类型，赋值时会共享底层数据
func f1() {
	m1 := map[string]int{"Alice": 25}
	m2 := m1
	m2["Alice"] = 30
	fmt.Println(m1["Alice"]) // 输出: 30，m1 和 m2 共享底层数据
}

// 并发：Go 的 map 在多个 goroutine 中并发读写时不是线程安全的。直接操作可能导致运行时崩溃。
// 使用 sync.Mutex 保护 map。
// 使用 sync.Map 提供线程安全的 map
func f2() {
	var m sync.Map
	m.Store("Alice", 25)
	go func(mm *sync.Map) {
		value, _ := mm.Load("Alice")
		fmt.Println(value)
	}(&m)
}

// nil Map 的陷阱
func f3() {
	var m map[string]int
	m["Alice"] = 25 // panic: assignment to entry in nil map
}

// 删除不存在的建
func f4() {
	m := map[string]int{"Alice": 25}
	delete(m, "Bob") // 没有 panic，不会报错
}

// 反序列化时，如果目标类型是 map[string]interface{}，注意数字类型可能会被解析为 float64
// 可能会，我没测出来
func f5() {
	jsonStr := `{"Alice":25.0,"Bob":30}`
	var m map[string]interface{}
	json.Unmarshal([]byte(jsonStr), &m)
	fmt.Println("map data :", m["Alice"]) // 输出为25；而不是输出：25.0，而不是 25
	fmt.Println("map data :", m["Bob"])   // 输出为30；而不是输出：30.0，而不是 30

	// 解决，显示的使用json.Number
	decoder := json.NewDecoder(strings.NewReader(`{"Age":25}`))
	decoder.UseNumber()
	var m2 map[string]interface{}
	_ = decoder.Decode(&m2)
	age, _ := m2["Age"].(json.Number).Int64()
	fmt.Println(age) // 输出: 25
}

// 遍历顺序不固定，每次遍历的顺序可能不同；因为 map 是无序的
func f6() {
	m := map[string]int{"00": 25, "01": 30, "02": 35, "03": 40}
	for k, v := range m {
		fmt.Printf("%s=%d\n", k, v)
	}
}

func f7() {
	fmt.Println("f7 -------------------- 并发：sync.Mutex --------------------")
	m := make(map[int]int)
	var mu sync.Mutex
	var wg sync.WaitGroup

	wg.Add(2)

	go func() {
		fmt.Println("goroutine 1: start writing")
		defer wg.Done()
		for i := 0; i < 10000; i++ {
			mu.Lock()
			m[i] = i * 2
			mu.Unlock()
		}
		fmt.Println("goroutine 1: finished writing")
	}()

	go func() {
		fmt.Println("goroutine 2: start reading")
		defer wg.Done()
		mu.Lock()
		fmt.Println(len(m), "keys in map")
		for k, v := range m {
			fmt.Println(k, v)
		}
		mu.Unlock()
		fmt.Println("goroutine 2: finished reading")
	}()

	wg.Wait()
	fmt.Println("f7: all goroutines finished")
}

func f8() {
	fmt.Println("f8 -------------------- 并发: sync.Map --------------------")
	var m sync.Map
	var wg sync.WaitGroup

	wg.Add(2)

	go func() {
		fmt.Println("goroutine 1: start writing")
		defer wg.Done()
		for i := 0; i < 10000; i++ {
			m.Store(i, i*2)
		}
		fmt.Println("goroutine 1: finished writing")
	}()

	go func() {
		fmt.Println("goroutine 2: start reading")
		defer wg.Done()
		m.Range(func(k, v interface{}) bool {
			fmt.Println(k, v)
			return true
		})
		fmt.Println("goroutine 1: finished reading")
	}()

	wg.Wait()
	fmt.Println("f8: all goroutines finished")
}
