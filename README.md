# the ClanMan

An all-in-one midi device kit powered by Raspberry Pi

Experimental WIP project warning

![WIP](https://github.com/bensinober/clanman/blob/main/docs/clanman.jpg?raw=true "first 3d print")
![WIP](https://github.com/bensinober/clanman/blob/main/docs/20210311_232658.jpg?raw=true "Used an old voice mod board")

Features:
* Raspberry Pi with PatchboxOS for real-time audio and midi setup
* Button panel to handle functions, instrument selections, effects, etc.
* Rotary switch to allow modifying
* OLED display

## Pinout

TODO: draw wiring diagram

OLED Display

    Pin  OLED Function RPi   (pin)
    -------------------
    1  -- GND          GND     (6)
    2  -- 3VV          3VV     (1)
    4  -- DC DataCmd   GPIO 25 (22)
    7  -- SCLK         GPIO 11 (23) # CLK
    8  -- DIN/DATA     GPIO 10 (19) # MOSI
    15 -- CS ChipSel   GPIO 8  (24) # CE0
    16 -- RST          GPIO 24 (18)

CONTROLLER

    led := gpioreg.ByName("GPIO23")  // pin 16

BUTTONS

    btnA := gpioreg.ByName("GPIO17") // pin 11
    btnB := gpioreg.ByName("GPIO18") // pin 12
    btnC := gpioreg.ByName("GPIO27") // pin 15
    btnD := gpioreg.ByName("GPIO22") // pin 13

ROTARY

    apin := gpioreg.ByName("GPIO5")  // pin 29
    bpin := gpioreg.ByName("GPIO6")  // pin 31

## Install

### Raspberry Pi 3+ or newer

Install PatchboxOS on raspberry flash disk

*mindless midi connection service*

    sudo apt install g++ make libasound2-dev
    git clone https://github.com/mzero/amidiminder.git
    cd amidiminder
    make
    sudo dpkg -i build/amidiminder.deb

*disable vnc*

    systemctl stop vncserver-x11-serviced.service
    systemctl disable vncserver-x11-serviced.service

*Enable SPI*

    sed -i "_#dtparam=spi=on_dtparam=spi=on_" /boot/config.txt
    reboot

*midi rules*

    /etc/amidiminder.rules
    amidiminder -C  # checks the rules and then quits
    sudo systemctl restart amidiminder

### Setup

Copy all script files into ./clanman/scripts folder, and the clanman binary to ./clanman

ln -s /home/patch/clanman/scripts/fluidsynth.service /etc/systemd/system/fluidsynth.service
ln -s /home/patch/clanman/scripts/clanman.service /etc/systemd/system/clanman.service

systemctl disable pisound-ctl.service
systemctl enable fluidsynth.service
systemctl enable clanman.service


### Cross-compile tools (HOST)

    sudo apt-get install libc6-armel-cross libc6-dev-armel-cross binutils-arm-linux-gnueabi libncurses5-dev build-essential bison flex libssl-dev bc
    sudo apt install gcc-arm-linux-gnueabihf

BUILD for pi:

    make

copy binary to raspberry and run

## Configuration

* Display

Chose SPI layout for speed, although I2C would be an easier configuration

* Buttons

4 mini pushbuttons with pulldown resistor. Debounced in software


* Rotary

Potentiometer needs some extra work for RaspberyPi since GPIO pins are digital only.
Either by measuring step response in software, or by adding an analog <-> digital signal converter.
