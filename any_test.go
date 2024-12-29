package lion

import (
	"fmt"
	"math/rand"
	"testing"
)

func TestAny(t *testing.T) {
	var x int64 = rand.Int63()
	var ptr = anytotype[int64](x)
	fmt.Println(*ptr == x)
}
