// Copyright 2024 LangVM Project
// This Source Code Form is subject to the terms of the Mozilla Public License, v. 2.0
// that can be found in the LICENSE file and https://mozilla.org/MPL/2.0/.

package parser

import (
	"fmt"
	scanner "github.com/langvm/cee-scanner"
)

type Node interface {
	GetPosRange() PosRange
	String() string
}

type PosRange struct {
	From, To scanner.Position
}

func (p PosRange) String() string { return fmt.Sprintln(p.From.String(), "->", p.To.String()) }

func (p PosRange) GetPosRange() PosRange { return p }

type Token struct {
	PosRange
	Kind, Format int
	Literal      string
}

func (token Token) String() string { return token.Literal }

type Ident struct {
	Token
}

type LiteralValue struct {
	Token
}

type List struct {
	PosRange
	Prefix   Ident
	Elements []Node
}

func (l List) String() string {
	return fmt.Sprintln(l.Elements)
}

func (l List) Map() (m map[string]Node) {
	for _, e := range l.Elements {
		m[e.(List).Prefix.Literal] = e
	}
	return
}
