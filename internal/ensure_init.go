package internal

import (
	"fmt"
	"reflect"
	"regexp"
	"runtime"
)

var (
	funcnameregexp = regexp.MustCompile(`^[a-zA-Z0-9/_.]*.init\.\d+\.func\d+$`)
)

func EnusreInInitFunc(fnc any) {
	fv := reflect.ValueOf(fnc)
	if fv.IsNil() {
		panic("nil func")
	}
	runtimefunc := runtime.FuncForPC(fv.Pointer())
	funcname := runtimefunc.Name()
	if !funcnameregexp.MatchString(funcname) {
		panic(fmt.Errorf("lion: `fnc` must be an anonymous function defined in the `init` function. `%s`", funcname))
	}
}
