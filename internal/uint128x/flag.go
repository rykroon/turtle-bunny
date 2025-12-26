package uint128x

import (
	"database/sql/driver"

	"github.com/spf13/pflag"
	"lukechampine.com/uint128"
)

type uint128Flag uint128.Uint128

func newUint128Flag(val uint128.Uint128, p *uint128.Uint128) *uint128Flag {
	*p = val
	return (*uint128Flag)(p)
}

func (u uint128Flag) Value() (driver.Value, error) {
	return u.String(), nil
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
