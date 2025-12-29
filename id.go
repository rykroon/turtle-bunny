package turtlebunny

import (
	ulid "github.com/oklog/ulid/v2"
	"lukechampine.com/uint128"
)

func ID() uint128.Uint128 {
	return uint128.FromBytesBE(ulid.Make().Bytes())
}
