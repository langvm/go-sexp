// Copyright 2024 LangVM Project
// This Source Code Form is subject to the terms of the Mozilla Public License, v. 2.0
// that can be found in the LICENSE file and https://mozilla.org/MPL/2.0/.

package parser

import "fmt"

type Node interface {
	GetPosRange() PosRange
}

type PosRange struct {
	From, To Position
}

func (p PosRange) String() string { return fmt.Sprintln(p.From.String(), "->", p.To.String()) }

func (p PosRange) GetPosRange() PosRange { return p }

type Token struct {
	PosRange
	Kind    int
	Literal string
}

type List struct {
	PosRange
	Prefix Token
	List   []Node
}
