package cli

import (
	"fmt"
	"strings"

	"github.com/spf13/pflag"
	"lukechampine.com/uint128"
)

type GenericFlag[T fmt.Stringer] struct {
	ptr     *T
	type_   string
	setFunc func(s string, p *T) error
}

func (gf GenericFlag[T]) String() string {
	if gf.ptr == nil {
		return "nil pointer"
	}
	v := *gf.ptr
	return v.String()
}

func (gf *GenericFlag[T]) Set(s string) error {
	return gf.setFunc(s, gf.ptr)
}

func (gf GenericFlag[T]) Type() string {
	return gf.type_
}

func NewUint128Flag(p *uint128.Uint128) *GenericFlag[uint128.Uint128] {
	return &GenericFlag[uint128.Uint128]{
		ptr:   p,
		type_: "uint128",
		setFunc: func(s string, p *uint128.Uint128) error {
			v, err := uint128.FromString(s)
			if err != nil {
				return err
			}
			*p = v
			return nil
		},
	}
}

// ~~~~~~~~~~~~~~~~

type uint128SliceFlag struct {
	value   *[]uint128.Uint128
	changed bool
}

func newUint128SliceFlag(val []uint128.Uint128, p *[]uint128.Uint128) *uint128SliceFlag {
	isv := new(uint128SliceFlag)
	isv.value = p
	*isv.value = val
	return isv
}

func (s *uint128SliceFlag) Set(val string) error {
	ss := strings.Split(val, ",")
	out := make([]uint128.Uint128, len(ss))
	for i, d := range ss {
		var err error
		out[i], err = uint128.FromString(d)
		if err != nil {
			return err
		}
	}
	if !s.changed {
		*s.value = out
	} else {
		*s.value = append(*s.value, out...)
	}
	s.changed = true
	return nil
}

func (s *uint128SliceFlag) Type() string {
	return "uint128Slice"
}

func (s *uint128SliceFlag) String() string {
	out := make([]string, len(*s.value))
	for i, d := range *s.value {
		out[i] = fmt.Sprintf("%d", d)
	}
	return "[" + strings.Join(out, ",") + "]"
}

// UintSliceVar defines a uintSlice flag with specified name, default value, and usage string.
// The argument p points to a []uint variable in which to store the value of the flag.
func Uint128SliceVar(f *pflag.FlagSet, p *[]uint128.Uint128, name string, value []uint128.Uint128, usage string) {
	f.VarP(newUint128SliceFlag(value, p), name, "", usage)
}

// UintSliceVarP is like UintSliceVar, but accepts a shorthand letter that can be used after a single dash.
func Uint128SliceVarP(f *pflag.FlagSet, p *[]uint128.Uint128, name, shorthand string, value []uint128.Uint128, usage string) {
	f.VarP(newUint128SliceFlag(value, p), name, shorthand, usage)
}

// UintSlice defines a []uint flag with specified name, default value, and usage string.
// The return value is the address of a []uint variable that stores the value of the flag.
func Uint128Slice(f *pflag.FlagSet, name string, value []uint128.Uint128, usage string) *[]uint128.Uint128 {
	p := []uint128.Uint128{}
	Uint128SliceVarP(f, &p, name, "", value, usage)
	return &p
}

// UintSliceP is like UintSlice, but accepts a shorthand letter that can be used after a single dash.
func Uint128SliceP(f *pflag.FlagSet, name, shorthand string, value []uint128.Uint128, usage string) *[]uint128.Uint128 {
	p := []uint128.Uint128{}
	Uint128SliceVarP(f, &p, name, shorthand, value, usage)
	return &p
}
