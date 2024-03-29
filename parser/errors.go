// Copyright 2024 LangVM Project
// This Source Code Form is subject to the terms of the Mozilla Public License, v. 2.0
// that can be found in the LICENSE file and https://mozilla.org/MPL/2.0/.

package parser

import (
	"fmt"
	scanner "github.com/langvm/cee-scanner"
)

type FormatError struct {
	Pos scanner.Position
}

func (e FormatError) Error() string {
	return fmt.Sprintln(e.Pos.String(), "format error")
}

type UnexpectedToken struct {
	Token
}

func (e UnexpectedToken) Error() string { return fmt.Sprintln(e.PosRange.String(), "unexpected token") }
