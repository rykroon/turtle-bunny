package turtlebunny

import (
	"lukechampine.com/uint128"
)

type Uint128 = uint128.Uint128

func NewUint128FromString(s string) (Uint128, error) {
	return uint128.FromString(s)
}
