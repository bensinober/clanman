package main

import (
  "fmt"
  "image"
  "log"

  "golang.org/x/image/font"
  "golang.org/x/image/font/basicfont"
  "golang.org/x/image/math/fixed"
  "periph.io/x/conn/gpio"
  "periph.io/x/conn/spi"
  "periph.io/x/devices/ssd1306"
  "periph.io/x/devices/ssd1306/image1bit"
)

// Created using dot2pic.com
// https://javl.github.io/image2cpp/ output vertical 1bit/pixel
var logo = []byte{0x00, 0x00, 0x00, 0x00, 0x80, 0xc0, 0xf0, 0xf8, 0x78, 0x78, 0x38, 0x38, 0x38, 0x38, 0x18, 0x18,
  0x30, 0xf0, 0xc0, 0x00, 0x00, 0x00, 0x00, 0x80, 0x80, 0x80, 0x80, 0x00, 0x00, 0x00, 0x00, 0x00,
  0x00, 0x00, 0x00, 0x00, 0x00, 0xc0, 0x60, 0x70, 0x30, 0x30, 0x30, 0x30, 0xf0, 0xf0, 0x00, 0x00,
  0x00, 0x00, 0x00, 0x00, 0x80, 0x80, 0xc0, 0xc0, 0xc0, 0xc0, 0xc0, 0x80, 0x00, 0x00, 0x00, 0x00,
  0x00, 0x00, 0x80, 0xc0, 0xe0, 0x30, 0x18, 0x18, 0x18, 0x38, 0x78, 0xf8, 0xf8, 0x80, 0x00, 0x00,
  0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x80, 0xc0, 0x40, 0x40, 0x40,
  0xc0, 0xc0, 0xc0, 0x00, 0x00, 0x00, 0xc0, 0x60, 0x20, 0x20, 0x20, 0xe0, 0xe0, 0xe0, 0xc0, 0x80,
  0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
  0x00, 0x00, 0x00, 0x00, 0xff, 0xff, 0x01, 0x01, 0x00, 0xe0, 0x38, 0x0c, 0x0c, 0x06, 0x02, 0x02,
  0x02, 0x03, 0x01, 0x00, 0x00, 0x00, 0xff, 0x83, 0x01, 0x03, 0x7f, 0xff, 0xe0, 0x00, 0x00, 0x00,
  0x00, 0xc0, 0xf8, 0x7e, 0x1f, 0x07, 0xc1, 0x40, 0xc0, 0x00, 0x00, 0x00, 0x3f, 0xff, 0xff, 0xc0,
  0x00, 0xc0, 0xf0, 0x3e, 0x03, 0x00, 0xc0, 0xc0, 0xc0, 0xc3, 0x87, 0x8f, 0x0f, 0x0e, 0x1c, 0x18,
  0xf0, 0x80, 0x7f, 0x3f, 0x00, 0x00, 0xc0, 0xc0, 0xc0, 0x80, 0x00, 0x01, 0x00, 0x07, 0x06, 0x83,
  0x03, 0x03, 0x0f, 0xff, 0xff, 0xfc, 0x00, 0x30, 0x3c, 0x1f, 0x0f, 0x03, 0x00, 0x70, 0x50, 0x70,
  0x00, 0x03, 0x0f, 0x1e, 0x70, 0xc0, 0x7f, 0x1f, 0x00, 0x00, 0xe0, 0x20, 0x20, 0xc1, 0xc3, 0x83,
  0x02, 0x0e, 0x1c, 0x1c, 0x18, 0x38, 0xf0, 0xe0, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
  0x00, 0x00, 0x00, 0x00, 0xff, 0xff, 0x00, 0x00, 0x00, 0x3f, 0x7f, 0x70, 0x60, 0x60, 0x60, 0x20,
  0x20, 0x20, 0xe0, 0xe0, 0xc0, 0x00, 0x1f, 0x3f, 0x78, 0xf0, 0xc0, 0x80, 0x03, 0x0e, 0x1c, 0x18,
  0x16, 0x1f, 0x31, 0xf0, 0x18, 0x08, 0x09, 0x09, 0x19, 0x38, 0xf8, 0xf0, 0x80, 0x80, 0x0f, 0x1f,
  0x3e, 0x7b, 0xe0, 0x20, 0x3e, 0x03, 0x01, 0x01, 0x01, 0x01, 0x07, 0x3f, 0xfc, 0xc0, 0x80, 0x80,
  0x80, 0x83, 0xfe, 0x80, 0xc0, 0x78, 0x0f, 0x07, 0x07, 0x0f, 0x0f, 0x0e, 0x0e, 0x0e, 0x3f, 0xff,
  0xf0, 0xc0, 0x80, 0x00, 0x01, 0x0f, 0xfc, 0x8c, 0x1e, 0x0e, 0x02, 0x02, 0x02, 0x06, 0x1c, 0xf8,
  0xf0, 0xc0, 0x80, 0x80, 0x80, 0x83, 0xfe, 0xfc, 0x00, 0x00, 0x7f, 0xe0, 0x00, 0x00, 0x01, 0x1f,
  0xff, 0xfc, 0xe0, 0x80, 0x00, 0x00, 0x03, 0x7f, 0xe0, 0xc0, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
  0x00, 0x00, 0x00, 0x00, 0x00, 0x03, 0x06, 0x0e, 0x0c, 0x08, 0x08, 0x08, 0x08, 0x08, 0x08, 0x08,
  0x08, 0x0c, 0x0f, 0x03, 0x03, 0x00, 0x00, 0x00, 0x00, 0x01, 0x03, 0x03, 0x07, 0x07, 0x06, 0x04,
  0x04, 0x04, 0x06, 0x03, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x07, 0x07, 0x07, 0x07, 0x06,
  0x06, 0x06, 0x07, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x01, 0x03, 0x03, 0x01,
  0x01, 0x01, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
  0x01, 0x01, 0x01, 0x02, 0x02, 0x02, 0x02, 0x03, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
  0x00, 0x01, 0x01, 0x01, 0x01, 0x01, 0x01, 0x03, 0x02, 0x02, 0x02, 0x03, 0x00, 0x00, 0x00, 0x00,
  0x00, 0x01, 0x03, 0x07, 0x07, 0x06, 0x06, 0x06, 0x03, 0x01, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00}

