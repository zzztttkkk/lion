package lion

import "reflect"

type _KindsNamespace struct{}

var Kinds _KindsNamespace

func (_ns _KindsNamespace) IsInt(v reflect.Kind) bool {
	switch v {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		{
			return true
		}
	default:
		{
			return false
		}
	}
}

func (_ns _KindsNamespace) IsUint(v reflect.Kind) bool {
	switch v {
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		{
			return true
		}
	default:
		{
			return false
		}
	}
}

func (_ns _KindsNamespace) IsFloat(v reflect.Kind) bool {
	switch v {
	case reflect.Float32, reflect.Float64:
		{
			return true
		}
	default:
		{
			return false
		}
	}
}

func (_ns _KindsNamespace) IsValue(v reflect.Kind) bool {
	if _ns.IsInt(v) || _ns.IsUint(v) || _ns.IsFloat(v) {
		return true
	}
	switch v {
	case reflect.Bool, reflect.String, reflect.Struct, reflect.Complex64, reflect.Complex128:
		{
			return true
		}
	default:
		{
			return false
		}
	}
}
