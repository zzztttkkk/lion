package lion

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

func BenchmarkGetFieldPtrByNormalReflect(b *testing.B) {
	val := A{}
	vv := reflect.ValueOf(&val).Elem()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		fptr := vv.FieldByIndex(a1Field.Index).Addr().Interface()
		noop(fptr)
	}
}

func BenchmarkGetFieldPtrByOffset(b *testing.B) {
	val := A{}
	valptr := unsafe.Pointer(&val)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		fptr := reflect.NewAt(a1Field.Type, unsafe.Add(valptr, a1offset)).Interface()
		noop(fptr)
	}
}

func BenchmarkGetFieldPtrByMethod(b *testing.B) {
	a1ptrgetter := FieldOf[A, struct{}](&(Ptr[A]().A1)).PtrGetter()
	val := A{}
	valptr := unsafe.Pointer(&val)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		fptr := a1ptrgetter(valptr)
		noop(fptr)
	}
}

func BenchmarkGetFieldPtrByOffsetAndTypecast(b *testing.B) {
	val := A{}
	valptr := unsafe.Pointer(&val)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		fptr := (*string)(unsafe.Add(valptr, a1offset))
		noop(fptr)
	}
}

func BenchmarkGetFieldPtrDirectly(b *testing.B) {
	val := A{}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		noop(&val.A1)
	}
}
