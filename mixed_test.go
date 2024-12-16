package reflectx

import (
	"fmt"
	"testing"
)

func init() {
	RegisterOf[EmptyMeta]().TagNames("json")
}

type Common struct {
	A string `json:"a"`
	B string `json:"b"`
}

func init() {
	ptr := Ptr[Common]()
	FieldOf[Common, EmptyMeta](&ptr.A).Name = "common_a"
}

type User struct {
	A string `json:"a1"`
	Common
}

func init() {
	ptr := Ptr[User]()
	TypeInfoOf[User, EmptyMeta]().Mixed(&ptr.Common, TypeInfoOf[Common, EmptyMeta]())
}

func TestMixed(t *testing.T) {
	ptr := Ptr[User]()
	field := FieldOf[User, EmptyMeta](&ptr.Common.A)
	fmt.Println(field.Name)
}
