# net conn

我打算來講一下net conn 到底做了什麼，這是go 封裝的一個庫
他的interface如下
```go
type Conn interface {

	Read(b []byte) (n int, err error)

	Write(b []byte) (n int, err error)

	Close() error

	LocalAddr() Addr

	RemoteAddr() Addr

	SetDeadline(t time.Time) error

	SetReadDeadline(t time.Time) error

	SetWriteDeadline(t time.Time) error
}
```
那我們找了一下他的實作是交給tcpsock
```go
type TCPConn struct {
	conn
}

// SyscallConn returns a raw network connection.
// This implements the [syscall.Conn] interface.
func (c *TCPConn) SyscallConn() (syscall.RawConn, error) {
	if !c.ok() {
		return nil, syscall.EINVAL
	}
	return newRawConn(c.fd), nil
}

// ReadFrom implements the [io.ReaderFrom] ReadFrom method.
func (c *TCPConn) ReadFrom(r io.Reader) (int64, error) {
	if !c.ok() {
		return 0, syscall.EINVAL
	}
	n, err := c.readFrom(r)
	if err != nil && err != io.EOF {
		err = &OpError{Op: "readfrom", Net: c.fd.net, Source: c.fd.laddr, Addr: c.fd.raddr, Err: err}
	}
	return n, err
}

// WriteTo implements the io.WriterTo WriteTo method.
func (c *TCPConn) WriteTo(w io.Writer) (int64, error) {
	if !c.ok() {
		return 0, syscall.EINVAL
	}
	n, err := c.writeTo(w)
	if err != nil && err != io.EOF {
		err = &OpError{Op: "writeto", Net: c.fd.net, Source: c.fd.laddr, Addr: c.fd.raddr, Err: err}
	}
	return n, err
}

// CloseRead shuts down the reading side of the TCP connection.
// Most callers should just use Close.
func (c *TCPConn) CloseRead() error {
	if !c.ok() {
		return syscall.EINVAL
	}
	if err := c.fd.closeRead(); err != nil {
		return &OpError{Op: "close", Net: c.fd.net, Source: c.fd.laddr, Addr: c.fd.raddr, Err: err}
	}
	return nil
}

// CloseWrite shuts down the writing side of the TCP connection.
// Most callers should just use Close.
func (c *TCPConn) CloseWrite() error {
	if !c.ok() {
		return syscall.EINVAL
	}
	if err := c.fd.closeWrite(); err != nil {
		return &OpError{Op: "close", Net: c.fd.net, Source: c.fd.laddr, Addr: c.fd.raddr, Err: err}
	}
	return nil
}
```

我們主要追蹤write跟read 了解他們從哪來又去了哪裡
```go
func (c *TCPConn) writeTo(w io.Writer) (int64, error) {
	if n, err, handled := spliceTo(w, c.fd); handled {
		return n, err
	}
	return genericWriteTo(c, w)
}
```
可以看到go 這邊呼叫了 spliceTo , 這主要是希望對OS 進行一個 ZERO COPY 的操作
```go
func spliceTo(w io.Writer, fd *netFD) (int64, error, bool) {
	if splicer, ok := w.(splice.Writer); ok {
		return splicer.SpliceTo(fd.pfd)
	}
	return 0, nil, false
}
```
那如果不行他就會走userspace正常的gen一個write出來

