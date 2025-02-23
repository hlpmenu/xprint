package xprint

import (
	"time"

	"gopkg.hlmpn.dev/pkg/go-logger"
)

// Timer struct to ensure consistent timing
type Timer struct {
	start  time.Time
	end    time.Time
	timing string
}

// StartTimer initializes and returns a new Timer
func StartTimer() *Timer {
	return &Timer{start: time.Now()}
}
func NewTimer() *Timer {
	return &Timer{}
}

// Elapsed returns the elapsed time since the timer started
func (t *Timer) elapsed() time.Duration {
	if t.end.IsZero() {
		logger.LogErrorf("Cant get elapsed time before timer is stopped")
	}
	return t.end.Sub(t.start)
}

func (t *Timer) Start() {
	t.start = time.Now()
}

func (t *Timer) Stop() {
	t.end = time.Now()
	t.timing = t.elapsed().String()
}

func (t *Timer) StopGetDuration() time.Duration {
	t.Stop()
	return t.elapsed()
}

func (t *Timer) Duration() time.Duration {
	return t.elapsed()
}

func (t *Timer) StopAndLog(s string) {
	t.Stop()
	logger.LogSuccessf("%s Timing: %s ", s, t.timing)
}
