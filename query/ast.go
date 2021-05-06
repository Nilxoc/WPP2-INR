package query

import (
	"fmt"
	"g1.wpp2.hsnr/inr/boolret/index"
	"github.com/alecthomas/participle/v2"
	"github.com/alecthomas/participle/v2/lexer/stateful"
	"github.com/alecthomas/repr"
)

type AstQuery struct {
	index *index.Index
	root  Expression
}

type AstQueryParser struct {
	index *index.Index
}

type astNode interface {
	nodeType() string
	terminal() bool
	evaluate(index *index.Index) (*index.PostingList, error)
}

func Evaluate(q *AstQuery) (*index.PostingList, error) {
	return nil, nil
}

// AST-Nodes

type Value struct {
	Text          *string     `parser:"( @Ident"`
	Phrase        *string     `parser:"| @String)"`
	Subexpression *Expression `parser:"| '(' @@ ')'"`
}

type Term struct {
	Left  *Value   `parser:"@@"`
	Right []*Value `@@*`
}

type OpTerm struct {
	OP   string  `parser:"( @BoolOp |"`
	K    *string `parser:"@Proxim )"`
	Term *Term   `parser:"@@"`
}

type Expression struct {
	Left  *Term     `parser:"@@"`
	Right []*OpTerm `parser:"@@*"`
}

//goland:noinspection GoVetStructTag
type Phrase struct {
	Value []*string `"\"" @Ident* "\""`
}

var lexer = stateful.MustSimple([]stateful.Rule{
	{"BoolOp", `(?i)AND\sNOT|AND|OR`, nil},
	{"Ident", `[a-zA-Z_]\w*`, nil},
	{"String", `"(\\"|[^"])*"`, nil},
	{`Proxim`, `/\d`, nil}, // MUST be over 'Punct'
	{"Punct", `[-[!@#$%^&*()+_={}\|:;"'<,>.?/]|]`, nil},
	{"Number", `[-+]?(\d*\.)?\d+`, nil},
	{"EOL", `[\n\r]+`, nil},
	{"whitespace", `[ \t]+`, nil},
})

var parser = participle.MustBuild(&Expression{},
	participle.Lexer(lexer),
	//participle.UseLookahead(3),
)

func Parse(query string) (*Expression, error) {
	gram := &Expression{}
	err := parser.ParseString("", query, gram)
	if err != nil {
		return nil, fmt.Errorf("cannot parse string: %v", err)
	}
	return gram, nil
}

func Print(r *Expression) {
	repr.Println(r, repr.Indent(" "), repr.OmitEmpty(true))
}
