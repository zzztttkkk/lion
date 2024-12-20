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

func init() {
	lion.RegisterOf[VldMetainfo]().TagNames("vld").Unexposed()
}

func VldField[T any](ptr any) *lion.Field[VldMetainfo] {
	return lion.FieldOf[T, VldMetainfo](ptr)
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
	VldField[User](&ptr.Age).UpdateMetainfo(&VldMetainfo{Regexp: "age"})
}

func TestTypeinfoOf(t *testing.T) {
	obj := &User{}
	objptr := unsafe.Pointer(obj)

	mptr := lion.Ptr[User]()

	VldField[User](&mptr.Age).ChangeInstance(objptr, 12)
	VldField[User](&mptr.Name).ChangeInstance(objptr, "ztk")
	VldField[User](&mptr.CreatedAt).ChangeInstance(objptr, int64(32))
	VldField[User](&mptr._DeletedAt).ChangeInstance(objptr, int64(45))

	fmt.Println(VldField[User](&mptr._DeletedAt).ValueOfInstance(objptr).(int64) == 45)

	deleted_at_ptr := VldField[User](&mptr._DeletedAt).PtrGetter()(objptr).(*int64)
	fmt.Println(deleted_at_ptr)
	fmt.Println(obj)
	*deleted_at_ptr = 455
	fmt.Println(obj)
}
