package uint128x

import (
	"errors"

	"lukechampine.com/uint128"
)

type ScannableUint128 struct {
	*uint128.Uint128
}

func NewScannableUint128(u *uint128.Uint128) *ScannableUint128 {
	return &ScannableUint128{u}
}

func (s *ScannableUint128) Scan(src any) error {
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
