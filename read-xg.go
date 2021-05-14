package svg

import (
	"errors"

	xg "github.com/adnsv/xmlgo"
)

func ParseXG(in string) (*Svg, error) {
	content := xg.Open(in)
	if !content.NextTag() {
		err := content.Err()
		if err == nil {
			err = errors.New("invalid file content")
		}
		return nil, err
	}
	if content.Name() != "svg" {
		return nil, content.MakeError("", "root tag must be 'svg'")
	}
	s := &Svg{}
	content.HandleTag(func(aa xg.AttributeList, cc *xg.Content) error {
		return s.read(&xgsourcer{aa, cc})
	})
	if content.Err() != nil {
		return nil, content.Err()
	}

	return s, nil
}

type xgsourcer struct {
	aa xg.AttributeList
	cc *xg.Content
}

func (x *xgsourcer) Attr(name string) (v string, exists bool) {
	return x.aa.Attr(name)
}

func (x *xgsourcer) ForEachChildNode(callback func(tag string, ch sourcer) error) error {
	if x.cc == nil {
		return nil
	}
	for x.cc.NextTag() {
		n := string(x.cc.Name())
		x.cc.HandleTag(func(aa xg.AttributeList, cc *xg.Content) error {
			return callback(n, &xgsourcer{aa, cc})
		})
		if x.cc.Err() != nil {
			return x.cc.Err()
		}
	}
	return nil
}
