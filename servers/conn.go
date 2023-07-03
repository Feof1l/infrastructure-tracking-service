package main

import (
	"io"
	"net"
	"time"
)

type conn struct {
	net.Conn
	Timeout       time.Duration
	IdleTimeout   time.Duration
	MaxReadBuffer int64
}

func (c *conn) Write(p []byte) (n int, err error) {
	c.SetWriteDeadline(time.Now().Add(c.Timeout))
	n, err = c.Conn.Write(p)
	return
}

func (c *conn) Read(b []byte) (n int, err error) {
	c.SetWriteDeadline(time.Now().Add(c.Timeout))
	r := io.LimitReader(c.Conn, c.MaxReadBuffer)
	n, err = r.Read(b)
	return
}
