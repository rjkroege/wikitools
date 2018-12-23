package wiki

import (
	"os"
	"time"
)

type SystemImpl int

func (s SystemImpl) Exists(path string) bool {
	_, err := os.Stat(path)
	if err == nil {
		return true
	}
	return false
}

func (s SystemImpl) Now() time.Time {
	return time.Now()
}
