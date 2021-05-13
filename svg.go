package svg

type Item interface {
	ID() string
}

type sourcer interface {
	Attr(name string) (v string, exists bool)
	ForEachChildNode(callback func(tag string, ch sourcer) error) error
}

type reader interface {
	Item
	read(src sourcer) error
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
		case "path":
			it = &Path{}
		}

		if it != nil {
			err := it.read(cs)
			if err != nil {
				return err
			}
			n.Items = append(n.Items, it)
		}
		return nil
	})
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
	return
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

func (svg *Svg) read(src sourcer) (err error) {

	svg.ViewBox = nil
	if s, ok := src.Attr("viewBox"); ok {
		svg.ViewBox = &ViewBox{}
		err = svg.ViewBox.Unmarshal(s)
		if err != nil {
			return
		}
	}

	svg.X, err = attrCoordinate(src, "x")
	if err != nil {
		return err
	}
	svg.Y, err = attrCoordinate(src, "y")
	if err != nil {
		return err
	}
	svg.Width, err = attrLength(src, "width")
	if err != nil {
		return err
	}
	svg.Height, err = attrLength(src, "height")
	if err != nil {
		return err
	}

	err = svg.Group.read(src)
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

func (l *Line) read(src sourcer) (err error) {
	err = l.Shape.read(src)
	if err != nil {
		return
	}
	if s, ok := src.Attr("x1"); ok {
		err = l.X1.Unmarshal(s)
		if err != nil {
			return
		}
	}
	if s, ok := src.Attr("y1"); ok {
		err = l.Y1.Unmarshal(s)
		if err != nil {
			return
		}
	}
	if s, ok := src.Attr("x2"); ok {
		err = l.X2.Unmarshal(s)
		if err != nil {
			return
		}
	}
	if s, ok := src.Attr("y2"); ok {
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

func (r *Rect) read(src sourcer) (err error) {
	err = r.Shape.read(src)
	if err != nil {
		return
	}
	if s, ok := src.Attr("x"); ok {
		err = r.X.Unmarshal(s)
		if err != nil {
			return
		}
	}
	if s, ok := src.Attr("y"); ok {
		err = r.Y.Unmarshal(s)
		if err != nil {
			return
		}
	}
	if s, ok := src.Attr("width"); ok {
		err = r.Width.Unmarshal(s)
		if err != nil {
			return
		}
	}
	if s, ok := src.Attr("height"); ok {
		err = r.Height.Unmarshal(s)
		if err != nil {
			return
		}
	}
	sx, okx := src.Attr("rx")
	sy, oky := src.Attr("ry")
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

func (c *Circle) read(src sourcer) (err error) {
	err = c.Shape.read(src)
	if err != nil {
		return
	}
	if s, ok := src.Attr("cx"); ok {
		err = c.Cx.Unmarshal(s)
		if err != nil {
			return
		}
	}
	if s, ok := src.Attr("cy"); ok {
		err = c.Cy.Unmarshal(s)
		if err != nil {
			return
		}
	}
	if s, ok := src.Attr("r"); ok {
		err = c.Radius.Unmarshal(s)
		if err != nil {
			return
		}
	}
	return
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
		err = e.Cx.Unmarshal(s)
		if err != nil {
			return
		}
	}
	if s, ok := src.Attr("cy"); ok {
		err = e.Cy.Unmarshal(s)
		if err != nil {
			return
		}
	}
	if s, ok := src.Attr("rx"); ok {
		err = e.Rx.Unmarshal(s)
		if err != nil {
			return
		}
	}
	if s, ok := src.Attr("ry"); ok {
		err = e.Ry.Unmarshal(s)
		if err != nil {
			return
		}
	}
	return
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
