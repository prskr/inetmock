package format

import (
	"errors"
	"fmt"
	"github.com/olekukonko/tablewriter"
	"reflect"
	"strconv"
)

type tblWriter struct {
	tableWriter *tablewriter.Table
}

func (t *tblWriter) Write(in interface{}) (err error) {
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

	return
}

func (t *tblWriter) getData(val reflect.Value, numberOfFields int) (data []string) {
	if val.Kind() == reflect.Ptr {
		val = val.Elem()
	}
	for i := 0; i < numberOfFields; i++ {
		data = append(data, value(val.Field(i)))
	}

	return
}

func value(val reflect.Value) string {
	switch val.Kind() {
	case reflect.Ptr:
		return value(val.Elem())
	case reflect.Struct, reflect.Interface:
		return "<not supported>"
	case reflect.Bool:
		return strconv.FormatBool(val.Bool())
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return strconv.FormatInt(val.Int(), 10)
	case reflect.Float32, reflect.Float64:
		return strconv.FormatFloat(val.Float(), 'f', 6, 64)
	default:
		return val.String()
	}
}

func headersForType(t reflect.Type, numberOfFields int) (headers []string) {
	for i := 0; i < numberOfFields; i++ {
		field := t.Field(i)
		if tableTag, ok := field.Tag.Lookup("table"); ok {
			headers = append(headers, tableTag)
		} else {
			headers = append(headers, field.Name)
		}
	}
	return
}
