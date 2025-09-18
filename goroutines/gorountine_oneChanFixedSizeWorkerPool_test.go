package goroutines

import (
	"fmt"
	"sync"
	"testing"
	"time"
)

// TestOneChanFixedSizeWorkerPool 是使用1条channel控制worker大小的workerPool
func TestOneChanFixedSizeWorkerPool(t *testing.T) {

	// request group -> 模拟有5条请求进来
	queue := make(chan int, 5)
	go func() {
		// 这里对queue这条chan进行获取，
		// 而当chan为空的时候，相当于block住
		for id := range queue {
			chan1Serve(id)
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

var chan1WorkerPool = make(chan int, 2)
var chan1Once sync.Once

// init1ChanWorkerPool is about init the 'id' for worker
func init1ChanWorkerPool() {
	chan1WorkerPool <- 1
	chan1WorkerPool <- 2
}

func chan1Process(workerId, taskId int) {
	fmt.Printf("Worker %d started task %d\n", workerId, taskId)
	time.Sleep(time.Second)
	fmt.Printf("Worker %d finished task %d\n", workerId, taskId)
}

func chan1Serve(taskId int) {
	// 初始化一次
	chan1Once.Do(init1ChanWorkerPool)
	fmt.Printf("Task %d waiting for worker...\n", taskId)
	// 阻塞的等待获取worker
	workerId := <-chan1WorkerPool
	fmt.Printf("Task %d assigned to worker %d\n", taskId, workerId)
	// 占用成功后，新建协程进行处理
	go func() {
		// 协程中执行主要logic
		chan1Process(workerId, taskId)
		// 执行完毕后，剔除chan的占用, 这里我们把workerId放回去
		chan1WorkerPool <- workerId
		fmt.Printf("Worker %d released by finished task %d\n", workerId, taskId)
	}()
}
