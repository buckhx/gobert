package tokenize

type Full struct {
	Basic     Basic
	Wordpiece Wordpiece
}

func (f Full) Tokenize(text string) []string {
	var toks []string
	for _, tok := range f.Basic.Tokenize(text) {
		toks = append(toks, f.Wordpiece.Tokenize(tok)...)
	}
	return toks
}
