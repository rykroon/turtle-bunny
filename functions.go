package turtlebunny

import (
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
