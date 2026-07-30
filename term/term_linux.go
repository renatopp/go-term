//go:build linux

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
			syscall.SIGWINCH,
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
				case syscall.SIGWINCH:
					width, height, err := GetScreenSize()
					if err == nil {
						bus.Publish(EventResize, ResizeEvent{Width: width, Height: height})
					}
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
