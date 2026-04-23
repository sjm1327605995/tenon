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

// ==================== 页面组件 ====================

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
	return newPageCard("Counter",
		components.NewText(fmt.Sprintf("Count: %d", c.count)).
			SetFontSize(tenon.GetTheme().FontSizeBase+2).
			SetMargin(yoga.EdgeBottom, 16),
		components.NewButton("Increment").SetWidth(140).SetHeight(40).SetOnClick(func() {
			c.SetState(func() { c.count++ })
		}),
	)
}

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
			d.SetState(func() { d.data = "Data loaded!" })
		}()
}

func (d *DataFetcher) Render() tenon.Component {
	return newPageCard("Lifecycle",
		components.NewText(d.data).SetFontSize(tenon.GetTheme().FontSizeBase + 2),
	)
}

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

	return newPageCard("Hooks",
		components.NewText(fmt.Sprintf("ID: %s", id)).SetFontSize(tenon.GetTheme().FontSizeBase).SetMargin(yoga.EdgeBottom, 4),
		components.NewText(fmt.Sprintf("Count: %d", count)).SetFontSize(tenon.GetTheme().FontSizeBase).SetMargin(yoga.EdgeBottom, 4),
		components.NewText(fmt.Sprintf("Doubled: %d", doubled)).SetFontSize(tenon.GetTheme().FontSizeBase).SetMargin(yoga.EdgeBottom, 16),
		components.NewView().SetFlexDirection(yoga.FlexDirectionRow).Add(
			components.NewButton("+1").SetWidth(80).SetHeight(36).SetMargin(yoga.EdgeRight, 8).SetOnClick(func() {
				setCount(count.(int) + 1)
			}),
			components.NewButton(fmt.Sprintf("Ref: %v", ref.Current)).SetWidth(140).SetHeight(36).SetOnClick(func() {
				ref.Current = fmt.Sprintf("clicked-%d", count)
				h.SetState(func() {})
			}),
		),
	)
}

type FormDemo struct {
	tenon.BaseWidget
	progress float32
	checked  bool
}

func NewFormDemo() *FormDemo {
	f := &FormDemo{progress: 0.3, checked: true}
	f.Init(f)
	return f
}

