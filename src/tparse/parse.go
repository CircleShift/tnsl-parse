package tparse

import (
	"bufio"
	"io"
	"os"
	"strings"
	"unicode"
	"unicode/utf8"
)

// Read in a number (may be a float)
func numericLiteral(r *bufio.Reader) Token {
	decimal := false
	run, _, err := r.ReadRune()

	b := strings.Builder{}

	for ; err == nil; run, _, err = r.ReadRune() {
		if (run == '.' || run == ',') && !decimal {
			decimal = true
		} else if !unicode.IsNumber(run) {
			break
		}
		b.WriteRune(run)
	}

	r.UnreadRune()

	return Token{Type: LITERAL, Data: b.String()}
}

// Parse a string (will escape \" only in this stage)
func stringLiteral(r *bufio.Reader) Token {
	escape := false
	run, _, err := r.ReadRune()

	if run != '"' {
		return Token{Type: LITERAL}
	}

	b := strings.Builder{}
	b.WriteRune(run)
	run, _, err = r.ReadRune()

	for ; err == nil; run, _, err = r.ReadRune() {
		b.WriteRune(run)
		if run == '\\' && !escape {
			escape = true
		} else if run == '"' && !escape {
			break
		}

	}

	return Token{Type: LITERAL, Data: b.String()}
}

// Parse a character in (escape \\ or \')
func charLiteral(r *bufio.Reader) Token {
	escape := false
	run, _, err := r.ReadRune()

	if run != '\'' {
		return Token{Type: LITERAL}
	}

	b := strings.Builder{}
	b.WriteRune(run)
	run, _, err = r.ReadRune()

	for ; err == nil; run, _, err = r.ReadRune() {
		b.WriteRune(run)
		if run == '\\' && !escape {
			escape = true
		} else if run == '\'' && !escape {
			break
		}

	}

	return Token{Type: LITERAL, Data: b.String()}
}

// Split reserved runes into rune groups
func splitResRunes(str string, max int) []Token {
	out := []Token{}

	rs := StringAsRunes(str)
	s, e := 0, max

	if max > len(rs) {
		e = len(rs)
	}

	for e <= len(rs) && s < len(rs) {
		if checkRuneGroup(RunesAsString(rs[s:e])) != -1 || e == s+1 {
			tmp := RunesAsString(rs[s:e])
			out = append(out, Token{Type: checkRuneGroup(tmp), Data: tmp})
			s = e
			if s+max < len(rs) {
				e = s + max
			} else {
				e = len(rs)
			}
		} else if e != s+1 {
			e--
		}
	}

	return out
}

// Remove block comments
func stripBlockComments(t []Token) []Token {
	out := []Token{}
	bc := false
	for _, tok := range t {
		if tok.Type == DELIMIT && tok.Data == "/#" {
			bc = true
			continue
		}

		if tok.Type == DELIMIT && tok.Data == "#/" {
			bc = false
			continue
		}

		if bc {
			continue
		}

		out = append(out, tok)
	}

	return out
}

// ParseFile tries to read a file and turn it into a series of tokens
func ParseFile(path string) []Token {
	out := []Token{}

	fd, err := os.Open(path)

	if err != nil {
		return out
	}

	read := bufio.NewReader(fd)

	b := strings.Builder{}

	max := maxResRunes()

	for r := rune(' '); ; r, _, err = read.ReadRune() {
		// If error in stream or EOF, break
		if err != nil {
			if err != io.EOF {
				out = append(out, Token{Type: -1})
			}
			break
		}

		// Checking for a space
		if unicode.IsSpace(r) {
			if b.String() != "" {
				out = append(out, Token{Type: checkToken(b.String()), Data: b.String()})
				b.Reset()
			}
			continue
		}

		if unicode.IsNumber(r) && b.String() == "" {
			read.UnreadRune()
			out = append(out, numericLiteral(read))

			continue
		}

		if r == '\'' {
			if b.String() != "" {
				out = append(out, Token{Type: checkToken(b.String()), Data: b.String()})
				b.Reset()
			}

			read.UnreadRune()
			out = append(out, charLiteral(read))

			continue
		}

		if r == '"' {
			if b.String() != "" {
				out = append(out, Token{Type: checkToken(b.String()), Data: b.String()})
				b.Reset()
			}

			read.UnreadRune()
			out = append(out, stringLiteral(read))

			continue
		}

		// Checking for a rune group
		if checkResRune(r) != -1 {
			if b.String() != "" {
				out = append(out, Token{Type: checkToken(b.String()), Data: b.String()})
				b.Reset()
			}

			for ; err == nil; r, _, err = read.ReadRune() {
				if checkResRune(r) == -1 {
					break
				}
				b.WriteRune(r)
			}

			read.UnreadRune()

			rgs := splitResRunes(b.String(), max)

			// Line Comments
			for i, rg := range rgs {
				if rg.Data == "#" {
					rgs = rgs[:i]
					read.ReadString('\n')
					break
				}
			}

			out = append(out, rgs...)

			b.Reset()

			continue
		}

		// Accumulate
		b.WriteRune(r)
	}

	return stripBlockComments(out)
}

// StringAsRunes returns a string as a rune slice
func StringAsRunes(s string) []rune {
	out := []rune{}
	for i, j := 0, 0; i < len(s); i += j {
		r, w := utf8.DecodeRuneInString(s[i:])
		out = append(out, r)
		j = w
	}
	return out
}

// BytesAsRunes returns a byte slice as a rune slice
func BytesAsRunes(b []byte) []rune {
	out := []rune{}
	for i, j := 0, 0; i < len(b); i += j {
		r, w := utf8.DecodeRune(b[i:])
		out = append(out, r)
		j = w
	}
	return out
}

// RunesAsString returns a string from a slice of runes
func RunesAsString(rs []rune) string {
	b := strings.Builder{}
	for _, r := range rs {
		b.WriteRune(r)
	}
	return b.String()
}
