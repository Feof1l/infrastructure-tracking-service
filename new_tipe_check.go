package main

import (
	"bufio"
	"bytes"
	"encoding/binary"
	"errors"
	"log"
	"math"
	"net"

	"google.golang.org/grpc/metadata"
)

type Conn struct {
	c    net.Conn
	meta metadata.MD
	scan *bufio.Scanner
}

func (c *Conn) Read() ([]byte, error) {
	var b []byte
	var err error
	if c.c == nil {
		return nil, errors.New("conn is nil")
	}

	if c.scan == nil {
		c.scan.Split(MySplitFunc)
	}

	if c.scan.Scan() { // конец чтения( eof, или иная ошбика)
		b = c.scan.Bytes()
	} else {
		c.c.Close()
		err = errors.New("conn scan error")
	}
	return b, err
}

// Заголовок метода записи состоит из модуля 0x0102 и длины len(данные)+4.
func (c *Conn) Write(data []byte) (err error) {
	if len(data) > math.MaxUint16 { // слишком большие значения данных
		return errors.New("data too big")
	}

	buf := bytes.Buffer{}
	binary.Write(&buf, binary.LittleEndian, []byte{0x01, 0x02}) //записывает двоичное представление данных в buf
	binary.Write(&buf, binary.LittleEndian, uint16(len(data)+4))
	binary.Write(&buf, binary.BigEndian, data)
	_, err = c.c.Write(buf.Bytes())
	return
}
func MySplitFunc(data []byte, atEOF bool) (advance int, token []byte, err error) {
	if len(data) < 4 || data[0] != 0x88 || data[1] != 0x94 {
		err = errors.New("protocol error")
		return
	}

	var header uint32
	err = binary.Read(bytes.NewReader(data[:4]), binary.BigEndian, &header)
	if err != nil {
		log.Println(err.Error())
		return
	}

	l := uint16(header)
	l = (l >> 8) | (l << 8)
	//Считывание данных по длине. advance - длина считывания, включая заголовок и данные.
	if int(l) <= len(data) {
		advance, token, err = int(l), data[:int(l)], nil
	}
	if atEOF {
		err = errors.New("EOF")
	}

	return
}
