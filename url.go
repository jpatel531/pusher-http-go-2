package pusher

import (
	"fmt"
	"github.com/pusher/pusher/requests"
	"github.com/pusher/pusher/signatures"
	"net/url"
	s "strings"
)

func unescapeURL(values *url.Values) string {
	unesc, _ := url.QueryUnescape(values.Encode())
	return unesc
}

func requestURL(p *Pusher, request *requests.Request, params *requests.Params) *url.URL {
	values := params.URLValues(p.key)

	var path string
	if params.Channel != "" {
		path = fmt.Sprintf(request.PathPattern, p.appID, params.Channel)
	} else {
		path = fmt.Sprintf(request.PathPattern, p.appID)
	}

	unsigned := s.Join([]string{request.Method, path, unescapeURL(values)}, "\n")
	signed := signatures.HMAC(unsigned, p.secret)
	values.Add("auth_signature", signed)

	host := "api.pusherapp.com"
	scheme := "http"

	if p.Host != "" {
		host = p.Host
	}

	if p.Cluster != "" {
		host = fmt.Sprintf("api-%s.pusher.com", p.Cluster)
	}

	if p.Secure {
		scheme = "https"
	}

	return &url.URL{
		Scheme:   scheme,
		Host:     host,
		Path:     path,
		RawQuery: values.Encode(),
	}
}
