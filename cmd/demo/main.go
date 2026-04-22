package main

import (
	"fmt"
	"image/color"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/sjm1327605995/tenon"
	"github.com/sjm1327605995/tenon/pkg/components"
	"github.com/sjm1327605995/tenon/pkg/fonts"
	"github.com/sjm1327605995/tenon/yoga"
)

// ==================== Counter：基础状态组件 ====================

type Counter struct {
	tenon.BaseWidget
	count int
}

func NewCounter() *Counter {
	c := &Counter{count: 0}
	c.Init(c)
	return c
}

func (c *Counter) Render() tenon.Component {
	return components.NewView().
		SetPadding(yoga.EdgeAll, 20).
		SetBackgroundColor(color.White).
		SetBorderRadius(12).
		SetBorder(yoga.EdgeAll, 1).
		SetBorderColor(color.RGBA{R: 222, G: 226, B: 230, A: 255}).
		SetMargin(yoga.EdgeBottom, 20).
		Add(
			components.NewText(fmt.Sprintf("Count: %d", c.count)).
				SetFontSize(18).
				SetColor(color.RGBA{R: 33, G: 37, B: 41, A: 255}).
				SetMargin(yoga.EdgeBottom, 15),
			components.NewButton("Increment").
				SetWidth(150).
				SetHeight(40).
				SetOnClick(func() {
					c.count++
					c.Invalidate()
				}),
		)
}

// ==================== DataFetcher：生命周期示例 ====================

type DataFetcher struct {
	tenon.BaseWidget
	data string
}

func NewDataFetcher() *DataFetcher {
	d := &DataFetcher{data: "Loading..."}
	d.Init(d)
	return d
}

func (d *DataFetcher) ComponentDidMount() {
	go func() {
		time.Sleep(1 * time.Second)
		d.data = "Data loaded!"
		d.Invalidate()
	}()
}

func (d *DataFetcher) Render() tenon.Component {
	return components.NewView().
		SetPadding(yoga.EdgeAll, 20).
		SetBackgroundColor(color.White).
		SetBorderRadius(12).
		SetBorder(yoga.EdgeAll, 1).
		SetBorderColor(color.RGBA{R: 222, G: 226, B: 230, A: 255}).
		SetMargin(yoga.EdgeBottom, 20).
		Add(
			components.NewText("Lifecycle Demo").SetFontSize(20).SetMargin(yoga.EdgeBottom, 10),
			components.NewText(d.data).SetFontSize(16),
		)
}

// ==================== HooksDemo：Hooks 综合示例 ====================

type HooksDemo struct {
	tenon.BaseWidget
}

func NewHooksDemo() *HooksDemo {
	h := &HooksDemo{}
	h.Init(h)
	return h
}

func (h *HooksDemo) Render() tenon.Component {
	count, setCount := h.UseState(0)
	doubled := h.UseMemo(func() any {
		return count.(int) * 2
	}, []any{count})
	ref := h.UseRef("initial")
	id := h.UseId()

	h.UseEffect(func() func() {
		fmt.Printf("Effect: count = %d\n", count)
		return func() {
			fmt.Printf("Cleanup: count was %d\n", count)
		}
	}, []any{count})

	return components.NewView().
		SetPadding(yoga.EdgeAll, 20).
		SetBackgroundColor(color.White).
		SetBorderRadius(12).
		SetBorder(yoga.EdgeAll, 1).
		SetBorderColor(color.RGBA{R: 222, G: 226, B: 230, A: 255}).
		SetMargin(yoga.EdgeBottom, 20).
		Add(
			components.NewText("Hooks Demo").SetFontSize(20).SetMargin(yoga.EdgeBottom, 10),
			components.NewText(fmt.Sprintf("ID: %s", id)).SetFontSize(14).SetMargin(yoga.EdgeBottom, 5),
			components.NewText(fmt.Sprintf("Count: %d", count)).SetFontSize(16).SetMargin(yoga.EdgeBottom, 5),
			components.NewText(fmt.Sprintf("Doubled: %d", doubled)).SetFontSize(16).SetMargin(yoga.EdgeBottom, 15),
			components.NewButton("+1").SetWidth(100).SetHeight(36).SetOnClick(func() {
				setCount(count.(int) + 1)
			}),
			components.NewButton(fmt.Sprintf("Ref: %v", ref.Current)).SetWidth(150).SetHeight(36).SetMargin(yoga.EdgeTop, 10).SetOnClick(func() {
				ref.Current = fmt.Sprintf("clicked-%d", count)
				h.Invalidate()
			}),
		)
}

