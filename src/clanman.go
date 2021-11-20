package main

/* the ClanMan
Build for RaspberryPi : GOOS=linux GOARCH=arm GOARM=6 go build clanman.go
*/
import (
	"flag"
	"fmt"
	"log"
	"time"

	"periph.io/x/conn/gpio"
	"periph.io/x/conn/gpio/gpioreg"
	"periph.io/x/conn/gpio/gpiotest"
	"periph.io/x/conn/spi/spireg"
	"periph.io/x/host"
)

type ClanMan struct {
	BtnA        *PushButton    // Menu function select
	BtnB        *PushButton    // Select
	BtnC        *PushButton    // Toggle C
	BtnD        *PushButton    // Toggle D
	Rot         *RotaryEncoder // Toggle Rot
	Led         *Led
	Display     *Display
	Menu        *Menu
	InputEvents chan InputEvent
	Fluid       *FluidSynth
}

type InputEvent struct {
	Name   string `json:"name"`
	Origin string `json:"origin"`
}

type Led struct {
	Pin gpio.PinIO
}

func NewClanMan(ba, bb, bc, bd *PushButton, led *Led, disp *Display, re *RotaryEncoder, m *Menu, ev chan InputEvent, f *FluidSynth) *ClanMan {
	return &ClanMan{
		BtnA:        ba,
		BtnB:        bb,
		BtnC:        bc,
		BtnD:        bd,
		Rot:         re,
		Display:     disp,
		Menu:        m,
		InputEvents: ev,
		Fluid:       f,
	}
}

/* update menu with positional state
   Line 1 - Function
   Line 2 - Select
*/
func (c *ClanMan) UpdateMenu(test bool) {
	fun := c.Menu.GetActiveFunction()
	p := c.Menu.GetcurrentPosition()
	sel := fun.Selects[p[1]]
	var funct, act, tog string
	switch fun.Type {
	case "instrumentSelector":
		if len(c.Fluid.Fonts) > 0 {
			fnt := c.Fluid.Fonts[p[2]] // menu pos 2 = font
			prg := fnt.Banks[0][p[3]]  // menu pos 3 = prg
			funct = fmt.Sprintf("%s %s", fun.Id, sel.Id)
			act = fmt.Sprintf("%d %s", p[2], fnt.Name)
			tog = fmt.Sprintf("%d %s", p[3], prg.Name)
		}
	default:
		if len(sel.Actions) > p[2] {
			act = sel.Actions[p[2]].Id
			if len(sel.Toggles) > p[3] {
				tog = sel.Toggles[p[3]].Id
			}
		}
	}
	fmt.Printf("FUNCT: %s\nSEL: %s\nACT: %s\nTOGG: %s\n", fun.Id, sel.Id, act, tog)
	//fmt.Printf("%#v\n", sel)
	if !test {
		c.Display.Clear()
		c.Display.DrawText(funct, TextTop)
		c.Display.DrawText(act, TextMiddle)
		c.Display.DrawText(tog, TextBottom)
	}
}

func (c *ClanMan) InputHandler(test bool) {
	for ev := range c.InputEvents {
		fmt.Printf("GOT EVENT: %s, BY %s\n", ev.Name, ev.Origin)
		fun := c.Menu.GetActiveFunction()
		switch ev.Origin {
		case "BtnA":
			c.Menu.NextFunction()
		case "BtnB":
			c.Menu.NextSelect()
		case "BtnC":
			if fun.Type == "instrumentSelector" {
				log.Println("Button C fontSelector")
				fontId, progId := c.Fluid.NextFont(c.Menu) // change chan 0 to new font
				c.Menu.SelectFontId(fontId)
				c.Menu.SelectProgId(progId)
			} else {
				c.Menu.NextAction()
			}
		case "BtnD":
			if fun.Type == "instrumentSelector" {
				log.Println("Button D progSelector")
				fontId := c.Fluid.NextInstrumentProg(c.Menu) // change chan 0 to new prog
				c.Menu.SelectProgId(fontId)
			} else {
				c.Menu.NextToggle()
			}
		}
		c.UpdateMenu(test)
	}
}

func main() {
	test := flag.Bool("test", false, "run testing mode")
	port := flag.String("addr", ":1984", "port to run server")
	flag.Parse()
	fmt.Println("Hello ClaNmAn!")

	var clan *ClanMan
	m := NewMenu()
	ev := make(chan InputEvent)
	fmt.Printf("#%v", m)
	//var spiPort spi.PortCloser
	//cc := gpiod.Chips() // to debug GPIO chip
	/* INITIALIZE SPI HOST */
	if *test != true {
		if _, err := host.Init(); err != nil {
			log.Fatal(err)
		}
		spiPort, err := spireg.Open("") // spireg.Open(fmt.Sprintf("/dev/spidev0.%d", index))
		if err != nil {
			log.Fatal(err)
		}
		defer spiPort.Close()
		// OLED
		dc := gpioreg.ByName("GPIO25")  // pin 22
		rst := gpioreg.ByName("GPIO24") // pin 18

		// LED
		lPin := gpioreg.ByName("GPIO23") // pin 16
		lPin.In(gpio.PullDown, gpio.BothEdges)
		led := &Led{Pin: lPin}
		// CONTROL BUTTONS
		btnA := gpioreg.ByName("GPIO17") // pin 11
		btnB := gpioreg.ByName("GPIO18") // pin 12
		btnC := gpioreg.ByName("GPIO27") // pin 13
		btnD := gpioreg.ByName("GPIO22") // pin 15
		ba := NewPushButton(btnA, "BtnA")
		bb := NewPushButton(btnB, "BtnB")
		bc := NewPushButton(btnC, "BtnC")
		bd := NewPushButton(btnD, "BtnD")

		go ba.Listen(ev)
		go bb.Listen(ev)
		go bc.Listen(ev)
		go bd.Listen(ev)
		//go re.Listen()
		// ROTARY
		apin := gpioreg.ByName("GPIO5") // pin 29
		bpin := gpioreg.ByName("GPIO6") // pin 31
		/* initialize */
		disp := NewDisplay(128, 32, spiPort, dc, rst)
		re := NewRotaryEncoder(apin, bpin)
		f := NewFluidSynth("patchbox:9800")
		clan = NewClanMan(ba, bb, bc, bd, led, disp, re, m, ev, f)
		time.Sleep(time.Second * 3)
		clan.Display.Clear()
		time.Sleep(time.Second * 3)
		clan.Display.PrintLogo()
		time.Sleep(time.Second * 3)
		clan.Display.Clear()
	} else {
		/*
			wr := io.WriteCloser(os.Stdout)
			spiPort = spitest.NewRecordRaw(wr)
			// fake spi not working
			edgesFake := make(chan gpio.Level, 0)
			fake := &gpiotest.Pin{N: "GPIO24", EdgesChan: edgesFake}
			clan = NewClanMan(fake, fake, fake, fake, nil, nil, nil, m)
		*/
		edgesFake := make(chan gpio.Level, 0)
		lPin := &gpiotest.Pin{N: "GPIO23", EdgesChan: edgesFake}
		led := &Led{Pin: lPin}
		ev := make(chan InputEvent)
		clan = NewClanMan(nil, nil, nil, nil, led, nil, nil, m, ev, nil)
	}

	go clan.InputHandler(*test)

	s := NewServer(*port, clan)
	if !*test {
		/* Load soundfonts from file and update menu */
		clan.Fluid.LoadFonts(clan.Display)
		clan.Fluid.ResetToInitialFont()
		clan.UpdateMenu(*test)
	}
	s.Run()
}
