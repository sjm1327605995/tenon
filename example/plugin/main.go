package main

import (
	"fmt"
	"time"

	ui "github.com/sjm1327605995/tenon/pkg/ui"
)

type particleProps struct{ Count int }

type particleAdapter struct{}

func (particleAdapter) Validate3D(node *ui.GameNode3D) error {
	props, ok := node.Native.(particleProps)
	if !ok || props.Count <= 0 {
		return fmt.Errorf("particles require a positive Count")
	}
	return nil
}

func (particleAdapter) Adapt3D(node *ui.GameNode3D) (ui.Game3DComponent, error) {
	return &particleSystem{count: node.Native.(particleProps).Count}, nil
}

type particleSystem struct {
	count    int
	updates  int
	shutdown int
}

func (p *particleSystem) Update(time.Duration) { p.updates++ }
func (p *particleSystem) Close()               { p.shutdown++ }

func main() {
	registry := ui.NewScene3DPluginRegistry()
	registry.RegisterPluginVersion("particles", 1, particleAdapter{})
	if _, err := registry.Build(ui.GameScene3D(
		ui.Plugin3DVersion("particles", 1, particleProps{}),
	)); err == nil {
		panic("invalid particle payload was accepted")
	}

	runtime, err := registry.Build(ui.GameScene3D(
		ui.Plugin3DVersion("particles", 1, particleProps{Count: 64}),
	))
	if err != nil {
		panic(err)
	}
	particles := runtime.Components()[0].(*particleSystem)
	runtime.Update(time.Second / 60)
	runtime.Close()
	runtime.Close()
	runtime.Activate()
	runtime.Update(time.Second / 60)
	runtime.Close()
	if particles.updates != 2 || particles.shutdown != 2 {
		panic(fmt.Sprintf("unexpected lifecycle: updates=%d shutdown=%d", particles.updates, particles.shutdown))
	}
	fmt.Printf("versioned plugin lifecycle passed: %d particles, 2 updates, 2 shutdowns\n", particles.count)
}
