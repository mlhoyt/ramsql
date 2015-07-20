package parser

import (
	"fmt"
)

func (p *parser) parseSelect(tokens []Token) (*Instruction, error) {
	i := &Instruction{}
	var err error

	// Create select decl
	selectDecl := NewDecl(tokens[p.index])
	i.Decls = append(i.Decls, selectDecl)

	// After select token, should be either
	// a StarToken
	// a list of table names + (StarToken Or Attribute)
	// a builtin func (COUNT, MAX, ...)
	if err = p.next(); err != nil {
		return nil, fmt.Errorf("SELECT token must be followed by attributes to select")
	}

	for {
		if p.is(CountToken) {
			attrDecl, err := p.parseBuiltinFunc()
			if err != nil {
				return nil, err
			}
			selectDecl.Add(attrDecl)
		} else {
			attrDecl, err := p.parseAttribute()
			if err != nil {
				return nil, err
			}
			selectDecl.Add(attrDecl)
		}

		// If comma, loop again.
		if p.is(CommaToken) {
			if err := p.next(); err != nil {
				return nil, err
			}
			continue
		}
		break
	}

	// Must be from now
	if tokens[p.index].Token != FromToken {
		return nil, fmt.Errorf("Syntax error near %v\n", tokens[p.index])
	}
	fromDecl := NewDecl(tokens[p.index])
	selectDecl.Add(fromDecl)

	// Now must be a list of table
	for {
		// string
		if err = p.next(); err != nil {
			return nil, fmt.Errorf("Unexpected end. Syntax error near %v\n", tokens[p.index])
		}
		if tokens[p.index].Token != StringToken {
			return nil, p.syntaxError()
		}
		tableNameDecl := NewDecl(tokens[p.index])
		fromDecl.Add(tableNameDecl)

		// If no next, then it's implicit where
		if err = p.next(); err != nil {
			addImplicitWhereAll(selectDecl)
			return i, nil
		}
		// if not comma, break
		if tokens[p.index].Token != CommaToken {
			break // No more table
		}
	}

	// JOIN OR ...?
	if p.is(JoinToken) {
		joinDecl, err := p.parseJoin()
		if err != nil {
			return nil, err
		}
		selectDecl.Add(joinDecl)
	}

	switch p.cur().Token {
	case WhereToken:
		err := p.parseWhere(selectDecl)
		if err != nil {
			return nil, err
		}
	case OrderToken:
		// WHERE clause is implicit
		addImplicitWhereAll(selectDecl)
		err := p.parseOrderBy(selectDecl)
		if err != nil {
			return nil, err
		}
	}

	return i, nil
}

func addImplicitWhereAll(decl *Decl) {

	whereDecl := &Decl{
		Token:  WhereToken,
		Lexeme: "where",
	}
	whereDecl.Add(&Decl{
		Token:  NumberToken,
		Lexeme: "1",
	})

	decl.Add(whereDecl)
}