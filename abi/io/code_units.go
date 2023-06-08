package io

const UTF16Tag = 1 << 31

type TaggedCodeUnits struct {
	CodeUnits uint32
	UTF16     bool
}

