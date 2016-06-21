package pusher

import (
	"net/http"
)

type Client interface {
	Trigger(channel string, eventName string, data interface{}) (*TriggerResponse, error)
	TriggerMulti(channels []string, eventName string, data interface{}) (*TriggerResponse, error)
	TriggerExclusive(channel string, eventName string, data interface{}, socketID string) (*TriggerResponse, error)
	TriggerMultiExclusive(channels []string, eventName string, data interface{}, socketID string) (*TriggerResponse, error)
	TriggerBatch(batch []Event) (*TriggerResponse, error)

	Channels(additionalQueries map[string]string) (*ChannelList, error)
	Channel(name string, additionalQueries map[string]string) (*Channel, error)
	ChannelUsers(name string) (*UserList, error)

	Authenticate(request AuthenticationRequest) (response []byte, err error)

	Webhook(header http.Header, body []byte) (*Webhook, error)
}