func (f *FormDemo) Render() tenon.Component {
	btnLabel := "Feature Disabled"
	if f.checked {
		btnLabel = "Feature Enabled"
	}
	return newPageCard("Form",
		components.NewText(fmt.Sprintf("Progress: %.0f%%", f.progress*100)).
			SetFontSize(tenon.GetTheme().FontSizeBase).
			SetMargin(yoga.EdgeBottom, 8),
		components.NewProgressBar().
			SetProgress(f.progress).
			SetWidthPercent(100).
			SetHeight(10).
			SetBorderRadius(5).
			SetMargin(yoga.EdgeBottom, 16),
		components.NewView().SetFlexDirection(yoga.FlexDirectionRow).SetAlignItems(yoga.AlignCenter).SetMargin(yoga.EdgeBottom, 12).Add(
			components.NewCheckbox("Enable Feature").
				SetChecked(f.checked).
				SetOnChange(func(checked bool) {
					f.SetState(func() { f.checked = checked })
				}),
		),
		components.NewButton(btnLabel).SetWidth(160).SetHeight(36).SetOnClick(func() {
			f.SetState(func() {
				f.progress += 0.1
				if f.progress > 1 {
					f.progress = 0
				}
			})
		}),
	)
}

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
	return newPageCard("Controls",
		components.NewView().SetFlexDirection(yoga.FlexDirectionRow).SetAlignItems(yoga.AlignCenter).SetMargin(yoga.EdgeBottom, 12).Add(
			components.NewText(fmt.Sprintf("Slider: %.0f", c.sliderVal)).SetFontSize(tenon.GetTheme().FontSizeBase).SetWidth(70),
			components.NewSlider(0, 100).SetValue(c.sliderVal).SetWidth(200).SetOnChange(func(v float32) {
				c.SetState(func() { c.sliderVal = v })
			}),
		),
		components.NewView().SetFlexDirection(yoga.FlexDirectionRow).SetAlignItems(yoga.AlignCenter).SetMargin(yoga.EdgeBottom, 12).Add(
			components.NewText("Switch:").SetFontSize(tenon.GetTheme().FontSizeBase).SetMargin(yoga.EdgeRight, 10),
			components.NewSwitch().SetChecked(c.switchOn).SetOnChange(func(on bool) {
				c.SetState(func() { c.switchOn = on })
			}),
			components.NewText(fmt.Sprintf("  %v", c.switchOn)).SetFontSize(tenon.GetTheme().FontSizeBase),
		),
		components.NewView().SetFlexDirection(yoga.FlexDirectionRow).SetAlignItems(yoga.AlignCenter).Add(
			components.NewRadio("Option A").SetSelected(c.radioSel == 0).SetOnChange(func(_ bool) {
				c.SetState(func() { c.radioSel = 0 })
			}).SetMargin(yoga.EdgeRight, 12),
			components.NewRadio("Option B").SetSelected(c.radioSel == 1).SetOnChange(func(_ bool) {
				c.SetState(func() { c.radioSel = 1 })
			}).SetMargin(yoga.EdgeRight, 12),
			components.NewRadio("Option C").SetSelected(c.radioSel == 2).SetOnChange(func(_ bool) {
				c.SetState(func() { c.radioSel = 2 })
			}),
		),
	)
}

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
	if i.nameInput == nil {
		i.nameInput = components.NewTextInput().
			SetWidth(300).
			SetPlaceholder("Enter your name").
			SetOnChange(func(v string) {
				i.SetState(func() { i.name = v })
			})
	}
	if i.emailInput == nil {
		i.emailInput = components.NewTextInput().
			SetWidth(300).
			SetPlaceholder("Enter your email").
			SetOnChange(func(v string) {
				i.SetState(func() { i.email = v })
			})
	}
	return newPageCard("Text Input",
		components.NewText("Name:").SetFontSize(tenon.GetTheme().FontSizeBase).SetColor(tenon.GetTheme().TextMutedColor).SetMargin(yoga.EdgeBottom, 4),
		i.nameInput.SetMargin(yoga.EdgeBottom, 12),
		components.NewText("Email:").SetFontSize(tenon.GetTheme().FontSizeBase).SetColor(tenon.GetTheme().TextMutedColor).SetMargin(yoga.EdgeBottom, 4),
		i.emailInput.SetMargin(yoga.EdgeBottom, 12),
		components.NewText(fmt.Sprintf("Name: %s | Email: %s", i.name, i.email)).
			SetFontSize(tenon.GetTheme().FontSizeBase).
			SetColor(tenon.GetTheme().TextMutedColor),
	)
}

type FocusDemo struct {
	tenon.BaseWidget
}

func NewFocusDemo() *FocusDemo {
	f := &FocusDemo{}
	f.Init(f)
	return f
}

func (f *FocusDemo) Render() tenon.Component {
	return newPageCard("Focus System (Press Tab)",
		components.NewView().SetFlexDirection(yoga.FlexDirectionRow).SetAlignItems(yoga.AlignCenter).SetFlexWrap(yoga.WrapWrap).Add(
			components.NewButton("Button A").SetWidth(90).SetHeight(36).SetMargin(yoga.EdgeRight, 8).SetMargin(yoga.EdgeBottom, 8),
			components.NewButton("Button B").SetWidth(90).SetHeight(36).SetMargin(yoga.EdgeRight, 8).SetMargin(yoga.EdgeBottom, 8),
			components.NewButton("Button C").SetWidth(90).SetHeight(36).SetMargin(yoga.EdgeRight, 16).SetMargin(yoga.EdgeBottom, 8),
			components.NewCheckbox("Check 1").SetMargin(yoga.EdgeRight, 12).SetMargin(yoga.EdgeBottom, 8),
			components.NewCheckbox("Check 2").SetMargin(yoga.EdgeBottom, 8),
		),
	)
}

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

	return newPageCard("Flex Layout",
		newSectionTitle("Row (fixed size)"),
		row1,
		newSectionTitle("Row (flex grow 1:2:1)"),
		row2,
		newSectionTitle("Row (wrap)"),
		row3,
	)
}

type StyleDemo struct {
	tenon.BaseWidget
}

func NewStyleDemo() *StyleDemo {
	s := &StyleDemo{}
	s.Init(s)
	return s
}

