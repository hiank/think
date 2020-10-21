package core

import (
	"bytes"
	"context"
	"strconv"
	"time"
)

//WithPort ip + port to string
func WithPort(ip string, port uint16) string {

	var buffer bytes.Buffer
	buffer.WriteString(ip)
	buffer.WriteByte(':')
	buffer.WriteString(strconv.FormatInt(int64(port), 10))
	return buffer.String()
}

//LoopRecv 循环读取消息
func LoopRecv(ticker *time.Ticker, delay time.Duration, recv func() (Message, error), call func(Message)) {

	for {
		msg, err := recv() //conn.Recv()
		if err != nil {
			ticker.Stop()
			break
		}
		ticker.Reset(delay)
		call(msg)
	}
}

//HealthLock 健康状态监控，如果是健康的话，会一直锁定
func HealthLock(ctx context.Context, ticker *time.Ticker, closeFunc func() error) error {

	select {
	case <-ctx.Done():
	case <-ticker.C:
	}
	return closeFunc()
}
