package reflectx_test

import (
	"fmt"
	"testing"
	"unsafe"

	"github.com/zzztttkkk/reflectx"
)

type VldMetainfo struct {
	Regexp string
}

func init() {
	reflectx.RegisterOf[VldMetainfo]().TagNames("vld")
}

func VldField[T any](ptr any) *reflectx.Field[VldMetainfo] {
	return reflectx.FieldOf[T, VldMetainfo](ptr)
}

type User struct {
	Name string `vld:"name"`
	Age  int    `vld:"age"`
}

func init() {
	ptr := reflectx.Ptr[User]()
	VldField[User](&ptr.Age).Meta = &VldMetainfo{Regexp: "age"}
}

func TestTypeinfoOf(t *testing.T) {
	obj := &User{}
	objptr := unsafe.Pointer(obj)

	VldField[User](&(reflectx.Ptr[User]()).Age).SetAny(objptr, 12)
	VldField[User](&(reflectx.Ptr[User]()).Name).SetAny(objptr, "ztk")

	fmt.Println(obj)
}
