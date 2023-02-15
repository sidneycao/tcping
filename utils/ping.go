package utils

import (
	"context"
	"fmt"
	"math"
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
	// 只关闭一次
	p.stopOnce.Do(func() {
		close(p.stopChan)
	})
}

func (p *Ping) Done() <-chan struct{} {
	return p.stopChan
}

func (p *Ping) logSts(sts *Stat) {
	if sts.Duration < p.minDuration {
		p.minDuration = sts.Duration
	}
	if sts.Duration > p.maxDuration {
		p.maxDuration = sts.Duration
	}
	p.totalDuration += sts.Duration

	cSts := "failed"
	if sts.Connected {
		cSts = "connected"
	}

	if sts.Error != nil {
		p.failed++
		fmt.Printf("Ping %s:%s(%s) %s(%s) - time=%s\n", p.target.host, p.target.port, sts.Address, cSts, sts.Error.Error(), round(sts.Duration, 3))
	} else {
		fmt.Printf("Ping %s:%s(%s) %s - time=%s\n", p.target.host, p.target.port, sts.Address, cSts, round(sts.Duration, 3))
	}

}

func (p *Ping) Summarize() {
	fmt.Printf(
		`--- tcping statistics ---
%d probes sent, %d successful, %d failed.
round-trip min/avg/max = %s/%s/%s
`, p.total, p.total-p.failed, p.failed, round(p.minDuration, 3), round(p.totalDuration/time.Duration(p.total), 3), round(p.maxDuration, 3))
}

func (p *Ping) Ping() {
	// 运行结束后调用Stop
	// 关闭stopChan
	defer p.Stop()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	/**
	go func() {
		<-p.Done()
		cancel()
	}()
	***/

	// 初始化一个1ns的定时器
	// 相当于立即执行 之后会重设为参数的值
	timer := time.NewTimer(1)
	defer timer.Stop()

	stop := false
	p.minDuration = time.Duration(math.MaxInt64)
	for !stop {
		select {
		// 定时器计时结束后，会向timer.C中发出信号
		// 此时会进入这个case
		case <-timer.C:
			// 发起tcp连接，获得状态参数
			sts := p.target.Connect(ctx)
			p.logSts(sts)
			if p.total++; p.counter > 0 && p.total > p.counter-1 {
				stop = true
			}
			// 将定时器设置为参数中的时间间隔
			timer.Reset(p.target.interval)
		// stopChan关闭后，会进入这个case
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
