package pusher

import (
	"encoding/json"
)

type TriggerResponse struct {
	EventIds map[string]string `json:"event_ids,omitempty"`
}

type Event struct {
	Name     string
	Channels []string
	Data     interface{}
	SocketID *string
}

type event struct {
	Name     string   `json:"name"`
	Channels []string `json:"channels"`
	Data     string   `json:"data"`
	SocketID *string  `json:"socket_id,omitempty"`
}

type batchRequest struct {
	Batch []Event `json:"batch"`
}

const maxDataSize = 10240

func (e *Event) toJSON() (body []byte, err error) {
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
		err = newError("Data must be smaller than 10kb")
		return
	}

	return json.Marshal(&event{
		Name:     e.Name,
		Channels: e.Channels,
		SocketID: e.SocketID,
		Data:     string(dataBytes),
	})
}
