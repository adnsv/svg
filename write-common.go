package svg

type targeter interface {
	Attr(name, value string)
	Child(tag string, callback func(tgt targeter))
}

type writer interface {
	write(tgt targeter)
}
