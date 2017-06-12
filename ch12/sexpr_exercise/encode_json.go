//e12.5

package sexpr

import (
	"bytes"
	"fmt"
	"reflect"
)

//这种实现结构上没有做到模块化，需改进

//JSONIndent ecode go value in Json format
//print unexport field, unlike json.Marshal
func JSONIndent(v interface{}) ([]byte, error) {
	var buf = new(bytes.Buffer)
	const pre = "  "
	var nPre = 0
	indent := func() {
		for i := 0; i < nPre; i++ {
			buf.WriteString(pre)
		}
	}
	var json func(v reflect.Value) error
	json = func(v reflect.Value) error {
		switch v.Kind() {
		case reflect.Invalid:
			buf.WriteString("null")

		case reflect.Int, reflect.Int8, reflect.Int16,
			reflect.Int32, reflect.Int64:
			fmt.Fprintf(buf, "%d", v.Int())

		case reflect.Uint, reflect.Uint8, reflect.Uint16,
			reflect.Uint32, reflect.Uint64, reflect.Uintptr:
			fmt.Fprintf(buf, "%d", v.Uint())

		case reflect.Float32, reflect.Float64:
			fmt.Fprintf(buf, "%g", v.Float())

		case reflect.Bool:
			if v.Bool() {
				fmt.Fprint(buf, "true")
			} else {
				fmt.Fprint(buf, "false")
			}
		case reflect.String:
			fmt.Fprintf(buf, "%q", v.String())

		case reflect.Ptr:
			return json(v.Elem())

		case reflect.Array, reflect.Slice: // (value ...)
			buf.WriteString("[")
			nPre++
			for i := 0; i < v.Len(); i++ {
				if i > 0 {
					buf.WriteString(",")
				}
				buf.WriteByte('\n')
				indent()
				if err := json(v.Index(i)); err != nil {
					return err
				}
			}
			nPre--
			if v.Len() > 0 {
				buf.WriteByte('\n')
				indent()
			}
			buf.WriteByte(']')

		case reflect.Struct: // ((name value) ...)
			buf.WriteString("{")
			nPre++
			for i := 0; i < v.NumField(); i++ {
				if i > 0 {
					buf.WriteString(",")
				}
				buf.WriteByte('\n')
				indent()
				fmt.Fprintf(buf, "%q: ", v.Type().Field(i).Name)
				if err := json(v.Field(i)); err != nil {
					return err
				}
			}
			nPre--
			if v.NumField() > 0 {
				buf.WriteByte('\n')
				indent()
			}
			buf.WriteByte('}')

		case reflect.Map: // ((key value) ...)
			buf.WriteString("{")
			nPre++
			for i, key := range v.MapKeys() {
				if i > 0 {
					buf.WriteString(",\n")
				}
				indent()
				if key.Kind() != reflect.String {
					return fmt.Errorf("map key type %s is not string", key.Type())
				}
				if err := json(key); err != nil {
					return err
				}
				buf.WriteString(": ")
				if err := json(v.MapIndex(key)); err != nil {
					return err
				}
			}
			nPre--
			if len(v.MapKeys()) > 0 {
				buf.WriteString("\n")
				indent()
			}
			buf.WriteString("}")

		default: // chan, func, interface
			return fmt.Errorf("unsupported type: %s", v.Type())
		}
		return nil
	}

	if err := json(reflect.ValueOf(v)); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}
