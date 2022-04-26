package testUtils

import (
	"app/clock"
	"time"
)

// SetFakeTime set fake time to clock package.
func SetFakeTime(t time.Time) {
	clock.Now = func() time.Time {
		return t
	}
}
