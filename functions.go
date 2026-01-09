package turtlebunny

import (
	"errors"
	"fmt"
	"regexp"
	"strconv"
	"time"

	"lukechampine.com/uint128"
)

var uintPattern *regexp.Regexp = regexp.MustCompile("^(0|[1-9][0-9]*)$")

func toUint(v any) (uint128.Uint128, error) {
	switch x := v.(type) {
	case string:
		if !uintPattern.MatchString(x) {
			return uint128.Zero, errors.New("not a valid unisgned integer")
		}
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
	if !uintPattern.MatchString(s) {
		return false
	}
	_, err := uint128.FromString(s)
	return err == nil
}

func isUint64(s string) bool {
	if !uintPattern.MatchString(s) {
		return false
	}
	_, err := strconv.ParseUint(s, 10, 64)
	return err == nil
}

func getUint128Max() string {
	return uint128.Max.String()
}

func unixNano() uint64 {
	return uint64(time.Now().UnixNano())
}
