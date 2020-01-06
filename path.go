package svg

// PathData represents a command in svg.Path D (Data) attribute
type PathData struct {
	Commands []PathCommand
	Vertices []Vertex
}

type PathCommand uint8

const (
	PathClose = PathCommand(iota)
	PathMoveTo
	PathLineTo
	PathCurveTo
)

type Vertex struct {
	X float64
	Y float64
}

func (p *PathData) Close() {
	p.Commands = append(p.Commands, PathClose)
}

func (p *PathData) MoveTo(v Vertex) {
	p.Commands = append(p.Commands, PathMoveTo)
	p.Vertices = append(p.Vertices, v)
}

func (p *PathData) LineTo(v Vertex) {
	p.Commands = append(p.Commands, PathLineTo)
	p.Vertices = append(p.Vertices, v)
}

func (p *PathData) CurveTo(c1, c2, v Vertex) {
	p.Commands = append(p.Commands, PathCurveTo)
	p.Vertices = append(p.Vertices, c1, c2, v)
}
