package server

import (
	"image"
	"image/draw"

	qrcode "github.com/skip2/go-qrcode"
)

func drawQrcode(dst *image.RGBA, qrcodeParam QrcodeParam, errch chan error) {
	qrcode, err := qrcode.New(qrcodeParam.Content, qrcode.Medium)
	if err != nil {
		errch <- err
		return
	}
	img := qrcode.Image(int(qrcodeParam.Size))
	if err != nil {
		errch <- err
		return
	}
	p := image.Pt(qrcodeParam.X, qrcodeParam.Y)
	draw.Draw(dst, img.Bounds().Add(p), img, image.Pt(0, 0), draw.Over)
	close(errch)
}
