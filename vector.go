package svg

import "math"

type Vector struct {
	X float64
	Y float64
}

func (v Vector) Norm() float64 {
	return v.X*v.X + v.Y*v.Y
}

func (v Vector) Length() float64 {
	return math.Sqrt(v.Norm())
}

func (v Vector) Normalized() Vector {
	return Mul(v, 1/v.Length())
}

func Add(vv ...Vector) Vector {
	ret := Vector{}
	for _, v := range vv {
		ret.X += v.X
		ret.Y += v.Y
	}
	return ret
}

func Sub(a, b Vector) Vector {
	return Vector{a.X - b.X, a.Y - b.Y}
}

func Mul(a Vector, b float64) Vector {
	return Vector{a.X * b, a.Y * b}
}

func Div(a Vector, b float64) Vector {
	return Vector{a.X / b, a.Y / b}
}

func Dot(a, b Vector) float64 {
	return a.X*b.X + a.Y*b.Y
}

func Cross(a, b Vector) float64 {
	return a.X*b.Y - a.Y*b.X
}
