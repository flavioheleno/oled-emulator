package animation

import (
	"time"
)

// Tween represents a tweened animation between two values
type Tween struct {
	from       float64
	to         float64
	duration   time.Duration
	elapsed    time.Duration
	easing     EasingFunc
	onComplete func()
	onUpdate   func(value float64)
}

// NewTween creates a new tween animation
func NewTween(from, to float64, duration time.Duration, easing EasingFunc) *Tween {
	if easing == nil {
		easing = Linear
	}

	return &Tween{
		from:     from,
		to:       to,
		duration: duration,
		elapsed:  0,
		easing:   easing,
	}
}

// SetOnComplete sets a callback when the tween completes
func (t *Tween) SetOnComplete(fn func()) *Tween {
	t.onComplete = fn
	return t
}

// SetOnUpdate sets a callback called each frame with the current value
func (t *Tween) SetOnUpdate(fn func(value float64)) *Tween {
	t.onUpdate = fn
	return t
}

// GetValue returns the current interpolated value
func (t *Tween) GetValue() float64 {
	if t.duration == 0 {
		return t.to
	}

	normalizedTime := float64(t.elapsed) / float64(t.duration)
	if normalizedTime > 1 {
		normalizedTime = 1
	}

	easedTime := t.easing(normalizedTime)
	return t.from + (t.to-t.from)*easedTime
}

// IsComplete returns whether the tween has finished
func (t *Tween) IsComplete() bool {
	return t.elapsed >= t.duration
}

// GetProgress returns the progress (0 to 1)
func (t *Tween) GetProgress() float64 {
	if t.duration == 0 {
		return 1
	}

	progress := float64(t.elapsed) / float64(t.duration)
	if progress > 1 {
		progress = 1
	}
	return progress
}

// Update updates the tween with delta time
func (t *Tween) Update(dt float64) bool {
	if t.IsComplete() {
		return true
	}

	t.elapsed += time.Duration(dt * float64(time.Second))

	if t.elapsed > t.duration {
		t.elapsed = t.duration
	}

	value := t.GetValue()
	if t.onUpdate != nil {
		t.onUpdate(value)
	}

	if t.IsComplete() {
		if t.onComplete != nil {
			t.onComplete()
		}
		return true
	}

	return false
}

// CreateTweenAnimation creates an AnimationFunc from a Tween
func CreateTweenAnimation(from, to float64, duration time.Duration, easing EasingFunc) AnimationFunc {
	tween := NewTween(from, to, duration, easing)

	return func(frame int, dt float64) bool {
		return tween.Update(dt)
	}
}

// ColorTween tweens between two RGB colors
type ColorTween struct {
	fromR, toR   byte
	fromG, toG   byte
	fromB, toB   byte
	duration     time.Duration
	elapsed      time.Duration
	easing       EasingFunc
	onComplete   func()
	onUpdate     func(r, g, b byte)
}

// NewColorTween creates a new color tween
func NewColorTween(fromR, fromG, fromB, toR, toG, toB byte, duration time.Duration, easing EasingFunc) *ColorTween {
	if easing == nil {
		easing = Linear
	}

	return &ColorTween{
		fromR:    fromR,
		toR:      toR,
		fromG:    fromG,
		toG:      toG,
		fromB:    fromB,
		toB:      toB,
		duration: duration,
		elapsed:  0,
		easing:   easing,
	}
}

// SetOnComplete sets a callback when the tween completes
func (ct *ColorTween) SetOnComplete(fn func()) *ColorTween {
	ct.onComplete = fn
	return ct
}

// SetOnUpdate sets a callback called each frame with the current color
func (ct *ColorTween) SetOnUpdate(fn func(r, g, b byte)) *ColorTween {
	ct.onUpdate = fn
	return ct
}

// GetColor returns the current interpolated color
func (ct *ColorTween) GetColor() (byte, byte, byte) {
	if ct.duration == 0 {
		return ct.toR, ct.toG, ct.toB
	}

	normalizedTime := float64(ct.elapsed) / float64(ct.duration)
	if normalizedTime > 1 {
		normalizedTime = 1
	}

	easedTime := ct.easing(normalizedTime)

	r := byte(float64(ct.fromR) + (float64(ct.toR)-float64(ct.fromR))*easedTime)
	g := byte(float64(ct.fromG) + (float64(ct.toG)-float64(ct.fromG))*easedTime)
	b := byte(float64(ct.fromB) + (float64(ct.toB)-float64(ct.fromB))*easedTime)

	return r, g, b
}

// IsComplete returns whether the tween has finished
func (ct *ColorTween) IsComplete() bool {
	return ct.elapsed >= ct.duration
}

// Update updates the tween with delta time
func (ct *ColorTween) Update(dt float64) bool {
	if ct.IsComplete() {
		return true
	}

	ct.elapsed += time.Duration(dt * float64(time.Second))

	if ct.elapsed > ct.duration {
		ct.elapsed = ct.duration
	}

	r, g, b := ct.GetColor()
	if ct.onUpdate != nil {
		ct.onUpdate(r, g, b)
	}

	if ct.IsComplete() {
		if ct.onComplete != nil {
			ct.onComplete()
		}
		return true
	}

	return false
}

// SequenceTween chains multiple tweens together
type SequenceTween struct {
	tweens        []*Tween
	currentIndex  int
	onComplete    func()
}

// NewSequenceTween creates a new sequence tween
func NewSequenceTween(tweens ...*Tween) *SequenceTween {
	return &SequenceTween{
		tweens:       tweens,
		currentIndex: 0,
	}
}

// SetOnComplete sets a callback when all tweens complete
func (st *SequenceTween) SetOnComplete(fn func()) *SequenceTween {
	st.onComplete = fn
	return st
}

// Update updates the sequence
func (st *SequenceTween) Update(dt float64) bool {
	if st.currentIndex >= len(st.tweens) {
		if st.onComplete != nil {
			st.onComplete()
		}
		return true
	}

	currentTween := st.tweens[st.currentIndex]
	if currentTween.Update(dt) {
		st.currentIndex++

		if st.currentIndex >= len(st.tweens) {
			if st.onComplete != nil {
				st.onComplete()
			}
			return true
		}
	}

	return false
}

// IsComplete returns whether all tweens have finished
func (st *SequenceTween) IsComplete() bool {
	return st.currentIndex >= len(st.tweens)
}

// ParallelTween runs multiple tweens in parallel
type ParallelTween struct {
	tweens     []*Tween
	onComplete func()
}

// NewParallelTween creates a new parallel tween
func NewParallelTween(tweens ...*Tween) *ParallelTween {
	return &ParallelTween{
		tweens: tweens,
	}
}

// SetOnComplete sets a callback when all tweens complete
func (pt *ParallelTween) SetOnComplete(fn func()) *ParallelTween {
	pt.onComplete = fn
	return pt
}

// Update updates all tweens
func (pt *ParallelTween) Update(dt float64) bool {
	allComplete := true

	for _, tween := range pt.tweens {
		if !tween.IsComplete() {
			tween.Update(dt)
			allComplete = false
		}
	}

	if allComplete && pt.onComplete != nil {
		pt.onComplete()
	}

	return allComplete
}

// IsComplete returns whether all tweens have finished
func (pt *ParallelTween) IsComplete() bool {
	for _, tween := range pt.tweens {
		if !tween.IsComplete() {
			return false
		}
	}
	return true
}
