package svg

func readOptLength(src sourcer, name string) (*Length, error) {
	attr, exists := src.Attr(name)
	if !exists {
		return nil, nil
	}
	ret := &Length{}
	err := ret.Unmarshal(attr)
	return ret, err
}

func readOptCoord(src sourcer, name string) (*Coordinate, error) {
	attr, exists := src.Attr(name)
	if !exists {
		return nil, nil
	}
	ret := &Coordinate{}
	err := ret.Unmarshal(attr)
	return ret, err
}
