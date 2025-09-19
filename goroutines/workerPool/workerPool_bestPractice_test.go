package workerPool

import (
	"fmt"
	"sync"
	"testing"
	"time"
)

// TestFixedWorkerPool 使用固定大小的workerPool限制最大同时的协程数量
func TestFixedWorkerPool(t *testing.T) {
	jobs := make(chan int, 10)
	var wg sync.WaitGroup

	// Init the workerPool, and the workLogicFunc
	for i := range 2 {
		go func(workerId int) {
			for job := range jobs {
				fmt.Printf("Worker %d processing job %d\n", workerId, job)
				time.Sleep(time.Second)
				fmt.Printf("Worker %d finished job %d\n", workerId, job)
				wg.Done()
			}
		}(i)
	}

	// !需要注意, 如果add跟wait存在race condition, 那么立刻调用的wait不会卡住
	// 这里一定要提前执行Add
	jobNum := 15
	wg.Add(jobNum)

	// send jobs
	go func() {
		for i := range jobNum {
			jobs <- i
			fmt.Printf("Filled job:%d\n", i)
		}
	}()

	// waiting for finished all jobs
	wg.Wait()
	time.Sleep(5 * time.Second)

	// gentle shutdown
	close(jobs)
	t.Log("TestFixedWorkerPool test completed")
}
