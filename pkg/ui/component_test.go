package ui

import (
	"context"
	"fmt"
	"strings"
	"testing"
)

type objectFixture struct {
	Component
	render func(*objectFixture) NodeSource
	value  string
}

func (p *objectFixture) Render() NodeSource { return p.render(p) }

func TestEmbeddedComponentContextUpdatesFieldsLocally(t *testing.T) {
	parentRenders, childRenders := 0, 0
	p := &objectFixture{value: "initial"}
	p.render = func(p *objectFixture) NodeSource {
		childRenders++
		return Button(Text(p.value), OnClick(func() {
			p.Context().Update(func() { p.value = "updated" })
		}))
	}
	h := MountDefault(View(func() *Node {
		parentRenders++
		return Div(Render(p))
	}))
	if p.Context() == nil || !p.Context().Mounted() || h.Root().Count() != 1 {
		t.Fatal("embedded base was not bound to a local context")
	}
	h.Root().ByText("initial").Click()
	if !h.Root().ByText("updated").Exists() || parentRenders != 1 || childRenders != 2 {
		t.Fatalf("field update did not rerender locally: parent=%d child=%d texts=%v", parentRenders, childRenders, h.Root().Texts())
	}
}

func TestObjectContextCancelsBeforeCleanupAndIgnoresLateWork(t *testing.T) {
	p := &objectFixture{value: "mounted"}
	cleanups := 0
	p.render = func(p *objectFixture) NodeSource {
		ctx := p.Context()
		ctx.Effect(func() Cleanup {
			return func() {
				cleanups++
				if ctx.Err() != context.Canceled {
					t.Error("component context was not canceled before effect cleanup")
				}
			}
		})
		return Text(p.value)
	}
	h := MountDefault(View(func() *Node {
		visible, set := UseState(true)
		return Div(If(visible, Render(p)), Button(Text("Remove"), OnClick(func() { set(false) })))
	}))
	ctx := p.Context()
	ctx.Post(func() { p.value = "queued before unmount" })
	h.Root().ByText("Remove").Click()
	if cleanups != 1 || ctx.Mounted() || ctx.Err() != context.Canceled {
		t.Fatal("object component lifetime did not end at unmount")
	}
	ctx.Update(func() { p.value = "late update" })
	ctx.Post(func() { p.value = "late post" })
	drainPosts()
	if p.value != "mounted" {
		t.Fatalf("unmounted component executed a late mutation: %q", p.value)
	}
}

func TestObjectParentLifetimeCancelsChildren(t *testing.T) {
	child := &objectFixture{}
	child.render = func(*objectFixture) NodeSource { return Text("Child") }
	parent := &objectFixture{}
	parent.render = func(*objectFixture) NodeSource { return Div(Render(child)) }
	h := MountDefault(View(func() *Node {
		visible, setVisible := UseState(true)
		return Div(If(visible, Render(parent)), Button(Text("Hide parent"), OnClick(func() { setVisible(false) })))
	}))
	parentContext, childContext := parent.Context(), child.Context()
	h.Root().ByText("Hide parent").Click()
	if parentContext.Err() != context.Canceled || childContext.Err() != context.Canceled {
		t.Fatal("unmounting an object parent did not cancel the context subtree")
	}
}

func TestObjectContextPostsFromBackground(t *testing.T) {
	p := &objectFixture{value: "before"}
	p.render = func(p *objectFixture) NodeSource { return Text(p.value) }
	h := MountDefault(Render(p))
	posted := make(chan struct{})
	go func() {
		p.Context().Post(func() { p.value = "after" })
		close(posted)
	}()
	<-posted
	if p.value != "before" {
		t.Fatal("background Post executed on the caller goroutine")
	}
	drainPosts()
	h.Step(0)
	if !h.Root().ByText("after").Exists() {
		t.Fatal("posted component mutation did not render")
	}
}

func TestObjectContextReadsProviderThemeAndViewport(t *testing.T) {
	key := CreateContext("fallback")
	p := &objectFixture{}
	p.render = func(p *objectFixture) NodeSource {
		ctx := p.Context()
		return Text(fmt.Sprintf("%s/%v/%.0f", ContextValueIn(ctx, key), ctx.Theme() == DarkTheme, ctx.Viewport().W))
	}
	h := Mount(View(func() *Node {
		changed, set := UseState(false)
		value, theme := "first", LightTheme
		if changed {
			value, theme = "second", DarkTheme
		}
		return ThemeProvider(theme, key.Provider(value,
			Div(Style(Fill), Render(p), Button(Text("Change"), OnClick(func() { set(true) })))))
	}), 400, 240)
	if !h.Root().ByText("first/false/400").Exists() {
		t.Fatal("context did not resolve its provider and window scope")
	}
	h.Root().ByText("Change").Click()
	h.Resize(300, 200)
	if !h.Root().ByText("second/true/300").Exists() {
		t.Fatalf("context dependencies did not update: %v", h.Root().Texts())
	}
}

type objectTypeOne struct{ Component }

func (p *objectTypeOne) Render() NodeSource {
	count, set := StateIn(p.Context(), 1)
	return Button(Text(fmt.Sprint(count)), OnClick(func() { set(count + 1) }))
}

type objectTypeTwo struct{ Component }

func (p *objectTypeTwo) Render() NodeSource {
	count, _ := StateIn(p.Context(), 99)
	return Text(fmt.Sprint(count))
}

func TestObjectConcreteTypesDoNotShareHookState(t *testing.T) {
	one, two := &objectTypeOne{}, &objectTypeTwo{}
	h := MountDefault(View(func() *Node {
		second, set := UseState(false)
		var child IComponent = one
		if second {
			child = two
		}
		return Div(Render(child), Button(Text("Switch"), OnClick(func() { set(true) })))
	}))
	h.Root().ByText("1").Click()
	h.Root().ByText("Switch").Click()
	if !h.Root().ByText("99").Exists() || one.Context().Mounted() {
		t.Fatal("object type change did not replace the component lifetime and state")
	}
}

func TestObjectRejectsSharedInstanceAndRestoresRenderScope(t *testing.T) {
	p := &objectFixture{}
	p.render = func(*objectFixture) NodeSource { return Text("Shared") }
	defer func() {
		if p.Context() != nil {
			p.Context().cancel()
		}
		if r := recover(); r == nil || !strings.Contains(fmt.Sprint(r), "two places") {
			t.Fatalf("shared instance was not rejected: %v", r)
		}
		if currentFiber != nil {
			t.Fatal("failed component render left a stale Hooks/context scope")
		}
	}()
	MountDefault(Div(Render(p), Render(p)))
}
