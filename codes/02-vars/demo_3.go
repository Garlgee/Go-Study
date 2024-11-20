// demo_3.go
package main

import (
	"fmt"
)

func main() {

	fmt.Println("---")
	fmt.Print("输出到控制台不换行")
	fmt.Printf("name=%s,age=%d\n", "Tom", 30)
	fmt.Printf("name=%s,age=%d,height=%v\n", "Tom", 30, fmt.Sprintf("%.2f", 180.567))
}
