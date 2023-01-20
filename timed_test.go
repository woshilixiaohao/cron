package cron

import (
	. "github.com/smartystreets/goconvey/convey"
	"testing"
	"time"
)

func TestTimedSchedule_Next(t *testing.T) {
	Convey("Given timed parser", t, func() {
		parser := NewParser(Timed)
		spec := "timed=2021-12-24 12:30:30"
		sched, err := parser.Parse(spec)
		if err != nil {
			t.Error(err)
		}
		date := time.Date(2021, time.December, 24, 12, 30, 30, 0, time.Local)
		Convey("When next time is before timed time", func() {
			before := time.Date(2021, time.December, 24, 12, 0, 0, 0, time.Local)
			next, expired := sched.Next(before)
			Convey("Then next should equal schedule time and expired is false", func() {
				So(next, ShouldEqual, date)
				So(expired, ShouldBeFalse)
			})
			next, expired = sched.Next(before)
			Convey("Then next again should equal schedule time and expired is true", func() {
				So(next, ShouldEqual, date)
				So(expired, ShouldBeTrue)
			})
		})
		Convey("When next time is after timed time", func() {
			after := time.Date(2021, time.December, 24, 13, 0, 0, 0, time.Local)
			next, expired := sched.Next(after)
			Convey("Then next should equal schedule time and expired is true", func() {
				So(next, ShouldEqual, date)
				So(expired, ShouldBeTrue)
			})
			next, expired = sched.Next(after)
			Convey("Then next again should equal schedule time and expired is true", func() {
				So(next, ShouldEqual, date)
				So(expired, ShouldBeTrue)
			})
		})
	})
}
