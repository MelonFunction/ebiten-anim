# ebiten-anim

Create a SpriteSheet to easily select parts of an image and create animations from them.


## Usage

[📖 Docs](https://pkg.go.dev/github.com/melonfunction/ebiten-anim)  
Look at [the example](https://github.com/melonfunction/ebiten-anim/tree/master/examples) to see how to use the library.

```go
// Create a spritesheet
SpriteSheet = anim.NewSpriteSheet(sprites, 8, 10, anim.SpriteSheetOptions{
    Scale: 16,
})

// Create an animation
frames := make([]anim.Frame, 5)
for x := 0; x <= 4; x++ {
    frames[x] = anim.NewFrame(SpriteSheet.GetSprite(x, 0), time.Second / 20)
}
Animation := anim.NewAnimation(frames)

// Update the animation
Animation.Update()

// Pause/play the animation
Animation.Pause()
Animation.Play()

// Draw the animation
Animation.Draw(screen, &ebiten.DrawImageOptions{})

// You can also just draw frames from the spritesheet without using animations
// This might be useful to avoid cluttering your Draw function with Image.SubImage
screen.DrawImage(SpriteSheet.GetSprite(0, 0), &ebiten.DrawImageOptions{})

```