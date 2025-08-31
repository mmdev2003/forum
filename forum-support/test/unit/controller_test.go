package unit

import (
	"fmt"
	"forum-support/internal/model"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestControllerCreateSupportRequest(t *testing.T) {
	testConfig.PrepareDB()

	title := "title 1"
	description := "desc 1"
	requestID, err := testConfig.SupportClient.CreateSupportRequest("user", title, description)
	assert.NoError(t, err, err)
	assert.NotEqual(t, requestID, -1)

	actualRequestId, err := testConfig.SupportClient.GetRequestById("user", requestID)

	assert.NoError(t, err, err)

	assert.Equal(t, requestID, actualRequestId.ID)
	assert.Equal(t, title, actualRequestId.Title)
	assert.Equal(t, description, actualRequestId.Description)
	assert.Nil(t, actualRequestId.Status)

	err = testConfig.SupportClient.OpenSupportRequest("user", requestID)
	assert.Error(t, err, err)

	err = testConfig.SupportClient.OpenSupportRequest("support", requestID)
	assert.NoError(t, err, err)

	actualRequestId, err = testConfig.SupportClient.GetRequestById("user", requestID)

	assert.NoError(t, err, err)

	assert.Equal(t, requestID, actualRequestId.ID)
	assert.Equal(t, title, actualRequestId.Title)
	assert.Equal(t, description, actualRequestId.Description)
	assert.Equal(t, string(model.OpenRequestStatus), *actualRequestId.Status)

	err = testConfig.SupportClient.CloseSupportRequest("user", requestID)
	assert.NoError(t, err, err)

	actualRequestId, err = testConfig.SupportClient.GetRequestById("user", requestID)

	assert.NoError(t, err, err)

	assert.Equal(t, requestID, actualRequestId.ID)
	assert.Equal(t, title, actualRequestId.Title)
	assert.Equal(t, description, actualRequestId.Description)
	assert.Equal(t, string(model.ClosedRequestStatus), *actualRequestId.Status)
}

func TestControllerCreateDialog(t *testing.T) {
	testConfig.PrepareDB()
	userID := 1

	dialogID, err := testConfig.DialogClient.CreateDialog(2, userID, "support")
	assert.Error(t, err)
	assert.Equal(t, -1, dialogID)

	title := "title 1"
	description := "desc 1"
	requestID, err := testConfig.SupportClient.CreateSupportRequest("user", title, description)
	assert.NoError(t, err)

	dialogID, err = testConfig.DialogClient.CreateDialog(requestID, userID, "user")
	assert.Error(t, err)
	assert.Equal(t, -1, dialogID)

	dialogID, err = testConfig.DialogClient.CreateDialog(requestID, userID, "support")
	assert.NoError(t, err)
	assert.NotEqual(t, -1, dialogID)

	dialogs, err := testConfig.DialogClient.GetDialogs("user")
	assert.NoError(t, err)

	supportDialogs, err := testConfig.DialogClient.GetDialogs("support")
	assert.NoError(t, err)

	assert.Equal(t, dialogID, dialogs[0].ID)
	assert.Equal(t, requestID, dialogs[0].SupportRequestID)
	assert.Equal(t, userID, dialogs[0].UserAccountID)

	assert.Equal(t, dialogID, supportDialogs[0].ID)
	assert.Equal(t, requestID, supportDialogs[0].SupportRequestID)
	assert.Equal(t, userID, supportDialogs[0].UserAccountID)

	dialogs, err = testConfig.DialogClient.GetDialogs("o_user")
	assert.NoError(t, err)

	assert.Empty(t, dialogs)
}

func TestWebSocketMessageExchange(t *testing.T) {
	testConfig.PrepareDB()

	conn1, err := CreateWsConnection(t, "user")
	assert.NoError(t, err)
	defer conn1.Close()

	conn2, err := CreateWsConnection(t, "support")
	assert.NoError(t, err)
	defer conn2.Close()

	dialogID := CreateDialog(t)

	msg := model.DialogWsMessage{
		DialogID:      dialogID,
		FromAccountID: 2,
		ToAccountID:   1,
		Text:          "ответ от саппорта",
	}

	err = conn2.WriteJSON(msg)
	assert.NoError(t, err)

	var receivedMsg model.DialogWsMessage

	fmt.Println("read message from")
	err = conn1.ReadJSON(&receivedMsg)
	assert.NoError(t, err)

	assert.Equal(t, msg.DialogID, receivedMsg.DialogID)
	assert.Equal(t, msg.FromAccountID, receivedMsg.FromAccountID)
	assert.Equal(t, msg.ToAccountID, receivedMsg.ToAccountID)
	assert.Equal(t, msg.Text, receivedMsg.Text)

	messages, err := testConfig.DialogClient.MessagesByDialogID(dialogID, "user")
	assert.NoError(t, err)

	assert.Equal(t, 1, len(messages))
	assert.Equal(t, msg.Text, messages[0].MessageText)
	assert.Equal(t, msg.FromAccountID, messages[0].FromAccountID)
	assert.Equal(t, msg.ToAccountID, messages[0].ToAccountID)
	assert.Equal(t, false, messages[0].IsRead)
}

func TestControllerMarkMessagesAsRead(t *testing.T) {
	testConfig.PrepareDB()

	dialogID := CreateDialog(t)

	conn1, err := CreateWsConnection(t, "user")
	assert.NoError(t, err)
	defer conn1.Close()

	conn2, err := CreateWsConnection(t, "support")
	assert.NoError(t, err)
	defer conn2.Close()

	msg := model.DialogWsMessage{
		DialogID:      dialogID,
		FromAccountID: 2,
		ToAccountID:   1,
		Text:          "привет от саппорта",
	}
	err = conn2.WriteJSON(msg)
	assert.NoError(t, err)

	var receivedMsg model.DialogWsMessage
	err = conn1.ReadJSON(&receivedMsg)
	assert.NoError(t, err)

	err = testConfig.DialogClient.MarkMessagesAsRead(dialogID, "user")
	assert.NoError(t, err)

	messages, err := testConfig.DialogClient.MessagesByDialogID(dialogID, "support")
	assert.NoError(t, err)

	assert.Equal(t, true, messages[0].IsRead)
}
