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

func init() {
	ruleLexer := lexer.Must(lexer.NewSimple([]lexer.SimpleRule{
		{Name: "Comment", Pattern: `(?:#|//)[^\n]*\n?`},
		{Name: `Module`, Pattern: `[a-z]{1}[A-z0-9]+`},
		{Name: `Ident`, Pattern: `[A-Z][a-zA-Z0-9_]*`},
		{Name: `CIDR`, Pattern: `(25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)(\.(25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)){3}/(3[0-2]|[1-2][0-9]|[1-9])`},
		{Name: `IP`, Pattern: `(25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)(\.(25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)){3}`},
		{Name: `Float`, Pattern: `\d+\.\d+`},
		{Name: `Int`, Pattern: `[-]?\d+`},
		{Name: `RawString`, Pattern: "`[^`]*`"},
		{Name: `String`, Pattern: `'[^']*'|"[^"]*"`},
		{Name: `Arrows`, Pattern: `(->|=>)`},
		{Name: "whitespace", Pattern: `\s+`},
		{Name: "Punct", Pattern: `[-[!@#$%^&*()+_={}\|:;\."'<,>?/]|]`},
	}))

	parsers = map[reflect.Type]*participle.Parser{
		reflect.TypeOf(new(SingleResponsePipeline)): participle.MustBuild(
			new(SingleResponsePipeline),
			participle.Lexer(ruleLexer),
			participle.Unquote("String"),
			participle.Unquote("RawString"),
		),
		reflect.TypeOf(new(ChainedResponsePipeline)): participle.MustBuild(
			new(ChainedResponsePipeline),
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
		reflect.TypeOf(new(CheckScript)): participle.MustBuild(
			new(CheckScript),
			participle.Lexer(ruleLexer),
			participle.Elide("Comment"),
			participle.Unquote("String"),
			participle.Unquote("RawString"),
		),
	}
}

// Parse takes a raw rule and parses it into the given target instance
// currently only SingleResponsePipeline and Check are supported for parsing
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

type FilteredPipeline interface {
	Filters() []Call
}

// SingleResponsePipeline describes how requests are handled that expect single response
// e.g. HTTP or DNS requests
// A SingleResponsePipeline is defined as an optional chain of Filters like: filter1() -> filter2()
// and a Response which determines how the request should be handled e.g. http.Status(204)
// a full chain might look like so: GET() -> Header("Accept", "application/json") -> http.Status(200).
type SingleResponsePipeline struct {
	FilterChain *Filters `parser:"@@*"`
	Response    *Call    `parser:"'=>' @@"`
}

func (p *SingleResponsePipeline) Filters() []Call {
	if p.FilterChain != nil {
		return p.FilterChain.Chain
	}
	return nil
}

// ChainedResponsePipeline describes how requests are handled that expect a chain of response handlers
// e.g. DHCP where one handler might set the IP, one sets the default gateway, one sets the DNS servers and so on
// A ChainedResponsePipeline is defined as an optional chain of Filters like: filter1() -> filter2()
// and a chain of Response handlers while at least one has to be present
// // a full chain might look like so: MatchMAC(`00:06:7C:.*`) => IP(3.3.6.6) => Router(1.2.3.4) => DNS(1.2.3.4, 4.5.6.7)
type ChainedResponsePipeline struct {
	FilterChain *Filters `parser:"@@*"`
	Response    []Call   `parser:"'=>' @@ ('=>' @@)*"`
}

func (p *ChainedResponsePipeline) Filters() []Call {
	if p.FilterChain != nil {
		return p.FilterChain.Chain
	}
	return nil
}

type Check struct {
	Initiator  *Call    `parser:"@@"`
	Validators *Filters `parser:"( '=>' @@)?"`
}

type CheckScript struct {
	Checks []Check `parser:"@@*"`
}

type Filters struct {
	Chain []Call `parser:"@@ ('->' @@)*"`
}

type Call struct {
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
