package ui

import (
	"errors"
	"fmt"
	"strings"
	"testing"

	values "github.com/sjm1327605995/tenon/pkg/value"
)

func TestValueAndBindSynchronizeTextAndInputWithoutRerender(t *testing.T) {
	name := values.New("Alice")
	renders := 0
	h := MountDefault(View(func() *Node {
		renders++
		return VStack(8, Text(values.Static("Profile")), Text(name), Input(BindInput(name)))
	}))
	name.Set("Bob")
	h.Step(0)
	if !h.Root().ByText("Bob").Exists() || h.Root().ByKind("input").Value() != "Bob" || renders != 1 {
		t.Fatal("Set did not directly synchronize bound controls")
	}
	h.Root().ByKind("input").Focus().SetValue("Carol")
	if name.Get() != "Carol" || !h.Root().ByText("Carol").Exists() || renders != 1 {
		t.Fatal("input did not write back to Bind without a component update")
	}
}

func TestReplacingAndRemovingValueReleasesSubscriptions(t *testing.T) {
	activeA, activeB := 0, 0
	a := values.New("A", values.Hooks[string]{OnSubscribe: func() { activeA++ }, OnUnsubscribe: func() { activeA-- }})
	b := values.New("B", values.Hooks[string]{OnSubscribe: func() { activeB++ }, OnUnsubscribe: func() { activeB-- }})
	h := MountDefault(View(func() *Node {
		state, set := UseState(0)
		var content *Node
		if state < 2 {
			source := a
			if state == 1 {
				source = b
			}
			content = VStack(0, Text(source), Input(BindInput(source)))
		} else {
			content = Text("Removed")
		}
		return VStack(0, content, Button(Text("Next"), OnClick(func() { set(state + 1) })))
	}))
	if activeA != 2 {
		t.Fatalf("initial active subscriptions=%d, want 2", activeA)
	}
	h.Root().ByText("Next").Click()
	if activeA != 0 || activeB != 2 {
		t.Fatal("switching sources retained stale subscriptions")
	}
	a.Set("Old")
	if h.Root().ByText("Old").Exists() || h.Root().ByKind("input").Value() != "B" {
		t.Fatal("old source overwrote the current binding")
	}
	h.Root().ByText("Next").Click()
	if activeB != 0 {
		t.Fatal("unmounted bindings did not release subscriptions")
	}
}

type externalStringValue struct {
	text    string
	last    func(string)
	stopped int
}

func (s *externalStringValue) Get() string { return s.text }
func (s *externalStringValue) Subscribe(fn func(string)) func() {
	s.last = fn
	fn(s.text) // custom sources may emit synchronously while subscribing
	return func() { s.stopped++ }
}

func TestCustomValueSupportsImmediateAndLateNotifications(t *testing.T) {
	source := &externalStringValue{text: "External"}
	h := MountDefault(View(func() *Node {
		bound, set := UseState(true)
		var content any = "Static"
		if bound {
			content = source
		}
		return VStack(0, Text(content), Button(Text("Unbind"), OnClick(func() { set(false) })))
	}))
	if !h.Root().ByText("External").Exists() {
		t.Fatal("custom Value was not accepted")
	}
	h.Root().ByText("Unbind").Click()
	source.last("Late")
	if source.stopped != 1 || !h.Root().ByText("Static").Exists() || h.Root().ByText("Late").Exists() {
		t.Fatal("released custom Value callback overwrote a literal")
	}
}

type functionStringValue struct {
	get       func() string
	subscribe func(func(string)) func()
}

func (v functionStringValue) Get() string                      { return v.get() }
func (v functionStringValue) Subscribe(fn func(string)) func() { return v.subscribe(fn) }

func TestNonComparableValueCanBeRebound(t *testing.T) {
	active := 0
	source := functionStringValue{get: func() string { return "Custom" }, subscribe: func(func(string)) func() {
		active++
		return func() { active-- }
	}}
	h := MountDefault(View(func() *Node {
		tick, set := UseState(0)
		return VStack(0, Text(source), Button(Text("Tick"), OnClick(func() { set(tick + 1) })))
	}))
	for range 5 {
		h.Root().ByText("Tick").Click()
	}
	if active != 1 || !h.Root().ByText("Custom").Exists() {
		t.Fatal("non-comparable custom Value leaked subscriptions")
	}
}

func TestBoundInputUsesCanonicalValueAfterValidationAndSubmit(t *testing.T) {
	name := values.New("A", values.Hooks[string]{BeforeSet: func(old, next string) (string, error) {
		if strings.Contains(next, "!") {
			return old, errors.New("invalid")
		}
		return strings.ToUpper(next), nil
	}})
	var submitted string
	h := MountDefault(Input(BindInput(name), OnSubmit(func(v string) { submitted = v })))
	q := h.Root().ByKind("input").Focus()
	q.SetValue("b")
	q.SetValue("invalid!")
	if q.Value() != "B" || name.Get() != "B" {
		t.Fatal("rejected input was not restored to the canonical Bind value")
	}
	old := input
	hi := &harnessInput{chars: []rune("c")}
	hi.just[keyEnter], hi.held[keyEnter] = true, true
	input = hi
	t.Cleanup(func() { input = old })
	h.g.handleKeyboardNav()
	h.g.editFocusedInput()
	h.settle()
	if q.Value() != "BC" || submitted != "BC" || name.Get() != "BC" {
		t.Fatalf("submit did not use normalized text: field=%q data=%q submitted=%q", q.Value(), name.Get(), submitted)
	}
}

func TestUseValueRerendersOnlyItsConsumer(t *testing.T) {
	count := values.New(0)
	parentRenders, childRenders := 0, 0
	child := func(_ struct{}) *Node {
		childRenders++
		return Text(fmt.Sprint(UseValue(count)))
	}
	h := MountDefault(View(func() *Node {
		parentRenders++
		return VStack(0, Use(child, struct{}{}), Button(Text("Increment"), OnClick(func() { count.Set(count.Get() + 1) })))
	}))
	h.Root().ByText("Increment").Click()
	if parentRenders != 1 || childRenders != 2 || !h.Root().ByText("1").Exists() {
		t.Fatal("UseValue did not schedule only its subscribing component")
	}
}

func TestClosingHarnessReleasesAllValueConsumers(t *testing.T) {
	active := 0
	name := values.New("Alice", values.Hooks[string]{OnSubscribe: func() { active++ }, OnUnsubscribe: func() { active-- }})
	h := MountDefault(VStack(0, Text(name), Input(BindInput(name))))
	h.Close()
	h.Close()
	name.Set("After close")
	if active != 0 || h.Root().Exists() || len(h.Overlays()) != 0 {
		t.Fatal("closing a root retained value subscribers or render nodes")
	}
}
