package main

import (
	"io"
	"net"
	"time"
)

type TelnetClient struct {
	conn   net.Conn
	in     io.Reader
	out    io.Writer
	closed bool
}

func NewTelnetClient(address string, timeout time.Duration, in io.ReadCloser, out io.Writer) (*TelnetClient, error) {
	conn, err := net.DialTimeout("tcp", address, timeout)
	if err != nil {
		return nil, err
	}
	return &TelnetClient{
		conn: conn,
		in:   in,
		out:  out,
	}, nil
}

func (tc *TelnetClient) Send() error {
	_, err := io.Copy(tc.conn, tc.in)
	return err
}

func (tc *TelnetClient) Receive() error {
	_, err := io.Copy(tc.out, tc.conn)
	return err
}

func (tc *TelnetClient) Close() error {
	tc.closed = true
	return tc.conn.Close()
}
