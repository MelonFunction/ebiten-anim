// Package anim provides a simple way to create animations for use with ebiten
package anim

import (
	"image"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
)

// SpriteSheet stores the image and information about the sizing of the SpriteSheet
type SpriteSheet struct {
	Image   *ebiten.Image
	Sprites []*ebiten.Image

	SpriteWidth  int // how big each sprite is
	SpriteHeight int
	SpritesWide  int // how many sprites are in the sheet
	SpritesHigh  int

	Scale float64 // convenience variable for storing scale for use with GeoM
}

// NewSpriteSheet returns a new SpriteSheet
func NewSpriteSheet(img *ebiten.Image, spriteWidth, spriteHeight int, scale float64) *SpriteSheet {
	w, h := img.Size()
	s := &SpriteSheet{
		Image:        img,
		SpriteWidth:  spriteWidth,
		SpriteHeight: spriteHeight,
		SpritesWide:  w / spriteWidth,
		SpritesHigh:  h / spriteHeight,
		Scale:        scale,
	}
	s.Sprites = make([]*ebiten.Image, s.SpritesWide*s.SpritesHigh)
	for x := 0; x < s.SpritesWide; x++ {
		for y := 0; y < s.SpritesHigh; y++ {
			s.Sprites[x+y*s.SpritesWide] = img.SubImage(
				image.Rect(
					x*s.SpriteWidth,
					y*s.SpriteHeight,
					x*s.SpriteWidth+s.SpriteWidth,
					y*s.SpriteHeight+s.SpriteHeight,
				)).(*ebiten.Image)
		}
	}

	return s
}

// GetSprite returns the sprite at the position x,y in the tilesheet
func (s *SpriteSheet) GetSprite(x, y int) *ebiten.Image {
	return s.Sprites[x+y*s.SpritesWide]
}

// Frame stores a single frame of an Animation. It contains an image and how long it should be drawn for
type Frame struct {
	Image    *ebiten.Image
	Duration time.Duration // how long to draw this frame for
}

// NewFrame returns a new Frame
func NewFrame(image *ebiten.Image, duration time.Duration) Frame {
	return Frame{
		Image:    image,
		Duration: duration,
	}
}

// Animation stores a list of Frames and other data regarding timing
type Animation struct {
	Frames        []Frame
	CurrentFrame  int
	LastFrameTime time.Time
	Paused        bool
}

// NewAnimation returns a new Animation
func NewAnimation(frames []Frame) *Animation {
	return &Animation{
		Frames: frames,
		Paused: false,
	}
}

// Update updates
func (a *Animation) Update() {
	if a.Paused {
		return
	}

	now := time.Now()
	if (now.Sub(a.LastFrameTime)) > a.Frames[a.CurrentFrame].Duration {
		a.LastFrameTime = now
		a.CurrentFrame++
		if a.CurrentFrame >= len(a.Frames) {
			a.CurrentFrame = 0
		}
	}
}

// Draw draws the animation to the surface with the provided DrawImageOptions
func (a *Animation) Draw(surface *ebiten.Image, op *ebiten.DrawImageOptions) {
	surface.DrawImage(a.Frames[a.CurrentFrame].Image, op)
}

// Pause pauses the animation
func (a *Animation) Pause() {
	a.Paused = true
}

// Play resumes the animation
func (a *Animation) Play() {
	a.Paused = false
}
