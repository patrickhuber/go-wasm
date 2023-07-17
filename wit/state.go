package wit

type State int

const (
	None State = iota
	BeginComment
	BlockComment
	BlockCommentStar
	LineComment
	String
	WhiteSpace
)
