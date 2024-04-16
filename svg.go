package svg

import (
	"fmt"
)

type Item interface {
	writer
	ID() string
}

type item struct {
	id string
}

func (it *item) ID() string {
	return it.id
}

func (it *item) read(src sourcer) (err error) {
	it.id, _ = src.Attr("id")
	return nil
}

func (it *item) write(tgt targeter) {
	if len(it.id) > 0 {
		tgt.Attr("id", it.id)
	}
}

type Node struct {
	item
	Items []Item
}

func (n *Node) read(src sourcer) error {
	err := n.item.read(src)
	if err != nil {
		return err
	}
	return src.ForEachChildNode(func(tag string, cs sourcer) error {
		var it reader

		switch tag {
		case "g":
			it = &Group{}
		case "line":
			it = &Line{}
		case "rect":
			it = &Rect{}
		case "circle":
			it = &Circle{}
		case "ellipse":
			it = &Ellipse{}
		case "polyline":
			it = &Polygon{}
		case "polygon":
			it = &Polygon{}
		case "path":
			it = &Path{}
		case "text":
			// todo: implement
		}

		if it != nil {
			err := it.read(cs)
			if err != nil {
				return fmt.Errorf("in <%s>: %s", tag, err)
			}
			n.Items = append(n.Items, it)
		}
		return nil
	})
}

