package format

import (
	"errors"
	"fmt"
	"reflect"
	"strconv"

	"github.com/olekukonko/tablewriter"
)

type tblWriter struct {
	tableWriter *tablewriter.Table
}

//nolint:exhaustive
func (t *tblWriter) Write(in any) (err error) {
	v := reflect.ValueOf(in)

	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}

	var vt reflect.Type
	var numberOfFields int

	data := make([][]string, 0)

	switch v.Kind() {
	case reflect.Interface:
		return errors.New("interface{} is not supported")
	case reflect.Slice, reflect.Array:
		length := v.Len()
		if length < 1 {
			return
		}

		vt = v.Index(0).Type()

		if vt.Kind() == reflect.Ptr {
			vt = vt.Elem()
		}

		if vt.Kind() != reflect.Struct {
			return fmt.Errorf("element type of array %v is not supported", vt.Kind())
		}

		numberOfFields = vt.NumField()

		for i := 0; i < length; i++ {
			data = append(data, t.getData(v.Index(i), numberOfFields))
		}
	case reflect.Struct:
		vt = v.Type()
		numberOfFields = vt.NumField()
		data = append(data, t.getData(v, numberOfFields))
	}

	t.tableWriter.SetHeader(headersForType(vt, numberOfFields))
	t.tableWriter.AppendBulk(data)
	t.tableWriter.Render()
	t.tableWriter.ClearRows()

	return err
}

func (t *tblWriter) getData(val reflect.Value, numberOfFields int) (data []string) {
	if val.Kind() == reflect.Ptr {
		val = val.Elem()
	}

	for i := 0; i < numberOfFields; i++ {
		field := val.Field(i)
		if !isPrivateField(val.Type(), i) {
			data = append(data, value(field))
		}
	}

	return
}

//nolint:exhaustive
func value(val reflect.Value) string {
	const (
		base10         = 10
		floatPrecision = 6
		longBitSize    = 64
	)

	if val.IsZero() || !val.IsValid() {
		return ""
	}

	if stringer, isStringer := val.Interface().(fmt.Stringer); isStringer {
		return stringer.String()
	}

	switch val.Kind() {
	case reflect.Ptr:
		return value(val.Elem())
	case reflect.Struct, reflect.Interface:
		return "<not supported>"
	case reflect.Bool:
		return strconv.FormatBool(val.Bool())
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return strconv.FormatInt(val.Int(), base10)
	case reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return strconv.FormatUint(val.Uint(), base10)
	case reflect.Float32, reflect.Float64:
		return strconv.FormatFloat(val.Float(), 'f', floatPrecision, longBitSize)
	default:
		return val.String()
	}
}

func headersForType(t reflect.Type, numberOfFields int) (headers []string) {
	for i := 0; i < numberOfFields; i++ {
		if isPrivateField(t, i) {
			continue
		}
		field := t.Field(i)
		if tableTag, ok := field.Tag.Lookup("table"); ok {
			headers = append(headers, tableTag)
		} else {
			headers = append(headers, field.Name)
		}
	}
	return
}

func isPrivateField(t reflect.Type, idx int) bool {
	varPrefix := t.Field(idx).Name[0:1]
	if varPrefix >= "a" && varPrefix <= "z" {
		return true
	}
	return false
}
