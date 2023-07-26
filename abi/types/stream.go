package types

type Stream interface {
	ValType
	Element() ValType
	End() ValType
}

type StreamImpl struct {
	ValTypeImpl
	element ValType
	end     ValType
}

// Element implements Stream.
func (s *StreamImpl) Element() ValType {
	return s.element
}

// End implements Stream.
func (s *StreamImpl) End() ValType {
	return s.end
}

func NewStream(element ValType, end ValType) Stream {
	return &StreamImpl{
		element: element,
		end:     end,
	}
}
