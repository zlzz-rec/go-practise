package main

import (
	"fmt"
	"math/rand"
	"time"
)

type Bucket struct {
	SuccessCount   uint32
	FailureCount   uint32
	TimeoutCount   uint32
	RejectionCount uint32
}

func (b *Bucket) Reset() {
	b.SuccessCount = 0
	b.FailureCount = 0
	b.TimeoutCount = 0
	b.RejectionCount = 0
}

type CircuitBreaker struct {
	Buckets []Bucket     // 桶
	Index   int          // 当前桶坐标
	Ticker  *time.Ticker // 内置定时器
	Volume  uint32       // 时间窗口内熔断器的生效起始请求数
	Rate    float64      // 失败率 大于此失败率请求会熔断
}

func NewCircuitBreaker(len int, volume uint32, rate float64) *CircuitBreaker {
	return &CircuitBreaker{
		Buckets: make([]Bucket, len+1), // 初始化桶数量+1, 多1个辅助桶方便环状桶数组旋转坐标时数据统计
		Index:   0,
		Ticker:  time.NewTicker(time.Second),
		Volume:  volume,
		Rate:    rate,
	}
}

// Start 内置定时器, 定时处理自己的桶坐标位移, 并清空辅助桶
func (c *CircuitBreaker) Start() {
	go func() {
		for range c.Ticker.C {
			// 计算下一个桶坐标
			next := c.Index + 1
			if next >= len(c.Buckets) {
				next = next % len(c.Buckets)
			}

			// 计算下下个桶坐标
			next2 := next + 1
			if next2 >= len(c.Buckets) {
				next2 = next2 % len(c.Buckets)
			}

			// 清空下下个桶计数
			c.Buckets[next2].Reset()
			// 切换当前桶坐标
			c.Index = next
		}
	}()
}

func (c *CircuitBreaker) Go(a func()) bool {
	a()
	// 0 1 2 3 4 5 6 7 8 9 10
	next := (c.Index%len(c.Buckets) + 1) % len(c.Buckets)

	totalCount, successCount, failureCount, timeoutCount, rejectionCount := uint32(0), uint32(0), uint32(0), uint32(0), uint32(0)
	for i, bucket := range c.Buckets {
		// 如果 index = 5, 取 0-5 + 7-10, 去掉6号辅助桶的统计
		if i == next {
			continue
		}
		totalCount += bucket.SuccessCount + bucket.FailureCount + bucket.TimeoutCount + bucket.RejectionCount
		successCount += bucket.SuccessCount
		failureCount += bucket.FailureCount
		timeoutCount += bucket.TimeoutCount
		rejectionCount += bucket.RejectionCount
	}
	fmt.Printf("%v\n", c.Buckets)

	// 判断当前时间窗口内请求数是否足够
	if totalCount < c.Volume {
		fmt.Println("请求数量不足, 默认允许请求")
		return false
	}

	// 判断当前时间窗口内请求失败率是否满足
	rate := float64(failureCount+timeoutCount+rejectionCount) / float64(totalCount)
	if rate < c.Rate {
		fmt.Printf("total:%d,success:%d,rate:%f, 允许请求\n", totalCount, successCount, rate)
		return false
	}
	fmt.Printf("total:%d,success:%d,rate:%f, 不允许请求\n", totalCount, successCount, rate)
	return true
}

func (c *CircuitBreaker) CurrentBucket() *Bucket {
	return &c.Buckets[c.Index]
}

func main() {
	c := NewCircuitBreaker(6, 5, 0.75)
	c.Start()

	ticker := time.NewTicker(100 * time.Millisecond)
	defer ticker.Stop()

	// 模拟每100ms进行一次请求
	for range ticker.C {
		if c.Go(func() {
			switch rand.Intn(4) {
			case 0:
				c.CurrentBucket().SuccessCount += 1
			case 1:
				c.CurrentBucket().FailureCount += 1
			case 2:
				c.CurrentBucket().TimeoutCount += 1
			case 3:
				c.CurrentBucket().RejectionCount += 1
			}
		}) {
			fmt.Println("do sth instead")
		}
	}

}
