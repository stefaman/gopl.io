// Copyright © 2016 Alan A. A. Donovan & Brian W. Kernighan.
// License: https://creativecommons.org/licenses/by-nc-sa/4.0/

// See page 349.

// Package params provides a reflection-based parser for URL parameters.
package params

import (
	"fmt"
	"net/http"
	"reflect"
	"strings"
)

//!+Unpack

// UnpackVal populates the fields of the struct pointed to by ptr
// from the HTTP request parameters in req.
//Check validity of field values
func UnpackVal(req *http.Request, ptr interface{}) error {
	if err := req.ParseForm(); err != nil {
		return err
	}

	type field struct {
		value    reflect.Value
		validity string
	}
	// Build map of fields keyed by effective name.
	fields := make(map[string]field)
	v := reflect.ValueOf(ptr).Elem() // the struct variable
	for i := 0; i < v.NumField(); i++ {
		fieldInfo := v.Type().Field(i) // a reflect.StructField
		tag := fieldInfo.Tag           // a reflect.StructTag
		name := tag.Get("http")
		if name == "" {
			name = strings.ToLower(fieldInfo.Name)
		}
		validity := tag.Get("validity")
		fields[name] = field{v.Field(i), validity}
	}

	// Update struct field for each parameter in the request.
	for name, values := range req.Form {
		f := fields[name].value
		if !f.IsValid() {
			continue // ignore unrecognized HTTP parameters
		}
		for _, value := range values {
			if err := check(value, fields[name].validity); err != nil {
				return fmt.Errorf("%s: %v", name, err)
			}
			if f.Kind() == reflect.Slice {
				elem := reflect.New(f.Type().Elem()).Elem()
				if err := populate(elem, value); err != nil {
					return fmt.Errorf("%s: %v", name, err)
				}
				f.Set(reflect.Append(f, elem))
			} else {
				if err := populate(f, value); err != nil {
					return fmt.Errorf("%s: %v", name, err)
				}
			}
		}
	}
	return nil
}

//!-Unpack

func check(value, validity string) error {
	switch validity {
	case "":
		return nil
	case "email":

	case "ZIP":
	case "creditCard":
	default:
		return fmt.Errorf("unsupported validity %s", validity)
	}
	return nil
}

//!-populate
