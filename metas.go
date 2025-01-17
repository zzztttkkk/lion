package lion

import (
	"reflect"

	"github.com/zzztttkkk/lion/internal"
)

func (filed *Field) getMetainfo(metatype reflect.Type) any {
	if filed.ref == nil {
		return filed.metas[metatype]
	}
	return filed.ref.getMetainfo(metatype)
}

func (filed *Field) updateMetainfo(metatype reflect.Type, meta any) {
	if filed.ref == nil {
		if filed.metas == nil {
			filed.metas = map[reflect.Type]any{}
		}
		filed.metas[metatype] = meta
		return
	}
	filed.ref.updateMetainfo(metatype, meta)
}

// MetaOf
// returns the meta information of the field.
func MetaOf[T any, M any](fptr any) *M {
	val := FieldOf[T](fptr).getMetainfo(Typeof[M]())
	if val == nil {
		return nil
	}
	return val.(*M)
}

// UpdateMetaScope
// update the meta information of T.
// you can only this function in `init`.
func UpdateMetaScope[T any, M any](fnc func(mptr *T, update func(ptr any, meta *M))) {
	internal.EnusreInInitFunc(fnc)
	fnc(Ptr[T](), func(fptr any, meta *M) {
		FieldOf[T](fptr).updateMetainfo(Typeof[M](), meta)
	})
}

// ReadMetaScope
// read the meta information of T.
func ReadMetaScope[T any, M any](fnc func(mptr *T, read func(fptr any) *M)) {
	fnc(Ptr[T](), func(fptr any) *M { return MetaOf[T, M](fptr) })
}
