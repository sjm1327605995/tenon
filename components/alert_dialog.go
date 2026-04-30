package components

import (
	"github.com/sjm1327605995/tenon/internal/core"
	"github.com/sjm1327605995/tenon/internal/native"
	"github.com/sjm1327605995/tenon/yoga"
)

// AlertDialog is a confirmation dialog with cancel and action buttons.
type AlertDialog struct {
	Modal
	titleText *native.Text
	descText  *native.Text
	cancelBtn *Button
	actionBtn *Button
	onConfirm func()
}

// NewAlertDialog creates an alert dialog.
func NewAlertDialog() *AlertDialog {
	theme := core.GetTheme()
	ad := &AlertDialog{}
	ad.Init(ad)
	ad.SetVisible(false)
	ad.SetPositionType(yoga.PositionTypeAbsolute)
	ad.SetPosition(yoga.EdgeLeft, 0)
	ad.SetPosition(yoga.EdgeTop, 0)
	ad.SetWidthPercent(100)
	ad.SetHeightPercent(100)
	ad.SetFlexDirection(yoga.FlexDirectionColumn)
	ad.SetJustifyContent(yoga.JustifyCenter)
	ad.SetAlignItems(yoga.AlignCenter)
	ad.closeOnMask = true
	ad.closeOnEsc = true

	panel := native.NewView()
	panel.SetWidth(400)
	panel.SetMinHeight(180)
	panel.SetBackgroundColor(theme.CardColor)
	panel.SetBorderRadius(theme.BorderRadius)
	panel.SetPadding(yoga.EdgeAll, 24)
	panel.SetFlexDirection(yoga.FlexDirectionColumn)
	panel.SetGap(yoga.GutterAll, 16)

	ad.titleText = native.NewText("").SetFontSize(18).SetColor(theme.TextColor)
	ad.descText = native.NewText("").SetFontSize(14).SetColor(theme.MutedForegroundColor)

	btnRow := native.NewView().
		SetFlexDirection(yoga.FlexDirectionRow).
		SetGap(yoga.GutterAll, 8).
		SetJustifyContent(yoga.JustifyFlexEnd)

	ad.cancelBtn = NewButton("Cancel").SetVariant(ButtonOutline)
	ad.cancelBtn.SetOnClick(func() { ad.Close() })
	ad.actionBtn = NewButton("Continue").SetVariant(ButtonDefault)
	ad.actionBtn.SetOnClick(func() {
		if ad.onConfirm != nil {
			ad.onConfirm()
		}
		ad.Close()
	})
	btnRow.Add(ad.cancelBtn, ad.actionBtn)

	panel.Add(ad.titleText, ad.descText, btnRow)
	ad.AppendChild(panel)
	ad.panel = panel
	return ad
}

// ElementType returns type identifier.
func (ad *AlertDialog) ElementType() string { return "AlertDialog" }

// SetTitle sets the dialog title.
func (ad *AlertDialog) SetTitle(title string) *AlertDialog {
	ad.titleText.SetContent(title)
	ad.titleText.SetVisible(title != "")
	return ad
}

// SetDescription sets the dialog description.
func (ad *AlertDialog) SetDescription(desc string) *AlertDialog {
	ad.descText.SetContent(desc)
	ad.descText.SetVisible(desc != "")
	return ad
}

// SetActionText sets the action button text.
func (ad *AlertDialog) SetActionText(text string) *AlertDialog {
	ad.actionBtn.SetText(text)
	return ad
}

// SetCancelText sets the cancel button text.
func (ad *AlertDialog) SetCancelText(text string) *AlertDialog {
	ad.cancelBtn.SetText(text)
	return ad
}

// SetOnConfirm sets the confirm callback.
func (ad *AlertDialog) SetOnConfirm(fn func()) *AlertDialog {
	ad.onConfirm = fn
	return ad
}

// Open shows the alert dialog.
func (ad *AlertDialog) Open() *AlertDialog {
	ad.SetVisible(true)
	if eng := ad.GetEngine(); eng != nil {
		eng.AddOverlay(ad)
	}
	return ad
}

// Close hides the alert dialog.
func (ad *AlertDialog) Close() {
	ad.SetVisible(false)
	if eng := ad.GetEngine(); eng != nil {
		eng.RemoveOverlay(ad)
	}
	if ad.onClose != nil {
		ad.onClose()
	}
}
