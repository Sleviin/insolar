// Code generated by "stringer -type=ReactivateMode"; DO NOT EDIT.

package slot

import "strconv"

const _ReactivateMode_name = "EmptyResponseTickSeqHead"

var _ReactivateMode_index = [...]uint8{0, 5, 13, 17, 24}

func (i ReactivateMode) String() string {
	if i < 0 || i >= ReactivateMode(len(_ReactivateMode_index)-1) {
		return "ReactivateMode(" + strconv.FormatInt(int64(i), 10) + ")"
	}
	return _ReactivateMode_name[_ReactivateMode_index[i]:_ReactivateMode_index[i+1]]
}