// ==================== FormDemo：表单组件示例 ====================

type FormDemo struct {
	tenon.BaseWidget
	progress  float32
	checked   bool
}

func NewFormDemo() *FormDemo {
	f := &FormDemo{progress: 0.3, checked: true}
	f.Init(f)
	return f
}

func (f *FormDemo) Render() tenon.Component {
	return components.NewView().
		SetPadding(yoga.EdgeAll, 20).
		SetBackgroundColor(color.White).
		SetBorderRadius(12).
		SetBorder(yoga.EdgeAll, 1).
		SetBorderColor(color.RGBA{R: 222, G: 226, B: 230, A: 255}).
		SetMargin(yoga.EdgeBottom, 20).
		Add(
			components.NewText("Form Components").SetFontSize(20).SetMargin(yoga.EdgeBottom, 15),
			components.NewText(fmt.Sprintf("Progress: %.0f%%", f.progress*100)).SetFontSize(14).SetMargin(yoga.EdgeBottom, 5),
			components.NewProgressBar().
				SetProgress(f.progress).
				SetWidth(300).
				SetHeight(10).
				SetBorderRadius(5).
				SetMargin(yoga.EdgeBottom, 15),
			components.NewView().SetFlexDirection(yoga.FlexDirectionRow).SetAlignItems(yoga.AlignCenter).Add(
				components.NewCheckbox("Enable Feature").
					SetChecked(f.checked).
					SetOnChange(func(checked bool) {
						f.checked = checked
						f.Invalidate()
					}),
			),
			components.NewButton(f.checkedStateText()).
				SetWidth(180).
				SetHeight(36).
				SetMargin(yoga.EdgeTop, 10).
				SetOnClick(func() {
					f.progress += 0.1
					if f.progress > 1 {
						f.progress = 0
					}
					f.Invalidate()
				}),
		)
}

func (f *FormDemo) checkedStateText() string {
	if f.checked {
		return "Feature Enabled"
	}
	return "Feature Disabled"
}

// ==================== ScrollDemo：滚动视图示例 ====================

type ScrollDemo struct {
	tenon.BaseWidget
	scrollView *components.ScrollView
}

func NewScrollDemo() *ScrollDemo {
	s := &ScrollDemo{}
	s.Init(s)
	return s
}

func (s *ScrollDemo) Render() tenon.Component {
	if s.scrollView == nil {
		s.scrollView = components.NewScrollView().
			SetWidth(400).
			SetHeight(200).
			SetBackgroundColor(color.White).
			SetBorderRadius(12).
			SetBorder(yoga.EdgeAll, 1).
			SetBorderColor(color.RGBA{R: 222, G: 226, B: 230, A: 255}).
			SetMargin(yoga.EdgeBottom, 20)

		content := s.scrollView.Content()
		content.SetFlexDirection(yoga.FlexDirectionColumn).SetPadding(yoga.EdgeAll, 15)
		for i := 1; i <= 20; i++ {
			content.Add(
				components.NewText(fmt.Sprintf("Scroll item #%d", i)).
					SetFontSize(14).
					SetMargin(yoga.EdgeBottom, 10),
			)
		}
	}
	return s.scrollView
}

// ==================== ImageDemo：图片组件示例 ====================

type ImageDemo struct {
	tenon.BaseWidget
}

func NewImageDemo() *ImageDemo {
	im := &ImageDemo{}
	im.Init(im)
	return im
}

