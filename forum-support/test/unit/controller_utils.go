package unit

import (
	"fmt"
	"net/http"
	"net/url"
	"testing"

	"github.com/gorilla/websocket"
	"github.com/stretchr/testify/assert"
)

const (
	supportWsURL = "ws://localhost:8003/api/support/ws"
)

func CreateWsConnection(t *testing.T, accessToken string) (*websocket.Conn, error) {
	wsURL, err := url.Parse(supportWsURL)
	assert.NoError(t, err, "Не удалось разобрать URL")

	headers := http.Header{}
	headers.Add("Cookie", fmt.Sprintf("Access-Token=%s", accessToken))

	dialer := websocket.DefaultDialer
	conn, _, err := dialer.Dial(wsURL.String(), headers)
	assert.NoError(t, err, "Не удалось установить WebSocket-соединение")

	return conn, nil
}

func CreateDialog(t *testing.T) int {
	title := "title 1"
	description := "desc 1"
	requestID, err := testConfig.SupportClient.CreateSupportRequest("user", title, description)
	assert.NoError(t, err)

	dialogID, err := testConfig.DialogClient.CreateDialog(requestID, 1, "support")
	assert.NoError(t, err)
	assert.NotEqual(t, -1, dialogID)

	return dialogID

}
