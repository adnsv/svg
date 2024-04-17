package svg

import (
	"fmt"
	"strconv"
)

func ParsePaint(s string) (*Paint, error) {
	if len(s) == 0 {
		return nil, fmt.Errorf("empty specs")
	}
	if s == "none" {
		return &Paint{Kind: PaintKindNone}, nil
	}
	if s[0] == '#' && len(s) == 7 || len(s) == 4 {
		// hex rgba
		v, err := strconv.ParseUint(s[1:], 16, 32)
		if err != nil {
			return nil, err
		}
		p := &Paint{Kind: PaintKindRGB}
		if len(s) == 7 {
			p.Color.R = uint8((v >> 16) & 0xff)
			p.Color.G = uint8((v >> 8) & 0xff)
			p.Color.B = uint8(v & 0xff)
		} else {
			p.Color.R = uint8((v>>8)&0xf) * 0x11
			p.Color.G = uint8((v>>4)&0xf) * 0x11
			p.Color.B = uint8(v&0xf) * 0x11
		}
		return p, nil
	}

	return nil, fmt.Errorf("unsupported specs")
}

func ParseOpacity(s string) (*float64, error) {
	if len(s) == 0 {
		return nil, fmt.Errorf("empty specs")
	}
	v, err := strconv.ParseFloat(s, 64)
	if err != nil {
		return nil, err
	}
	return &v, nil
}
