package server

import (
	"net"
	"sync"
)

type Server struct {
	clients []Client
	mutex   sync.Mutex
}

type Client struct {
	conn   net.Conn
	Pseudo string
	// Messages []string
}

var Log []Historic

type Historic struct {
	Time    string
	Pseudo  string
	Message string
}

const (
	IP   = "localhost"
	PORT = "8081"
)
