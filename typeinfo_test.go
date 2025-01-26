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

type _Common struct {
	CreatedAt  int64  `vld:"created_at"`
	_DeletedAt int64  `vld:"deleted_at"`
	Name       string `vld:"cname"`
	XX         bool   `vld:"-"`
}

type User struct {
	_Common
	Name string `vld:"name"`
	Age  int    `vld:"age"`
}

func init() {
	lion.UpdateMetaScope(func(mptr *User, update func(ptr any, meta *VldMetainfo)) {
		update(
			&mptr.Age,
			&VldMetainfo{Regexp: "age"},
		)
	})
}

func TestTypeinfoOf(t *testing.T) {
	obj := &User{}
	objptr := unsafe.Pointer(obj)

	mptr := lion.Ptr[User]()

	lion.FieldOf[User](&mptr.Age).AssignTo(objptr, 12)
	lion.FieldOf[User](&mptr.Name).AssignTo(objptr, "ztk")
	lion.FieldOf[User](&mptr.CreatedAt).AssignTo(objptr, int64(32))
	lion.FieldOf[User](&mptr._DeletedAt).AssignTo(objptr, int64(45))

	deleted_at_ptr := lion.FieldOf[User](&mptr._DeletedAt).PtrOf(objptr).(*int64)
	fmt.Println(deleted_at_ptr)
	fmt.Println(obj)
	*deleted_at_ptr = 455
	fmt.Println(obj)

	for v := range lion.TypeInfoOf[User]().Fields(nil) {
		fmt.Println(v)
	}
}
