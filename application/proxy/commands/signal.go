package commands

import (
	"bufio"
	"errors"
	"fmt"
	"infrastructure/proxy/port"
	"strings"
)

// SignalCommand represents a command to send a signal to the proxy control port.
type SignalCommand struct {
	c      *port.Connection // The connection to the proxy control port.
	signal string           // The signal to be sent (e.g., "NEWNYM").
}

// NewSignalCommand creates a new SignalCommand command instance.
func NewSignalCommand(c *port.Connection, signal string) *SignalCommand {
	return &SignalCommand{c: c, signal: signal}
}

// Execute sends the signal command to the proxy control port and processes the response.
func (cmd *SignalCommand) Execute() error {
	if err := cmd.validate(); err != nil {
		return err
	}

	// Connect to the proxy control port
	if err := cmd.c.Connect(); err != nil {
		return fmt.Errorf("could not connect to proxy: %w", err)
	}

	// Send the signal command
	if err := cmd.sendCommand(); err != nil {
		return fmt.Errorf("could not send signal command: %w", err)
	}

	// Process the response from the proxy
	res, err := cmd.c.GetReader().ReadString('\n')
	return cmd.processResponse(res, err)
}

// sendCommand sends the signal command to the proxy.
func (cmd *SignalCommand) sendCommand() error {
	c := cmd.c.GetConnection()
	if c == nil {
		return errors.New("no active connection to send the command")
	}

	writer := bufio.NewWriter(c)
	command := fmt.Sprintf("SIGNAL %s\r\n", cmd.signal)
	bytesWritten, err := writer.WriteString(command)
	if err != nil {
		return fmt.Errorf("failed to send signal command: %w", err)
	}

	if err = writer.Flush(); err != nil {
		return fmt.Errorf("failed to flush authentication command: %w", err)
	}

	fmt.Printf("Signal command sent: Bytes written: %d\r\n", bytesWritten)
	return nil
}

// processResponse processes the response from the proxy after the SIGNAL command.
func (cmd *SignalCommand) processResponse(response string, err error) error {
	if err != nil {
		return fmt.Errorf("failed to read signal response: %w", err)
	}
	fmt.Printf("Signal response: %s\n", response)

	// Check the response prefix to determine the outcome
	switch {
	case strings.HasPrefix(response, "250"): // "250" indicates success
		return nil
	case strings.HasPrefix(response, "514"): // "514" indicates authentication failure
		return errors.New("authentication required")
	default:
		return fmt.Errorf("unexpected signal response: %s", response)
	}
}

// validate checks whether the SignalCommand is correctly initialized.
func (cmd *SignalCommand) validate() error {
	if cmd.c == nil {
		return errors.New("proxy connection is not initialized")
	}
	if cmd.signal == "" {
		return errors.New("signal is required")
	}
	return nil
}
