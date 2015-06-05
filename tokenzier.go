package mockingbird

import (
	"regexp"
	"strings"

	"github.com/lazywei/mockingbird/scanner"
)

func multiLineCmtPairs() map[string]string {
	return map[string]string{}
}

var (
	reShebang   = regexp.MustCompile(`^#!.+`)
	reEndOfLine = regexp.MustCompile(`\n|\z`)

	reSingleLineComment = regexp.MustCompile(`\s*\/\/ |\s*\-\- |\s*# |\s*% |\s*" `)
	reMultiLineComment  = regexp.MustCompile(`\/\*|<!--|{-|\(\*|"""|'''`)

	reMultiLineCommentPairs = map[string]*regexp.Regexp{
		`/*`:   regexp.MustCompile(`\*\/`), // C
		`<!--`: regexp.MustCompile(`-->`),  // XML
		`{-`:   regexp.MustCompile(`-}`),   // Haskell
		`(*`:   regexp.MustCompile(`\*\)`), // Coq
		`"""`:  regexp.MustCompile(`"""`),  // Python
		`'''`:  regexp.MustCompile(`'''`),  // Python
	}

	reQuote    = regexp.MustCompile(`"`)
	reQuoteEnd = regexp.MustCompile(`[^\\]"`)

	reSQuote    = regexp.MustCompile(`'`)
	reSQuoteEnd = regexp.MustCompile(`[^\\]'`)

	reNumberLiteral = regexp.MustCompile(`(0x)?\d(\d|\.)*`)
	reSgml          = regexp.MustCompile(`<[^\s<>][^<>]*>`)
	rePunctuation   = regexp.MustCompile(`;|\{|\}|\(|\)|\[|\]`)
	reRegularToken  = regexp.MustCompile(`[\w\.@#\/\*]+`)
	reOperators     = regexp.MustCompile(`<<?|\+|\-|\*|\/|%|&&?|\|\|?`)

	reMLTag    = regexp.MustCompile(`<\/?[^\s>]+`)
	reMLAssign = regexp.MustCompile(`\w+=`)

	reIdentifier = regexp.MustCompile(`\w+`)
)

func ExtractTokens(data string) []string {
	s := scanner.NewScanner(data)
	tokens := []string{}

	for s.IsEos() != true {
		if token, ok := s.Scan(reShebang); ok {
			name, ok := extractShebang(token)
			if ok {
				tokens = append(tokens, "SHEBANG#!"+name)
			}
		} else if s.IsBol() && scanOrNot(s, reSingleLineComment) {

			s.SkipUntil(reEndOfLine)

		} else if startToken, ok := s.Scan(reMultiLineComment); ok {

			closeToken := reMultiLineCommentPairs[startToken]
			s.SkipUntil(closeToken)

		} else if scanOrNot(s, reQuote) {
			if s.Peek(1) == `"` {
				s.Getch()
			} else {
				s.ScanUntil(reQuoteEnd)
			}

		} else if scanOrNot(s, reSQuote) {

			if s.Peek(1) == `'` {
				s.Getch()
			} else {
				s.ScanUntil(reSQuoteEnd)
			}

		} else if scanOrNot(s, reNumberLiteral) {
			// Skip number literals

		} else if rtn, ok := s.Scan(reSgml); ok {

			for _, tkn := range extractSgmlTokens(rtn) {
				tokens = append(tokens, tkn)
			}

		} else if rtn, ok := s.Scan(rePunctuation); ok {
			// Common programming punctuation
			tokens = append(tokens, rtn)

		} else if rtn, ok := s.Scan(reRegularToken); ok {
			// Regular token
			tokens = append(tokens, rtn)

		} else if rtn, ok := s.Scan(reOperators); ok {
			// Common operators
			tokens = append(tokens, rtn)

		} else {
			s.Getch()
		}
	}

	return tokens
}

// This function is silly ... silly Go ...
func scanOrNot(s *scanner.Scanner, re *regexp.Regexp) bool {
	_, ok := s.Scan(re)
	return ok
}

func extractSgmlTokens(token string) []string {
	s := scanner.NewScanner(token)
	tokens := []string{}

	for s.IsEos() != true {
		if token, ok := s.Scan(reMLTag); ok {
			tokens = append(tokens, token+">")
		} else if token, ok := s.Scan(reMLAssign); ok {
			tokens = append(tokens, token)

			// Then skip over attribute value

			if scanOrNot(s, reQuote) {
				s.SkipUntil(reQuoteEnd)
			} else if scanOrNot(s, reSQuote) {
				s.SkipUntil(reSQuoteEnd)
			} else {
				s.SkipUntil(reIdentifier)
			}
		} else if token, ok := s.Scan(reIdentifier); ok {
			tokens = append(tokens, token)
		} else if scanOrNot(s, reIdentifier) {
			s.Terminate()
		} else {
			s.Getch()
		}
	}

	return tokens
}

var (
	reShebangContent = regexp.MustCompile(`^#!\s*\S+`)
	reShebangSpace   = regexp.MustCompile(`\s+`)
	reShebangNSpace  = regexp.MustCompile(`\S+`)
	reShebangAssign  = regexp.MustCompile(`.*=[^\s]+\s+`)
	reShebangName    = regexp.MustCompile(`[^\d]+`)
)

func extractShebang(token string) (string, bool) {
	s := scanner.NewScanner(token)

	path, ok := s.Scan(reShebangContent)

	if !ok {
		return "", false
	}

	paths := strings.Split(path, `/`)

	if len(paths) == 0 {
		return "", false
	}

	name := paths[len(paths)-1]

	if name == `env` {
		s.Scan(reShebangSpace)
		s.Scan(reShebangAssign)
		name, ok = s.Scan(reShebangNSpace)
	}

	if ok {
		name = reShebangName.FindString(name)
		return name, true
	} else {
		return "", false
	}
}
