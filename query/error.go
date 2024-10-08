package query

import "fmt"

type _Error struct {
	Ln      int
	Start   int
	End     int
	Message string
}

func newError(ln int, start int, end int, message string) _Error {
	return _Error{
		Ln:      ln,
		Start:   start,
		End:     end,
		Message: message,
	}
}

func (self _Error) Error() string {
	return self.String()
}

func (self _Error) String() string {
	return fmt.Sprintf(
		"[ln: %d, start: %d, end: %d] => %s",
		self.Ln,
		self.Start,
		self.End,
		self.Message,
	)
}