type Display struct {
  Opts *ssd1306.Opts
  Dev  *ssd1306.Dev
  Img  *image1bit.VerticalLSB
  Text *font.Drawer
}

type TextPlacement *fixed.Point26_6

var (
  TextTop    fixed.Point26_6 = fixed.P(4, 12)
  TextMiddle fixed.Point26_6 = fixed.P(12, 22)
  TextBottom fixed.Point26_6 = fixed.P(12, 32)
)

/* Add display with basic font and monochrome image support */
func NewDisplay(w, h int, port spi.Port, dc, rst gpio.PinIO) *Display {
  //fmt.Println(ssd1306.DefaultOpts)
  rst.In(gpio.PullUp, gpio.BothEdges) // need to pull rst pin high to control display
  opts := ssd1306.Opts{W: w, H: h}    // + rotated etc
  if p, ok := port.(spi.Pins); ok {
    log.Printf("Using pins CLK: %s  MOSI: %s  CS: %s", p.CLK(), p.MOSI(), p.CS())
  }

  //c, err := p.Connect(physic.MegaHertz, spi.Mode3, 8)
  //c, err := port.Connect(4800*physic.MegaHertz, spi.Mode0, 8)
  //if err != nil {
  //  log.Fatal(err)
  //}

  dev, err := ssd1306.NewSPI(port, dc, &opts)
  if err != nil {
    log.Fatalf("failed to initialize ssd1306: %v", err)
  }
  //fmt.Printf("%#+v", dev)
  log.Printf("Display Bounds: %#+v", dev.Bounds())
  img := image1bit.NewVerticalLSB(dev.Bounds())

  face := basicfont.Face7x13
  text := font.Drawer{
    Dst:  img,
    Src:  &image.Uniform{image1bit.On},
    Face: face,
    Dot:  fixed.P(0, 32), // start bottom left
  }
  return &Display{
    Opts: &opts,
    Dev:  dev,
    Img:  img,
    Text: &text,
  }
}

func (d *Display) DrawText(txt string, dot fixed.Point26_6) {
  //point := fixed.Point26_6{fixed.Int26_6(x * 32), fixed.Int26_6(y * 32)}
  d.Text.Dot = dot
  d.Text.DrawString(txt)
  if err := d.Dev.Draw(d.Dev.Bounds(), d.Img, image.Point{}); err != nil {
    log.Println(err)
  }
}

func (d *Display) DrawImg(img image.Image, dot fixed.Point26_6) {
  //point := fixed.Point26_6{fixed.Int26_6(x * 32), fixed.Int26_6(y * 32)}
  d.Text.Dot = dot
  if err := d.Dev.Draw(d.Dev.Bounds(), img, image.Point{}); err != nil {
    log.Println(err)
  }
}

func (d *Display) Clear() {
  fmt.Println("CLEAR")
  c := make([]byte, d.Opts.W*d.Opts.H/8)
  if _, err := d.Dev.Write(c); err != nil {
    log.Println(err)
  }
  if _, err := d.Dev.Write(logo); err != nil {
    fmt.Println(err)
  }
}

func (d *Display) Scroll() {
  if err := d.Dev.Scroll(ssd1306.Left, ssd1306.FrameRate2, 16, -1); err != nil {
    log.Println(err)
  }
  //dev.StopScroll()
}

// http://dot2pic.com/
func (d *Display) PrintLogo() {
  if _, err := d.Dev.Write(logo); err != nil {
    fmt.Println(err)
  }
}
