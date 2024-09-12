package server

import (
	"net"
	"sync"
)

type Server struct {
	clients []Client
	mutex   sync.Mutex
	IP      string
	PORT    string
}

type Client struct {
	conn   net.Conn
	Pseudo string
	Color  string
}

var Log []Historic

var Txt []string

var connected string

type Historic struct {
	Time    string
	Pseudo  string
	Message string
}

const (
	IP   = "localhost"
	PORT = "8081"
)
