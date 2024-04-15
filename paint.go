package svg

import (
	"errors"
)

type RGB struct {
	R uint8
	G uint8
	B uint8
}

type PaintKind int

const (
	PaintKindNone = PaintKind(iota)
	PaintKindRGB
	PaintKindGradient
)

type Paint struct {
	Kind     PaintKind
	Color    RGB
	Gradient *Gradient
}

type GradientUnits int

const (
	GradientUnitsUnspecified = GradientUnits(iota)
	GradientUnitsUserSpaceOnUse
	GradientUnitsObjectBoundingBox
)

func (gu GradientUnits) String() string {
	switch gu {
	case GradientUnitsUserSpaceOnUse:
		return "userSpaceOnUse"
	case GradientUnitsObjectBoundingBox:
		return "objectBoundingBox"
	default:
		return ""
	}
}

func (gu GradientUnits) Unmarshal(s string) error {
	switch s {
	case "userSpaceOnUse":
		gu = GradientUnitsUserSpaceOnUse
	case "objectBoundingBox":
		gu = GradientUnitsObjectBoundingBox
	default:
		return errors.New("invalid gradient-units value")
	}
	return nil
}

type SpreadMethod int

const (
	SpreadMethodUnspecified = SpreadMethod(iota)
	SpreadMethodPad
	SpreadMethodReflect
	SpreadMethodRepeat
)

func (sm SpreadMethod) String() string {
	switch sm {
	case SpreadMethodPad:
		return "pad"
	case SpreadMethodReflect:
		return "reflect"
	case SpreadMethodRepeat:
		return "repeat"
	default:
		return ""
	}
}

func (sm SpreadMethod) Unmarshal(s string) error {
	switch s {
	case "pad":
		sm = SpreadMethodPad
	case "reflect":
		sm = SpreadMethodReflect
	case "round":
		sm = SpreadMethodRepeat
	case "":
		sm = SpreadMethodUnspecified
	default:
		return errors.New("invalid spread-method value")
	}
	return nil
}

type Gradient struct {
	Units             GradientUnits
	Stops             []GradientStop
	GradientTransform Transform
	SpreadMethod      SpreadMethod
}

type GradientStop struct {
	Offset  float64
	Color   RGB
	Opacity float64
}

type LinearGradient struct {
	Gradient
	X1 Coordinate
	Y1 Coordinate
	X2 Coordinate
	Y2 Coordinate
}

type RadialGradient struct {
	Gradient
	Cx     Coordinate
	Cy     Coordinate
	Radius Length
	Fx     Coordinate
	Fy     Coordinate
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
