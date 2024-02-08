// Copyright 2024 LangVM Project
// This Source Code Form is subject to the terms of the Mozilla Public License, v. 2.0
// that can be found in the LICENSE file and https://mozilla.org/MPL/2.0/.

package sexp

import "sexp/parser"

// Integers with tags are processed as enums.

type Marshaller struct {
}

func (m *Marshaller) Marshal(v any) parser.Node {

	return nil
}
