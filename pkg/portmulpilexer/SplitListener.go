package portmulpilexer

import (
	"bufio"
	"bytes"
	"crypto/tls"
	"log"
	"net"
)

// here's a buffered conn for peeking into the connection
type Conn struct {
	net.Conn
	buf *bufio.Reader
}

func (c *Conn) Read(b []byte) (int, error) {
	return c.buf.Read(b)
}

type SplitListener struct {
	net.Listener
	config *tls.Config
}

func (l *SplitListener) Accept() (net.Conn, error) {
	c, err := l.Listener.Accept()
	if err != nil {
		return nil, err
	}

	// buffer reads on our conn
	bconn := &Conn{
		Conn: c,
		buf:  bufio.NewReader(c),
	}

	// inspect the first few bytes
	hdr, err := bconn.buf.Peek(1)
	if err != nil {
		_ = bconn.Close()
		return nil, err
	}

	// I don't remember what the TLS handshake looks like,
	// but this works as a POC
	if bytes.Equal(hdr, []byte{22}) {
		log.Println("Https")
		return tls.Server(bconn, l.config), nil
	}
	log.Println("Http")
	return bconn, nil
}
func NewListener(inner net.Listener, config *tls.Config) net.Listener {
	return &SplitListener{Listener: inner, config: config}
}
