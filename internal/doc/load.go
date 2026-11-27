package doc

import (
	"strings"

	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/ast"
	"github.com/yuin/goldmark/text"
	"github.com/yuin/goldmark/util"
)

// clean turns markdown escapes and entities into the real characters.
func clean(b []byte) []byte {
	b = util.UnescapePunctuations(b)     // \*  → *
	b = util.ResolveNumericReferences(b) // &#233; → é
	b = util.ResolveEntityNames(b)       // &amp; → &
	return b
}

// inline returns the text inside block without markdown symbols, and the
// bold, italic and code ranges inside that text. Line breaks become spaces.
func inline(block ast.Node, src []byte) (string, []Span) {
	var b strings.Builder
	var spans []Span
	var stack []Style
	cur := Plain

	ast.Walk(block, func(n ast.Node, entering bool) (ast.WalkStatus, error) {
		if n == block {
			return ast.WalkContinue, nil
		}

		if t, ok := n.(*ast.Text); ok && entering {
			start := b.Len()
			b.Write(clean(t.Segment.Value(src)))

			if t.SoftLineBreak() || t.HardLineBreak() {
				b.WriteByte(' ')
			}

			if cur != Plain {
				spans = appendSpan(spans, Span{
					Start: start,
					End:   b.Len(),
					Style: cur,
				})
			}
			return ast.WalkContinue, nil
		}
		add := nodeStyle(n)

		if add == Plain {
			return ast.WalkContinue, nil
		}

		if entering {
			stack = append(stack, cur)
			cur |= add
		} else {
			cur = stack[len(stack)-1]    //get the last element
			stack = stack[:len(stack)-1] //pop the element
		}

		return ast.WalkContinue, nil

	})
	return b.String(), spans

}

// Load reads a markdown document into a Doc.
func Load(src []byte) Doc {
	root := goldmark.DefaultParser().Parse(text.NewReader(src))
	var d Doc

	add := func(kind Kind, level int, block ast.Node) {
		text, spans := inline(block, src)
		d.Paras = append(d.Paras, Para{
			Kind:  kind,
			Level: level,
			Text:  text,
			Spans: spans,
		})
	}

	for n := root.FirstChild(); n != nil; n = n.NextSibling() {
		switch v := n.(type) {
		case (*ast.Heading):
			add(KindHeading, v.Level, v)
		case (*ast.Paragraph):
			add(KindBody, 0, v)
		case (*ast.Blockquote):
			for c := v.FirstChild(); c != nil; c = c.NextSibling() {
				add(KindQuote, 0, c)
			}
		case (*ast.List):
			for c := v.FirstChild(); c != nil; c = c.NextSibling() {
				add(KindListItem, 0, c)
			}
		}
	}
	return d
}
