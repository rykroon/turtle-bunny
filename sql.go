package turtlebunny

import (
	"errors"

	"lukechampine.com/uint128"
)

type GenericScanner[T any] struct {
	ptr      *T
	scanFunc func(any, *T) error
}

func (gs *GenericScanner[T]) Scan(src any) error {
	return gs.scanFunc(src, gs.ptr)
}

func NewScannableUint128(p *uint128.Uint128) *GenericScanner[uint128.Uint128] {
	return &GenericScanner[uint128.Uint128]{
		ptr: p,
		scanFunc: func(src any, p *uint128.Uint128) error {
			srcString, ok := src.(string)
			if !ok {
				return errors.New("src value is not a string")
			}
			v, err := uint128.FromString(srcString)
			if err != nil {
				return err
			}
			*p = v
			return nil
		},
	}
}
