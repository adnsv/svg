package svg

func attrLength(src sourcer, name string) (*Length, error) {
	attr, exists := src.Attr(name)
	if !exists {
		return nil, nil
	}
	ret := &Length{}
	err := ret.Unmarshal(attr)
	return ret, err
}

func attrCoordinate(src sourcer, name string) (*Coordinate, error) {
	attr, exists := src.Attr(name)
	if !exists {
		return nil, nil
	}
	ret := &Coordinate{}
	err := ret.Unmarshal(attr)
	return ret, err
}

/*
func attrLen(aa xg.AttributeList, name string) (*Length, error) {
	attr, exists := aa.Attr(name)
	if !exists {
		return nil, nil
	}
	ret := &Length{}
	err := ret.Unmarshal(attr)
	return ret, err
}

func attrCoord(aa xg.AttributeList, name string) (*Coordinate, error) {
	attr, exists := aa.Attr(name)
	if !exists {
		return nil, nil
	}
	ret := &Coordinate{}
	err := ret.Unmarshal(attr)
	return ret, err
}
*/
