package turtlebunny

import (
	"time"

	"github.com/shopspring/decimal"
)

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

func decimalBetween(x, y, z string) (bool, error) {
	dx, err := decimal.NewFromString(x)
	if err != nil {
		return false, err
	}

	dy, err := decimal.NewFromString(y)
	if err != nil {
		return false, err
	}

	dz, err := decimal.NewFromString(z)
	if err != nil {
		return false, err
	}

	return dx.GreaterThanOrEqual(dy) && dx.LessThanOrEqual(dz), nil

}

func isUint128(s string) (bool, error) {
	d, err := decimal.NewFromString(s)
	if err != nil {
		return false, err
	}

	maxUint128, err := decimal.NewFromString("340282366920938463463374607431768211455")
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

	maxUint128, err := decimal.NewFromString("18446744073709551615")
	if err != nil {
		return false, err
	}
	return d.GreaterThanOrEqual(decimal.Zero) && d.LessThanOrEqual(maxUint128), nil
}

func unixNano() int64 {
	return time.Now().UnixNano()
}

func decimalCmp(x, y string) (int8, error) {
	dx, err := decimal.NewFromString(x)
	if err != nil {
		return 0, err
	}

	dy, err := decimal.NewFromString(y)
	if err != nil {
		return 0, err
	}

	if dx.LessThan(dy) {
		return -1, nil
	} else if dx.Equal(dy) {
		return 0, nil
	} else {
		return 1, nil
	}
}
