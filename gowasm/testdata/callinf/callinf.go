// Code generated by webidlgenerator. DO NOT EDIT.

// +build !js

package callinf

import js "github.com/gowebapi/webapi/core/failjs"

// using following types:

// ReleasableApiResource is used to release underlaying
// allocated resources.
type ReleasableApiResource interface {
	Release()
}

type releasableApiResourceList []ReleasableApiResource

func (a releasableApiResourceList) Release() {
	for _, v := range a {
		v.Release()
	}
}

// workaround for compiler error
func unused(value interface{}) {
	// TODO remove this method
}

type Union struct {
	Value js.Value
}

func (u *Union) JSValue() js.Value {
	return u.Value
}

func UnionFromJS(value js.Value) *Union {
	return &Union{Value: value}
}

const Test1_Foo1 string = 2

// Foo1 is a callback interface.
type Foo1 interface {
	Test2(b int, c int)
}

// Foo1Value is javascript reference value for callback interface Foo1.
// This is holding the underlaying javascript object.
type Foo1Value struct {
	// Value is the underlying javascript object or function.
	Value js.Value
	// Functions is the underlying function objects that is allocated for the interface callback
	Functions [1]js.Func
	// Go interface to invoke
	impl      Foo1
	function  func(b int, c int)
	useInvoke bool
}

// JSValue is returning the javascript object that implements this callback interface
func (t *Foo1Value) JSValue() js.Value {
	return t.Value
}

// Release is releasing all resources that is allocated.
func (t *Foo1Value) Release() {
	for i := range t.Functions {
		if t.Functions[i].Type() != js.TypeUndefined {
			t.Functions[i].Release()
		}
	}
}

// NewFoo1 is allocating a new javascript object that
// implements Foo1.
func NewFoo1(callback Foo1) *Foo1Value {
	ret := &Foo1Value{impl: callback}
	ret.Value = js.Global().Get("Object").New()
	ret.Functions[0] = ret.allocateTest2()
	ret.Value.Set("test2", ret.Functions[0])
	return ret
}

// NewFoo1Func is allocating a new javascript
// function is implements
// Foo1 interface.
func NewFoo1Func(f func(b int, c int)) *Foo1Value {
	// single function will result in javascript function type, not an object
	ret := &Foo1Value{function: f}
	ret.Functions[0] = ret.allocateTest2()
	ret.Value = ret.Functions[0].Value
	return ret
}

// Foo1FromJS is taking an javascript object that reference to a
// callback interface and return a corresponding interface that can be used
// to invoke on that element.
func Foo1FromJS(value js.Wrapper) *Foo1Value {
	input := value.JSValue()
	if input.Type() == js.TypeObject {
		return &Foo1Value{Value: input}
	}
	if input.Type() == js.TypeFunction {
		return &Foo1Value{Value: input, useInvoke: true}
	}
	panic("unsupported type")
}

func (t *Foo1Value) allocateTest2() js.Func {
	return js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		var (
			_p0 int // javascript: long b
			_p1 int // javascript: long c
		)
		_p0 = (args[0]).Int()
		_p1 = (args[1]).Int()
		if t.function != nil {
			t.function(_p0, _p1)
		} else {
			t.impl.Test2(_p0, _p1)
		}
		// returning no return value
		return nil
	})
}

func (_this *Foo1Value) Test2(b int, c int) {
	if _this.function != nil {
		_this.function(b, c)
	}
	if _this.impl != nil {
		_this.impl.Test2(b, c)
	}
	var (
		_args [2]interface{}
		_end  int
	)
	_p0 := b
	_args[0] = _p0
	_end++
	_p1 := c
	_args[1] = _p1
	_end++
	if _this.useInvoke {
		// invoke a javascript function
		_this.Value.Invoke(_args[0:_end]...)
	} else {
		_this.Value.Call("test2", _args[0:_end]...)
	}
	return
}

// Foo2 is a callback interface.
type Foo2 interface {
	Test3(a string, b js.Value, c *Union, d int, e *A, f *B) (_result bool)
}

// Foo2Value is javascript reference value for callback interface Foo2.
// This is holding the underlaying javascript object.
type Foo2Value struct {
	// Value is the underlying javascript object or function.
	Value js.Value
	// Functions is the underlying function objects that is allocated for the interface callback
	Functions [1]js.Func
	// Go interface to invoke
	impl      Foo2
	function  func(a string, b js.Value, c *Union, d int, e *A, f *B) (_result bool)
	useInvoke bool
}

// JSValue is returning the javascript object that implements this callback interface
func (t *Foo2Value) JSValue() js.Value {
	return t.Value
}

