package lion_test

import (
	"fmt"
	"testing"

	"github.com/zzztttkkk/lion"
)

func init() {
	lion.RegisterOf[lion.EmptyMeta]().TagNames("json").Unexposed()
}

type _CommonA struct {
	A string `json:"a"`
	B string `json:"b"`
}

func init() {
	ptr := lion.Ptr[_CommonA]()
	lion.FieldOf[_CommonA, lion.EmptyMeta](&ptr.A).SetName("common_a")
}

type _X struct {
	_CommonA
}

type UserA struct {
	A string `json:"a1"`
	_CommonA
	_X
}

func TestMixed(t *testing.T) {
	ptr := lion.Ptr[UserA]()

	field1 := lion.FieldOf[UserA, lion.EmptyMeta](&ptr._CommonA.A)
	fmt.Println(field1.StructField(), field1.Offset())
	field2 := lion.FieldOf[UserA, lion.EmptyMeta](&(ptr._X._CommonA.A))
	fmt.Println(field2.StructField(), field2.Offset())
	fmt.Println(field1.StructField() == field2.StructField())

	for _, v := range lion.TypeInfoOf[UserA, lion.EmptyMeta]().Fields {
		fmt.Println(v.String())
	}
}
