package components

import (
	"testing"

	"github.com/sjm1327605995/tenon/pkg/v2/core"
	"github.com/sjm1327605995/tenon/yoga"
)

type testWidgetRoot struct {
	core.BaseWidget
	root core.Element
}

func (w *testWidgetRoot) Build() core.Element { return w.root }

func newTestEvent(tpe core.EventType, x, y float32) *core.Event {
	return &core.Event{
		Type: tpe,
		X:    x,
		Y:    y,
	}
}

// ==================== Button ====================

func TestButton_ClickEvent(t *testing.T) {
	clicked := false
	btn := NewButton("Test").SetOnClick(func() {
		clicked = true
	})

	event := newTestEvent(core.EventClick, 10, 10)
	btn.HandleEvent(event)

	if !clicked {
		t.Fatal("Button onClick should be triggered")
	}
}

func TestButton_ChainAPI(t *testing.T) {
	btn := NewButton("Test")
	if btn.SetOnClick(func() {}) != btn {
		t.Error("SetOnClick should return *Button")
	}
	if btn.SetDisabled(true) != btn {
		t.Error("SetDisabled should return *Button")
	}
	if btn.SetText("hello") != btn {
		t.Error("SetText should return *Button")
	}
	if btn.SetLabel("world") != btn {
		t.Error("SetLabel should return *Button")
	}
	if btn.SetLoading(true) != btn {
		t.Error("SetLoading should return *Button")
	}
	if btn.SetColors(nil, nil, nil) != btn {
		t.Error("SetColors should return *Button")
	}
}

// ==================== Text ====================

func TestText_ChainAPI(t *testing.T) {
	txt := NewText("hello")
	if txt.SetContent("world") != txt {
		t.Error("SetContent should return *Text")
	}
	if txt.SetFontSize(20) != txt {
		t.Error("SetFontSize should return *Text")
	}
	if txt.SetColor(nil) != txt {
		t.Error("SetColor should return *Text")
	}
	if txt.SetWhiteSpace(WhiteSpaceNormal) != txt {
		t.Error("SetWhiteSpace should return *Text")
	}
	if txt.SetWordBreak(WordBreakNormal) != txt {
		t.Error("SetWordBreak should return *Text")
	}
	if txt.SetLineHeight(24) != txt {
		t.Error("SetLineHeight should return *Text")
	}
}

// ==================== ScrollView ====================

func TestScrollView_MaxScrollY(t *testing.T) {
	sv := NewScrollView()
	sv.SetWidth(200)
	sv.SetHeight(100)
	child := NewView().SetWidth(50).SetHeight(50)
	sv.Content().AppendChild(child)

	tw := &testWidgetRoot{root: sv}
	tw.Init(tw)
	engine := core.NewEngine(tw, 200, 100)
	engine.Mount()
	engine.Update()

	if sv.maxScrollY != 0 {
		t.Fatalf("expected maxScrollY=0 when content is smaller, got %f", sv.maxScrollY)
	}
}

func TestScrollView_ScrollBounds(t *testing.T) {
	sv := NewScrollView()
	sv.SetWidth(200)
	sv.SetHeight(100)
	child := NewView().SetWidth(50).SetHeight(300)
	sv.Content().AppendChild(child)

	tw := &testWidgetRoot{root: sv}
	tw.Init(tw)
	engine := core.NewEngine(tw, 200, 100)
	engine.Mount()
	engine.Update()

	if sv.maxScrollY <= 0 {
		t.Fatal("maxScrollY should be > 0 when content is larger")
	}

	// scrollY within bounds should stay unchanged
	sv.scrollY = 100
	sv.clampScroll()
	if sv.scrollY != 100 {
		t.Fatalf("scrollY should stay 100 within bounds, got %f", sv.scrollY)
	}

	// scrollY beyond max should be clamped
	sv.scrollY = 9999
	sv.clampScroll()
	if sv.scrollY != sv.maxScrollY {
		t.Fatalf("scrollY should be clamped to maxScrollY(%f), got %f", sv.maxScrollY, sv.scrollY)
	}

	// negative scrollY should be clamped to 0
	sv.scrollY = -9999
	sv.clampScroll()
	if sv.scrollY != 0 {
		t.Fatalf("scrollY should be clamped to 0 when negative, got %f", sv.scrollY)
	}
}

// ==================== Checkbox ====================

