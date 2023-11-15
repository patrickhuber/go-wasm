package token

type Type string

const (
	LineComment  Type = "line_comment"
	BlockComment Type = "block_comment"
	Whitespace   Type = "whitespace"
	OpenParen    Type = "open_paren"
	CloseParen   Type = "close_paren"
	String       Type = "string"
	Id           Type = "id"
	Keyword      Type = "keyword"
	Reserved     Type = "reserved"
	Integer      Type = "integer"
	Float        Type = "float"
	EndOfStream  Type = "eof"
	Module       Type = "module"
	Component    Type = "component"
)

type Token struct {
	Type     Type
	Position int
	Column   int
	Line     int
	Capture  string
}
