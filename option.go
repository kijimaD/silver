package silver

import (
	"io"
	"time"
)

type TaskOption func(*Task)

func TaskWithWriter(w io.Writer) TaskOption {
	return func(t *Task) {
		t.w = w
	}
}

type StatsOption func(*Stats)

func StatsWithIdx(currentIdx int, allLen int) func(*Stats) {
	return func(s *Stats) {
		s.CurrentIdx = currentIdx
		s.AllLen = allLen
	}
}

func StatsWithStart() func(*Stats) {
	return func(s *Stats) {
		s.StartedAt = time.Now()
	}
}
