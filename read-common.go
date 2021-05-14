package svg

type sourcer interface {
	Attr(name string) (v string, exists bool)
	ForEachChildNode(callback func(tag string, ch sourcer) error) error
}

type reader interface {
	Item
	read(src sourcer) error
}
