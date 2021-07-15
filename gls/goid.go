package gls

import (
	reflect "github.com/modern-go/reflect2"
	"unsafe"
)

// offset for go1.4
var goidOffset uintptr = 128

func init() {
	gType := reflect.TypeByName("runtime.g").(reflect.StructType)
	if gType == nil {
		panic("failed to get runtime.g type")
	}
	goidField := gType.FieldByName("goid")
	goidOffset = goidField.Offset()
}

// GoID returns the goroutine id of current goroutine
func GoID() int64 {
	g := getg()
	p_goid := (*int64)(unsafe.Pointer(g + goidOffset))
	return *p_goid
}

func getg() uintptr
