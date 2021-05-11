package util

import (
	"context"
	"log"
	"runtime"
)

func goWithRecovery(ctx context.Context, df func(ctx context.Context), f func()) {
	go func() {
		defer df(ctx)
		f()
	}()
}

func defaultF(ctx context.Context) {
	if e := recover(); e != nil {
		const size = 64 << 10
		buf := make([]byte, size)
		buf = buf[:runtime.Stack(buf, false)]
		log.Fatalf("goroutine panic %s", e)
	}
}

func GoParallel(ctx context.Context, f func()) {
	goWithRecovery(ctx, defaultF, f)
}
