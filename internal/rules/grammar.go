package rules

import (
	"errors"
	"fmt"
	"reflect"

	"github.com/alecthomas/participle/v2"
	"github.com/alecthomas/participle/v2/lexer"
	"github.com/alecthomas/participle/v2/lexer/stateful"
)

var (
	ErrNoParser     = errors.New("no parser available for given type")
	ErrTypeMismatch = errors.New("param has a different type")
	parsers         map[reflect.Type]*participle.Parser
)

func init() {
	ruleLexer := lexer.Must(stateful.NewSimple([]stateful.Rule{
		{Name: `Module`, Pattern: `[a-z]+`, Action: nil},
		{Name: `Ident`, Pattern: `[A-Z][a-zA-Z0-9_]*`, Action: nil},
		{Name: `Float`, Pattern: `\d+\.\d+`, Action: nil},
		{Name: `Int`, Pattern: `[-]?\d+`, Action: nil},
		{Name: `RawString`, Pattern: "`[^`]*`", Action: nil},
		{Name: `String`, Pattern: `'[^']*'|"[^"]*"`, Action: nil},
		{Name: `Arrows`, Pattern: `(->|=>)`, Action: nil},
		{Name: "whitespace", Pattern: `\s+`, Action: nil},
		{Name: "Punct", Pattern: `[-[!@#$%^&*()+_={}\|:;\."'<,>?/]|]`, Action: nil},
	}))

	parsers = map[reflect.Type]*participle.Parser{
		reflect.TypeOf(new(Routing)): participle.MustBuild(
			new(Routing),
			participle.Lexer(ruleLexer),
			participle.Unquote("String"),
			participle.Unquote("RawString"),
		),
		reflect.TypeOf(new(Check)): participle.MustBuild(
			new(Check),
			participle.Lexer(ruleLexer),
			participle.Unquote("String"),
			participle.Unquote("RawString"),
		),
	}
}

func Parse(rule string, target interface{}) error {
	parser, available := parsers[reflect.TypeOf(target)]
	if !available {
		return ErrNoParser
	}
	if err := parser.ParseString("", rule, target); err != nil {
		return err
	}
	return nil
}

type Routing struct {
	Filters    *Filters `parser:"@@*"`
	Terminator *Method  `parser:"'=>' @@"`
}

type Check struct {
	Initiator  *Method  `parser:"@@"`
	Validators *Filters `parser:"( '=>' @@)?"`
}

type Filters struct {
	Chain []Method `parser:"@@ ('->' @@)*"`
}

type Method struct {
	Module string  `parser:"(@Module'.')?"`
	Name   string  `parser:"@Ident"`
	Params []Param `parser:"'(' @@? ( ',' @@ )*')'"`
}

type Param struct {
	String *string  `parser:"@String | @RawString"`
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
	if p.Float == nil {
		return 0, fmt.Errorf("float is nil %w", ErrTypeMismatch)
	}
	return *p.Float, nil
}
