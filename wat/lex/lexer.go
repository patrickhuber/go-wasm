package lex

import (
	"fmt"
	"unicode"

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
	input     []rune
	offset    int
	position  int
	column    int
	line      int
	peekToken *token.Token
}

var runeMap = map[rune]token.Type{
	'(': token.OpenParen,
	')': token.CloseParen,
}

func New(input []rune) *Lexer {
	// order matters here
	// two rules could match the same input of the same length
	// if that happens the one listed first wins
	rules := []Rule{
		whitespace(),
		lineComment(),
		blockComment(),
		str(),
		integer(),
		identifier(),
		reserved(),
	}
	for r, ty := range runeMap {
		rules = append(rules, &RuneSliceRule{
			Runes:     []rune{r},
			TokenType: ty,
		})
	}
	return &Lexer{
		input: input,
		Rules: rules,
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
		Runes:    l.input[l.offset:l.position],
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

func (l *Lexer) peekRune() (op types.Option[rune]) {
	if l.position >= len(l.input) {
		return option.None[rune]()
	}
	r := l.input[l.position]
	return option.Some(r)
}

func (l *Lexer) readRune() types.Option[rune] {
	if l.position >= len(l.input) {
		return option.None[rune]()
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
		Func: unicode.IsSpace,
		Node: end,
	})
	end.Edges = append(end.Edges, &FuncEdge{
		Func: unicode.IsSpace,
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
	start.Edges = append(start.Edges, &RuneEdge{
		Rune: ';',
		Node: semi,
	})
	semi.Edges = append(semi.Edges, &RuneEdge{
		Rune: ';',
		Node: semi2,
	})
	semi2.Edges = append(semi2.Edges, &RuneEdge{
		Rune: '\n',
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

	start.Edges = append(start.Edges, &RuneEdge{
		Rune: '(',
		Node: openParen,
	})
	openParen.Edges = append(openParen.Edges, &RuneEdge{
		Rune: ';',
		Node: semi,
	})
	semi.Edges = append(semi.Edges, &RuneEdge{
		Rune: ';',
		Node: semiEnd,
	}, &FuncEdge{
		Func: not(';'),
		Node: semi,
	})
	semiEnd.Edges = append(semiEnd.Edges, &RuneEdge{
		Rune: ')',
		Node: closeParen,
	}, &RuneEdge{
		Rune: ';',
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
	start.Edges = append(start.Edges, &RuneEdge{
		Rune: '"',
		Node: doubleQuote,
	})
	doubleQuote.Edges = append(doubleQuote.Edges, &RuneEdge{
		Rune: '"',
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

// integer ~ [+-]?[0-9](_*[0-9])*
func integer() Rule {
	start := &Node{}
	plusOrMinus := &Node{}
	firstNumber := &Node{Final: true}
	underscore := &Node{}
	lastNumber := &Node{}

	// ( start ) -- [+|-] --> ( plusOrMinus )
	// ( start ) -- [0-9] --> ( firstNumber )
	start.Edges = append(start.Edges, &RuneEdge{
		Rune: '+',
		Node: plusOrMinus,
	}, &RuneEdge{
		Rune: '-',
		Node: plusOrMinus,
	}, &FuncEdge{
		Func: unicode.IsDigit,
		Node: firstNumber,
	})
	// ( plusOrMinus ) -- [0-9] --> ( firstNumber )
	plusOrMinus.Edges = append(plusOrMinus.Edges, &FuncEdge{
		Func: unicode.IsDigit,
		Node: firstNumber,
	})
	// ( firstNumber ) -- [0-9] --> ( lastNumber )
	// ( firstNumber ) -- _ --> ( underscore )
	firstNumber.Edges = append(firstNumber.Edges, &FuncEdge{
		Func: unicode.IsDigit,
		Node: lastNumber,
	}, &RuneEdge{
		Rune: '_',
		Node: underscore,
	})
	// ( underscore ) -- _ --> ( underscore )
	// ( underscore ) -- [0-9] --> ( lastNumber )
	underscore.Edges = append(underscore.Edges, &RuneEdge{
		Rune: '_',
		Node: underscore,
	}, &FuncEdge{
		Func: unicode.IsDigit,
		Node: lastNumber,
	})
	// ( lastNumber ) -- _ --> ( underscore )
	// ( lastNumber ) -- [0-9] --> ( lastNumber )
	lastNumber.Edges = append(lastNumber.Edges, &RuneEdge{
		Rune: '_',
		Node: underscore,
	}, &FuncEdge{
		Func: unicode.IsDigit,
		Node: lastNumber,
	})
	return &DfaRule{
		Dfa: &Dfa{
			Start: start,
		},
		TokenType: token.Integer,
	}
}

// identifier ~ $([\w]|[^ ",;\[\]])+
func identifier() Rule {
	start := &Node{}
	dollar := &Node{}
	idchar := &Node{Final: true}
	start.Edges = []Edge{
		&RuneEdge{Rune: '$', Node: dollar},
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

// hex ~ [+-]?0x[a-fA-F](_*[a-fA-F])*
func hex() Rule {
	return nil
}

var idCharMap = map[rune]struct{}{
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

func isIdChar(ch rune) bool {
	_, ok := idCharMap[ch]
	if ok {
		return true
	}
	switch {
	case unicode.IsSpace(ch):
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

func not(chars ...rune) func(ch rune) bool {
	return func(ch rune) bool {
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
