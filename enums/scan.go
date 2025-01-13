// https://stackoverflow.com/a/75610059/6683474
// https://go.dev/play/p/yEz6G1Kqoe3

package enums

import (
	"go/ast"
	"go/token"
)

type _EnumInfo struct {
	Name     string
	TypeName string // base type
	Consts   []_ConstValue
}

type _ConstValue struct {
	Name  string
	Value any // int or string
}

func _PopulateEnumInfo(enumTypesMap map[string]*_EnumInfo, file *ast.File) {
	// phase 1: iterate scope objects to get the values
	var nameValues = make(map[string]any)

	for _, object := range file.Scope.Objects {
		if object.Kind == ast.Con {
			nameValues[object.Name] = object.Data
		}
	}

	// phase 2: iterate decls to get the type and names in order
	for _, decl := range file.Decls {
		genDecl, ok := decl.(*ast.GenDecl)
		if !ok {
			continue
		}
		if genDecl.Tok != token.CONST {
			continue
		}
		var enumInfo *_EnumInfo
		for _, spec := range genDecl.Specs {
			valSpec := spec.(*ast.ValueSpec)
			if valSpec.Type == nil {
				if len(valSpec.Values) > 0 {
					fv, ok := valSpec.Values[0].(*ast.CallExpr)
					if ok {
						fun, ok := fv.Fun.(*ast.Ident)
						if ok {
							enumInfo = enumTypesMap[fun.Name]
						}
					}
				}
			} else {
				if typeIdent, ok := valSpec.Type.(*ast.Ident); ok {
					enumInfo = enumTypesMap[typeIdent.String()]
				}
			}
			if enumInfo != nil {
				for _, nameIdent := range valSpec.Names {
					name := nameIdent.String()
					if name == "_" {
						continue
					}
					value := nameValues[name]
					enumInfo.Consts = append(enumInfo.Consts, _ConstValue{
						Name:  name,
						Value: value,
					})
				}
			}
		}
	}
}
