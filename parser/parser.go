// Copyright 2024 LangVM Project
// This Source Code Form is subject to the terms of the Mozilla Public License, v. 2.0
// that can be found in the LICENSE file and https://mozilla.org/MPL/2.0/.

package parser

type Parser struct {
	Scanner

	ReachedEOF bool

	Token Token
}

func (p *Parser) Scan() {
	defer func() {
		switch v := recover().(type) {
		case nil:
			return
		case EOFError:
			if p.ReachedEOF {
				panic(v)
			} else {
				p.ReachedEOF = true
			}
		default:
			panic(v)
		}
	}()

	begin, kind, lit := p.Scanner.ScanToken()
	p.Token = Token{
		PosRange: PosRange{begin, p.Position},
		Kind:     kind,
		Literal:  lit,
	}
}

func (p *Parser) MatchTerm(term int) {
	tok := p.Token

	p.Scan()

	if tok.Kind != term {
		panic(UnexpectedToken{Token: p.Token})
	}
}

func (p *Parser) ExpectIdent() Token {
	tok := p.Token
	p.MatchTerm(IDENT)
	return tok
}

func (p *Parser) ExpectList() List {
	var list List

	p.MatchTerm(LPAREN)

	list.Prefix = p.ExpectIdent()

	for p.Token.Kind != RPAREN {
		list.List = append(list.List, p.ExpectNode())
	}
	p.MatchTerm(RPAREN)

	return list
}

func (p *Parser) ExpectNode() Node {
	switch p.Token.Kind {
	case LPAREN:
		return p.ExpectList()
	default:
		tok := p.Token
		p.Scan()
		return tok
	}
}

func Parse(data []rune) List {
	p := Parser{
		Scanner: Scanner{
			BufferScanner{Buffer: data},
		},
		Token: Token{},
	}

	p.Scan()

	return p.ExpectList()
}
