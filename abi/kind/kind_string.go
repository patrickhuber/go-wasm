// Code generated by "stringer -type=Kind"; DO NOT EDIT.

package kind

import "strconv"

func _() {
	// An "invalid array index" compiler error signifies that the constant values have changed.
	// Re-run the stringer command to generate them again.
	var x [1]struct{}
	_ = x[U32-0]
	_ = x[U64-1]
	_ = x[Float32-2]
	_ = x[Float64-3]
}

const _Kind_name = "U32U64Float32Float64"

var _Kind_index = [...]uint8{0, 3, 6, 13, 20}

func (i Kind) String() string {
	if i < 0 || i >= Kind(len(_Kind_index)-1) {
		return "Kind(" + strconv.FormatInt(int64(i), 10) + ")"
	}
	return _Kind_name[_Kind_index[i]:_Kind_index[i+1]]
}
