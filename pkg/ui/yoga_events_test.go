package ui

import (
	"testing"

	"github.com/sjm1327605995/tenon/yoga"
)

func TestUnmountPublishesYogaDeallocation(t *testing.T) {
	yoga.ResetEventSubscribers()
	t.Cleanup(yoga.ResetEventSubscribers)
	deallocations := 0
	yoga.Subscribe(yoga.EventNodeDeallocation, func(yoga.Event) { deallocations++ })

	h := MountDefault(Div(Text("child")))
	h.Close()
	if deallocations != 2 {
		t.Fatalf("deallocation events=%d, want 2 for root and text", deallocations)
	}
	h.Close()
	if deallocations != 2 {
		t.Fatal("closing a harness twice published duplicate deallocations")
	}
}
