package io

const UTF16Tag = 1 << 31

type TaggedCodeUnits struct {
	CodeUnits uint32
	UTF16     bool
}

func UInt32ToTaggedCodeUnits(i uint32) TaggedCodeUnits {
	return TaggedCodeUnits{
		CodeUnits: i &^ UTF16Tag,
		UTF16:     i&UTF16Tag != 0,
	}
}

func (tcu TaggedCodeUnits) ToUInt32() uint32 {
	i := tcu.CodeUnits
	if tcu.UTF16 {
		i = i | UTF16Tag
	}
	return i
}
