package turtlebunny

import (
	"errors"
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

func unixEpoch(s string) (float64, error) {
	if s != "subsec" && s != "subsecond" {
		return 0, errors.New("arg not one of ('subsec', 'subsecond')")
	}
	return float64(time.Now().UnixNano()) / 1e9, nil
}
