package twenty48

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
)

const (
	ScreenWidth  = 420
	ScreenHeight = 600
	boardSize    = 4
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

type game struct {
	scene *scene

	newGameChanel chan DataEvent
}

func NewGame() (*game, error) {
	g := &game{}
	scene, err := newScene(0)
	if err != nil {
		return nil, err
	}
	g.scene = scene

	g.newGameChanel = make(chan DataEvent)
	Subscribe("new_game", g.newGameChanel)

	return g, nil
}

func (g *game) Update() error {
	g.scene.update()

	select {
	case res := <-g.newGameChanel:
		fmt.Println("Response from `new_game`", res)
		bestScore := g.scene.menu.best.value
		if g.scene.menu.score.value > bestScore {
			bestScore = g.scene.menu.score.value
		}

		newScene, _ := newScene(bestScore)
		newScene.menu.scoreChanel = g.scene.menu.scoreChanel
		g.scene = newScene
	default:
	}

	return nil
}

func (g *game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return ScreenWidth, ScreenHeight
}

func (g *game) Draw(screen *ebiten.Image) {
	g.scene.draw(screen)
}

type scene struct {
	input *Input
	menu  *menu
	board *board
}

func newScene(bestScore int) (*scene, error) {
	s := &scene{
		input: NewInput(),
		menu:  newMenu(bestScore),
	}

	var err error
	s.board, err = NewBoard(boardSize)
	if err != nil {
		return nil, err
	}

	return s, nil
}

func (s *scene) update() {
	s.input.Update()
	s.menu.update()

	if err := s.board.Update(s.input); err != nil {
		return
	}
}

func (s *scene) draw(screen *ebiten.Image) {
	screen.Fill(backgroundColor)

	s.menu.draw(screen)

	w, h := s.board.Size()
	boardImage := ebiten.NewImage(w, h)
	s.board.Draw(boardImage)
	op := &ebiten.DrawImageOptions{}
	x := 40
	y := 220
	op.GeoM.Translate(float64(x), float64(y))
	screen.DrawImage(boardImage, op)
}
