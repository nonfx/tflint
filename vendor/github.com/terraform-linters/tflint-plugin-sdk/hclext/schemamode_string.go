// Code generated by "stringer -type=SchemaMode"; DO NOT EDIT.

package hclext

import "strconv"

func _() {
	// An "invalid array index" compiler error signifies that the constant values have changed.
	// Re-run the stringer command to generate them again.
	var x [1]struct{}
	_ = x[SchemaDefaultMode-0]
	_ = x[SchemaJustAttributesMode-1]
}

const _SchemaMode_name = "SchemaDefaultModeSchemaJustAttributesMode"

var _SchemaMode_index = [...]uint8{0, 17, 41}

func (i SchemaMode) String() string {
	if i < 0 || i >= SchemaMode(len(_SchemaMode_index)-1) {
		return "SchemaMode(" + strconv.FormatInt(int64(i), 10) + ")"
	}
	return _SchemaMode_name[_SchemaMode_index[i]:_SchemaMode_index[i+1]]
}