func TestCheckbox_Toggle(t *testing.T) {
	changed := false
	var newChecked bool
	cb := NewCheckbox("Option").SetOnChange(func(checked bool) {
		changed = true
		newChecked = checked
	})

	if cb.checked {
		t.Fatal("Checkbox should start unchecked")
	}

	event := newTestEvent(core.EventClick, 10, 10)
	cb.HandleEvent(event)

	if !cb.checked {
		t.Fatal("Checkbox should be checked after click")
	}
	if !changed {
		t.Fatal("onChange should be triggered")
	}
	if !newChecked {
		t.Fatal("onChange should receive true")
	}
}

func TestCheckbox_ChainAPI(t *testing.T) {
	cb := NewCheckbox("Option")
	if cb.SetChecked(true) != cb {
		t.Error("SetChecked should return *Checkbox")
	}
	if cb.SetOnChange(func(bool) {}) != cb {
		t.Error("SetOnChange should return *Checkbox")
	}
	if cb.SetBoxSize(20) != cb {
		t.Error("SetBoxSize should return *Checkbox")
	}
	if cb.SetBorderColor(nil) != cb {
		t.Error("SetBorderColor should return *Checkbox")
	}
	if cb.SetFillColor(nil) != cb {
		t.Error("SetFillColor should return *Checkbox")
	}
	if cb.SetCheckColor(nil) != cb {
		t.Error("SetCheckColor should return *Checkbox")
	}
	if cb.SetTextColor(nil) != cb {
		t.Error("SetTextColor should return *Checkbox")
	}
	if cb.SetFontSize(16) != cb {
		t.Error("SetFontSize should return *Checkbox")
	}
}

// ==================== Radio ====================

func TestRadio_Select(t *testing.T) {
	changed := false
	var newSelected bool
	r := NewRadio("Option").SetOnChange(func(selected bool) {
		changed = true
		newSelected = selected
	})

	if r.selected {
		t.Fatal("Radio should start unselected")
	}

	event := newTestEvent(core.EventClick, 10, 10)
	r.HandleEvent(event)

	if !r.selected {
		t.Fatal("Radio should be selected after click")
	}
	if !changed {
		t.Fatal("onChange should be triggered")
	}
	if !newSelected {
		t.Fatal("onChange should receive true")
	}
}

func TestRadio_ChainAPI(t *testing.T) {
	r := NewRadio("Option")
	if r.SetSelected(true) != r {
		t.Error("SetSelected should return *Radio")
	}
	if r.SetOnChange(func(bool) {}) != r {
		t.Error("SetOnChange should return *Radio")
	}
	if r.SetBoxSize(20) != r {
		t.Error("SetBoxSize should return *Radio")
	}
	if r.SetBorderColor(nil) != r {
		t.Error("SetBorderColor should return *Radio")
	}
	if r.SetFillColor(nil) != r {
		t.Error("SetFillColor should return *Radio")
	}
	if r.SetInnerColor(nil) != r {
		t.Error("SetInnerColor should return *Radio")
	}
	if r.SetTextColor(nil) != r {
		t.Error("SetTextColor should return *Radio")
	}
	if r.SetFontSize(16) != r {
		t.Error("SetFontSize should return *Radio")
	}
}

// ==================== Switch ====================

func TestSwitch_Toggle(t *testing.T) {
	changed := false
	var newChecked bool
	sw := NewSwitch().SetOnChange(func(checked bool) {
		changed = true
		newChecked = checked
	})

	if sw.checked {
		t.Fatal("Switch should start unchecked")
	}

	event := newTestEvent(core.EventClick, 10, 10)
	sw.HandleEvent(event)

	if !sw.checked {
		t.Fatal("Switch should be checked after click")
	}
	if !changed {
		t.Fatal("onChange should be triggered")
	}
	if !newChecked {
		t.Fatal("onChange should receive true")
	}
}

func TestSwitch_ChainAPI(t *testing.T) {
	sw := NewSwitch()
	if sw.SetChecked(true) != sw {
		t.Error("SetChecked should return *Switch")
	}
	if sw.SetOnChange(func(bool) {}) != sw {
		t.Error("SetOnChange should return *Switch")
	}
	if sw.SetOffColor(nil) != sw {
		t.Error("SetOffColor should return *Switch")
	}
	if sw.SetOnColor(nil) != sw {
		t.Error("SetOnColor should return *Switch")
	}
	if sw.SetThumbColor(nil) != sw {
		t.Error("SetThumbColor should return *Switch")
	}
	if sw.SetTrackSize(50, 30) != sw {
		t.Error("SetTrackSize should return *Switch")
	}
}

