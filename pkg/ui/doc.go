// Package ui is a declarative, React-style GUI toolkit built on the yoga
// flexbox engine and the Gio renderer.
//
// The model has three layers: an immutable Node tree (the description you
// return from components), a persistent Fiber tree (identity, hooks, and
// reconciliation), and a renderNode tree (yoga nodes + painting). A state
// setter re-renders only its owning component, giving automatic, local
// updates without manual invalidation.
//
// Components use View for no props, or Use/Memo for typed props:
//
//	func Counter() *ui.Node {
//		count, setCount := ui.UseState(0)
//		return ui.HStack(16,
//			ui.Button(ui.Text("-")).OnClick(func() { setCount(count - 1) }),
//			ui.Text(fmt.Sprintf("%d", count)).FontSize(32),
//			ui.Button(ui.Text("+")).OnClick(func() { setCount(count + 1) }),
//		)
//	}
//
//	func main() { ui.Run(ui.View(Counter)) }
//
// Elements (Div, Span, Button, Input, Img, Text, ScrollView, Fragment,
// Portal) take attributes and children in any order. Styling is expressed
// with Style(...StyleOpt); hooks (UseState, UseEffect, UseReducer, UseMemo,
// UseCallback, UseRef, UseContext) follow React semantics, and animation is
// available via UseTween / UseTransition and the Animated layout style.
//
// See README.md in this directory for the full guide and the example/accordion
// program for a runnable demo.
//
// Rendering is single-threaded; trigger state changes from callbacks and
// effects, or use Post from background goroutines.
package ui
