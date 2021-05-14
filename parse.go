package svg

type element struct {
	name       string
	attributes map[string]string
	children   []*element
}

func (e *element) Attr(name string) (v string, exists bool) {
	v, exists = e.attributes[name]
	return
}

func (e *element) ForEachChildNode(callback func(tag string, ch sourcer) error) error {
	if callback == nil {
		return nil
	}
	for _, c := range e.children {
		err := callback(c.name, c)
		if err != nil {
			return err
		}
	}
	return nil
}
