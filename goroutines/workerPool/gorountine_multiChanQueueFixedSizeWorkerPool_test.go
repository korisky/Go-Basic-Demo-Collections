package workerPool

import (
	"fmt"
	"sync"
	"testing"
	"time"
)

// TestMultiChanFixedSizeWorkerPool 同样使用固定的worker数量, 但这里的job和result分channel更贴合真实使用
func TestMultiChanFixedSizeWorkerPool(t *testing.T) {

	// 声明job和result channels, 分别对应获取task的queue与处理完给出结果的queue
	jobs := make(chan int, 5)
	results := make(chan int, 5)

	// 异步发送任务到jobs
	go func() {
		for j := 1; j <= 5; j++ {
			fmt.Printf("Sending job %d\n", j)
			jobs <- j
		}
	}()

	// 异步处理jobs
	go func() {
		for i := 1; i <= 5; i++ {
			serveWithChan(jobs, results)
		}
	}()

	// 异步(阻塞)收集结果
	go func() {
		for a := 1; a <= 5; a++ {
			result := <-results
			fmt.Printf("Received result: %d\n", result)
		}
	}()

	// 等待足够的时间让所有工作完成
	time.Sleep(8 * time.Second)

	// 在主线程中安全地关闭channels (可选，对于demo来说)
	close(jobs)
	close(results)

	t.Log("Pool test completed")
}

var workerPool = make(chan int, 2)
var once sync.Once

// init1ChanWorkerPool is about init the 'id' for workerPool
func initWorkerPool() {
	workerPool <- 1
	workerPool <- 2
}

// processWithChan 对具体的workerId处理单个任务, 结果放入results
func processWithChan(workerId int, job int, results chan<- int) {
	fmt.Printf("Worker %d started job %d\n", workerId, job)
	time.Sleep(time.Second)
	fmt.Printf("Worker %d finished job %d\n", workerId, job)
	results <- job * 2
}

// serveWithChan 获取一个worker处理一个job (go中chan自然被passed为ref)
func serveWithChan(jobs <-chan int, results chan<- int) {
	// 初始化一次
	once.Do(initWorkerPool)

	// 先获取job, 然后再分配worker, 更合理
	// 获取一个job
	job, ok := <-jobs
	if !ok {
		return // 如果jobs channel已关闭，直接返回
	}
	// 阻塞的等待获取worker
	workerId := <-workerPool
	fmt.Printf("Worker %d assigned to job %d\n", workerId, job)

	// 占用成功后，新建协程进行处理
	go func() {
		// 协程中执行主要logic，处理单个job
		processWithChan(workerId, job, results)
		// 执行完毕后，剔除chan的占用, 这里我们把workerId放回去
		workerPool <- workerId
		fmt.Printf("Worker %d released\n", workerId)
	}()
}
