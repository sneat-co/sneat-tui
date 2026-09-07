package commands

// ExitCodeError wraps a command error with the process exit code main.go
// should use instead of the default 1. Most command failures don't need
// this: SilenceErrors already routes RunE's returned error to stderr with
// exit code 1. It exists for the one case in the Action Protocol where the
// exit code itself is part of the machine-readable contract — an agent
// scripting `sneat action commit` needs to tell "the server refused because
// validation still needs input" (2) apart from any other failure (1)
// without parsing the error text.
type ExitCodeError struct {
	Code int
	Err  error
}

// Error implements the error interface by delegating to the wrapped error.
func (e *ExitCodeError) Error() string { return e.Err.Error() }

// Unwrap lets errors.Is/errors.As see through to the wrapped error.
func (e *ExitCodeError) Unwrap() error { return e.Err }
