package ui

import (
	"testing"

	"gioui.org/f32"
	"gioui.org/io/event"
	"gioui.org/io/pointer"
)

type scriptedInputSource struct {
	events []event.Event
}

func (s *scriptedInputSource) Event(...event.Filter) (event.Event, bool) {
	if len(s.events) == 0 {
		return nil, false
	}
	e := s.events[0]
	s.events = s.events[1:]
	return e, true
}

func TestGioPointerPressReachesRuntimeClick(t *testing.T) {
	clicked := false
	tag := new(int)
	input := &gioInput{tag: tag, clip: &gioClipboardState{tag: tag}}
	runtime := NewBackendRuntime(Button(Style(Width(100), Height(40)), OnClick(func() { clicked = true }), Text("Click")), windowConfig{w: 100, h: 40}, input)
	defer runtime.Close()
	runtime.Frame(0)
	input.resetFrame()
	source := &scriptedInputSource{events: []event.Event{pointer.Event{
		Kind: pointer.Press, Position: f32.Pt(20, 20), Buttons: pointer.ButtonPrimary,
	}}}
	input.process(source)
	runtime.Frame(0)
	if !clicked {
		t.Fatal("Gio pointer press did not reach the runtime click callback")
	}
}
