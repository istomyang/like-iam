package util

import "github.com/MakeNowJust/heredoc/v2"

// Normalize normalizes the raw string on different demands.
type Normalize struct {
	raw string
}

func NewNormalize(raw string) *Normalize {
	return &Normalize{raw: raw}
}

// Heredoc do like that:
//
// doc := heredoc.Doc(`
//
//	Foo
//	Bar
//
// `)
// Output: "Foo\nBar\n"
func (n *Normalize) Heredoc() *Normalize {
	n.raw = heredoc.Doc(n.raw)
	return n
}

func (n *Normalize) String() string {
	return n.raw
}
