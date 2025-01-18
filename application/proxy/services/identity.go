package services

import (
	"application/proxy/commands"
	"application/proxy/commands/control"
	"application/proxy/strategies"
	"errors"
	"fmt"
	"infrastructure/proxy/port"
	"time"
)

// Identity manages interactions with the proxy control port, including authentication, signaling, and status retrieval.
type Identity struct {
	conn          *port.Connection             // Connection to the proxy control port.
	authCommand   *control.AuthenticateCommand // Command to authenticate with the proxy.
	signalCommand *control.SignalCommand       // Command to signal a new circuit.
	statusCommand *commands.StatusCommand      // Command to check the proxy circuit.
	retry         strategies.RetryStrategy     // Retry strategy for handling circuit change attempts.
	url           string                       // URL for checking proxy status or connectivity.
	attempts      int
}

// NewIdentity creates a new Identity instance.
func NewIdentity(
	conn *port.Connection,
	authCommand *control.AuthenticateCommand,
	signalCommand *control.SignalCommand,
	statusCommand *commands.StatusCommand,
	retry strategies.RetryStrategy,
	url string,
) *Identity {
	return &Identity{
		conn:          conn,
		authCommand:   authCommand,
		signalCommand: signalCommand,
		statusCommand: statusCommand,
		retry:         retry,
		url:           url,
		attempts:      7,
	}
}

// Change requests a new proxy circuit.
func (i *Identity) Change() (err error) {
	var circuitStatus string

	// Current circuit status
	if circuitStatus, err = i.statusCommand.Execute(i.url); err != nil {
		return fmt.Errorf("execute status command: %w", err)
	}
	fmt.Printf("circuit status: %v\n", circuitStatus)
	defer i.Close()

	// Attempt to change the circuit, up to i.attempts times
	return i.do(circuitStatus)
}

// do attempts to signal a new circuit and verifies the circuit change.
func (i *Identity) do(status string) (err error) {
	// Authenticate
	if err := i.authenticate(); err != nil {
		return err
	}

	for try := 1; try <= i.attempts+1; try++ {
		if err = i.trySignalAndCheck(status); err == nil {
			// Success
			return nil
		}
		fmt.Printf("Attempt #%d failed: %v\n", try, err)

		// Get the back-off duration
		var sleepTime time.Duration
		if sleepTime, err = i.retry.WaitDuration(try); err != nil {
			return fmt.Errorf("wait duration: %w", err)
		}

		// Sleep before the next attempt
		fmt.Printf("Attempt #%d: waiting %s before retrying...\n", try+1, sleepTime)
		time.Sleep(sleepTime)
	}

	return fmt.Errorf("maximum retry attempts exceeded")
}

// trySignalAndCheck tries to send the signal command and verifies if the circuit actually changed.
func (i *Identity) trySignalAndCheck(oldStatus string) (err error) {
	var newStatus string
	// Send signal command
	if err = i.signalCommand.Execute(); err != nil {
		return fmt.Errorf("signal command: %w", err)
	}

	// Get new circuit status
	if newStatus, err = i.statusCommand.Execute(i.url); err != nil {
		return fmt.Errorf("execute status command: %w", err)
	}
	fmt.Printf("new circuit status: %v\n", newStatus)

	// Compare
	if newStatus == oldStatus {
		return errors.New("circuit didn't change")
	}
	fmt.Println("circuit successfully changed")
	return nil
}

// authenticate performs the authentication process with the proxy via the control port.
func (i *Identity) authenticate() (err error) {
	if err = i.conn.Connect(); err != nil {
		return fmt.Errorf("connect to the control port: %w", err)
	}
	if err = i.authCommand.Execute(); err != nil {
		return fmt.Errorf("authenticate command: %w", err)
	}
	return nil
}

// Close terminates the connection to the proxy control port.
func (i *Identity) Close() {
	if err := i.conn.Close(); err != nil {
		fmt.Printf("close connection: %v\n", err)
	}
}
