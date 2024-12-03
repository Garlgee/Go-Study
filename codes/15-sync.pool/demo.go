package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"sync"
	"sync/atomic"

	"math/rand"
)

func main() {
	pool := sync.Pool{
		New: func() interface{} {
			fmt.Println("Creating new object")
			return "new object"
		},
	}
	// 获取对象
	obj := pool.Get().(string)
	fmt.Println("Get:", obj)

	// 放回对象
	pool.Put("reused object")

	// 再次获取对象
	obj2 := pool.Get().(string)
	fmt.Println("Get:", obj2)

	// 获取一个新对象
	obj3 := pool.Get().(string)
	fmt.Println("Get:", obj3)

	f1()    // 缓冲区复用
	f2()    // 数据库连接池
	f3()    // JSON 编解码复用
	go f4() // 网络连接中的对象复用
	testHTTPReq()
}

// 场景1： 数据缓冲区复用
//
//	网络服务器缓冲区
//	数据处理中的临时存储
func f1() {
	fmt.Println("data buffer pool ------------")
	pool := sync.Pool{
		New: func() interface{} {
			return bytes.NewBuffer(make([]byte, 0, 1024))
		},
	}

	// 获取缓冲区
	buf := pool.Get().(*bytes.Buffer)
	buf.Reset() // 确保状态干净

	buf.WriteString("hello sync pool!")
	fmt.Printf("buffer[%p]: %s \n", buf, buf.String())

	// 使用完后放回池
	pool.Put(buf)

	// 再次获取缓冲区，避免重新分配内存
	buf2 := pool.Get().(*bytes.Buffer)
	fmt.Printf("Reused buffer[%p]:%s capacity: %d \n", buf2, buf2, buf2.Cap()) // 复用了之前分配的缓冲区
}

// 数据库连接池: sync.Pool 可用作轻量级对象池，比如存储数据库连接或客户端。
//
//	轻量级的对象池管理
//	临时需要高频获取和释放的资源（如短连接）
type DBConnection struct {
	ID int
}

func f2() {
	fmt.Println("DB connection pool -------------")
	pool := sync.Pool{
		New: func() interface{} {
			fmt.Println("Creating new DB connection")
			return &DBConnection{ID: rand.Intn(1000)} // 随机生成一个 ID 模拟连接
		},
	}
	// 获取连接
	conn := pool.Get().(*DBConnection)
	fmt.Println("Acquired DB connection ID:", conn.ID)

	// 使用后放回池
	pool.Put(conn)

	// 再次获取连接
	conn2 := pool.Get().(*DBConnection)
	fmt.Println("Reused DB connection ID:", conn2.ID)
}

// JSON 编解码复用：在频繁进行 JSON 编码解码的场景下，可以使用 sync.Pool 来复用 json.Encoder 或 json.Decoder，避免每次都创建新的实例。
// 适用场景：
//
//	高性能 HTTP API 的 JSON 数据处理
//	日志系统中 JSON 格式化

// 创建一个对象池，同时复用 encoder 和 buffer
type EncoderWithBuffer struct {
	Encoder *json.Encoder
	Buffer  *bytes.Buffer
}

func f3() {
	fmt.Println("JSON encoding/decoding pool ------------")

	// 创建 Encoder 的对象池
	pool := sync.Pool{
		New: func() interface{} {
			buf := &bytes.Buffer{}
			return &EncoderWithBuffer{
				Encoder: json.NewEncoder(buf),
				Buffer:  buf,
			}
		},
	}

	// 使用 Encoder
	pooled := pool.Get().(*EncoderWithBuffer)
	pooled.Encoder.SetIndent("", "  ") // 设置编码格式

	data := map[string]string{"key": "value"}
	pooled.Buffer.Reset()       // 确保 buffer 是干净的
	pooled.Encoder.Encode(data) // 编码到 buffer 中

	fmt.Println(pooled.Buffer.String()) // 打印 buffer 内容

	// 将对象放回池，供后续复用
	pool.Put(pooled)
}

// 网络连接的对象复用
// 处理高并发网络请求时，用 sync.Pool 来存储临时对象（如 HTTP 请求或响应的处理结构）。
type RequestHandler struct {
	RequestID int
}

// curl -X GET -d "example data" http://localhost:8088
func f4() {
	fmt.Println("HTTP request pool ------------")

	// RequestHandler计数器
	var count int32 = 0

	// 创建 RequestHandler 的对象池
	pool := sync.Pool{
		New: func() interface{} {
			fmt.Println("Creating new RequestHandler", atomic.AddInt32(&count, 1))
			return &RequestHandler{}
		},
	}

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		handler := pool.Get().(*RequestHandler)
		handler.RequestID = int(r.ContentLength) // 模拟处理 ID

		body, _ := io.ReadAll(r.Body)
		// 输出结果
		fmt.Fprintf(w, "Handling request with handler: %p, body: %s", handler, body)

		// 清理状态并放回池
		handler.RequestID = 0
		pool.Put(handler)
	})

	fmt.Println("Server running at :8088")
	go http.ListenAndServe(":8088", nil)
}

func testHTTPReq() {
	const url = "http://localhost:8088"
	var wg sync.WaitGroup

	// 模拟 100 个并发请求
	for i := 0; i < 100; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			resp, err := http.Post(url, "text/plain", bytes.NewBufferString(fmt.Sprintf("Request %d", id)))
			if err != nil {
				fmt.Println("Error:", err)
				return
			}
			body, _ := io.ReadAll(resp.Body)
			fmt.Printf("Response %d: '%s'\n", id, body)
			resp.Body.Close()
		}(i)
	}

	wg.Wait()
}
