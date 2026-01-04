package turtlebunny

import (
	"errors"
	"fmt"
	"strconv"
	"time"

	"lukechampine.com/uint128"
)

func toUint(v any) (uint128.Uint128, error) {
	switch x := v.(type) {
	case string:
		return uint128.FromString(x)
	case int64:
		if x < 0 {
			return uint128.Zero, fmt.Errorf("cannot convert negative int64 (%d) to uint64", x)
		}
		return uint128.From64(uint64(x)), nil
	default:
		return uint128.Zero, fmt.Errorf("cannot convert %v of type %T to uint128", x, x)
	}
}

func uintAdd(x, y any) (string, error) {
	ux, err := toUint(x)
	if err != nil {
		return "", err
	}

	uy, err := toUint(y)
	if err != nil {
		return "", err
	}

	return ux.Add(uy).String(), nil
}

func uintSub(x, y any) (string, error) {
	ux, err := toUint(x)
	if err != nil {
		return "", err
	}

	uy, err := toUint(y)
	if err != nil {
		return "", err
	}

	return ux.Sub(uy).String(), nil
}

func uintCmp(x, y any) (int, error) {
	ux, err := toUint(x)
	if err != nil {
		return 0, err
	}

	uy, err := toUint(y)
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
