package value_test

import (
	"errors"
	"fmt"
	"reflect"
	"strings"
	"testing"

	v "github.com/sjm1327605995/tenon/pkg/value"
)

func TestBindNormalizesRejectsAndNotifies(t *testing.T) {
	var events []string
	b := v.New("initial", v.Hooks[string]{
		BeforeSet: func(old, next string) (string, error) {
			if strings.Contains(next, "!") {
				return old, errors.New("invalid")
			}
			return strings.TrimSpace(next), nil
		},
		OnChange: func(old, next string) { events = append(events, old+"->"+next) },
		OnError:  func(err error) { events = append(events, err.Error()) },
	})
	stop := b.Subscribe(func(next string) { events = append(events, "subscriber:"+next) })
	b.Set(" Alice ")
	b.Set("Alice")
	b.Set("invalid!")
	stop()
	stop()
	b.Set("Bob")
	want := []string{"initial->Alice", "subscriber:Alice", "invalid", "Alice->Bob"}
	if b.Get() != "Bob" || !reflect.DeepEqual(events, want) {
		t.Fatalf("binding hooks or notification order mismatch: value=%q events=%v", b.Get(), events)
	}
}

func TestSubscriptionsHandleCancellationAndReentrantWrites(t *testing.T) {
	bound, unbound := 0, 0
	b := v.New(0, v.Hooks[int]{OnSubscribe: func() { bound++ }, OnUnsubscribe: func() { unbound++ }})
	var events []string
	stopA := b.Subscribe(func(next int) {
		events = append(events, fmt.Sprintf("a%d", next))
		if next == 1 {
			b.Set(2)
		}
		if b.Get() != next {
			t.Error("reentrant Set changed the snapshot during an active notification")
		}
	})
	stopB := b.Subscribe(func(next int) { events = append(events, fmt.Sprintf("b%d", next)) })
	b.Set(1)
	stopA()
	stopA()
	stopB()
	if !reflect.DeepEqual(events, []string{"a1", "b1", "a2", "b2"}) || bound != 2 || unbound != 2 {
		t.Fatalf("reentrant notifications or lifecycle hooks mismatch: %v bound=%d unbound=%d", events, bound, unbound)
	}
}

func TestCancelDuringNotificationAndStaticMappedValues(t *testing.T) {
	b := v.New(0)
	var events []string
	var stopB func()
	b.Subscribe(func(next int) { events = append(events, "a"); stopB() })
	stopB = b.Subscribe(func(int) { events = append(events, "b") })
	b.Set(1)
	if !reflect.DeepEqual(events, []string{"a"}) {
		t.Fatal("canceled subscriber received the in-progress notification")
	}
	constant := v.Static("constant")
	stop := constant.Subscribe(func(string) { t.Error("static value unexpectedly emitted") })
	stop()
	stop()
	derived := v.Map(b, func(n int) string { return fmt.Sprintf("Count %d", n) })
	var got string
	unsub := derived.Subscribe(func(s string) { got = s })
	b.Set(3)
	if constant.Get() != "constant" || derived.Get() != "Count 3" || got != "Count 3" {
		t.Fatal("static or derived value did not follow the shared Value protocol")
	}
	unsub()
}
