package render

import (
	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/ast"
	"github.com/yuin/goldmark/text"
)

type Span struct {
	Start int
	End   int
	Style Style
}

func textSpans(src []byte) []Span {
	doc := goldmark.DefaultParser().Parse(text.NewReader(src))
	out := []Span{}
	stack := []Style{}
	cur := Plain

	ast.Walk(doc, func(n ast.Node, entering bool) (ast.WalkStatus, error) {

		if t, ok := n.(*ast.Text); ok && entering {
			out = append(out, Span{
				Start: t.Segment.Start,
				End:   t.Segment.Stop,
				Style: cur,
			})
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
	return out

}

func nodeStyle(n ast.Node) Style {
	switch v := n.(type) {
	case *ast.Heading:
		return Heading
	case *ast.Blockquote:
		return Quote
	case *ast.CodeSpan:
		return Code
	case *ast.Emphasis:
		if v.Level == 2 {
			return Bold
		}
		return Italic
	}
	return Plain
}

// appendSpan adds s to the end of spans, keeping the list clean:
// empty spans are ignored, and a span with the same style as the last
// one, starting right where it ends, is merged into it.
func appendSpan(spans []Span, s Span) []Span {

	if s.End <= s.Start {
		return spans
	}

	if len(spans) > 0 {
		last := len(spans) - 1
		if spans[last].Style == s.Style && spans[last].End == s.Start {
			spans[last].End = s.End
			return spans
		}
	}
	return append(spans, s)
}
