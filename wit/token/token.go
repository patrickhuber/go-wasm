package token

type Token struct {
	Type     TokenType
	Position int
	Column   int
	Line     int
	Capture  string
}
