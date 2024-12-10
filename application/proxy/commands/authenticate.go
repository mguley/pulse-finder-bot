package commands

import (
	"bufio"
	"errors"
	"fmt"
	"infrastructure/proxy/port"
	"strings"
)

// AuthenticateCommand represents a command for authenticating with a proxy control port.
type AuthenticateCommand struct {
	c *port.Connection // The connection to the proxy control port.
}

// NewAuthenticateCommand creates a new AuthenticateCommand command instance.
func NewAuthenticateCommand(c *port.Connection) *AuthenticateCommand {
	return &AuthenticateCommand{c: c}
}

// Execute performs the authentication command on the proxy control port.
func (cmd *AuthenticateCommand) Execute() error {
	if cmd.c == nil {
		return errors.New("proxy connection is not initialized")
	}

	// Connect to the proxy control port
	if err := cmd.c.Connect(); err != nil {
		return fmt.Errorf("could not connect to proxy: %w", err)
	}

	// Send the AUTHENTICATE command
	if err := cmd.sendCommand(); err != nil {
		return fmt.Errorf("could not send authentication command: %w", err)
	}

	// Process the response from the proxy
	res, err := cmd.c.GetReader().ReadString('\n')
	return cmd.processResponse(res, err)
}

// sendCommand sends the authentication command to the proxy.
func (cmd *AuthenticateCommand) sendCommand() error {
	c := cmd.c.GetConnection()
	if c == nil {
		return errors.New("no active connection to send the command")
	}

	writer := bufio.NewWriter(c)
	command := fmt.Sprintf("AUTHENTICATE %q\n", cmd.c.GetPassword())
	bytesWritten, err := writer.WriteString(command)
	if err != nil {
		return fmt.Errorf("failed to send authentication command: %w", err)
	}

	if err = writer.Flush(); err != nil {
		return fmt.Errorf("failed to flush authentication command: %w", err)
	}

	fmt.Printf("Authentication command sent. Bytes written: %d\n", bytesWritten)
	return nil
}

// processResponse processes the response from the proxy after the AUTHENTICATE command.
func (cmd *AuthenticateCommand) processResponse(response string, err error) error {
	if err != nil {
		return fmt.Errorf("failed to read response: %w", err)
	}
	fmt.Printf("Authentication response: %s\n", response)

	// Check the response prefix to determine the outcome
	switch {
	case strings.HasPrefix(response, "250"): // "250" indicates success
		return nil
	case strings.HasPrefix(response, "515"): // "515" indicates authentication failure
		return errors.New("authentication failed: incorrect password")
	default:
		return fmt.Errorf("unexpected response: %s", response)
	}
}
