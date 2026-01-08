package main

import (
	"log"
	"time"
)

// Scheduler 定时任务调度器
type Scheduler struct {
	pusher   *Pusher
	interval time.Duration
	once     bool
	stopCh   chan struct{}
}

// NewScheduler 创建新的调度器
func NewScheduler(pusher *Pusher, interval time.Duration, once bool) *Scheduler {
	return &Scheduler{
		pusher:   pusher,
		interval: interval,
		once:     once,
		stopCh:   make(chan struct{}),
	}
}

// Start 启动调度器
func (s *Scheduler) Start() {
	log.Printf("调度器启动，更新间隔: %v", s.interval)

	// 立即执行一次
	s.execute()

	if s.once {
		log.Println("单次执行模式，退出")
		return
	}

	// 创建定时器
	ticker := time.NewTicker(s.interval)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			s.execute()
		case <-s.stopCh:
			log.Println("调度器停止")
			return
		}
	}
}

// Stop 停止调度器
func (s *Scheduler) Stop() {
	close(s.stopCh)
}

// execute 执行推送任务
func (s *Scheduler) execute() {
	if err := s.pusher.FetchAndPush(); err != nil {
		log.Printf("执行推送任务失败: %v", err)
	}
}
