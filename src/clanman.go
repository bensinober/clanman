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
}

type InputEvent struct {
	Name   string `json:"name"`
	Origin string `json:"origin"`
}

type Led struct {
	Pin gpio.PinIO
}

func NewClanMan(ba, bb, bc, bd *PushButton, led *Led, disp *Display, re *RotaryEncoder, m *Menu, ev chan InputEvent) *ClanMan {
	return &ClanMan{
		BtnA:        ba,
		BtnB:        bb,
		BtnC:        bc,
		BtnD:        bd,
		Rot:         re,
		Display:     disp,
		Menu:        m,
		InputEvents: ev,
	}
}

func (c *ClanMan) UpdateMenu(test bool) {
	p := c.Menu.currentPosition
	fun := c.Menu.Functions[p[0]]
	sel := fun.Selects[p[1]]
	// TODO: error checking here if json is not complete
	act := sel.Actions[p[2]]
	//fmt.Printf("%#v\n", sel)
	fmt.Printf("LINE1: %s\nLINE2: %s\nLINE3: %s\n", fun.Id, sel.Id, act.Id)
	if !test {
		c.Display.Clear()
		c.Display.DrawText(fun.Id, TextTop)
		c.Display.DrawText(sel.Id, TextMiddle)
		c.Display.DrawText(act.Id, TextBottom)
	}
}

func (c *ClanMan) InputHandler(test bool) {
	for ev := range c.InputEvents {
		fmt.Printf("GOT EVENT: %s, BY %s\n", ev.Name, ev.Origin)
		switch ev.Origin {
		case "BtnA":
			c.Menu.ToggleFunction()
		case "BtnB":
			c.Menu.ToggleSelect()
		case "BtnC":
			c.Menu.ToggleAction()
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
		clan = NewClanMan(ba, bb, bc, bd, led, disp, re, m, ev)
		time.Sleep(time.Second * 3)
		clan.Display.Clear()
		time.Sleep(time.Second * 3)
		clan.Display.PrintLogo()
		time.Sleep(time.Second * 3)
		clan.Display.Clear()
		clan.UpdateMenu(*test)
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
		clan = NewClanMan(nil, nil, nil, nil, led, nil, nil, m, ev)
	}

	go clan.InputHandler(*test)
	s := NewServer(*port, clan)
	s.Run()
}
