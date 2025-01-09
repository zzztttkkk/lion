// Code generated by "github.com/zzztttkkk/lion/enums", DO NOT EDIT
// Code generated @ 1736395692

package enums_test

import "fmt"

import "database/sql/driver"


import "encoding/json"


func (ev X) String() string {
	switch(ev){
		
		case X_1 : {
			return "_XX1"
		}
		case X_2 : {
			return "_XX2"
		}
		case X_3 : {
			return "_XX3"
		}
		case X_ : {
			return "_XXX_"
		}
		case X_4 : {
			return "_XX4"
		}
		case X_6 : {
			return "Six"
		}
		case X_7 : {
			return "_XX7"
		}
		case X_8 : {
			return "_XX8"
		}
		case X_9 : {
			return "_XX9"
		}
		case X_10 : {
			return "_XX10"
		}
		default: {
			panic(fmt.Errorf("enums_test.X: unknown enum value, %d", ev))
		} 
	}
}


func init(){
	
	AllXValues = append(AllXValues, X_1)
	
	AllXValues = append(AllXValues, X_2)
	
	AllXValues = append(AllXValues, X_3)
	
	AllXValues = append(AllXValues, X_)
	
	AllXValues = append(AllXValues, X_4)
	
	AllXValues = append(AllXValues, X_6)
	
	AllXValues = append(AllXValues, X_7)
	
	AllXValues = append(AllXValues, X_8)
	
	AllXValues = append(AllXValues, X_9)
	
	AllXValues = append(AllXValues, X_10)
	
}


var (
	_EnumXNameMap = map[string]X{}
)
func _getEnumXByName(name string) (X, error) {
	v, ok := _EnumXNameMap[name]
	if ok {
		return v, nil
	}
	return (X)(0), fmt.Errorf("enums_test.X: invalid enum name, %s", name)
}
func init(){
	_EnumXNameMap["_XX1"] = X_1
	_EnumXNameMap["_XX2"] = X_2
	_EnumXNameMap["_XX3"] = X_3
	_EnumXNameMap["_XXX_"] = X_
	_EnumXNameMap["_XX4"] = X_4
	_EnumXNameMap["Six"] = X_6
	_EnumXNameMap["_XX7"] = X_7
	_EnumXNameMap["_XX8"] = X_8
	_EnumXNameMap["_XX9"] = X_9
	_EnumXNameMap["_XX10"] = X_10

}


// JSON impl
func (ev X) MarshalJSON() ([]byte, error) {
	return json.Marshal(ev.String())
}
func (ev *X) UnmarshalJSON(bs []byte) error {
	var name string
	if err := json.Unmarshal(bs, &name); err != nil{
		return nil
	}
	emv, err := _getEnumXByName(name)
	if err != nil {
		return err
	}
	*ev = emv
	return nil
}


// Sql impl
func (ev X) Value() (driver.Value, error) {
	return ev.String(), nil
}
func (ev *X) Scan(val any) error {
	if val == nil {
		return nil
	}
	var name string
	switch tv := val.(type) {
		case string: {
			name = tv
		}
		case []byte: {
			name = string(tv)
		}
		default: {
			return fmt.Errorf("enums_test.X: invalid value, %v", val)
		}
	}
	emv, err := _getEnumXByName(name)
	if err != nil{
		return err
	}
	*ev = emv
	return nil
}

