package turtlebunny

import (
	"time"

	"github.com/shopspring/decimal"
)

func toDecimal(v any) decimal.Decimal {
	switch x := v.(type) {
	case int64:
		return decimal.NewFromInt(x)
	case float64:
		return decimal.NewFromFloat(x)
	case string:
		d, err := decimal.NewFromString(x)
		if err != nil {
			return decimal.Zero
		}
		return d
	default:
		return decimal.Zero
	}
}

func decimalAdd(x, y any) string {
	dx := toDecimal(x)
	dy := toDecimal(y)
	return dx.Add(dy).String()
}

func decimalSub(x, y any) string {
	dx := toDecimal(x)
	dy := toDecimal(y)
	return dx.Sub(dy).String()
}

func decimalMul(x, y any) string {
	dx := toDecimal(x)
	dy := toDecimal(y)
	return dx.Mul(dy).String()
}

func decimalCmp(x, y any) int {
	dx := toDecimal(x)
	dy := toDecimal(y)
	return dx.Cmp(dy)
}

/*
This is a tough problem that is probably not worth my time, but I feel addicted

first thing to process is the time value.
The time value can be many things:
- It can be a datetime string, which will then require date time parsing.
- It can be an int/float, which can be interpreted as a unix timestamp or Julian Day
- It can be the subsecond modifier.
*/

func julianDaystoUnixSecs(julianDays decimal.Decimal) decimal.Decimal {
	return julianDays.Sub(decimal.NewFromFloat(2440587.5)).Mul(decimal.NewFromInt(86400))
}

func unixEpoch(timeValue any, modifiers ...string) any {
	subsec := false
	auto := false
	unixepoch := false

	for _, mod := range modifiers {
		switch mod {
		case "subsec", "subsecond":
			subsec = true
		case "auto":
			auto = true
		case "unixepoch":
			unixepoch = true
		default:
			return nil
		}
	}

	t := time.Now()

	switch tv := timeValue.(type) {
	case string:
		switch tv {
		case "subsec", "subsecond":
			subsec = true
		case "now":
			t = time.Now()
		default:
			var err error
			// check multiple formats
			formats := []string{
				"2006-01-02",
				"2006-01-02 15:04",
				"2006-01-02 15:04:05",
				"2006-01-02 15:04:05.000",
				"2006-01-02T15:04",
				"2006-01-02T15:04:05",
				"2006-01-02T15:04:05.000",
				"15:04",
				"15:04:05",
				"15:04:05.000",
				// "DDDDDDDDDD",
			}
			success := false
			for _, format := range formats {
				t, err = time.Parse(format, tv)
				if err == nil {
					success = true
					break
				}
			}

			if !success {
				return nil
			}
		}

	case int64, float64:
		d := decimal.Zero
		if i64, ok := tv.(int64); ok {
			d = decimal.NewFromInt(i64)
		}

		if f64, ok := tv.(float64); ok {
			d = decimal.NewFromFloat(f64)
		}

		if !auto && !unixepoch {
			d = julianDaystoUnixSecs(d)
		}

		t = time.Unix(d.IntPart(), 0)

	default:
		return nil
	}

	if subsec {
		return float64(t.UnixNano()) / 1000000000
	} else {
		return t.Unix()
	}
}
