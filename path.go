package svg

import (
	"fmt"
	"strconv"
)

// PathCommand represents a command in svg.Path D (Data) attribute
type PathCommand struct {
	Cmd    byte
	Values []float64
}

type token struct {
	offset int
	cmd    byte
	num    float64
}

func tokenize(s string) ([]token, error) {
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

func ParsePath(s string) ([]PathCommand, error) {

	tokens, err := tokenize(s)
	if err != nil {
		return nil, err
	}

	ret := []PathCommand{}

	cur, last := 0, len(tokens)
	for cur < last {
		if tokens[cur].cmd == '#' {
			return nil, fmt.Errorf("unexpected number at %d", tokens[cur].offset)
		}
		cmd := tokens[cur].cmd
		offset := tokens[cur].offset
		cur++
		values := []float64{}
		for cur < last && tokens[cur].cmd == '#' {
			values = append(values, tokens[cur].num)
			cur++
		}

		n := len(values)

		switch cmd {
		case 'm':
			if n < 2 || n%2 != 0 {
				return nil, fmt.Errorf("invalid num arguments in '%c' command at %d", cmd, offset)
			}
			ret = append(ret, PathCommand{'m', []float64{values[0], values[1]}})
			values = values[2:]
			for len(values) >= 2 {
				ret = append(ret, PathCommand{'l', []float64{values[0], values[1]}})
				values = values[2:]
			}

		case 'M':
			if n < 2 || n%2 != 0 {
				return nil, fmt.Errorf("invalid num arguments in '%c' command at %d", cmd, offset)
			}
			ret = append(ret, PathCommand{'M', []float64{values[0], values[1]}})
			values = values[2:]
			for len(values) >= 2 {
				ret = append(ret, PathCommand{'L', []float64{values[0], values[1]}})
				values = values[2:]
			}

		case 'z':
			if n != 0 {
				return nil, fmt.Errorf("invalid num arguments in '%c' command at %d", cmd, offset)
			}
		case 'Z':
			if n != 0 {
				return nil, fmt.Errorf("invalid num arguments in '%c' command at %d", cmd, offset)
			}
		case 'L':
			if n < 2 || n%2 != 0 {
				return nil, fmt.Errorf("invalid num arguments in '%c' command at %d", cmd, offset)
			}
			for len(values) >= 2 {
				ret = append(ret, PathCommand{'L', values[:2]})
				values = values[2:]
			}
		case 'l':
			if n < 2 || n%2 != 0 {
				return nil, fmt.Errorf("invalid num arguments in '%c' command at %d", cmd, offset)
			}
			for len(values) >= 2 {
				ret = append(ret, PathCommand{'l', values[:2]})
				values = values[2:]
			}
		case 'H':
			if n == 0 {
				return nil, fmt.Errorf("invalid num arguments in '%c' command at %d", cmd, offset)
			}
			for len(values) >= 1 {
				ret = append(ret, PathCommand{'H', values[:1]})
				values = values[1:]
			}
		case 'h':
			if n == 0 {
				return nil, fmt.Errorf("invalid num arguments in '%c' command at %d", cmd, offset)
			}
			for len(values) >= 1 {
				ret = append(ret, PathCommand{'h', values[:1]})
				values = values[1:]
			}
		case 'V':
			if n == 0 {
				return nil, fmt.Errorf("invalid num arguments in '%c' command at %d", cmd, offset)
			}
			for len(values) >= 1 {
				ret = append(ret, PathCommand{'V', values[:1]})
				values = values[1:]
			}
		case 'v':
			if n == 0 {
				return nil, fmt.Errorf("invalid num arguments in '%c' command at %d", cmd, offset)
			}
			for len(values) >= 1 {
				ret = append(ret, PathCommand{'v', values[:1]})
				values = values[1:]
			}
		case 'C':
			if n < 6 || n%6 != 0 {
				return nil, fmt.Errorf("invalid num arguments in '%c' command at %d", cmd, offset)
			}
			for len(values) >= 6 {
				ret = append(ret, PathCommand{'C', values[:6]})
				values = values[6:]
			}
		case 'c':
			if n < 6 || n%6 != 0 {
				return nil, fmt.Errorf("invalid num arguments in '%c' command at %d", cmd, offset)
			}
			for len(values) >= 6 {
				ret = append(ret, PathCommand{'c', values[:6]})
				values = values[6:]
			}
		case 'S':
			if n < 4 || n%4 != 0 {
				return nil, fmt.Errorf("invalid num arguments in '%c' command at %d", cmd, offset)
			}
			for len(values) >= 4 {
				ret = append(ret, PathCommand{'S', values[:4]})
				values = values[4:]
			}
		case 's':
			if n < 4 || n%4 != 0 {
				return nil, fmt.Errorf("invalid num arguments in '%c' command at %d", cmd, offset)
			}
			for len(values) >= 4 {
				ret = append(ret, PathCommand{'s', values[:4]})
				values = values[4:]
			}
		case 'Q':
			if n < 4 || n%4 != 0 {
				return nil, fmt.Errorf("invalid num arguments in '%c' command at %d", cmd, offset)
			}
			for len(values) >= 4 {
				ret = append(ret, PathCommand{'Q', values[:4]})
				values = values[4:]
			}
		case 'q':
			if n < 4 || n%4 != 0 {
				return nil, fmt.Errorf("invalid num arguments in '%c' command at %d", cmd, offset)
			}
			for len(values) >= 4 {
				ret = append(ret, PathCommand{'q', values[:4]})
				values = values[4:]
			}
		case 'T':
			if n < 2 || n%2 != 0 {
				return nil, fmt.Errorf("invalid num arguments in '%c' command at %d", cmd, offset)
			}
			for len(values) >= 2 {
				ret = append(ret, PathCommand{'T', values[:2]})
				values = values[2:]
			}
		case 't':
			if n < 2 || n%2 != 0 {
				return nil, fmt.Errorf("invalid num arguments in '%c' command at %d", cmd, offset)
			}
			for len(values) >= 2 {
				ret = append(ret, PathCommand{'t', values[:2]})
				values = values[2:]
			}
		case 'A':
			if n < 7 || n%7 != 0 {
				return nil, fmt.Errorf("invalid num arguments in '%c' command at %d", cmd, offset)
			}
			for len(values) >= 7 {
				ret = append(ret, PathCommand{'A', values[:7]})
				values = values[7:]
			}
		case 'a':
			if n < 7 || n%7 != 0 {
				return nil, fmt.Errorf("invalid num arguments in '%c' command at %d", cmd, offset)
			}
			for len(values) >= 7 {
				ret = append(ret, PathCommand{'a', values[:7]})
				values = values[7:]
			}
		default:
			return nil, fmt.Errorf("invalid path command '%c' at %d", cmd, offset)
		}
	}
	return ret, nil
}

func IteratePath(pcs []PathCommand, iterator func(index int, cmd byte, s []float64, x, y float64)) {
	lastX := float64(0)
	lastY := float64(0)
	startX := float64(0)
	startY := float64(0)

	for i, pc := range pcs {
		iterator(i, pc.Cmd, pc.Values, lastX, lastY)
		switch pc.Cmd {
		case 'm':
			lastX += pc.Values[0]
			lastY += pc.Values[1]
			startX = lastX
			startY = lastY
		case 'M':
			lastX = pc.Values[0]
			lastY = pc.Values[1]
			startX = lastX
			startY = lastY
		case 'h':
			lastX += pc.Values[0]
		case 'H':
			lastX = pc.Values[0]
		case 'v':
			lastY += pc.Values[0]
		case 'V':
			lastY = pc.Values[0]
		case 'z':
			lastX = startX
			lastY = startY
		case 'Z':
			lastX = startX
			lastY = startY
		case 'l':
			lastX += pc.Values[0]
			lastY += pc.Values[1]
		case 'L':
			lastX = pc.Values[0]
			lastY = pc.Values[1]
		case 'c':
			lastX += pc.Values[4]
			lastY += pc.Values[5]
		case 'C':
			lastX = pc.Values[4]
			lastY = pc.Values[5]
		case 's':
			lastX += pc.Values[2]
			lastY += pc.Values[3]
		case 'S':
			lastX = pc.Values[2]
			lastY = pc.Values[3]
		case 'q':
			lastX += pc.Values[2]
			lastY += pc.Values[3]
		case 'Q':
			lastX = pc.Values[2]
			lastY = pc.Values[3]
		case 't':
			lastX += pc.Values[0]
			lastY += pc.Values[1]
		case 'T':
			lastX = pc.Values[0]
			lastY = pc.Values[1]
		case 'a':
			lastX += pc.Values[5]
			lastY += pc.Values[6]
		case 'A':
			lastX = pc.Values[5]
			lastY = pc.Values[6]
		}
	}

}

func TransformPath(pcs []PathCommand, t *Transform) []PathCommand {
	ret := []PathCommand{}

	IteratePath(pcs, func(index int, cmd byte, s []float64, lx, ly float64) {
		if cmd == 'm' && index == 0 {
			cmd = 'M'
		}
		switch cmd {
		case 'v':
			dx, dy := t.Calc(0, s[0], true)
			ret = append(ret, PathCommand{'l', []float64{dx, dy}})
		case 'V':
			x, y := t.Calc(lx, s[0], false)
			ret = append(ret, PathCommand{'L', []float64{x, y}})
		case 'h':
			dx, dy := t.Calc(s[0], 0, true)
			ret = append(ret, PathCommand{'l', []float64{dx, dy}})
		case 'H':
			x, y := t.Calc(s[0], ly, false)
			ret = append(ret, PathCommand{'L', []float64{x, y}})
		}
	})

	return ret
}
