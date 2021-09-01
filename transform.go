package svg

import (
	"errors"
	"fmt"
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

func UnitTransform() *Transform {
	return &Transform{
		A: 1.0,
		D: 1.0,
	}
}

func Translation(tx, ty float64) *Transform {
	return &Transform{
		A: 1.0, D: 1.0, E: tx, F: ty,
	}
}

func Scaling(sx, sy float64) *Transform {
	return &Transform{
		A: sx, D: sy,
	}
}

func Rotation(a float64) *Transform {
	s, c := math.Sincos(a)
	return &Transform{
		A: c, B: s, C: -s, D: c,
	}
}

func SkewX(a float64) *Transform {
	return &Transform{
		A: 1.0, C: math.Tan(a), D: 1.0,
	}
}

func SkewY(a float64) *Transform {
	return &Transform{
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

func (t *Transform) CalcAbs(x, y float64) (float64, float64) {
	return t.A*x + t.C*y + t.E, t.B*x + t.D*y + t.F
}

func (t *Transform) CalcRel(x, y float64) (float64, float64) {
	return t.A*x + t.C*y, t.B*x + t.D*y
}

func (t *Transform) Calc(x, y float64, relative bool) (float64, float64) {
	if relative {
		return t.A*x + t.C*y, t.B*x + t.D*y
	}
	return t.A*x + t.C*y + t.E, t.B*x + t.D*y + t.F
}

func (t *Transform) Unmarshal(s string) (err error) {

	skipWSP := func() {
		for len(s) > 0 && (s[0] == ' ' || s[0] == '\t' || s[0] == 'r' || s[0] == '\f') {
			s = s[1:]
		}
	}

	getIdent := func() string {
		skipWSP()
		i := 0
		for i < len(s) && ((s[i] >= 'a' && s[i] <= 'z') || (s[i] >= 'A' && s[i] <= 'Z')) {
			i++
		}
		ret := s[:i]
		s = s[i:]
		skipWSP()
		return ret
	}

	tts := []*Transform{}

	skipWSP()
	for len(s) > 0 {
		cmd := getIdent()
		if cmd == "" {
			return errors.New("invalid transform command")
		}

		if len(s) > 0 && s[0] == '(' {
			s = s[1:]
			skipWSP()
		} else {
			return errors.New("missing opening parenthesis")
		}
		cp := strings.Index(s, ")")
		if cp < 0 {
			return errors.New("missing closing parenthesis")
		}
		argstr := s[:cp]
		s = s[cp+1:]
		argstr = strings.ReplaceAll(argstr, ",", " ")
		argstrs := strings.Split(argstr, " ")
		args := []float64{}
		for _, arg := range argstrs {
			if arg != "" {
				v, err := strconv.ParseFloat(arg, 64)
				if err != nil {
					return err
				}
				args = append(args, v)
			}
		}

		switch cmd {
		case "translate":
			if len(args) == 1 {
				tts = append(tts, Translation(args[0], 0))
			} else if len(args) == 2 {
				tts = append(tts, Translation(args[0], args[1]))
			} else {
				return errors.New("invalid number of arguments in 'translate' transform")
			}
		case "scale":
			if len(args) == 1 {
				tts = append(tts, Scaling(args[0], args[0]))
			} else if len(args) == 2 {
				tts = append(tts, Scaling(args[0], args[1]))
			} else {
				return errors.New("invalid number of arguments in 'scale' transform")
			}
		case "rotate":
			if len(args) == 1 {
				tts = append(tts,
					Rotation(args[0]))
			} else if len(args) == 3 {
				tts = append(tts,
					Translation(-args[1], -args[2]),
					Rotation(args[0]),
					Translation(args[1], args[2]))
			} else {
				return errors.New("invalid number of arguments in 'rotate' transform")
			}
		case "skewX":
			if len(args) == 1 {
				tts = append(tts, SkewX(args[0]))
			} else {
				return errors.New("invalid number of arguments in 'skewX' transform")
			}
		case "skewY":
			if len(args) == 1 {
				tts = append(tts, SkewY(args[0]))
			} else {
				return errors.New("invalid number of arguments in 'skewY' transform")
			}
		case "matrix":
			if len(args) == 6 {
				tts = append(tts, &Transform{args[0], args[1], args[2], args[3], args[4], args[5]})
			} else {
				return errors.New("invalid number of arguments in 'matrix' transform")
			}
		default:
			return fmt.Errorf("unknown transform '%s'", cmd)
		}
	}

	if len(tts) > 0 {

	}

	return nil
}

type ViewBox string

func (b ViewBox) Parse() (x, y, w, h float64, err error) {
	s := strings.Fields(strings.ReplaceAll(string(b), ",", " "))
	if len(s) != 4 {
		err = errors.New("invalid viewBox attribute")
		return
	}
	x, err = strconv.ParseFloat(s[0], 64)
	if err != nil {
		return
	}
	y, err = strconv.ParseFloat(s[1], 64)
	if err != nil {
		return
	}
	w, err = strconv.ParseFloat(s[2], 64)
	if err != nil {
		return
	}
	h, err = strconv.ParseFloat(s[3], 64)
	return
}
