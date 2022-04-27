package testutils

import (
	"app/clock"
	"time"
)

func SetFakeTime(t time.Time) {
	clock.Now = func() time.Time {
		return t
	}
}
