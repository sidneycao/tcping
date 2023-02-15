package utils

import (
	"context"
	"fmt"
	"sync"
	"time"
)

type Ping struct {
	target   Target        //tcp target
	stopOnce sync.Once     //stop once
	stopChan chan struct{} // stop chan

	minDuration   time.Duration // ping的最小时间
	maxDuration   time.Duration // ping的最大时间
	totalDuration time.Duration // ping的总花费时间
	counter       int           // 需要ping的次数
	total         int           // ping的总次数
	failed        int           // ping失败的次数
}

func (p *Ping) Stop() {
	p.stopOnce.Do(func() {
		close(p.stopChan)
	})
}

func (p *Ping) Done() <-chan struct{} {
	return p.stopChan
}

func (p *Ping) Ping() {
	defer p.Stop()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go func() {
		<-p.Done()
		cancel()
	}()

	timer := time.NewTimer(1)
	defer timer.Stop()

	stop := false
	for !stop {
		select {
		case <-timer.C:
			sts := p.target.Connect(ctx)
			fmt.Println(sts)
			if p.total++; p.counter > 0 && p.total > p.counter-1 {
				stop = true
			}
			timer.Reset(p.target.interval)
		case <-p.Done():
			stop = true
		}

	}
}

func NewPing(target Target, counter int) *Ping {
	return &Ping{
		target:   target,
		counter:  counter,
		stopChan: make(chan struct{}),
	}
}
