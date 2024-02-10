// Copyright 2024 LangVM Project
// This Source Code Form is subject to the terms of the Mozilla Public License, v. 2.0
// that can be found in the LICENSE file and https://mozilla.org/MPL/2.0/.

package sexp

import (
	"fmt"
	"sexp/parser"
	"testing"
)

func TestUnmarshalStruct(t *testing.T) {
	v := struct {
		IntA    int
		StructA struct {
			IntB, IntC int
			ArrA       []int
		}
	}{}

	list, err := parser.Parse([]rune(`
(a (IntA 1) (StructA (IntB 2) (IntC 3) (ArrA 1 2 3 4 5 6)))
`))
	if err != nil {
		t.Fatal(err)
	}

	t.Log("Parsed successfully.")

	err = Unmarshal(&v, list)
	if err != nil {
		t.Fatal(err)
	}

	fmt.Println(v.IntA, v.StructA.IntB, v.StructA.IntC, v.StructA.ArrA)
}
