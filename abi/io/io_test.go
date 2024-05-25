package io_test

import (
	"bytes"
	"errors"
	"fmt"
	"math"
	"reflect"
	"strconv"
	"testing"

	"github.com/patrickhuber/go-wasm/abi/io"
	"github.com/patrickhuber/go-wasm/abi/kind"
	"github.com/patrickhuber/go-wasm/abi/types"
	"github.com/patrickhuber/go-wasm/abi/values"
	"github.com/patrickhuber/go-wasm/encoding"
	"github.com/stretchr/testify/require"
)

func Test(t *testing.T) {
	type testCase struct {
		name       string
		t          types.ValType
		valsToLift []any
		v          any
	}
	tests := []testCase{
		{"record", Record(), []any{}, map[string]any{}},
		{"record", Record(Field("x", U8()), Field("y", U16()), Field("z", U32())),
			[]any{uint32(1), uint32(2), uint32(3)}, map[string]any{"x": uint8(1), "y": uint16(2), "z": uint32(3)}},
		{"tuple", Tuple(
			Tuple(U8(), U8()),
			U8()), []any{uint32(1), uint32(2), uint32(3)}, map[string]any{"0": map[string]any{"0": uint8(1), "1": uint8(2)}, "1": uint8(3)}},
		{"flags", Flags(), []any{}, map[string]any{}},
		{"flags", Flags("a", "b"), []any{uint32(0)}, map[string]any{"a": false, "b": false}},
		{"flags", Flags("a", "b"), []any{uint32(2)}, map[string]any{"a": false, "b": true}},
		{"flags", Flags("a", "b"), []any{uint32(3)}, map[string]any{"a": true, "b": true}},
		{"flags", Flags("a", "b"), []any{uint32(4)}, map[string]any{"a": false, "b": false}},
		{"flags", Flags(Apply(Range(0, 33), strconv.Itoa)...), []any{uint32(math.MaxUint32), uint32(0x1)}, Zip(Apply(Range(0, 33), strconv.Itoa), Repeat[any](true, 33))},
		{"variant", Variant(Case("x", U8()), Case("y", Float32()), Case("z", nil)), []any{uint32(0), uint32(42)}, map[string]any{"x": uint8(42)}},
		{"variant", Variant(Case("x", U8()), Case("y", Float32()), Case("z", nil)), []any{uint32(0), uint32(256)}, map[string]any{"x": uint8(0)}},
		{"variant", Variant(Case("x", U8()), Case("y", Float32()), Case("z", nil)), []any{uint32(1), uint32(0x4048f5c3)}, map[string]any{"y": float32(3.140000104904175)}},
		{"variant", Variant(Case("x", U8()), Case("y", Float32()), Case("z", nil)), []any{uint32(2), uint32(0xffffffff)}, map[string]any{"z": nil}},
		{"union", Union(U32(), U64()), []any{uint32(0), uint64(42)}, map[string]any{"0": uint32(42)}},
		{"union", Union(U32(), U64()), []any{uint32(0), uint64(1 << 35)}, map[string]any{"0": uint32(0)}},
		{"union", Union(U32(), U64()), []any{uint32(1), uint64(1 << 35)}, map[string]any{"1": uint64(1 << 35)}},
		{"union", Union(Float32(), U64()), []any{uint32(0), uint64(0x4048f5c3)}, map[string]any{"0": float32(3.140000104904175)}},
		{"union", Union(Float32(), U64()), []any{uint32(0), uint64(1 << 35)}, map[string]any{"0": float32(0)}},
		{"union", Union(Float32(), U64()), []any{uint32(1), uint64(1 << 35)}, map[string]any{"1": uint64(1 << 35)}},
		{"union", Union(Float64(), U64()), []any{uint32(0), uint64(0x40091EB851EB851F)}, map[string]any{"0": float64(3.14)}},
		{"union", Union(Float64(), U64()), []any{uint32(0), uint64(1 << 35)}, map[string]any{"0": float64(1.69759663277e-313)}},
		{"union", Union(Float64(), U64()), []any{uint32(1), uint64(1 << 35)}, map[string]any{"1": uint64(1 << 35)}},
		{"union", Union(U8()), []any{uint32(0), uint32(42)}, map[string]any{"0": uint8(42)}},
		{"union", Union(U8()), []any{uint32(1), uint32(256)}, nil},
		{"union", Union(U8()), []any{uint32(0), uint32(256)}, map[string]any{"0": uint8(0)}},
		{"option", Option(Float32()), []any{uint32(0), float32(3.14)}, map[string]any{"none": nil}},
		{"option", Option(Float32()), []any{uint32(1), float32(3.14)}, map[string]any{"some": float32(3.14)}},
		{"result", Result(U8(), U32()), []any{uint32(0), uint32(42)}, map[string]any{"ok": uint8(42)}},
		{"result", Result(U8(), U32()), []any{uint32(1), uint32(1000)}, map[string]any{"error": uint32(1000)}},
	}
	vt := Variant(
		Case("w", U8()),
		CaseWith("x", U8(), "w"),
		Case("y", U8()),
		CaseWith("z", U8(), "x"))
	tests = append(tests,
		testCase{"variant", vt, []any{uint32(0), uint32(42)}, map[string]any{"w": uint8(42)}},
		testCase{"variant", vt, []any{uint32(1), uint32(42)}, map[string]any{"x|w": uint8(42)}},
		testCase{"variant", vt, []any{uint32(2), uint32(42)}, map[string]any{"y": uint8(42)}},
		testCase{"variant", vt, []any{uint32(3), uint32(42)}, map[string]any{"z|x|w": uint8(42)}},
	)

	for _, oneTest := range tests {
		t.Run(oneTest.name, func(t *testing.T) {
			cxt := Context()
			err := test(oneTest.t, oneTest.valsToLift, oneTest.v, cxt, encoding.UTF8, nil, nil)
			require.Nil(t, err)
		})
	}
}

