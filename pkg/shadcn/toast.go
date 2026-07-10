package shadcn

import (
	"fmt"
	"time"

	ui "github.com/sjm1327605995/tenon/pkg/ui"
)

type toastItem struct {
	id          int
	title, desc string
	variant     AlertVariant
	expiry      time.Time
}

// 全局通知存储；Toaster 挂载时注册 bump 以在变化时重渲染。
var toastState = struct {
	items  []toastItem
	nextID int
	bump   func()
}{}

// Toast 弹出一条通知（约 3 秒后自动消失）。可在任意回调中调用。
func Toast(title, desc string) { addToast(title, desc, AlertDefault) }

// ToastError 弹出一条错误样式通知。
func ToastError(title, desc string) { addToast(title, desc, AlertDestructive) }

func addToast(title, desc string, v AlertVariant) {
	toastState.nextID++
	toastState.items = append(toastState.items, toastItem{
		id: toastState.nextID, title: title, desc: desc, variant: v,
		expiry: time.Now().Add(3 * time.Second),
	})
	if toastState.bump != nil {
		toastState.bump()
	}
}

func removeToast(id int) {
	out := toastState.items[:0]
	for _, it := range toastState.items {
		if it.id != id {
			out = append(out, it)
		}
	}
	toastState.items = out
	if toastState.bump != nil {
		toastState.bump()
	}
}

// Toaster 应挂在应用根部；渲染活动通知（屏幕右下角），自动消失。
func Toaster() *ui.Node { return ui.Use(toaster, struct{}{}) }

func toaster(_ struct{}) *ui.Node {
	_, setV := ui.UseState(0)
	ver := ui.UseRef(0)
	ui.UseEffect(func() ui.Cleanup {
		toastState.bump = func() { *ver++; setV(*ver) }
		return func() { toastState.bump = nil }
	})
	if len(toastState.items) == 0 {
		return nil
	}
	stack := []*ui.Node{ui.Style(ui.Absolute, ui.Bottom(20), ui.Right(20),
		ui.Column, ui.Gap(10), ui.ItemsEnd)}
	for _, it := range toastState.items {
		item := it
		stack = append(stack, ui.Keyed(fmt.Sprintf("%d", item.id),
			ui.Use(toastView, toastViewProps{item: item, onExpire: func() { removeToast(item.id) }})))
	}
	return ui.Portal(ui.Div(stack...))
}

type toastViewProps struct {
	item     toastItem
	onExpire func()
}

func toastView(p toastViewProps) *ui.Node {
	th := ui.UseTheme()
	t := ui.UseElapsed() // 每帧检查过期
	fired := ui.UseRef(false)
	ui.UseEffect(func() ui.Cleanup {
		if !*fired && time.Now().After(p.item.expiry) {
			*fired = true
			p.onExpire()
		}
		return nil
	}, t)

	fg, border := th.Foreground, th.Border
	if p.item.variant == AlertDestructive {
		fg, border = th.Destructive, th.Destructive
	}
	kids := []*ui.Node{ui.Style(ui.Column, ui.Gap(4), ui.Padding(14), ui.MinWidth(240),
		ui.Bg(th.Popover), ui.TextColor(fg), ui.Border(1, border), ui.Radius(th.Radius))}
	kids = append(kids, ui.Text(p.item.title, ui.FontSize(14)))
	if p.item.desc != "" {
		kids = append(kids, ui.Text(p.item.desc, ui.FontSize(13), ui.TextColor(th.MutedForeground)))
	}
	return ui.Div(kids...)
}
