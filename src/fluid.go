package main

import (
  "bufio"
  "encoding/json"
  "fmt"
  "io/ioutil"
  "log"
  "net"
  "os"
  "strconv"
  "strings"
  "time"
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
const soundFontResultString = "loaded SoundFont has ID "

type FluidSynth struct {
  conn  net.Conn
  host  string
  Fonts []Font
}

type Font struct {
  Id    int
  Name  string
  File  string
  Banks map[int][]Prog
}

type Prog struct {
  Id   int
  Name string
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

func (f *FluidSynth) LoadFonts(d *Display) {
  fonts := make([]Font, 0)
  file, err := os.Open("fonts.json")
  if err != nil {
    log.Fatal("Failed opening fonts file")
  }
  defer file.Close()
  bs, _ := ioutil.ReadAll(file)
  fs := []struct{ Name, File string }{} // intermediate struct for json font input file
  if err := json.Unmarshal(bs, &fs); err != nil {
    log.Fatalf("Failed parsing fonts: %s", err)
  }

  for _, font := range fs {
    d.Clear()
    d.DrawText("Loading font...", TextTop)
    d.DrawText(font.Name, TextMiddle)
    l := fmt.Sprintf("load /home/patch/SF2/%s", font.File)
    f.Send(l)
    time.Sleep(time.Millisecond * 3000)
    // Need to sleep
    res, err := f.ReceiveLines()
    if err != nil {
      log.Printf("Failed loading font: %s", err)
      return
    }
    if strings.HasPrefix(res, soundFontResultString) {
      idStr := strings.TrimSpace(strings.Replace(res, soundFontResultString, "", -1))
      log.Printf("loaded font got id: %s", idStr)
      id, err := strconv.Atoi(idStr)
      if err != nil {
        log.Println(err)
        continue
      }
      bnks := f.ParseFontBanks(id)
      fonts = append(fonts, Font{
        Id:    id,
        Name:  font.Name,
        File:  font.File, // dont need this?
        Banks: bnks,
      })
    } else {
      log.Printf("failed getting font id from res: %s", res)
    }
  }
  f.Fonts = fonts
  d.Clear()
  d.DrawText("Done loading...", TextMiddle)
  time.Sleep(time.Second * 1)

}

/* parse instruments from font
   eg 000-012 rock organ
   disabled as loading is async in fluidsynth and we have no callback to ensure order
*/
func (f FluidSynth) ParseFontBanks(id int) map[int][]Prog {
  msg := fmt.Sprintf("inst %d", id)
  banks := make(map[int][]Prog, 0) // intermediary map of bank ids
  f.Send(msg)
  res, err := f.ReceiveLines()
  if err != nil {
    log.Printf("Error: %s", err)
    return banks
  }
  log.Println(res)

  /* scanner splits by newlines */
  scanner := bufio.NewScanner(strings.NewReader(res))
  for scanner.Scan() {
    ln := scanner.Text()
    bnk, _ := strconv.Atoi(ln[0:3])
    prg, _ := strconv.Atoi(ln[4:7])
    nam := ln[8:len(ln)]
    banks[bnk] = append(banks[bnk], Prog{prg, nam})
  }
  return banks
}

func (f *FluidSynth) Send(msg string) {
  log.Printf("<< %s", msg)
  //fmt.Fprintf(f.conn, msg+"\n")
  f.conn.Write([]byte(msg))
  f.conn.Write([]byte("\n"))
}

/* read one line until newline */
func (f *FluidSynth) ReceiveLine() (string, error) {
  msg, err := bufio.NewReader(f.conn).ReadString('\n')
  if err != nil {
    return "", err
  }
  log.Printf(">> %s", msg)
  return msg, nil
}

/* read one or more lines until no more is received */
func (f *FluidSynth) ReceiveLines() (string, error) {
  scanner := bufio.NewScanner(f.conn)
  var out string
  for {
    f.conn.SetReadDeadline(time.Now().Add(time.Second * 1))
    if ok := scanner.Scan(); !ok {
      break
    }
    res := scanner.Text()
    log.Printf(">> %s", res)
    out += res + "\n"
  }
  /*buf := make([]byte, 1024)
    n, err := f.conn.Read(buf)
    if err != nil {
      return "", err
    }
    res := string(buf[:n])
    log.Printf(">> %s", res)
  */
  return out, nil
}

func (f *FluidSynth) SetSampleRate(rate int) {
  msg := fmt.Sprintf("set synth.sample-rate %d", rate)
  f.Send(msg)
}

func (f *FluidSynth) ResetToInitialFont() {
  font := f.Fonts[0]
  prog := font.Banks[0][0]
  log.Printf("Starting initial font %d: %s : %d %s\n", font.Id, font.Name, prog.Id, prog.Name)
  msg := fmt.Sprintf("select 0 %d 0 %d", font.Id, prog.Id)
  f.Send(msg)
}

// font ids are integers returned from fluidsynth, not guaranteed in sequence
func (f *FluidSynth) NextFont(m *Menu) (int, int) {
  chanId := m.currentPosition[1]
  fontId := m.currentPosition[2]
  if fontId == len(f.Fonts)-1 {
    fontId = 0
  } else {
    fontId++
  }
  font := f.Fonts[fontId]
  prog := font.Banks[0][0]
  log.Printf("Choosing first instrument prog of next font: %s, prog %s for channel %d\n", font.Name, prog.Name, chanId)
  msg := fmt.Sprintf("select %d %d 0 %d", chanId, font.Id, prog.Id)
  f.Send(msg)
  return fontId, 0 // reset prog Id
}

/* increase prog number, or restart with first */
func (f *FluidSynth) NextInstrumentProg(m *Menu) int {
  chanId := m.currentPosition[1]
  fontId := m.currentPosition[2]
  progId := m.currentPosition[3]
  font := f.Fonts[fontId]
  // ignore banks for now, just take first and increase prog
  if progId == len(font.Banks[0])-1 {
    progId = 0
  } else {
    progId++
  }
  prog := font.Banks[0][progId]
  log.Printf("Choosing next instrument prog: %s of font: %s for channel: %d\n", font.Name, prog.Name, chanId)
  msg := fmt.Sprintf("select %d %d 0 %d", chanId, font.Id, prog.Id)
  f.Send(msg)
  return progId
}
