package main

import (
	"bufio"
	"errors"
	"io"
	"net"
	"time"
)

var (
	ErrEOF              = errors.New("...EOF")
	ErrConnectionClosed = errors.New("connection closed by peer")
)

type TelnetClient interface {
	Connect() error
	Close() error
	Send() error
	Receive() error
}

type Telnet struct {
	address           string
	timeout           time.Duration
	in                io.ReadCloser
	out               io.Writer
	connection        net.Conn
	stdinScanner      *bufio.Scanner
	connectionScanner *bufio.Scanner
}

func NewTelnetClient(address string, timeout time.Duration, in io.ReadCloser, out io.Writer) TelnetClient {
	return &Telnet{
		address: address,
		timeout: timeout,
		in:      in,
		out:     out,
	}
}

func (t *Telnet) Connect() error {
	var err error
	t.connection, err = net.DialTimeout("tcp", t.address, timeout)
	if err != nil {
		return err
	}
	// STDIN программы должен записываться в сокет
	t.stdinScanner = bufio.NewScanner(t.in)
	// данные, полученные из сокета, должны выводиться в STDOUT
	t.connectionScanner = bufio.NewScanner(t.connection)
	return nil
}

func (t *Telnet) Close() error {
	if t.connection != nil {
		err := t.connection.Close()
		if err != nil {
			return err
		}
	}
	return nil
}

func (t *Telnet) Send() error {
	if t.connection != nil {
		if t.stdinScanner.Scan() {
			bytes := append(t.stdinScanner.Bytes(), '\n')
			_, err := t.connection.Write(bytes)
			if err != nil {
				return err
			}
		} else {
			return ErrEOF
		}
	}
	return nil
}

func (t *Telnet) Receive() error {
	if t.connection != nil {
		if t.connectionScanner.Scan() {
			bytes := append(t.connectionScanner.Bytes(), '\n')
			_, err := t.out.Write(bytes)
			if err != nil {
				return err
			}
		} else {
			return ErrConnectionClosed
		}
	}
	return nil
}
