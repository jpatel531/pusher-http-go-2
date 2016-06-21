package pusher

import (
	"encoding/json"
	"fmt"
	"github.com/pusher/pusher/requests"
	"github.com/pusher/pusher/signatures"
	"net/http"
	"time"
)

type Pusher struct {
	appID, key, secret string
	*Options
}

type Options struct {
	Host       string // host or host:port pair
	Secure     bool   // true for HTTPS
	Cluster    string
	HttpClient *http.Client
}

func New(appID, key, secret string) Client {
	return &Pusher{
		appID:  appID,
		key:    key,
		secret: secret,
		Options: &Options{
			Host:       "api.pusherapp.com",
			Secure:     true,
			HttpClient: &http.Client{Timeout: time.Second * 5},
		},
	}
}

func NewWithOptions(appID, key, secret string, options *Options) Client {
	return &Pusher{
		appID:   appID,
		key:     key,
		secret:  secret,
		Options: options,
	}
}

func (p *Pusher) Trigger(channel string, eventName string, data interface{}) (*TriggerResponse, error) {
	event := &Event{
		Channels: []string{channel},
		Name:     eventName,
		Data:     data,
	}
	return p.trigger(event)
}

func (p *Pusher) TriggerMulti(channels []string, eventName string, data interface{}) (*TriggerResponse, error) {
	event := &Event{
		Channels: channels,
		Name:     eventName,
		Data:     data,
	}
	return p.trigger(event)
}

func (p *Pusher) TriggerExclusive(channel string, eventName string, data interface{}, socketID string) (*TriggerResponse, error) {
	event := &Event{
		Channels: []string{channel},
		Name:     eventName,
		Data:     data,
		SocketID: &socketID,
	}
	return p.trigger(event)
}

func (p *Pusher) TriggerMultiExclusive(channels []string, eventName string, data interface{}, socketID string) (*TriggerResponse, error) {
	event := &Event{
		Channels: channels,
		Name:     eventName,
		Data:     data,
		SocketID: &socketID,
	}
	return p.trigger(event)
}

func (p *Pusher) trigger(event *Event) (response *TriggerResponse, err error) {
	var (
		eventJSON    []byte
		byteResponse []byte
	)

	if len(event.Channels) > 10 {
		err = newError("You cannot trigger on more than 10 channels at once")
		return
	}

	if err = validateChannels(event.Channels); err != nil {
		return
	}

	if err = validateSocketID(event.SocketID); err != nil {
		return
	}

	if eventJSON, err = event.toJSON(); err != nil {
		return
	}

	params := &requests.Params{
		Body: eventJSON,
	}

	if byteResponse, err = p.sendRequest(requests.Trigger, params); err != nil {
		return
	}

	err = json.Unmarshal(byteResponse, &response)
	return
}

func (p *Pusher) TriggerBatch(batch []Event) (response *TriggerResponse, err error) {
	var (
		byteResponse []byte
		batchJSON    []byte
	)

	if batchJSON, err = json.Marshal(&batch); err != nil {
		return
	}

	params := &requests.Params{
		Body: batchJSON,
	}

	if byteResponse, err = p.sendRequest(requests.TriggerBatch, params); err != nil {
		return
	}

	err = json.Unmarshal(byteResponse, &response)
	return
}

func (p *Pusher) Channels(additionalQueries map[string]string) (response *ChannelList, err error) {
	var byteResponse []byte

	params := &requests.Params{
		Queries: additionalQueries,
	}

	if byteResponse, err = p.sendRequest(requests.Channels, params); err != nil {
		return
	}

	fmt.Println(string(byteResponse))
	err = json.Unmarshal(byteResponse, &response)
	return
}

func (p *Pusher) Channel(name string, additionalQueries map[string]string) (response *Channel, err error) {
	var byteResponse []byte

	params := &requests.Params{
		Channel: name,
		Queries: additionalQueries,
	}

	if byteResponse, err = p.sendRequest(requests.Channel, params); err != nil {
		return
	}
	fmt.Println(string(byteResponse))
	err = json.Unmarshal(byteResponse, &response)
	return
}

func (p *Pusher) ChannelUsers(name string) (response *UserList, err error) {
	var byteResponse []byte

	params := &requests.Params{
		Channel: name,
	}

	if byteResponse, err = p.sendRequest(requests.ChannelUsers, params); err != nil {
		return
	}
	fmt.Println(string(byteResponse))
	err = json.Unmarshal(byteResponse, &response)
	return

}

func (p *Pusher) Authenticate(request AuthenticationRequest) (response []byte, err error) {
	var unsigned string
	if unsigned, err = request.StringToSign(); err != nil {
		return
	}
	authSignature := signatures.HMAC(unsigned, p.secret)

	return json.Marshal(map[string]string{
		"auth": fmt.Sprintf("%s: %s", p.key, authSignature),
	})
}

func (p *Pusher) Webhook(header http.Header, body []byte) (webhook *Webhook, err error) {
	for _, token := range header["X-Pusher-Key"] {
		if token == p.key && signatures.CheckHMAC(header.Get("X-Pusher-Signature"), p.secret, body) {
			return newWebhook(body)
		}
	}
	err = newError("Invalid webhook")
	return
}

func (p *Pusher) sendRequest(request *requests.Request, params *requests.Params) (response []byte, err error) {
	url := requestURL(p, request, params)
	return request.Do(p.HttpClient, url, params.Body)
}
