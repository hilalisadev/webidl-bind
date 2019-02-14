package types

import (
	"github.com/gowebapi/webidlparser/ast"
)

type Dictionary struct {
	standardType
	basic    BasicInfo
	source   *ast.Dictionary
	Inherits *Dictionary

	Members []*DictMember
}

// Dictionary need to implement Type
var _ Type = &Dictionary{}

type DictMember struct {
	nameAndLink
	Type     TypeRef
	Required bool
}

func (t *extractTypes) convertDictionary(in *ast.Dictionary) (*Dictionary, bool) {
	ref := createRef(in, t)
	t.assertTrue(len(in.Annotations) == 0, ref, "unsupported annotations")
	// t.assertTrue(in.Inherits == "", ref , "unsupported dictionary inherites of %s", in.Inherits)
	ret := &Dictionary{
		standardType: standardType{
			ref:         ref,
			needRelease: false,
		},
		basic:  fromIdlToTypeName(t.main.setup.Package, in.Name, "dictionary"),
		source: in,
	}
	for _, mi := range in.Members {
		mo := t.convertDictMember(mi)
		ret.Members = append(ret.Members, mo)
	}
	return ret, in.Partial
}

func (conv *extractTypes) convertDictMember(in *ast.Member) *DictMember {
	ref := createRef(in, conv)
	conv.assertTrue(!in.Readonly, ref, "read only not allowed")
	conv.assertTrue(in.Attribute, ref, "must be an attribute")
	conv.assertTrue(!in.Static, ref, "static is not allowed")
	conv.assertTrue(!in.Const, ref, "const is not allowed")
	conv.assertTrue(len(in.Parameters) == 0, ref, "parameters on member is not allowed (or not supported)")
	conv.assertTrue(len(in.Specialization) == 0, ref, "specialization on member is not allowed (or not supported)")
	conv.warningTrue(!in.Required, ref, "required value not implemented yet, report this as a bug :)")
	for _, a := range in.Annotations {
		ref := createRef(a, conv)
		conv.warning(ref, "dictionary member: annotation '%s' is not supported", a.Name)
	}
	if in.Init != nil {
		conv.warning(ref, "dictionary: default value for dictionary not implemented yet")
		// parser.Dump(os.Stdout, in)
	}
	return &DictMember{
		nameAndLink: nameAndLink{
			ref:  createRef(in, conv),
			name: fromIdlToMethodName(in.Name),
		},
		Type:     convertType(in.Type, conv),
		Required: in.Required,
	}
}

func (t *Dictionary) Basic() BasicInfo {
	return TransformBasic(t, t.basic)
}

func (t *Dictionary) DefaultParam() (info *TypeInfo, inner TypeRef) {
	return t.Param(false, false, false)
}

func (t *Dictionary) key() string {
	return t.basic.Idl
}

func (t *Dictionary) link(conv *Convert, inuse inuseLogic) TypeRef {
	if t.inuse {
		return t
	}
	t.inuse = true
	inner := make(inuseLogic)
	for _, m := range t.Members {
		m.Type = m.Type.link(conv, inner)
	}
	return t
}

func (t *Dictionary) merge(partial *Dictionary, conv *Convert) {
	conv.assertTrue(partial.source.Inherits == "", partial, "unsupported dictionary inherites on partial")
	// TODO member elemination logic with duplicate is detected
	t.Members = append(t.Members, partial.Members...)
}

func (t *Dictionary) NeedRelease() bool {
	need := false
	for _, v := range t.Members {
		need = need || v.Type.NeedRelease()
	}
	return need
}

func (t *Dictionary) Param(nullable, option, variadic bool) (info *TypeInfo, inner TypeRef) {
	return newTypeInfo(t.Basic(), nullable, option, variadic, true, false, false), t
}

func (t *Dictionary) SetBasic(basic BasicInfo) {
	t.basic = basic
}


func (t *Dictionary) TypeID() TypeID {
	return TypeDictionary
}
