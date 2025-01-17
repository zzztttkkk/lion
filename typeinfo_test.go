package lion_test

import (
	"fmt"
	"testing"
	"unsafe"

	"github.com/zzztttkkk/lion"
)

type VldMetainfo struct {
	Regexp string
}

func VldField[T any](ptr any) *lion.Field {
	return lion.FieldOf[T](ptr)
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
	ptr := lion.Ptr[User]()
	lion.UpdateMetainfo[User](&ptr.Age, &VldMetainfo{Regexp: "age"})
}

func TestTypeinfoOf(t *testing.T) {
	obj := &User{}
	objptr := unsafe.Pointer(obj)

	mptr := lion.Ptr[User]()

	VldField[User](&mptr.Age).AssignTo(objptr, 12)
	VldField[User](&mptr.Name).AssignTo(objptr, "ztk")
	VldField[User](&mptr.CreatedAt).AssignTo(objptr, int64(32))
	VldField[User](&mptr._DeletedAt).AssignTo(objptr, int64(45))

	deleted_at_ptr := VldField[User](&mptr._DeletedAt).PtrOf(objptr).(*int64)
	fmt.Println(deleted_at_ptr)
	fmt.Println(obj)
	*deleted_at_ptr = 455
	fmt.Println(obj)

	for v := range lion.TypeInfoOf[User]().Fields(nil) {
		fmt.Println(v)
	}
}
