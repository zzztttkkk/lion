package enums

import (
	"fmt"
	"go/parser"
	"go/token"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"reflect"
	"regexp"
	"runtime"
	"slices"
	"strings"
	"text/template"
	"time"

	"github.com/zzztttkkk/lion"
)

type Options[T lion.IntType] struct {
	SkilModTimeCheck bool

	RemoveCommonPrefix bool
	AddPrefix          string
	NameOverwrites     map[T]string

	AllSlice              bool
	AllSliceName          string
	AllSliceNotPreDefined bool
	AllSliceHidens        []T

	Sql  bool
	JSON bool
}

var (
	funcnameregexp = regexp.MustCompile(`^\.init\.\d+\.func\d+$`)
)

func should_re_gen(dir string, targetfp string) bool {
	target_stat, err := os.Stat(targetfp)
	if err != nil {
		return true
	}

	target_modtime := target_stat.ModTime()

	entries, err := os.ReadDir(dir)
	if err != nil {
		return true
	}

	for _, v := range entries {
		if v.IsDir() {
			continue
		}
		name := v.Name()
		if !strings.HasSuffix(name, ".go") {
			continue
		}
		if strings.HasSuffix(name, "_test.go") {
			continue
		}
		es, err := os.Stat(filepath.Join(dir, name))
		if err != nil || es.ModTime().After(target_modtime) {
			return true
		}
	}
	return false
}

func common_prefix(ts []string) string {
	tmp := []rune{}

outer:
	for i := 0; ; i++ {
		var c rune = 0
		for _, txt := range ts {
			chars := []rune(txt)
			if i >= len(chars) {
				break outer
			}
			if c == 0 {
				c = chars[i]
				continue
			}
			if c == chars[i] {
				continue
			}
			break outer
		}
		tmp = append(tmp, c)
	}
	return string(tmp)
}

func Generate[T lion.IntType](fnc func() *Options[T]) {
	enumtype := lion.Typeof[T]()
	enumpkgpath := enumtype.PkgPath()
	enumpkgname := path.Base(enumpkgpath)

	runtimefunc := runtime.FuncForPC(reflect.ValueOf(fnc).Pointer())
	funcname := runtimefunc.Name()
	if !strings.HasPrefix(funcname, enumpkgpath) {
		panic(fmt.Errorf("lion.enums: fnc's pkg is not same as enum's pkg"))
	}
	if !funcnameregexp.MatchString(funcname[len(enumpkgpath):]) {
		panic(fmt.Errorf("lion.enums: `fuc` must be an anonymous function defined in the `init` function. `%s`", funcname))
	}
	filename, _ := runtimefunc.FileLine(0)
	dirname := filepath.Dir(filename)

	outs, err := exec.Command("go", "env", "GOMODCACHE").Output()
	if err != nil {
		panic(fmt.Errorf("lion.enums: exec `go env GOMODCACHE` failed, %s", err))
	}
	if strings.HasPrefix(filename, strings.TrimSpace(string(outs))) {
		return
	}

	testsuf := ""
	if strings.HasSuffix(filename, "_test.go") {
		testsuf = "_test"
	}
	targetfp := fmt.Sprintf("%s/lion.enums.generate.%s%s.go", dirname, lion.Typeof[T]().Name(), testsuf)

	opts := fnc()
	if opts == nil {
		opts = &Options[T]{}
	}

	if opts.SkilModTimeCheck {
	} else {
		if !should_re_gen(dirname, targetfp) {
			return
		}
	}

	fset := token.NewFileSet()

	pkgs, err := parser.ParseDir(fset, dirname, nil, 0)
	if err != nil {
		panic(err)
	}

	pkg, ok := pkgs[enumpkgname]
	if !ok {
		return
	}

	var enums = map[string]*_EnumInfo{}
	var info = &_EnumInfo{}
	enums[enumtype.Name()] = info
	for _, file := range pkg.Files {
		_PopulateEnumInfo(enums, file)
	}
	if len(info.Consts) < 1 {
		panic(fmt.Errorf(`lion.enums: failed to scan enum values of '%s', you must define enum values like:
"""
const (
		A EnumType = iota
		B
		C
		.... 
)
"""
not:
"""
const (
		A = EnumType(iota)
		B
		C
		.... 
)
"""`, enumtype))
	}

	genGoCode(info, enumpkgname, targetfp, opts)
}

