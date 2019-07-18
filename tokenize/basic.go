package tokenize

import (
	"unicode"

	"strings"

	"golang.org/x/text/unicode/norm"
)

// Basic is a BasicTokenizer that runs Runs basic tokenize (punctuation splitting, lower casing, etc.).
type Basic struct {
	// Lower will apply a lower case filter to input
	Lower bool
}

// NewBasic returns a basic tokenizer. Method is supplied to match constructor of other tokenizers
func NewBasic() Basic {
	return Basic{Lower: true}
}

// Tokenize will segment a texxt into individual tokens. Follows algorithm from ref-imp
// Clean, PadChinese, Whitespace Split, Lower?, SplitPunc, Whitespace Split
func (bt Basic) Tokenize(text string) []string {
	// TODO assert text is unicode
	// text = unicode(text), from python impl
	text = clean(text)
	text = padChinese(text)
	var toks []string
	for _, tok := range tokenizeWhitespace(text) {
		if bt.Lower {
			tok = strings.ToLower(tok)
			tok = stripAccents(tok) // why does lower strip accents?
		}
		toks = append(toks, splitPunc(tok)...)
	}
	toks = tokenizeWhitespace(strings.Join(toks, " "))
	return toks
}

func clean(text string) string {
	var b strings.Builder
	for _, c := range text {
		if c == 0 || c == 0xfffd || isControl(c) {
			continue
		} else if isWhitespace(c) {
			b.WriteRune(' ')
		} else {
			b.WriteRune(c)
		}
	}
	return b.String()
}

func stripAccents(text string) string {
	// TODO test
	var b strings.Builder
	for _, c := range norm.NFD.String(text) {
		if !unicode.Is(unicode.Mn, c) {
			b.WriteRune(c)
		}
	}
	return b.String()
}

func splitPunc(text string) []string {
	// TODO test
	var toks []string
	var b strings.Builder
	for _, c := range text {
		if isPunctuation(c) {
			toks = append(toks, b.String())
			toks = append(toks, string(c))
			b.Reset()
		} else {
			b.WriteRune(c)
		}
	}
	if b.Len() > 0 {
		toks = append(toks, b.String())
	}
	return toks
}

//tokenizeWhitespace splits text into tokeens by whitespace, per python semantics empty strings are not included
func tokenizeWhitespace(text string) []string {
	split := strings.Split(text, " ")
	var toks []string
	for _, tok := range split {
		if tok != "" {
			toks = append(toks, tok)
		}
	}
	return toks

}

//padChinese will add space padding around all CJK chars
// This implementation matches BasicTokenizer._tokenize_chinese_chars
func padChinese(text string) string {
	var b strings.Builder
	for _, c := range text {
		if isChinese(c) {
			b.WriteRune(' ')
			b.WriteRune(c)
			b.WriteRune(' ')
		} else {
			b.WriteRune(c)
		}
	}
	return b.String()
}
