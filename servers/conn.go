package main

import (
	"io"
	"net"
	"time"
)

type conn struct {
	net.Conn

	IdleTimeout   time.Duration
	MaxReadBuffer int64
}

func (c *conn) Write(p []byte) (n int, err error) {
	n, err = c.Conn.Write(p)
	return
}

func (c *conn) Read(b []byte) (n int, err error) {
	r := io.LimitReader(c.Conn, c.MaxReadBuffer)
	n, err = r.Read(b)
	return
}
