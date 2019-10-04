// Code generated by "stringer -type=ErrorCode"; DO NOT EDIT.

package payload

import "strconv"

func _() {
	// An "invalid array index" compiler error signifies that the constant values have changed.
	// Re-run the stringer command to generate them again.
	var x [1]struct{}
	_ = x[CodeUnknown-0]
	_ = x[CodeDeactivated-1]
	_ = x[CodeFlowCanceled-2]
	_ = x[CodeNotFound-3]
	_ = x[CodeNoPendings-4]
	_ = x[CodeNoStartPulse-5]
	_ = x[CodeRequestNotFound-6]
	_ = x[CodeRequestInvalid-7]
	_ = x[CodeRequestNonClosedOutgoing-8]
	_ = x[CodeRequestNonOldestMutable-9]
	_ = x[CodeReasonIsWrong-10]
	_ = x[CodeNonActivated-11]
	_ = x[CodeLoopDetected-12]
}

const _ErrorCode_name = "CodeUnknownCodeDeactivatedCodeFlowCanceledCodeNotFoundCodeNoPendingsCodeNoStartPulseCodeRequestNotFoundCodeRequestInvalidCodeRequestNonClosedOutgoingCodeRequestNonOldestMutableCodeReasonIsWrongCodeNonActivatedCodeLoopDetected"

var _ErrorCode_index = [...]uint8{0, 11, 26, 42, 54, 68, 84, 103, 121, 149, 176, 193, 209, 225}

func (i ErrorCode) String() string {
	if i >= ErrorCode(len(_ErrorCode_index)-1) {
		return "ErrorCode(" + strconv.FormatInt(int64(i), 10) + ")"
	}
	return _ErrorCode_name[_ErrorCode_index[i]:_ErrorCode_index[i+1]]
}
