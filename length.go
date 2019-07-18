package svg

import (
	"errors"
	"fmt"
	"strconv"
)

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

// Length corresponds to SVG <length> data type
type Length struct {
	Value float64
	Units Units
}

// String implements stringer interface for <length>
func (v *Length) String() string {
	return fmt.Sprintf("%g%s", v.Value, v.Units)
}

// Unmarshal implements unmarshals string for <length>
func (v *Length) Unmarshal(s string) (err error) {
	n := len(s)
	if n == 0 {
		return errors.New("empty value")
	}
	if s[n-1] == '%' {
		v.Units = UnitPercent
		s = s[:n-1]
	} else if n > 2 {
		switch s[n-2:] {
		case "em":
			v.Units = UnitEM
		case "ex":
			v.Units = UnitEX
		case "px":
			v.Units = UnitPX
		case "in":
			v.Units = UnitIN
		case "cm":
			v.Units = UnitCM
		case "mm":
			v.Units = UnitMM
		case "pt":
			v.Units = UnitPT
		case "pc":
			v.Units = UnitPC
		default:
			v.Units = UnitNone
		}
		if v.Units != UnitNone {
			s = s[:n-2]
		}
	}

	v.Value, err = strconv.ParseFloat(s, 64)
	return
}

// Coordinate corresponds to SVG <coordinate> type
type Coordinate Length

func (c *Coordinate) Unmarshal(s string) error {
	return (*Length)(c).Unmarshal(s)
}
