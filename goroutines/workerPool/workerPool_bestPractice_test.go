package workerPool

import (
	"fmt"
	"sync"
	"testing"
	"time"
)

// TestFixedWorkerPool 使用固定大小的workerPool限制最大同时的协程数量
// 这里的限制，无关是否同时执行, 哪怕内部task可能要抢锁等待等，都不允许新建goroutine了
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

// TestDynamicWithWaitGroupWorkerPool 使用WaitGroup来动态控制'正在执行'的goRoutine
// 与固定大小的TestFixedWorkerPool不同，WaitGroup使用Semaphore的版本，是控制
// 允许多少goroutine同时执行，但当有新的task要处理，同样新建goroutine，而不是等待
func TestDynamicWithWaitGroupWorkerPool(t *testing.T) {

	// limit to 2 concurrent with semaphore
	sem := make(chan struct{}, 2)
	var wg sync.WaitGroup

	// 声明需要执行的task的logic
	submitJob := func(jobId int) {
		wg.Add(1)
		go func() {
			defer wg.Done()

			// 这里semaphore控制的是，同时能执行的goroutine数量
			// logging中可以看到，通常ProcessingJob可以有2个并发
			sem <- struct{}{}
			defer func() { <-sem }()

			fmt.Printf("Processing job %d\n", jobId)
			time.Sleep(time.Second)
			fmt.Printf("Finished job %d\n", jobId)
		}()
	}

	// 模拟一次性投敌15个job
	for i := range 15 {
		submitJob(i)
	}

	// waiting for finished all sem
	wg.Wait()

	// gentle shutdown
	close(sem)
	t.Log("TestFixedWorkerPool test completed")
}
