package main

import (
	"fmt"
	"sync"
)

func main() {
	var m sync.Map

	m.Store("name", "Alice")
	m.Store("age", 20)
	m.Store("sex", "female")

	if v, ok := m.Load("name"); ok {
		println(v)
	} else {
		println("not found")
	}

	val, loaded := m.LoadOrStore("name", "Bob")
	if !loaded {
		fmt.Println("name not found, store new value:", val)
	} else {
		fmt.Println("name found :", val)
	}

	val1, loaded1 := m.LoadOrStore("height", 175)
	if !loaded1 {
		fmt.Println("height not found, store new value:", val1)
	} else {
		fmt.Println("height found :", val1)
	}

	_, ok := m.Load("age")
	fmt.Println("Before delete, age exist =", ok)
	m.Delete("age")
	_, ok = m.Load("age")
	fmt.Println("After delete, age exist =", ok)

	m.Range(func(k, v interface{}) bool {
		fmt.Println("Key:", k, "Value:", v)
		return true // return false to stop the Range loope
	})

	f1()
	f2()
}

func f1() {
	sm := sync.Map{}
	sm.Store("key", 123)
	if value, ok := sm.Load("key"); ok {
		//fmt.Println(value + 1) // 会报错
		fmt.Println(value.(int) + 1)
	}

}

// LoadOrStore 的原子性
func f2() {
	fmt.Println("f2 ------------LoadOrStore 的原子性------------------")
	var sm sync.Map
	var wg sync.WaitGroup

	wg.Add(2)

	go func() {
		defer wg.Done()
		sm.LoadOrStore("key", "value1")
		fmt.Println("goroutine 1")
	}()

	go func() {
		defer wg.Done()
		sm.LoadOrStore("key", "value2")
		fmt.Println("goroutine 2")
	}()

	wg.Wait()

	// 最终值为先执行的 goroutine 的值。
	if value, ok := sm.Load("key"); ok {
		fmt.Println("sm[key]=", value)
	}
}
