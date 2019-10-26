package parser

import (
	"github.com/mlhoyt/ramsql/engine/parser/lexer"
)

func (p *Parser) parseDelete() (*Instruction, error) {
	i := &Instruction{}

	// Set DELETE decl
	deleteDecl, err := p.consumeToken(lexer.DeleteToken)
	if err != nil {
		return nil, err
	}
	i.Decls = append(i.Decls, deleteDecl)

	// should be From
	fromDecl, err := p.consumeToken(lexer.FromToken)
	if err != nil {
		return nil, err
	}
	deleteDecl.Add(fromDecl)

	// Should be a table name
	nameDecl, err := p.parseQuotedToken()
	if err != nil {
		return nil, err
	}
	fromDecl.Add(nameDecl)

	// MAY be WHERE  here
	debug("WHERE ? %v", p.tokens[p.index])
	if !p.hasNext() {
		return i, nil
	}

	err = p.parseWhere(deleteDecl)
	if err != nil {
		return nil, err
	}

	return i, nil
}
