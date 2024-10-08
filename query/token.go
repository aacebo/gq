package query

import "strconv"

type _Token struct {
	Kind  _TokenKind
	Ln    int
	Start int
	End   int
	Value []byte
}

func newToken(kind _TokenKind, ln int, start int, end int, value []byte) _Token {
	return _Token{
		Kind:  kind,
		Ln:    ln,
		Start: start,
		End:   end,
		Value: value,
	}
}

func (self _Token) Byte() byte {
	return self.Value[0]
}

func (self _Token) String() string {
	return string(self.Value)
}

func (self _Token) Int() (int, error) {
	return strconv.Atoi(string(self.Value))
}

func (self _Token) Float() (float64, error) {
	return strconv.ParseFloat(string(self.Value), 64)
}

func (self _Token) Bool() (bool, error) {
	return strconv.ParseBool(string(self.Value))
}
