package main

import "fmt"

func main() {
	s := []int{0, 1, 2, 3, 4, 5, 6}
	fmt.Printf("len=%d cap=%d slice=%v\n", len(s), cap(s), s)

	fmt.Println("s[1]=", s[1])
	fmt.Println("s[:]=", s[:])
	fmt.Println("s[1:]=", s[1:])
	fmt.Println("s[:4]=", s[:4])

	fmt.Println("s[0:3]=", s[0:3])
	fmt.Printf("len=%d cap=%d slice=%v\n", len(s[0:3]), cap(s[0:3]), s[0:3])

	// s[low:high:max]
	// low: 起始索引（包含），切片从该索引位置开始。
	// high: 结束索引（不包含），切片到该索引位置结束。
	// max: 用于限制切片的容量。 7 (max): 最大容量为 7，即从 s[0] 开始，最多到 s[6]。
	fmt.Println("s[0:3:7] ==", s[0:3:7])
	fmt.Printf("len=%d cap=%d slice=%v\n", len(s[0:3:7]), cap(s[0:3:7]), s[0:3:7])

	ss := []int{0, 1, 2, 3, 4}
	subSlice := ss[1:3:4] // 切片内容 [1, 2]，容量 4-1=3
	fmt.Printf("len=%d cap=%d slice=%v\n", len(subSlice), cap(subSlice), subSlice)

	f1()
	f2()
	f3()
	f4()
	f5()

}

// sub 是 s 的切片，修改 sub 会影响底层数组。
func f1() {
	s := []int{1, 2, 3, 4, 5}
	sub := s[1:3]
	sub[0] = 99
	fmt.Println(s)

}

// 切片扩容机制: 当切片容量不足时，会分配一个新的底层数组。通常扩容为原容量的 2 倍，但具体策略可能会因实现而异。
func f2() {
	s := []int{1, 2, 3, 4, 5}
	fmt.Printf("len=%d cap=%d slice=%v %p\n", len(s), cap(s), s, &s)
}

// 如何限制切片容量？为什么要这样做？使用 low:high:max 语法，可以限制切片 sub 的容量，防止影响原数组。
func f3() {
	s := []int{1, 2, 3, 4, 5}
	sub := s[1:3:3] // 限制容量为，即从 s[1] 开始，最多到 s[2]。
	fmt.Printf("len=%d cap=%d slice=%v\n", len(sub), cap(sub), sub)
}

// 切片深浅拷贝
// 使用 copy 函数实现深拷贝，不会共享底层数组。
func f4() {
	original := []int{1, 2, 3}
	copySlice := make([]int, len(original))
	copy(copySlice, original)
	fmt.Printf("original=%v copySlice=%v\n", original, copySlice)

	copySlice[0] = 99
	fmt.Printf("original=%v copySlice=%v\n", original, copySlice)

}

// append陷阱
func f5() {
	s := []int{1, 2, 3, 4} // 底层数组容量为 4
	sub := s[:2]           // sub = [1, 2]，长度为 2，容量为 4
	fmt.Printf("s len=%d cap=%d slice=%v\n", len(s), cap(s), s)
	fmt.Printf("sub len=%d cap=%d slice=%v\n", len(sub), cap(sub), sub)

	// 第一次 append，不超出容量，修改底层数组
	sub = append(sub, 99) // sub = [1, 2, 99]，底层数组变为 [1, 2, 99, 4]
	fmt.Println("After first append:")
	fmt.Println("s:", s)     // 输出: s: [1, 2, 99, 4]
	fmt.Println("sub:", sub) // 输出: sub: [1, 2, 99]

	// 第二次 append，超出容量，分配新数组
	sub = append(sub, 88, 77) // sub = [1, 2, 99, 88, 77]，分配新数组
	fmt.Println("\nAfter second append:")
	fmt.Println("s:", s)     // 输出: s: [1, 2, 99, 4]，不受影响
	fmt.Println("sub:", sub) // 输出: sub: [1, 2, 99, 88, 77]
}
