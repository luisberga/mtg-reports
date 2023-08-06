package timer

import "time"

type timer struct{}

func New() *timer {
	return &timer{}
}

func (t *timer) Now() string {
	return time.Now().UTC().Format("2006-01-02 15:04:05")
}
