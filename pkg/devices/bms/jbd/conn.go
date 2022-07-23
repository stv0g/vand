package jbd

import (
	"net"

	"golang.org/x/sys/unix"
)

type Conn struct {
	net.Conn
}

func (c *Conn) Read(b []byte) (n int, err error) {
	if n, err = c.Conn.Read(b); err == unix.EPIPE {
		if err = c.Reconnect(); err != nil {
			return 0, err
		}

		return c.Conn.Read(b)
	} else {
		return n, nil
	}
}

func (c *Conn) Write(b []byte) (n int, err error) {
	if n, err = c.Conn.Write(b); err == unix.EPIPE {
		if err = c.Reconnect(); err != nil {
			return 0, err
		}

		return c.Conn.Write(b)
	} else {
		return n, nil
	}
}

func (c *Conn) Reconnect() (err error) {
	if c.Conn, err = net.Dial(c.RemoteAddr().Network(), c.RemoteAddr().String()); err != nil {
		return err
	}

	return nil
}
