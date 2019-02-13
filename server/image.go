package server

import (
	"errors"
	"image"
	"image/color"
	"image/draw"
	"image/gif"
	"image/jpeg"
	"image/png"
	"net/http"
	"strings"

	"github.com/nfnt/resize"
)

type circle struct {
	p image.Point
	r int
}

func (c *circle) ColorModel() color.Model {
	return color.AlphaModel
}

func (c *circle) Bounds() image.Rectangle {
	return image.Rect(c.p.X-c.r, c.p.Y-c.r, c.p.X+c.r, c.p.Y+c.r)
}

func (c *circle) At(x, y int) color.Color {
	xx, yy, rr := float64(x-c.p.X)+0.5, float64(y-c.p.Y)+0.5, float64(c.r)
	if xx*xx+yy*yy < rr*rr {
		return color.Alpha{255}
	}
	return color.Alpha{0}
}

func composeImage(m *Model) (resultImg *image.RGBA, err error) {
	// get base image
	baseImg, err := getImageFromURL(m.BaseImage)
	if err != nil {
		return
	}
	resultImg = image.NewRGBA(baseImg.Bounds())
	drawBg(resultImg, baseImg)

	// draw image
	var errchs = make([]chan error, len(m.Images)+len(m.Texts)+len(m.Qrcodes))
	for n, imageParam := range m.Images {
		errchs[n] = make(chan error)
		go drawImage(resultImg, imageParam, errchs[n])
	}

	// draw text
	for i, textParam := range m.Texts {
		errchs[len(m.Images)+i] = make(chan error)
		go drawText(resultImg, textParam, errchs[len(m.Images)+i])
	}

	// draw qrcode
	for o, qrcodeParam := range m.Qrcodes {
		errchs[len(m.Images)+len(m.Texts)+o] = make(chan error)
		go drawQrcode(resultImg, qrcodeParam, errchs[len(m.Images)+len(m.Texts)+o])
	}

	for _, errch := range errchs {
		if err = <-errch; err != nil {
			return
		}
	}

	return
}

func getImageFromURL(url string) (img image.Image, err error) {
	resp, err := http.Get(url)
	if err != nil {
		return
	}
	contentType := resp.Header.Get("Content-Type")
	switch {
	case strings.Contains(contentType, "jpg") || strings.Contains(contentType, "jpeg"):
		img, err = jpeg.Decode(resp.Body)
	case strings.Contains(contentType, "png"):
		img, err = png.Decode(resp.Body)
	case strings.Contains(contentType, "gif"):
		img, err = gif.Decode(resp.Body)
	default:
		err = errors.New("Get " + url + ": unsupported content type \"" + contentType + "\", it need image")
	}
	defer resp.Body.Close()
	return
}

func drawBg(dst *image.RGBA, bg image.Image) {
	draw.Draw(dst, dst.Bounds(), bg, image.Pt(0, 0), draw.Over)
}

func drawImage(dst *image.RGBA, imgParam ImageParam, errch chan error) {
	img, err := getImageFromURL(imgParam.URL)
	if err != nil {
		errch <- err
		return
	}
	img = resize.Resize(imgParam.Width, imgParam.Height, img, resize.NearestNeighbor)
	p := image.Pt(imgParam.X, imgParam.Y)

	switch imgParam.Clip {
	case 0:
		draw.Draw(dst, img.Bounds().Add(p), img, image.Pt(0, 0), draw.Over)
	case 1:
		draw.DrawMask(dst, img.Bounds().Add(p), img, image.Pt(0, 0), &circle{image.Pt(int(imgParam.Width/2), int(imgParam.Height/2)), int(imgParam.Width / 2)}, image.Pt(0, 0), draw.Over)
	default:
		errch <- errors.New("image: \"" + imgParam.URL + "\" incorrect clip")
		return
	}
	close(errch)
}
