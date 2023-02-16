package utils

import (
	"context"
	"fmt"
	"math"
	"sync"
	"time"
)

type Ping struct {
	Target   Target        //tcp target
	stopOnce sync.Once     //stop once
	stopChan chan struct{} // stop chan

	MinDuration   time.Duration // ping的最小时间
	MaxDuration   time.Duration // ping的最大时间
	TotalDuration time.Duration // ping的总花费时间
	Counter       int           // 需要ping的次数
	Total         int           // ping的总次数
	Failed        int           // ping失败的次数
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
	if sts.Duration < p.MinDuration {
		p.MinDuration = sts.Duration
	}
	if sts.Duration > p.MaxDuration {
		p.MaxDuration = sts.Duration
	}
	p.TotalDuration += sts.Duration

	cSts := "failed"
	if sts.Connected {
		cSts = "connected"
	}

	if sts.Error != nil {
		p.Failed++
		fmt.Printf("Ping %s:%s(%s) %s(%s) - time=%s\n", p.Target.host, p.Target.port, sts.Address, cSts, sts.Error.Error(), Round(sts.Duration, 3))
	} else {
		fmt.Printf("Ping %s:%s(%s) %s - time=%s\n", p.Target.host, p.Target.port, sts.Address, cSts, Round(sts.Duration, 3))
	}

}

func (p *Ping) Summarize() {
	fmt.Printf(
		`--- tcping statistics ---
%d probes sent, %d successful, %d failed.
round-trip min/avg/max = %s/%s/%s
`, p.Total, p.Total-p.Failed, p.Failed, Round(p.MinDuration, 3), Round(p.TotalDuration/time.Duration(p.Total), 3), Round(p.MaxDuration, 3))
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
	p.MinDuration = time.Duration(math.MaxInt64)
	for !stop {
		select {
		// 定时器计时结束后，会向timer.C中发出信号
		// 此时会进入这个case
		case <-timer.C:
			// 发起tcp连接，获得状态参数
			sts := p.Target.Connect(ctx)
			p.logSts(sts)
			if p.Total++; p.Counter > 0 && p.Total > p.Counter-1 {
				stop = true
			}
			// 将定时器设置为参数中的时间间隔
			timer.Reset(p.Target.interval)
		// stopChan关闭后，会进入这个case
		case <-p.Done():
			stop = true
		}

	}
}

func NewPing(target Target, counter int) *Ping {
	return &Ping{
		Target:   target,
		Counter:  counter,
		stopChan: make(chan struct{}),
	}
}
