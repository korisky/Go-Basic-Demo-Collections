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
	t.Log("TestDynamicWithWaitGroupWorkerPool test completed")
}

// ReusablePool struct for reusable goroutine
// 可复用goroutine池的设计：
// workers: 存储可用worker的channel队列，每个元素是worker的work channel
// jobs: 接收待处理任务的channel
// quit: 用于优雅关闭所有goroutine的信号channel
// wg: 确保所有任务完成后才能关闭池的WaitGroup
type ReusablePool struct {
	workers chan chan int  // 可用worker队列，每个worker通过自己的channel注册
	jobs    chan int       // 任务队列
	quit    chan bool      // 关闭信号
	wg      sync.WaitGroup // 任务计数器
}

// dispatch help for dispatch the job / shutdown
// 调度器的复用机制核心：
// 1. 从jobs channel接收新任务
// 2. 从workers channel获取一个可用的worker（阻塞等待直到有worker可用）
// 3. 将任务发送给该worker处理
// 4. worker处理完后会重新注册到workers channel，实现循环复用
func (p *ReusablePool) dispatch() {
	for {
		select {
		case job, ok := <-p.jobs:
			// 检查jobs channel是否已关闭
			if !ok {
				return // jobs channel已关闭，退出调度器
			}
			// 获取一个可用worker - 如果没有可用worker会阻塞等待
			workerChan := <-p.workers
			// 将任务发送给该worker
			workerChan <- job
		case <-p.quit:
			return
		}
	}
}

// NewReusablePool return a ReusablePool
// 创建固定数量的worker goroutine实现复用：
// 1. 创建maxWorker个长期运行的goroutine
// 2. 每个goroutine处理完任务后不会销毁，而是重新注册等待新任务
// 3. 这样避免了频繁创建/销毁goroutine的性能开销
func NewReusablePool(maxWorker int) *ReusablePool {
	// 初始化池，workers channel的缓冲区大小等于worker数量
	pool := &ReusablePool{
		workers: make(chan chan int, maxWorker), // 缓冲区确保所有worker都能注册
		jobs:    make(chan int, 100),            // 修复并发问题：使用缓冲channel避免Submit阻塞
		quit:    make(chan bool),
	}
	// 启动调度器goroutine
	go pool.dispatch()
	// 创建固定数量的worker goroutine
	for i := range maxWorker {
		worker := NewReusableWorker(i+1, pool.workers, pool.quit, pool)
		worker.WorkerStart() // 每个worker启动一个长期运行的goroutine
	}
	return pool
}

// Submit job into the ReusablePool
// 修复并发问题的关键：
// 之前的bug：wg.Add(1)在channel send之前调用，如果send阻塞，
// wg计数器已增加但job未实际提交，导致wg.Wait()永久阻塞
// 现在通过缓冲channel确保Submit调用不会阻塞
func (p *ReusablePool) Submit(job int) {
	p.wg.Add(1)
	p.jobs <- job // 现在不会阻塞，因为jobs channel有缓冲区
}

// Close the whole ReusablePool
// 修复关闭顺序问题：
// 1. 先关闭jobs channel停止接收新任务
// 2. 等待所有当前任务完成
// 3. 最后关闭quit channel让worker退出
func (p *ReusablePool) Close() {
	close(p.jobs) // 停止接收新任务
	p.wg.Wait()   // 等待所有任务完成
	close(p.quit) // 然后让workers优雅退出
}

// ReusableWorker for ReusablePool
type ReusableWorker struct {
	ID      int
	work    chan int
	workers chan chan int
	quit    chan bool
	pool    *ReusablePool
}

// NewReusableWorker return new ReusableWorker
func NewReusableWorker(id int, workers chan chan int, quit chan bool, pool *ReusablePool) *ReusableWorker {
	return &ReusableWorker{
		ID:      id,
		work:    make(chan int),
		workers: workers,
		quit:    quit,
		pool:    pool,
	}
}

// WorkerStart 开启每个单独的ReusableWorker
// 关键的goroutine复用机制：
// 1. 每个worker在一个长期运行的goroutine中循环
// 2. worker完成任务后不会退出，而是重新注册自己到workers channel
// 3. 这样同样的goroutine可以被重复分配新任务，避免频繁创建/销毁goroutine的开销
func (r *ReusableWorker) WorkerStart() {
	go func() {
		for {
			// 注册可用worker - 这是复用的关键：worker完成任务后重新进入可用队列
			r.workers <- r.work
			// 持续处理
			select {
			case job := <-r.work:
				fmt.Printf("Worker %d processing job %d\n", r.ID, job)
				time.Sleep(time.Second)
				fmt.Printf("Worker %d finished job %d\n", r.ID, job)
				// 处理完任务后调用Done()
				r.pool.wg.Done()
			case <-r.quit:
				return
			}
		}
	}()
}

// TestReusableWorkerPool 重复使用的WorkerPool
// 测试goroutine复用机制：
// 1. 创建5个worker goroutine的池
// 2. 提交5个任务，观察同样的goroutine如何被重复使用
// 3. 任务数等于worker数，每个worker处理一个任务后重新等待
// 4. 验证所有任务完成后能正确关闭
func TestReusableWorkerPool(t *testing.T) {
	// 创建包含5个可复用worker的池
	pool := NewReusablePool(5)

	// 提交15个任务 - 测试worker复用机制
	for i := range 15 {
		pool.Submit(i)
	}

	// 等待所有任务完成并关闭池
	pool.Close()
	t.Log("TestReusableWorkerPool test completed")
}
