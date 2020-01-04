package svg

import (
	"errors"
	"math"
	"strconv"
	"strings"
)

type Transform struct {
	A float64
	B float64
	C float64
	D float64
	E float64
	F float64
}

func UnitTransform() Transform {
	return Transform{
		A: 1.0,
		D: 1.0,
	}
}

func Translation(tx, ty float64) Transform {
	return Transform{
		A: 1.0, D: 1.0, E: tx, F: ty,
	}
}

func Scaling(sx, sy float64) Transform {
	return Transform{
		A: sx, D: sy,
	}
}

func Rotation(a float64) Transform {
	s, c := math.Sincos(a)
	return Transform{
		A: c, B: s, C: -s, D: c,
	}
}

func SkewX(a float64) Transform {
	return Transform{
		A: 1.0, C: math.Tan(a), D: 1.0,
	}
}

func SkewY(a float64) Transform {
	return Transform{
		A: 1.0, B: math.Tan(a), D: 1.0,
	}
}

func Concatenate(a, b Transform) Transform {
	return Transform{
		A: a.A*b.A + a.C*b.B,
		B: a.B*b.A + a.D*b.B,
		C: a.A*b.C + a.C*b.D,
		D: a.B*b.C + a.D*b.D,
		E: a.A*b.E + a.C*b.F + a.E,
		F: a.B*b.E + a.D*b.F + a.F,
	}
}

func (t *Transform) TransformAbs(x, y float64) (float64, float64) {
	return t.A*x + t.C*y + t.E, t.B*x + t.D*y + t.F
}

func (t *Transform) TransformRel(x, y float64) (float64, float64) {
	return t.A*x + t.C*y, t.B*x + t.D*y
}

func (t *Transform) Calc(x, y float64, relative bool) (float64, float64) {
	if relative {
		return t.A*x + t.C*y, t.B*x + t.D*y
	}
	return t.A*x + t.C*y + t.E, t.B*x + t.D*y + t.F
}

type ViewBox struct {
	MinX   float64
	MinY   float64
	Width  float64
	Height float64
}

func (b *ViewBox) Unmarshal(s string) (err error) {
	subs := strings.Fields(strings.ReplaceAll(s, ",", " "))
	if len(subs) != 4 {
		return errors.New("invalid viewBox attribute")
	}
	b.MinX, err = strconv.ParseFloat(subs[0], 64)
	if err != nil {
		return
	}
	b.MinY, err = strconv.ParseFloat(subs[1], 64)
	if err != nil {
		return
	}
	b.Width, err = strconv.ParseFloat(subs[2], 64)
	if err != nil {
		return
	}
	b.Height, err = strconv.ParseFloat(subs[3], 64)
	return
}
