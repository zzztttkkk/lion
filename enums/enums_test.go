package enums_test

import (
	"fmt"
	"testing"

	"github.com/zzztttkkk/lion/enums"
)

type X int

const (
	X1 X = (iota)
	X2
	X3
	X4
	X6
)

var (
	AllXValues []X
)

func init() {
	enums.Generate(func() *enums.EnumOptions[X] {
		return &enums.EnumOptions[X]{
			NamePrefix: "X",
			NameOverwrites: map[X]string{
				X6: "Six",
			},
			GenAllSlice: true,
			// AllSliceNotPreDefined: true,
			AllSliceName: "AllXValues",
		}
	})
}

func TestEnums(t *testing.T) {
	fmt.Println(AllXValues)
}
