package lion

import (
	"testing"
	"unsafe"
)

type MemSafeGet struct {
	Obj [1280]int64
}

var (
	msgobjf *Field[struct{}]
)

func init() {
	mptr := Ptr[MemSafeGet]()
	msgobjf = FieldOf[MemSafeGet, struct{}](&mptr.Obj)
}

// this is ok
func getobjfptr1(idx int64) *[1280]int64 {
	var obj = new(MemSafeGet)
	obj.Obj[8] = idx
	return (msgobjf.PtrOf(unsafe.Pointer(obj)).(*[1280]int64))
}

// this is not ok
func getobjfptr2(idx int64) *[1280]int64 {
	var obj = new(MemSafeGet)
	obj.Obj[8] = idx
	ptrnum := uintptr(unsafe.Pointer(obj))
	return (msgobjf.PtrOf(unsafe.Pointer(ptrnum)).(*[1280]int64))
}

func TestMemSafeForGetObjFptr1(t *testing.T) {
	fptr := getobjfptr1(-1)
	if fptr[8] != -1 {
		t.Fail()
	}
	for i := 0; i < 1000000; i++ {
		getobjfptr1(int64(i))
		if fptr[8] != -1 {
			t.Fail()
			break
		}
	}
}

func TestMemSafeForGetObjFptr2(t *testing.T) {
	fptr := getobjfptr2(-1)
	if fptr[8] != -1 {
		t.Fail()
	}
	for i := 0; i < 1000000; i++ {
		getobjfptr2(int64(i))
		if fptr[8] != -1 {
			t.Fail()
			break
		}
	}
}

func TestMemCopy(t *testing.T) {
	var numa = int64(0)
	var numb = int64(455)
	memcopy(unsafe.Pointer(&numa), unsafe.Pointer(&numb), 8)
	if numa != numb {
		t.Error("1")
	}
	numa = 777
	var ptra = &numa
	var ptrb = &numb
	memcopy(unsafe.Pointer(&ptra), unsafe.Pointer(&ptrb), 8)
	if ptra != ptrb {
		t.Error("2")
	}
	if *ptra != numb {
		t.Error("3")
	}
}

type MemSafeSet struct {
	Obj *[1280]int64
}

var (
	mssobjf *Field[struct{}]
)

func init() {
	mptr := Ptr[MemSafeSet]()
	mssobjf = FieldOf[MemSafeSet, struct{}](&mptr.Obj)
}

func setobjfptr1(v int64) *MemSafeSet {
	ele := MemSafeSet{}
	objptr := &[1280]int64{}
	objptr[8] = v
	mssobjf.AssignTo(unsafe.Pointer(&ele), objptr)
	return &ele
}

func TestMemSafeForSet(t *testing.T) {
	obj := setobjfptr1(-1)
	if obj.Obj[8] != -1 {
		t.Fail()
	}
	for i := 0; i < 100000; i++ {
		setobjfptr1(int64(i))
		if obj.Obj[8] != -1 {
			t.Fail()
		}
	}
}
