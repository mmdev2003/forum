package unit

import (
	"forum-dialog/internal/model"
	"github.com/stretchr/testify/assert"
	"path/filepath"
	"strings"
	"testing"
)

func TestControllerCreateDialog(t *testing.T) {
	testConfig.PrepareDB()

	response, err := testConfig.dialogClient.CreateDialog(2, "user1")
	assert.NoError(t, err)

	dialog, err := testConfig.dialogClient.DialogsByAccountID("user1")
	assert.NoError(t, err)

	assert.Equal(t, response.DialogID, dialog.Dialogs[0].ID)
	assert.Equal(t, 1, dialog.Dialogs[0].Account1ID)
	assert.Equal(t, 2, dialog.Dialogs[0].Account2ID)
	assert.Equal(t, false, dialog.Dialogs[0].IsStarredByAccount1)
	assert.Equal(t, false, dialog.Dialogs[0].IsStarredByAccount2)
}

func TestControllerDeleteDialog(t *testing.T) {
	testConfig.PrepareDB()

	response, err := testConfig.dialogClient.CreateDialog(2, "user1")
	assert.NoError(t, err)

	err = testConfig.dialogClient.DeleteDialog(response.DialogID)
	assert.NoError(t, err)

	dialog, err := testConfig.dialogClient.DialogsByAccountID("user1")
	assert.NoError(t, err)

	assert.Equal(t, 0, len(dialog.Dialogs))
}

func TestWebSocketMessageExchange(t *testing.T) {
	testConfig.PrepareDB()

	response, err := testConfig.dialogClient.CreateDialog(2, "user1")

	conn1, err := CreateWsConnection(t, "user1")
	assert.NoError(t, err)
	defer conn1.Close()

	conn2, err := CreateWsConnection(t, "user2")
	assert.NoError(t, err)
	defer conn2.Close()

	msg := model.DialogWsMessage{
		DialogID:      response.DialogID,
		FromAccountID: 2,
		ToAccountID:   1,
		Text:          "привет от клиента 2",
	}
	err = conn2.WriteJSON(msg)
	assert.NoError(t, err)

	var receivedMsg model.DialogWsMessage
	err = conn1.ReadJSON(&receivedMsg)
	assert.NoError(t, err)

	assert.Equal(t, msg.DialogID, receivedMsg.DialogID)
	assert.Equal(t, msg.FromAccountID, receivedMsg.FromAccountID)
	assert.Equal(t, msg.ToAccountID, receivedMsg.ToAccountID)
	assert.Equal(t, msg.Text, receivedMsg.Text)

	messages, err := testConfig.dialogClient.MessagesByDialogID(response.DialogID, "user1")
	assert.NoError(t, err)

	assert.Equal(t, 1, len(messages.Messages))
	assert.Equal(t, msg.Text, messages.Messages[0].MessageText)
	assert.Equal(t, msg.FromAccountID, messages.Messages[0].FromAccountID)
	assert.Equal(t, msg.ToAccountID, messages.Messages[0].ToAccountID)
	assert.Equal(t, false, messages.Messages[0].IsRead)
}

