package main

import (
	"bytes"
	"container/list"
	"github.com/Allenxuxu/gev"
	"log"
)

type server struct {
	conn   *list.List
	server *gev.Server

	buf bytes.Buffer
}

func NewServer(ip, port string) (*server, error) {
	var err error
	s := new(server)
	s.conn = list.New()

	s.server, err = gev.NewServer(s,
		gev.Network("tcp"),
		gev.Address(ip+":"+port),
		gev.NumLoops(-1),
		gev.MetricsServer("", ":9091"),
	)
	if err != nil {
		return nil, err
	}

	return s, nil
}

func (s *server) OnConnect(c *gev.Connection) {
	log.Println(" OnConnect ： ", c.PeerAddr())
}
func (s *server) OnMessage(c *gev.Connection, ctx interface{}, data []byte) (out interface{}) {
	lastDelimIndex := bytes.LastIndexByte(data, '\n')
	if lastDelimIndex == -1 { // '\n' not present in 'data'
		s.buf.Write(data)
		return
	}

	if lastDelimIndex == len(data)-1 { // data ends with '\n'
		s.buf.Write(data)
		out = s.buf.Bytes()
		s.buf.Reset()
		return
	}

	if s.buf.Len() != 0 { // if buf not empty
		out = append(s.buf.Bytes(), data[:lastDelimIndex+1]...)
	} else { // buf empty (previous element ends with '\n')
		out = data[:lastDelimIndex+1]
	}

	s.buf.Reset()
	s.buf.Write(data[lastDelimIndex+1:])
	return
}

func (s *server) OnClose(c *gev.Connection) {
	log.Println("OnClose")
}

func main() {
	s, err := NewServer("", "1833")
	if err != nil {
		panic(err)
	}
	defer s.server.Stop()

	s.server.Start()
}
