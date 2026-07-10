// Package ui is a declarative, React-style GUI toolkit built on the yoga
// flexbox engine and the Ebiten renderer (vector + text/v2).
//
// The model has three layers: an immutable Node tree (the description you
// return from components), a persistent Fiber tree (identity, hooks, and
// reconciliation), and a renderNode tree (yoga nodes + painting). A state
// setter re-renders only its owning component, giving automatic, local
// updates without manual invalidation.
//
// Components are functions of a typed props struct, mounted with Use or Memo:
//
//	func Counter(_ struct{}) *ui.Node {
//		count, setCount := ui.UseState(0)
//		return ui.Div(
//			ui.Style(ui.Row, ui.Gap(16), ui.ItemsCenter),
//			ui.Button(ui.OnClick(func() { setCount(count - 1) }), ui.Text("-")),
//			ui.Text(fmt.Sprintf("%d", count), ui.FontSize(32)),
//			ui.Button(ui.OnClick(func() { setCount(count + 1) }), ui.Text("+")),
//		)
//	}
//
//	func main() { ui.Run(ui.Use(Counter, struct{}{})) }
//
// Elements (Div, Span, Button, Input, Img, Text, ScrollView, Fragment,
// Portal) take attributes and children in any order. Styling is expressed
// with Style(...StyleOpt); hooks (UseState, UseEffect, UseReducer, UseMemo,
// UseCallback, UseRef, UseContext) follow React semantics, and animation is
// available via UseTween / UseTransition and the Animated layout style.
//
// See README.md in this directory for the full guide and the example/hooks-*
// programs for runnable demos.
//
// Rendering is single-threaded (Ebiten's Update/Draw loop); trigger state
// changes from callbacks and effects rather than other goroutines.
package ui
