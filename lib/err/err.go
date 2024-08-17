package err

import "fmt"

func Wrap(err error, message string) error {
	return fmt.Errorf("%s: %w", message, err)
}

func WrapIfError(err error, message string) error {
	if err == nil {
		return nil
	}

	return Wrap(err, message)
}
