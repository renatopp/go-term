//go:build windows

package term

import (
	"os"
	"os/signal"
	"syscall"
)

func startSignalListener(sig ...os.Signal) (stop func()) {
	if len(sig) == 0 {
		sig = []os.Signal{
			os.Interrupt,
			syscall.SIGTERM,
		}
	}
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, sig...)

	done := make(chan struct{})
	go func() {
		for {
			select {
			case s := <-sigs:
				switch s {
				default:
					bus.Publish(EventSignal, SignalEvent{Signal: s})
				}
			case <-done:
				signal.Stop(sigs)
				return
			}
		}
	}()
	return func() { close(done) }
}
