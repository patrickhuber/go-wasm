package types

type Field struct {
	Label string
	Type  ValType
}

type Record interface {
	ValType
	Fields() []Field
	record()
}

type RecordImpl struct {
	ValTypeImpl
	fields []Field
}

func (*RecordImpl) record() {}

func (r *RecordImpl) Fields() []Field {
	return r.fields
}

func NewRecord(fields ...Field) Record {
	return &RecordImpl{
		fields: fields,
	}
}
