package rules

import (
	"github.com/alecthomas/participle/v2"
)

var (
	parser *participle.Parser
)

func init() {
	parser = participle.MustBuild(
		new(Routing),
		participle.Unquote("String"),
	)
}

func Parse(rule string) (*Routing, error) {
	routing := new(Routing)
	if err := parser.ParseString("", rule, routing); err != nil {
		return nil, err
	}
	return routing, nil
}

type Routing struct {
	Filters    *Filters `parser:"@@*"`
	Terminator *Method  `parser:"'=''>' @@"`
}

type Filters struct {
	Chain []Method `parser:"@@ ('-''>' @@)*"`
}

type Method struct {
	Name   string  `parser:"@Ident"`
	Params []Param `parser:"'(' @@ ( ',' @@ )*')'"`
}

type Param struct {
	String *string `parser:"@String|RawString"`
	Int    *int    `parser:"| @Int"`
}
