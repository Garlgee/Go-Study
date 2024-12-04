package main

import "fmt"

func main() {
	// nil slice
	var s_1 []int
	fmt.Printf("len=%d cap=%d slice=%v\n", len(s_1), cap(s_1), s_1)

	// empty slice
	var s_2 = []int{}
	fmt.Printf("len=%d cap=%d slice=%v\n", len(s_2), cap(s_2), s_2)

	var s_3 = []int{1, 2}
	fmt.Printf("len=%d cap=%d slice=%v\n", len(s_3), cap(s_3), s_3)

	s_4 := []int{1, 2}
	fmt.Printf("len=%d cap=%d slice=%v\n", len(s_4), cap(s_4), s_4)

	var s_5 []int = make([]int, 5, 8)
	fmt.Printf("len=%d cap=%d slice=%v\n", len(s_5), cap(s_5), s_5)

	s_6 := make([]int, 5, 9)
	fmt.Printf("len=%d cap=%d slice=%v\n", len(s_6), cap(s_6), s_6)

}

// vim: set noet ts=4 sw=4:
