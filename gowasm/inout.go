package gowasm

import (
	"bytes"
	"fmt"
	"io"
	"strconv"
	"strings"
	"text/template"

	"github.com/gowebapi/webidl-bind/types"
)

const inoutToTmplInput = `
{{define "start"}}
{{if .ReleaseHdl}}	var _releaseList releasableApiResourceList {{end}}
{{end}}
{{define "end"}}
{{end}}

{{define "param-start"}}
	{{if .Optional}}
		{{if and .AnyType (not .UseIn)}}
			if {{.In}}.Type() != js.TypeUndefined {
		{{else}}
			if {{.In}} != nil {
		{{end}}
	{{end}}
{{end}}
{{define "param-end"}}
	{{.Assign}}
	{{if .Optional}}
		}
	{{end}}
{{end}}

{{define "type-primitive"}}
	{{if .Info.Pointer}}
		var {{.Out}} interface{}
		if {{.In}} != nil {
			{{.Out}} = *( {{.In}} )
		} else {
			{{.Out}} = nil
		}
	{{else}}
		{{.Out}} := {{.In}}
	{{end}}
{{end}}
{{define "type-dictionary"}}	{{.Out}} := {{.In}}.JSValue() {{end}}
{{define "type-interface"}}		{{.Out}} := {{.In}}.JSValue() {{end}}

{{define "type-callback"}}
	var __callback{{.Idx}} js.Value
	if {{.In}} != nil {
		__callback{{.Idx}} = ( * {{.In}} ).Value
	} else {
		__callback{{.Idx}} = js.Null()
	}
	{{.Out}} := __callback{{.Idx}}
{{end}}
{{define "type-enum"}}      {{.Out}} := {{.In}}.JSValue() {{end}}
{{define "type-union"}}	{{.Out}} := {{.In}}.JSValue() {{end}}
{{define "type-any"}}    {{.Out}} := {{.In}} {{end}}
{{define "type-typedarray"}} {{.Out}} := jsarray.{{.GoFunc}}ToJS( {{.In}} ) {{end}}
{{define "type-parametrized"}}	{{.Out}} := {{.In}}.JSValue() {{end}}
{{define "type-rawjs"}}    {{.Out}} := {{.In}} {{end}}

{{define "type-sequence"}} 
	{{.Out}} := js.Global().Get("Array").New(len( {{if .Info.Pointer}}*{{end}} {{.In}} ))
	for __idx{{.Idx}} , __seq_in{{.Idx}} := range {{if .Info.Pointer}}*{{end}} {{.In}} {
		{{.Inner}}
		{{.Out}} .SetIndex( __idx{{.Idx}} , __seq_out{{.Idx}} )
	}
{{end}}

{{define "type-variadic"}}
	for _, __in := range {{.In}} {
		{{.Inner}}
		_args[_end] = __out
		_end++
	} 
{{end}}

`

