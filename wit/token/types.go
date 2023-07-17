package token

type TokenType string

const (
	None        TokenType = "nil"
	OpenParen   TokenType = "("
	CloseParen  TokenType = ")"
	String      TokenType = "\\w+"
	Whitespace  TokenType = "\\s+"
	EndOfStream TokenType = "EOF"
	Use         TokenType = "use"
	Type        TokenType = "type"
	Resource    TokenType = "resource"
	Func        TokenType = "func"
	Record      TokenType = "record"
	Enum        TokenType = "enum"
	Flags       TokenType = "flags"
	Variant     TokenType = "variant"
	Union       TokenType = "union"
	Static      TokenType = "static"
	Interface   TokenType = "interface"
	World       TokenType = "world"
	Import      TokenType = "import"
	Export      TokenType = "export"
	Package     TokenType = "package"
	Include     TokenType = "include"
)