func (im *ImageDemo) Render() tenon.Component {
	img := ebiten.NewImage(120, 120)
	img.Fill(color.RGBA{R: 100, G: 180, B: 255, A: 255})

	return components.NewView().
		SetPadding(yoga.EdgeAll, 20).
		SetBackgroundColor(color.White).
		SetBorderRadius(12).
		SetBorder(yoga.EdgeAll, 1).
		SetBorderColor(color.RGBA{R: 222, G: 226, B: 230, A: 255}).
		SetMargin(yoga.EdgeBottom, 20).
		Add(
			components.NewText("Image Component").SetFontSize(20).SetMargin(yoga.EdgeBottom, 15),
			components.NewView().SetFlexDirection(yoga.FlexDirectionRow).SetAlignItems(yoga.AlignCenter).Add(
				components.NewImage().SetEbitenImage(img).SetWidth(120).SetHeight(120).SetMargin(yoga.EdgeRight, 20),
				components.NewView().SetFlexDirection(yoga.FlexDirectionColumn).SetJustifyContent(yoga.JustifyCenter).Add(
					components.NewText("Ebiten Image").SetFontSize(14).SetColor(color.RGBA{R: 100, G: 100, B: 100, A: 255}),
					components.NewText("Rendered in Tenon").SetFontSize(14).SetColor(color.RGBA{R: 100, G: 100, B: 100, A: 255}),
				),
			),
		)
}

// ==================== StyleDemo：样式展示 ====================

type StyleDemo struct {
	tenon.BaseWidget
}

func NewStyleDemo() *StyleDemo {
	s := &StyleDemo{}
	s.Init(s)
	return s
}

func (s *StyleDemo) Render() tenon.Component {
	return components.NewView().
		SetPadding(yoga.EdgeAll, 20).
		SetBackgroundColor(color.White).
		SetBorderRadius(12).
		SetBorder(yoga.EdgeAll, 1).
		SetBorderColor(color.RGBA{R: 222, G: 226, B: 230, A: 255}).
		SetMargin(yoga.EdgeBottom, 20).
		Add(
			components.NewText("Styles & Effects").SetFontSize(20).SetMargin(yoga.EdgeBottom, 15),
			components.NewView().SetFlexDirection(yoga.FlexDirectionRow).SetJustifyContent(yoga.JustifySpaceBetween).Add(
				components.NewView().
					SetWidth(80).SetHeight(80).
					SetBackgroundColor(color.RGBA{R: 255, G: 99, B: 71, A: 255}).
					SetBorderRadius(40).
					SetMargin(yoga.EdgeRight, 10),
				components.NewView().
					SetWidth(80).SetHeight(80).
					SetBackgroundColor(color.RGBA{R: 50, G: 205, B: 50, A: 255}).
					SetBorderRadius(8).
					SetShadow(color.RGBA{A: 80}, 12, 4, 4).
					SetMargin(yoga.EdgeRight, 10),
				components.NewView().
					SetWidth(80).SetHeight(80).
					SetBackgroundColor(color.RGBA{R: 255, G: 215, B: 0, A: 255}).
					SetBorderRadius(0).
					SetBorder(yoga.EdgeAll, 3).
					SetBorderColor(color.RGBA{R: 255, G: 140, B: 0, A: 255}).
					SetMargin(yoga.EdgeRight, 10),
				components.NewView().
					SetWidth(80).SetHeight(80).
					SetBackgroundColor(color.RGBA{R: 138, G: 43, B: 226, A: 255}).
					SetBorderRadius4(20, 4, 20, 4).
					SetShadow(color.RGBA{A: 60}, 8, 0, 6),
			),
		)
}

// ==================== FocusDemo：焦点系统示例 ====================

type FocusDemo struct {
	tenon.BaseWidget
}

func NewFocusDemo() *FocusDemo {
	f := &FocusDemo{}
	f.Init(f)
	return f
}

func (f *FocusDemo) Render() tenon.Component {
	return components.NewView().
		SetPadding(yoga.EdgeAll, 20).
		SetBackgroundColor(color.White).
		SetBorderRadius(12).
		SetBorder(yoga.EdgeAll, 1).
		SetBorderColor(color.RGBA{R: 222, G: 226, B: 230, A: 255}).
		SetMargin(yoga.EdgeBottom, 20).
		Add(
			components.NewText("Focus System (Press Tab)").SetFontSize(20).SetMargin(yoga.EdgeBottom, 15),
			components.NewView().SetFlexDirection(yoga.FlexDirectionRow).SetAlignItems(yoga.AlignCenter).Add(
				components.NewButton("Button A").SetWidth(100).SetHeight(36).SetMargin(yoga.EdgeRight, 10),
				components.NewButton("Button B").SetWidth(100).SetHeight(36).SetMargin(yoga.EdgeRight, 10),
				components.NewButton("Button C").SetWidth(100).SetHeight(36).SetMargin(yoga.EdgeRight, 20),
				components.NewCheckbox("Check 1").SetMargin(yoga.EdgeRight, 15),
				components.NewCheckbox("Check 2"),
			),
		)
}

