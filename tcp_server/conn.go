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
	c.updateDeadline()
	n, err = c.Conn.Write(p)
	return
}

func (c *conn) Read(b []byte) (n int, err error) {
	c.updateDeadline()
	//ограничиваем количество данных, которые мы читаем за один раз, но только в рамках одного соединения и одной операции чтения
	r := io.LimitReader(c.Conn, c.MaxReadBuffer)
	n, err = r.Read(b)
	return
}

func (c *conn) Close() (err error) {
	err = c.Conn.Close()
	return
}

// Соединение будет убито, если мы заново не обновим дедлайн.
// Каждый раз после удачной операции необходимо увеличивать дедлайн.
func (c *conn) updateDeadline() {
	idleDeadline := time.Now().Add(c.IdleTimeout)
	c.Conn.SetDeadline(idleDeadline)
}
