package debug

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/sjm1327605995/tenon/internal/core"
)

// RunWithDebug 启动 Tenon 应用并附带远程调试器。
func RunWithDebug(root core.Widget, width, height int, debugPort int) *Debugger {
	engine := core.NewEngine(root, width, height)

	d := NewDebugger(engine, debugPort)
	engine.SetDebugger(d)
	if err := d.Start(); err != nil {
		panic("failed to start debugger: " + err.Error())
	}

	engine.Mount()
	ebiten.SetWindowSize(width, height)
	ebiten.SetWindowResizingMode(ebiten.WindowResizingModeEnabled)
	if err := ebiten.RunGame(engine); err != nil {
		d.Stop()
		panic(err)
	}
	return d
}