// ==================== LayoutDemo：Flex 布局示例 ====================

type LayoutDemo struct {
	tenon.BaseWidget
}

func NewLayoutDemo() *LayoutDemo {
	l := &LayoutDemo{}
	l.Init(l)
	return l
}

func (l *LayoutDemo) Render() tenon.Component {
	row1 := components.NewView().SetFlexDirection(yoga.FlexDirectionRow).SetMargin(yoga.EdgeBottom, 10)
	for _, c := range []color.Color{
		color.RGBA{R: 255, G: 100, B: 100, A: 255},
		color.RGBA{R: 100, G: 255, B: 100, A: 255},
		color.RGBA{R: 100, G: 100, B: 255, A: 255},
	} {
		row1.Add(components.NewView().SetWidth(60).SetHeight(30).SetBackgroundColor(c).SetMargin(yoga.EdgeRight, 8))
	}

	row2 := components.NewView().SetFlexDirection(yoga.FlexDirectionRow).SetMargin(yoga.EdgeBottom, 10)
	row2.Add(components.NewView().SetWidth(40).SetHeight(30).SetBackgroundColor(color.RGBA{R: 255, G: 180, B: 100, A: 255}).SetFlexGrow(1).SetMargin(yoga.EdgeRight, 8))
	row2.Add(components.NewView().SetWidth(40).SetHeight(30).SetBackgroundColor(color.RGBA{R: 100, G: 180, B: 255, A: 255}).SetFlexGrow(2).SetMargin(yoga.EdgeRight, 8))
	row2.Add(components.NewView().SetWidth(40).SetHeight(30).SetBackgroundColor(color.RGBA{R: 180, G: 100, B: 255, A: 255}).SetFlexGrow(1))

	row3 := components.NewView().SetFlexDirection(yoga.FlexDirectionRow).SetFlexWrap(yoga.WrapWrap)
	for i := 0; i < 6; i++ {
		row3.Add(components.NewView().SetWidth(80).SetHeight(30).
			SetBackgroundColor(color.RGBA{R: uint8(100 + i*25), G: 150, B: 200, A: 255}).
			SetMargin(yoga.EdgeRight, 8).SetMargin(yoga.EdgeBottom, 8))
	}

	return components.NewView().
		SetPadding(yoga.EdgeAll, 20).
		SetBackgroundColor(color.White).
		SetBorderRadius(12).
		SetBorder(yoga.EdgeAll, 1).
		SetBorderColor(color.RGBA{R: 222, G: 226, B: 230, A: 255}).
		SetMargin(yoga.EdgeBottom, 20).
		Add(
			components.NewText("Flex Layout").SetFontSize(20).SetMargin(yoga.EdgeBottom, 15),
			components.NewText("Row (fixed size)").SetFontSize(12).SetColor(color.RGBA{R: 120, G: 120, B: 120, A: 255}).SetMargin(yoga.EdgeBottom, 5),
			row1,
			components.NewText("Row (flex grow 1:2:1)").SetFontSize(12).SetColor(color.RGBA{R: 120, G: 120, B: 120, A: 255}).SetMargin(yoga.EdgeBottom, 5),
			row2,
			components.NewText("Row (wrap)").SetFontSize(12).SetColor(color.RGBA{R: 120, G: 120, B: 120, A: 255}).SetMargin(yoga.EdgeBottom, 5),
			row3,
		)
}

// ==================== InputDemo：文本输入示例 ====================

type InputDemo struct {
	tenon.BaseWidget
	nameInput  *components.TextInput
	emailInput *components.TextInput
	name       string
	email      string
}

func NewInputDemo() *InputDemo {
	i := &InputDemo{name: "", email: ""}
	i.Init(i)
	return i
}

