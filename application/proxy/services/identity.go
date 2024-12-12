package services

import (
	"application/proxy/commands"
	"application/proxy/strategies"
	"fmt"
	"infrastructure/proxy/port"
	"time"
)

// Identity manages interactions with the proxy control port, including authentication, signaling, and status retrieval.
type Identity struct {
	conn      *port.Connection              // Connection to the proxy control port.
	authCmd   *commands.AuthenticateCommand // Command to authenticate with the proxy.
	signalCmd *commands.SignalCommand       // Command to signal a new circuit.
	statusCmd *commands.StatusCommand       // Command to check the proxy circuit.
	strategy  strategies.RetryStrategy      // Retry strategy for handling circuit change attempts.
	url       string                        // URL for checking proxy status or connectivity.
}

// NewIdentity creates a new Identity instance.
func NewIdentity(
	conn *port.Connection,
	authCmd *commands.AuthenticateCommand,
	signalCmd *commands.SignalCommand,
	statusCmd *commands.StatusCommand,
	strategy strategies.RetryStrategy,
	url string,
) *Identity {
	return &Identity{
		conn:      conn,
		authCmd:   authCmd,
		signalCmd: signalCmd,
		statusCmd: statusCmd,
		strategy:  strategy,
		url:       url,
	}
}

// Request handles the process of authenticating, checking the current status, and requesting a new circuit.
func (i *Identity) Request() error {
	if err := i.authenticate(); err != nil {
		return err
	}
	defer i.close()

	// Retrieve the initial circuit status
	status, err := i.statusCmd.Execute(i.url)
	if err != nil {
		return fmt.Errorf("could not get identity status: %w", err)
	}
	fmt.Printf("Circuit status: %s\n", status)

	// Attempt to request a new circuit
	return i.retry(status, i.url, 5)
}

// retry attempts to signal a new circuit and verifies the circuit change.
func (i *Identity) retry(status, url string, attempts int) error {
	for attempt := 0; attempt < attempts; attempt++ {
		if err := i.signalCmd.Execute(); err == nil {
			newStatus, err := i.statusCmd.Execute(url)
			if err != nil {
				return fmt.Errorf("failed to retrieve new circuit status: %w", err)
			}
			fmt.Printf("New circuit status: %s\n", newStatus)

			if status != newStatus {
				fmt.Println("Circuit successfully changed")
				return nil
			}
			fmt.Println("Circuit did not change; retrying...")
		}

		delay, err := i.strategy.WaitDuration(attempt)
		if err != nil {
			return fmt.Errorf("wait retry strategy failed: %w", err)
		}
		fmt.Printf("Attempt #%d: waiting %s before retrying...\n", attempt, delay)
		time.Sleep(delay)
	}
	return fmt.Errorf("maximum retry attempts exceeded")
}

// authenticate handles the process of authenticating with the proxy.
func (i *Identity) authenticate() error {
	if err := i.connect(); err != nil {
		return err
	}
	if err := i.authCmd.Execute(); err != nil {
		return fmt.Errorf("authentication error: %w", err)
	}
	return nil
}

// connect establishes a connection to the proxy control port.
func (i *Identity) connect() error {
	if err := i.conn.Connect(); err != nil {
		return fmt.Errorf("connect command failed: %w", err)
	}
	return nil
}

// close terminates the connection to the proxy control port.
func (i *Identity) close() {
	if err := i.conn.Close(); err != nil {
		fmt.Printf("close connection failed: %v\n", err)
	}
}
