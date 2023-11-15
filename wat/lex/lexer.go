package lex

import (
	"fmt"
	"strings"

	"github.com/patrickhuber/go-types"
	"github.com/patrickhuber/go-types/handle"
	"github.com/patrickhuber/go-types/option"
	"github.com/patrickhuber/go-types/result"
	"github.com/patrickhuber/go-wasm/wat/token"
)

// Lexer consists of rules that are enumerated during a prefix search.
// Rules are greedy so they are matched until they can no longer match.
// At that point, a token is emitted
type Lexer struct {
	Rules     []Rule
	input     string
	offset    int
	position  int
	column    int
	line      int
	peekToken *token.Token
}

var runeMap = map[byte]token.Type{
	'(': token.OpenParen,
	')': token.CloseParen,
}

func New(input string) *Lexer {
	// order matters here
	// two rules could match the same input of the same length
	// if that happens the one listed first wins
	rules := []Rule{
		whitespace(),
		lineComment(),
		blockComment(),
		str(),
		identifier(),
		reserved(),
	}
	for r, ty := range runeMap {
		rules = append(rules, &StringRule{
			String:    string(r),
			TokenType: ty,
		})
	}
	return &Lexer{
		input: input,
		Rules: rules,
	}
}

func (l *Lexer) Clone() *Lexer {
	return &Lexer{
		input:     l.input,
		offset:    l.offset,
		position:  l.position,
		column:    l.column,
		line:      l.line,
		peekToken: l.peekToken,
		Rules:     l.Rules,
	}
}

func (l *Lexer) Peek() (*token.Token, error) {
	// always return the peek token if it exists
	if l.peekToken != nil {
		return l.peekToken, nil
	}

	l.peekToken = l.next().Unwrap()
	return l.peekToken, nil
}

func (l *Lexer) Next() (*token.Token, error) {
	// any peek token?
	if l.peekToken == nil {
		return l.next().Deconstruct()
	}

	tok := l.peekToken
	l.peekToken = nil

	return tok, nil
}

func (l *Lexer) next() (res types.Result[*token.Token]) {
	defer handle.Error(&res)

	r, ok := l.readRune().Deconstruct()

	if !ok {
		return l.token(token.EndOfStream)
	}

	var lexemes []Lexeme
	// find any matching rules
	for _, rule := range l.Rules {
		if rule.Check(r) {
			factory := registry[rule.RuleType()]
			lexeme := factory.Lexeme(rule)
			lexeme.Eat(r)
			lexemes = append(lexemes, lexeme)
		}
	}

	if len(lexemes) == 0 {
		return result.Error[*token.Token](l.lexerError())
	}

	for {
		r, ok := l.peekRune().Deconstruct()
		if !ok {
			for _, lexeme := range lexemes {
				// emit any accepted state
				if lexeme.Accepted() {
					return l.token(lexeme.Rule().Type())
				}
			}

			// process any lexemes in an accepted state and emit a token
			return l.token(token.EndOfStream) // returned if nothing is matched
		}

		var matched []Lexeme
		for _, lexeme := range lexemes {
			if lexeme.Eat(r) {
				matched = append(matched, lexeme)
			}
		}

		if len(matched) == 0 {
			for _, lexeme := range lexemes {
				// emit any accepted state
				if lexeme.Accepted() {
					return l.token(lexeme.Rule().Type())
				}
			}
			return result.Error[*token.Token](l.lexerError())
		}

		_ = l.readRune().Unwrap()

		// move only winners
		lexemes = matched
	}
}

func (l *Lexer) token(ty token.Type) types.Result[*token.Token] {

	// snapshot the state for the current token
	tok := &token.Token{
		Type:     ty,
		Position: l.offset,
		Column:   l.column,
		Line:     l.line,
		Capture:  l.input[l.offset:l.position],
	}

	// classify the reserved token
	if tok.Type == token.Reserved {
		if isInteger(tok.Capture) {
			tok.Type = token.Integer
		} else if isFloat(tok.Capture) {
			tok.Type = token.Float
		} else if kw, ok := keyword(tok.Capture); ok {
			tok.Type = kw
		}
	}

	// fast forward updating metrics
	for i := l.offset; i < l.position; i++ {
		ch := l.input[i]
		if ch == '\n' {
			l.line++
			l.column = 0
		} else {
			l.column++
		}
	}

	// update the current offset to the position
	l.offset = l.position

	return result.Ok(tok)
}

