package svg

import "errors"

type Paint struct {
}

// FillRule implements SVG <fill-rule> type
type FillRule int

const (
	FillRuleInherit = FillRule(iota)
	FillRuleNonZero
	FillRuleEvenOdd
)

func (fr FillRule) String() string {
	switch fr {
	case FillRuleNonZero:
		return "nonzero"
	case FillRuleEvenOdd:
		return "evenodd"
	default:
		return ""
	}
}

func (fr FillRule) UnmarshalText(text []byte) error {
	s := string(text)
	switch s {
	case "nonzero":
		fr = FillRuleNonZero
	case "evenodd":
		fr = FillRuleEvenOdd
	case "inherit":
		fr = FillRuleInherit
	default:
		return errors.New("invalid fill-rule value")
	}
	return nil
}

type LineCap int

const (
	LineCapInherit = LineCap(iota)
	LineCapButt
	LineCapRound
	LineCapSquare
)

func (lc LineCap) String() string {
	switch lc {
	case LineCapInherit:
		return "inherit"
	case LineCapButt:
		return "butt"
	case LineCapRound:
		return "round"
	case LineCapSquare:
		return "square"
	default:
		return ""
	}
}

func (lc LineCap) UnmarshalText(text []byte) error {
	s := string(text)
	switch s {
	case "inherit":
		lc = LineCapInherit
	case "butt":
		lc = LineCapButt
	case "round":
		lc = LineCapRound
	case "square":
		lc = LineCapSquare
	default:
		return errors.New("invalid stroke-linecap value")
	}
	return nil
}

type LineJoin int

const (
	LineJoinInerit = LineJoin(iota)
	LineJoinMiter
	LineJoinRound
	LineJoinBevel
)

func (lj LineJoin) UnmarshalText(text []byte) error {
	s := string(text)
	switch s {
	case "inherit":
		lj = LineJoinInerit
	case "miter":
		lj = LineJoinMiter
	case "round":
		lj = LineJoinRound
	case "bevel":
		lj = LineJoinBevel
	default:
		return errors.New("invalid stroke-linejoin value")
	}
	return nil
}
