package common

import "time"

func NowIfZero(t time.Time) time.Time {
	if t.IsZero() {
		return time.Now()
	}

	return t
}