func (l *Lexer) peekRune() (op types.Option[byte]) {
	if l.position >= len(l.input) {
		return option.None[byte]()
	}
	r := l.input[l.position]
	return option.Some(r)
}

func (l *Lexer) readRune() types.Option[byte] {
	if l.position >= len(l.input) {
		return option.None[byte]()
	}
	r := l.input[l.position]
	l.position++
	return option.Some(r)
}

func whitespace() Rule {
	start := &Node{}
	end := &Node{
		Final: true,
	}
	start.Edges = append(start.Edges, &FuncEdge{
		Func: isSpace,
		Node: end,
	})
	end.Edges = append(end.Edges, &FuncEdge{
		Func: isSpace,
		Node: end,
	})
	return &DfaRule{
		Dfa: &Dfa{
			Start: start,
		},
		TokenType: token.Whitespace,
	}
}

func lineComment() Rule {
	start := &Node{}
	semi := &Node{}
	semi2 := &Node{}
	newLine := &Node{
		Final: true,
	}
	start.Edges = append(start.Edges, &ByteEdge{
		Byte: ';',
		Node: semi,
	})
	semi.Edges = append(semi.Edges, &ByteEdge{
		Byte: ';',
		Node: semi2,
	})
	semi2.Edges = append(semi2.Edges, &ByteEdge{
		Byte: '\n',
		Node: newLine,
	})
	semi2.Edges = append(semi2.Edges, &FuncEdge{
		Func: not('\n'),
		Node: semi2,
	})
	return &DfaRule{
		Dfa: &Dfa{
			Start: start,
		},
		TokenType: token.LineComment,
	}
}

func blockComment() Rule {
	start := &Node{}
	openParen := &Node{}
	semi := &Node{}
	semiEnd := &Node{}
	closeParen := &Node{Final: true}

	start.Edges = append(start.Edges, &ByteEdge{
		Byte: '(',
		Node: openParen,
	})
	openParen.Edges = append(openParen.Edges, &ByteEdge{
		Byte: ';',
		Node: semi,
	})
	semi.Edges = append(semi.Edges, &ByteEdge{
		Byte: ';',
		Node: semiEnd,
	}, &FuncEdge{
		Func: not(';'),
		Node: semi,
	})
	semiEnd.Edges = append(semiEnd.Edges, &ByteEdge{
		Byte: ')',
		Node: closeParen,
	}, &ByteEdge{
		Byte: ';',
		Node: semiEnd,
	}, &FuncEdge{
		Func: not(')', ';'),
		Node: semi,
	})

	return &DfaRule{
		Dfa: &Dfa{
			Start: start,
		},
		TokenType: token.BlockComment,
	}
}

func str() Rule {
	start := &Node{}
	doubleQuote := &Node{}
	end := &Node{Final: true}
	start.Edges = append(start.Edges, &ByteEdge{
		Byte: '"',
		Node: doubleQuote,
	})
	doubleQuote.Edges = append(doubleQuote.Edges, &ByteEdge{
		Byte: '"',
		Node: end,
	}, &FuncEdge{
		Func: not('"'),
		Node: doubleQuote,
	})
	return &DfaRule{
		Dfa: &Dfa{
			Start: start,
		},
		TokenType: token.String,
	}
}

func reserved() Rule {
	start := &Node{}
	idchar := &Node{Final: true}
	start.Edges = []Edge{
		&FuncEdge{Func: isIdChar, Node: idchar},
	}
	idchar.Edges = []Edge{
		&FuncEdge{Func: isIdChar, Node: idchar},
	}
	return &DfaRule{
		Dfa: &Dfa{
			Start: start,
		},
		TokenType: token.Reserved,
	}
}

