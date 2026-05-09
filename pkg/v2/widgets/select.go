package widgets

import (
	"image/color"

	"github.com/sjm1327605995/tenon/pkg/v2/render"
	"github.com/sjm1327605995/tenon/pkg/v2/ui"
)

func colorToRender(c color.Color) render.Color {
	if c == nil {
		return render.Color{}
	}
	return *render.NewColorFrom(c)
}

// SelectOption 下拉选项。
type SelectOption struct {
	Value    string
	Label    string
	Disabled bool
}

// SelectWidget 下拉选择器。
type SelectWidget struct {
	ui.BaseWidget
	Options     []SelectOption
	Value       string
	Placeholder string
	OnChange    func(string)
	Disabled    bool
	Width       float32
}

// Select 创建下拉选择器。
func Select(options []SelectOption) SelectWidget {
	return SelectWidget{Options: options, Placeholder: "请选择...", Width: 200}
}

func (s SelectWidget) WithValue(v string) SelectWidget           { return s.SelectedValue(v) }
func (s SelectWidget) SelectedValue(v string) SelectWidget       { s.Value = v; return s }
func (s SelectWidget) WithPlaceholder(v string) SelectWidget     { return s.PlaceholderText(v) }
func (s SelectWidget) PlaceholderText(v string) SelectWidget     { s.Placeholder = v; return s }
func (s SelectWidget) WithOnChange(fn func(string)) SelectWidget { return s.OnChanged(fn) }
func (s SelectWidget) OnChanged(fn func(string)) SelectWidget    { s.OnChange = fn; return s }
func (s SelectWidget) WithDisabled(v bool) SelectWidget          { return s.SetDisabled(v) }
func (s SelectWidget) SetDisabled(v bool) SelectWidget           { s.Disabled = v; return s }
func (s SelectWidget) WithWidth(v float32) SelectWidget          { return s.Width_(v) }
func (s SelectWidget) Width_(v float32) SelectWidget             { s.Width = v; return s }

func (s SelectWidget) CreateElement() ui.Element { return ui.NewStatefulElement(s) }
func (s SelectWidget) CreateState() ui.State {
	st := &selectState{}
	st.Init(st)
	return st
}

type selectState struct {
	ui.BaseStateOf[SelectWidget]
	open        bool
	dismisserID int
	value       string
}

func (s *selectState) InitState() {
	s.dismisserID = ui.RegisterPopupDismisser(func() {
		if s.open {
			s.open = false
			s.SetState(nil)
		}
	})
	w := s.Widget()
	s.value = w.Value
}
func (s *selectState) Dispose() {
	ui.UnregisterPopupDismisser(s.dismisserID)
}
func (s *selectState) DidUpdateWidget(old ui.Widget) {
	oldW := old.(SelectWidget)
	w := s.Widget()
	if oldW.Value != w.Value {
		s.value = w.Value
	}
}

func (s *selectState) Build(ctx ui.BuildContext) ui.Widget {
	w := s.Widget()
	theme := ui.ThemeOf(ctx)
	trigger := s.buildTrigger(w, theme)
	if !s.open {
		return Stack(trigger).Z(0)
	}
	return Stack(
		trigger,
		Positioned(s.buildDropdown(w, theme)).T(40).L(0).W(w.Width).Z(999),
	).Z(0)
}

func (s *selectState) buildTrigger(w SelectWidget, theme *ui.Theme) ui.Widget {
	text := w.Placeholder
	clr := theme.PlaceholderColor
	if s.value != "" {
		text = findOptLabel(s.value, w.Options)
		clr = theme.TextColor
	}
	arrow := IconArrowDown
	if s.open {
		arrow = IconArrowUp
	}
	return Container(
		Row(
			Text(text).FontSize(theme.FontSizeBase).Color(clr),
			Icon(arrow).FontSize(10).Color(theme.TextMutedColor),
		).Gapf(8).JustifyContent(ui.JustifySpaceBetween).AlignItems(ui.AlignCenter),
	).
		Background(colorToRender(theme.BackgroundColor)).
		Border(colorToRender(theme.BorderColor), 1).
		Radius(theme.BorderRadius).
		Pad(ui.EdgeInsets{Top: 8, Right: 12, Bottom: 8, Left: 12}).
		W(w.Width).
		H(40).
		Focusable(true).
		OnTap(func() {
			if !w.Disabled {
				if !s.open {
					ui.DismissAllPopups()
				}
				s.open = !s.open
				s.SetState(nil)
			}
		})
}

