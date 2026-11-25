package main

import (
	"github.com/sjm1327605995/tenon/example/projectedparity/shared"
	ui "github.com/sjm1327605995/tenon/pkg/ui"
)

func main() {
	ui.WindowTitle("Tenon — Gio projected image")
	ui.WindowSize(420, 320)
	ui.Run(shared.Root())
}
