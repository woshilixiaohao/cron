package cron

import "time"

type TimedSchedule struct {
	Timing     time.Time
	isExecuted bool
}

func (d *TimedSchedule) Next(t time.Time) (time.Time, bool) {
	if d.Timing.Before(t) || d.isExecuted {
		return d.Timing, true
	}
	d.isExecuted = true
	return d.Timing, false
}
