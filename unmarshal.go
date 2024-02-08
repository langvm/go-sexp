// Copyright 2024 LangVM Project
// This Source Code Form is subject to the terms of the Mozilla Public License, v. 2.0
// that can be found in the LICENSE file and https://mozilla.org/MPL/2.0/.

package sexp

import (
	"reflect"
	"sexp/parser"
	"strconv"
)

func UnmarshalStruct(v reflect.Value, node parser.List) error {
	var (
		set     = reflect.Indirect(v)
		typ     = set.Type()
		fields  = reflect.VisibleFields(typ)
		nodeMap = node.Map()
	)

	for i, field := range fields {
		if node := nodeMap[field.Name]; node != nil {
			switch node := node.(type) {
			case parser.List:
				v := set.Field(i)
				switch field.Type.Kind() {
				case reflect.Int:
					fallthrough
				case reflect.Int8:
					fallthrough
				case reflect.Int16:
					fallthrough
				case reflect.Int32:
					fallthrough
				case reflect.Int64:
					if len(node.Elements) != 1 {
						return ListLengthMismatchError{}
					}
					switch lit := node.Elements[0].(type) {
					case parser.LiteralValue:
						i, err := strconv.ParseInt(lit.Literal, lit.Format, 64)
						if err != nil {
							return err
						}
						v.SetInt(i)
					default:
						return TypeMismatchError{}
					}
				case reflect.Uint:
					fallthrough
				case reflect.Uint8:
					fallthrough
				case reflect.Uint16:
					fallthrough
				case reflect.Uint32:
					fallthrough
				case reflect.Uint64:
					if len(node.Elements) != 1 {
						return ListLengthMismatchError{}
					}
					switch lit := node.Elements[0].(type) {
					case parser.LiteralValue:
						i, err := strconv.ParseUint(lit.Literal, lit.Format, 64)
						if err != nil {
							return err
						}
						v.SetUint(i)
					default:
						return TypeMismatchError{}
					}
				case reflect.Float32:
					fallthrough
				case reflect.Float64:
					f, err := strconv.ParseFloat("", 64)
					if err != nil {
						return err
					}
					v.SetFloat(f)
				case reflect.String:
					if len(node.Elements) != 1 {
						return ListLengthMismatchError{}
					}
					switch lit := node.Elements[0].(type) {
					case parser.LiteralValue:
						v.SetString(lit.Literal)
					default:
						return TypeMismatchError{}
					}
				case reflect.Slice:
					// TODO
				case reflect.Struct:
					err := UnmarshalStruct(v, node)
					if err != nil {
						return err
					}
				default:
					return UnsupportedType{reflect.TypeOf(v).String()}
				}
			default:
				return UnsupportedType{Type: "fuck you"}
			}
		}
	}

	return nil
}

func Unmarshal(v any, list parser.List) error {
	return UnmarshalStruct(reflect.ValueOf(v), list)
}
