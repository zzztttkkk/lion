package reflectx_test

import (
	"reflect"
	"testing"
	"unsafe"

	"github.com/zzztttkkk/reflectx"
)

type SetTest struct {
	A int
	B int64
}

func BenchmarkDirectlySet(b *testing.B) {
	var obj = SetTest{}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		obj.A = 12
	}
}

func BenchmarkUpdate(b *testing.B) {
	mptr := reflectx.Ptr[SetTest]()
	fieldOfA := reflectx.FieldOf[SetTest, reflectx.EmptyMeta](&mptr.A)

	var obj SetTest
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		reflectx.Update(&obj, fieldOfA, 12)
	}
}

func BenchmarkFieldPtrSet(b *testing.B) {
	mptr := reflectx.Ptr[SetTest]()
	fieldOfA := reflectx.FieldOf[SetTest, reflectx.EmptyMeta](&mptr.A)

	var obj SetTest
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		ptr := fieldOfA.PtrOf(unsafe.Pointer(&obj)).(*int)
		*ptr = 12
	}
}

func BenchmarkNormalReflectSet(b *testing.B) {
	obj := &SetTest{}
	objv := reflect.ValueOf(obj).Elem()
	fv := objv.FieldByName("A")

	b.ResetTimer()
	numv := reflect.ValueOf(345)
	for i := 0; i < b.N; i++ {
		fv.Set(numv)
	}
}
