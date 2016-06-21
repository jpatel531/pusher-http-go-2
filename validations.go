package pusher

import (
	"fmt"
	"github.com/pusher/pusher/errors"
	"regexp"
	s "strings"
)

var channelValidationRegex = regexp.MustCompile("^[-a-zA-Z0-9_=@,.;]+$")
var socketIDValidationRegex = regexp.MustCompile(`\A\d+\.\d+\z`)

func validateChannels(channels []string) (err error) {
	channelErrors := []string{}
	for _, channel := range channels {
		if len(channel) > 200 {
			channelErrors = append(channelErrors, channelTooLong(channel))
			continue
		}

		if !channelValidationRegex.MatchString(channel) {
			channelErrors = append(channelErrors, channelHasIllegalCharacters(channel))
			continue
		}
	}

	if len(channelErrors) > 0 {
		message := s.Join(channelErrors, ". ")
		err = errors.New(message)
	}

	return
}

func validateSocketID(socketID *string) (err error) {
	if (socketID == nil) || socketIDValidationRegex.MatchString(*socketID) {
		return
	}
	return errors.New("socket_id invalid")
}

func channelTooLong(channel string) string {
	return fmt.Sprintf("%s is over 200 characters.", channel)
}

func channelHasIllegalCharacters(channel string) string {
	return fmt.Sprintf("%s has illegal characters.", channel)
}
