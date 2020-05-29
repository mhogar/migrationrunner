package migrationrunner

import "errors"

// ChainError will combine the error message and the message together in an easy to read manner.
func ChainError(message string, err error) error {
	return errors.New(message + "\n\t" + err.Error())
}
