package svg

import (
	"encoding/xml"
)

type Item struct {
	ID string `xml:"id,attr,omitempty"`
}

type Node struct {
	Items []*Item
}

type Shape struct {
	Item
	Fill           *Paint     `xml:"fill,attr,omitempty"`
	FillRule       *FillRule  `xml:"fill-rule,attr,omitempty"`
	FillOpacity    *float64   `xml:"fill-opacity,attr,omitempty"`
	Stroke         *Paint     `xml:"stroke,attr,omitempty"`
	StrokeWidth    *Length    `xml:"stroke-width,attr,omitempty"`
	StrokeOpacity  *float64   `xml:"stroke-opacity,attr,omitempty"`
	StrokeLineCap  *LineCap   `xml:"stroke-linecap,attr,omitempty"`
	StrokeLineJoin *LineJoin  `xml:"stroke-linejoin,attr,omitempty"`
	Opacity        *float64   `xml:"opacity,attr,omitempty"`
	Transform      *Transform `xml:"transform,attr,omitempty"`
}

type Group struct {
	XMLName xml.Name `xml:"g"`
	Node
	Transform *Transform `xml:"transform,attr,omitempty"`
	Items     []*Item
}

type Rect struct {
	XMLName xml.Name `xml:"shape"`
	Shape
	X      *Coordinate `xml:"x,attr,omitempty"`
	Y      *Coordinate `xml:"y,attr,omitempty"`
	Width  *Length     `xml:"width,attr,omitempty"`
	Height *Length     `xml:"height,attr,omitempty"`
	Rx     *Length     `xml:"rx,attr,omitempty"`
	Ry     *Length     `xml:"ry,attr,omitempty"`
}

// Svg implements SVG <svg> element
type Svg struct {
	XMLName xml.Name `xml:"svg"`
	Group
	Version     string      `xml:"version,attr,omitempty"`
	BaseProfile string      `xml:"baseProfile,attr,omitempty"`
	ViewBox     string      `xml:"viewBox,attr,omitempty"`
	X           *Coordinate `xml:"x,attr,omitempty"`
	Y           *Coordinate `xml:"y,attr,omitempty"`
	Width       *Length     `xml:"width,attr,omitempty"`
	Height      *Length     `xml:"height,attr,omitempty"`
}
