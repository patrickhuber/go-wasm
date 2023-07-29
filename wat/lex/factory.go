package lex

type Factory interface {
	Lexeme(Rule) Lexeme
}
