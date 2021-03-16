package main

import (
	"fmt"
	"log"
	"time"

	"periph.io/x/conn/gpio"
	"periph.io/x/conn/gpio/gpioutil"
)

type PushButton struct {
	Pin  gpio.PinIO
	Name string
}

func NewPushButton(p gpio.PinIO, n string) *PushButton {
	p.In(gpio.PullUp, gpio.BothEdges) // pull high before usage, listen for both edges
	d, err := gpioutil.Debounce(p, 30*time.Millisecond, 30*time.Millisecond, gpio.BothEdges)
	if err != nil {
		log.Fatal(err)
	}
	return &PushButton{
		Pin:  d,
		Name: n,
	}
}

func (b *PushButton) Listen(ev chan InputEvent) {
	for {
		b.Pin.WaitForEdge(-1)
		//fmt.Println(b.Pin.Read())
		if b.Pin.Read() == gpio.High {
			p := InputEvent{Name: "PRESS", Origin: b.Name}
			ev <- p
		}
	}
}

/* rotary encoder
https://gist.github.com/toxygene/6ee54127aa1c133574da2f0d9bb0e8c2
*/
type RotaryEncoder struct {
	PinA gpio.PinIO
	PinB gpio.PinIO
	//previousEncoderState uint8
	//m                    *sync.Mutex
}

func NewRotaryEncoder(pinA gpio.PinIO, pinB gpio.PinIO) *RotaryEncoder {
	pinA.In(gpio.PullUp, gpio.BothEdges) // pull pin up before reading
	pinB.In(gpio.PullUp, gpio.BothEdges) // pull pin up before reading

	return &RotaryEncoder{
		PinA: pinA,
		PinB: pinB,
		//m:    &sync.Mutex{},
	}
}

func (re *RotaryEncoder) Read() int {
	// discharge first for 5ms
	re.PinA.In(gpio.PullNoChange, gpio.NoEdge)
	re.PinB.Out(gpio.Low)

	//    GPIO.output(b_pin, False)
	time.Sleep(time.Millisecond * 5)
	// then measure time
	c := 0
	re.PinB.In(gpio.PullNoChange, gpio.NoEdge)
	re.PinA.Out(gpio.High)
	for {
		if re.PinB.Read() == gpio.Low {
			c++
			continue
		}
		break
	}
	return c
}

/*
func (re *RotaryEncoder) Read() int {
	// discharge first for 5ms
	(*re.PinA).In(gpio.PullNoChange, gpio.NoEdge)
	(*re.PinB).Out(gpio.Low)

	//    GPIO.output(b_pin, False)
	time.Sleep(time.Millisecond * 5)
	// then measure time
	c := 0
	(*re.PinB).In(gpio.PullNoChange, gpio.NoEdge)
	(*re.PinA).Out(gpio.High)
	for {
		if (*re.PinB).Read() == gpio.Low {
			c++
			continue
		}
		break
	}
	return c
}
*/
func (re *RotaryEncoder) Listen() {
	for {
		r := re.Read()
		fmt.Println(r)
	}
}
