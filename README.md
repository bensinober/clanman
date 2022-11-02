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

### Raspberry Pi 3+ or newer with Raspbian OS

*linuxsampler*

    linuxsampler is an API for alsa midi that supports sfz, gig and sf2 sound banks

    sudo apt install libsndfile-dev libaudiofile-dev

    wget https://download.linuxsampler.org/packages/libgig-4.3.0.tar.bz2
    tar xvjf libgig-4.3.0.tar.bz2
    cd libgig-4.3.0
    ./configure
    make
    make install

    wget https://download.linuxsampler.org/packages/linuxsampler-2.2.0.tar.bz2
    tar xvjf linuxsampler-2.2.0.tar.bz2
    cd linuxsampler-2.2.0
    ./configure
    make
    make install


    alternatives: install from deb packages: https://download.linuxsampler.org/packages/debian/

*fluidsynth et al*

    sudo apt install fluidsynth

*mindless midi connection service*

    sudo apt install git g++ make libasound2-dev
    git clone https://github.com/mzero/amidiminder.git
    cd amidiminder
    make
    sudo dpkg -i build/amidiminder.deb

*Enable SPI*

    sed -i "_#dtparam=spi=on_dtparam=spi=on_" /boot/config.txt
    reboot

*midi rules*

    /etc/amidiminder.rules
    amidiminder -C  # checks the rules and then quits
    sudo systemctl restart amidiminder

### Setup

Copy all script files into ./clanman/scripts folder, and the clanman binary to ./clanman

    ln -s /home/pi/clanman/scripts/fluidsynth.service /etc/systemd/system/fluidsynth.service
    ln -s /home/pi/clanman/scripts/linuxsampler.service /etc/systemd/system/linuxsampler.service
    ln -s /home/pi/clanman/scripts/clanman.service /etc/systemd/system/clanman.service

    systemctl enable fluidsynth.service
    systemctl enable linuxsampler.service
    systemctl enable clanman.service


### Cross-compile tools (on HOST)

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
