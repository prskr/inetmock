package rules

import (
	"errors"
	"fmt"

	"github.com/alecthomas/participle/v2"
	"github.com/alecthomas/participle/v2/lexer"
	"github.com/alecthomas/participle/v2/lexer/stateful"
)

var (
	ErrTypeMismatch = errors.New("param has a different type")
	parser          *participle.Parser
)

func init() {
	sqlLexer := lexer.Must(stateful.NewSimple([]stateful.Rule{
		{Name: `Ident`, Pattern: `[a-zA-Z_][a-zA-Z0-9_]*`, Action: nil},
		{Name: `Float`, Pattern: `\d+.\d+`, Action: nil},
		{Name: `Int`, Pattern: `[-]?\d+`, Action: nil},
		{Name: `String`, Pattern: `'[^']*'|"[^"]*"`, Action: nil},
		{Name: `Arrows`, Pattern: `(->|=>)`, Action: nil},
		{Name: "whitespace", Pattern: `\s+`, Action: nil},
		{Name: "Punct", Pattern: `[-[!@#$%^&*()+_={}\|:;\."'<,>?/]|]`, Action: nil},
	}))

	parser = participle.MustBuild(
		new(Routing),
		participle.Lexer(sqlLexer),
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
	Terminator *Method  `parser:"'=>' @@"`
}

type Filters struct {
	Chain []Method `parser:"@@ ('->' @@)*"`
}

type Method struct {
	Name   string  `parser:"@Ident"`
	Params []Param `parser:"'(' @@ ( ',' @@ )*')'"`
}

type Param struct {
	String *string  `parser:"@String"`
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
