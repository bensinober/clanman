package main

import (
  "encoding/json"
  "io/ioutil"
  "log"
  "os"
  "sync"
)

type Menu struct {
  Functions       []Function
  mu              sync.Mutex
  currentPosition [4]int // position
}

type Function struct {
  Id      string
  Type    string
  Selects []Select
}

type Select struct {
  Id       string
  Actions  []Action
  Toggles  []Toggle    // TODO
  RotLeft  interface{} // TODO
  RotRight interface{} // TODO
}

type Action struct {
  Id     string
  Type   string
  Action string
}

type Toggle struct {
  Id      string
  Feature string
}

func NewMenu() *Menu {
  var m Menu
  f, err := os.Open("menu.json")
  if err != nil {
    log.Fatal("Failed opening menu file")
  }
  defer f.Close()
  bs, _ := ioutil.ReadAll(f)
  if err := json.Unmarshal(bs, &m); err != nil {
    log.Fatalf("Failed parsing menu: %s", err)
  }
  m.mu = sync.Mutex{}
  // positional buttons [Function, Select, Action, Toggle]
  m.currentPosition = [4]int{0, 0, 0, 0}
  return &m
}

func (m *Menu) GetActiveFunction() Function {
  m.mu.Lock()
  f := m.Functions[m.currentPosition[0]]
  m.mu.Unlock()
  return f
}

func (m *Menu) GetcurrentPosition() [4]int {
  m.mu.Lock()
  p := m.currentPosition
  m.mu.Unlock()
  return p
}

/* Generic next button functions */
func (m *Menu) NextFunction() {
  m.mu.Lock()
  if m.currentPosition[0] == len(m.Functions)-1 {
    m.currentPosition[0] = 0
  } else {
    m.currentPosition[0]++
  }
  m.currentPosition[1], m.currentPosition[2] = 0, 0 // reset submenus
  m.mu.Unlock()
}

func (m *Menu) NextSelect() {
  m.mu.Lock()
  if m.currentPosition[1] == len(m.Functions[m.currentPosition[0]].Selects)-1 {
    m.currentPosition[1] = 0
  } else {
    m.currentPosition[1]++
  }
  m.mu.Unlock()
}

func (m *Menu) NextAction() {
  m.mu.Lock()
  if m.currentPosition[2] == len(m.Functions[m.currentPosition[0]].Selects[m.currentPosition[1]].Actions)-1 {
    m.currentPosition[2] = 0
  } else {
    m.currentPosition[2]++
  }
  m.mu.Unlock()
}

func (m *Menu) NextToggle() {
  m.mu.Lock()
  if m.currentPosition[3] == len(m.Functions[m.currentPosition[0]].Selects[m.currentPosition[1]].Actions)-1 {
    m.currentPosition[3] = 0
  } else {
    m.currentPosition[3]++
  }
  m.mu.Unlock()
}

/* specific instrument toggle buttons
   btnC toggles font
   btnD toggles patch
*/

func (m *Menu) SelectFontId(id int) {
  m.mu.Lock()
  m.currentPosition[2] = id
  m.mu.Unlock()
}

func (m *Menu) SelectProgId(id int) {
  m.mu.Lock()
  m.currentPosition[3] = id
  m.mu.Unlock()
}
