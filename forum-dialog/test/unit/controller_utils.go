package unit

import (
	"fmt"
	"github.com/gorilla/websocket"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/url"
	"testing"
)

const (
	dialogWsURL = "ws://localhost:8003/api/dialog/ws"
)

func CreateWsConnection(t *testing.T, accessToken string) (*websocket.Conn, error) {
	wsURL, err := url.Parse(dialogWsURL)
	assert.NoError(t, err, "Не удалось разобрать URL")

	headers := http.Header{}
	headers.Add("Cookie", fmt.Sprintf("Access-Token=%s", accessToken))

	dialer := websocket.DefaultDialer
	conn, _, err := dialer.Dial(wsURL.String(), headers)
	assert.NoError(t, err, "Не удалось установить WebSocket-соединение")

	return conn, nil
}
