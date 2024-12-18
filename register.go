package lion

import "reflect"

type _Register[M any] struct {
	tagnames  []string
	unexposed bool
}

var (
	registers = map[reflect.Type]any{}
)

func RegisterOf[M any]() *_Register[M] {
	gotype := Typeof[M]()

	v, ok := registers[gotype]
	if ok {
		return v.(*_Register[M])
	}
	obj := &_Register[M]{}
	registers[gotype] = obj
	return obj
}

func (reg *_Register[M]) TagNames(names ...string) *_Register[M] {
	reg.tagnames = append(reg.tagnames, names...)
	return reg
}

func (reg *_Register[M]) Unexposed() *_Register[M] {
	reg.unexposed = true
	return reg
}
