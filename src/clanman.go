package main

/* the ClanMan
Build for RaspberryPi : GOOS=linux GOARCH=arm GOARM=6 go build clanman.go
*/
import (
	"fmt"
	"log"
	"os"
	"time"

	"periph.io/x/conn/gpio"
	"periph.io/x/conn/gpio/gpioreg"
	"periph.io/x/conn/spi/spireg"
	"periph.io/x/host"
)

type Server struct {
	port int
}

type ClanMan struct {
	BtnA    *PushButton    // Menu function select
	BtnB    *PushButton    // Select
	BtnC    *PushButton    // Toggle C
	BtnD    *PushButton    // Toggle D
	Rot     *RotaryEncoder // Toggle Rot
	Led     *gpio.PinIO
	Display *Display
}

func NewClanMan(ba, bb, bc, bd, led gpio.PinIO, disp *Display, re *RotaryEncoder) *ClanMan {
	a := NewPushButton(ba, "BtnA")
	b := NewPushButton(bb, "BtnB")
	c := NewPushButton(bc, "BtnC")
	d := NewPushButton(bd, "BtnD")
	go a.Listen()
	go b.Listen()
	go c.Listen()
	go d.Listen()
	//go re.Listen()
	led.In(gpio.PullDown, gpio.BothEdges)
	return &ClanMan{
		BtnA:    a,
		BtnB:    b,
		BtnC:    c,
		BtnD:    d,
		Rot:     re,
		Display: disp,
	}
}

func main() {
	fmt.Println("Hello ClanMan!")
	//cc := gpiod.Chips() // to debug GPIO chip

	/* INITIALIZE SPI HOST */
	if _, err := host.Init(); err != nil {
		log.Fatal(err)
	}
	spiPort, err := spireg.Open("") // spireg.Open(fmt.Sprintf("/dev/spidev0.%d", index))
	if err != nil {
		log.Fatal(err)
	}
	defer spiPort.Close()

	/* define GPIO PINS*/
	dc := gpioreg.ByName("GPIO25")  // pin 22
	rst := gpioreg.ByName("GPIO24") // pin 18

	led := gpioreg.ByName("GPIO23") // pin 16

	btnA := gpioreg.ByName("GPIO17") // pin 11
	btnB := gpioreg.ByName("GPIO18") // pin 12
	btnC := gpioreg.ByName("GPIO27") // pin 13
	btnD := gpioreg.ByName("GPIO22") // pin 15

	apin := gpioreg.ByName("GPIO5") // pin 29
	bpin := gpioreg.ByName("GPIO6") // pin 31

	/* initialize */
	disp := NewDisplay(128, 32, spiPort, dc, rst)
	re := NewRotaryEncoder(&apin, &bpin)
	clan := NewClanMan(btnA, btnB, btnC, btnD, led, disp, re)

	clan.Display.DrawText("Top line ---~", TextTop)
	clan.Display.DrawText("Center line --->", TextMiddle)
	clan.Display.DrawText("Bottom line ---^", TextBottom)
	time.Sleep(time.Second * 3)
	clan.Display.Clear()
	time.Sleep(time.Second * 3)
	clan.Display.PrintLogo()
	time.Sleep(time.Second * 3)
	clan.Display.Clear()
	time.Sleep(time.Second * 10)
	os.Exit(0)
}
