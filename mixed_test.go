package reflectx

import (
	"fmt"
	"testing"
)

func init() {
	RegisterOf[EmptyMeta]().TagNames("json").Unexposed()
}

type _Common struct {
	a string
	B string `json:"b"`
}

func init() {
	ptr := Ptr[_Common]()
	FieldOf[_Common, EmptyMeta](&ptr.a).Name = "common_a"
}

type User struct {
	A string `json:"a1"`
	_Common
}

func init() {
	ptr := Ptr[User]()
	TypeInfoOf[User, EmptyMeta]().Mix(&ptr._Common, TypeInfoOf[_Common, EmptyMeta]())
}

func TestMixed(t *testing.T) {
	ptr := Ptr[User]()
	field := FieldOf[User, EmptyMeta](&ptr._Common.a)
	fmt.Println(field)
}
