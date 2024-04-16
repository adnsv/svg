package svg

import (
	"fmt"
	"strconv"
)

func tokenizePoints(s string) ([]float64, error) {
	ret := []float64{}
	cur, last := 0, len(s)

	isDigit := func(c byte) bool {
		return c >= '0' && c <= '9'
	}

	for cur < last {
		if s[cur] <= ' ' || s[cur] == ',' {
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
		ret = append(ret, v)
	}
	return ret, nil
}

func ParsePoints(s string) ([]Vertex, error) {
	vv := []Vertex{}

	dd, err := tokenizePoints(s)
	if err != nil {
		return nil, err
	}

	if len(dd)&1 != 0 {
		return nil, fmt.Errorf("number of coordinates is odd")
	}

	for i := 0; i < len(dd)/2; i++ {
		x := dd[i*2]
		y := dd[i*2+1]
		vv = append(vv, Vertex{X: x, Y: y})
	}
	return vv, nil
}
