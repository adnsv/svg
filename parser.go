package svg

import (
	"bytes"
	"encoding/xml"
	"strings"
)

type element struct {
	name       string
	attributes map[string]string
	children   []*element
	content    string
}

func newElement(token xml.StartElement) *element {
	element := &element{}
	attributes := make(map[string]string)
	for _, attr := range token.Attr {
		attributes[attr.Name.Local] = attr.Value
	}
	element.name = token.Name.Local
	element.attributes = attributes
	return element
}

func decodeFirst(decoder *xml.Decoder) (*element, error) {
	for {
		token, err := decoder.Token()
		if token == nil {
			break
		}
		if err != nil {
			return nil, err
		}

		switch element := token.(type) {
		case xml.StartElement:
			return newElement(element), nil
		}
	}
	return &element{}, nil
}

func (e *element) find(id string) *element {
	for _, child := range e.children {
		if childId, ok := child.attributes["id"]; ok && childId == id {
			return child
		}
		if element := child.find(id); element != nil {
			return element
		}
	}
	return nil
}

func (e *element) decode(decoder *xml.Decoder) error {
	for {
		token, err := decoder.Token()
		if token == nil {
			break
		}

		if err != nil {
			return err
		}

		switch element := token.(type) {
		case xml.StartElement:
			nextElement := newElement(element)
			err := nextElement.decode(decoder)
			if err != nil {
				return err
			}

			e.children = append(e.children, nextElement)

		case xml.CharData:
			data := strings.TrimSpace(string(element))
			if data != "" {
				e.content = string(element)
			}

		case xml.EndElement:
			if element.Name.Local == e.name {
				return nil
			}
		}
	}
	return nil
}

func parse(raw []byte) (*element, error) {
	decoder := xml.NewDecoder(bytes.NewReader(raw))
	element, err := decodeFirst(decoder)
	if err != nil {
		return nil, err
	}
	if err := element.decode(decoder); err != nil {
		return nil, err
	}
	return element, nil
}