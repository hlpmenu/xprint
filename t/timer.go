package xprint

import (
	"time"

	"gopkg.hlmpn.dev/pkg/go-logger"
)

// Timer struct to ensure consistent timing
type Timer struct {
	start  time.Time
	Timing string
}

// StartTimer initializes and returns a new Timer
func StartTimer() *Timer {
	return &Timer{start: time.Now()}
}

// Elapsed returns the elapsed time since the timer started
func (t Timer) Elapsed() time.Duration {
	return time.Since(t.start)
}

func (t Timer) Stop(s string) {
	t.Timing = t.Elapsed().String()
	logger.LogSuccessf("%s Timing: %s ", s, t.Timing)
}
