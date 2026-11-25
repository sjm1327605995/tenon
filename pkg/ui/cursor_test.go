package ui

import "testing"

func TestBackendRuntimePointerCursor(t *testing.T) {
	root := Div(Style(Column),
		Button(Style(Width(100), Height(30)), OnClick(func() {}), Text("button")),
		Input(Style(Width(100), Height(30)), Placeholder("input")),
		Button(Style(Width(100), Height(30)), Disabled(true), OnClick(func() {}), Text("disabled")),
	)
	runtime := NewBackendRuntime(root, windowConfig{w: 300, h: 200}, &harnessInput{})
	defer runtime.Close()
	runtime.Frame(0)
	for _, test := range []struct {
		x, y float32
		want CursorRole
	}{
		{x: 50, y: 15, want: CursorPointer},
		{x: 50, y: 45, want: CursorText},
		{x: 50, y: 75, want: CursorDefault},
		{x: 200, y: 150, want: CursorDefault},
	} {
		if got := runtime.PointerCursor(test.x, test.y); got != test.want {
			t.Errorf("cursor at (%v,%v) = %v, want %v", test.x, test.y, got, test.want)
		}
	}
	runtime.Close()
	if got := runtime.PointerCursor(50, 15); got != CursorDefault {
		t.Fatalf("closed runtime cursor = %v, want default", got)
	}
}
