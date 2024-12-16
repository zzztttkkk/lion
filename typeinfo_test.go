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

type _Common struct {
	CreatedAt  int64 `vld:"created_at"`
	_DeletedAt int64 `vld:"deleted_at"`
}

type User struct {
	_Common
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
	VldField[User](&(reflectx.Ptr[User]()).CreatedAt).SetAny(objptr, int64(34))
	VldField[User](&(reflectx.Ptr[User]())._DeletedAt).SetAny(objptr, int64(134))

	deleted_at_ptr := VldField[User](&(reflectx.Ptr[User]())._DeletedAt).PtrGetter()(objptr).(*int64)
	fmt.Println(deleted_at_ptr)
	fmt.Println(obj)
	*deleted_at_ptr = 455
	fmt.Println(obj)
}
