package lion_test

import (
	"fmt"
	"testing"
	"unsafe"

	"github.com/zzztttkkk/lion"
)

type MLT struct {
	Obj [1280]int64
}

var (
	mltfobj *lion.Field[struct{}]
)

func init() {
	mptr := lion.Ptr[MLT]()
	mltfobj = lion.FieldOf[MLT, struct{}](&mptr.Obj)
}

func objptr(idx int64) *[1280]int64 {
	var obj = new(MLT)
	obj.Obj[8] = idx
	var uptr = uint64(uintptr(unsafe.Pointer(obj)))
	return (mltfobj.PtrGetter()(unsafe.Pointer(uintptr(uptr))).(*[1280]int64))
}

func TestMemLeakForPtrGet(t *testing.T) {
	fptr := objptr(-1)
	fmt.Println(fptr[8])
	for i := 0; i < 1000000; i++ {
		objptr(int64(i))
		if fptr[8] != -1 {
			fmt.Println("!!!!!!!!!!!!!!", fptr[8])
			break
		}
	}
}
