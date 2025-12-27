package turtlebunny

import (
	"strconv"
	"time"

	"lukechampine.com/uint128"
)

func isUint128(s string) (bool, error) {
	_, err := uint128.FromString(s)
	if err != nil {
		return false, err
	}

	return true, nil
}

func isUint64(s string) (bool, error) {
	_, err := strconv.ParseUint(s, 10, 64)
	if err != nil {
		return false, err
	}

	return true, nil
}

func unixNano() string {
	return strconv.Itoa(int(time.Now().UnixNano()))
}

func decimalAdd(x, y string) (string, error) {
	a, err := uint128.FromString(x)
	if err != nil {
		return "", err
	}

	b, err := uint128.FromString(y)
	if err != nil {
		return "", err
	}

	return a.Add(b).String(), nil
}

func decimalSub(x, y string) (string, error) {
	a, err := uint128.FromString(x)
	if err != nil {
		return "", err
	}

	b, err := uint128.FromString(y)
	if err != nil {
		return "", err
	}

	return a.Sub(b).String(), nil
}
