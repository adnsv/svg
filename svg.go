package svg

import (
	"encoding/xml"
	"errors"
	"io"
)

type Item interface {
	ID() string
	read(elt *element) error
}

type item struct {
	id string
}

func (it *item) ID() string {
	return it.id
}

func (it *item) read(elt *element) error {
	it.id = elt.attributes["id"]
	return nil
}

type Node struct {
	item
	Items []Item
}

func (n *Node) read(elt *element) error {
	err := n.item.read(elt)
	if err != nil {
		return err
	}
	for _, ce := range elt.children {

		var it Item

		switch ce.name {
		case "g":
			it = &Group{}
		case "line":
			it = &Line{}
		case "rect":
			it = &Rect{}
		case "circle":
			it = &Circle{}
		case "path":
			it = &Path{}
		}

		if it != nil {
			err := it.read(ce)
			if err != nil {
				return err
			}
			n.Items = append(n.Items, it)
		}
	}
	return nil
}

type Shape struct {
	item
	Fill           *Paint
	FillRule       *FillRule
	FillOpacity    *float64
	Stroke         *Paint
	StrokeWidth    *Length
	StrokeOpacity  *float64
	StrokeLineCap  *LineCap
	StrokeLineJoin *LineJoin
	Opacity        *float64
	Transform      *Transform
}

func (s *Shape) read(elt *element) (err error) {
	err = s.item.read(elt)
	if err != nil {
		return
	}
	return
}

type Group struct {
	Node
	Transform *Transform
}

func (g *Group) read(elt *element) (err error) {
	err = g.Node.read(elt)
	if err != nil {
		return
	}
	return
}

type Defs struct {
	Node
	Transform *Transform
}

// Svg implements SVG <svg> element
type Svg struct {
	Group
	ViewBox *ViewBox
	X       *Coordinate
	Y       *Coordinate
	Width   *Length
	Height  *Length
}

func (svg *Svg) read(elt *element) (err error) {

	err = svg.Group.read(elt)
	if err != nil {
		return err
	}

	svg.ViewBox = nil
	if s, ok := elt.attributes["viewBox"]; ok {
		svg.ViewBox = &ViewBox{}
		err = svg.ViewBox.Unmarshal(s)
		if err != nil {
			return
		}
	}

	svg.X, err = attrCoordinate(elt, "x")
	if err != nil {
		return err
	}
	svg.Y, err = attrCoordinate(elt, "y")
	if err != nil {
		return err
	}
	svg.Width, err = attrLength(elt, "width")
	if err != nil {
		return err
	}
	svg.Height, err = attrLength(elt, "height")
	if err != nil {
		return err
	}

	return nil
}

type Line struct {
	Shape
	X1 Coordinate
	Y1 Coordinate
	X2 Coordinate
	Y2 Coordinate
}

func (l *Line) read(elt *element) (err error) {
	err = l.Shape.read(elt)
	if err != nil {
		return
	}
	if s, ok := elt.attributes["x1"]; ok {
		err = l.X1.Unmarshal(s)
		if err != nil {
			return
		}
	}
	if s, ok := elt.attributes["y1"]; ok {
		err = l.Y1.Unmarshal(s)
		if err != nil {
			return
		}
	}
	if s, ok := elt.attributes["x2"]; ok {
		err = l.X2.Unmarshal(s)
		if err != nil {
			return
		}
	}
	if s, ok := elt.attributes["y2"]; ok {
		err = l.Y2.Unmarshal(s)
		if err != nil {
			return
		}
	}
	return
}

type Rect struct {
	Shape
	X      Coordinate
	Y      Coordinate
	Width  Length
	Height Length
	Rx     Length
	Ry     Length
}

func (r *Rect) read(elt *element) (err error) {
	err = r.Shape.read(elt)
	if err != nil {
		return
	}
	if s, ok := elt.attributes["x"]; ok {
		err = r.X.Unmarshal(s)
		if err != nil {
			return
		}
	}
	if s, ok := elt.attributes["y"]; ok {
		err = r.Y.Unmarshal(s)
		if err != nil {
			return
		}
	}
	if s, ok := elt.attributes["width"]; ok {
		err = r.Width.Unmarshal(s)
		if err != nil {
			return
		}
	}
	if s, ok := elt.attributes["height"]; ok {
		err = r.Height.Unmarshal(s)
		if err != nil {
			return
		}
	}
	sx, okx := elt.attributes["rx"]
	sy, oky := elt.attributes["ry"]
	if okx {
		err = r.Width.Unmarshal(sx)
		if err != nil {
			return err
		}
	}
	if oky {
		err = r.Width.Unmarshal(sy)
		if err != nil {
			return err
		}
	}
	if okx && !oky {
		r.Ry = r.Rx
	} else if !okx && oky {
		r.Rx = r.Ry
	}

	return
}

type Circle struct {
	Shape
	Cx     Coordinate
	Cy     Coordinate
	Radius Length
}

type Path struct {
	Shape
	D string
}

func (p *Path) read(elt *element) (err error) {
	err = p.Shape.read(elt)
	if err != nil {
		return
	}
	if s, ok := elt.attributes["d"]; ok {
		p.D = s
	}
	return
}

func Parse(in io.Reader) (*Svg, error) {
	decoder := xml.NewDecoder(in)
	element, err := decodeFirst(decoder)
	if err != nil {
		return nil, err
	}
	if err := element.decode(decoder); err != nil {
		return nil, err
	}
	if element == nil || element.name != "svg" {
		return nil, errors.New("invalid root element")
	}
	document := &Svg{}
	err = document.read(element)
	if err != nil {
		return nil, err
	}
	return document, nil
}
