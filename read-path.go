package svg

import (
	"fmt"
	"math"
	"strconv"
)

type token struct {
	offset int
	cmd    byte
	num    float64
}

func tokenizePath(s string) ([]token, error) {
	ret := []token{}
	cur, last := 0, len(s)

	isDigit := func(c byte) bool {
		return c >= '0' && c <= '9'
	}

	for cur < last {
		if s[cur] <= ' ' || s[cur] == ',' {
			cur++
			continue
		}
		if (s[cur] >= 'a' && s[cur] <= 'z') || (s[cur] >= 'A' && s[cur] <= 'Z') {
			ret = append(ret, token{offset: cur, cmd: s[cur]})
			cur++
			continue
		}
		start := cur
		if s[cur] == '+' || s[cur] == '-' {
			cur++
		}
		for cur < last && isDigit(s[cur]) {
			cur++
		}
		if cur < last && s[cur] == '.' {
			cur++
			for cur < last && isDigit(s[cur]) {
				cur++
			}
		}
		if cur != start && cur < last && (s[cur] == 'e' || s[cur] == 'E') {
			cur++
			if cur < last && (s[cur] == '+' || s[cur] == '-') {
				cur++
			}
			for cur < last && isDigit(s[cur]) {
				cur++
			}
		}
		if cur == start {
			return nil, fmt.Errorf("invalid content at %d", cur)
		}
		v, err := strconv.ParseFloat(s[start:cur], 64)
		if err != nil {
			return nil, fmt.Errorf("invalid number at %d, %w", cur, err)
		}
		ret = append(ret, token{offset: start, cmd: '#', num: v})
	}
	return ret, nil
}

