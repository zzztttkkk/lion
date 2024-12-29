package lion

import "unsafe"

// https://github.com/goccy/go-json/tree/master?tab=readme-ov-file#elimination-of-reflection

type emptyInterface struct {
	typ unsafe.Pointer
	ptr unsafe.Pointer
}

func anytotype[T any](v any) *T {
	iface := (*emptyInterface)(unsafe.Pointer(&v))
	return (*T)(iface.ptr)
}
