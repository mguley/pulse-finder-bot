package control

import (
	"bufio"
	"errors"
	"fmt"
	"infrastructure/proxy/port"
	"strings"
)

// AuthenticateCommand represents a command for authenticating with a proxy control port.
type AuthenticateCommand struct {
	conn *port.Connection // The connection to the proxy control port.
}

// NewAuthenticateCommand creates a new AuthenticateCommand command instance.
func NewAuthenticateCommand(conn *port.Connection) *AuthenticateCommand {
	return &AuthenticateCommand{conn: conn}
}

// Execute performs the authentication command on the proxy control port.
func (c *AuthenticateCommand) Execute() (err error) {
	if c.conn == nil {
		return errors.New("proxy connection is not initialized")
	}

	// Connect to the proxy control port
	if err = c.conn.Connect(); err != nil {
		return fmt.Errorf("could not connect to proxy: %w", err)
	}

	// Send the AUTHENTICATE command
	if err = c.sendCommand(); err != nil {
		return fmt.Errorf("could not send authentication command: %w", err)
	}

	res, err := c.conn.GetReader().ReadString('\n')
	return c.processResponse(res, err)
}

// sendCommand sends the authentication command to the proxy.
func (c *AuthenticateCommand) sendCommand() (err error) {
	conn := c.conn.GetConnection()
	if conn == nil {
		return errors.New("no active connection to send the command")
	}

	writer := bufio.NewWriter(conn)
	command := fmt.Sprintf("AUTHENTICATE %q\n", c.conn.GetPassword())
	bytesWritten, err := writer.WriteString(command)
	if err != nil {
		return fmt.Errorf("could not write command: %w", err)
	}

	if err = writer.Flush(); err != nil {
		return fmt.Errorf("could not flush command: %w", err)
	}

	fmt.Printf("AUTHENTICATE: wrote %d bytes to server\n", bytesWritten)
	return nil
}

// processResponse processes the response from the proxy after the AUTHENTICATE command.
func (c *AuthenticateCommand) processResponse(response string, err error) error {
	if err != nil {
		return fmt.Errorf("could not process response: %w", err)
	}
	fmt.Printf("AUTHENTICATE: %s\n", response)

	switch {
	case strings.HasPrefix(response, "250"): // "250" indicates success
		return nil
	case strings.HasPrefix(response, "515"): // "515" indicates authentication failure
		return errors.New("authentication failed: incorrect password")
	default:
		return fmt.Errorf("unexpected response: %s", response)
	}
}
