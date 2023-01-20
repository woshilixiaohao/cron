package cron

import "time"

type RoutineSchedule struct {
	ExecuteTime time.Time
	Min, Max    time.Time
	Location    *time.Location
	Interval    int
	Spec        int
	RoutineType RoutineType
}

func (r *RoutineSchedule) Next(t time.Time) (time.Time, bool) {
	origLocation := t.Location()
	loc := r.Location
	if loc == time.Local {
		loc = t.Location()
	}
	if r.Location != time.Local {
		t = t.In(r.Location)
	}

	var nextTime time.Time
	switch r.RoutineType {
	case RoutineDay:
		nextTime = r.getDayTime(t, loc, 1)
	case RoutineDow:
		nextTime = r.getDayTime(t, loc, 7)
	case RoutineDom:
		nextTime = r.getMonthdayTime(t, loc)
	}
	return nextTime.In(origLocation), nextTime.After(r.Max)
}

// getDayTime 每n天
func (r *RoutineSchedule) getDayTime(t time.Time, loc *time.Location, week int) time.Time {
	nextTime := time.Date(r.Min.Year(), r.Min.Month(), r.Min.Day(), r.ExecuteTime.Hour(), r.ExecuteTime.Minute(), r.ExecuteTime.Second(), 0, loc)
	if week == 7 {
		wd := r.Spec - int(nextTime.Weekday())
		nextTime = nextTime.AddDate(0, 0, wd)
	}
	if r.Min.After(nextTime) {
		nextTime = nextTime.AddDate(0, 0, r.Interval*week)
	}
	if t.After(nextTime) {
		days := int(t.Sub(nextTime).Hours() / 24)
		days = days / r.Interval / week
		nextTime = nextTime.AddDate(0, 0, (days+1)*r.Interval*week)
	}
	return nextTime
}

// getMonthdayTime 每n月几号
func (r *RoutineSchedule) getMonthdayTime(t time.Time, loc *time.Location) time.Time {
	// last day of each month
	if r.Spec == 0 {
		return r.getLastMonthDayTime(t, loc)
	}

	return r.getNormalMonthDayTime(t, loc)
}

func (r *RoutineSchedule) getLastMonthDayTime(t time.Time, loc *time.Location) time.Time {
	firstDay := time.Date(r.Min.Year(), r.Min.Month()+1, 1, r.ExecuteTime.Hour(), r.ExecuteTime.Minute(), r.ExecuteTime.Second(), 0, loc)
	lastDayOfMin := firstDay.AddDate(0, 0, -1)
	if r.Min.After(lastDayOfMin) {
		lastDayOfMin = getLastDay(firstDay, r.Interval)
	}
	if t.After(lastDayOfMin) {
		tYear, tMonth := t.Year(), t.Month()
		nYear, nMonth := r.Min.Year(), r.Min.Month()
		monthInterval := ((tYear-nYear)*12 + int(tMonth-nMonth)) / r.Interval
		lastDayOfMin = getLastDay(firstDay, monthInterval*r.Interval)
		if t.After(lastDayOfMin) {
			lastDayOfMin = getLastDay(firstDay, r.Interval)
		}
	}
	return lastDayOfMin
}

func getLastDay(firstDay time.Time, addMonth int) time.Time {
	firstDay = firstDay.AddDate(0, addMonth, 0)
	return firstDay.AddDate(0, 0, -1)
}

func (r *RoutineSchedule) getNormalMonthDayTime(t time.Time, loc *time.Location) time.Time {
	firstTime := time.Date(r.Min.Year(), r.Min.Month(), r.Spec, r.ExecuteTime.Hour(), r.ExecuteTime.Minute(), r.ExecuteTime.Second(), 0, loc)
	for n := 1; r.Spec != firstTime.Day() && r.Max.After(firstTime); n++ {
		firstTime = time.Date(r.Min.Year(), r.Min.Month()+time.Month(r.Interval*n), r.Spec, r.ExecuteTime.Hour(), r.ExecuteTime.Minute(), r.ExecuteTime.Second(), 0, loc)
	}
	nextTime := firstTime
	if r.Min.After(firstTime) {
		nextTime = time.Date(firstTime.Year(), firstTime.Month()+time.Month(r.Interval), firstTime.Day(), firstTime.Hour(), firstTime.Minute(), firstTime.Second(), 0, loc)
		for n := 2; r.Spec != nextTime.Day() && r.Max.After(nextTime); n++ {
			nextTime = time.Date(firstTime.Year(), firstTime.Month()+time.Month(r.Interval*n), firstTime.Day(), firstTime.Hour(), firstTime.Minute(), firstTime.Second(), 0, loc)
		}
	}

	if t.After(nextTime) {
		tYear, tMonth := t.Year(), t.Month()
		nYear, nMonth := nextTime.Year(), nextTime.Month()
		monthInterval := ((tYear-nYear)*12 + int(tMonth-nMonth)) / r.Interval

		intervalNextTime := nextTime
		for n := 0; r.Spec != intervalNextTime.Day() || t.After(intervalNextTime); n++ {
			intervalNextTime = time.Date(nextTime.Year(), nextTime.Month()+time.Month((monthInterval+n)*r.Interval), nextTime.Day(), nextTime.Hour(), nextTime.Minute(), nextTime.Second(), 0, loc)
		}
		return intervalNextTime
	}
	return nextTime
}
