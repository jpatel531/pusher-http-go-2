package pusher

import (
	"github.com/pusher/pusher/requests"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"testing"
)

type mockDispatcher struct {
	mock.Mock
}

func (m mockDispatcher) sendRequest(p *Pusher, request *requests.Request, params *requests.Params) (response []byte, err error) {
	args := m.Called(p, request, params)
	return args.Get(0).([]byte), args.Error(1)
}

func TestSimpleTrigger(t *testing.T) {
	mDispatcher := &mockDispatcher{}
	p := &Pusher{
		appID:      "id",
		key:        "key",
		secret:     "secret",
		dispatcher: mDispatcher,
	}

	expectedParams := &requests.Params{
		Body: []byte(`{"name":"my-event","channels":["test-channel"],"data":"data"}`),
	}

	mDispatcher.
		On("sendRequest", p, requests.Trigger, expectedParams).
		Return([]byte("{}"), nil)

	_, err := p.Trigger("test-channel", "my-event", "data")
	assert.NoError(t, err)

	mDispatcher.AssertExpectations(t)
}
