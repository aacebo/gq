package query

import (
	"slices"
)

type _Scanner struct {
	src   []byte
	ln    int
	left  int
	right int
}

func newScanner(src []byte) *_Scanner {
	return &_Scanner{
		src:   src,
		ln:    0,
		left:  0,
		right: 0,
	}
}

func (self *_Scanner) Next() (_Token, error) {
	if self.right >= len(self.src) {
		return self.create(EOF), nil
	}

	self.left = self.right
	b := self.src[self.right]
	self.right++

	switch b {
	case ' ':
	case '\r':
	case '\t':
		// ignore whitespace
		break
	case '\n':
		self.ln++
		break
	case '(':
		return self.create(LEFT_PAREN), nil
	case ')':
		return self.create(RIGHT_PAREN), nil
	case '{':
		return self.create(LEFT_BRACE), nil
	case '}':
		return self.create(RIGHT_BRACE), nil
	case ',':
		return self.create(COMMA), nil
	case ':':
		return self.create(COLON), nil
	case '\'':
		return self.onByte()
	case '"':
		return self.onString()
	default:
		if self.isInt(b) {
			return self.onNumeric()
		} else if self.isAlpha(b) {
			return self.onIdentifier()
		}

		return _Token{}, self.error("unexpected character")
	}

	return self.Next()
}

func (self *_Scanner) onByte() (_Token, error) {
	self.right++

	if self.peek() != '\'' {
		return _Token{}, self.error("unterminated byte")
	}

	self.left++
	token := self.create(BYTE)
	self.right++
	return token, nil
}

func (self *_Scanner) onString() (_Token, error) {
	for self.peek() != '"' && self.peek() != 0 {
		if self.peek() == '\n' {
			self.ln++
		} else if self.peek() == '\\' {
			err := self.onEscape()

			if err != nil {
				return _Token{}, err
			}
		}

		self.right++
	}

	if self.right == len(self.src) {
		return _Token{}, self.error("unterminated string")
	}

	self.left++
	token := self.create(STRING)
	self.right++
	return token, nil
}

func (self *_Scanner) onEscape() error {
	self.right++

	defer func() {
		self.right--
	}()

	switch self.peek() {
	case 'a': // bell
		self.src = slices.Replace(self.src, self.right-1, self.right+1, '\a')
	case 'b': // backspace
		self.src = slices.Replace(self.src, self.right-1, self.right+1, '\b')
	case 'f': // form feed
		self.src = slices.Replace(self.src, self.right-1, self.right+1, '\f')
	case 'n': // new line
		self.src = slices.Replace(self.src, self.right-1, self.right+1, '\n')
	case 'r': // carriage return
		self.src = slices.Replace(self.src, self.right-1, self.right+1, '\r')
	case 't': // horizontal tab
		self.src = slices.Replace(self.src, self.right-1, self.right+1, '\t')
	case 'v': // verical tab
		self.src = slices.Replace(self.src, self.right-1, self.right+1, '\v')
	case '\'': // single quote
		self.src = slices.Replace(self.src, self.right-1, self.right+1, '\'')
	case '"': // double quote
		self.src = slices.Replace(self.src, self.right-1, self.right+1, '"')
	case '\\': // back slash
		self.src = slices.Replace(self.src, self.right-1, self.right+1, '\\')
	default:
		return self.error("unknown escape sequence")
	}

	return nil
}

func (self *_Scanner) onNumeric() (_Token, error) {
	kind := INT

	for self.isInt(self.peek()) {
		self.right++
	}

	if self.peek() == '.' {
		kind = FLOAT
		self.right++

		for self.isInt(self.peek()) {
			self.right++
		}
	}

	return self.create(kind), nil
}

func (self *_Scanner) onIdentifier() (_Token, error) {
	for self.isAlpha(self.peek()) || self.isInt(self.peek()) {
		self.right++
	}

	name := self.src[self.left:self.right]

	if kind, ok := Keywords[string(name)]; ok {
		return self.create(kind), nil
	}

	return self.create(IDENTIFIER), nil
}

func (self _Scanner) peek() byte {
	if self.right >= len(self.src) {
		return 0
	}

	return self.src[self.right]
}

func (self _Scanner) isInt(b byte) bool {
	return b >= '0' && b <= '9'
}

func (self _Scanner) isAlpha(b byte) bool {
	return (b >= 'a' && b <= 'z') ||
		(b >= 'A' && b <= 'Z') ||
		(b == '_')
}

func (self _Scanner) create(kind _TokenKind) _Token {
	return newToken(
		kind,
		self.ln,
		self.left,
		self.right,
		self.src[self.left:self.right],
	)
}

func (self _Scanner) error(message string) error {
	return newError(
		self.ln,
		self.left,
		self.right,
		message,
	)
}
