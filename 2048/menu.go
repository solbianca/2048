package twenty48

import (
	"image/color"
	"strconv"

	"github.com/fogleman/gg"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text"
	"golang.org/x/image/font"
)

type menu struct {
	img *ebiten.Image

	title         *title
	newGameButton *newGameButton
	best          *best
	score         *score

	scoreChanel chan DataEvent
}

func newMenu(bestScore int) *menu {
	dc := gg.NewContext(340, 160)
	dc.DrawRoundedRectangle(0, 0, 340, 160, 4)

	img := ebiten.NewImageFromImage(dc.Image())

	m := &menu{
		best:          newBest(bestScore),
		score:         newScore(),
		img:           img,
		newGameButton: newNewGameButton(),
		title:         newTitle(),
	}

	m.scoreChanel = make(chan DataEvent)
	Subscribe("score", m.scoreChanel)

	return m
}

func (m *menu) draw(screen *ebiten.Image) {
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(40, 20)
	screen.DrawImage(m.img, op)

	m.best.draw(screen)
	m.score.draw(screen)
	m.newGameButton.draw(screen)
	m.title.draw(screen)
}

func (m *menu) update() {
	select {
	case res := <-m.scoreChanel:
		scoreValue := res.Data.(int)
		m.score.value += scoreValue
		m.score.items = append(m.score.items, newScoreItem(scoreValue))
	default:
	}

	m.score.update()
	m.newGameButton.update()
}

type best struct {
	x, y          int
	width, height int
	img           *ebiten.Image
	value         int
}

func newBest(bestScore int) *best {
	b := &best{value: bestScore, x: 300, y: 34, width: 60, height: 60}

	dc := gg.NewContext(b.width, b.height)
	dc.DrawRoundedRectangle(0, 0, float64(b.width), float64(b.height), 4)
	dc.SetColor(hexToColor("bbada0", 200))
	dc.Fill()

	b.img = ebiten.NewImageFromImage(dc.Image())

	fontFace, _ := GetFont("mplus", 18)
	str := "best"
	bound, _ := font.BoundString(fontFace, str)
	strWidth := (bound.Max.X - bound.Min.X).Ceil()
	strHeight := (bound.Max.Y - bound.Min.Y).Ceil()
	x := ((b.width - strWidth) / 2)
	y := ((b.height-strHeight)/2 + strHeight) - 16
	text.Draw(b.img, str, fontFace, x, y, color.Black)

	fontFace, _ = GetFont("mplus", 24)
	str = strconv.Itoa(b.value)
	switch {
	case 3 < len(str):
		fontFace, _ = GetFont("mplus", 12)
	case 2 < len(str):
		fontFace, _ = GetFont("mplus", 18)
	}
	bound, _ = font.BoundString(fontFace, str)
	strWidth = (bound.Max.X - bound.Min.X).Ceil()
	strHeight = (bound.Max.Y - bound.Min.Y).Ceil()
	x = ((b.width - strWidth) / 2)
	y = ((b.height-strHeight)/2 + strHeight) + 10
	text.Draw(b.img, str, fontFace, x, y, color.Black)

	return b
}

func (b *best) draw(screen *ebiten.Image) {
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(300, 34)

	screen.DrawImage(b.img, op)
}

type score struct {
	x, y          int
	width, height int
	img           *ebiten.Image
	value         int

	items []*scoreItem
}

func newScore() *score {
	s := &score{value: 0, items: []*scoreItem{}, x: 220, y: 34, width: 60, height: 60}

	dc := gg.NewContext(s.width, s.height)
	dc.DrawRoundedRectangle(0, 0, float64(s.width), float64(s.height), 4)
	dc.SetColor(hexToColor("bbada0", 200))
	dc.Fill()

	s.img = ebiten.NewImageFromImage(dc.Image())

	return s
}

func (s *score) update() {
	changedScoreItems := s.items
	for index, score := range s.items {
		if score.alpha <= 0 {
			changedScoreItems = append(changedScoreItems[:index], changedScoreItems[index+1:]...)
		}
	}

	s.items = changedScoreItems
}

