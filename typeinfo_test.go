package reflectx_test

import (
	"fmt"
	"math/rand"
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
	VldField[User](&ptr.Age).Meta = &VldMetainfo{Regexp: "age"}
}

func TestTypeinfoOf(t *testing.T) {
	obj := &User{}
	objptr := unsafe.Pointer(obj)

	VldField[User](&(reflectx.Ptr[User]()).Age).Set(objptr, 12)
	VldField[User](&(reflectx.Ptr[User]()).Name).Set(objptr, "ztk")
	VldField[User](&(reflectx.Ptr[User]()).CreatedAt).Set(objptr, int64(34))
	VldField[User](&(reflectx.Ptr[User]())._DeletedAt).Set(objptr, int64(134))

	deleted_at_ptr := VldField[User](&(reflectx.Ptr[User]())._DeletedAt).PtrGetter()(objptr).(*int64)
	fmt.Println(deleted_at_ptr)
	fmt.Println(obj)
	*deleted_at_ptr = 455
	fmt.Println(obj)
}

type Pair struct {
	Key int64
	Val int64
}

var (
	lstmap  = []Pair{}
	hashmap = map[int64]int64{}
)

func init() {
	for i := 0; i < 15; i++ {
		key := rand.Int63()
		val := rand.Int63()
		lstmap = append(lstmap, Pair{Key: key, Val: val})
		hashmap[key] = val
	}

}

func noop(v any) {}

func BenchmarkReadLstMap(b *testing.B) {
	for i := 0; i < b.N; i++ {
		key := rand.Int63()
		for idx := range lstmap {
			ptr := &lstmap[idx]
			if ptr.Key == key {
				noop(ptr.Val)
				break
			}
		}
	}
}

func BenchmarkReadHashMap(b *testing.B) {
	for i := 0; i < b.N; i++ {
		key := rand.Int63()
		v, ok := hashmap[key]
		if ok {
			noop(v)
		}
	}
}
