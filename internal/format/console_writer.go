package format

import (
	"encoding/json"
	"github.com/olekukonko/tablewriter"
	"gopkg.in/yaml.v2"
	"io"
	"strings"
)

type consoleWriterFactory func(io.Writer) ConsoleWriter

var (
	writers = map[string]consoleWriterFactory{
		"table": func(writer io.Writer) ConsoleWriter {
			tw := tablewriter.NewWriter(writer)
			tw.SetBorders(tablewriter.Border{Left: true, Top: false, Right: true, Bottom: false})
			tw.SetCenterSeparator("|")

			return &tblWriter{
				tableWriter: tw,
			}
		},
		"json": func(writer io.Writer) ConsoleWriter {
			return &jsonWriter{
				encoder: json.NewEncoder(writer),
			}
		},
		"yaml": func(writer io.Writer) ConsoleWriter {
			return &yamlWriter{
				encoder: yaml.NewEncoder(writer),
			}
		},
	}
)

func Writer(format string, writer io.Writer) ConsoleWriter {
	if cw, ok := writers[strings.ToLower(format)]; ok {
		return cw(writer)
	}
	return writers["table"](writer)
}

type ConsoleWriter interface {
	Write(in interface{}) (err error)
}
