package turtlebunny

import (
	"errors"

	"github.com/shopspring/decimal"
)

func toDecimal(v any) (string, error) {
	switch x := v.(type) {
	case string:
		d, err := decimal.NewFromString(x)
		return d.String(), err
	case int64:
		d := decimal.NewFromInt(x)
		return d.String(), nil
	case float64:
		d := decimal.NewFromFloat(x)
		return d.String(), nil
	default:
		return "", errors.New("unable to parse as number")
	}
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
