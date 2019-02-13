package helpers

import "errors"

// Error ...
func Error(msg string) error {
	return errors.New(msg)
}
