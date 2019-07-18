package svg

func attrLength(elt *element, name string) (*Length, error) {
	attr, exists := elt.attributes[name]
	if !exists {
		return nil, nil
	}
	ret := &Length{}
	err := ret.Unmarshal(attr)
	return ret, err
}

func attrCoordinate(elt *element, name string) (*Coordinate, error) {
	attr, exists := elt.attributes[name]
	if !exists {
		return nil, nil
	}
	ret := &Coordinate{}
	err := ret.Unmarshal(attr)
	return ret, err
}