func (s *score) draw(screen *ebiten.Image) {
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(float64(s.x), float64(s.y))

	img := ebiten.NewImageFromImage(s.img)

	fontFace, _ := GetFont("mplus", 18)
	str := "score"
	bound, _ := font.BoundString(fontFace, str)
	strWidth := (bound.Max.X - bound.Min.X).Ceil()
	strHeight := (bound.Max.Y - bound.Min.Y).Ceil()
	x := (s.width - strWidth) / 2
	y := ((s.height-strHeight)/2 + strHeight) - 16
	text.Draw(img, str, fontFace, x, y, color.Black)

	fontFace, _ = GetFont("mplus", 24)
	str = strconv.Itoa(s.value)
	switch {
	case 3 < len(str):
		fontFace, _ = GetFont("mplus", 12)
	case 2 < len(str):
		fontFace, _ = GetFont("mplus", 18)
	}
	bound, _ = font.BoundString(fontFace, str)
	strWidth = (bound.Max.X - bound.Min.X).Ceil()
	strHeight = (bound.Max.Y - bound.Min.Y).Ceil()
	x = (s.width - strWidth) / 2
	y = ((s.height-strHeight)/2 + strHeight) + 10
	text.Draw(img, strconv.Itoa(s.value), fontFace, x, y, color.Black)

	screen.DrawImage(img, op)

	for _, score := range s.items {
		op := &ebiten.DrawImageOptions{}

		x := s.x + 24 + score.offsetX
		y := s.y + 25 + score.offsetY

		op.GeoM.Translate(float64(x), float64(y))
		op.ColorM.Scale(1, 1, 1, score.alpha)

		score.move()

		screen.DrawImage(score.img, op)
	}
}

type scoreItem struct {
	value            int
	img              *ebiten.Image
	alpha            float64
	offsetX, offsetY int
}

func newScoreItem(value int) *scoreItem {
	img := ebiten.NewImage(42, 24)
	valueStr := strconv.Itoa(value)

	fontFace, _ := GetFont("mplus", 24)
	text.Draw(img, valueStr, fontFace, 0, 21, color.Black)

	return &scoreItem{value: value, img: img, alpha: 1}
}

func (s *scoreItem) move() {
	s.alpha -= 0.05
	s.offsetY -= 3
}

type title struct {
	x, y          int
	width, height int
	img           *ebiten.Image
}

func newTitle() *title {
	t := &title{x: 50, y: 34, width: 140, height: 60}
	dc := gg.NewContext(t.width, t.height)
	dc.DrawRoundedRectangle(0, 0, float64(t.width), float64(t.height), 4)
	dc.SetColor(hexToColor("bbada0", 200))
	dc.Fill()

	img := ebiten.NewImageFromImage(dc.Image())
	fontFace, _ := GetFont("mplus", 46)

	str := "2048"
	bound, _ := font.BoundString(fontFace, str)
	strWidth := (bound.Max.X - bound.Min.X).Ceil()
	strHeight := (bound.Max.Y - bound.Min.Y).Ceil()
	x := (t.width - strWidth) / 2
	y := (t.height-strHeight)/2 + strHeight
	text.Draw(img, str, fontFace, x, y, color.Black)

	t.img = img

	return t
}

func (t *title) draw(screen *ebiten.Image) {
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(float64(t.x), float64(t.y))
	screen.DrawImage(t.img, op)
}

type newGameButton struct {
	x, y          int
	width, height int

	img       *ebiten.Image
	isPressed bool
}

func newNewGameButton() *newGameButton {
	btn := &newGameButton{isPressed: false, x: 50, y: 106, width: 140, height: 60}

	dc := gg.NewContext(btn.width, btn.height)
	dc.DrawRoundedRectangle(0, 0, float64(btn.width), float64(btn.height), 4)
	dc.SetColor(hexToColor("bbada0", 200))
	dc.Fill()

	img := ebiten.NewImageFromImage(dc.Image())
	fontFace, _ := GetFont("mplus", 18)
	str := "New Game"
	bound, _ := font.BoundString(fontFace, str)
	strWidth := (bound.Max.X - bound.Min.X).Ceil()
	strHeight := (bound.Max.Y - bound.Min.Y).Ceil()

	x := (btn.width - strWidth) / 2
	y := (btn.height-strHeight)/2 + strHeight

	text.Draw(img, str, fontFace, x, y, color.Black)

	btn.img = img

	return btn
}

func (n *newGameButton) draw(screen *ebiten.Image) {
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(float64(n.x), float64(n.y))
	screen.DrawImage(n.img, op)
}

func (n *newGameButton) update() {
	positionX, positionY := ebiten.CursorPosition()
	x1, y1 := n.x, n.y
	x2, y2 := n.x+n.width, n.y+n.height

	cursorOverButton := positionX > x1 && positionX < x2 && positionY > y1 && positionY < y2
	if !cursorOverButton {
		return
	}

	if ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft) && n.isPressed == false {
		n.press()
	} else if !ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft) && n.isPressed == true {
		n.release()
	}
}

func (n *newGameButton) press() {
	n.isPressed = true
}

func (n *newGameButton) release() {
	n.isPressed = false
	Publish("new_game", nil)
}