func (n *Node) write(tgt targeter) {
	n.item.write(tgt)
	for _, it := range n.Items {
		tag := ""
		switch it.(type) {
		case *Group:
			tag = "g"
		case *Line:
			tag = "line"
		case *Rect:
			tag = "rect"
		case *Circle:
			tag = "circle"
		case *Ellipse:
			tag = "ellipse"
		case *Polyline:
			tag = "polyline"
		case *Polygon:
			tag = "polygon"
		case *Path:
			tag = "path"
		default:
			panic("inknown element tag")
		}
		tgt.Child(tag, func(t targeter) {
			it.write(t)
		})
	}
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

func (s *Shape) read(src sourcer) (err error) {
	err = s.item.read(src)
	if err != nil {
		return
	}

	if v, exists := src.Attr("fill"); exists {
		s.Fill, err = ParsePaint(v)
		if err != nil {
			return fmt.Errorf("invalid fill: %w", err)
		}
	}
	if v, exists := src.Attr("fill-rule"); exists {
		r := FillRuleInherit
		err = r.UnmarshalText([]byte(v))
		if err != nil {
			return fmt.Errorf("invalid fill-rule: %w", err)
		}
	}

	if v, exists := src.Attr("fill-opacity"); exists {
		s.FillOpacity, err = ParseOpacity(v)
		if err != nil {
			return fmt.Errorf("invalid fill-opacity: %w", err)
		}
	}

	if v, exists := src.Attr("opacity"); exists {
		s.Opacity, err = ParseOpacity(v)
		if err != nil {
			return fmt.Errorf("invalid opacity: %w", err)
		}
	}

	return
}

func (s *Shape) write(tgt targeter) {
	s.item.write(tgt)
}

type Group struct {
	Node
	Transform *Transform
}

func (g *Group) read(src sourcer) (err error) {
	err = g.Node.read(src)
	if err != nil {
		return
	}
	return
}

func (g *Group) write(tgt targeter) {
	g.Node.write(tgt)
}

type Defs struct {
	Node
	Transform *Transform
}

// Svg implements SVG <svg> element
type Svg struct {
	Group
	ViewBox ViewBox
	X       Coordinate
	Y       Coordinate
	Width   Length
	Height  Length
}

func (svg *Svg) read(src sourcer) (err error) {
	if s, ok := src.Attr("viewBox"); ok {
		svg.ViewBox = ViewBox(s)
	}
	if s, ok := src.Attr("x"); ok {
		svg.X = Coordinate(s)
	}
	if s, ok := src.Attr("y"); ok {
		svg.Y = Coordinate(s)
	}
	if s, ok := src.Attr("width"); ok {
		svg.Width = Length(s)
	}
	if s, ok := src.Attr("height"); ok {
		svg.Height = Length(s)
	}
	err = svg.Group.read(src)
	if err != nil {
		return err
	}
	return nil
}

func (svg *Svg) write(tgt targeter) {
	tgt.Attr("viewBox", string(svg.ViewBox))
	tgt.Attr("x", string(svg.X))
	tgt.Attr("y", string(svg.Y))
	tgt.Attr("width", string(svg.Width))
	tgt.Attr("height", string(svg.Height))
	svg.Group.write(tgt)
}

type Line struct {
	Shape
	X1 Coordinate
	Y1 Coordinate
	X2 Coordinate
	Y2 Coordinate
}

func (l *Line) read(src sourcer) (err error) {
	err = l.Shape.read(src)
	if err != nil {
		return
	}
	if s, ok := src.Attr("x1"); ok {
		l.X1 = Coordinate(s)
	}
	if s, ok := src.Attr("y1"); ok {
		l.Y1 = Coordinate(s)
	}
	if s, ok := src.Attr("x2"); ok {
		l.X2 = Coordinate(s)
	}
	if s, ok := src.Attr("y2"); ok {
		l.Y2 = Coordinate(s)
	}
	return
}

func (l *Line) write(tgt targeter) {
	l.Shape.write(tgt)
	tgt.Attr("x1", string(l.X1))
	tgt.Attr("y1", string(l.Y1))
	tgt.Attr("x2", string(l.X2))
	tgt.Attr("y2", string(l.Y2))
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

func (r *Rect) read(src sourcer) (err error) {
	err = r.Shape.read(src)
	if err != nil {
		return
	}
	if s, ok := src.Attr("x"); ok {
		r.X = Coordinate(s)
	}
	if s, ok := src.Attr("y"); ok {
		r.Y = Coordinate(s)
	}
	if s, ok := src.Attr("width"); ok {
		r.Width = Length(s)
	}
	if s, ok := src.Attr("height"); ok {
		r.Height = Length(s)
	}
	if s, ok := src.Attr("rx"); ok {
		r.Rx = Length(s)
	}
	if s, ok := src.Attr("ry"); ok {
		r.Ry = Length(s)
	}
	return
}

func (r *Rect) write(tgt targeter) {
	r.Shape.write(tgt)
	tgt.Attr("x", string(r.X))
	tgt.Attr("y", string(r.Y))
	tgt.Attr("width", string(r.Width))
	tgt.Attr("height", string(r.Height))
	tgt.Attr("rx", string(r.Rx))
	tgt.Attr("ry", string(r.Ry))
}

type Circle struct {
	Shape
	Cx     Coordinate
	Cy     Coordinate
	Radius Length
}

func (c *Circle) read(src sourcer) (err error) {
	err = c.Shape.read(src)
	if err != nil {
		return
	}
	if s, ok := src.Attr("cx"); ok {
		c.Cx = Coordinate(s)
	}
	if s, ok := src.Attr("cy"); ok {
		c.Cy = Coordinate(s)
	}
	if s, ok := src.Attr("r"); ok {
		c.Radius = Length(s)
	}
	return
}

func (c *Circle) write(tgt targeter) {
	c.Shape.write(tgt)
	tgt.Attr("cx", string(c.Cx))
	tgt.Attr("cy", string(c.Cy))
	tgt.Attr("r", string(c.Radius))
}

type Ellipse struct {
	Shape
	Cx Coordinate
	Cy Coordinate
	Rx Length
	Ry Length
}

func (e *Ellipse) read(src sourcer) (err error) {
	err = e.Shape.read(src)
	if err != nil {
		return
	}
	if s, ok := src.Attr("cx"); ok {
		e.Cx = Coordinate(s)
	}
	if s, ok := src.Attr("cy"); ok {
		e.Cy = Coordinate(s)
	}
	if s, ok := src.Attr("rx"); ok {
		e.Rx = Length(s)
	}
	if s, ok := src.Attr("ry"); ok {
		e.Ry = Length(s)
	}
	return
}

func (e *Ellipse) write(tgt targeter) {
	e.Shape.write(tgt)
	tgt.Attr("cx", string(e.Cx))
	tgt.Attr("cy", string(e.Cy))
	tgt.Attr("rx", string(e.Rx))
	tgt.Attr("ry", string(e.Ry))
}

type Polyline struct {
	Shape
	Points string
}

func (p *Polyline) read(src sourcer) (err error) {
	err = p.Shape.read(src)
	if err != nil {
		return
	}
	if s, ok := src.Attr("points"); ok {
		p.Points = s
	}
	return
}

func (p *Polyline) write(tgt targeter) {
	p.Shape.write(tgt)
	tgt.Attr("points", p.Points)
}

type Polygon struct {
	Shape
	Points string
}

func (p *Polygon) read(src sourcer) (err error) {
	err = p.Shape.read(src)
	if err != nil {
		return
	}
	if s, ok := src.Attr("points"); ok {
		p.Points = s
	}
	return
}

func (p *Polygon) write(tgt targeter) {
	p.Shape.write(tgt)
	tgt.Attr("points", p.Points)
}

type Path struct {
	Shape
	D string
}

func (p *Path) read(src sourcer) (err error) {
	err = p.Shape.read(src)
	if err != nil {
		return
	}
	if s, ok := src.Attr("d"); ok {
		p.D = s
	}
	return
}

func (p *Path) write(tgt targeter) {
	p.Shape.write(tgt)
	tgt.Attr("d", p.D)
}
