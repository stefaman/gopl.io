package params

import (
	"bytes"
	"fmt"
	"io"
	"reflect"
	"strings"
)

//Pack encode a struct that storege URL form
func Pack(data interface{}) (_ string, err error) {
	defer func() {
		if v := recover(); v != nil {
			err = fmt.Errorf("Pack: %v", v)
		}
	}()
	v := reflect.ValueOf(data)
	if v.Kind() == reflect.Ptr { //can pass a pointer to struct
		v = v.Elem()
	}
	buf := new(bytes.Buffer)
	typ := v.Type()
	for i := 0; i < typ.NumField(); i++ {
		if i > 0 {
			fmt.Fprintf(buf, "%c", '&')
		}
		name := getName(typ.Field(i))
		// buf.WriteString(name + "=")
		err := parse(buf, name, v.Field(i))
		if err != nil {
			return "", err
		}
	}
	return buf.String(), nil
}

func getName(i reflect.StructField) (name string) {
	if name = i.Tag.Get("http"); name == "" {
		name = strings.ToLower(i.Name)
	}
	return
}

func parse(buf io.Writer, name string, v reflect.Value) error {
	switch v.Kind() {
	// case reflect.Int:
	// 	fmt.Fprintf(buf, "%s=%d", name, v.Int())
	// case reflect.Float64:
	// 		fmt.Fprintf(buf, "%s=%f", name, v.Float())
	case reflect.Slice, reflect.Array:
		for i := 0; i < v.Len(); i++ {
			if i > 0 {
				fmt.Fprintf(buf, "%c", '&') //ignore error
			}
			if err := parse(buf, name, v.Index(i)); err != nil {
				return err
			}
		}

	case reflect.Func, reflect.Chan, reflect.Complex64, reflect.Complex128, reflect.Map, reflect.Ptr:
		return fmt.Errorf("unsupported type %s", v.Type())

	default: //int, float, bool
		fmt.Fprintf(buf, "%s=%v", name, v)
	}
	return nil
}