const inoutFromTmplInput = `
{{define "start"}}
	var (
	{{range .ParamList}}
		{{.Out}} {{.Var}} // javascript: {{.Info.Idl}} {{.Name}}
	{{end}}
	)
{{end}}
{{define "end"}}{{end}}

{{define "param-start"}}
	{{if .Optional}}
		if len(args) > {{.Idx}} {
	{{end}}
	{{if .Nullable}}
		if {{.In}}.Type() != js.TypeNull && {{.In}}.Type() != js.TypeUndefined {
	{{end}}
{{end}}
{{define "param-end"}}
	{{if .Optional}}
		}
	{{end}}
	{{if .Nullable}}
		}
	{{end}}
{{end}}

{{define "type-primitive"}}	
	{{if .Info.Pointer}}__tmp := {{else}} {{.Out}} = {{end}} {{if .Type.Cast}}( {{.Type.Lang}} ) ( {{end}} ( {{.In}} ) . {{.Type.JsMethod}} () {{if .Type.Cast}} ) {{end}}
	{{if .Info.Pointer}} {{.Out}} = &__tmp {{end}}
{{end}}
{{define "type-callback"}}	{{.Out}} = {{.Info.Def}}FromJS( {{.In}} ) {{end}}
{{define "type-enum"}}
	{{if .Info.Pointer}}__tmp := {{else}} {{.Out}} = {{end}} {{.Info.Def}}FromJS( {{.In}} )
	{{if .Info.Pointer}} {{.Out}} = &__tmp {{end}}
{{end}}
{{define "type-interface"}}	{{.Out}} = {{.Info.Def}}FromJS( {{.In}} ) {{end}}
{{define "type-union"}}  {{.Out}} = {{.Info.Def}}FromJS( {{.In}} ) {{end}}
{{define "type-any"}}    {{.Out}} = {{.In}} {{end}}
{{define "type-typedarray"}} {{.Out}} =  jsarray.{{.GoFunc}}ToGo ( {{.In}} ) {{end}}
{{define "type-parametrized"}}	{{.Out}} = {{.Info.Def}}FromJS( {{.In}} ) {{end}}
{{define "type-dictionary"}}	{{.Out}} = {{.Info.Def}}FromJS( {{.In}} ) {{end}}
{{define "type-rawjs"}}    {{.Out}} = {{.In}} {{end}}

{{define "type-sequence"}}
	__length{{.Idx}} := {{.In}}.Length()
	__array{{.Idx}} := make( {{.Var}} , __length{{.Idx}}, __length{{.Idx}} )
	for __idx{{.Idx}} := 0; __idx{{.Idx}} < __length{{.Idx}} ; __idx{{.Idx}} ++ {
		var __seq_out{{.Idx}} {{.VarInner}}
		__seq_in{{.Idx}} := {{.In}}.Index( __idx{{.Idx}} )
		{{.Inner}}
		__array{{.Idx}}[ __idx{{.Idx}} ] = __seq_out{{.Idx}}
	}
	{{.Out}} = {{if .Info.Pointer}} & {{end}} __array{{.Idx}}
{{end}}
{{define "type-variadic"}}
	{{.Out}} = make( {{.Var}} , 0, len( {{.In}} ))
	for _, __in := range {{.In}} {
		var __out {{.VarInner}}
		{{.Inner}}
		{{.Out}} = append({{.Out}}, __out)
	} 
{{end}}
`

var inoutToTmpl = template.Must(template.New("inout-to").Parse(inoutToTmplInput))
var inoutFromTmpl = template.Must(template.New("inout-from").Parse(inoutFromTmplInput))

type useInOut int

const (
	useIn useInOut = iota
	useOut
)

type inoutData struct {
	Params    string
	ParamList []inoutParam
	AllOut    string

	// ReleaseHdl indicate that some input parameter require a returning
	// release handle
	ReleaseHdl bool
}

type inoutParam struct {
	// IDl variable name
	Name string
	// Info about the type
	Info *types.TypeInfo
	// template name
	Tmpl string
	// input variable during convert to/from wasm
	In string
	// output variable during convert to/from wasm
	Out string

	// Param references input parameter
	Param *types.Parameter

	// Inner type definintion
	Type types.TypeRef

	Var string
}

func parameterArgumentLine(input []*types.Parameter) (all string, list []string) {
	for _, value := range input {
		info, _ := value.Type.Param(false, value.Optional, value.Variadic)
		name := value.Name + " " + info.Output
		list = append(list, name)
	}
	all = strings.Join(list, ", ")
	return
}

func setupInOutWasmData(params []*types.Parameter, in, out string, use useInOut) *inoutData {
	paramTextList := []string{}
	paramList := []inoutParam{}
	allout := []string{}
	releaseHdl := false
	for idx, pi := range params {
		po := inoutParam{
			Name:  pi.Name,
			Param: pi,
			In:    setupVarName(in, idx, pi.Name, pi.Variadic),
			Out:   setupVarName(out, idx, pi.Name, pi.Variadic),
		}
		out := po.Out
		// if pi.Variadic {
		// 	out = setupVarName(out, idx, pi.Name, false) + "..."
		// }
		po.Info, po.Type = pi.Type.Param(false, pi.Optional, pi.Variadic)
		po.Tmpl = po.Info.Template
		po.Var = po.Info.VarOut
		if use == useIn {
			po.Var = po.Info.VarIn
		}
		releaseHdl = releaseHdl || pi.Type.NeedRelease()
		paramList = append(paramList, po)
		paramType := po.Info.Output
		if use == useIn {
			paramType = po.Info.Input
		}
		paramTextList = append(paramTextList, fmt.Sprint(pi.Name, " ", paramType))
		allout = append(allout, out)
	}
	return &inoutData{
		ParamList:  paramList,
		Params:     strings.Join(paramTextList, ", "),
		ReleaseHdl: releaseHdl,
		AllOut:     strings.Join(allout, ", "),
	}
}

