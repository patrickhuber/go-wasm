package lex

type Registry map[RuleType]Factory

// registry is used for processing incoming runes after a rule matches
var registry = Registry{
	StringRuleType: &StringFactory{},
	DfaRuleType:    &DfaFactory{},
}
