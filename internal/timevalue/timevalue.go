package timevalue

import (
	"errors"
	"fmt"
	"strconv"
	"time"
)

// https://sqlite.org/lang_datefunc.html

type TimeValue time.Time

var zero = TimeValue{}

func Now() TimeValue {
	return TimeValue(time.Now())
}

func NewFromString(s string, mods ...string) (TimeValue, error) {
	if s == "now" {
		return Now(), nil
	}

	dateLayouts := []string{
		"2006-01-02",
		"2006-01-02 15:04",
		"2006-01-02 15:04:05",
		"2006-01-02 15:04:05.999999999",
		"2006-01-02T15:04",
		"2006-01-02T15:04:05",
		"2006-01-02T15:04:05.999999999",
	}

	for _, layout := range dateLayouts {
		t, err := time.Parse(layout, s)
		if err == nil {
			return TimeValue(t), nil
		}
	}

	timeLayouts := []string{
		"15:04",
		"15:04:05",
		"15:04:05.999999999",
	}

	for _, layout := range timeLayouts {
		t, err := time.Parse(layout, s)
		if err == nil {
			// https://sqlite.org/lang_datefunc.html#tmval
			// Formats 8 through 10 that specify only a time assume a date of 2000-01-01
			y2k := time.Date(2000, 1, 1, t.Hour(), t.Minute(), t.Second(), t.Nanosecond(), t.Location())
			return TimeValue(y2k), nil
		}
	}

	f, err := strconv.ParseFloat(s, 64)
	if err != nil {
		return zero, err
	}

	tv, err := NewFromFloat(f, mods...)
	if err != nil {
		return zero, err
	}

	return tv, nil
}

func NewFromFloat(f float64, mods ...string) (TimeValue, error) {
	modifiers := Modifiers(mods)
	auto := modifiers.Auto()
	unixepoch := modifiers.UnixEpoch()

	if auto && unixepoch {
		return zero, errors.New("cannot specify both 'auto' and 'unixepoch'")
	}

	julianday := !auto && !unixepoch

	var numMod NumberModifier

	if julianday {
		numMod = JulianDayMod
	} else if auto {
		numMod = AutoMod
	} else if unixepoch {
		numMod = UnixEpochMod
	} else {
		return zero, errors.New("")
	}

	return numMod(f)
}

func (tv TimeValue) Date() string {
	return time.Time(tv).Format(time.DateOnly)
}

func (tv TimeValue) Time(subsec bool) string {
	s := time.Time(tv).Format(time.TimeOnly)
	if !subsec {
		return s
	}
	ns := time.Time(tv).Nanosecond()
	return fmt.Sprintf("%s.%d", s, ns)
}

func (tv TimeValue) DateTime(subsec bool) string {
	s := time.Time(tv).Format(time.DateTime)
	if !subsec {
		return s
	}
	ns := time.Time(tv).Nanosecond()
	return fmt.Sprintf("%s.%d", s, ns)
}

func (tv TimeValue) UnixEpoch(subsec bool) any {
	if !subsec {
		return time.Time(tv).Unix()
	}
	return float64(time.Time(tv).UnixNano()) / 1_000_000_000
}
