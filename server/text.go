package server

import (
	"errors"
	"image"
	"image/color"
	"io/ioutil"
	"math"
	"path/filepath"

	"golang.org/x/image/math/fixed"

	"github.com/golang/freetype"
	"github.com/golang/freetype/truetype"
)

func drawText(dst *image.RGBA, textParam TextParam, errch chan error) {
	fontPath, err := filepath.Abs("./server/font/" + textParam.Font + ".ttf")
	if err != nil {
		errch <- errors.New("Text " + textParam.Content + ", " + textParam.Font + " font not exists")
		return
	}

	fontTypes, err := ioutil.ReadFile(fontPath)
	if err != nil {
		errch <- err
		return
	}

	font, err := freetype.ParseFont(fontTypes)
	if err != nil {
		errch <- err
		return
	}

	// type context
	c := freetype.NewContext()
	c.SetClip(dst.Bounds())
	c.SetDPI(72)
	c.SetFont(font)
	c.SetFontSize(textParam.Size)
	c.SetDst(dst)

	// text color
	var opacity byte
	if len(textParam.Color) == 4 {
		opacity = byte(math.Floor(float64(textParam.Color[3])*255/100 + 0.5))
	} else {
		opacity = 255
	}
	c.SetSrc(image.NewUniform(&color.NRGBA{textParam.Color[0], textParam.Color[1], textParam.Color[2], opacity}))

	var x fixed.Int26_6
	switch textParam.Anchor {
	case 0:
		x = fixed.I(textParam.X)
	case 1:
		// Truetype stuff
		opts := &truetype.Options{}
		opts.Size = textParam.Size
		face := truetype.NewFace(font, opts)

		// calculate the width
		var textWidth fixed.Int26_6
		for _, x := range textParam.Content {
			width, ok := face.GlyphAdvance(rune(x))
			if ok != true {
				errch <- errors.New("the face does not contain a glyph for " + string(x))
				return
			}
			textWidth += width
		}
		x = fixed.I(textParam.X) - (textWidth / 2)
	default:
		errch <- errors.New("Text " + textParam.Content + ", anchor incorrect" + "[" + string(textParam.Anchor) + "]")
		return
	}

	y := fixed.I(textParam.Y) + c.PointToFixed(textParam.Size)
	pt := fixed.Point26_6{X: x, Y: y}

	_, err = c.DrawString(textParam.Content, pt)
	if err != nil {
		errch <- err
		return
	}

	close(errch)
}