func TestWithLower(t *testing.T) {
	type testCase struct {
		name       string
		t          types.ValType
		valsToLift []any
		v          any
		lowerT     types.ValType
		lowerV     any
	}
	vt := Variant(
		Case("w", U8()),
		CaseWith("x", U8(), "w"),
		Case("y", U8()),
		CaseWith("z", U8(), "x"))

	vt2 := Variant(Case("w", U8()))
	tests := []testCase{
		{"variant", vt, []any{uint32(0), uint32(42)}, map[string]any{"w": uint8(42)}, vt2, map[string]any{"w": uint8(42)}},
		{"variant", vt, []any{uint32(1), uint32(42)}, map[string]any{"x|w": uint8(42)}, vt2, map[string]any{"w": uint8(42)}},
		{"variant", vt, []any{uint32(3), uint32(42)}, map[string]any{"z|x|w": uint8(42)}, vt2, map[string]any{"w": uint8(42)}},
	}
	for _, oneTest := range tests {
		t.Run(oneTest.name, func(t *testing.T) {
			cxt := Context() // &bytes.Buffer{}, encoding.UTF8, nil, nil)
			err := test(oneTest.t, oneTest.valsToLift, oneTest.v, cxt, encoding.UTF8, oneTest.lowerT, oneTest.lowerV)
			require.Nil(t, err)
		})
	}
}

type ContextOption func(*types.CallContext)
type CanonicalOptionsOption func(*types.CanonicalOptions)
type ComponentInstanceOption func(*types.ComponentInstance)

func CanonicalOptions(options ...CanonicalOptionsOption) ContextOption {
	return func(cc *types.CallContext) {
		if cc == nil {
			return
		}
		if cc.Options == nil && len(options) > 0 {
			cc.Options = Options()
		}
		for _, op := range options {
			op(cc.Options)
		}
	}
}

func ComponentInstance(options ...ComponentInstanceOption) ContextOption {
	return func(cc *types.CallContext) {
		if cc == nil || cc.Instance == nil {
			return
		}
		if cc.Instance == nil && len(options) > 0 {
			cc.Instance = Instance()
		}
		for _, op := range options {
			op(cc.Instance)
		}
	}
}

func Memory(memory *bytes.Buffer) CanonicalOptionsOption {
	return func(op *types.CanonicalOptions) {
		op.Memory = memory
	}
}

