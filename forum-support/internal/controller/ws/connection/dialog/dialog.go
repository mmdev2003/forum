package dialog

import (
	"context"
	"encoding/json"
	"forum-support/internal/model"
	forum_authorization "forum-support/pkg/api/forum-authorization"
	"log/slog"
	"net/http"

	"github.com/gorilla/websocket"
	"github.com/labstack/echo/v4"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin:     func(r *http.Request) bool { return true },
}

func RegisterDialogConnection(
	dialogService model.IDialogService,
	wsConnManager model.IWsConnManager,
) echo.HandlerFunc {
	return func(request echo.Context) error {
		ctx := request.Request().Context()
		authorizationData := request.Get(model.AuthorizationDataKey).(*forum_authorization.AuthorizationData)

		conn, err := upgrader.Upgrade(request.Response(), request.Request(), nil)
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
		}

		err = wsConnManager.AddConnection(ctx, authorizationData.AccountID, conn)
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
		}

		go handleWebSocketConnection(authorizationData.AccountID, conn, wsConnManager, dialogService)
		return nil
	}
}

func handleWebSocketConnection(
	accountID int,
	conn *websocket.Conn,
	wsConnManager model.IWsConnManager,
	dialogService model.IDialogService,
) {
	defer wsConnManager.RemoveConnection(accountID)
	for {
		_, msg, err := conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err) {
				slog.Info("Клиент отключился", "accountID", accountID)
			} else {
				slog.Error("Ошибка чтения сообщения", "error", err, "accountID", accountID)
			}
			break
		}
		go HandleWsMsgForDialog(dialogService, msg)
	}
}

func HandleWsMsgForDialog(
	dialogService model.IDialogService,
	wsMsg []byte,
) {
	var dialogMessage model.DialogWsMessage

	err := json.Unmarshal(wsMsg, &dialogMessage)
	if err != nil {
		slog.Error("Ошибка распаковки сообщения", "error", err)
		return
	}

	_, err = dialogService.AddMessageToDialog(context.TODO(),
		dialogMessage.DialogID,
		dialogMessage.FromAccountID,
		dialogMessage.ToAccountID,
		dialogMessage.Text,
	)
	if err != nil {
		slog.Error("Ошибка отправки сообщения в диалог", "error", err, "dialogID", dialogMessage.DialogID)
	}
	return
}
