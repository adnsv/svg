package svg

import (
	"io"

	xg "github.com/adnsv/xmlgo"
)

type xgwriter struct {
	out *xg.Writer
}

func (x *xgwriter) Attr(k, v string) {
	x.out.OptStringAttr(k, v)
}

func (x *xgwriter) Child(tag string, callback func(tgt targeter)) {
	x.out.OTag(tag)
	callback(x)
	x.out.CTag()
}

func WriteXG(w io.Writer, s *Svg) {
	xgw := xgwriter{out: xg.NewWriter(w)}
	xgw.Child("svg", func(tgt targeter) { s.write(tgt) })
}
