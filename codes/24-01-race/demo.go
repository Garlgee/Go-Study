package main

import "fmt"

/*
go run -race .\demo.go
*/
func main() {
	//f1()
	f2()
}

func f1() { //未同步的变量访问：
	var counter int
	go func() {
		counter++
	}()
	fmt.Println(counter)
}

/*
==================
WARNING: DATA RACE
Write at 0x00c00012c078 by goroutine 7:
  main.f1.func1()
      D:/cath/code/Go-Study/codes/18-race/demo.go:15 +0x44

Previous read at 0x00c00012c078 by main goroutine:
  main.f1()
      D:/cath/code/Go-Study/codes/18-race/demo.go:17 +0xbe
  main.main()
      D:/cath/code/Go-Study/codes/18-race/demo.go:9 +0x24

Goroutine 7 (running) created at:
  main.f1()
      D:/cath/code/Go-Study/codes/18-race/demo.go:14 +0xae
  main.main()
      D:/cath/code/Go-Study/codes/18-race/demo.go:9 +0x24
==================
Found 1 data race(s)
*/

func f2() { //错误的切片访问：
	var data []int
	go func() {
		data = append(data, 1)
	}()
	fmt.Println(len(data))
}

/*
==================
WARNING: DATA RACE
Write at 0x00c000120060 by goroutine 7:
  main.f2.func1()
      D:/cath/code/Go-Study/codes/18-race/demo.go:46 +0xa6

Previous read at 0x00c000120060 by main goroutine:
  main.f2()
      D:/cath/code/Go-Study/codes/18-race/demo.go:48 +0xce
  main.main()
      D:/cath/code/Go-Study/codes/18-race/demo.go:10 +0x24

Goroutine 7 (running) created at:
  main.f2()
      D:/cath/code/Go-Study/codes/18-race/demo.go:45 +0xc4
  main.main()
      D:/cath/code/Go-Study/codes/18-race/demo.go:10 +0x24
==================
Found 1 data race(s)
*/