func sign(s string, i int) (int, bool) {
	if i >= len(s) {
		return i, false
	}
	if s[i] == '+' || s[i] == '-' {
		return i + 1, true
	}
	return i, false
}

func num(s string, i int) (int, bool) {
	if i >= len(s) {
		return i, false
	}
	if !isDigit(s[i]) {
		return i, false
	}

	for ; i < len(s); i++ {
		if s[i] == '_' {
			continue
		}
		if isDigit(s[i]) {
			continue
		}
		return i, false
	}
	return i, true
}

func hexnum(s string, i int) (int, bool) {
	if i >= len(s) {
		return i, false
	}
	if !isHex(s[i]) {
		return i, false
	}
	for ; i < len(s); i++ {
		if s[i] == '_' {
			continue
		}
		if isHex(s[i]) {
			continue
		}
		return i, true
	}
	return i, true
}

func frac(s string, i int) (int, bool) {
	if i >= len(s) {
		return i, false
	}
	if !isDigit(s[i]) {
		return i, false
	}
	for ; i < len(s); i++ {
		if s[i] == '_' {
			continue
		}
		if isDigit(s[i]) {
			continue
		}
		return i, false
	}
	return i, true
}

func hexfrac(s string, i int) (int, bool) {
	if i >= len(s) {
		return i, false
	}
	if !isHex(s[i]) {
		return i, false
	}
	for ; i < len(s); i++ {
		if s[i] == '_' {
			continue
		}
		if isHex(s[i]) {
			continue
		}
		return i, true
	}
	return i, true
}

func float(s string, i int) (int, bool) {
	n, ok := num(s, i)
	if !ok {
		return i, false
	}
	if n == len(s) {
		return n, ok
	}
	i = n
	// p:num '.'?
	if s[i] == '.' && i == len(s)-1 {
		return i + 1, true
	}
	// p:num '.'?('E'|'e')('+'|'-')e:num
	if (s[i] == '.' && (s[i+1] == 'E' || s[i+1] == 'e')) ||
		(s[i] == 'e' || s[i] == 'E') {
		if s[i] == '.' {
			i++ // '.'
		}
		i++ // 'E' | 'e'
		n, ok = sign(s, i)
		if !ok {
			return i, false
		}
		i = n
		n, ok = num(s, i)
		if !ok {
			return i, false
		}
		return n, true
	}
	// p:num '.' q:frac
	n, ok = frac(s, i)
	if !ok {
		return i, false
	}
	i = n
	if i == len(s) {
		return i, true
	}
	// p:num '.' q:frac ('E'|'e')('+'|'-')e:num
	if s[i] != 'e' && s[i] != 'E' {
		return i, false
	}
	i++
	n, ok = sign(s, i)
	if !ok {
		return i, false
	}
	i = n
	return num(s, i)
}

func hexfloat(s string, i int) (int, bool) {
	if !strings.HasPrefix(s, "0x") {
		return i, false
	}
	i += 2
	n, ok := hexnum(s, i)
	if !ok {
		return i, false
	}
	if n == len(s) {
		return n, ok
	}
	i = n
	// p:hexnum '.'?
	if s[i] == '.' && i == len(s)-1 {
		return i + 1, true
	}
	// p:hexnum '.'?('P'|'p')('+'|'-')e:num
	if (s[i] == '.' && (s[i+1] == 'P' || s[i+1] == 'p')) ||
		(s[i] == 'p' || s[i] == 'P') {
		if s[i] == '.' {
			i++ // '.'
		}
		i++ // 'P' | 'p'
		n, ok = sign(s, i)
		if !ok {
			return i, false
		}
		i = n
		n, ok = num(s, i)
		if !ok {
			return i, false
		}
		return n, true
	}
	i++ // '.'
	// p:hexnum '.' q:hexfrac
	n, ok = hexfrac(s, i)
	if !ok {
		return i, false
	}
	i = n
	if i == len(s) {
		return i, true
	}
	// p:num '.' q:frac ('P'|'p')('+'|'-')e:num
	if s[i] != 'p' && s[i] != 'P' {
		return i, false
	}
	i++
	n, ok = sign(s, i)
	if !ok {
		return i, false
	}
	i = n
	return num(s, i)
}

