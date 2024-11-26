package main

import (
	"fmt"
	"math/rand"
	"sync"
	"time"
)

type Item struct {
	ID    int
	Value int
}

func producer(producerID int, itemCh chan<- Item, wg *sync.WaitGroup) {
	defer wg.Done()

	for i := 0; i < 10; i++ {
		item := Item{ID: producerID*10 + i, Value: rand.Intn(100)}
		fmt.Println("Producer[", producerID, "] produce item: {ID: ", item.ID, ", Value: ", item.Value, "}")
		itemCh <- item
		time.Sleep(time.Millisecond * time.Duration(rand.Intn(50)))
	}
}

func consumer(consumerID int, itemCh <-chan Item, wg *sync.WaitGroup) {
	defer wg.Done()
	for item := range itemCh { // 需要等待通道关闭后，才能正常退出；所以生产结束后，需要关闭通道
		fmt.Println("							Consumer[", consumerID, "] consume item: {ID: ", item.ID, ", Value: ", item.Value, "}")
		time.Sleep(time.Millisecond * time.Duration(rand.Intn(100)))
	}
}

func main() {
	rand.Seed(time.Now().UnixNano())

	var nProducer = 2
	var nConsumer = 10
	itemCh := make(chan Item, 10)

	var wg1 sync.WaitGroup
	for i := 0; i < nProducer; i++ {
		wg1.Add(1)
		go producer(i, itemCh, &wg1)
	}

	var wg2 sync.WaitGroup
	for i := 0; i < nConsumer; i++ {
		wg2.Add(1)
		go consumer(i, itemCh, &wg2)
	}

	go func() {
		wg1.Wait()
		close(itemCh)
	}()

	wg2.Wait()

	fmt.Println("All producers and consumers have completed.")

}
