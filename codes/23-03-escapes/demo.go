package main

import "fmt"

/*
非逃逸：变量的生命周期在当前函数范围内即可管理（分配在栈上）。
逃逸：	变量的生命周期超出当前函数，或者被其他 Goroutine、闭包引用（分配在堆上）。
如何检查：
go build -gcflags="-m -l"
go run -gcflags="-m -l"
*/
func main() {
	f1()
	f2()
	f3() //闭包引用变量，逃逸
	f4() //接口类型导致的逃逸
	f5() //切片潜在的逃逸
}

// ---------------- 示例 1：变量逃逸到堆 ----------------
// x 的地址被返回，生命周期超出了函数的作用域，因此必须分配到堆。
func escapeExample() *int {
	x := 42   //	moved to heap: x
	return &x // 	返回指针，x 逃逸到堆
}

// go run -gcflags="-m"
// go build -gcflags="-m"
func f1() {
	ptr := escapeExample()
	fmt.Println(*ptr) // 输出 42;*ptr escapes to heap
}

// ---------------- 示例 2：变量不逃逸 ----------------
func noEscapeExample() int {
	x := 42
	return x // 没有返回指针，不逃逸
}
func f2() {
	x := noEscapeExample()
	fmt.Println(x) // 输出 42；x escapes to heap
}

// ------------------ 示例 3：闭包引用变量 ----------------
func closureExample() func() int {
	x := 42
	return func() int { return x } // 闭包引用变量，逃逸; func literal escapes to heap
}

// 优化：避免逃逸. 通过参数传递 x，避免了逃逸。
func closureExample2(x int) func() int {
	return func() int { return x }
}
func f3() {
	f := closureExample()
	fmt.Println(f()) // 输出 42

	f2 := closureExample2(42)
	fmt.Println(f2())
}

// -------------------- 示例 4：接口类型导致的逃逸 ----------------
// interface{} 赋值，会发生逃逸，优化方案是将类型设置为固定类型
func interfaceEscape() interface{} {
	x := 42
	return x // x 被装箱为接口类型，逃逸到堆 x escapes to heap

}
func f4() {
	x := interfaceEscape()
	fmt.Println(x) // 输出 42; x escapes to heap
}

// ----- 示例 5：切片的潜在逃逸 -----
func sliceEscape() []int {
	x := []int{1, 2, 3}
	return x[:] // 切片指向底层数组，数组逃逸到堆;[]int{...} escapes to heap
}
func f5() {
	x := sliceEscape()
	fmt.Println(x) // 输出 [1 2 3]; x escapes to heap
}
