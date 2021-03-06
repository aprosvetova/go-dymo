package dymo

import (
	"github.com/zserge/hid"
	"image"
	"image/color"
	"time"
)

const minWidth = 150
const printerHeight = 64
const margin = 112

const dymoVid uint16 = 0x0922
const dymoPid uint16 = 0x1001

const esc = 0x1B
const syn = 0x16

func FindPrinter() *Printer {
	var p Printer
	var found bool
	hid.UsbWalk(func(d hid.Device) {
		info := d.Info()
		if info.Vendor == dymoVid && info.Product == dymoPid {
			found = true
			p.d = d
		}
	})
	if !found {
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

	err := p.d.Open()
	if err != nil {
		return err
	}

	p.d.Write(data, 3*time.Second)

	p.d.Close()

	return nil
}

type Printer struct {
	d hid.Device
}
