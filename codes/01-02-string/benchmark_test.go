// string 一些常见操作的基准测试

package main

import (
	"bytes"
	"fmt"
	"strings"
	"testing"
)

/*
shell:
	cd 11-string
	go mod init 11-string
	go mod tidy

	go test -bench=. -benchmem
	go test -bench=. -run=^$ -benchmem
	go test -bench=BenchmarkStringConcatPlus -benchmem
	# -run=^$：表示不运行普通测试函数，只运行基准测试函数。
	# -benchmem：报告内存分配情况。
*/

/*
字符串拼接的性能测试

BenchmarkStringConcatPlus-8               128858             10026 ns/op           21080 B/op         99 allocs/op
BenchmarkStringConcatBuilder-8           1944698               649.8 ns/op          1016 B/op          7 allocs/op
BenchmarkStringConcatJoin-8               371589              2881 ns/op            4496 B/op          9 allocs/op
BenchmarkStringConcatFormat-8              50535             23403 ns/op           22681 B/op        199 allocs/op
BenchmarkStringConcatBuffer-8            1138202              1079 ns/op            1072 B/op          4 allocs/op

1. + 拼接：简单但效率低。
2. strings.Builder：高效的推荐方法。
3. bytes.Buffer 与 strings.Builder：性能接近。
3. fmt.Sprintf：功能强大，但性能相对较差。
4. append：性能不及 strings.Builder，比+ 或 sprintf 好点。

*/

func BenchmarkStringConcatPlus(b *testing.B) {
	for i := 0; i < b.N; i++ {
		s := ""
		for j := 0; j < 100; j++ {
			s += "test"
		}
	}
}

// strings.Builder 性能最优
func BenchmarkStringConcatBuilder(b *testing.B) {
	for i := 0; i < b.N; i++ {
		var builder strings.Builder
		for j := 0; j < 100; j++ {
			builder.WriteString("test")
		}
	}
}

func BenchmarkStringConcatJoin(b *testing.B) {
	for i := 0; i < b.N; i++ {
		var s []string
		for j := 0; j < 100; j++ {
			s = append(s, "test")
		}
		strings.Join(s, "")

	}
}

func BenchmarkStringConcatFormat(b *testing.B) {
	for i := 0; i < b.N; i++ {
		var s string
		for j := 0; j < 100; j++ {
			s = fmt.Sprintf("%s%s", s, "test")
		}
	}
}

// 使用 bytes.NewBufferString 进行拼接
func BenchmarkStringConcatBuffer(b *testing.B) {
	for i := 0; i < b.N; i++ {
		var buffer bytes.Buffer
		for j := 0; j < 100; j++ {
			buffer.WriteString("test")
		}
	}
}

/*
string 与 []byte 转换

BenchmarkStringToBytes-8                172994492                6.332 ns/op           0 B/op          0 allocs/op
BenchmarkBytesToString-8                222415846                5.166 ns/op           0 B/op          0 allocs/op
*/
func BenchmarkStringToBytes(b *testing.B) {
	s := "test string for conversion"
	for i := 0; i < b.N; i++ {
		_ = []byte(s)
	}
}

func BenchmarkBytesToString(b *testing.B) {
	bs := []byte("test string for conversion")
	for i := 0; i < b.N; i++ {
		_ = string(bs)
	}
}

/*
字符串分片
BenchmarkStringSlice-8                  1000000000               0.3889 ns/op          0 B/op          0 allocs/op
BenchmarkBytesWithCopy-8                103547361               11.26 ns/op            0 B/op          0 allocs/op
*/

func BenchmarkStringSlice(b *testing.B) {
	s := "This is a very long string for testing slicing performance."
	for i := 0; i < b.N; i++ {
		_ = s[:10]
	}
}

func BenchmarkBytesWithCopy(b *testing.B) {
	s := "This is a very long string for testing slicing performance."
	for i := 0; i < b.N; i++ {
		sub := s[:10]
		_ = string([]byte(sub)) // 通过复制避免原字符串引用
	}
}

/*
多字节操作
BenchmarkStringIndexAccess-8            1000000000               0.3636 ns/op          0 B/op          0 allocs/op
BenchmarkStringRangeLoop-8              49653664                28.30 ns/op            0 B/op          0 allocs/op
*/
func BenchmarkStringIndexAccess(b *testing.B) {
	s := "你好，世界！" // UTF-8 编码
	for i := 0; i < b.N; i++ {
		_ = s[0] // 直接按字节访问
	}
}

func BenchmarkStringRangeLoop(b *testing.B) {
	s := "你好，世界！" // UTF-8 编码
	for i := 0; i < b.N; i++ {
		for _, c := range s { // 循环遍历
			_ = c
		}
	}
}

/*
字符串查找和比较
BenchmarkStringContains-8               130739504                9.070 ns/op           0 B/op          0 allocs/op
BenchmarkStringCompare-8                293177896                4.976 ns/op           0 B/op          0 allocs/op
*/
func BenchmarkStringContains(b *testing.B) {
	s := "This is a test string for performance benchmarking."
	sub := "test"
	for i := 0; i < b.N; i++ {
		_ = strings.Contains(s, sub)
	}
}

func BenchmarkStringCompare(b *testing.B) {
	s1 := "Hello, World!"
	s2 := "Hello, Go!"
	for i := 0; i < b.N; i++ {
		_ = strings.Compare(s1, s2)
	}
}
