package enums

import (
	"fmt"
	"testing"
)

//go:generate stringer -type XInt -trimprefix "X"
type XInt int

const (
	XOne XInt = iota
	XTwo
	XThree
)

var (
	AllXInts    = All(XOne, 100)
	AllXIntsMap = Map(XOne, 100)
)

func TestXxx(t *testing.T) {
	fmt.Println(AllXInts, AllXIntsMap)

	fmt.Println(NewEnumNamesIndex(XOne, 100, true).Find("one"))
}
