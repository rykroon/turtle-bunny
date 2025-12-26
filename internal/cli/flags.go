package cli

import (
	"fmt"
	"strings"

	"github.com/spf13/pflag"
	"lukechampine.com/uint128"
)

type uint128Flag uint128.Uint128

func newUint128Flag(val uint128.Uint128, p *uint128.Uint128) *uint128Flag {
	*p = val
	return (*uint128Flag)(p)
}

func (u uint128Flag) String() string {
	return uint128.Uint128(u).String()
}

func (uf *uint128Flag) Set(s string) error {
	u, err := uint128.FromString(s)
	if err != nil {
		return err
	}
	*uf = uint128Flag(u)
	return nil
}

func (u uint128Flag) Type() string {
	return "uint128"
}

func Uint128VarP(f *pflag.FlagSet, p *uint128.Uint128, name, shorthand string, value uint128.Uint128, usage string) {
	f.VarP(newUint128Flag(value, p), name, shorthand, usage)
}

func Uint128Var(f *pflag.FlagSet, p *uint128.Uint128, name string, value uint128.Uint128, usage string) {
	Uint128VarP(f, p, name, "", value, usage)
}

func Uint128P(f *pflag.FlagSet, name, shorthand string, value uint128.Uint128, usage string) {
	p := new(uint128.Uint128)
	Uint128VarP(f, p, name, shorthand, value, usage)
}

func Uint128(f *pflag.FlagSet, name string, value uint128.Uint128, usage string) {
	p := new(uint128.Uint128)
	Uint128VarP(f, p, name, "", value, usage)
}

// ~~~~~~~~~~~~~~~~

type uint128SliceValue struct {
	value   *[]uint128.Uint128
	changed bool
}

func newUint128SliceValue(val []uint128.Uint128, p *[]uint128.Uint128) *uint128SliceValue {
	isv := new(uint128SliceValue)
	isv.value = p
	*isv.value = val
	return isv
}

func (s *uint128SliceValue) Set(val string) error {
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

func (s *uint128SliceValue) Type() string {
	return "uint128Slice"
}

func (s *uint128SliceValue) String() string {
	out := make([]string, len(*s.value))
	for i, d := range *s.value {
		out[i] = fmt.Sprintf("%d", d)
	}
	return "[" + strings.Join(out, ",") + "]"
}

// UintSliceVar defines a uintSlice flag with specified name, default value, and usage string.
// The argument p points to a []uint variable in which to store the value of the flag.
func Uint128SliceVar(f *pflag.FlagSet, p *[]uint128.Uint128, name string, value []uint128.Uint128, usage string) {
	f.VarP(newUint128SliceValue(value, p), name, "", usage)
}

// UintSliceVarP is like UintSliceVar, but accepts a shorthand letter that can be used after a single dash.
func Uint128SliceVarP(f *pflag.FlagSet, p *[]uint128.Uint128, name, shorthand string, value []uint128.Uint128, usage string) {
	f.VarP(newUint128SliceValue(value, p), name, shorthand, usage)
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
