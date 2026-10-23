package emulator

import (
	"fmt"
	"os"
)

// captureAuthResponse is explicitly opt-in: responses can contain account secrets.
// Store only the server body, never request credentials or cookie headers.
func captureAuthResponse(body []byte) error {
	directory := os.Getenv("IPATOOL_AUTH_RESPONSE_DIR")
	if directory == "" {
		return nil
	}

	if err := os.MkdirAll(directory, 0o700); err != nil {
		return fmt.Errorf("create authentication diagnostic directory: %w", err)
	}

	info, err := os.Lstat(directory)
	if err != nil {
		return fmt.Errorf("inspect authentication diagnostic directory: %w", err)
	}
	if !info.IsDir() || info.Mode().Perm()&0o077 != 0 {
		return fmt.Errorf("authentication diagnostic directory must be a private directory")
	}

	file, err := os.CreateTemp(directory, "auth-response-*.body")
	if err != nil {
		return fmt.Errorf("create authentication diagnostic: %w", err)
	}
	_, writeErr := file.Write(body)
	closeErr := file.Close()
	if writeErr != nil {
		return fmt.Errorf("write authentication diagnostic: %w", writeErr)
	}
	if closeErr != nil {
		return fmt.Errorf("close authentication diagnostic: %w", closeErr)
	}

	return nil
}
