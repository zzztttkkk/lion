package lion

import (
	"reflect"

	"github.com/zzztttkkk/lion/internal"
)

// Metainfo
// the meta information of the metatype on this field.
func (filed *Field) Metainfo(metatype reflect.Type) any {
	if filed.ref == nil {
		return filed.metas[metatype]
	}
	return filed.ref.Metainfo(metatype)
}

// UpdateMetainfo
// update the meta information of the metatype on this field.
func (filed *Field) UpdateMetainfo(metatype reflect.Type, meta any) {
	if filed.ref == nil {
		if filed.metas == nil {
			filed.metas = map[reflect.Type]any{}
		}
		filed.metas[metatype] = meta
		return
	}
	filed.ref.UpdateMetainfo(metatype, meta)
}

// MetaOf
// returns the meta information of the field.
func MetaOf[T any, M any](fptr any) *M {
	val := FieldOf[T](fptr).Metainfo(Typeof[M]())
	if val == nil {
		return nil
	}
	return val.(*M)
}

// UpdateMetaFor
// update the meta information of the field.
func UpdateMetaFor[T any, M any](fptr any, meta *M) {
	FieldOf[T](fptr).UpdateMetainfo(Typeof[M](), meta)
}

// UpdateMetaScope
// update the meta information of T.
// you can only this function in `init`.
func UpdateMetaScope[T any, M any](fnc func(mptr *T, update func(ptr any, meta *M))) {
	internal.EnusreInInitFunc(fnc)
	fnc(Ptr[T](), func(fptr any, meta *M) {
		FieldOf[T](fptr).UpdateMetainfo(Typeof[M](), meta)
	})
}

// ReadMetaScope
// read the meta information of T.
func ReadMetaScope[T any, M any](fnc func(mptr *T, read func(fptr any) *M)) {
	fnc(Ptr[T](), func(fptr any) *M { return MetaOf[T, M](fptr) })
}
