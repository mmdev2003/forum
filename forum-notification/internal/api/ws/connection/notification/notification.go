package notification

import (
	"forum-notification/internal/model"
	"forum-notification/pkg/api/forum-authorization"
	"log/slog"
	"net/http"
)

import (
	"github.com/gorilla/websocket"
	"github.com/labstack/echo/v4"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin:     func(r *http.Request) bool { return true },
}

func RegisterNotificationConnection(
	wsConnManager model.IWsConnManager,
) echo.HandlerFunc {
	return func(request echo.Context) error {
		authorizationData := request.Get(model.AuthorizationDataKey).(*forum_authorization.AuthorizationData)

		conn, err := upgrader.Upgrade(request.Response(), request.Request(), nil)
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
		}

		wsConnManager.AddConnection(authorizationData.AccountID, conn)

		go handleWebSocketConnection(authorizationData.AccountID, conn, wsConnManager)
		return nil
	}
}

func handleWebSocketConnection(
	accountID int,
	conn *websocket.Conn,
	wsConnManager model.IWsConnManager,
) {
	defer wsConnManager.RemoveConnection(accountID)
	for {
		_, _, err := conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err) {
				slog.Info("Клиент отключился", "accountID", accountID)
			} else {
				slog.Error("Ошибка чтения сообщения", "error", err, "accountID", accountID)
			}
			break
		}
	}
}
