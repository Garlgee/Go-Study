// defer
package main

import (
	"fmt"
	"log"
	"os"
	"sync"
)

func main() {

	f0()
	f1()
	f11()
	result := f2()
	fmt.Println("f2 result:", result)
	result = f3()
	fmt.Println("f3 result:", result)

}

// 1. 资源清理
func readFile(fileName string) {
	file, err := os.Open(fileName)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()
}

// 2. 锁的释放
var mu sync.Mutex

func criticalSection() {
	mu.Lock()
	defer mu.Unlock() // 确保锁在函数结束时释放

	// 临界区逻辑
}

// 3. 日志打印
func example() {
	fmt.Println("Start")
	defer fmt.Println("End") // 无论函数如何退出，都打印 "End"

	// 中间逻辑
	fmt.Println("Doing something")
}

// 3. panic 处理
func recoverPanic() {
	defer func() {
		err := recover()
		if err != nil {
			log.Println("Recovered from panic:", err)
		}
	}()

	// 即使发生 panic，也会执行 defer 中的代码
	panic("This is a panic")
}

// defer 捕获变量的时机；
// 闭包获取变量相当于引用传递，而非值传递。
func f0() {
	fmt.Println("---------------- f0 -------------------")
	var i int // 定义一个变量 i

	for i = 0; i < 3; i++ {
		defer func() {
			fmt.Println("Deferred:", i) // 捕获的是变量 i 的引用
		}()
		// 区别
		// defer fmt.Println("Deferred: ", i)
	}
	// 输出：
	//Deferred: 3
	//Deferred: 3
	//Deferred: 3
}

func f1() {
	fmt.Println("---------------- f1 -------------------")
	for i := 0; i < 3; i++ {
		defer fmt.Println("Deferred: ", i)
	}
	// 输出：
	//Deferred:  2
	//Deferred:  1
	//Deferred:  0

}

func f11() {
	fmt.Println("---------------- f11 -------------------")
	for i := 0; i < 3; i++ {
		defer func(i int) {
			fmt.Println("Deferred: ", i)
		}(i)
	}
	// 输出：
	//Deferred:  2
	//Deferred:  1
	//Deferred:  0
}

// 与 return 交互
// 函数返回值 result 是命名的，defer 中对 result 的修改会影响返回值。
// f2 result: 15
func f2() (result int) {
	fmt.Println("---------------- f2 -------------------")
	defer func() {
		result += 10 // 修改返回值，result=5+10=15
	}()
	return 5 // result = 5
}

// 如果返回值是匿名的，defer 中的修改不会生效
// f3 result: 5
func f3() int {
	fmt.Println("---------------- f3 -------------------")
	result := 5
	defer func() {
		result += 10
		fmt.Println("defer result:", result)
	}()
	return result
}