func (s *selectState) buildDropdown(w SelectWidget, theme *ui.Theme) ui.Widget {
	opts := make([]ui.Widget, 0, len(w.Options))
	for _, opt := range w.Options {
		opts = append(opts, s.buildOpt(opt, w, theme))
	}
	list := Scroll(Column(opts...).Gapf(0)).W(w.Width).MaxH(200)
	return Container(list).
		Background(colorToRender(theme.BackgroundColor)).
		Border(colorToRender(theme.BorderColor), 1).
		Radius(theme.BorderRadius).
		Pad(ui.EdgeInsets{Top: 4, Bottom: 4}).
		W(w.Width)
}

func (s *selectState) buildOpt(opt SelectOption, w SelectWidget, theme *ui.Theme) ui.Widget {
	label := opt.Label
	if label == "" {
		label = opt.Value
	}
	bg := theme.BackgroundColor
	txtClr := theme.TextColor
	if opt.Disabled {
		txtClr = theme.MutedColor
	}
	if opt.Value == s.value {
		bg = theme.AccentColor
	}
	return Container(Text(label).FontSize(theme.FontSizeBase).Color(txtClr)).
		Background(colorToRender(bg)).
		Pad(ui.EdgeInsets{Top: 6, Right: 12, Bottom: 6, Left: 12}).
		W(w.Width).
		OnTap(func() {
			if !opt.Disabled {
				s.open = false
				s.value = opt.Value
				if w.OnChange != nil {
					w.OnChange(opt.Value)
				}
				s.SetState(nil)
			}
		})
}

func findOptLabel(value string, opts []SelectOption) string {
	for _, o := range opts {
		if o.Value == value {
			if o.Label != "" {
				return o.Label
			}
			return o.Value
		}
	}
	return value
}

// ==================== MultiSelect ====================

// MultiSelectWidget 多选下拉。
type MultiSelectWidget struct {
	ui.BaseWidget
	Options     []SelectOption
	Values      []string
	Placeholder string
	OnChange    func([]string)
	Disabled    bool
	Width       float32
}

// MultiSelect 创建多选下拉。
func MultiSelect(options []SelectOption) MultiSelectWidget {
	return MultiSelectWidget{Options: options, Placeholder: "请选择...", Width: 200}
}

func (m MultiSelectWidget) WithValues(v []string) MultiSelectWidget     { m.Values = v; return m }
func (m MultiSelectWidget) WithPlaceholder(v string) MultiSelectWidget  { m.Placeholder = v; return m }
func (m MultiSelectWidget) WithOnChange(fn func([]string)) MultiSelectWidget { m.OnChange = fn; return m }
func (m MultiSelectWidget) WithDisabled(v bool) MultiSelectWidget       { m.Disabled = v; return m }
func (m MultiSelectWidget) WithWidth(v float32) MultiSelectWidget       { m.Width = v; return m }

func (m MultiSelectWidget) CreateElement() ui.Element { return ui.NewStatefulElement(m) }
func (m MultiSelectWidget) CreateState() ui.State {
	st := &multiSelectState{}
	st.Init(st)
	return st
}

type multiSelectState struct {
	ui.BaseStateOf[MultiSelectWidget]
	open        bool
	dismisserID int
	values      []string
}

