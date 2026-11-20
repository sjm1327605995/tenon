package ui

import (
	"sync"
	"sync/atomic"
	"testing"
)

func TestPostConcurrentWakeRegistration(t *testing.T) {
	drainPosts()
	setBackendWake(nil)
	t.Cleanup(func() {
		setBackendWake(nil)
		drainPosts()
	})

	const count = 128
	var wakes, executed atomic.Int64
	wake := func() { wakes.Add(1) }
	var wg sync.WaitGroup
	wg.Add(2)
	go func() {
		defer wg.Done()
		for i := 0; i < count; i++ {
			setBackendWake(wake)
		}
	}()
	go func() {
		defer wg.Done()
		for i := 0; i < count; i++ {
			Post(func() { executed.Add(1) })
		}
	}()
	wg.Wait()
	setBackendWake(nil)
	drainPosts()

	if got := executed.Load(); got != count {
		t.Fatalf("executed posts=%d, want %d", got, count)
	}
}
