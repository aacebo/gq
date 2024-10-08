package query

type _Parser struct {
	curr    *_Token
	prev    *_Token
	errs    []error
	scanner *_Scanner
}

func Parser(src []byte) *_Parser {
	self := &_Parser{
		curr:    nil,
		prev:    nil,
		errs:    []error{},
		scanner: newScanner(src),
	}

	self.next()
	return self
}

func (self *_Parser) Parse() (Query, error) {
	query := Query{
		Args:   QueryArgs{},
		Fields: map[string]Query{},
	}

	if _, err := self.consume(LEFT_BRACE, "expected '{'"); err != nil {
		return query, err
	}

	for {
		if self.curr.Kind == EOF {
			break
		}

		name, subQuery, err := self.parseField()

		if err != nil {
			return query, err
		}

		query.Fields[name] = subQuery

		if !self.match(COMMA) {
			break
		}
	}

	return query, nil
}

func (self *_Parser) parseField() (string, Query, error) {
	query := Query{
		Args:   QueryArgs{},
		Fields: map[string]Query{},
	}

	name, err := self.consume(IDENTIFIER, "expected field name")

	if err != nil {
		return "", query, err
	}

	if self.match(LEFT_PAREN) {
		query.Args, err = self.parseArgs()

		if err != nil {
			return "", query, err
		}
	}

	if self.curr.Kind == LEFT_BRACE {
		q, err := self.Parse()

		if err != nil {
			return "", query, err
		}

		query.Fields = q.Fields
	}

	return name.String(), query, nil
}

func (self *_Parser) parseArgs() (QueryArgs, error) {
	args := QueryArgs{}

	for self.curr.Kind != RIGHT_PAREN {
		param, err := self.consume(IDENTIFIER, "expected parameter name")

		if err != nil {
			return args, err
		}

		if _, err := self.consume(COLON, "expected ':'"); err != nil {
			return args, err
		}

		if self.match(STRING) {
			args[param.String()] = self.prev.String()
		} else if self.match(INT) {
			args[param.String()], err = self.prev.Int()
		} else if self.match(FLOAT) {
			args[param.String()], err = self.prev.Float()
		} else if self.match(BOOL) {
			args[param.String()], err = self.prev.Bool()
		} else if self.match(BYTE) {
			args[param.String()] = self.prev.Byte()
		} else if self.match(NULL) {
			args[param.String()] = nil
		} else {
			return args, newError(
				self.curr.Ln,
				self.curr.Start,
				self.curr.End,
				"invalid parameter value",
			)
		}

		if err != nil {
			return args, err
		}

		if self.match(COMMA) {
			continue
		}
	}

	self.next()
	return args, nil
}

func (self *_Parser) next() bool {
	self.prev = self.curr
	t, err := self.scanner.Next()

	if err != nil {
		self.errs = append(self.errs, err)
		return self.next()
	}

	self.curr = &t

	if t.Kind == EOF {
		return false
	}

	return true
}

func (self *_Parser) match(kind _TokenKind) bool {
	if self.curr.Kind != kind {
		return false
	}

	self.next()
	return true
}

func (self *_Parser) consume(kind _TokenKind, message string) (*_Token, error) {
	if self.curr.Kind == kind {
		self.next()
		return self.prev, nil
	}

	return nil, newError(
		self.curr.Ln,
		self.curr.Start,
		self.curr.End,
		message,
	)
}
