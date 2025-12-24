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

func decimalCmp(x, y string) (int, error) {
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

func unixNano() int64 {
	return time.Now().UnixNano()
}