// ==================== Slider ====================

func TestSlider_ValueBounds(t *testing.T) {
	s := NewSlider(0, 100).SetValue(50)
	if s.value != 50 {
		t.Fatalf("expected value=50, got %f", s.value)
	}

	s.SetValue(-10)
	if s.value != 0 {
		t.Fatalf("expected value=0 after underflow, got %f", s.value)
	}

	s.SetValue(200)
	if s.value != 100 {
		t.Fatalf("expected value=100 after overflow, got %f", s.value)
	}
}

func TestSlider_UpdateValueFromX(t *testing.T) {
	s := NewSlider(0, 100)
	bounds := core.LayoutBounds{X: 0, Y: 0, Width: 200, Height: 24}

	s.updateValueFromX(0, bounds)
	if s.value != 0 {
		t.Fatalf("expected value=0 at left edge, got %f", s.value)
	}

	s.updateValueFromX(200, bounds)
	if s.value != 100 {
		t.Fatalf("expected value=100 at right edge, got %f", s.value)
	}

	s.updateValueFromX(100, bounds)
	if s.value != 50 {
		t.Fatalf("expected value=50 at middle, got %f", s.value)
	}

	s.updateValueFromX(-50, bounds)
	if s.value != 0 {
		t.Fatalf("expected value=0 beyond left, got %f", s.value)
	}

	s.updateValueFromX(300, bounds)
	if s.value != 100 {
		t.Fatalf("expected value=100 beyond right, got %f", s.value)
	}
}

func TestSlider_ChainAPI(t *testing.T) {
	s := NewSlider(0, 100)
	if s.SetValue(50) != s {
		t.Error("SetValue should return *Slider")
	}
	if s.SetOnChange(func(float32) {}) != s {
		t.Error("SetOnChange should return *Slider")
	}
	if s.SetTrackColor(nil) != s {
		t.Error("SetTrackColor should return *Slider")
	}
	if s.SetFillColor(nil) != s {
		t.Error("SetFillColor should return *Slider")
	}
	if s.SetThumbColor(nil) != s {
		t.Error("SetThumbColor should return *Slider")
	}
	if s.SetThumbRadius(10) != s {
		t.Error("SetThumbRadius should return *Slider")
	}
	if s.SetTrackHeight(6) != s {
		t.Error("SetTrackHeight should return *Slider")
	}
}

// ==================== TextInput ====================

func TestTextInput_ChainAPI(t *testing.T) {
	ti := NewTextInput()
	if ti.SetText("hello") != ti {
		t.Error("SetText should return *TextInput")
	}
	if ti.SetPlaceholder("hint") != ti {
		t.Error("SetPlaceholder should return *TextInput")
	}
	if ti.SetMultiline(true) != ti {
		t.Error("SetMultiline should return *TextInput")
	}
	if ti.SetOnChange(func(string) {}) != ti {
		t.Error("SetOnChange should return *TextInput")
	}
	if ti.SetOnSubmit(func(string) {}) != ti {
		t.Error("SetOnSubmit should return *TextInput")
	}
	if ti.SetFontSize(16) != ti {
		t.Error("SetFontSize should return *TextInput")
	}
	if ti.SetTextColor(nil) != ti {
		t.Error("SetTextColor should return *TextInput")
	}
	if ti.SetPlaceholderColor(nil) != ti {
		t.Error("SetPlaceholderColor should return *TextInput")
	}
	if ti.SetBorderColor(nil) != ti {
		t.Error("SetBorderColor should return *TextInput")
	}
	if ti.SetFocusBorderColor(nil) != ti {
		t.Error("SetFocusBorderColor should return *TextInput")
	}
	if ti.SetBackgroundColor(nil) != ti {
		t.Error("SetBackgroundColor should return *TextInput")
	}
	if ti.SetInnerPadding(10) != ti {
		t.Error("SetInnerPadding should return *TextInput")
	}
}

// ==================== View layout chain ====================

func TestView_ChainLayoutAPI(t *testing.T) {
	v := NewView()
	_ = v.SetWidth(100)
	_ = v.SetHeight(50)
	_ = v.SetMargin(yoga.EdgeAll, 10)
	_ = v.SetPadding(yoga.EdgeAll, 8)
	_ = v.SetBackgroundColor(nil)
	_ = v.SetBorderRadius(4)
}
