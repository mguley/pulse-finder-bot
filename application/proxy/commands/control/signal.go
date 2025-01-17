package control

import (
	"bufio"
	"errors"
	"fmt"
	"infrastructure/proxy/port"
	"strings"
)

// SignalCommand represents a command to send a signal to the proxy control port.
type SignalCommand struct {
	conn   *port.Connection // The connection to the proxy control port.
	signal string           // The signal to be sent (e.g., "NEWNYM").
}

// NewSignalCommand creates a new SignalCommand command instance.
func NewSignalCommand(conn *port.Connection, signal string) *SignalCommand {
	return &SignalCommand{conn: conn, signal: signal}
}

// Execute sends the signal command to the proxy control port and processes the response.
func (c *SignalCommand) Execute() (err error) {
	if err := c.validate(); err != nil {
		return err
	}

	// Connect to the proxy control port
	if err = c.conn.Connect(); err != nil {
		return fmt.Errorf("could not connect to proxy: %w", err)
	}

	// Send the signal command
	if err = c.sendCommand(); err != nil {
		return fmt.Errorf("could not send signal command: %w", err)
	}

	res, err := c.conn.GetReader().ReadString('\n')
	return c.processResponse(res, err)
}

// sendCommand sends the signal command to the proxy.
func (c *SignalCommand) sendCommand() (err error) {
	conn := c.conn.GetConnection()
	if conn == nil {
		return errors.New("no active connection to send the command")
	}

	writer := bufio.NewWriter(conn)
	command := fmt.Sprintf("SIGNAL %s\r\n", c.signal)
	bytesWritten, err := writer.WriteString(command)
	if err != nil {
		return fmt.Errorf("could not send command: %s", err)
	}

	if err = writer.Flush(); err != nil {
		return fmt.Errorf("could not flush command: %s", err)
	}

	fmt.Printf("SIGNAL: wrote %d bytes to server\n", bytesWritten)
	return nil
}

// processResponse processes the response from the proxy after the SIGNAL command.
func (c *SignalCommand) processResponse(response string, err error) error {
	if err != nil {
		return fmt.Errorf("could not process response: %w", err)
	}
	fmt.Printf("SIGNAL: %s\n", response)

	switch {
	case strings.HasPrefix(response, "250"): // "250" indicates success
		return nil
	case strings.HasPrefix(response, "514"): // "514" indicates authentication failure
		return errors.New("authentication required")
	default:
		return fmt.Errorf("unexpected response: %s", response)
	}
}

// validate checks whether the SignalCommand is correctly initialized.
func (c *SignalCommand) validate() (err error) {
	if c.conn == nil {
		return errors.New("proxy connection is not initialized")
	}
	if c.signal == "" {
		return errors.New("signal is required")
	}
	return nil
}
