package main

import (
  "log"
  "net/http"
)

type Server struct {
  addr  string
  clan  *ClanMan
  fluid *FluidSynth
}

func NewServer(a string, c *ClanMan, f *FluidSynth) *Server {
  return &Server{
    addr:  a,
    clan:  c,
    fluid: f,
  }
}

func (s *Server) statusHandler(w http.ResponseWriter, r *http.Request) {
  w.Write([]byte("OK"))
}

func (s *Server) eventHandler(w http.ResponseWriter, r *http.Request) {
  name, ok := r.URL.Query()["name"]
  origin, ok := r.URL.Query()["origin"]
  if !ok || len(name[0]) < 1 || len(origin[0]) < 1 {
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
