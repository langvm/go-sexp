// Copyright 2024 LangVM Project
// This Source Code Form is subject to the terms of the Mozilla Public License, v. 2.0
// that can be found in the LICENSE file and https://mozilla.org/MPL/2.0/.

package sexp

import (
	"reflect"
	"sexp/parser"
	"strconv"
)

func GetTheOnlyElement(list parser.List) (parser.Node, error) {
	if len(list.Elements) != 1 {
		return nil, ListLengthMismatchError{}
	}
	node := list.Elements[0]
	return node, nil
}

func UnmarshalValue(v reflect.Value, node parser.Node) error {
	switch lit := node.(type) {
	case parser.LiteralValue:
		switch v.Type().Kind() {
		case reflect.Int:
			fallthrough
		case reflect.Int8:
			fallthrough
		case reflect.Int16:
			fallthrough
		case reflect.Int32:
			fallthrough
		case reflect.Int64:
			i, err := strconv.ParseInt(lit.Literal, lit.Format, 64)
			if err != nil {
				return err
			}
			v.SetInt(i)
		case reflect.Uint:
			fallthrough
		case reflect.Uint8:
			fallthrough
		case reflect.Uint16:
			fallthrough
		case reflect.Uint32:
			fallthrough
		case reflect.Uint64:
			u, err := strconv.ParseUint(lit.Literal, lit.Format, 64)
			if err != nil {
				return err
			}
			v.SetUint(u)
		case reflect.Float32:
			fallthrough
		case reflect.Float64:
			panic("not implemented")
		case reflect.String:
			v.SetString(lit.Literal)
		case reflect.Bool:
			switch lit.Literal {
			case "true":
				v.SetBool(true)
			case "false":
				v.SetBool(false)
			default:
				return parser.UnexpectedToken{} // TODO
			}
		default:
			return UnsupportedType{}
		}
		return nil
	default:
		return TypeMismatchError{}
	}
}

func UnmarshalField(v reflect.Value, node parser.Node) error {
	switch list := node.(type) {
	case parser.List:
		node, err := GetTheOnlyElement(list)
		if err != nil {
			return err
		}
		return UnmarshalValue(v, node)
	case parser.LiteralValue:
		return UnmarshalValue(v, node)
	default:
		return TypeMismatchError{}
	}
}

func UnmarshalNode(v reflect.Value, node parser.Node) error {
	switch v.Type().Kind() {
	case reflect.Struct:
		return UnmarshalStruct(v, node)
	case reflect.Slice:
		return UnmarshalSlice(v, node)
	default:
	}
	return UnmarshalField(v, node)
}

func UnmarshalSlice(v reflect.Value, node parser.Node) error {
	switch list := node.(type) {
	case parser.List:
		arr := reflect.MakeSlice(v.Type(), len(list.Elements), cap(list.Elements))
		for i, e := range list.Elements {
			if err := UnmarshalNode(arr.Index(i), e); err != nil {
				return err
			}
		}
		v.Set(arr)
	default:
		return TypeMismatchError{}
	}
	return nil
}

func UnmarshalStruct(v reflect.Value, node parser.Node) error {
	switch list := node.(type) {
	case parser.List:
		var (
			structValue = reflect.Indirect(v)
			structType  = structValue.Type()
			fields      = reflect.VisibleFields(structType)
			elementsMap = list.Map()
		)

		for i, field := range fields {
			if node := elementsMap[field.Name]; node != nil {
				if err := UnmarshalNode(structValue.Field(i), node); err != nil {
					return err
				}
			}
		}

		return nil
	default:
		return TypeMismatchError{}
	}
}

func Unmarshal(v any, list parser.List) error {
	return UnmarshalStruct(reflect.ValueOf(v), list)
}
