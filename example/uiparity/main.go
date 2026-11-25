package main

import (
	"github.com/sjm1327605995/tenon/example/uiparity/shared"
	ui "github.com/sjm1327605995/tenon/pkg/ui"
)

func main() {
	ui.WindowTitle("Tenon — Gio UI parity")
	ui.WindowSize(480, 320)
	ui.Run(shared.Root(shared.NewState()))
}