func (i *InputDemo) Render() tenon.Component {
	// 缓存并复用 TextInput，避免 Invalidate 后重建导致输入状态丢失
	if i.nameInput == nil {
		i.nameInput = components.NewTextInput().
			SetWidth(300).
			SetPlaceholder("Enter your name").
			SetOnChange(func(v string) {
				i.name = v
				i.Invalidate()
			})
	}
	if i.emailInput == nil {
		i.emailInput = components.NewTextInput().
			SetWidth(300).
			SetPlaceholder("Enter your email").
			SetOnChange(func(v string) {
				i.email = v
				i.Invalidate()
			})
	}
	return components.NewView().
		SetPadding(yoga.EdgeAll, 20).
		SetBackgroundColor(color.White).
		SetBorderRadius(12).
		SetBorder(yoga.EdgeAll, 1).
		SetBorderColor(color.RGBA{R: 222, G: 226, B: 230, A: 255}).
		SetMargin(yoga.EdgeBottom, 20).
		Add(
			components.NewText("Text Input").SetFontSize(20).SetMargin(yoga.EdgeBottom, 15),
			components.NewText("Name:").SetFontSize(14).SetColor(color.RGBA{R: 100, G: 100, B: 100, A: 255}).SetMargin(yoga.EdgeBottom, 4),
			i.nameInput.SetMargin(yoga.EdgeBottom, 12),
			components.NewText("Email:").SetFontSize(14).SetColor(color.RGBA{R: 100, G: 100, B: 100, A: 255}).SetMargin(yoga.EdgeBottom, 4),
			i.emailInput.SetMargin(yoga.EdgeBottom, 12),
			components.NewText(fmt.Sprintf("Name: %s | Email: %s", i.name, i.email)).SetFontSize(14).SetColor(color.RGBA{R: 100, G: 100, B: 100, A: 255}),
		)
}

// ==================== ControlDemo：控件组件示例 ====================

type ControlDemo struct {
	tenon.BaseWidget
	sliderVal float32
	switchOn  bool
	radioSel  int
}

func NewControlDemo() *ControlDemo {
	c := &ControlDemo{sliderVal: 50, switchOn: false, radioSel: 0}
	c.Init(c)
	return c
}

func (c *ControlDemo) Render() tenon.Component {
	return components.NewView().
		SetPadding(yoga.EdgeAll, 20).
		SetBackgroundColor(color.White).
		SetBorderRadius(12).
		SetBorder(yoga.EdgeAll, 1).
		SetBorderColor(color.RGBA{R: 222, G: 226, B: 230, A: 255}).
		SetMargin(yoga.EdgeBottom, 20).
		Add(
			components.NewText("Controls").SetFontSize(20).SetMargin(yoga.EdgeBottom, 15),
			components.NewView().SetFlexDirection(yoga.FlexDirectionRow).SetAlignItems(yoga.AlignCenter).SetMargin(yoga.EdgeBottom, 12).Add(
				components.NewText(fmt.Sprintf("Slider: %.0f", c.sliderVal)).SetFontSize(14).SetWidth(80),
				components.NewSlider(0, 100).SetValue(c.sliderVal).SetWidth(200).SetOnChange(func(v float32) {
					c.sliderVal = v
					c.Invalidate()
				}),
			),
			components.NewView().SetFlexDirection(yoga.FlexDirectionRow).SetAlignItems(yoga.AlignCenter).SetMargin(yoga.EdgeBottom, 12).Add(
				components.NewText("Switch:").SetFontSize(14).SetMargin(yoga.EdgeRight, 10),
				components.NewSwitch().SetChecked(c.switchOn).SetOnChange(func(on bool) {
					c.switchOn = on
					c.Invalidate()
				}),
				components.NewText(fmt.Sprintf("  %v", c.switchOn)).SetFontSize(14),
			),
			components.NewView().SetFlexDirection(yoga.FlexDirectionRow).SetAlignItems(yoga.AlignCenter).Add(
				components.NewRadio("Option A").SetSelected(c.radioSel == 0).SetOnChange(func(_ bool) {
					c.radioSel = 0
					c.Invalidate()
				}).SetMargin(yoga.EdgeRight, 15),
				components.NewRadio("Option B").SetSelected(c.radioSel == 1).SetOnChange(func(_ bool) {
					c.radioSel = 1
					c.Invalidate()
				}).SetMargin(yoga.EdgeRight, 15),
				components.NewRadio("Option C").SetSelected(c.radioSel == 2).SetOnChange(func(_ bool) {
					c.radioSel = 2
					c.Invalidate()
				}),
			),
		)
}

// ==================== TextWrapDemo：文本换行示例 ====================

type TextWrapDemo struct {
	tenon.BaseWidget
}