func TestWebSocketMessageExchangeWithFile(t *testing.T) {
	testConfig.PrepareDB()

	response, err := testConfig.dialogClient.CreateDialog(2, "user1")

	conn1, err := CreateWsConnection(t, "user1")
	assert.NoError(t, err)
	defer conn1.Close()

	conn2, err := CreateWsConnection(t, "user2")
	assert.NoError(t, err)
	defer conn2.Close()

	files := [][]byte{[]byte("file111"), []byte("file2")}
	filesNames := []string{"file11.jpg", "file22.jpg"}
	var filesURLs []string
	for i, _ := range files {
		uploadResponse, err := testConfig.dialogClient.UploadFile(files[i], filesNames[i], "user1")
		assert.NoError(t, err)
		filesURLs = append(filesURLs, uploadResponse.FileURL)
	}
	messages, err := testConfig.dialogClient.MessagesByDialogID(response.DialogID, "user1")
	assert.NoError(t, err)
	assert.Equal(t, 0, len(messages.Files))

	msg := model.DialogWsMessage{
		DialogID:      response.DialogID,
		FromAccountID: 2,
		ToAccountID:   1,
		Text:          "привет от клиента 2",
		FilesURLs:     filesURLs,
	}
	err = conn2.WriteJSON(msg)
	assert.NoError(t, err)

	var receivedMsg model.DialogWsMessage
	err = conn1.ReadJSON(&receivedMsg)
	assert.NoError(t, err)

	messages, err = testConfig.dialogClient.MessagesByDialogID(response.DialogID, "user1")
	assert.NoError(t, err)
	assert.Equal(t, 2, len(messages.Files))

	for i, _ := range messages.Files {
		assert.Equal(t, filesURLs[i], messages.Files[i].URL)

		extension := filepath.Ext(filesNames[i])
		name := strings.TrimSuffix(filesNames[i], extension)

		assert.Equal(t, name, messages.Files[i].Name)
		assert.Equal(t, extension, messages.Files[i].Extension)
		assert.Equal(t, len(files[i]), messages.Files[i].Size)

		file, err := testConfig.dialogClient.DownloadFile(messages.Files[i].URL)
		assert.NoError(t, err)
		assert.Equal(t, files[i], file)
	}
}

func TestControllerMarkDialogAsStarred(t *testing.T) {
	testConfig.PrepareDB()

	response, err := testConfig.dialogClient.CreateDialog(2, "user1")
	assert.NoError(t, err)

	err = testConfig.dialogClient.MarkDialogAsStarred(response.DialogID, "user1")
	assert.NoError(t, err)

	dialog, err := testConfig.dialogClient.DialogsByAccountID("user1")
	assert.NoError(t, err)

	assert.Equal(t, true, dialog.Dialogs[0].IsStarredByAccount1)

	err = testConfig.dialogClient.MarkDialogAsStarred(response.DialogID, "user1")
	assert.NoError(t, err)

	dialog, err = testConfig.dialogClient.DialogsByAccountID("user1")
	assert.NoError(t, err)

	assert.Equal(t, false, dialog.Dialogs[0].IsStarredByAccount1)

	err = testConfig.dialogClient.MarkDialogAsStarred(response.DialogID, "user2")
	assert.NoError(t, err)

	dialog, err = testConfig.dialogClient.DialogsByAccountID("user2")
	assert.NoError(t, err)

	assert.Equal(t, true, dialog.Dialogs[0].IsStarredByAccount2)

	err = testConfig.dialogClient.MarkDialogAsStarred(response.DialogID, "user2")
	assert.NoError(t, err)

	dialog, err = testConfig.dialogClient.DialogsByAccountID("user2")
	assert.NoError(t, err)

	assert.Equal(t, false, dialog.Dialogs[0].IsStarredByAccount2)
}

func TestControllerMarkMessagesAsRead(t *testing.T) {
	testConfig.PrepareDB()

	response, err := testConfig.dialogClient.CreateDialog(2, "user1")

	conn1, err := CreateWsConnection(t, "user1")
	assert.NoError(t, err)
	defer conn1.Close()

	conn2, err := CreateWsConnection(t, "user2")
	assert.NoError(t, err)
	defer conn2.Close()

	msg := model.DialogWsMessage{
		DialogID:      response.DialogID,
		FromAccountID: 2,
		ToAccountID:   1,
		Text:          "привет от клиента 2",
	}
	err = conn2.WriteJSON(msg)
	assert.NoError(t, err)

	var receivedMsg model.DialogWsMessage
	err = conn1.ReadJSON(&receivedMsg)
	assert.NoError(t, err)

	err = testConfig.dialogClient.MarkMessagesAsRead(response.DialogID, "user1")
	assert.NoError(t, err)

	messages, err := testConfig.dialogClient.MessagesByDialogID(response.DialogID, "user1")
	assert.NoError(t, err)

	assert.Equal(t, true, messages.Messages[0].IsRead)
}
