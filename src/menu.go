package main

import (
  "encoding/json"
  "io/ioutil"
  "log"
  "os"
  "sync"
)

type Menu struct {
  Functions       []MenuItem `json:"functions"`
  mu              sync.Mutex
  currentPosition [3]int // position
}

type MenuItem struct {
  Id      string
  Selects []Select
}

type Select struct {
  Id       string
  ToggleC  interface{}
  ToggleD  interface{}
  RotLeft  interface{}
  RotRight interface{}
}

func NewMenu() *Menu {
  var m Menu
  f, err := os.Open("./docs/menu.json")
  if err != nil {
    log.Fatal("Failed opening menu file")
  }
  defer f.Close()
  bs, _ := ioutil.ReadAll(f)
  if err := json.Unmarshal(bs, &m); err != nil {
    log.Fatalf("Failed parsing menu: %s", err)
  }
  m.mu = sync.Mutex{}
  m.currentPosition = [3]int{0, 0, 0}
  return &m
}

func (m *Menu) ToggleFunction() {
  m.mu.Lock()
  if m.currentPosition[0] == len(m.Functions)-1 {
    m.currentPosition[0] = 0
  } else {
    m.currentPosition[0]++
  }
  m.mu.Unlock()
}

/*
const (
  Hammond Instrument = iota
  Rhodes
  Organ
  Wurlitzer
  Moog
  Piano
)

func (i Instrument) String() string {
  return [...]string{"Hammond", "Rhodes", "Organ", "Wurlitzer", "Moog", "Piano"}[i]
}

const (
  Effect1 Effect = iota
  Effect2
  Effect3
)

func (e Effect) String() string {
  return [...]string{"Effect1", "Effect2", "Effect3"}[e]
}

const (
  Volume Mixer = iota
  Balance
)

func (m Mixer) String() string {
  return [...]string{"Volume", "Balance"}[m]
}
*/
