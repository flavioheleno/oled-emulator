package animation

import (
	"sync"
	"time"
)

// AnimationFunc defines a callback function for animations
// frame: current frame number
// dt: delta time in seconds since last frame
// Returns true when animation is complete
type AnimationFunc func(frame int, dt float64) bool

// Animator manages frame-based animations
type Animator struct {
	mu           sync.Mutex
	fps          int
	targetDt     float64
	ticker       *time.Ticker
	running      bool
	frameCount   int
	animations   []AnimationFunc
	lastTime     time.Time
	stopChan     chan struct{}
	onFrame      func(frame int, dt float64)
}

// NewAnimator creates a new animator with the specified FPS
func NewAnimator(fps int) *Animator {
	if fps <= 0 {
		fps = 60
	}

	return &Animator{
		fps:      fps,
		targetDt: 1.0 / float64(fps),
		running:  false,
		stopChan: make(chan struct{}),
	}
}

// SetFrameRate sets the target frame rate
func (a *Animator) SetFrameRate(fps int) {
	if fps <= 0 {
		fps = 60
	}

	a.mu.Lock()
	defer a.mu.Unlock()

	a.fps = fps
	a.targetDt = 1.0 / float64(fps)
}

// AddAnimation adds an animation function
func (a *Animator) AddAnimation(fn AnimationFunc) {
	a.mu.Lock()
	defer a.mu.Unlock()

	a.animations = append(a.animations, fn)
}

// SetOnFrame sets a callback called every frame
func (a *Animator) SetOnFrame(fn func(frame int, dt float64)) {
	a.mu.Lock()
	defer a.mu.Unlock()

	a.onFrame = fn
}

// Start begins animation
func (a *Animator) Start() {
	a.mu.Lock()
	if a.running {
		a.mu.Unlock()
		return
	}

	a.running = true
	a.frameCount = 0
	a.lastTime = time.Now()
	a.ticker = time.NewTicker(time.Duration(float64(time.Second) / float64(a.fps)))
	a.mu.Unlock()

	// Run animation loop in goroutine
	go a.loop()
}

// Stop halts animation
func (a *Animator) Stop() {
	a.mu.Lock()
	if !a.running {
		a.mu.Unlock()
		return
	}

	a.running = false
	if a.ticker != nil {
		a.ticker.Stop()
		a.ticker = nil
	}
	a.mu.Unlock()

	// Signal the loop to stop
	select {
	case a.stopChan <- struct{}{}:
	default:
	}
}

// loop runs the animation update loop
func (a *Animator) loop() {
	for {
		a.mu.Lock()
		if a.ticker == nil {
			a.mu.Unlock()
			return
		}
		ticker := a.ticker
		a.mu.Unlock()

		select {
		case <-ticker.C:
			a.update()

		case <-a.stopChan:
			return
		}
	}
}

// update processes animations for the current frame
func (a *Animator) update() {
	a.mu.Lock()
	defer a.mu.Unlock()

	if !a.running {
		return
	}

	now := time.Now()
	dt := now.Sub(a.lastTime).Seconds()
	a.lastTime = now

	// Call onFrame callback if set
	if a.onFrame != nil {
		a.onFrame(a.frameCount, dt)
	}

	// Update animations
	completedCount := 0
	activeAnimations := make([]AnimationFunc, 0, len(a.animations))

	for _, anim := range a.animations {
		isComplete := anim(a.frameCount, dt)
		if isComplete {
			completedCount++
		} else {
			activeAnimations = append(activeAnimations, anim)
		}
	}

	// Remove completed animations
	a.animations = activeAnimations

	a.frameCount++
}

// IsRunning returns whether animations are currently running
func (a *Animator) IsRunning() bool {
	a.mu.Lock()
	defer a.mu.Unlock()

	return a.running
}

// GetFrameCount returns the number of frames processed
func (a *Animator) GetFrameCount() int {
	a.mu.Lock()
	defer a.mu.Unlock()

	return a.frameCount
}

// GetAnimationCount returns the number of active animations
func (a *Animator) GetAnimationCount() int {
	a.mu.Lock()
	defer a.mu.Unlock()

	return len(a.animations)
}

// Clear removes all animations
func (a *Animator) Clear() {
	a.mu.Lock()
	defer a.mu.Unlock()

	a.animations = nil
}

// WaitForCompletion blocks until all animations are complete
func (a *Animator) WaitForCompletion(timeout time.Duration) bool {
	deadline := time.Now().Add(timeout)

	for {
		a.mu.Lock()
		if len(a.animations) == 0 {
			a.mu.Unlock()
			return true
		}
		a.mu.Unlock()

		if time.Now().After(deadline) {
			return false
		}

		time.Sleep(1 * time.Millisecond)
	}
}