func ParsePath(s string) (*PathData, error) {

	tokens, err := tokenizePath(s)
	if err != nil {
		return nil, err
	}

	pd := &PathData{}
	pt := Vertex{}
	first := Vertex{0, 0}
	last := Vertex{0, 0}
	cpt := [3]Vertex{}

	lccmd := byte(0)
	prevlccmd := byte(0)

	curToken, lastToken := 0, len(tokens)
	for curToken < lastToken {
		if tokens[curToken].cmd == '#' {
			return nil, fmt.Errorf("unexpected number at %d", tokens[curToken].offset)
		}
		cmd := tokens[curToken].cmd
		offset := tokens[curToken].offset
		curToken++
		values := []float64{}
		for curToken < lastToken && tokens[curToken].cmd == '#' {
			values = append(values, tokens[curToken].num)
			curToken++
		}

		prevlccmd = lccmd
		lccmd = cmd
		rel := true
		if lccmd < 'a' {
			lccmd += 'a' - 'A'
			rel = false
		}

		n := len(values)

		switch lccmd {

		case 'm':
			if n < 2 || n%2 != 0 {
				return nil, fmt.Errorf("invalid num arguments in '%c' command at %d", cmd, offset)
			}
			pt.X = values[0]
			pt.Y = values[1]
			if rel {
				pt.X += last.X
				pt.Y += last.Y
			}
			values = values[2:]
			first = pt
			last = pt
			pd.MoveTo(pt)
			for len(values) >= 2 {
				pt.X = values[0]
				pt.Y = values[1]
				if rel {
					pt.X += last.X
					pt.Y += last.Y
				}
				values = values[2:]
				last = pt
				pd.LineTo(pt)
			}

		case 'z':
			if n != 0 {
				return nil, fmt.Errorf("invalid num arguments in '%c' command at %d", cmd, offset)
			}
			last = first
			pd.Close()

		case 'l':
			if n < 2 || n%2 != 0 {
				return nil, fmt.Errorf("invalid num arguments in '%c' command at %d", cmd, offset)
			}
			for len(values) >= 2 {
				pt.X = values[0]
				pt.Y = values[1]
				if rel {
					pt.X += last.X
					pt.Y += last.Y
				}
				values = values[2:]
				last = pt
				pd.LineTo(pt)
			}

		case 'h':
			if n == 0 {
				return nil, fmt.Errorf("invalid num arguments in '%c' command at %d", cmd, offset)
			}
			for len(values) >= 1 {
				pt.X = values[0]
				pt.Y = last.Y
				if rel {
					pt.X += last.X
				}
				values = values[1:]
				last = pt
				pd.LineTo(pt)
			}

		case 'v':
			if n == 0 {
				return nil, fmt.Errorf("invalid num arguments in '%c' command at %d", cmd, offset)
			}
			for len(values) >= 1 {
				pt.X = last.X
				pt.Y = values[0]
				if rel {
					pt.Y += last.Y
				}
				values = values[1:]
				last = pt
				pd.LineTo(pt)
			}

		case 'c':
			if n < 6 || n%6 != 0 {
				return nil, fmt.Errorf("invalid num arguments in '%c' command at %d", cmd, offset)
			}
			for len(values) >= 6 {
				cpt[0].X = values[0]
				cpt[0].Y = values[1]
				cpt[1].X = values[2]
				cpt[1].Y = values[3]
				cpt[2].X = values[4]
				cpt[2].Y = values[5]
				if rel {
					cpt[0].X += last.X
					cpt[0].Y += last.Y
					cpt[1].X += last.X
					cpt[1].Y += last.Y
					cpt[2].X += last.X
					cpt[2].Y += last.Y
				}
				values = values[6:]
				last = cpt[2]
				pd.CurveTo(cpt[0], cpt[1], cpt[2])
			}

		case 's':
			if n < 4 || n%4 != 0 {
				return nil, fmt.Errorf("invalid num arguments in '%c' command at %d", cmd, offset)
			}
			if prevlccmd != 's' && prevlccmd != 'c' {
				cpt[1] = last
			}
			for len(values) >= 4 {
				cpt[0].X = last.X*2 - cpt[1].X
				cpt[0].Y = last.Y*2 - cpt[1].Y
				cpt[1].X = values[0]
				cpt[1].Y = values[1]
				cpt[2].X = values[2]
				cpt[2].Y = values[3]
				if rel {
					cpt[1].X += last.X
					cpt[1].Y += last.Y
					cpt[2].X += last.X
					cpt[2].Y += last.Y
				}
				values = values[4:]
				last = cpt[2]
				pd.CurveTo(cpt[0], cpt[1], cpt[2])
			}

		case 'q':
			if n < 4 || n%4 != 0 {
				return nil, fmt.Errorf("invalid num arguments in '%c' command at %d", cmd, offset)
			}
			for len(values) >= 4 {
				cpt[0].X = values[0]
				cpt[0].Y = values[1]
				cpt[2].X = values[2]
				cpt[2].Y = values[3]
				if rel {
					cpt[0].X += last.X
					cpt[0].Y += last.Y
					cpt[2].X += last.X
					cpt[2].Y += last.Y
				}
				cpt[1] = cpt[0]
				values = values[4:]
				last = cpt[2]
				pd.CurveTo(cpt[0], cpt[1], cpt[2])
			}

		case 't':
			if n < 2 || n%2 != 0 {
				return nil, fmt.Errorf("invalid num arguments in '%c' command at %d", cmd, offset)
			}
			if prevlccmd != 't' && prevlccmd != 'q' {
				cpt[1] = last
			}
			for len(values) >= 2 {
				cpt[2].X = values[0]
				cpt[2].Y = values[1]
				cpt[0].X = last.X*2 - cpt[1].X
				cpt[0].Y = last.Y*2 - cpt[1].Y
				cpt[1] = cpt[0]
				if rel {
					cpt[2].X += last.X
					cpt[2].Y += last.Y
				}
				values = values[2:]
				last = cpt[2]
				pd.CurveTo(cpt[0], cpt[1], cpt[2])
			}

		case 'a':
			if n < 7 || n%7 != 0 {
				return nil, fmt.Errorf("invalid num arguments in '%c' command at %d", cmd, offset)
			}
			for len(values) >= 7 {
				r := Vertex{values[0], values[1]}
				axisRotation := values[2] * math.Pi / 180.0
				largeArc := values[3] != 0
				sweep := values[4] != 0
				pt.X = values[5]
				pt.Y = values[6]
				values = values[7:]

				pd.Arc(last.X, last.Y, r.X, r.Y, axisRotation, largeArc, sweep, pt.X, pt.Y)
				last = pt
			}
		default:
			return nil, fmt.Errorf("invalid path command '%c' at %d", cmd, offset)
		}
	}
	return pd, nil
}