func genGoCode[T lion.IntType](enuminfo *_EnumInfo, pkgname string, targetfp string, opts *Options[T]) {
	ftpl := `// Code generated by "github.com/zzztttkkk/lion/enums", DO NOT EDIT
// Code generated @ {{.time}}

package {{.pkgname}}

import "fmt"
{{if .gensql }}
import "database/sql/driver"
{{end}}
{{if .genjson}}
import "encoding/json"
{{end}}

func (ev {{.enumtypename}}) String() string {
	switch(ev){
		{{range .items}}
		case {{.vname}} : {
			return "{{.string}}"
		}{{end}}
		default: {
			panic(fmt.Errorf("{{.pkgname}}.{{.enumtypename}}: unknown enum value, %d", ev))
		} 
	}
}
{{if .defineallslice }}
var(
	{{.allslicename}} []{{.enumtypename}}
)
{{end}}
{{if .appendtoallslice}}
func init(){
	{{range .appenitems}}
	{{$.allslicename}} = append({{$.allslicename}}, {{.}})
	{{end}}
}
{{end}}
{{if or .gensql .genjson   }}
var (
	_Enum{{.enumtypename}}NameMap = map[string]{{.enumtypename}}{}
)
func _getEnum{{$.enumtypename}}ByName(name string) ({{$.enumtypename}}, error) {
	v, ok := _Enum{{.enumtypename}}NameMap[name]
	if ok {
		return v, nil
	}
	return ({{$.enumtypename}})(0), fmt.Errorf("{{.pkgname}}.{{.enumtypename}}: invalid enum name, %s", name)
}
func init(){
{{range .items}}	_Enum{{$.enumtypename}}NameMap["{{.string}}"] = {{.vname}}
{{end}}
}
{{end}}
{{if .genjson}}
// JSON impl
func (ev {{.enumtypename}}) MarshalJSON() ([]byte, error) {
	return json.Marshal(ev.String())
}
func (ev *{{.enumtypename}}) UnmarshalJSON(bs []byte) error {
	var name string
	if err := json.Unmarshal(bs, &name); err != nil{
		return nil
	}
	emv, err := _getEnum{{$.enumtypename}}ByName(name)
	if err != nil {
		return err
	}
	*ev = emv
	return nil
}
{{end}}
{{if .gensql}}
// Sql impl
func (ev {{.enumtypename}}) Value() (driver.Value, error) {
	return ev.String(), nil
}
func (ev *{{.enumtypename}}) Scan(val any) error {
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
			return fmt.Errorf("{{.pkgname}}.{{.enumtypename}}: invalid value, %v", val)
		}
	}
	emv, err := _getEnum{{$.enumtypename}}ByName(name)
	if err != nil{
		return err
	}
	*ev = emv
	return nil
}
{{end}}
`
	items := []map[string]string{}
	prefix := ""
	if opts.RemoveCommonPrefix {
		names := []string{}
		for _, cv := range enuminfo.Consts {
			names = append(names, cv.Name)
		}
		prefix = common_prefix(names)
	}

	for _, cv := range enuminfo.Consts {
		item := map[string]string{
			"vname":  cv.Name,
			"string": cv.Name,
		}
		if prefix != "" && strings.HasPrefix(cv.Name, prefix) {
			new_name := cv.Name[len(prefix):]
			if new_name == "" {
				new_name = cv.Name
			}
			item["string"] = new_name
		}
		if opts.AddPrefix != "" {
			item["string"] = fmt.Sprintf("%s%s", opts.AddPrefix, item["string"])
		}
		if opts.NameOverwrites != nil {
			ow := opts.NameOverwrites[T(reflect.ValueOf(cv.Value).Int())]
			if ow != "" {
				item["string"] = ow
			}
		}

		items = append(items, item)
	}
	val := map[string]any{
		"time":         time.Now().Unix(),
		"pkgname":      pkgname,
		"enumtypename": lion.Typeof[T]().Name(),
		"items":        items,
		"gensql":       opts.Sql,
		"genjson":      opts.JSON,
	}

	if opts.AllSlice {
		val["appendtoallslice"] = true

		name := opts.AllSliceName
		if name == "" {
			name = fmt.Sprintf("All%ss", lion.Typeof[T]().Name())
		}
		val["allslicename"] = name
		if opts.AllSliceNotPreDefined {
			val["defineallslice"] = true
		}

		appenditems := []string{}
		for _, cv := range enuminfo.Consts {
			if slices.IndexFunc(opts.AllSliceHidens, func(v T) bool {
				return reflect.ValueOf(cv.Value).Int() == reflect.ValueOf(v).Int()
			}) > -1 {
				continue
			}
			appenditems = append(appenditems, cv.Name)
		}

		val["appenitems"] = appenditems
	}

	sb := strings.Builder{}
	tpl := template.Must(template.New("").Parse(ftpl))
	err := tpl.Execute(&sb, val)
	if err != nil {
		panic(err)
	}

	f, err := os.OpenFile(targetfp, os.O_WRONLY|os.O_CREATE, 0o0755)
	if err != nil {
		panic(err)
	}
	defer f.Close()
	f.Truncate(0)
	_, err = f.WriteString(sb.String())
	if err != nil {
		panic(err)
	}
	fmt.Printf("lion.enums: the file `%s` has been generated, you need to recompile.\r\n", targetfp)
}
