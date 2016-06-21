package pusher

import (
	"encoding/json"
	"github.com/pusher/pusher/errors"
)

type TriggerResponse struct {
	EventIds map[string]string `json:"event_ids,omitempty"`
}

type Event struct {
	Name     string  `json:"name"`
	Channel  string  `json:"channel"`
	Data     string  `json:"data"`
	SocketID *string `json:"socket_id,omitempty"`
}

type eventAnyData struct {
	Name     string      `json:"name"`
	Channels []string    `json:"-"`
	Data     interface{} `json:"data"`
	SocketID *string     `json:"socket_id,omitempty"`
}

type eventStringData struct {
	Name     string   `json:"name"`
	Channels []string `json:"channels"`
	Data     string   `json:"data"`
	SocketID *string  `json:"socket_id,omitempty"`
}

type batchRequest struct {
	Batch []Event `json:"batch"`
}

const maxDataSize = 10240

func (e *eventAnyData) toJSON() (body []byte, err error) {
	var dataBytes []byte

	switch d := e.Data.(type) {
	case []byte:
		dataBytes = d
	case string:
		dataBytes = []byte(d)
	default:
		if dataBytes, err = json.Marshal(e.Data); err != nil {
			return
		}
	}

	if len(dataBytes) > maxDataSize {
		err = errors.New("Data must be smaller than 10kb")
		return
	}

	return json.Marshal(&eventStringData{
		Name:     e.Name,
		Channels: e.Channels,
		SocketID: e.SocketID,
		Data:     string(dataBytes),
	})
}
