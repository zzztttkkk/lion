package lion

import (
	"fmt"
	"testing"
)

type _CommonA struct {
	A string `json:"a"`
	B string `json:"b"`
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
	ptr := Ptr[UserA]()

	field1 := FieldOf[UserA](&ptr._CommonA.A)
	fmt.Println(field1.StructField(), field1.Offset())
	field2 := FieldOf[UserA](&(ptr._X._CommonA.A))
	fmt.Println(field2.StructField(), field2.Offset())
	fmt.Println(field1.StructField() == field2.StructField())

	for _, v := range TypeInfoOf[UserA]().fields {
		fmt.Println(v.String())
	}
}
