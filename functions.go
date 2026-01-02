package turtlebunny

import (
	"errors"
	"strconv"
	"time"

	"lukechampine.com/uint128"
)

func uintAdd(x, y string) (string, error) {
	ux, err := uint128.FromString(x)
	if err != nil {
		return "", err
	}

	uy, err := uint128.FromString(y)
	if err != nil {
		return "", err
	}

	return ux.Add(uy).String(), nil
}

func uintSub(x, y string) (string, error) {
	ux, err := uint128.FromString(x)
	if err != nil {
		return "", err
	}

	uy, err := uint128.FromString(y)
	if err != nil {
		return "", err
	}

	return ux.Sub(uy).String(), nil
}

func uintCmp(x, y string) (int, error) {
	ux, err := uint128.FromString(x)
	if err != nil {
		return 0, err
	}

	uy, err := uint128.FromString(y)
	if err != nil {
		return 0, err
	}

	return ux.Cmp(uy), nil
}

func isUint128(s string) bool {
	_, err := uint128.FromString(s)
	return err == nil
}

func isUint64(s string) bool {
	_, err := strconv.ParseUint(s, 10, 64)
	return err == nil
}

func unixNano() string {
	return strconv.Itoa(int(time.Now().UnixNano()))
}

func unixEpoch(s string) (float64, error) {
	if s != "subsec" && s != "subsecond" {
		return 0, errors.New("arg not one of ('subsec', 'subsecond')")
	}
	return float64(time.Now().UnixNano()) / 1e9, nil
}
