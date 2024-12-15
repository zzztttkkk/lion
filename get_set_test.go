package reflectx

import (
	"reflect"
	"testing"
	"unsafe"
)

type A struct {
	A1 string `json:"a1"`
}

func noop(v any) {
	_ = v.(*string)
}

var a1offset int64
var a1Field reflect.StructField

func init() {
	a1Field, _ = reflect.TypeOf(A{}).FieldByName("A1")

	var ptr = &A{}
	a1offset = int64(uintptr(unsafe.Pointer(&ptr.A1))) - int64(uintptr(unsafe.Pointer(ptr)))
}

func BenchmarkNormalGetFieldPtrByIndex(b *testing.B) {
	val := A{}

	vv := reflect.ValueOf(&val).Elem()
	for i := 0; i < b.N; i++ {
		fv := vv.FieldByIndex(a1Field.Index).Addr().Interface()
		noop(fv)
	}
}

func BenchmarkGetFieldPtrByOffset(b *testing.B) {
	val := A{}
	valptr := unsafe.Pointer(&val)

	for i := 0; i < b.N; i++ {
		fv := reflect.NewAt(a1Field.Type, unsafe.Add(valptr, a1offset)).Interface()
		noop(fv)
	}
}

func BenchmarkGetFieldPtrByOffsetAndTypecast(b *testing.B) {
	val := A{}
	valptr := unsafe.Pointer(&val)

	for i := 0; i < b.N; i++ {
		fv := (*string)(unsafe.Add(valptr, a1offset))
		noop(fv)
	}
}
