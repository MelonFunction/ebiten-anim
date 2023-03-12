// Package anim provides a simple way to create animations for use with ebiten
package anim

import (
	"image"
	"image/color"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
)

// SpriteSheet stores the image and information about the sizing of the SpriteSheet
type SpriteSheet struct {
	Image       *ebiten.Image // original image which was passed on creation
	PaddedImage *ebiten.Image
	Sprites     []*ebiten.Image

	SpriteWidth  int // how big each sprite is
	SpriteHeight int
	SpritesWide  int // how many sprites are in the sheet
	SpritesHigh  int

	Scale float64 // convenience variable for storing scale for use with GeoM
}

// NewSpriteSheet returns a new SpriteSheet
// Use padding==true if you're creating a tilesheet
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

	p := 2
	paddedImg := ebiten.NewImage(w+(s.SpritesWide+1)*p, h+(s.SpritesHigh+1)*p)
	eraser := ebiten.NewImage(spriteWidth, spriteHeight)
	eraser.Fill(color.RGBA{255, 255, 255, 255})
	// paddedImg.Fill(color.RGBA{255, 0, 255, 255})

	s.Sprites = make([]*ebiten.Image, s.SpritesWide*s.SpritesHigh)
	for x := 0; x < s.SpritesWide; x++ {
		for y := 0; y < s.SpritesHigh; y++ {
			dx := float64(spriteWidth)*float64(x) + float64(p)*(float64(x)+1)
			dy := float64(spriteHeight)*float64(y) + float64(p)*(float64(y)+1)

			// draw padding first
			d := func(op *ebiten.DrawImageOptions) {
				paddedImg.DrawImage(img.SubImage(
					image.Rect(
						x*s.SpriteWidth,
						y*s.SpriteHeight,
						(x+1)*s.SpriteWidth,
						(y+1)*s.SpriteHeight,
					)).(*ebiten.Image), op)
			}
			for zx := -p / 2; zx <= p/2; zx++ {
				if zx != 0 {
					op := &ebiten.DrawImageOptions{}
					op.GeoM.Translate(dx+float64(zx), dy)
					d(op)
				}
			}
			for zy := -p / 2; zy <= p/2; zy++ {
				if zy != 0 {
					op := &ebiten.DrawImageOptions{}
					op.GeoM.Translate(dx, dy+float64(zy))
					d(op)
				}
			}

			// clear area, if a tile isn't full width, it'll be the wrong size (2px will be increased to 4px wide!)
			op := &ebiten.DrawImageOptions{}
			op.GeoM.Translate(dx, dy)
			op.CompositeMode = ebiten.CompositeModeClear
			paddedImg.DrawImage(eraser, op)

			// draw the sprite itself
			op = &ebiten.DrawImageOptions{}
			op.GeoM.Translate(
				dx, dy)
			paddedImg.DrawImage(img.SubImage(
				image.Rect(
					x*s.SpriteWidth,
					y*s.SpriteHeight,
					(x+1)*s.SpriteWidth,
					(y+1)*s.SpriteHeight,
				)).(*ebiten.Image), op)

			// save subimage/reference
			s.Sprites[x+y*s.SpritesWide] = paddedImg.SubImage(
				image.Rect(
					int(dx),
					int(dy),
					int(dx)+spriteWidth,
					int(dy)+spriteHeight,
				)).(*ebiten.Image)
		}
	}

	s.PaddedImage = paddedImg

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
