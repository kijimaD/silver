package silver

import "io"

type TaskOption func(*Task)

func WithWriter(w io.Writer) TaskOption {
	return func(t *Task) {
		t.w = w
	}
}
