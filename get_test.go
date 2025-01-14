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

	a1f := FieldOf[A, struct{}](&mptr.A1)
	a2f := FieldOf[A, struct{}](&mptr.A2)
	a3f := FieldOf[A, struct{}](&mptr.A3)

	obj := A{
		A1: "a1",
		A2: new(GTInt),
		A3: []GTInt{7, 888, 9},
	}
	*obj.A2 = 12

	fmt.Println(*(a1f.PtrGetter()(unsafe.Pointer(&obj)).(*string)))
	fmt.Println(**(a2f.PtrGetter()(unsafe.Pointer(&obj)).(**GTInt)))
	fmt.Println(a3f.PtrGetter()(unsafe.Pointer(&obj)).(*[]GTInt))
}

func TestGetFieldValue(t *testing.T) {
	mptr := Ptr[A]()

	a1f := FieldOf[A, struct{}](&mptr.A1)
	a2f := FieldOf[A, struct{}](&mptr.A2)
	a3f := FieldOf[A, struct{}](&mptr.A3)

	obj := A{
		A1: "a1",
		A2: new(GTInt),
		A3: []GTInt{7, 888, 9},
	}
	*obj.A2 = 12

	fmt.Println(a1f.Getter()(unsafe.Pointer(&obj)))
	fmt.Println(*((a2f.Getter()(unsafe.Pointer(&obj))).(*GTInt)))
	fmt.Println(a3f.Getter()(unsafe.Pointer(&obj)))
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
	a1ptrgetter := FieldOf[A, struct{}](&(Ptr[A]().A1)).PtrGetter()
	val := A{}
	valptr := unsafe.Pointer(&val)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		fptr := a1ptrgetter(valptr)
		noopptr(fptr)
	}
}

func BenchmarkGetFieldValueByMethod(b *testing.B) {
	a1ptrgetter := FieldOf[A, struct{}](&(Ptr[A]().A1)).Getter()
	val := A{}
	valptr := unsafe.Pointer(&val)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		fptr := a1ptrgetter(valptr)
		noopval(fptr)
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
