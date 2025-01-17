package lion

import (
	"fmt"
	"reflect"
	"testing"
	"unsafe"
)

type GTInt int64

type A struct {
	A1 string `json:"a1"`
	A2 *GTInt
	A3 []GTInt
}

func TestGetFieldPtr(t *testing.T) {
	mptr := Ptr[A]()

	a1f := FieldOf[A](&mptr.A1)
	a2f := FieldOf[A](&mptr.A2)
	a3f := FieldOf[A](&mptr.A3)

	obj := A{
		A1: "a1",
		A2: new(GTInt),
		A3: []GTInt{7, 888, 9},
	}
	*obj.A2 = 12

	fmt.Println(*(a1f.PtrOf(unsafe.Pointer(&obj)).(*string)))
	fmt.Println(**(a2f.PtrOf(unsafe.Pointer(&obj)).(**GTInt)))
	fmt.Println(a3f.PtrOf(unsafe.Pointer(&obj)).(*[]GTInt))
}

func noopptr(v any) {
	_ = v.(*string)
}

func noopval(v any) {
	_ = v.(string)
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
		noopptr(fptr)
	}
}

func BenchmarkGetFieldValueByNormalReflect(b *testing.B) {
	val := A{}
	vv := reflect.ValueOf(&val).Elem()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		fptr := vv.FieldByIndex(a1Field.Index).Interface()
		noopval(fptr)
	}
}

func BenchmarkGetFieldPtrByOffset(b *testing.B) {
	val := A{}
	valptr := unsafe.Pointer(&val)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		fptr := reflect.NewAt(a1Field.Type, unsafe.Add(valptr, a1offset)).Interface()
		noopptr(fptr)
	}
}

func BenchmarkGetFieldPtrByMethod(b *testing.B) {
	a1f := FieldOf[A](&(Ptr[A]().A1))
	val := A{}
	valptr := unsafe.Pointer(&val)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		fptr := a1f.PtrOf(valptr)
		noopptr(fptr)
	}
}

func BenchmarkGetFieldPtrByOffsetAndTypecast(b *testing.B) {
	val := A{}
	valptr := unsafe.Pointer(&val)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		fptr := (*string)(unsafe.Add(valptr, a1offset))
		noopptr(fptr)
	}
}

func BenchmarkGetFieldPtrDirectly(b *testing.B) {
	val := A{}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		noopptr(&val.A1)
	}
}
