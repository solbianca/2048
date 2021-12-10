package main

import (
	"log"

	"github.com/hajimehoshi/ebiten/v2"
	twenty48 "twenty48/2048"
)

func main() {
	game, err := twenty48.NewGame()
	if err != nil {
		log.Fatal(err)
	}
	ebiten.SetWindowSize(twenty48.ScreenWidth, twenty48.ScreenHeight)
	ebiten.SetWindowTitle("2048 (Ebiten Demo)")
	if err := ebiten.RunGame(game); err != nil {
		log.Fatal(err)
	}
}
