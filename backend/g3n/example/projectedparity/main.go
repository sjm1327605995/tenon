package main

import (
	"github.com/sjm1327605995/tenon/backend/g3n"
	"github.com/sjm1327605995/tenon/example/projectedparity/shared"
	ui "github.com/sjm1327605995/tenon/pkg/ui"
)

func main() {
	g3n.Use(g3n.Options{})
	ui.WindowTitle("Tenon — G3N projected image")
	ui.WindowSize(420, 320)
	ui.Run(shared.Root())
}
