package main

import (
	"github.com/boombuler/hid"
	"image"
	"image/color"
)

const minWidth = 150
const printerHeight = 64
const margin = 112

const dymoVid uint16 = 0x0922
const dymoPid uint16 = 0x1001

const esc = 0x1B
const syn = 0x16

func FindPrinter() *Printer {
	ch := hid.FindDevices(dymoVid, dymoPid)
	var p Printer
	for d := range ch {
		p.d = d
	}
	if p.d == nil {
		return nil
	}
	return &p
}

func (p *Printer) Print(img image.Image, c color.Color) error {
	img = convertToMonochrome(img, c)
	img = fixImageDimensions(img)

	cols := convertToColumns(img)
	for i := 0; i < margin; i++ {
		cols = append(cols, make([]byte, 8))
	}

	data := []byte{esc, 'C', 0, esc, 'B', 0, esc, 'D', 8}
	for _, col := range cols {
		data = append(data, syn)
		data = append(data, col...)
	}

	dev, err := p.d.Open()
	if err != nil {
		return err
	}

	dev.Write(data)

	dev.Close()

	return nil
}

type Printer struct {
	d *hid.DeviceInfo
}
