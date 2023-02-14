package utils

import (
	"context"
	"fmt"
	"net"
	"strconv"
	"time"
)

type Target struct {
	host     string
	port     string
	dialer   *net.Dialer
	timeout  time.Duration
	interval time.Duration
	counter  int
}

type Stat struct {
	Connected bool
	Error     error
	Duration  time.Duration
	Address   string
}

func (t *Target) Ping(ctx context.Context) *Stat {
	var sts Stat
	timeout := DTimeout
	if t.timeout > 0 {
		timeout = t.timeout
	}
	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()
	start := time.Now()
	con, err := t.dialer.DialContext(ctx, "tcp", fmt.Sprintf("%s:%s", t.host, t.port))
	sts.Duration = time.Since(start)
	if err != nil {
		sts.Error = err
		if oe, ok := err.(*net.OpError); ok && oe.Addr != nil {
			sts.Address = oe.Addr.String()
		}
	} else {
		sts.Connected = true
		sts.Address = con.RemoteAddr().String()
	}
	return &sts
}

func NewTarget(host string, port string, timeout string, interval string, counter int) *Target {
	t := ParseDuartion(timeout)
	i := ParseDuartion(interval)
	return &Target{
		host:     host,
		port:     port,
		timeout:  t,
		interval: i,
		counter:  counter,
	}
}

func ParseDuartion(dString string) time.Duration {
	dInt, _ := strconv.Atoi(dString)
	return time.Duration(dInt) * time.Second
}
