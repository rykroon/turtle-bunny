package turtlebunny

import (
	"errors"

	"lukechampine.com/uint128"
)

type scannableUint128 struct {
	*uint128.Uint128
}

func newScannableUint128(u *uint128.Uint128) *scannableUint128 {
	return &scannableUint128{u}
}

func (s *scannableUint128) Scan(src any) error {
	srcString, ok := src.(string)
	if !ok {
		return errors.New("src value is not a string")
	}
	u, err := uint128.FromString(srcString)
	if err != nil {
		return err
	}
	*s.Uint128 = u
	return nil
}
