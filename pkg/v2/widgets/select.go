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

func (s SelectWidget) WithValue(v string) SelectWidget           { s.Value = v; return s }
func (s SelectWidget) WithPlaceholder(v string) SelectWidget     { s.Placeholder = v; return s }
func (s SelectWidget) WithOnChange(fn func(string)) SelectWidget { s.OnChange = fn; return s }
func (s SelectWidget) WithDisabled(v bool) SelectWidget          { s.Disabled = v; return s }
func (s SelectWidget) WithWidth(v float32) SelectWidget          { s.Width = v; return s }

func (s SelectWidget) CreateElement() ui.Element { return ui.NewStatefulElement(s) }
func (s SelectWidget) CreateState() ui.State {
	st := &selectState{}
	st.Init(st)
	return st
}

type selectState struct {
	ui.BaseState
	open       bool
	dismisserID int
}

func (s *selectState) InitState() {
	s.dismisserID = ui.RegisterPopupDismisser(func() {
		if s.open {
			s.open = false
			s.SetState(nil)
		}
	})
}
func (s *selectState) Dispose() {
	ui.UnregisterPopupDismisser(s.dismisserID)
}
func (s *selectState) DidUpdateWidget(old ui.Widget) {}

func (s *selectState) Build(ctx ui.BuildContext) ui.Widget {
	w := s.GetWidget().(SelectWidget)
	theme := ui.ThemeOf(ctx)
	trigger := s.buildTrigger(w, theme)
	if !s.open {
		return trigger
	}
	// 使用 Stack + Positioned 让 dropdown 悬浮显示，不占用布局空间
	return Stack(
		trigger,
		Positioned(s.buildDropdown(w, theme)).T(40).L(0).W(w.Width).Z(100),
	).Z(100)
}

func (s *selectState) buildTrigger(w SelectWidget, theme *ui.Theme) ui.Widget {
	text := w.Placeholder
	clr := theme.PlaceholderColor
	if w.Value != "" {
		text = findOptLabel(w.Value, w.Options)
		clr = theme.TextColor
	}
	arrow := "▼"
	if s.open {
		arrow = "▲"
	}
	return Container(
		Row(
			Text(text).FontSize(theme.FontSizeBase).Color(clr),
			Text(arrow).FontSize(10).Color(theme.TextMutedColor),
		).Gapf(8).JustifyContent(ui.JustifySpaceBetween).AlignItems(ui.AlignCenter),
	).
		Background(colorToRender(theme.BackgroundColor)).
		Border(colorToRender(theme.BorderColor), 1).
		Radius(theme.BorderRadius).
		Pad(ui.EdgeInsets{Top: 8, Right: 12, Bottom: 8, Left: 12}).
		W(w.Width).
		H(40).
		OnTap(func() {
			if !w.Disabled {
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
	return Container(Column(opts...).Gapf(0)).
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
	if opt.Value == w.Value {
		bg = theme.AccentColor
	}
	return Container(Text(label).FontSize(theme.FontSizeBase).Color(txtClr)).
		Background(colorToRender(bg)).
		Pad(ui.EdgeInsets{Top: 6, Right: 12, Bottom: 6, Left: 12}).
		W(w.Width).
		OnTap(func() {
			if !opt.Disabled {
				s.open = false
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
	ui.BaseState
	open       bool
	dismisserID int
}

func (s *multiSelectState) InitState() {
	s.dismisserID = ui.RegisterPopupDismisser(func() {
		if s.open {
			s.open = false
			s.SetState(nil)
		}
	})
}
func (s *multiSelectState) Dispose() {
	ui.UnregisterPopupDismisser(s.dismisserID)
}
func (s *multiSelectState) DidUpdateWidget(old ui.Widget) {}

func (s *multiSelectState) Build(ctx ui.BuildContext) ui.Widget {
	w := s.GetWidget().(MultiSelectWidget)
	theme := ui.ThemeOf(ctx)

	display := w.Placeholder
	if len(w.Values) > 0 {
		display = "已选 " + itoa(len(w.Values)) + " 项"
	}

	trigger := Container(
		Row(
			Text(display).FontSize(theme.FontSizeBase).Color(theme.TextColor),
			Text("▼").FontSize(10).Color(theme.TextMutedColor),
		).Gapf(8).JustifyContent(ui.JustifySpaceBetween).AlignItems(ui.AlignCenter),
	).
		Background(colorToRender(theme.BackgroundColor)).
		Border(colorToRender(theme.BorderColor), 1).
		Radius(theme.BorderRadius).
		Pad(ui.EdgeInsets{Top: 8, Right: 12, Bottom: 8, Left: 12}).
		W(w.Width).
		OnTap(func() {
			if !w.Disabled {
				s.open = !s.open
				s.SetState(nil)
			}
		})

	if !s.open {
		return trigger
	}

	opts := make([]ui.Widget, 0, len(w.Options))
	for _, opt := range w.Options {
		opt := opt
		checked := strContains(w.Values, opt.Value)
		label := opt.Label
		if label == "" {
			label = opt.Value
		}
		prefix := "☐"
		if checked {
			prefix = "☑"
		}
		opts = append(opts, Container(
			Text(prefix+" "+label).FontSize(theme.FontSizeBase).Color(theme.TextColor),
		).
			Background(colorToRender(theme.BackgroundColor)).
			Pad(ui.EdgeInsets{Top: 6, Right: 12, Bottom: 6, Left: 12}).
			W(w.Width).
			OnTap(func() {
				if !opt.Disabled {
					if checked {
						w.Values = strRemove(w.Values, opt.Value)
					} else {
						w.Values = append(w.Values, opt.Value)
					}
					if w.OnChange != nil {
						w.OnChange(w.Values)
					}
					s.SetState(nil)
				}
			}))
	}

	return Column(trigger, Container(Column(opts...).Gapf(0)).
		Background(colorToRender(theme.BackgroundColor)).
		Border(colorToRender(theme.BorderColor), 1).
		Radius(theme.BorderRadius).
		Pad(ui.EdgeInsets{Top: 4, Bottom: 4}).
		W(w.Width)).Gapf(2)
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
