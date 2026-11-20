package decl_test

import (
	"fmt"
	"testing"

	d "github.com/sjm1327605995/tenon/pkg/decl"
	ui "github.com/sjm1327605995/tenon/pkg/ui"
	"github.com/sjm1327605995/tenon/pkg/value"
)

type profileComponent struct {
	d.Component
	Title string
	Name  string
}

func (p *profileComponent) Render() d.View {
	ctx := p.Context()
	return d.Column{Width: 320, Gap: 8, Background: ctx.Theme().Background, Children: []d.View{
		d.Text{Value: p.Title, FontSize: 20},
		d.Input{Value: p.Name, Placeholder: "Name", OnChange: func(value string) {
			ctx.Update(func() { p.Name = value })
		}},
		d.Button{Text: "Clear", OnClick: func() { p.SetState(func() { p.Name = "" }) }},
		d.Text{Value: "Hello " + p.Name},
	}}
}

func TestEmbeddedComponentIsAcceptedAtRootAndAsChild(t *testing.T) {
	p := &profileComponent{Title: "Profile"}
	var component d.IComponent = p
	h := d.Mount(component, 400, 240)
	h.Root().ByPlaceholder("Name").Focus().Type("Alice")
	if p.Name != "Alice" || !h.Root().ByText("Hello Alice").Exists() || p.Context() == nil {
		t.Fatal("component interface did not bind context and update fields")
	}
	h.Root().ByText("Clear").Click()
	if p.Name != "" || h.Root().ByPlaceholder("Name").Value() != "" {
		t.Fatal("embedded SetState did not update the component")
	}
	child := &profileComponent{Title: "Nested"}
	nested := d.Mount(func() d.View { return d.Column{Children: []d.View{child}} }, 400, 240)
	if !nested.Root().ByText("Nested").Exists() || nested.Root().Count() != 1 {
		t.Fatal("automatic child dispatch invoked the embedded base instead of the concrete Render")
	}
}

type statefulPropsComponent struct {
	d.Component
	Label string
}

func (p *statefulPropsComponent) Render() d.View {
	count, set := d.State(p.Context(), 0)
	return d.Button{Text: fmt.Sprintf("%s:%d", p.Label, count), OnClick: func() { set(count + 1) }}
}

func TestContextStateSurvivesFreshSameTypeProps(t *testing.T) {
	h := d.Mount(func() d.View {
		changed, set := ui.UseState(false)
		label := "Before"
		if changed {
			label = "After"
		}
		return d.Column{Children: []d.View{
			&statefulPropsComponent{Label: label},
			d.Button{Text: "Change props", OnClick: func() { set(true) }},
		}}
	}, 400, 240)
	h.Root().ByText("Before:0").Click()
	h.Root().ByText("Change props").Click()
	if !h.Root().ByText("After:1").Exists() {
		t.Fatalf("fresh props reset context state: %v", h.Root().Texts())
	}
}

var _ d.IComponent = (*profileComponent)(nil)

func TestDeclarationsAcceptValueAndBind(t *testing.T) {
	name := value.New("Alice")
	h := d.Mount(d.Column{Width: 300, Children: []d.View{
		d.Text{Bind: name}, d.Input{Bind: name, Placeholder: "Bound"}, d.Textarea{Bind: name, Placeholder: "Notes"},
	}}, 400, 300)
	h.Root().ByPlaceholder("Bound").SetValue("Bob")
	if name.Get() != "Bob" || !h.Root().ByText("Bob").Exists() || h.Root().ByPlaceholder("Notes").Value() != "Bob" {
		t.Fatal("declaration Bind fields did not share the same data source")
	}
}
