package enums

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/zzztttkkk/lion"
)

type IEnum interface {
	lion.IntType
	fmt.Stringer
}

var (
	enumStringRegex = regexp.MustCompile("^[a-zA-Z_][a-zA-Z0-9_]*$")
)

func All[T IEnum](min T, max T) []T {
	items := []T{}
	for i := min; i <= max; i++ {
		if enumStringRegex.MatchString(i.String()) {
			items = append(items, i)
		}
	}
	return items
}

func Map[T IEnum](min T, max T) map[string]T {
	items := map[string]T{}
	for i := min; i <= max; i++ {
		if enumStringRegex.MatchString(i.String()) {
			items[i.String()] = i
		}
	}
	return items
}

type EnumNamesIndex[T IEnum] struct {
	mapv       map[string]T
	ignorecase bool
}

func (eni *EnumNamesIndex[T]) Find(name string) (T, bool) {
	if !eni.ignorecase {
		v, ok := eni.mapv[name]
		return v, ok
	}
	v, ok := eni.mapv[strings.ToUpper(name)]
	return v, ok
}

func NewEnumNamesIndex[T IEnum](min T, max T, ignorecase bool) *EnumNamesIndex[T] {
	obj := &EnumNamesIndex[T]{
		mapv:       map[string]T{},
		ignorecase: ignorecase,
	}

	if ignorecase {
		for k, v := range Map(min, max) {
			obj.mapv[strings.ToUpper(k)] = v
		}
	} else {
		obj.mapv = Map(min, max)
	}
	return obj
}
