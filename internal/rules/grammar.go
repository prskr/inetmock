package rules

import (
	"errors"
	"fmt"

	"github.com/alecthomas/participle/v2"
)

var (
	ErrTypeMismatch = errors.New("param has a different type")
	parser          *participle.Parser
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
	String *string  `parser:"@String|RawString"`
	Int    *int     `parser:"| @Int"`
	Float  *float64 `parser:"| @Float"`
}

func (p Param) AsString() (string, error) {
	if p.String == nil {
		return "", fmt.Errorf("string is nil %w", ErrTypeMismatch)
	}
	return *p.String, nil
}

func (p Param) AsInt() (int, error) {
	if p.Int == nil {
		return 0, fmt.Errorf("int is nil %w", ErrTypeMismatch)
	}
	return *p.Int, nil
}

func (p Param) AsFloat() (float64, error) {
	if p.Int == nil {
		return 0, fmt.Errorf("float is nil %w", ErrTypeMismatch)
	}
	return *p.Float, nil
}