func Encoding(enc encoding.Encoding) CanonicalOptionsOption {
	return func(op *types.CanonicalOptions) {
		op.StringEncoding = enc
	}
}

func Realloc(realloc types.ReallocFunc) CanonicalOptionsOption {
	return func(op *types.CanonicalOptions) {
		op.Realloc = realloc
	}
}

func PostReturn(postReturn types.PostReturnFunc) CanonicalOptionsOption {
	return func(op *types.CanonicalOptions) {
		op.PostReturn = postReturn
	}
}

func Context(options ...ContextOption) *types.CallContext {
	cx := &types.CallContext{
		BorrowCount: 0,
		Options:     Options(),
		Instance:    Instance(),
	}
	for _, op := range options {
		op(cx)
	}
	return cx
}

func Options(options ...CanonicalOptionsOption) *types.CanonicalOptions {
	opt := &types.CanonicalOptions{
		StringEncoding: encoding.UTF8,
	}
	for _, op := range options {
		op(opt)
	}
	if opt.Memory == nil {
		opt.Memory = &bytes.Buffer{}
	}
	return opt
}

func Instance(options ...ComponentInstanceOption) *types.ComponentInstance {
	inst := &types.ComponentInstance{
		MayEnter: true,
		MayLeave: true,
		Handles: types.HandleTables{
			ResourceTypeToTable: map[types.ResourceType]*types.HandleTable{},
		},
	}
	for _, op := range options {
		op(inst)
	}
	return inst
}

func Bool() types.Bool { return types.NewBool() }

func S8() types.S8   { return types.NewS8() }
func S16() types.S16 { return types.NewS16() }
func S32() types.S32 { return types.NewS32() }
func S64() types.S64 { return types.NewS64() }

func U8() types.U8   { return types.NewU8() }
func U16() types.U16 { return types.NewU16() }
func U32() types.U32 { return types.NewU32() }
func U64() types.U64 { return types.NewU64() }

func Float32() types.F32 { return types.NewF32() }
func Float64() types.F64 { return types.NewF64() }

func Char() types.Char     { return types.NewChar() }
func String() types.String { return types.NewString() }

func Enum(labels ...string) types.Enum {
	return types.NewEnum(labels...)
}

func Flags(labels ...string) types.Flags {
	return types.NewFlags(labels...)
}

func Tuple(t ...types.ValType) types.Tuple {
	return types.NewTuple(t...)
}

func Variant(c ...types.Case) types.Variant {
	return types.NewVariant(c...)
}

func Case(label string, val types.ValType) types.Case {
	return types.NewCase(label, val)
}

func CaseWith(label string, val types.ValType, refines string) types.Case {
	return types.NewCaseRefines(label, val, refines)
}

func Union(valTypes ...types.ValType) types.Union {
	return types.NewUnion(valTypes...)
}

func Option(valType types.ValType) types.Option {
	return types.NewOption(valType)
}

func Result(ok, err types.ValType) types.Result {
	return types.NewResult(ok, err)
}

func Record(fields ...types.Field) types.Record {
	return types.NewRecord(fields...)
}

func Field(label string, t types.ValType) types.Field {
	return types.Field{
		Label: label,
		Type:  t,
	}
}

func List(vt types.ValType) types.List {
	return types.NewList(vt)
}

func FuncType(params []types.ValType, results []types.ValType) types.FuncType {
	toParameters := func(valTypes []types.ValType) []types.Parameter {
		parameters := []types.Parameter{}
		for i, vt := range valTypes {
			param := types.Parameter{Name: strconv.Itoa(i), Type: vt}
			parameters = append(parameters, param)
		}
		return parameters
	}
	return types.NewFuncType(toParameters(params), toParameters(results))
}

func Own(rt types.ResourceType) types.Own {
	return types.NewOwn(rt)
}

func Borrow(rt types.ResourceType) types.Borrow {
	return types.NewBorrow(rt)
}

func Range(low int, high int) []int {
	var result []int
	for i := low; i < high; i++ {
		result = append(result, i)
	}
	return result
}

// Cross creates the cross product of the two slices as a slice of maps
func Cross[TKey comparable, TValue any](keys []TKey, values []TValue) []any {
	results := []any{}

	for _, value := range values {
		result := map[TKey]TValue{}
		for _, key := range keys {
			result[key] = value
		}
		results = append(results, result)
	}
	return results
}