// Release is releasing all resources that is allocated.
func (t *Foo2Value) Release() {
	for i := range t.Functions {
		if t.Functions[i].Type() != js.TypeUndefined {
			t.Functions[i].Release()
		}
	}
}

// NewFoo2 is allocating a new javascript object that
// implements Foo2.
func NewFoo2(callback Foo2) *Foo2Value {
	ret := &Foo2Value{impl: callback}
	ret.Value = js.Global().Get("Object").New()
	ret.Functions[0] = ret.allocateTest3()
	ret.Value.Set("test3", ret.Functions[0])
	return ret
}

// NewFoo2Func is allocating a new javascript
// function is implements
// Foo2 interface.
func NewFoo2Func(f func(a string, b js.Value, c *Union, d int, e *A, f *B) (_result bool)) *Foo2Value {
	// single function will result in javascript function type, not an object
	ret := &Foo2Value{function: f}
	ret.Functions[0] = ret.allocateTest3()
	ret.Value = ret.Functions[0].Value
	return ret
}

// Foo2FromJS is taking an javascript object that reference to a
// callback interface and return a corresponding interface that can be used
// to invoke on that element.
func Foo2FromJS(value js.Wrapper) *Foo2Value {
	input := value.JSValue()
	if input.Type() == js.TypeObject {
		return &Foo2Value{Value: input}
	}
	if input.Type() == js.TypeFunction {
		return &Foo2Value{Value: input, useInvoke: true}
	}
	panic("unsupported type")
}

func (t *Foo2Value) allocateTest3() js.Func {
	return js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		var (
			_p0 string   // javascript: DOMString a
			_p1 js.Value // javascript: any b
			_p2 *Union   // javascript: Union c
			_p3 int      // javascript: long d
			_p4 *A       // javascript: A e
			_p5 *B       // javascript: B f
		)
		_p0 = (args[0]).String()
		_p1 = args[1]
		_p2 = UnionFromJS(args[2])
		_p3 = (args[3]).Int()
		_p4 = AFromJS(args[4])
		_p5 = BFromJS(args[5])
		var _returned bool
		if t.function != nil {
			_returned = t.function(_p0, _p1, _p2, _p3, _p4, _p5)
		} else {
			_returned = t.impl.Test3(_p0, _p1, _p2, _p3, _p4, _p5)
		}
		_converted := _returned
		return _converted
	})
}

func (_this *Foo2Value) Test3(a string, b js.Value, c *Union, d int, e *A, f *B) (_result bool) {
	if _this.function != nil {
		return _this.function(a, b, c, d, e, f)
	}
	if _this.impl != nil {
		return _this.impl.Test3(a, b, c, d, e, f)
	}
	var (
		_args [6]interface{}
		_end  int
	)
	_p0 := a
	_args[0] = _p0
	_end++
	_p1 := b
	_args[1] = _p1
	_end++
	_p2 := c.JSValue()
	_args[2] = _p2
	_end++
	_p3 := d
	_args[3] = _p3
	_end++
	_p4 := e.JSValue()
	_args[4] = _p4
	_end++
	_p5 := f.JSValue()
	_args[5] = _p5
	_end++
	var _returned js.Value
	if _this.useInvoke {
		// invoke a javascript function
		_returned = _this.Value.Invoke(_args[0:_end]...)
	} else {
		_returned = _this.Value.Call("test3", _args[0:_end]...)
	}
	var (
		_converted bool // javascript: boolean _what_return_name
	)
	_converted = (_returned).Bool()
	_result = _converted
	return
}

// Foo3 is a callback interface.
type Foo3 interface {
	Test1(a string)
	Test2(b int) (_result int)
	Test3(c bool) (_result bool)
	Test4(d js.Value) (_result js.Value)
}

// Foo3Value is javascript reference value for callback interface Foo3.
// This is holding the underlaying javascript object.
type Foo3Value struct {
	// Value is the underlying javascript object or function.
	Value js.Value
	// Functions is the underlying function objects that is allocated for the interface callback
	Functions [4]js.Func
	// Go interface to invoke
	impl Foo3
}

// JSValue is returning the javascript object that implements this callback interface
func (t *Foo3Value) JSValue() js.Value {
	return t.Value
}

// Release is releasing all resources that is allocated.
func (t *Foo3Value) Release() {
	for i := range t.Functions {
		if t.Functions[i].Type() != js.TypeUndefined {
			t.Functions[i].Release()
		}
	}
}

