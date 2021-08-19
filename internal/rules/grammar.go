package rules

import (
	"errors"
	"fmt"
	"net"
	"reflect"

	"github.com/alecthomas/participle/v2"
	"github.com/alecthomas/participle/v2/lexer"
)

var (
	ErrNoParser     = errors.New("no parser available for given type")
	ErrTypeMismatch = errors.New("param has a different type")
	parsers         map[reflect.Type]*participle.Parser
)

// nolint:lll
func init() {
	ruleLexer := lexer.Must(lexer.NewSimple([]lexer.Rule{
		{Name: `Module`, Pattern: `[a-z]+`, Action: nil},
		{Name: `Ident`, Pattern: `[A-Z][a-zA-Z0-9_]*`, Action: nil},
		{Name: `CIDR`, Pattern: `(25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)(\.(25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)){3}/(3[0-2]|[1-2][0-9]|[1-9])`, Action: nil},
		{Name: `IP`, Pattern: `(25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)(\.(25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)){3}`, Action: nil},
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

// Parse takes a raw rule and parses it into the given target instance
// currently only Routing and Check are supported for parsing
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

// Routing describes how DNS or HTTP requests are handled
// A routing is described as a optional chain of Filters like: filter1() -> filter2()
// and a Terminator which determines how the request should be handled e.g. http.Status(204)
// a full chain might look like so: GET() -> Header("Accept", "application/json") -> http.Status(200).
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
	IP     net.IP   `parser:"| @IP"`
	CIDR   *CIDR    `parser:"| @CIDR"`
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

func (p Param) AsIP() (net.IP, error) {
	if p.IP == nil {
		return nil, fmt.Errorf("IP is nil %w", ErrTypeMismatch)
	}
	return p.IP, nil
}

func (p Param) AsCIDR() (*CIDR, error) {
	if p.CIDR == nil {
		return nil, fmt.Errorf("IP is nil %w", ErrTypeMismatch)
	}
	return p.CIDR, nil
}
