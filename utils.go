package main

import (
	"time"
)

func Sleep(duration time.Duration) {
	if duration < 0 {
		return
	}
	time.Sleep(duration)
}
