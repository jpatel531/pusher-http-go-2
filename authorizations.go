package pusher

import (
	"encoding/json"
	"fmt"
	"net/url"
)

type AuthenticationRequest interface {
	StringToSign() (string, error)
	UserData() (string, error)
}

type PrivateChannelRequest struct {
	Body []byte
}

func (p *PrivateChannelRequest) StringToSign() (unsigned string, err error) {
	var (
		params             url.Values
		keyExists          bool
		channelNameWrapper []string
		socketIDWrapper    []string
	)

	if params, err = url.ParseQuery(string(p.Body)); err != nil {
		return
	}

	if channelNameWrapper, keyExists = params["channel_name"]; !keyExists || len(channelNameWrapper) == 0 {
		err = newError("Channel param not found")
		return
	}

	if socketIDWrapper, keyExists = params["socket_id"]; !keyExists || len(socketIDWrapper) == 0 {
		err = newError("Socket_id not found")
		return
	}

	channelName := channelNameWrapper[0]
	socketID := socketIDWrapper[0]

	if err = validateSocketID(&socketID); err != nil {
		return
	}

	unsigned = fmt.Sprintf("%s:%s", socketID, channelName)
	return
}

func (p *PrivateChannelRequest) UserData() (userData string, err error) {
	return
}

type PresenceChannelRequest struct {
	Body   []byte
	Member *Member
}

type Member struct {
	UserId   string            `json:"user_id"`
	UserInfo map[string]string `json:"user_info,omitempty"`
}

func (p *PresenceChannelRequest) StringToSign() (unsigned string, err error) {
	privateChannelRequest := &PrivateChannelRequest{p.Body}
	if unsigned, err = privateChannelRequest.StringToSign(); err != nil {
		return
	}
	var userData string
	if userData, err = p.UserData(); err != nil {
		return
	}

	unsigned = fmt.Sprintf("%s:%s", unsigned, userData)
	return
}

func (p *PresenceChannelRequest) UserData() (userData string, err error) {
	var userDataBytes []byte
	if userDataBytes, err = json.Marshal(p.Member); err != nil {
		return
	}
	userData = string(userDataBytes)
	return
}