func (s *StyleDemo) Render() tenon.Component {
	return newPageCard("Styles & Effects",
		components.NewView().SetFlexDirection(yoga.FlexDirectionRow).SetJustifyContent(yoga.JustifySpaceBetween).Add(
			components.NewView().SetWidth(70).SetHeight(70).SetBackgroundColor(tenon.GetTheme().PrimaryColor).SetBorderRadius(35).SetMargin(yoga.EdgeRight, 10),
			components.NewView().SetWidth(70).SetHeight(70).SetBackgroundColor(color.RGBA{R: 50, G: 205, B: 50, A: 255}).SetBorderRadius(tenon.GetTheme().BorderRadius).SetShadow(tenon.GetTheme().ShadowColor, 12, 4, 4).SetMargin(yoga.EdgeRight, 10),
			components.NewView().SetWidth(70).SetHeight(70).SetBackgroundColor(color.RGBA{R: 255, G: 215, B: 0, A: 255}).SetBorderRadius(0).SetBorder(yoga.EdgeAll, 3).SetBorderColor(tenon.GetTheme().BorderColor).SetMargin(yoga.EdgeRight, 10),
			components.NewView().SetWidth(70).SetHeight(70).SetBackgroundColor(color.RGBA{R: 138, G: 43, B: 226, A: 255}).SetBorderRadius4(16, 4, 16, 4).SetShadow(tenon.GetTheme().ShadowColor, 8, 0, 6),
		),
	)
}

type ImageDemo struct {
	tenon.BaseWidget
}

func NewImageDemo() *ImageDemo {
	im := &ImageDemo{}
	im.Init(im)
	return im
}

func (im *ImageDemo) Render() tenon.Component {
	img := ebiten.NewImage(100, 100)
	img.Fill(tenon.GetTheme().PrimaryColor)

	return newPageCard("Image",
		components.NewView().SetFlexDirection(yoga.FlexDirectionRow).SetAlignItems(yoga.AlignCenter).Add(
			components.NewImage().SetEbitenImage(img).SetWidth(100).SetHeight(100).SetMargin(yoga.EdgeRight, 16),
			components.NewView().SetFlexDirection(yoga.FlexDirectionColumn).SetJustifyContent(yoga.JustifyCenter).Add(
				components.NewText("Ebiten Image").SetFontSize(tenon.GetTheme().FontSizeBase),
				components.NewText("Rendered in Tenon").SetFontSize(tenon.GetTheme().FontSizeBase).SetColor(tenon.GetTheme().TextMutedColor),
			),
		),
	)
}

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

	return newPageCard("Text Wrapping",
		newSectionTitle("Normal (auto wrap)"),
		components.NewText(longText).SetFontSize(tenon.GetTheme().FontSizeBase).SetWidth(300).SetWhiteSpace(components.WhiteSpaceNormal).SetMargin(yoga.EdgeBottom, 10),
		newSectionTitle("BreakAll"),
		components.NewText(longText).SetFontSize(tenon.GetTheme().FontSizeBase).SetWidth(300).SetWhiteSpace(components.WhiteSpaceNormal).SetWordBreak(components.WordBreakBreakAll).SetMargin(yoga.EdgeBottom, 10),
		newSectionTitle("CJK Normal"),
		components.NewText(cjkText).SetFontSize(tenon.GetTheme().FontSizeBase).SetWidth(260).SetWhiteSpace(components.WhiteSpaceNormal).SetMargin(yoga.EdgeBottom, 10),
		newSectionTitle("Pre-wrap (preserve \\n)"),
		components.NewText("Line 1\nLine 2\nLine 3").SetFontSize(tenon.GetTheme().FontSizeBase).SetWidth(260).SetWhiteSpace(components.WhiteSpacePreWrap),
	)
}

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
			SetHeight(200).
			SetBackgroundColor(tenon.GetTheme().SurfaceColor).
			SetBorderRadius(tenon.GetTheme().BorderRadius).
			SetBorder(yoga.EdgeAll, 1).
			SetBorderColor(tenon.GetTheme().BorderColor)

		content := s.scrollView.Content()
		content.SetFlexDirection(yoga.FlexDirectionColumn).SetPadding(yoga.EdgeAll, 12)
		for i := 1; i <= 20; i++ {
			content.Add(
				components.NewText(fmt.Sprintf("Scroll item #%d", i)).
					SetFontSize(tenon.GetTheme().FontSizeBase).
					SetMargin(yoga.EdgeBottom, 8),
			)
		}
	}
	return s.scrollView
}

