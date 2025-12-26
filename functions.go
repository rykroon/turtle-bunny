package turtlebunny

import (
	"errors"
	"strconv"
	"time"

	"github.com/shopspring/decimal"
)

func isUint128(s string) (bool, error) {
	d, err := decimal.NewFromString(s)
	if err != nil {
		return false, err
	}

	if !d.IsInteger() {
		return false, errors.New("not an integer")
	}

	if d.IsNegative() {
		return false, errors.New("not an unsigned integer")
	}

	if d.BigInt().BitLen() > 128 {
		return false, errors.New("not an unsigned 128 bit integer")
	}

	return true, nil
}

func isUint64(s string) (bool, error) {
	d, err := decimal.NewFromString(s)
	if err != nil {
		return false, err
	}

	if !d.IsInteger() {
		return false, errors.New("not an integer")
	}

	if d.IsNegative() {
		return false, errors.New("not an unsigned integer")
	}

	if d.BigInt().BitLen() > 64 {
		return false, errors.New("not an unsigned 64 bit integer")
	}

	return true, nil
}

func unixNano() string {
	return strconv.Itoa(int(time.Now().UnixNano()))
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
