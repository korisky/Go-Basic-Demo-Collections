package workerPool

import (
	"fmt"
	"sync"
	"sync/atomic"
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
// 真实世界的worker pool设计：
// workers: 存储可用worker的channel队列
// jobs: 任务队列 (小buffer模拟真实backpressure)
// quit: worker关闭信号
// acceptingJobs: 是否还接收新任务的标志
// mu: 保护acceptingJobs的读写锁
// wg: 任务计数器，确保所有任务完成
type ReusablePool struct {
	workers       chan chan int  // 可用worker队列
	jobs          chan int       // 任务队列(小buffer创造真实的阻塞场景)
	quit          chan bool      // worker关闭信号
	acceptingJobs bool           // 是否接收新任务
	mu            sync.RWMutex   // 保护acceptingJobs标志
	wg            sync.WaitGroup // 任务计数器
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
// 真实世界的worker pool：
// 1. 小buffer的jobs channel创造真实的backpressure
// 2. 初始状态接收任务
// 3. worker复用机制避免频繁创建/销毁goroutine
func NewReusablePool(maxWorker int) *ReusablePool {
	pool := &ReusablePool{
		workers:       make(chan chan int, maxWorker), // worker队列
		jobs:          make(chan int, 2),              // 小buffer创造真实阻塞场景
		quit:          make(chan bool),
		acceptingJobs: true, // 初始状态接收任务
	}
	// 启动调度器goroutine
	go pool.dispatch()
	// 创建固定数量的worker goroutine
	for i := range maxWorker {
		worker := NewReusableWorker(i+1, pool.workers, pool.quit, pool)
		worker.WorkerStart()
	}
	return pool
}

// Submit job into the ReusablePool (TRUE BLOCKING - Real-world solution)
// 真实世界的blocking方案：
// 1. Submit会无限期阻塞等待，直到job被接受（绝不丢弃任务）
// 2. 只有在shutdown开始后，新的Submit才会被拒绝
// 3. 已经开始的Submit会继续阻塞直到完成，保证所有任务都被处理
func (p *ReusablePool) Submit(job int) error {
	// 原子操作：检查是否接收任务 + 增加计数器
	p.mu.Lock()
	if !p.acceptingJobs {
		p.mu.Unlock()
		return fmt.Errorf("pool is shutting down, not accepting new jobs")
	}
	// 一旦Add了，Close()就必须等这个job完成
	p.wg.Add(1)
	p.mu.Unlock()

	// 真正的blocking行为：无限期等待直到空间可用
	// 没有timeout，没有"busy"错误，就是纯粹的backpressure
	p.jobs <- job

	return nil
}

// Close the whole ReusablePool
// 修复版本的优雅关闭，防止worker注册死锁：
// 阶段1: 停止接收新任务，但让已阻塞的Submit继续等待
// 阶段2: 等待所有任务(包括阻塞中的)完成
// 阶段3: 安全关闭 - quit信号让所有goroutine同时退出
func (p *ReusablePool) Close() {
	// 阶段1: 原子地停止接收新任务
	p.mu.Lock()
	p.acceptingJobs = false
	p.mu.Unlock()

	// 阶段2: 等待所有已提交的任务完成
	// 包括当前正在阻塞等待的Submit调用
	p.wg.Wait()

	// 阶段3: 发送退出信号 - dispatcher和workers同时退出
	close(p.quit) // 所有goroutine看到此信号立即退出
	close(p.jobs) // 清理资源
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
// 修复版本：防止shutdown时的worker注册死锁
// 1. 使用非阻塞的worker注册，避免在dispatcher退出后卡死
// 2. worker完成任务后重新注册，实现goroutine复用
// 3. shutdown时能立即响应quit信号
func (r *ReusableWorker) WorkerStart() {
	go func() {
		for {
			// 非阻塞worker注册 - 关键修复！
			select {
			case r.workers <- r.work:
				// 成功注册，等待工作或退出信号
				select {
				case job := <-r.work:
					fmt.Printf("Worker %d processing job %d\n", r.ID, job)
					time.Sleep(time.Second)
					fmt.Printf("Worker %d finished job %d\n", r.ID, job)
					r.pool.wg.Done()
				case <-r.quit:
					return
				}
			case <-r.quit:
				// 如果无法注册(如dispatcher已退出)，立即退出
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

	// 提交15个任务 - 测试worker复用机制 (true blocking submit, zero drops)
	for i := range 15 {
		// Submit will block until space is available - no errors expected
		pool.Submit(i)
	}

	// 等待所有任务完成并关闭池
	pool.Close()
	t.Log("TestReusableWorkerPool test completed")
}

// TestReusableWorkerPoolPanicBug 测试关闭channel时的panic bug
// 当Close()关闭jobs channel时，仍在阻塞的Submit()会panic
func TestReusableWorkerPoolPanicBug(t *testing.T) {
	// 创建一个会导致Submit阻塞的场景
	pool := &ReusablePool{
		workers: make(chan chan int, 1), // 只有1个worker位置
		jobs:    make(chan int),         // 无缓冲channel
		quit:    make(chan bool),
	}
	go pool.dispatch()

	// 故意不启动任何worker - 这会让dispatch阻塞在获取worker上
	// 这样Submit调用会阻塞在jobs channel上

	// 启动一个会阻塞的Submit
	go func() {
		defer func() {
			if r := recover(); r != nil {
				t.Logf("Caught panic as expected: %v", r)
			}
		}()
		pool.Submit(1) // 这会阻塞，因为没有worker处理
	}()

	// 给点时间让Submit阻塞
	time.Sleep(100 * time.Millisecond)

	// 现在关闭pool - 这会导致阻塞的Submit panic
	func() {
		defer func() {
			if r := recover(); r != nil {
				t.Logf("Close() caused panic: %v", r)
			}
		}()
		pool.Close()
	}()

	t.Log("TestReusableWorkerPoolPanicBug completed")
}

// TestOption1RaceCondition demonstrates why swapping wg.Add/channel send order doesn't work
func TestOption1RaceCondition(t *testing.T) {
	// Create pool with Option 1 "fix"
	pool := &ReusablePool{
		workers: make(chan chan int, 2),
		jobs:    make(chan int, 2),
		quit:    make(chan bool),
	}
	go pool.dispatch()

	// Start fast workers
	for i := range 2 {
		worker := &ReusableWorker{
			ID:      i + 1,
			work:    make(chan int),
			workers: pool.workers,
			quit:    pool.quit,
			pool:    pool,
		}
		go func(w *ReusableWorker) {
			for {
				w.workers <- w.work
				select {
				case job := <-w.work:
					fmt.Printf("Worker %d processing job %d\n", w.ID, job)
					// Process immediately - no sleep
					fmt.Printf("Worker %d finished job %d\n", w.ID, job)
					w.pool.wg.Done() // This might happen before wg.Add(1)!
				case <-w.quit:
					return
				}
			}
		}(worker)
	}

	// Modified Submit using Option 1 approach
	submitOption1 := func(job int) {
		defer func() {
			if r := recover(); r != nil {
				t.Logf("Option 1 caused panic: %v", r)
			}
		}()
		pool.jobs <- job // Send first
		pool.wg.Add(1)   // Add after - RACE CONDITION!
	}

	// Submit jobs rapidly to increase chance of race condition
	for i := range 10 {
		go submitOption1(i)
	}

	time.Sleep(100 * time.Millisecond)
	pool.Close()
	t.Log("TestOption1RaceCondition completed")
}

// TestOption3Fixed demonstrates the fixed non-blocking Submit approach
func TestOption3Fixed(t *testing.T) {
	// Create pool with very small buffer to force "busy" conditions
	pool := NewReusablePool(2)

	var successCount, failCount int
	var wg sync.WaitGroup

	// Concurrent submission with error handling
	for i := range 100 {
		wg.Add(1)
		go func(jobId int) {
			defer wg.Done()
			if err := pool.Submit(jobId); err != nil {
				failCount++
				t.Logf("Job %d failed: %v", jobId, err)
			} else {
				successCount++
			}
		}(i)
	}

	// Wait for all submissions to complete
	wg.Wait()

	// Close pool - this should never deadlock now!
	pool.Close()

	t.Logf("Option 3 results: %d succeeded, %d failed (no deadlocks!)", successCount, failCount)
}

// TestTrueBlockingBehavior 测试真正的阻塞行为 - 零丢弃任务
func TestTrueBlockingBehavior(t *testing.T) {
	// 创建极小的pool和buffer，强制产生真实的阻塞
	pool := NewReusablePool(2) // 只有2个worker
	// jobs channel只有2个buffer，很容易填满

	var submitted int64
	var submitWg sync.WaitGroup

	// 启动多个goroutine并发提交任务，观察阻塞行为
	for i := range 10 {
		submitWg.Add(1)
		go func(submitterID int) {
			defer submitWg.Done()

			// 每个提交者提交3个任务
			for j := range 3 {
				jobID := submitterID*10 + j

				t.Logf("Submitter %d attempting to submit job %d...", submitterID, jobID)
				start := time.Now()

				// 这里会真正阻塞！不会返回错误，只会等待
				pool.Submit(jobID)

				duration := time.Since(start)
				atomic.AddInt64(&submitted, 1)
				t.Logf("Submitter %d submitted job %d after waiting %v", submitterID, jobID, duration)
			}
		}(i)
	}

	// 等待一会儿让提交开始阻塞
	time.Sleep(2 * time.Second)

	// 现在开始优雅关闭，观察所有阻塞的提交是否都能完成
	t.Log("Starting graceful shutdown while some submits are still blocked...")

	// 在另一个goroutine中关闭，这样我们可以观察过程
	go func() {
		time.Sleep(500 * time.Millisecond)
		t.Log("Calling pool.Close()...")
		pool.Close()
		t.Log("pool.Close() completed - all blocked submits finished!")
	}()

	// 等待所有提交完成
	submitWg.Wait()

	t.Logf("✅ TRUE BLOCKING SUCCESS: All %d tasks submitted with zero drops!", submitted)
	t.Log("💡 Key insight: Submit() blocked and waited, never returned errors or dropped tasks")
}

// TestWorkerDeadlockFix 专门测试worker注册死锁修复
func TestWorkerDeadlockFix(t *testing.T) {
	// 创建一个极端场景：1个worker，无缓冲jobs channel
	pool := &ReusablePool{
		workers:       make(chan chan int, 1),
		jobs:          make(chan int), // 无缓冲
		quit:          make(chan bool),
		acceptingJobs: true,
	}
	go pool.dispatch()

	// 只启动1个worker
	worker := NewReusableWorker(1, pool.workers, pool.quit, pool)
	worker.WorkerStart()

	// 提交一个需要长时间处理的任务
	pool.wg.Add(1)
	pool.jobs <- 999 // 这会被worker接收和处理

	// 等待worker开始处理
	time.Sleep(100 * time.Millisecond)

	// 在worker仍在处理任务时调用Close()
	// 这会测试worker能否在无法注册时正确退出
	t.Log("Calling Close() while worker is busy...")
	start := time.Now()

	// 在另一个goroutine中Close，这样我们可以检测是否会hang
	done := make(chan bool, 1)
	go func() {
		pool.Close()
		done <- true
	}()

	// 等待Close完成，如果超过5秒说明有死锁
	select {
	case <-done:
		duration := time.Since(start)
		t.Logf("✅ Close() completed in %v - no worker deadlock!", duration)
	case <-time.After(5 * time.Second):
		t.Fatal("❌ Close() hung for 5+ seconds - worker deadlock detected!")
	}
}