// ==================== 布局辅助函数 ====================

func newPageCard(title string, children ...tenon.Component) tenon.Component {
	theme := tenon.GetTheme()
	card := components.NewView().
		SetPadding(yoga.EdgeAll, 24).
		SetBackgroundColor(theme.SurfaceColor).
		SetBorderRadius(theme.BorderRadius).
		SetBorder(yoga.EdgeAll, 1).
		SetBorderColor(theme.BorderColor).
		SetMargin(yoga.EdgeBottom, 16)

	card.Add(
		components.NewText(title).
			SetFontSize(theme.FontSizeLG).
			SetColor(theme.TextColor).
			SetMargin(yoga.EdgeBottom, 16),
	)
	for _, c := range children {
		card.AddChild(c)
	}
	return card
}

func newSectionTitle(text string) tenon.Component {
	return components.NewText(text).
		SetFontSize(tenon.GetTheme().FontSizeSM).
		SetColor(tenon.GetTheme().TextMutedColor).
		SetMargin(yoga.EdgeBottom, 4)
}

// ==================== App：根组件 ====================

type App struct {
	tenon.BaseWidget
	currentPage string
}

func NewApp() *App {
	a := &App{currentPage: "counter"}
	a.Init(a)
	return a
}

func (a *App) Render() tenon.Component {
	menu := components.NewMenu().SetItems([]components.MenuItemData{
		{Key: "counter", Label: "Counter"},
		{Key: "lifecycle", Label: "Lifecycle"},
		{Key: "hooks", Label: "Hooks"},
		{Key: "form", Label: "Form"},
		{Key: "controls", Label: "Controls"},
		{Key: "input", Label: "Text Input"},
		{Key: "focus", Label: "Focus"},
		{Key: "layout", Label: "Layout"},
		{Key: "styles", Label: "Styles"},
		{Key: "image", Label: "Image"},
		{Key: "text", Label: "Text"},
		{Key: "scroll", Label: "Scroll"},
	}).SetSelectedKey(a.currentPage).SetOnSelect(func(key string) {
		a.SetState(func() { a.currentPage = key })
	})

	content := components.NewView().
		SetFlexGrow(1).
		SetPadding(yoga.EdgeAll, 24).
		SetBackgroundColor(tenon.GetTheme().BackgroundColor)

	switch a.currentPage {
	case "counter":
		content.AddChild(NewCounter())
	case "lifecycle":
		content.AddChild(NewDataFetcher())
	case "hooks":
		content.AddChild(NewHooksDemo())
	case "form":
		content.AddChild(NewFormDemo())
	case "controls":
		content.AddChild(NewControlDemo())
	case "input":
		content.AddChild(NewInputDemo())
	case "focus":
		content.AddChild(NewFocusDemo())
	case "layout":
		content.AddChild(NewLayoutDemo())
	case "styles":
		content.AddChild(NewStyleDemo())
	case "image":
		content.AddChild(NewImageDemo())
	case "text":
		content.AddChild(NewTextWrapDemo())
	case "scroll":
		content.AddChild(NewScrollDemo())
	}

	scrollView := components.NewScrollView().
		SetWidthPercent(100).
		SetHeightPercent(100)
	scrollView.Content().AddChild(content)

	return components.NewView().
		SetFlexDirection(yoga.FlexDirectionRow).
		SetWidthPercent(100).
		SetHeightPercent(100).
		Add(menu, scrollView)
}

// ==================== Main ====================

func main() {
	if err := fonts.InitDefaultFont(); err != nil {
		panic("Failed to init default font: " + err.Error())
	}
	if err := fonts.LoadFontFromFile(fonts.FontFamilySans, "font/OPPOSans-Medium.ttf"); err == nil {
		fonts.SetDefaultFontFamily(fonts.FontFamilySans)
	}

	// 使用 Ant Design 主题
	tenon.SetTheme(tenon.DefaultAntTheme())

	tenon.Run(NewApp(), 900, 650)
}
