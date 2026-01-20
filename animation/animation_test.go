package animation

import (
	"testing"
	"time"
)

func TestEasingFunctions(t *testing.T) {
	easingTests := []struct {
		name     string
		fn       EasingFunc
		t        float64
		expected float64
		epsilon  float64
	}{
		{"Linear 0", Linear, 0, 0, 0.001},
		{"Linear 0.5", Linear, 0.5, 0.5, 0.001},
		{"Linear 1", Linear, 1, 1, 0.001},

		{"EaseInQuad 0", EaseInQuad, 0, 0, 0.001},
		{"EaseInQuad 0.5", EaseInQuad, 0.5, 0.25, 0.001},
		{"EaseInQuad 1", EaseInQuad, 1, 1, 0.001},

		{"EaseOutQuad 0", EaseOutQuad, 0, 0, 0.001},
		{"EaseOutQuad 0.5", EaseOutQuad, 0.5, 0.75, 0.001},
		{"EaseOutQuad 1", EaseOutQuad, 1, 1, 0.001},

		{"EaseInCubic 0.5", EaseInCubic, 0.5, 0.125, 0.001},
		{"EaseOutCubic 0.5", EaseOutCubic, 0.5, 0.875, 0.001},
	}

	for _, test := range easingTests {
		result := test.fn(test.t)
		if diff := result - test.expected; diff < 0 {
			diff = -diff
		} else if diff > test.epsilon {
			t.Errorf("%s: expected %.3f, got %.3f", test.name, test.expected, result)
		}
	}
}

func TestTweenBasic(t *testing.T) {
	tween := NewTween(0, 100, 1*time.Second, Linear)

	// At start
	if tween.GetValue() != 0 {
		t.Errorf("start value should be 0, got %v", tween.GetValue())
	}

	// Halfway through
	tween.Update(0.5)
	if tween.GetValue() < 49 || tween.GetValue() > 51 {
		t.Errorf("halfway value should be ~50, got %v", tween.GetValue())
	}

	// At end
	tween.Update(0.5)
	if tween.GetValue() != 100 {
		t.Errorf("end value should be 100, got %v", tween.GetValue())
	}

	if !tween.IsComplete() {
		t.Error("tween should be complete")
	}
}

func TestTweenEasing(t *testing.T) {
	tween := NewTween(0, 100, 1*time.Second, EaseInQuad)

	// At 0.5 seconds, with EaseInQuad, we should be at 0.5^2 = 0.25
	tween.Update(0.5)
	expected := 25.0
	actual := tween.GetValue()

	if actual < expected-1 || actual > expected+1 {
		t.Errorf("eased value should be ~%v, got %v", expected, actual)
	}
}

func TestTweenCallback(t *testing.T) {
	completed := false
	tween := NewTween(0, 100, 100*time.Millisecond, Linear)
	tween.SetOnComplete(func() {
		completed = true
	})

	tween.Update(0.15)
	if tween.IsComplete() {
		if !completed {
			t.Error("onComplete callback should have been called")
		}
	}
}

func TestColorTween(t *testing.T) {
	ct := NewColorTween(0, 0, 0, 255, 255, 255, 1*time.Second, Linear)

	// At start
	r, g, b := ct.GetColor()
	if r != 0 || g != 0 || b != 0 {
		t.Errorf("start color should be (0,0,0), got (%d,%d,%d)", r, g, b)
	}

	// Halfway
	ct.Update(0.5)
	r, g, b = ct.GetColor()
	if r < 120 || r > 140 {
		t.Errorf("halfway R should be ~127, got %d", r)
	}

	// At end
	ct.Update(0.5)
	r, g, b = ct.GetColor()
	if r != 255 || g != 255 || b != 255 {
		t.Errorf("end color should be (255,255,255), got (%d,%d,%d)", r, g, b)
	}
}

func TestSequenceTween(t *testing.T) {
	t1 := NewTween(0, 100, 100*time.Millisecond, Linear)
	t2 := NewTween(100, 0, 100*time.Millisecond, Linear)

	seq := NewSequenceTween(t1, t2)

	// First tween should progress
	seq.Update(0.05)
	if t1.GetValue() < 40 || t1.GetValue() > 60 {
		t.Errorf("first tween should be ~50, got %v", t1.GetValue())
	}

	// Complete first tween
	seq.Update(0.06)
	if !t1.IsComplete() {
		t.Error("first tween should be complete")
	}

	// Second tween should have started
	seq.Update(0.05)
	if t2.GetValue() < 40 || t2.GetValue() > 60 {
		t.Errorf("second tween should be ~50, got %v", t2.GetValue())
	}
}

func TestParallelTween(t *testing.T) {
	t1 := NewTween(0, 100, 100*time.Millisecond, Linear)
	t2 := NewTween(100, 0, 100*time.Millisecond, Linear)

	par := NewParallelTween(t1, t2)

	// Both tweens should progress together
	par.Update(0.05)

	if t1.GetValue() < 40 || t1.GetValue() > 60 {
		t.Errorf("tween 1 should be ~50, got %v", t1.GetValue())
	}

	if t2.GetValue() < 40 || t2.GetValue() > 60 {
		t.Errorf("tween 2 should be ~50, got %v", t2.GetValue())
	}

	// Complete both
	par.Update(0.06)
	if !par.IsComplete() {
		t.Error("parallel tween should be complete")
	}
}

func TestAnimator(t *testing.T) {
	animator := NewAnimator(60)

	frameCount := 0
	animator.AddAnimation(func(frame int, dt float64) bool {
		frameCount++
		return frameCount >= 3
	})

	animator.Start()

	// Wait for animations to complete
	if !animator.WaitForCompletion(1 * time.Second) {
		t.Error("animator should have completed")
	}

	if frameCount != 3 {
		t.Errorf("expected 3 frames, got %d", frameCount)
	}

	animator.Stop()
}

func TestAnimatorRemovesComplete(t *testing.T) {
	animator := NewAnimator(60)

	// Add animation that completes immediately
	animator.AddAnimation(func(frame int, dt float64) bool {
		return true
	})

	animator.Start()
	time.Sleep(50 * time.Millisecond) // Give it time to process

	if animator.GetAnimationCount() != 0 {
		t.Errorf("expected 0 animations after completion, got %d", animator.GetAnimationCount())
	}

	animator.Stop()
}

func TestAnimatorFrameRate(t *testing.T) {
	animator := NewAnimator(120)
	animator.SetFrameRate(60)

	// Just test that it doesn't crash
	animator.Start()
	animator.Stop()
}
