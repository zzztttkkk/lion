package enums_test

import (
	"fmt"
	"path/filepath"
	"testing"

	"github.com/zzztttkkk/lion/enums"
)

type X int

const (
	X_1 = X(iota)
	X_2
	X_3
	X_
	X_4
	X_6
	X_7
	X_8
	X_9
	X_10
)

var (
	AllXValues []X
)

func init() {
	enums.Generate(func() *enums.Options[X] {
		return &enums.Options[X]{
			SkilModTimeCheck:   true,
			RemoveCommonPrefix: true,
			AddPrefix:          "_XX",
			NameOverwrites: map[X]string{
				X_6: "Six",
			},
			AllSlice:     true,
			AllSliceName: "AllXValues",
			Sql:          true,
			JSON:         true,
		}
	})
}

func TestEnums(t *testing.T) {
	fmt.Println(filepath.Abs("d:/xxx/aaa"))
	fmt.Println(AllXValues)
}
