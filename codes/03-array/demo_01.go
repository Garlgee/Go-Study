package main

// 数组
// 数组不可动态变化问题，一旦声明了，其长度就是固定的。
// 数组是值类型问题，在函数中传递的时候是传递的值，如果传递数组很大，这对内存是很大开销。
// 数组赋值问题，同样类型的数组（长度一样且每个元素类型也一样）才可以相互赋值，反之不可以。
// len() 和 cap() 返回结果始终一样。
/*
var arr =  [5] int {1, 2, 3, 4, 5}
var arr_1 [5] int = arr
var arr_2 [6] int = arr // error
*/

import "fmt"

func main() {
	//  一维数组
	var arr_1 [5]int
	fmt.Println(arr_1)

	var arr_2 = [5]int{1, 3, 4, 6, 7}
	fmt.Println(arr_2)

	arr_3 := [5]int{1, 3, 4, 6, 7}
	fmt.Println(arr_3)

	arr_4 := [...]int{1, 3, 4, 6, 7}
	fmt.Println(arr_4)

	arr_5 := [5]int{0: 3, 1: 5, 4: 6}
	fmt.Println(arr_5)

	// 二维数组
	var arr_6 = [2][3]int{
		{1, 2, 3},
		{4, 5, 6},
	}
	fmt.Println(arr_6)

	arr_7 := [5][3]int{
		{1, 2, 3},
		{4, 5, 6},
	}
	fmt.Println(arr_7)
	fmt.Println("cap(arr_7)=", cap(arr_7))
	fmt.Println("len(arr_7)=", len(arr_7))

	arr_8 := [...][3]int{
		{1, 2, 3},
		{4, 5, 6},
	}
	fmt.Println(arr_8)

	fmt.Println("cap(arr_8)=", cap(arr_8))
	fmt.Println("len(arr_8)=", len(arr_8))

	//如果需要存储不同类型，可以使用 interface{}，但这会牺牲一些类型安全性和性能。
	arr1 := [2][3]interface{}{
		{1, "string", 3.14}, // 允许不同数据类型
		{true, 'A', nil},
	}

	fmt.Println(arr1)

}

// 传值
func modifyArr_01(a [5]int) {
	a[1] = 20
}

// 传指针
func modifyArr_02(a *[5]int) {
	a[1] = 20
}
