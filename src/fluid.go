package main

import (
  "fmt"
  "log"
  "net"
)

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
