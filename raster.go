package braster

import (
	"image"
	"image/color"
	"io"
)

type RasterEncoder struct {
	img    image.Image
	bounds image.Rectangle
	line   int
	dots   int
	buffer []byte

	dpi                                     DPI
	marginLeftDots, marginRightDots         int
	userMarginLeftDots, userMarginRightDots int
}
type RasterEncodeConfig struct {
	Model     ModelCode
	MediaInfo *MediaInfo
}

func NewRasterEncoder(img image.Image, config *RasterEncodeConfig) *RasterEncoder {
	re := RasterEncoder{
		img:    img,
		bounds: img.Bounds(),
		dpi:    DPI(180),
	}
	re.line = re.bounds.Dx()

	re.dots = re.bounds.Dy()
	if config != nil {
		re.dots = config.Model.TotalDots()
		re.dpi = config.Model.DPI()

		if config.MediaInfo != nil {
			margin := re.dpi.MillimetreToDots(config.MediaInfo.PageMargin)
			re.marginLeftDots = margin
			re.marginRightDots = margin
		}
	}

	bufSize := re.dots / 8
	if re.dots%8 > 0 {
		bufSize++
	}
	re.buffer = make([]byte, bufSize)

	return &re
}

func (re *RasterEncoder) SetMargin(left, right float32) {
	re.userMarginLeftDots = re.dpi.MillimetreToDots(left)
	re.userMarginRightDots = re.dpi.MillimetreToDots(right)
}

func (re *RasterEncoder) Next() bool {
	if re.line > 0 {
		re.line--
		return true
	}
	return false
}
func (re *RasterEncoder) EncodeLine(b []byte) []byte {
	return re.encode(re.line)
}
func (re *RasterEncoder) encode(x int) []byte {
	b := re.buffer
	if x < 0 || x >= re.bounds.Dx() {
		for i := range b {
			b[i] = 0
		}
		return b
	}

	var pmin = re.marginLeftDots + re.userMarginLeftDots
	var pmax = re.dots - (re.marginRightDots + re.userMarginRightDots)

	var nbyte = -1
	var bit byte
	var y = re.bounds.Min.Y
	for i := 0; i < re.dots; i++ {
		bit = bit << 1
		if bit == 0 {
			nbyte++
			bit = 1
			b[nbyte] = 0
		}

		if i < pmin || i >= pmax {
			continue
		}

		if y > re.bounds.Max.Y {
			continue
		}
		if color.GrayModel.Convert(re.img.At(x+re.bounds.Min.X, y)).(color.Gray).Y > 127 {
			b[nbyte] |= byte(bit)
		}
	}
	return b
}
func (re *RasterEncoder) Transfer(w io.Writer) error {
	cb := CommandBuilder{w}
	dx := re.bounds.Dx()
	for x := 0; x < dx; x++ {
		err := cb.RasterGraphicsTransfer(re.encode(x))
		if err != nil {
			return err
		}
	}
	return nil
}
