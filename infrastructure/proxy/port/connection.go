package port

import (
	"bufio"
	"fmt"
	"net"
	"time"
)

// Connection represents a connection to a Proxy control port.
type Connection struct {
	address  string        // Address of the Proxy control port.
	timeout  time.Duration // Timeout duration for the connection.
	password string        // Password for authenticating with the Proxy control port.
	conn     net.Conn      // Established connection
	reader   *bufio.Reader // Buffered reader for reading responses from the Proxy control port.
}

// NewConnection creates a new Connection instance.
func NewConnection(address, password string, timeout time.Duration) *Connection {
	return &Connection{
		address:  address,
		timeout:  timeout,
		password: password,
		conn:     nil,
		reader:   nil,
	}
}

// validate checks whether the required fields are valid.
func (c *Connection) validate() error {
	if c.address == "" {
		return fmt.Errorf("address is required but not provided")
	}
	if c.timeout <= 0 {
		return fmt.Errorf("timeout must be greater than zero")
	}
	if c.password == "" {
		return fmt.Errorf("password is required but not provided")
	}
	return nil
}

// GetPassword retrieves the password for authenticating with the Proxy control port.
func (c *Connection) GetPassword() string {
	return c.password
}

// GetConnection retrieves the established network connection.
func (c *Connection) GetConnection() net.Conn {
	return c.conn
}

// GetReader retrieves the buffered reader for reading responses from the Proxy control port.
func (c *Connection) GetReader() *bufio.Reader {
	return c.reader
}

// Connect establishes a connection to the Proxy control port.
func (c *Connection) Connect() error {
	if err := c.validate(); err != nil {
		return fmt.Errorf("connection validation: %w", err)
	}
	if c.conn != nil {
		return nil
	}
	conn, err := net.DialTimeout("tcp", c.address, c.timeout)
	if err != nil {
		return fmt.Errorf("could not connect to control port %s: %s", c.address, err)
	}

	c.conn = conn
	c.reader = bufio.NewReader(conn)
	fmt.Printf("connected to control port %s\n", c.address)
	return nil
}

// Close terminates the connection to the Proxy control port.
func (c *Connection) Close() error {
	if c.conn == nil {
		return nil
	}

	if err := c.conn.Close(); err != nil {
		return fmt.Errorf("could not close connection to %s: %w", c.address, err)
	}

	c.conn = nil
	c.reader = nil
	fmt.Printf("connection to %s closed\n", c.address)
	return nil
}
