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

func (c *Conn) Read() (b []byte, err error) {
	if c.c == nil {
		return nil, errors.New("conn is nil")
	}

	if c.scan == nil {
		c.scan.Split(MySplitFunc)
	}

	if c.scan.Scan() {
		b = c.scan.Bytes()
	} else {
		c.c.Close()
		err = errors.New("conn scan error")
	}
	return
}

// The head of write method consists of modulus 0x0102 and length len(data)+4.
// Follow the data behind your head
func (c *Conn) Write(data []byte) (err error) {
	if len(data) > math.MaxUint16 {
		return errors.New("data too big")
	}

	buf := bytes.Buffer{}
	binary.Write(&buf, binary.LittleEndian, []byte{0x01, 0x02})
	binary.Write(&buf, binary.LittleEndian, uint16(len(data)+4))
	binary.Write(&buf, binary.BigEndian, data)
	_, err = c.c.Write(buf.Bytes())
	return
}
func MySplitFunc(data []byte, atEOF bool) (advance int, token []byte, err error) {
	//Judge length, magic number
	if len(data) < 4 || data[0] != 0x88 || data[1] != 0x94 {
		err = errors.New("protocol error")
		return
	}

	var header uint32
	//Get the header
	err = binary.Read(bytes.NewReader(data[:4]), binary.BigEndian, &header)
	if err != nil {
		log.Println(err.Error())
		return
	}

	l := uint16(header)
	//Making Size-to-Size Conversion
	l = (l >> 8) | (l << 8)
	//Read data by length. advance is read length, including header and data. Data is read data.
	if int(l) <= len(data) {
		advance, token, err = int(l), data[:int(l)], nil
	}
	if atEOF {
		err = errors.New("EOF")
	}

	return
}
