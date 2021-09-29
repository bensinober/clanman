package main

import (
  "fmt"
  "log"
  "net"
  "strconv"
)

/*
fluidsynth API
    noteon chan key vel        Send noteon
    noteoff chan key           Send noteoff
    pitch_bend chan offset     Bend pitch
    pitch_bend chan range      Set bend pitch range
    cc chan ctrl value         Send control-change message
    prog chan num              Send program-change message
    select chan sfont bank prog  Combination of bank-select and program-change
    load file [reset] [bankofs] Load SoundFont (reset=0|1, def 1; bankofs=n, def 0)
    unload id [reset]          Unload SoundFont by ID (reset=0|1, default 1)
    reload id                  Reload the SoundFont with the specified ID
    fonts                      Display the list of loaded SoundFonts
    inst font                  Print out the available instruments for the font
    channels [-verbose]        Print out preset of all channels
    interp num                 Choose interpolation method for all channels
    interpc chan num           Choose interpolation method for one channel
    rev_preset num             Load preset num into the reverb unit
    rev_setroomsize num        Change reverb room size
    rev_setdamp num            Change reverb damping
    rev_setwidth num           Change reverb width
    rev_setlevel num           Change reverb level
    reverb [0|1|on|off]        Turn the reverb on or off
    cho_set_nr n               Use n delay lines (default 3)
    cho_set_level num          Set output level of each chorus line to num
    cho_set_speed num          Set mod speed of chorus to num (Hz)
    cho_set_depth num          Set chorus modulation depth to num (ms)
    chorus [0|1|on|off]        Turn the chorus on or off
    gain value                 Set the master gain (0 < gain < 5)
    voice_count                Get number of active synthesis voices
    tuning name bank prog      Create a tuning with name, bank number,
                               and program number (0 <= bank,prog <= 127)
    tune bank prog key pitch   Tune a key
    settuning chan bank prog   Set the tuning for a MIDI channel
    resettuning chan           Restore the default tuning of a MIDI channel
    tunings                    Print the list of available tunings
    dumptuning bank prog       Print the pitch details of the tuning
    reset                      System reset (all notes off, reset controllers)
    set name value             Set the value of a controller or settings
    get name                   Get the value of a controller or settings
    info name                  Get information about a controller or settings
    settings                   Print out all settings
    echo arg                   Print arg
    ladspa_clear               Resets LADSPA effect unit to bypass state
    ladspa_add lib plugin n1 <- p1 n2 -> p2 ... Loads and connects LADSPA plugin
    ladspa_start               Starts LADSPA effect unit
    ladspa_declnode node value Declares control node `node' with value `value'
    ladspa_setnode node value  Assigns `value' to `node'
    router_clear               Clears all routing rules from the midi router
    router_default             Resets the midi router to default state
    router_begin [note|cc|prog|pbend|cpress|kpress]: Starts a new routing rule
    router_chan min max mul add      filters and maps midi channels on current rule
    router_par1 min max mul add      filters and maps parameter 1 (key/ctrl nr)
    router_par2 min max mul add      filters and maps parameter 2 (vel/cc val)
    router_end                 closes and commits the current routing rule


*/
type FluidSynth struct {
  conn net.Conn
  host string
}

func NewFluidSynth(host string) *FluidSynth {
  conn, err := net.Dial("tcp", host)
  if err != nil {
    log.Println(err)
  }
  return &FluidSynth{
    conn: conn,
    host: host,
  }
}

func (f FluidSynth) LoadFonts(ss []string) {
  for _, s := range ss {
    l := fmt.Sprintf("load /home/patch/SF2/%s", s)
    f.Send(l)
  }
}

func (f FluidSynth) Send(msg string) {
  f.conn.Write([]byte(msg))
  f.conn.Write([]byte("\n"))
  log.Printf("Send: %s", msg)

  buff := make([]byte, 1024)
  n, _ := f.conn.Read(buff)
  log.Printf("Receive: %s", buff[:n])
}

func (f FluidSynth) SetSampleRate(rate int) {
  msg := fmt.Sprintf("set synth.sample-rate %s", rate)
  f.Send(msg)
}

func (f FluidSynth) PutFontInFront(font int) {
  r := strconv.Itoa(rate)
  // select chan sfont bank prog
  // Just put the first prog of the first bank on channel 0
  msg := fmt.Sprintf("select 0 %s 0 0", font)
  f.Send(msg)
}
