package token

type TokenType string

const (
	ILLEGAL = "ILLEGAL"
	EOF     = "EOF"

	//Identifiers + literals
	IDENT = "IDENT"
	INT   = "INT"

	//operators
	ASSIGN   = "="
	PLUS     = "+"
	MINUS    = "-"
	BANG     = "!"
	ASTERISK = "*"
	SLASH    = "/"

	LT = "<"
	GT = ">"

	//bools
	EQ     = "=="
	NOT_EQ = "!="

	//delimiters
	COMMA     = ","
	SEMICOLON = ";"

	LPAREN = "("
	RPAREN = ")"
	LBRACE = "{"
	RBRACE = "}"

	//keywords
	FUNCTION = "FUNCTION"
	LET      = "LET"
	FALSE    = "false"
	TRUE     = "true"
	RETURN   = "return"
	IF       = "if"
	ELSE     = "else"
)

type Token struct {
	Type    TokenType
	Literal string
}

var keywords = map[string]TokenType{
	"let":    LET,
	"fn":     FUNCTION,
	"false":  FALSE,
	"true":   TRUE,
	"return": RETURN,
	"if":     IF,
	"else":   ELSE,
}

func LookupIdent(input string) TokenType {
	if tok, ok := keywords[input]; ok {
		return tok
	}

	return IDENT
}
