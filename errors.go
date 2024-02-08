// Copyright 2024 LangVM Project
// This Source Code Form is subject to the terms of the Mozilla Public License, v. 2.0
// that can be found in the LICENSE file and https://mozilla.org/MPL/2.0/.

package sexp

import (
	"fmt"
	"reflect"
	"sexp/parser"
)

type ListLengthMismatchError struct {
	Field string
}

func (e ListLengthMismatchError) Error() string {
	return "list length mismatch"
}

type UndefinedEnumTypeError struct {
	Found parser.Token
}

func (e UndefinedEnumTypeError) Error() string {
	return fmt.Sprintln(e.Found.PosRange.String(), "undefined enum type:", e.Found.Literal)
}

type UndefinedIdentifierError struct {
	Token parser.Token
}

func (e UndefinedIdentifierError) Error() string {
	return fmt.Sprintln(e.Token.PosRange.String(), "undefined identifier:", e.Token.Literal)
}

type TypeMismatchError struct {
	parser.PosRange
	Want reflect.Kind
	Have reflect.Kind
}

func (e TypeMismatchError) Error() string {
	return fmt.Sprintln(e.PosRange.String(), "type mismatch, want:", e.Want, "but have:", e.Have)
}

type UnsupportedType struct {
	Type string
}

func (e UnsupportedType) Error() string { return e.Type }
