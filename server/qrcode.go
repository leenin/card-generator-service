package server

import (
	"image"
	"image/draw"

	"github.com/boombuler/barcode"
	"github.com/boombuler/barcode/qr"
)

func drawQrcode(dst *image.RGBA, qrcodeParam QrcodeParam, errch chan error) {
	img, err := qr.Encode(qrcodeParam.Content, qr.M, qr.Auto)
	if err != nil {
		errch <- err
		return
	}
	img, err = barcode.Scale(img, int(qrcodeParam.Size), int(qrcodeParam.Size))
	if err != nil {
		errch <- err
		return
	}
	p := image.Pt(qrcodeParam.X, qrcodeParam.Y)
	draw.Draw(dst, img.Bounds().Add(p), img, image.Pt(0, 0), draw.Over)
	close(errch)
}
