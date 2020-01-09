package svg

import (
	"testing"
)

func TestParser(t *testing.T) {

	data := `
		<svg version="1.1" xmlns="http://www.w3.org/2000/svg" xmlns:xlink="http://www.w3.org/1999/xlink" id="svg-root" width="100%" height="100%" viewBox="0 0 480 360">
			<title id="test-title">color-prop-01-b</title>
			<desc id="test-desc">Test that viewer has the basic capability to process the color property</desc>
			<rect id="test-frame" x="1" y="1" width="478" height="358" fill="none" stroke="#000000"/>
		</svg>
		`

	/*doc := Svg{}

	err := xml.Unmarshal([]byte(data), &doc)
	if err != nil {
		t.Error(err)
		return
	}
	t.Logf("ver: %s", doc.Version)*/

	element, _ := parse([]byte(data))

	if element == nil {
		return
	}
}
