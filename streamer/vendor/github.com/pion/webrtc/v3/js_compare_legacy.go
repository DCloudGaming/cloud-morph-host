// +build js,!go1.14

package webrtc

import "syscall/js"

func jsValueIsUndefined(v js.Value) bool {
	return v == js.Undefined()
}

func jsValueIsNull(v js.Value) bool {
	return v == js.Null()
}