func isFloat(s string) bool {
	if len(s) == 0 {
		return false
	}
	//negative := s[0] == '-'
	if s[0] == '+' || s[0] == '-' {
		s = s[1:]
	}
	if s == "inf" {
		return true
	}
	if s == "nan" {
		return true
	}
	if s == "nan:arithmetic" {
		return true
	}
	if s == "nan:canonical" {
		return true
	}
	if strings.HasPrefix(s, "nan:0x") {
		s = strings.TrimPrefix(s, "nan:0x")
		_, ok := hexnum(s, 0)
		return ok
	}
	if _, ok := float(s, 0); ok {
		return true
	}
	if _, ok := hexfloat(s, 0); ok {
		return true
	}
	return false
}

func isInteger(s string) bool {
	if len(s) == 0 {
		return false
	}
	// negative := s[0] == '-'
	if s[0] == '+' || s[0] == '-' {
		s = s[1:]
	}
	if _, ok := num(s, 0); ok {
		return true
	}
	if strings.HasPrefix(s, "0x") {
		s = strings.TrimPrefix(s, "0x")
		if i, ok := hexnum(s, 0); ok {
			// if we reached the end of the string this is an integer
			return i == len(s)
		}
	}
	return false
}

var keywordMap = map[string]token.Type{
	string(token.Module):    token.Module,
	string(token.Component): token.Component,
}

func keyword(s string) (token.Type, bool) {
	kw, ok := keywordMap[s]
	return kw, ok
}

// identifier ~ $([\w]|[^ ",;\[\]])+
func identifier() Rule {
	start := &Node{}
	dollar := &Node{}
	idchar := &Node{Final: true}
	start.Edges = []Edge{
		&ByteEdge{Byte: '$', Node: dollar},
	}
	dollar.Edges = []Edge{
		&FuncEdge{Func: isIdChar, Node: idchar},
	}
	idchar.Edges = []Edge{
		&FuncEdge{Func: isIdChar, Node: idchar},
	}
	return &DfaRule{
		Dfa: &Dfa{
			Start: start,
		},
		TokenType: token.Id,
	}
}

var idCharMap = map[byte]struct{}{
	'!':  {},
	'#':  {},
	'$':  {},
	'%':  {},
	'&':  {},
	'\'': {},
	'*':  {},
	'+':  {},
	'-':  {},
	'.':  {},
	'/':  {},
	':':  {},
	'<':  {},
	'=':  {},
	'>':  {},
	'?':  {},
	'@':  {},
	'\\': {},
	'^':  {},
	'_':  {},
	'`':  {},
	'|':  {},
	'~':  {},
}

func isIdChar(ch byte) bool {
	_, ok := idCharMap[ch]
	if ok {
		return true
	}
	switch {
	case isSpace(ch):
		return false
	case '0' <= ch && ch <= '9':
		return true
	case 'A' <= ch && ch <= 'Z':
		return true
	case 'a' <= ch && ch <= 'z':
		return true
	}
	return false
}

func isSpace(ch byte) bool {
	switch ch {
	case ' ':
	case '\t':
	case '\n':
	case '\f':
	case '\r':
	case '\v':
	default:
		return false
	}
	return true
}

func isHex(ch byte) bool {
	switch {
	case isDigit(ch):
	case 'a' <= ch && ch <= 'f':
	case 'A' <= ch && ch <= 'F':
	default:
		return false
	}
	return true
}

func isDigit(ch byte) bool {
	return '0' <= ch && ch <= '9'
}

func not(chars ...byte) func(ch byte) bool {
	return func(ch byte) bool {
		for _, r := range chars {
			if ch == r {
				return false
			}
		}
		return true
	}
}

func (l *Lexer) lexerError() error {
	return fmt.Errorf("error parsing at line: %d column: %d position: %d, '%s'", l.line, l.column, l.position, string(l.input[l.offset:l.position]))
}
