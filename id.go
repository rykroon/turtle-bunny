package turtlebunny

import (
	"fmt"

	ulid "github.com/oklog/ulid/v2"
	"lukechampine.com/uint128"
)

func ID() uint128.Uint128 {
	id := ulid.Make()
	fmt.Println(id.String())
	return uint128.FromBytesBE(id.Bytes())
}
