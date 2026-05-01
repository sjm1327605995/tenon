package ui

import (
	"math"
	"testing"
	"time"
)

func TestAnimationControllerForwardTick(t *testing.T) {
	ac := &AnimationController{
		Duration:   1 * time.Second,
		LowerBound: 0,
		UpperBound: 1,
	}

	ac.Forward()
	if ac.Status != AnimationForward {
		t.Fatalf("expected status Forward, got %d", ac.Status)
	}

	// Tick half duration
	changed := ac.Tick(500 * time.Millisecond)
	if !changed {
		t.Fatal("expected changed=true")
	}
	if math.Abs(ac.Value-0.5) > 1e-9 {
		t.Fatalf("expected value 0.5, got %f", ac.Value)
	}

	// Tick remaining
	changed = ac.Tick(500 * time.Millisecond)
	if !changed {
		t.Fatal("expected changed=true")
	}
	if ac.Value != 1.0 {
		t.Fatalf("expected value 1.0, got %f", ac.Value)
	}
	if ac.Status != AnimationCompleted {
		t.Fatalf("expected status Completed, got %d", ac.Status)
	}
}

func TestAnimationControllerReverse(t *testing.T) {
	ac := &AnimationController{
		Duration:   1 * time.Second,
		LowerBound: 0,
		UpperBound: 1,
		Value:      1,
	}

	ac.Reverse()
	changed := ac.Tick(500 * time.Millisecond)
	if !changed {
		t.Fatal("expected changed=true")
	}
	if math.Abs(ac.Value-0.5) > 1e-9 {
		t.Fatalf("expected value 0.5, got %f", ac.Value)
	}

	changed = ac.Tick(500 * time.Millisecond)
	if !changed {
		t.Fatal("expected changed=true")
	}
	if ac.Value != 0 {
		t.Fatalf("expected value 0, got %f", ac.Value)
	}
	if ac.Status != AnimationDismissed {
		t.Fatalf("expected status Dismissed, got %d", ac.Status)
	}
}

func TestAnimationControllerStop(t *testing.T) {
	ac := &AnimationController{
		Duration:   1 * time.Second,
		LowerBound: 0,
		UpperBound: 1,
	}
	ac.Forward()
	ac.Stop()
	changed := ac.Tick(100 * time.Millisecond)
	if changed {
		t.Fatal("expected changed=false after stop")
	}
}

func TestAnimationControllerTickNoRun(t *testing.T) {
	ac := &AnimationController{
		Duration:   1 * time.Second,
		LowerBound: 0,
		UpperBound: 1,
	}
	changed := ac.Tick(100 * time.Millisecond)
	if changed {
		t.Fatal("expected changed=false when not running")
	}
}

func TestAnimationControllerZeroDuration(t *testing.T) {
	ac := &AnimationController{
		Duration:   0,
		LowerBound: 0,
		UpperBound: 1,
	}
	ac.Forward()
	changed := ac.Tick(1 * time.Millisecond)
	if !changed {
		t.Fatal("expected changed=true")
	}
	if ac.Value != 1.0 {
		t.Fatalf("expected value 1.0, got %f", ac.Value)
	}
	if ac.Status != AnimationCompleted {
		t.Fatalf("expected status Completed, got %d", ac.Status)
	}
}

func TestTweenFloat64(t *testing.T) {
	tw := &Tween[float64]{Begin: 0, End: 100}
	if tw.Evaluate(0.5) != 50 {
		t.Fatalf("expected 50, got %f", tw.Evaluate(0.5))
	}
	if tw.Evaluate(0) != 0 {
		t.Fatal("expected 0")
	}
	if tw.Evaluate(1) != 100 {
		t.Fatal("expected 100")
	}
}

func TestTweenFloat32(t *testing.T) {
	tw := &Tween[float32]{Begin: 0, End: 100}
	if tw.Evaluate(0.5) != 50 {
		t.Fatalf("expected 50, got %f", tw.Evaluate(0.5))
	}
}

func TestLinearCurve(t *testing.T) {
	c := LinearCurve{}
	if c.Transform(0.5) != 0.5 {
		t.Fatal("expected 0.5")
	}
	if c.Transform(0) != 0 {
		t.Fatal("expected 0")
	}
	if c.Transform(1) != 1 {
		t.Fatal("expected 1")
	}
}

func TestEaseInOutCurve(t *testing.T) {
	c := EaseInOutCurve{}
	if c.Transform(0) != 0 {
		t.Fatal("expected 0")
	}
	if c.Transform(1) != 1 {
		t.Fatal("expected 1")
	}
	mid := c.Transform(0.5)
	if math.Abs(mid-0.5) > 1e-9 {
		t.Fatalf("expected 0.5 at midpoint, got %f", mid)
	}
	// Ease-in-out should be slower at start and end
	if c.Transform(0.25) >= 0.25 {
		t.Fatal("expected eased value < 0.25 at t=0.25")
	}
	if c.Transform(0.75) <= 0.75 {
		t.Fatal("expected eased value > 0.75 at t=0.75")
	}
}

func TestAnimationValue(t *testing.T) {
	ac := &AnimationController{
		Duration:   1 * time.Second,
		LowerBound: 0,
		UpperBound: 1,
	}
	ac.Forward()
	ac.Tick(500 * time.Millisecond)

	anim := &Animation[float64]{
		Controller: ac,
		Tween:      &Tween[float64]{Begin: 0, End: 200},
		Curve:      LinearCurve{},
	}

	val := anim.Value()
	if math.Abs(val-100) > 1e-9 {
		t.Fatalf("expected animation value 100, got %f", val)
	}
}

func TestAnimationValueFloat32(t *testing.T) {
	ac := &AnimationController{
		Duration:   1 * time.Second,
		LowerBound: 0,
		UpperBound: 1,
	}
	ac.Forward()
	ac.Tick(250 * time.Millisecond)

	anim := &Animation[float32]{
		Controller: ac,
		Tween:      &Tween[float32]{Begin: 0, End: 100},
		Curve:      EaseInOutCurve{},
	}

	val := anim.Value()
	// progress = 0.25, eased value < 0.25, so val < 25
	if val >= 25 {
		t.Fatalf("expected eased float32 value < 25, got %f", val)
	}
}
