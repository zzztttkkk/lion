package enums

import (
	"fmt"
	"testing"
)

func TestCommonPrefix(t *testing.T) {
	fmt.Println(common_prefix([]string{"x1", "x2", "xx3"}))
}
