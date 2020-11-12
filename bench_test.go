package taskq_test

import (
	"context"
	"sync"
	"testing"

	"github.com/airbrake/taskq/v2"
	"github.com/airbrake/taskq/v2/memqueue"
	"github.com/airbrake/taskq/v2/redisq"
)

func BenchmarkConsumerMemq(b *testing.B) {
	benchmarkConsumer(b, memqueue.NewFactory())
}

func BenchmarkConsumerRedisq(b *testing.B) {
	benchmarkConsumer(b, redisq.NewFactory())
}

var (
	once sync.Once
	q    taskq.Queue
	task *taskq.Task
	wg   sync.WaitGroup
)

func benchmarkConsumer(b *testing.B, factory taskq.Factory) {
	c := context.Background()

	once.Do(func() {
		q = factory.RegisterQueue(&taskq.QueueOptions{
			Name:  "bench",
			Redis: redisRing(),
		})

		task = taskq.RegisterTask(&taskq.TaskOptions{
			Name: "bench",
			Handler: func() {
				wg.Done()
			},
		})

		_ = q.Consumer().Start(c)
	})

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		for j := 0; j < 100; j++ {
			wg.Add(1)
			_ = q.Add(task.WithArgs(c))
		}
		wg.Wait()
	}
}