func (s *multiSelectState) InitState() {
	s.dismisserID = ui.RegisterPopupDismisser(func() {
		if s.open {
			s.open = false
			s.SetState(nil)
		}
	})
	w := s.Widget()
	s.values = append([]string(nil), w.Values...)
}
func (s *multiSelectState) Dispose() {
	ui.UnregisterPopupDismisser(s.dismisserID)
}
func (s *multiSelectState) DidUpdateWidget(old ui.Widget) {
	oldW := old.(MultiSelectWidget)
	w := s.Widget()
	if !sliceEqual(oldW.Values, w.Values) {
		s.values = append([]string(nil), w.Values...)
	}
}

func (s *multiSelectState) Build(ctx ui.BuildContext) ui.Widget {
	w := s.Widget()
	theme := ui.ThemeOf(ctx)

	display := w.Placeholder
	if len(s.values) > 0 {
		display = "已选 " + itoa(len(s.values)) + " 项"
	}

	trigger := Container(
		Row(
			Text(display).FontSize(theme.FontSizeBase).Color(theme.TextColor),
			Icon(IconArrowDown).FontSize(10).Color(theme.TextMutedColor),
		).Gapf(8).JustifyContent(ui.JustifySpaceBetween).AlignItems(ui.AlignCenter),
	).
		Background(colorToRender(theme.BackgroundColor)).
		Border(colorToRender(theme.BorderColor), 1).
		Radius(theme.BorderRadius).
		Pad(ui.EdgeInsets{Top: 8, Right: 12, Bottom: 8, Left: 12}).
		W(w.Width).
		Focusable(true).
		OnTap(func() {
			if !w.Disabled {
				if !s.open {
					ui.DismissAllPopups()
				}
				s.open = !s.open
				s.SetState(nil)
			}
		})

	opts := make([]ui.Widget, 0, len(w.Options))
	for _, opt := range w.Options {
		opt := opt
		checked := strContains(s.values, opt.Value)
		label := opt.Label
		if label == "" {
			label = opt.Value
		}
		prefix := IconCheckboxEmpty
		if checked {
			prefix = IconCheckboxChecked
		}
		bg := theme.BackgroundColor
		if checked {
			bg = theme.AccentColor
		}
		opts = append(opts, Container(
			Row(
				Icon(prefix).Size(theme.FontSizeBase).Color(theme.TextColor),
				Text(label).FontSize(theme.FontSizeBase).Color(theme.TextColor),
			).Gapf(6).AlignItems(ui.AlignCenter),
		).
			Background(colorToRender(bg)).
			Pad(ui.EdgeInsets{Top: 6, Right: 12, Bottom: 6, Left: 12}).
			W(w.Width).
			OnTap(func() {
				if !opt.Disabled {
					var newValues []string
					if checked {
						newValues = strRemove(s.values, opt.Value)
					} else {
						newValues = append([]string(nil), s.values...)
						newValues = append(newValues, opt.Value)
					}
					s.values = newValues
					if w.OnChange != nil {
						w.OnChange(newValues)
					}
					s.open = true
					s.SetState(nil)
				}
			}))
	}

	list := Scroll(Column(opts...).Gapf(0)).W(w.Width).MaxH(200)
	dropdown := Container(list).
		Background(colorToRender(theme.BackgroundColor)).
		Border(colorToRender(theme.BorderColor), 1).
		Radius(theme.BorderRadius).
		Pad(ui.EdgeInsets{Top: 4, Bottom: 4}).
		W(w.Width)

	if !s.open {
		return Stack(trigger).Z(0)
	}

	return Stack(
		trigger,
		Positioned(dropdown).T(40).L(0).W(w.Width).Z(999),
	).Z(0)
}

func sliceEqual(a, b []string) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}

func strContains(ss []string, t string) bool {
	for _, s := range ss {
		if s == t {
			return true
		}
	}
	return false
}

func strRemove(ss []string, t string) []string {
	r := make([]string, 0, len(ss))
	for _, s := range ss {
		if s != t {
			r = append(r, s)
		}
	}
	return r
}

func itoa(n int) string {
	if n == 0 {
		return "0"
	}
	s := ""
	for n > 0 {
		s = string(rune('0'+n%10)) + s
		n /= 10
	}
	return s
}
