package unixepoch

import (
	"errors"
	"fmt"
	"math"
	"slices"
	"strconv"
	"time"
)

// https://sqlite.org/lang_datefunc.html

type TimeValue time.Time

var zero = TimeValue{}

func Now() TimeValue {
	return TimeValue(time.Now())
}

func NewFromString(s string, mods Modifiers) (TimeValue, error) {
	if s == "now" {
		tv := TimeValue(time.Now())
		return tv, nil
	}

	layouts := []string{
		"2006-01-02",
		"2006-01-02 15:04",
		"2006-01-02 15:04:05",
		"2006-01-02 15:04:05.999999999",
		"2006-01-02T15:04",
		"2006-01-02T15:04:05",
		"2006-01-02T15:04:05.999999999",
		// "15:04"
		// "15:04:05"
		// "15:04:05.999999999"
	}

	for _, layout := range layouts {
		tv, err := time.Parse(layout, s)
		if err == nil {
			return TimeValue(tv), nil
		}
	}

	f, err := strconv.ParseFloat(s, 64)
	if err != nil {
		return zero, err
	}

	tv, err := NewFromFloat(f, mods)
	if err != nil {
		return zero, err
	}

	return tv, nil
}

func NewFromFloat(f float64, mods Modifiers) (TimeValue, error) {
	auto := mods.Auto()
	unixepoch := mods.UnixEpoch()

	if auto && unixepoch {
		return zero, errors.New("cannot specify both 'auto' and 'unixepoch'")
	}

	julianday := !auto && !unixepoch

	if auto {
		if 0 <= f && f <= 5373484.499999 {
			julianday = true
		} else if -210866760000 <= f && f <= 253402300799 {
			unixepoch = true
		}
	}

	if !julianday && !unixepoch {
		return zero, errors.New("invalid range")
	}

	if julianday {
		f = julianDaystoUnixSecs(f)
	}

	sec := int64(math.Floor(f))
	frac := f - float64(sec)
	nsec := int64(math.Round(frac * 1e9))
	return TimeValue(time.Unix(sec, nsec)), nil
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

func julianDaystoUnixSecs(julianDays float64) float64 {
	return (julianDays - 2440587.5) * 86400
}

func unixSecsToJulianDay(unixSecs float64) float64 {
	return unixSecs/86400 + 2440587.5
}

func resolveArgs(args []any) (any, []string, error) {
	if len(args) == 0 {
		return "now", []string{}, nil
	}

	timeValue := args[0]
	modifiers, err := anySliceToStringSlice(args[1:])
	if err != nil {
		return "", []string{}, err
	}

	s, ok := timeValue.(string)
	if ok && (s == "subsec" || s == "subsecond") {
		// if the first arg is subsec then append to to modifiers
		// and set the time value to now
		timeValue = "now"
		modifiers = append(modifiers, s)
	}

	return timeValue, modifiers, nil
}

func UnixEpoch(args ...any) (any, error) {
	tv, mods, err := resolveArgs(args)
	if err != nil {
		return nil, err
	}

	var timeValue TimeValue
	modifiers := Modifiers(mods)
	subsec := modifiers.SubSec()

	switch tvTyped := tv.(type) {
	case string:
		timeValue, err = NewFromString(tvTyped, modifiers)
		if err != nil {
			return nil, err
		}
	case int64:
		timeValue, err = NewFromFloat(float64(tvTyped), modifiers)
	case float64:
		timeValue, err = NewFromFloat(tvTyped, modifiers)
	}

	return timeValue.UnixEpoch(subsec), nil
}

type Modifiers []string

func (m Modifiers) SubSec() bool {
	return slices.ContainsFunc(m, func(mod string) bool {
		return mod == "subsec" || mod == "subsecond"
	})
}

func (m Modifiers) Auto() bool {
	return slices.Contains(m, "auto")
}

func (m Modifiers) UnixEpoch() bool {
	return slices.Contains(m, "unixepoch")
}

func (m Modifiers) JulianDay() bool {
	return slices.Contains(m, "julianday")
}

func anySliceToStringSlice(a []any) ([]string, error) {
	s := make([]string, len(a))
	for i, v := range a {
		str, ok := v.(string)
		if !ok {
			return nil, fmt.Errorf("element %d is not a string (type %T)", i, v)
		}
		s[i] = str
	}
	return s, nil
}
