package lion

import (
	"fmt"
	"reflect"
	"testing"
	"unsafe"
)

type GTInt int64

type IVVV interface {
	vvv()
}

type IfceByPtr struct {
	Num int64
}

func (vvv *IfceByPtr) vvv() {
	fmt.Printf("IfceByPtr: %p, %d\r\n", unsafe.Pointer(vvv), vvv.Num)
}

type IfceByVal struct {
	Num int64
}

func (vvv IfceByVal) vvv() {
	fmt.Printf("IfceByVal: %d\r\n", vvv.Num)
}

type A struct {
	A1 string `json:"a1"`
	A2 *GTInt
	A3 []GTInt
	A4 IVVV
	A5 func()
	A6 chan int
	A7 unsafe.Pointer
}

func TestGetFieldPtr(t *testing.T) {
	mptr := Ptr[A]()

	a1f := FieldOf[A](&mptr.A1)
	a2f := FieldOf[A](&mptr.A2)
	a3f := FieldOf[A](&mptr.A3)
	a4f := FieldOf[A](&mptr.A4)
	a5f := FieldOf[A](&mptr.A5)
	a6f := FieldOf[A](&mptr.A6)
	a7f := FieldOf[A](&mptr.A7)

	obj := A{
		A1: "a1",
		A2: new(GTInt),
		A3: []GTInt{7, 888, 9},
		A4: &IfceByPtr{Num: 567},
		A5: func() {
			fmt.Println(">>>>>>>>>>> A5")
		},
		A6: make(chan int),
		A7: unsafe.Pointer(mptr),
	}
	*obj.A2 = 12
	objuptr := unsafe.Pointer(&obj)

	fmt.Println(*(a1f.PtrOf(objuptr).(*string)))
	fmt.Println(**(a2f.PtrOf(objuptr).(**GTInt)))
	fmt.Println(a3f.PtrOf(objuptr).(*[]GTInt))

	(*(a4f.PtrOf(objuptr).(*IVVV))).vvv()

	fmt.Println(a1f.ValueOf(objuptr), obj.A1)
	fmt.Println(a2f.ValueOf(objuptr), obj.A2)
	fmt.Println(a3f.ValueOf(objuptr), obj.A3)

	fmt.Println(">>>", a4f.ValueOf(objuptr).(IVVV))

	obj.A4 = IfceByVal{Num: 789}
	(*(a4f.PtrOf(objuptr).(*IVVV))).vvv()
	fmt.Println(">>>", a4f.ValueOf(objuptr).(IVVV))

	a5fnc := a5f.ValueOf(objuptr).(func())
	a5fnc()
	fmt.Println(a6f.ValueOf(objuptr).(chan int))

	fmt.Println(a7f.ValueOf(objuptr), unsafe.Pointer(mptr))
}

func TestSetInterface(t *testing.T) {
	mptr := Ptr[A]()

	a4f := FieldOf[A](&mptr.A4)

	obj := A{
		A1: "a1",
		A2: new(GTInt),
		A3: []GTInt{7, 888, 9},
		A4: nil,
	}
	*obj.A2 = 12
	objuptr := unsafe.Pointer(&obj)

	a4f.AssignTo(objuptr, IfceByVal{566})
	a4f.ValueOf(objuptr).(IVVV).vvv()
	a4f.AssignTo(objuptr, &IfceByPtr{Num: 765})
	a4f.ValueOf(objuptr).(IVVV).vvv()
	a4f.AssignTo(objuptr, nil)
	fmt.Println(a4f.ValueOf(objuptr))
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

func BenchmarkGetFieldValByMethod(b *testing.B) {
	a1f := FieldOf[A](&(Ptr[A]().A1))
	val := A{}
	valptr := unsafe.Pointer(&val)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		fptr := a1f.ValueOf(valptr)
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
