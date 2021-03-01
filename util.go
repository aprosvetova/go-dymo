package dymo

import (
	"image"
	"image/color"
	"image/draw"
)

func convertToColumns(img image.Image) (cols [][]byte) {
	for x := 0; x < img.Bounds().Dx(); x++ {
		col := make([]byte, 8)
		for y := 0; y < printerHeight; y++ {
			bit := 0
			if img.At(x, y).(color.Gray).Y < 127 {
				bit = 1
			}
			col[y>>3] = col[y>>3] | (byte(bit) << (7 - (y & 0x7)))
		}
		cols = append(cols, col)
	}
	for i, j := 0, len(cols)-1; i < j; i, j = i+1, j-1 {
		cols[i], cols[j] = cols[j], cols[i]
	}
	return
}

func convertToMonochrome(in image.Image, c color.Color) image.Image {
	if in.ColorModel() == color.GrayModel {
		return in
	}
	gray := image.NewGray(in.Bounds())
	for y := 0; y < in.Bounds().Dy(); y++ {
		for x := 0; x < in.Bounds().Dx(); x++ {
			ir, ig, ib, _ := in.At(x, y).RGBA()
			r, g, b, _ := c.RGBA()
			if ir == r && ig == g && ib == b {
				gray.Set(x, y, color.Black)
			} else {
				gray.Set(x, y, color.White)
			}
		}
	}
	return gray
}

func fixImageDimensions(in image.Image) image.Image {
	width := in.Bounds().Dx()
	height := in.Bounds().Dy()
	if width > minWidth && height == printerHeight {
		return in
	}

	newWidth := max(width, minWidth)
	x := max((minWidth-width)/2, 0)
	y := max((printerHeight-height)/2, 0)

	img := image.NewGray(image.Rect(0, 0, newWidth, printerHeight))
	draw.Draw(img, img.Bounds(), &image.Uniform{C: color.White}, image.Point{}, draw.Src)
	draw.Draw(img, image.Rect(x, y, newWidth, printerHeight), in, image.Point{}, draw.Src)

	return img
}

func max(x, y int) int {
	if x < y {
		return y
	}
	return x
}
