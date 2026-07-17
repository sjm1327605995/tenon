package ui

import "testing"

// 单行 Input 按回车必须触发 OnSubmit。此前回车在两条路上都被挡掉：
// editFocusedInput 只处理 multiline 的回车，activateFocused 又显式跳过 rnInput，
// 于是聊天框只能点按钮发送。
func TestSingleLineInputSubmitsOnEnter(t *testing.T) {
	var submitted []string
	app := func(_ struct{}) *Node {
		v, setV := UseState("")
		return Div(Style(Width(200), Height(80)),
			Input(Value(v), OnChange(setV),
				OnSubmit(func(s string) { submitted = append(submitted, s) }),
				Style(Width(180), Height(30))),
		)
	}
	h := Mount(Use(app, struct{}{}), 200, 80)
	in := h.Root().ByKind("input")
	if !in.Exists() {
		t.Fatal("没有 input 节点")
	}
	b := in.Bounds()
	if h.MouseDown(b.X+b.W/2, b.Y+b.H/2).Kind() != "input" { // MouseDown 才聚焦，ClickAt 只触发 onClick
		t.Fatal("按下没有落在 input 上")
	}
	in.Type("hi")

	if len(submitted) != 0 {
		t.Fatalf("还没按回车就提交了: %v", submitted)
	}
	h.Enter()
	if len(submitted) != 1 || submitted[0] != "hi" {
		t.Fatalf("回车后 submitted = %v, 期望 [hi]", submitted)
	}
}

// 多行 Input 的回车是换行，绝不能提交。
func TestMultilineInputDoesNotSubmit(t *testing.T) {
	submitted := 0
	app := func(_ struct{}) *Node {
		v, setV := UseState("")
		return Div(Style(Width(200), Height(80)),
			Input(Value(v), OnChange(setV), Multiline(),
				OnSubmit(func(string) { submitted++ }),
				Style(Width(180), Height(60))),
		)
	}
	h := Mount(Use(app, struct{}{}), 200, 80)
	in := h.Root().ByKind("input")
	b := in.Bounds()
	h.MouseDown(b.X+b.W/2, b.Y+b.H/2)
	in.Type("a")
	h.Enter()
	if submitted != 0 {
		t.Errorf("多行 Input 的回车触发了 %d 次提交 —— 那里回车该是换行", submitted)
	}
}

// 没挂 OnSubmit 时回车必须是安全的空操作（不能 panic、不能误激活别的东西）。
func TestEnterWithoutOnSubmitIsNoop(t *testing.T) {
	clicked := 0
	app := func(_ struct{}) *Node {
		v, setV := UseState("")
		return Div(Style(Width(200), Height(80)),
			Input(Value(v), OnChange(setV), Style(Width(180), Height(30))),
			Div(Style(Width(50), Height(20)), OnClick(func() { clicked++ }), Text("btn")),
		)
	}
	h := Mount(Use(app, struct{}{}), 200, 80)
	b := h.Root().ByKind("input").Bounds()
	h.MouseDown(b.X+b.W/2, b.Y+b.H/2)
	if ok := h.Enter(); ok {
		t.Error("没挂 OnSubmit 的 Input 上回车却报告「已激活」")
	}
	if clicked != 0 {
		t.Errorf("输入框上的回车误触发了别处的 OnClick %d 次", clicked)
	}
}

// OnSubmit 拿到的必须是「按回车前刚输入的那个字符也算数」的值。
// 提交若排在 onChange 之前，受控组件的值会晚一帧，最后一个字符就丢了。
func TestSubmitSeesLatestValue(t *testing.T) {
	var got string
	app := func(_ struct{}) *Node {
		v, setV := UseState("")
		return Div(Style(Width(200), Height(80)),
			Input(Value(v), OnChange(setV),
				OnSubmit(func(s string) { got = s }), Style(Width(180), Height(30))),
		)
	}
	h := Mount(Use(app, struct{}{}), 200, 80)
	in := h.Root().ByKind("input")
	b := in.Bounds()
	h.MouseDown(b.X+b.W/2, b.Y+b.H/2)
	in.Type("abc")
	h.Enter()
	if got != "abc" {
		t.Errorf("OnSubmit 收到 %q, 期望 \"abc\"", got)
	}
}
