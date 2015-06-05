// This package provides Ruby's StringScanner-like functions.
package scanner

import "regexp"

type Scanner struct {
	input string
	pos   int
}

func NewScanner(input string) *Scanner {
	return &Scanner{
		input: input,
		pos:   0,
	}
}

func (scn *Scanner) Scan(re *regexp.Regexp) (rtn string, ok bool) {
	if scn.IsEos() {
		return "", false
	}

	strForScan := scn.input[scn.pos:]
	loc := re.FindStringIndex(strForScan)

	if loc == nil || loc[0] != 0 {
		return "", false
	}

	rtn = scn.input[scn.pos+loc[0] : scn.pos+loc[1]]
	ok = true
	scn.pos = loc[1] + scn.pos
	return
}

func (scn *Scanner) ScanUntil(re *regexp.Regexp) (rtn string, ok bool) {
	if scn.IsEos() {
		return "", false
	}

	strForScan := scn.input[scn.pos:]
	loc := re.FindStringIndex(strForScan)

	if loc == nil {
		return "", false
	}

	rtn = scn.input[scn.pos : scn.pos+loc[1]]
	ok = true
	scn.pos = loc[1] + scn.pos
	return
}

func (scn *Scanner) SkipUntil(re *regexp.Regexp) (rtn int, ok bool) {
	rtnStr, ok := scn.ScanUntil(re)
	if !ok {
		return 0, false
	}
	return len(rtnStr), true
}

func (scn *Scanner) Getch() (rtn string, ok bool) {
	if scn.IsEos() {
		return "", false
	}

	rtn = scn.input[scn.pos : scn.pos+1]
	ok = true
	scn.pos = scn.pos + 1
	return
}

func (scn *Scanner) IsEos() bool {
	return scn.pos >= len(scn.input)
}
func (scn *Scanner) IsBol() bool {

	if scn.pos == 0 {
		return true
	} else if scn.pos >= 2 && scn.input[scn.pos-1:scn.pos] == "\n" {
		return true
	} else {
		return false
	}
}

func (scn *Scanner) Peek(length int) string {
	if scn.IsEos() {
		return ""
	} else {
		end := scn.pos + length

		if end >= len(scn.input) {
			end = len(scn.input)
		}

		return scn.input[scn.pos:end]
	}
}

func (scn *Scanner) Terminate() {
	scn.pos = len(scn.input)
}
