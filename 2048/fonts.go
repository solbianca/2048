package twenty48

import (
	"fmt"
	"io/ioutil"
	"strconv"

	"github.com/golang/freetype/truetype"
	"golang.org/x/image/font"
)

const dpi = 72

var (
	fonts           *fontStorage
	mplusSmallFont  font.Face
	mplusNormalFont font.Face
	mplusBigFont    font.Face
)

func init() {
	fonts = newFontStorage()

	LoadFont("mplus", "resources/font/mplus-regular.ttf")
	LoadFont("arial", "resources/font/arial-bold.ttf")

	mplusSmallFont, _ = GetFont("mplus", 24)
	mplusNormalFont, _ = GetFont("mplus", 32)
	mplusBigFont, _ = GetFont("mplus", 48)
}

type fontStorage struct {
	ttf   map[string]*truetype.Font
	faces map[string]font.Face
}

func newFontStorage() *fontStorage {
	return &fontStorage{ttf: map[string]*truetype.Font{}, faces: map[string]font.Face{}}
}

func GetFont(alias string, size int) (font.Face, error) {
	ttfFont, ok := fonts.ttf[alias]

	if !ok {
		return nil, fmt.Errorf("font by alias [%s] not found", alias)
	}

	fontKey := fmt.Sprintf("%s:%s", alias, strconv.Itoa(size))
	fontFace, ok := fonts.faces[fontKey]

	if !ok {
		fontFace = truetype.NewFace(
			ttfFont, &truetype.Options{
				Size:    float64(size),
				DPI:     dpi,
				Hinting: font.HintingFull,
			},
		)
		fonts.faces[fontKey] = fontFace
	}

	return fontFace, nil
}

func LoadFont(alias, path string) {
	if _, ok := fonts.ttf[alias]; ok {
		panic(fmt.Errorf("font by alias [%s] already loaded", alias))
	}

	fontData := loadFileOrPanic(path)

	ttfFont, err := truetype.Parse(fontData)
	if err != nil {
		panic(fmt.Errorf("font not parsed by path [%a], %v", path, err))
	}

	fonts.ttf[alias] = ttfFont
}

func loadFileOrPanic(path string) []byte {
	file, err := ioutil.ReadFile(path)
	if err != nil {
		panic(fmt.Errorf("file not founded by path [%a], %v", path, err))
	}

	return file
}
