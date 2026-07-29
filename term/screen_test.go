package term

import (
	"errors"
	"sync"
	"testing"
	"time"
)

func withScreenSizeGetter(t *testing.T, fn func() (int, int, error)) {
	t.Helper()
	original := screenSizeGetter
	originalInterval := resizePollInterval
	screenSizeGetter = fn
	resizePollInterval = time.Millisecond
	t.Cleanup(func() {
		screenSizeGetter = original
		resizePollInterval = originalInterval
	})
}

func TestOnResizeInvokesCallbackOnChange(t *testing.T) {
	var mu sync.Mutex
	sizes := []struct{ w, h int }{{80, 24}, {80, 24}, {100, 30}}
	i := 0
	withScreenSizeGetter(t, func() (int, int, error) {
		mu.Lock()
		defer mu.Unlock()
		s := sizes[min(i, len(sizes)-1)]
		i++
		return s.w, s.h, nil
	})

	got := make(chan [2]int, 1)
	stop := OnResize(func(width, height int) {
		got <- [2]int{width, height}
	})
	defer stop()

	select {
	case size := <-got:
		if size != [2]int{100, 30} {
			t.Fatalf("got %v, want [100 30]", size)
		}
	case <-time.After(time.Second):
		t.Fatal("timed out waiting for resize callback")
	}
}

func TestOnResizeStopStopsPolling(t *testing.T) {
	var mu sync.Mutex
	calls := 0
	toggle := false
	withScreenSizeGetter(t, func() (int, int, error) {
		mu.Lock()
		defer mu.Unlock()
		calls++
		toggle = !toggle
		if toggle {
			return 1, 1, nil
		}
		return 2, 2, nil
	})

	var fnCalls int
	var fnMu sync.Mutex
	stop := OnResize(func(width, height int) {
		fnMu.Lock()
		fnCalls++
		fnMu.Unlock()
	})

	time.Sleep(20 * time.Millisecond)
	stop()

	fnMu.Lock()
	after := fnCalls
	fnMu.Unlock()

	time.Sleep(20 * time.Millisecond)

	fnMu.Lock()
	defer fnMu.Unlock()
	if fnCalls != after {
		t.Fatalf("callback invoked after stop: before %d, after %d", after, fnCalls)
	}
}

func TestOnResizeIgnoresErrors(t *testing.T) {
	withScreenSizeGetter(t, func() (int, int, error) {
		return 0, 0, errors.New("size unavailable")
	})

	called := make(chan struct{}, 1)
	stop := OnResize(func(width, height int) {
		called <- struct{}{}
	})
	defer stop()

	select {
	case <-called:
		t.Fatal("callback should not fire when the getter errors")
	case <-time.After(50 * time.Millisecond):
	}
}
