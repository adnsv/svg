package svg

import (
	"errors"
	"strconv"
)

// Length corresponds to SVG <length> data type, also
// can hold <percentage>, "auto", etc
type Length string

func (l Length) AsNumeric() (v float64, u Units, err error) {
	return parseLengthOrPercentage(string(l))
}

type Coordinate = Length // alias, for convenience

// Units defines length and coordinate units
type Units int

// Valid unit values
const (
	UnitNone = Units(iota)
	UnitEM
	UnitEX
	UnitPX
	UnitIN
	UnitCM
	UnitMM
	UnitPT
	UnitPC
	UnitPercent
)

// String implements unit stringer interface
func (u Units) String() string {
	switch u {
	case UnitEM:
		return "em"
	case UnitEX:
		return "ex"
	case UnitPX:
		return "px"
	case UnitIN:
		return "in"
	case UnitCM:
		return "cm"
	case UnitMM:
		return "mm"
	case UnitPT:
		return "pt"
	case UnitPC:
		return "pc"
	case UnitPercent:
		return "%"
	default:
		return ""
	}
}

var ErrEmptyValue = errors.New("empty value")

func parseLengthOrPercentage(s string) (v float64, u Units, e error) {
	n := len(s)
	if n == 0 {
		e = errors.New("missing or empty value")
		return
	}
	if s[n-1] == '%' {
		u = UnitPercent
		s = s[:n-1]
	} else if n > 2 {
		switch s[n-2:] {
		case "em":
			u = UnitEM
		case "ex":
			u = UnitEX
		case "px":
			u = UnitPX
		case "in":
			u = UnitIN
		case "cm":
			u = UnitCM
		case "mm":
			u = UnitMM
		case "pt":
			u = UnitPT
		case "pc":
			u = UnitPC
		default:
			u = UnitNone
		}
		if u != UnitNone {
			s = s[:n-2]
		}
	}

	v, e = strconv.ParseFloat(s, 64)
	return
}
