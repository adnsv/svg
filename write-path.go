package svg

import "fmt"

func (pd *PathData) String() string {
	v := 0
	s := ""
	for _, c := range pd.Commands {
		switch c {
		case PathClose:
			s += "z"
		case PathMoveTo:
			s += fmt.Sprintf("M%g,%g", pd.Vertices[v].X, pd.Vertices[v].Y)
			v++
		case PathLineTo:
			s += fmt.Sprintf("L%g,%g", pd.Vertices[v].X, pd.Vertices[v].Y)
			v++
		case PathCurveTo:
			s += fmt.Sprintf("C%g,%g,%g,%g,%g,%g", pd.Vertices[v].X, pd.Vertices[v].Y, pd.Vertices[v+1].X, pd.Vertices[v+1].Y, pd.Vertices[v+2].X, pd.Vertices[v+2].Y)
			v += 3
		}
	}
	return s
}
