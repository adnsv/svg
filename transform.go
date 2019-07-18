package svg

import (
	"errors"
	"strconv"
	"strings"
)

type Transform struct {
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
