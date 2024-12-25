package lion

import "reflect"

type SingedInt interface {
	~int | ~int8 | ~int16 | ~int32 | ~int64
}

type UnsignedInt interface {
	~uint | ~uint8 | ~uint16 | ~uint32 | ~uint64
}

type IntType interface {
	SingedInt | UnsignedInt
}

func IsUnsignedInt[T any]() bool {
	switch Typeof[T]().Kind() {
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		{
			return true
		}
	}
	return false
}
