package main

import (
	"os"

	"github.com/sjm1327605995/tenon/backend/g3n"
	"github.com/sjm1327605995/tenon/example/uiparity/shared"
	ui "github.com/sjm1327605995/tenon/pkg/ui"
)

func main() {
	state := shared.NewState()
	options := g3n.Options{}
	if os.Getenv("TENON_UI_SMOKE") == "1" {
		options.Launch = func(game *g3n.Game) { scheduleNativeInteraction(game, state) }
	}
	g3n.Use(options)
	ui.WindowTitle("Tenon — G3N UI parity")
	ui.WindowSize(480, 320)
	ui.Run(shared.Root(state))
}