func setupInOutWasmForOne(param *types.Parameter, in, out string, use useInOut) *inoutData {
	idx := 0
	pi := param
	po := inoutParam{
		Name:  pi.Name,
		Param: pi,
		In:    setupVarName(in, idx, pi.Name, pi.Variadic),
		Out:   setupVarName(out, idx, pi.Name, pi.Variadic),
	}
	po.Info, po.Type = pi.Type.Param(false, pi.Optional, pi.Variadic)
	po.Tmpl = po.Info.Template
	po.Var = po.Info.VarOut
	if use == useIn {
		po.Var = po.Info.VarIn
	}
	return &inoutData{
		ParamList:  []inoutParam{po},
		Params:     fmt.Sprint(pi.Name, " ", po.Info.Input),
		ReleaseHdl: pi.Type.NeedRelease(),
		AllOut:     po.Out,
	}
}
func setupInOutWasmForType(t types.TypeRef, name, in, out string, use useInOut) *inoutData {
	pi := types.Parameter{
		Name:     name,
		Optional: false,
		Variadic: false,
		Type:     t,
	}
	return setupInOutWasmForOne(&pi, in, out, use)
}

func setupVarName(value string, idx int, name string, variadic bool) string {
	value = strings.Replace(value, "@name@", name, -1)
	if variadic {
		value = strings.Replace(value, "@variadicSlice@", ":", -1)
	} else {
		value = strings.Replace(value, "@variadicSlice@", "", -1)
	}
	count := strings.Count(value, "%")
	switch count {
	case 0:
	case 1:
		value = fmt.Sprintf(value, idx)
	case 2:
		value = fmt.Sprintf(value, idx, idx)
	default:
		panic("invalid count")
	}
	return value
}

func writeInOutToWasm(data *inoutData, assign string, use useInOut, dst io.Writer) error {
	return writeInOutLoop(data, assign, use, inoutToTmpl, dst)
}

func writeInOutFromWasm(data *inoutData, assign string, use useInOut, dst io.Writer) error {
	return writeInOutLoop(data, assign, use, inoutFromTmpl, dst)
}

func writeInOutLoop(data *inoutData, assign string, use useInOut, tmpl *template.Template, dst io.Writer) error {
	for _, p := range data.ParamList {
		p.Var = p.Info.VarOut
		if use == useIn {
			p.Var = p.Info.VarIn
		}
	}
	if err := tmpl.ExecuteTemplate(dst, "start", data); err != nil {
		return err
	}
	for idx, p := range data.ParamList {
		start := inoutParamStart(p.Type, p.Info, p.Out, p.In, idx, use, tmpl)
		if _, err := io.WriteString(dst, start); err != nil {
			return err
		}
		code := inoutGetToFromWasm(p.Type, p.Info, p.Out, p.In, idx, use, tmpl)
		if _, err := io.WriteString(dst, code); err != nil {
			return err
		}
		av := setupVarName(assign, idx, p.Name, false)
		end := inoutParamEnd(p.Info, av, tmpl)
		if _, err := io.WriteString(dst, end); err != nil {
			return err
		}
	}
	if err := tmpl.ExecuteTemplate(dst, "end", data); err != nil {
		return err
	}
	return nil
}

