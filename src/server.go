package main

import (
  "encoding/json"
  "log"
  "net/http"
)

type Server struct {
  addr string
  clan *ClanMan
}

type ServerStatus struct {
  Fonts []Font
  Menu  *Menu
}

func NewServer(a string, c *ClanMan) *Server {
  return &Server{
    addr: a,
    clan: c,
  }
}

func (s *Server) ServerStatus() *ServerStatus {
  //now := time.Now()
  //uptime := now.Sub(s.startTime)
  return &ServerStatus{
    Fonts: s.clan.Fluid.Fonts,
    Menu:  s.clan.Menu,
  }
}

func (s *Server) statusHandler(w http.ResponseWriter, r *http.Request) {
  b, err := json.Marshal(s.ServerStatus())
  if err != nil {
    http.Error(w, err.Error(), http.StatusInternalServerError)
    return
  }
  w.Header().Set("Content-Type", "application/json")
  w.Write(b)
}

func (s *Server) eventHandler(w http.ResponseWriter, r *http.Request) {
  name, okName := r.URL.Query()["name"]
  origin, okOrigin := r.URL.Query()["origin"]
  if !okName || !okOrigin || len(name[0]) < 1 || len(origin[0]) < 1 {
    http.Error(w, "Url Param 'name' or 'origin' is missing", http.StatusBadRequest)
    return
  }
  ie := InputEvent{name[0], origin[0]}
  s.clan.InputEvents <- ie
  w.Write([]byte("OK"))
}

func (s *Server) Run() {
  mux := http.NewServeMux()
  mux.HandleFunc("/.status", s.statusHandler)
  mux.HandleFunc("/event", s.eventHandler)
  log.Printf("Starting web server at port %s\n", s.addr)
  log.Fatal(http.ListenAndServe(s.addr, mux))
}