// NewFoo3 is allocating a new javascript object that
// implements Foo3.
func NewFoo3(callback Foo3) *Foo3Value {
	ret := &Foo3Value{impl: callback}
	ret.Value = js.Global().Get("Object").New()
	ret.Functions[0] = ret.allocateTest1()
	ret.Value.Set("test1", ret.Functions[0])
	ret.Functions[1] = ret.allocateTest2()
	ret.Value.Set("test2", ret.Functions[1])
	ret.Functions[2] = ret.allocateTest3()
	ret.Value.Set("test3", ret.Functions[2])
	ret.Functions[3] = ret.allocateTest4()
	ret.Value.Set("test4", ret.Functions[3])
	return ret
}

// Foo3FromJS is taking an javascript object that reference to a
// callback interface and return a corresponding interface that can be used
// to invoke on that element.
func Foo3FromJS(value js.Wrapper) *Foo3Value {
	input := value.JSValue()
	if input.Type() == js.TypeObject {
		return &Foo3Value{Value: input}
	}
	if input.Type() == js.TypeFunction {
		return &Foo3Value{Value: input, useInvoke: true}
	}
	panic("unsupported type")
}

func (t *Foo3Value) allocateTest1() js.Func {
	return js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		var (
			_p0 string // javascript: DOMString a
		)
		_p0 = (args[0]).String()
		t.impl.Test1(_p0)
		// returning no return value
		return nil
	})
}

func (t *Foo3Value) allocateTest2() js.Func {
	return js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		var (
			_p0 int // javascript: long b
		)
		_p0 = (args[0]).Int()
		var _returned int
		_returned = t.impl.Test2(_p0)
		_converted := _returned
		return _converted
	})
}

func (t *Foo3Value) allocateTest3() js.Func {
	return js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		var (
			_p0 bool // javascript: boolean c
		)
		_p0 = (args[0]).Bool()
		var _returned bool
		_returned = t.impl.Test3(_p0)
		_converted := _returned
		return _converted
	})
}

func (t *Foo3Value) allocateTest4() js.Func {
	return js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		var (
			_p0 js.Value // javascript: any d
		)
		_p0 = args[0]
		var _returned js.Value
		_returned = t.impl.Test4(_p0)
		_converted := _returned
		return _converted
	})
}

func (_this *Foo3Value) Test1(a string) {
	if _this.impl != nil {
		_this.impl.Test1(a)
	}
	var (
		_args [1]interface{}
		_end  int
	)
	_p0 := a
	_args[0] = _p0
	_end++
	_this.Value.Call("test1", _args[0:_end]...)
	return
}

func (_this *Foo3Value) Test2(b int) (_result int) {
	if _this.impl != nil {
		return _this.impl.Test2(b)
	}
	var (
		_args [1]interface{}
		_end  int
	)
	_p0 := b
	_args[0] = _p0
	_end++
	var _returned js.Value
	_returned = _this.Value.Call("test2", _args[0:_end]...)
	var (
		_converted int // javascript: long _what_return_name
	)
	_converted = (_returned).Int()
	_result = _converted
	return
}

func (_this *Foo3Value) Test3(c bool) (_result bool) {
	if _this.impl != nil {
		return _this.impl.Test3(c)
	}
	var (
		_args [1]interface{}
		_end  int
	)
	_p0 := c
	_args[0] = _p0
	_end++
	var _returned js.Value
	_returned = _this.Value.Call("test3", _args[0:_end]...)
	var (
		_converted bool // javascript: boolean _what_return_name
	)
	_converted = (_returned).Bool()
	_result = _converted
	return
}

func (_this *Foo3Value) Test4(d js.Value) (_result js.Value) {
	if _this.impl != nil {
		return _this.impl.Test4(d)
	}
	var (
		_args [1]interface{}
		_end  int
	)
	_p0 := d
	_args[0] = _p0
	_end++
	var _returned js.Value
	_returned = _this.Value.Call("test4", _args[0:_end]...)
	var (
		_converted js.Value // javascript: any _what_return_name
	)
	_converted = _returned
	_result = _converted
	return
}

// interface: A
type A struct {
	// Value_JS holds a reference to a javascript value
	Value_JS js.Value
}

func (_this *A) JSValue() js.Value {
	return _this.Value_JS
}

// AFromJS is casting a js.Wrapper into A.
func AFromJS(value js.Wrapper) *A {
	input := value.JSValue()
	if input.Type() == js.TypeNull {
		return nil
	}
	ret := &A{}
	ret.Value_JS = input
	return ret
}

// interface: B
type B struct {
	// Value_JS holds a reference to a javascript value
	Value_JS js.Value
}

func (_this *B) JSValue() js.Value {
	return _this.Value_JS
}

// BFromJS is casting a js.Wrapper into B.
func BFromJS(value js.Wrapper) *B {
	input := value.JSValue()
	if input.Type() == js.TypeNull {
		return nil
	}
	ret := &B{}
	ret.Value_JS = input
	return ret
}
