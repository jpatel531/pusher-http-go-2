package pusher

import (
	"fmt"
)

const ERROR_TAG = "[pusher-http-go]"

func newError(message string) error {
	return fmt.Errorf("%s: %s", ERROR_TAG, message)
}
