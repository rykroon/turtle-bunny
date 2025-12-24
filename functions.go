package turtlebunny

import (
	"time"

	"github.com/shopspring/decimal"
)

var maxUint128 decimal.Decimal
var maxUint64 decimal.Decimal

func init() {
	maxUint128, _ = decimal.NewFromString("340282366920938463463374607431768211455")
	maxUint64, _ = decimal.NewFromString("18446744073709551615")
}

func decimalAdd(x, y string) (string, error) {
	dx, err := decimal.NewFromString(x)
	if err != nil {
		return "", err
	}

	dy, err := decimal.NewFromString(y)
	if err != nil {
		return "", err
	}

	return dx.Add(dy).String(), nil
}

func decimalSub(x, y string) (string, error) {
	dx, err := decimal.NewFromString(x)
	if err != nil {
		return "", err
	}

	dy, err := decimal.NewFromString(y)
	if err != nil {
		return "", err
	}

	return dx.Sub(dy).String(), nil
}

func isUint128(s string) (bool, error) {
	d, err := decimal.NewFromString(s)
	if err != nil {
		return false, err
	}

	return d.GreaterThanOrEqual(decimal.Zero) && d.LessThanOrEqual(maxUint128), nil
}

func isUint64(s string) (bool, error) {
	d, err := decimal.NewFromString(s)
	if err != nil {
		return false, err
	}

	return d.GreaterThanOrEqual(decimal.Zero) && d.LessThanOrEqual(maxUint64), nil
}

func unixNano() int64 {
	return time.Now().UnixNano()
}
