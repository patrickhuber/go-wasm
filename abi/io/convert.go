package io

func IntToBool(i uint32) bool {
	return i != 0
}

func UInt32ToTaggedCodeUnits(i uint32) TaggedCodeUnits {
	return TaggedCodeUnits{
		CodeUnits: i &^ UTF16Tag,
		UTF16:     i&UTF16Tag != 0,
	}
}

func TaggedCodeUnitsToUint32(tcu TaggedCodeUnits) uint32 {
	i := tcu.CodeUnits
	if tcu.UTF16 {
		i = i | UTF16Tag
	}
	return i
}
