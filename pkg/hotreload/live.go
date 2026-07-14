// Package hotreload provides an in-process live preview of Tenon UI written in
// plain Go, interpreted with yaegi. Edit a component source file and the running
// window updates immediately — no rebuild, no restart.
//
// The previewed file must declare `func View() *ui.Node` (any package name):
//
//	package view
//	import ui "github.com/sjm1327605995/tenon/pkg/ui"
//	func View() *ui.Node { return ui.Text("hello") }
//
// Constraints (yaegi): the interpreted code may call any NON-generic API in
// pkg/ui and pkg/shadcn (Div/Text/Button/Style/… and every shadcn component,
// whose generic internals stay compiled), but must NOT use generics itself —
// so no ui.UseState/ui.Use in the previewed file (it renders a stateless view).
package hotreload

import (
	"fmt"
	"go/parser"
	"go/token"
	"os"
	"time"

	"github.com/traefik/yaegi/interp"
	"github.com/traefik/yaegi/stdlib"

	ui "github.com/sjm1327605995/tenon/pkg/ui"
)

// Symbols holds the extracted pkg/ui + pkg/shadcn symbol tables (populated by the
// generated *_ui.go / *_shadcn.go init functions), exposed to interpreted code.
var Symbols = interp.Exports{}

// Interpret reads file, interprets it, and returns the node produced by its
// View() function. Any interpret error or panic (e.g. use of generics) is
// returned as an error rather than crashing the host.
func Interpret(file string) (node *ui.Node, err error) {
	defer func() {
		if r := recover(); r != nil {
			node, err = nil, fmt.Errorf("%v", r)
		}
	}()
	src, err := os.ReadFile(file)
	if err != nil {
		return nil, err
	}
	pf, err := parser.ParseFile(token.NewFileSet(), file, src, parser.PackageClauseOnly)
	if err != nil {
		return nil, err
	}
	i := interp.New(interp.Options{})
	if err := i.Use(stdlib.Symbols); err != nil {
		return nil, err
	}
	if err := i.Use(Symbols); err != nil {
		return nil, err
	}
	if _, err := i.Eval(string(src)); err != nil {
		return nil, err
	}
	fn, err := i.Eval(pf.Name.Name + ".View")
	if err != nil {
		return nil, fmt.Errorf("no `func View() *ui.Node`: %w", err)
	}
	out := fn.Call(nil)
	n, ok := out[0].Interface().(*ui.Node)
	if !ok || n == nil {
		return nil, fmt.Errorf("View() must return a non-nil *ui.Node")
	}
	return n, nil
}

// LivePreview renders the View() from file and hot-reloads it whenever the file
// changes on disk. Interpret errors are shown as an overlay (the last good view
// stays visible underneath).
func LivePreview(file string) *ui.Node { return ui.Use(livePreview, file) }

func livePreview(file string) *ui.Node {
	tree, setTree := ui.UseState[*ui.Node](nil)
	errMsg, setErr := ui.UseState("")

	ui.UseEffect(func() ui.Cleanup {
		reload := func() {
			if n, err := Interpret(file); err != nil {
				setErr(err.Error())
			} else {
				setErr("")
				setTree(n)
			}
		}
		reload() // 首次
		stop := make(chan struct{})
		go watch(file, stop, func() { ui.Post(reload) })
		return func() { close(stop) }
	})

	if errMsg != "" {
		return errorOverlay(errMsg, tree)
	}
	if tree == nil {
		return ui.Text("加载中…")
	}
	return tree
}

// watch 轮询文件修改时间，变化时回调（简单可靠、跨平台，无需 fsnotify 依赖）。
func watch(file string, stop <-chan struct{}, onChange func()) {
	last := modTime(file)
	t := time.NewTicker(250 * time.Millisecond)
	defer t.Stop()
	for {
		select {
		case <-stop:
			return
		case <-t.C:
			if m := modTime(file); !m.Equal(last) {
				last = m
				onChange()
			}
		}
	}
}

func modTime(file string) time.Time {
	if fi, err := os.Stat(file); err == nil {
		return fi.ModTime()
	}
	return time.Time{}
}

func errorOverlay(msg string, last *ui.Node) *ui.Node {
	kids := []*ui.Node{
		ui.Style(ui.Column, ui.Fill, ui.Gap(12), ui.Padding(16)),
		ui.Div(ui.Style(ui.Column, ui.Gap(6), ui.Padding(14), ui.Radius(8),
			ui.Bg(ui.Hex("#7f1d1d")), ui.Border(1, ui.Hex("#b91c1c"))),
			ui.Text("编译/解释错误", ui.Bold, ui.TextColor(ui.White)),
			ui.Text(msg, ui.FontSize(13), ui.TextColor(ui.Hex("#fecaca"))),
		),
	}
	if last != nil {
		kids = append(kids, ui.Div(ui.Style(ui.Opacity(0.5)), last))
	}
	return ui.Div(kids...)
}