func inoutGetToFromWasm(t types.TypeRef, info *types.TypeInfo, out, in string, idx int, use useInOut, tmpl *template.Template) string {
	if info == nil {
		panic("null")
		// info = t.DefaultParam()
	}

	// convert current
	data := struct {
		In, Out string
		Type    types.TypeRef
		Info    *types.TypeInfo
		Idx     int
		Inner   string
		GoFunc  string

		InnerInfo *types.TypeInfo
		InnerType types.TypeRef

		Var      string
		VarInner string
	}{
		In:       in,
		Type:     t,
		Out:      out,
		Info:     info,
		Idx:      idx,
		Var:      info.VarOut,
		VarInner: info.VarOutInner,
		GoFunc:   fixGoFuncName(t),
	}

	if use == useIn {
		data.Var, data.VarInner = info.VarIn, info.VarInInner
	}
	// sequence types need conversion of inner type
	if seq, ok := t.(*types.SequenceType); ok {
		sp := strconv.Itoa(idx)
		data.InnerInfo, data.InnerType = seq.Elem.DefaultParam()
		data.Inner = inoutGetToFromWasm(data.InnerType, data.InnerInfo, "__seq_out"+sp, "__seq_in"+sp, idx+1, use, tmpl)
	}
	if data.Info.Variadic {
		copy := *data.Info
		copy.Variadic = false
		data.Inner = inoutGetToFromWasm(data.Type, &copy, "__out", "__in", idx+1, use, tmpl)
		t = types.ChangeTemplateName(t, "variadic")
	}
	return convertType(t, data, tmpl) + "\n"
}

func inoutParamStart(t types.TypeRef, info *types.TypeInfo, out, in string, idx int, use useInOut, tmpl *template.Template) string {
	data := struct {
		Nullable bool
		Optional bool
		Info     *types.TypeInfo
		Type     types.TypeRef
		In, Out  string
		Idx      int
		AnyType  bool
		UseIn    bool
	}{
		Nullable: info.Nullable,
		Optional: info.Option,
		Info:     info,
		Type:     t,
		In:       in,
		Out:      out,
		Idx:      idx,
		UseIn:    use == useIn,
	}
	_, data.AnyType = t.(*types.AnyType)
	return executeTemplateToString("param-start", data, true, tmpl)
}

func inoutParamEnd(info *types.TypeInfo, assign string, tmpl *template.Template) string {
	if info.Variadic {
		assign = ""
	}
	data := struct {
		Nullable bool
		Optional bool
		Info     *types.TypeInfo
		Assign   string
	}{
		Nullable: info.Nullable,
		Optional: info.Option,
		Info:     info,
		Assign:   assign,
	}
	return executeTemplateToString("param-end", data, true, tmpl)
}

func executeTemplateToString(name string, data interface{}, newLine bool, tmpl *template.Template) string {
	var buf bytes.Buffer
	if err := tmpl.ExecuteTemplate(&buf, name, data); err != nil {
		panic(err)
	}
	out := buf.String()
	// out = strings.Replace(out, "\n", " ", -1)
	out = strings.TrimSpace(out)
	if newLine || strings.Contains(out, "\n") {
		out += "\n"
	}
	return out
}

func inoutDictionaryVariableStart(dict *dictionaryData, from useInOut, tmpl *template.Template) string {
	type elem struct {
		Name types.MethodName
		In   string
		Out  string
		Info *types.TypeInfo
		Var  string
	}
	data := struct {
		ParamList  []*elem
		ReleaseHdl bool
	}{}
	for _, m := range dict.Members {
		v := &elem{
			Name: m.Name,
			In:   m.toIn,
			Out:  m.toOut,
			Info: m.Type,
			Var:  m.Type.VarOut,
		}
		if from == useIn {
			v.Var = m.Type.VarIn
			v.In, v.Out = m.fromIn, m.fromOut
		}
		data.ParamList = append(data.ParamList, v)
	}
	return executeTemplateToString("start", data, true, tmpl)
}

func fixGoFuncName(t types.TypeRef) string {
	if array, ok := t.(*types.TypedArrayType); ok {
		name := array.Elem.Lang
		name = strings.ToUpper(name[0:1]) + name[1:]
		return name
	}
	return ""
}
