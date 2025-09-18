package goroutines

import (
	"fmt"
	"sync"
	"testing"
	"time"
)

func TestGoRoutinePoolWithChannel(t *testing.T) {

	// request group -> 模拟有5条请求进来
	queue := make(chan int, 5)
	go func() {
		// 这里对queue这条chan进行获取，
		// 而当chan为空的时候，相当于block住
		for id := range queue {
			Serve(id)
		}
	}()

	// 模拟并发5条请求，由于consumer-group至多3个goroutine处理，
	// 所以会卡一下另外的请求
	var wg sync.WaitGroup
	for i := range 5 {
		// 告知有一个事情要等待
		wg.Add(1)
		go func(reqId int) {
			// 处理完毕告知done
			defer wg.Done()
			fmt.Printf("Sending request %d\n", reqId)
			// 这里就是真正的发送request到queue
			queue <- i
		}(i)
	}

	// 有wg, 主线程进行等待所有
	wg.Wait()
	close(queue)
	time.Sleep(6 * time.Second)

	t.Log("Pool test completed")
}

var poolChan = make(chan int, 3)

func process(id int) {
	fmt.Printf("%d started\n", id)
	time.Sleep(time.Second)
	fmt.Printf("%d finished\n", id)
}

func Serve(id int) {
	// 外层, 我们向chan进行一个占用
	fmt.Printf("Request %d waiting for worker...\n", id)
	poolChan <- 1
	// 占用成功后，新建协程进行处理
	go func(reqId int) {
		// 协程中执行主要logic
		process(reqId)
		// 执行完毕后，剔除chan的占用
		<-poolChan
	}(id)
}
