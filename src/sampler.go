package main

import (
  "bufio"
  "encoding/json"
  "fmt"
  "io/ioutil"
  "log"
  "net"
  "os"
  "time"
)

/*
  linuxsampler API
  https://download.linuxsampler.org/packages/


  SET VOLUME <float>                                  Set general volume float 0.1...3.0
  CREATE AUDIO_OUTPUT_DEVICE ALSA FRAGMENTSIZE=1024   Create output device
  CREATE MIDI_INPUT_DEVICE ALSA                       Create input device

  *for each channel you want to be able to load an instrument into, add a channel*
  ADD CHANNEL                       Add sampler channel
  LOAD ENGINE gig <id>              Set sampler channel font engine
  ADD CHANNEL MIDI_INPUT <sampler-channel> <midi-device-id>             Set sampler channel midi input device
  SET CHANNEL AUDIO_OUTPUT_DEVICE <sampler-channel> <audio-device-id>   Set sampler channel audio output device

  ex:
  ADD CHANNEL MIDI_INPUT 0 0
  SET CHANNEL AUDIO_OUTPUT_DEVICE 0 0

  *insert instrument directly into channel*
  LOAD INSTRUMENT [NON_MODAL] '<filename>' <instr-index> <sampler-channel>
  ex:
  LOAD INSTRUMENT "/home/pi/GIG/clangig.gig" 0 0
  LOAD INSTRUMENT "/home/pi/SF2/Dance_Organs.sf2" 0 1

  *assign instrument to maps*
  ADD MIDI_INSTRUMENT_MAP [<name>]
  MAP MIDI_INSTRUMENT [NON_MODAL] <map id> <midi_bank> <midi_prog> <engine_name> <filename> <instrument_index> <volume_value> [<instr_load_mode>] [<name>]

  ex:
  ADD MIDI_INSTRUMENT_MAP "danceorgans"
  MAP MIDI_INSTRUMENT 0 0 0 sf2 '/home/pi/SF2/Dance_Organs.sf2' 0 0.9

  SET CHANNEL MIDI_INSTRUMENT_MAP <sampler-channel> <map>

  *test sound on channel*
  SEND CHANNEL MIDI_DATA <midi-msg> <sampler-chan> <arg1> <arg2>
  SEND CHANNEL MIDI_DATA NOTE_ON 0 56 112
  SEND CHANNEL MIDI_DATA NOTE_OFF 0 56 112


*/

type Sampler struct {
  conn   net.Conn
  host   string
  Groups []Group
}

// Font != soundfont, but rather instrument group
type Group struct {
  Id          int
  Instruments []Instrument
  Name        string
}

// for now we just add one channel per engine, as we load instrument on change
// sampler channel 0 = sf2
// sampler channel 1 = gig
type Instrument struct {
  Engine string
  File   string
  Name   string
  Preset int
  Volume float32
}

func NewSampler(host string) *Sampler {
  conn, err := net.Dial("tcp", host)
  if err != nil {
    log.Println(err)
  }
  return &Sampler{
    conn: conn,
    host: host,
  }
}

func (s *Sampler) Init(d *Display) {
  d.Clear()
  d.DrawText("Init sampler...", TextTop)
  cmds := []string{
    "SET VOLUME 1",
    "CREATE AUDIO_OUTPUT_DEVICE ALSA FRAGMENTSIZE=1024",
    "CREATE MIDI_INPUT_DEVICE ALSA",
    // gig
    "ADD CHANNEL",
    "LOAD ENGINE gig 0",
    "ADD CHANNEL MIDI_INPUT 0 0",
    "SET CHANNEL AUDIO_OUTPUT_DEVICE 0 0",
    // sf2
    //"ADD CHANNEL",
    //"LOAD ENGINE sf2 1",
    //"ADD CHANNEL MIDI_INPUT 1 0",
    //"SET CHANNEL AUDIO_OUTPUT_DEVICE 1 0",
  }
  for _, cmd := range cmds {
    s.Send(cmd)
    time.Sleep(time.Millisecond * 100)
    // Need to sleep
    res, err := s.ReceiveLine()
    if err != nil {
      log.Printf("Failed initializing: %s", err)
      return
    }
    log.Println(res)
  }
}

func (s *Sampler) LoadFonts(d *Display) {
  file, err := os.Open("fonts.json")
  if err != nil {
    log.Fatal("Failed opening fonts file")
  }
  defer file.Close()
  bs, _ := ioutil.ReadAll(file)
  var gs []Group
  if err := json.Unmarshal(bs, &gs); err != nil {
    log.Fatalf("Failed parsing font groups file: %s", err)
  }

  for _, g := range gs {
    d.Clear()
    d.DrawText("Loading font...", TextTop)
    d.DrawText(g.Name, TextMiddle)
    time.Sleep(time.Millisecond * 100)
  }
  s.Groups = gs
  fmt.Printf("%+v\n\n", gs)
  d.Clear()
  d.DrawText("Done loading...", TextMiddle)
  time.Sleep(time.Second * 1)

}

func (s *Sampler) Send(msg string) {
  log.Printf("<< %s\n", msg)
  //fmt.Fprintf(s.conn, msg+"\n")
  s.conn.Write([]byte(msg))
  s.conn.Write([]byte("\n"))
}

/* read one line until newline */
func (s *Sampler) ReceiveLine() (string, error) {
  msg, err := bufio.NewReader(s.conn).ReadString('\n')
  if err != nil {
    return "", err
  }
  log.Printf(">> %s\n", msg)
  return msg, nil
}

func (s *Sampler) ResetToInitialFont() {
  g := s.Groups[0]
  inst := g.Instruments[0]
  log.Printf("Setting initial font group %d: %s : %s\n", g.Id, g.Name, inst.Name)
  s.LoadInstrument(inst)
}

func (s *Sampler) LoadInstrument(i Instrument) {
  cmds := []string{
    fmt.Sprintf(`SET CHANNEL VOLUME 0 %.1f`, i.Volume),
    fmt.Sprintf(`LOAD ENGINE %s 0`, i.Engine),
    fmt.Sprintf(`LOAD INSTRUMENT '%s' %d 0`, i.File, i.Preset),
  }
  for _, cmd := range cmds {
    s.Send(cmd)
    time.Sleep(time.Millisecond * 100)
    // Need to sleep
    res, err := s.ReceiveLine()
    if err != nil {
      log.Printf("Failed loading instrument: %s", err)
      return
    }
    log.Println(res)
  }
}

// Group ids are integers, just simply increase until end then start over
func (s *Sampler) NextGroup(m *Menu) (int, int) {
  chanId := m.currentPosition[1]
  groupId := m.currentPosition[2]
  if groupId == len(s.Groups)-1 {
    groupId = 0
  } else {
    groupId++
  }
  group := s.Groups[groupId]
  inst := group.Instruments[0]
  log.Printf("Choosing first instrument from next group: %s for channel: %d\n", group.Name, chanId)
  s.LoadInstrument(inst)
  return groupId, 0 // reset prog Id
}

/* increase instrument preset number, or restart with first */
func (s *Sampler) NextInstrument(m *Menu) int {
  chanId := m.currentPosition[1]
  groupId := m.currentPosition[2]
  instId := m.currentPosition[3]
  group := s.Groups[groupId]
  inst := group.Instruments[instId]
  if instId == len(group.Instruments)-1 {
    instId = 0
  } else {
    instId++
  }
  inst = group.Instruments[instId]
  log.Printf("Choosing next instrument: %s from group: %s for channel: %d\n", inst.Name, group.Name, chanId)
  s.LoadInstrument(inst)
  return instId
}
