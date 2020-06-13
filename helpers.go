package migrationrunner

import "errors"

func chainError(message string, err error) error {
	return errors.New(message + "\n\t" + err.Error())
}
