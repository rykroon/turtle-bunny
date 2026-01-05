package timevalue

import (
	"errors"
	"math"
	"slices"
	"time"
)

// unixepoch
// julianday
// auto
// must immediately follow the initial time-value which must be of the form DDDDDDDDD
type NumberModifier func(float64) (TimeValue, error)

func UnixEpochMod(f float64) (TimeValue, error) {
	sec := int64(math.Floor(f))
	frac := f - float64(sec)
	nsec := int64(math.Round(frac * 1e9))
	return TimeValue(time.Unix(sec, nsec)), nil
}

func JulianDayMod(f float64) (TimeValue, error) {
	return UnixEpochMod(julianDaystoUnixSecs(f))
}

func julianDaystoUnixSecs(julianDays float64) float64 {
	return (julianDays - 2440587.5) * 86400
}

func unixSecsToJulianDay(unixSecs float64) float64 {
	return unixSecs/86400 + 2440587.5
}

func AutoMod(f float64) (TimeValue, error) {
	if 0 <= f && f <= 5373484.499999 {
		return JulianDayMod(f)
	} else if -210866760000 <= f && f <= 253402300799 {
		return UnixEpochMod(f)
	} else {
		return zero, errors.New("invalid time-value range")
	}
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
