package render

import (
	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/text"
)

func Dump(src []byte) string {
	doc := goldmark.DefaultParser().Parse(text.NewReader(src))
	doc.Dump(src, 0)
	return ""
}
