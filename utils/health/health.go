package health

import "context"

//MonitorHealth 监测健康
func MonitorHealth(ctx context.Context, handle Handle) {

	<-ctx.Done()
	handle()
}

//Handle 健康检查有问题后调用
type Handle func()
