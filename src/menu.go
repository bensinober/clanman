package main

import (
  "encoding/json"
  "io/ioutil"
  "log"
  "os"
  "sync"
)

type Menu struct {
  Functions       []MenuItem
  mu              sync.Mutex
  currentPosition [3]int // position
}

type MenuItem struct {
  Id      string
  Selects []Select
}

type Select struct {
  Id       string
  Actions  []Action
  ToggleD  []Action
  RotLeft  interface{}
  RotRight interface{}
}

type Action struct {
  Id     string
  Type   string
  Action string
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
  m.currentPosition[1], m.currentPosition[2] = 0, 0 // reset submenus
  m.mu.Unlock()
}

func (m *Menu) ToggleSelect() {
  m.mu.Lock()
  if m.currentPosition[1] == len(m.Functions[m.currentPosition[0]].Selects)-1 {
    m.currentPosition[1] = 0
  } else {
    m.currentPosition[1]++
  }
  m.mu.Unlock()
}

func (m *Menu) ToggleAction() {
  m.mu.Lock()
  if m.currentPosition[2] == len(m.Functions[m.currentPosition[0]].Selects[m.currentPosition[1]].Actions)-1 {
    m.currentPosition[2] = 0
  } else {
    m.currentPosition[2]++
  }
  m.mu.Unlock()
}