func Repeat[TValue any](value TValue, times int) []TValue {
	var result []TValue
	for i := 0; i < times; i++ {
		result = append(result, value)
	}
	return result
}

func Apply[TInput any, TOutput any](input []TInput, transform func(k TInput) TOutput) []TOutput {
	var slice []TOutput
	for _, i := range input {
		slice = append(slice, transform(i))
	}
	return slice
}

func Zip[TKey comparable, TValue any](keys []TKey, values []TValue) map[TKey]TValue {
	length := len(keys)
	if len(values) < length {
		length = len(values)
	}
	result := make(map[TKey]TValue)
	for i := 0; i < length; i++ {
		result[keys[i]] = values[i]
	}
	return result
}

func Slice[T any](values ...T) []T {
	return values
}

func test(t types.ValType, valsToLift []any, v any,
	cx *types.CallContext,
	dstEncoding encoding.Encoding, lowerT types.ValType, lowerV any) error {

	flattened, err := io.FlattenType(t)
	if err != nil {
		return err
	}
	vs, err := zip(flattened, valsToLift)
	if err != nil {
		return err
	}

	vi := values.NewIterator(vs...)

	// this error handling logic is strange
	// but basically if v is nil,
	// handle the return from the function differently because we are expectig failure
	got, err := io.LiftFlat(cx, vi, t)
	if v == nil {
		if errors.Is(err, types.Trap()) {
			return nil
		}
		return fmt.Errorf("expected trap, but found %v", got)
	}

	if err != nil {
		return err
	}

	if vi.Index() != vi.Length() {
		return types.TrapWith("value iterator index %d exceeds length %d", vi.Index(), vi.Length())
	}
	if !reflect.DeepEqual(got, v) {
		return fmt.Errorf("initial lift_flat() expected %v but got %v", v, got)
	}

	lowerT = coalesce(lowerT, t)
	lowerV = coalesce(lowerV, v)

	heap := NewHeap(5 * cx.Options.Memory.Len())
	if dstEncoding == encoding.None {
		dstEncoding = cx.Options.StringEncoding
	}

	cx = Context(
		CanonicalOptions(
			Memory(heap.Memory),
			Encoding(dstEncoding),
			Realloc(heap.ReAllocate),
			PostReturn(cx.Options.PostReturn)))

	loweredValues, err := io.LowerFlat(cx, v, lowerT)
	if err != nil {
		return err
	}
	// assert here with lowerT

	vi = values.NewIterator(loweredValues...)
	got, err = io.LiftFlat(cx, vi, lowerT)
	if err != nil {
		return err
	}
	if !equalModuloStringEncoding(got, lowerV) {
		return fmt.Errorf("re-lift expected %v but got %v", lowerV, got)
	}
	return nil
}

// zip emulates python zip but only used here
func zip(types []kind.Kind, v []any) ([]values.Value, error) {
	vs := []values.Value{}

	if len(v) != len(types) {
		return nil, fmt.Errorf("expected len(values)=%d to equal len(types)=%d", len(v), len(types))
	}
	for i := 0; i < len(v); i++ {
		t := types[i]
		var vals []values.Value
		var err error
		switch t {
		case kind.U32:
			vals, err = io.LowerU32(v[i])
		case kind.U64:
			vals, err = io.LowerU64(v[i])
		case kind.Float32:
			vals, err = io.LowerFloat32(v[i])
		case kind.Float64:
			vals, err = io.LowerFloat64(v[i])
		default:
			err = fmt.Errorf("invalid kind.%s", t.String())
		}
		if err != nil {
			return nil, err
		}
		vs = append(vs, vals...)
	}
	return vs, nil
}

func equalModuloStringEncoding(s any, v any) bool {
	if s == nil && v == nil {
		return true
	}
	if s == nil {
		return false
	}
	if v == nil {
		return false
	}

	return reflect.DeepEqual(s, v)
}

func coalesce[T comparable](v T, other ...T) T {
	var zero T
	if v != zero {
		return v
	}
	for _, o := range other {
		if o != zero {
			return o
		}
	}
	return zero
}