func NewTextWrapDemo() *TextWrapDemo {
	t := &TextWrapDemo{}
	t.Init(t)
	return t
}

func (t *TextWrapDemo) Render() tenon.Component {
	longText := "The quick brown fox jumps over the lazy dog. 这是一段用于测试文本换行功能的中英文混合长文本内容。"
	cjkText := "这是一个很长的中文文本用于测试自动换行功能，中文应该在字符边界处正确断行。"

	return components.NewView().
		SetPadding(yoga.EdgeAll, 20).
		SetBackgroundColor(color.White).
		SetBorderRadius(12).
		SetBorder(yoga.EdgeAll, 1).
		SetBorderColor(color.RGBA{R: 222, G: 226, B: 230, A: 255}).
		SetMargin(yoga.EdgeBottom, 20).
		Add(
			components.NewText("Text Wrapping").SetFontSize(20).SetMargin(yoga.EdgeBottom, 15),
			components.NewText("Normal (auto wrap):").SetFontSize(12).SetColor(color.RGBA{R: 120, G: 120, B: 120, A: 255}).SetMargin(yoga.EdgeBottom, 4),
			components.NewText(longText).SetFontSize(14).SetWidth(350).SetWhiteSpace(components.WhiteSpaceNormal).SetMargin(yoga.EdgeBottom, 12),
			components.NewText("BreakAll:").SetFontSize(12).SetColor(color.RGBA{R: 120, G: 120, B: 120, A: 255}).SetMargin(yoga.EdgeBottom, 4),
			components.NewText(longText).SetFontSize(14).SetWidth(350).SetWhiteSpace(components.WhiteSpaceNormal).SetWordBreak(components.WordBreakBreakAll).SetMargin(yoga.EdgeBottom, 12),
			components.NewText("CJK Normal:").SetFontSize(12).SetColor(color.RGBA{R: 120, G: 120, B: 120, A: 255}).SetMargin(yoga.EdgeBottom, 4),
			components.NewText(cjkText).SetFontSize(14).SetWidth(300).SetWhiteSpace(components.WhiteSpaceNormal).SetMargin(yoga.EdgeBottom, 12),
			components.NewText("Pre-wrap (preserve \\n):").SetFontSize(12).SetColor(color.RGBA{R: 120, G: 120, B: 120, A: 255}).SetMargin(yoga.EdgeBottom, 4),
			components.NewText("Line 1\nLine 2\nLine 3").SetFontSize(14).SetWidth(300).SetWhiteSpace(components.WhiteSpacePreWrap),
		)
}

// ==================== App：根组件 ====================

type App struct {
	tenon.BaseWidget
	scrollView *components.ScrollView
}

func NewApp() *App {
	a := &App{}
	a.Init(a)
	return a
}

func (a *App) Render() tenon.Component {
	if a.scrollView == nil {
		a.scrollView = components.NewScrollView().
			SetWidth(800).
			SetHeight(600).
			SetBackgroundColor(color.RGBA{R: 240, G: 240, B: 240, A: 255})

		content := a.scrollView.Content()
		content.SetFlexDirection(yoga.FlexDirectionColumn).SetPadding(yoga.EdgeAll, 20)
		content.Add(
			components.NewText("Tenon UI Framework").
				SetFontSize(24).
				SetColor(color.RGBA{R: 51, G: 51, B: 51, A: 255}).
				SetMargin(yoga.EdgeBottom, 30),
			NewCounter(),
			NewDataFetcher(),
			NewHooksDemo(),
			NewFormDemo(),
			NewImageDemo(),
			NewStyleDemo(),
			NewFocusDemo(),
			NewLayoutDemo(),
			NewInputDemo(),
			NewControlDemo(),
			NewTextWrapDemo(),
			NewScrollDemo(),
		)
	}
	return a.scrollView
}

// ==================== Main ====================

func main() {
	if err := fonts.InitDefaultFont(); err != nil {
		panic("Failed to init default font: " + err.Error())
	}
	if err := fonts.LoadFontFromFile(fonts.FontFamilySans, "font/OPPOSans-Medium.ttf"); err == nil {
		fonts.SetDefaultFontFamily(fonts.FontFamilySans)
	}

	// 直接运行 App，App 内部是 ScrollView
	tenon.Run(NewApp(), 800, 600)
}
