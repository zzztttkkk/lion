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
	reflectx.RegisterOf[VldMetainfo]().TagNames("vld").Unexposed()
}

func VldField[T any](ptr any) *reflectx.Field[VldMetainfo] {
	return reflectx.FieldOf[T, VldMetainfo](ptr)
}

type _Common struct {
	CreatedAt  int64  `vld:"created_at"`
	_DeletedAt int64  `vld:"deleted_at"`
	Name       string `vld:"cname"`
}

type User struct {
	_Common
	Name string `vld:"name"`
	Age  int    `vld:"age"`
}

func init() {
	ptr := reflectx.Ptr[User]()
	VldField[User](&ptr.Age).UpdateMetainfo(&VldMetainfo{Regexp: "age"})
}

func TestTypeinfoOf(t *testing.T) {
	obj := &User{}
	objptr := unsafe.Pointer(obj)

	mptr := reflectx.Ptr[User]()

	reflectx.Update(obj, VldField[User](&mptr.Age), 12)
	reflectx.Update(obj, VldField[User](&mptr.Name), "ztk")
	reflectx.Update(obj, VldField[User](&mptr.CreatedAt), int64(23))
	reflectx.Update(obj, VldField[User](&mptr._DeletedAt), int64(485))

	deleted_at_ptr := VldField[User](&mptr._DeletedAt).PtrGetter()(objptr).(*int64)
	fmt.Println(deleted_at_ptr)
	fmt.Println(obj)
	*deleted_at_ptr = 455
	fmt.Println(obj)
}