func (pd *PathData) ArcSegment(xc, yc, th0, th1, rx, ry, axisRotation float64) {
	h := 0.5 * (th1 - th0)
	t := (8.0 / 3.0) * math.Sin(h*0.5) * math.Sin(h*0.5) / math.Sin(h)
	x1 := xc + math.Cos(th0) - t*math.Sin(th0)
	y1 := yc + math.Sin(th0) + t*math.Cos(th0)
	x3 := xc + math.Cos(th1)
	y3 := yc + math.Sin(th1)
	x2 := x3 + t*math.Sin(th1)
	y2 := y3 - t*math.Cos(th1)

	s, c := math.Sincos(axisRotation)
	a00 := c * rx
	a01 := -s * ry
	a10 := s * rx
	a11 := c * ry

	pd.CurveTo(
		Vertex{a00*x1 + a01*y1, a10*x1 + a11*y1},
		Vertex{a00*x2 + a01*y2, a10*x2 + a11*y2},
		Vertex{a00*x3 + a01*y3, a10*x3 + a11*y3})
}

func (pd *PathData) Arc(lastX, lastY, rx, ry, axisRotation float64, largeArc, sweep bool, x, y float64) {
	rx = math.Abs(rx)
	ry = math.Abs(ry)

	curx := lastX
	cury := lastY

	sin_th, cos_th := math.Sincos(axisRotation * (math.Pi / 180.0))

	dx := (curx - x) / 2.0
	dy := (cury - y) / 2.0
	dx1 := cos_th*dx + sin_th*dy
	dy1 := -sin_th*dx + cos_th*dy
	Pr1 := rx * rx
	Pr2 := ry * ry
	Px := dx1 * dx1
	Py := dy1 * dy1

	// Spec : check if radii are large enough
	check := Px/Pr1 + Py/Pr2
	if check > 1 {
		rx = rx * math.Sqrt(check)
		ry = ry * math.Sqrt(check)
	}

	a00 := cos_th / rx
	a01 := sin_th / rx
	a10 := -sin_th / ry
	a11 := cos_th / ry
	x0 := a00*curx + a01*cury
	y0 := a10*curx + a11*cury
	x1 := a00*x + a01*y
	y1 := a10*x + a11*y
	/* (x0, y0) is current point in transformed coordinate space.
	   (x1, y1) is new point in transformed coordinate space.

	   The arc fits a unit-radius circle in this space. */

	d := (x1-x0)*(x1-x0) + (y1-y0)*(y1-y0)
	sfactor_sq := 1.0/d - 0.25
	if sfactor_sq < 0 {
		sfactor_sq = 0
	}
	sfactor := math.Sqrt(sfactor_sq)
	if sweep == largeArc {
		sfactor = -sfactor
	}

	// arc center
	xc := 0.5*(x0+x1) - sfactor*(y1-y0)
	yc := 0.5*(y0+y1) + sfactor*(x1-x0)

	th0 := math.Atan2(y0-yc, x0-xc)
	th1 := math.Atan2(y1-yc, x1-xc)

	th_arc := th1 - th0
	if (th_arc < 0) && sweep {
		th_arc = th_arc + 2*math.Pi
	} else if th_arc > 0 && !sweep {
		th_arc = th_arc - 2*math.Pi
	}

	/* XXX: I still need to evaluate the math performed in this
	   function. The critical behavior desired is that the arc must be
	   approximated within an arbitrary error tolerance, (which the
	   user should be able to specify as well). I don't yet know the
	   bounds of the error from the following computation of
	   n_segs. Plus the "+ 0.001" looks just plain fishy. -cworth */
	n_segs := math.Round(0.5 + math.Abs(th_arc/(math.Pi*0.5+0.001)))

	for i := 0; i < int(n_segs); i++ {
		pd.ArcSegment(xc, yc,
			th0+float64(i)*th_arc/n_segs,
			th0+float64(i+1)*th_arc/n_segs,
			rx, ry, axisRotation)
	}
}
