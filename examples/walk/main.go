// Package main 👍
package main

import (
	"errors"
	"log"

	"github.com/hajimehoshi/ebiten/v2"
)

// vars
var (
	WindowWidth  = 640 * 2
	WindowHeight = 480 * 2

	ErrNormalExit = errors.New("Normal exit")
)

// Game implements ebiten.Game interface.
type Game struct{}

// Update proceeds the game state.
// Update is called every tick (1/60 [s] by default).
func (g *Game) Update() error {
	if ebiten.IsKeyPressed(ebiten.KeyEscape) {
		return ErrNormalExit
	}

	return nil
}

// Draw draws the game screen.
// Draw is called every frame (typically 1/60[s] for 60Hz display).
func (g *Game) Draw(screen *ebiten.Image) {
}

// Layout sets window size
func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	if WindowWidth != outsideWidth || WindowHeight != outsideHeight {
		log.Println("resize", outsideWidth, outsideHeight)
		WindowWidth = outsideWidth
		WindowHeight = outsideHeight
	}
	return outsideWidth, outsideHeight
}

func main() {
	game := &Game{}
	ebiten.SetWindowSize(WindowWidth, WindowHeight)
	ebiten.SetWindowTitle("Animation example")
	ebiten.SetWindowResizable(true)

	if err := ebiten.RunGame(game); err != nil {
		log.Fatal(err)
	}
}
