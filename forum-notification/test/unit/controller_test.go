package unit

import (
	"encoding/json"
	"forum-notification/internal/api/http/handler/notification"
	"forum-notification/internal/model"
	service "forum-notification/internal/service/notification"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCreateLikeMessageFromTopicNotification(t *testing.T) {
	testConfig.PrepareDB()
	testConfig.notificationClient.WebsocketConnect("user1")
	defer testConfig.notificationClient.WebsocketDisconnect()

	expectedNotification := notification.NewLikeMessageFromTopicNotificationRequest{
		MessageOwnerAccountID: 1,
		SenderMessageID:       100,
		LikerAccountID:        2,
		TopicID:               10,
		SenderMessageText:     "Test message",
		TopicName:             "Test topic",
		LikerLogin:            "user2",
	}
	notificationID, err := testConfig.notificationClient.CreateNotification(
		model.LikeMessageFromTopicType, expectedNotification,
	)

	assert.NoError(t, err)

	notifications, err := testConfig.notificationClient.GetNotifications("user1")
	assert.NoError(t, err)
	assert.NotEmpty(t, notifications)

	var actualNotification model.NewLikeMessageFromTopicNotification
	err = json.Unmarshal(notifications[0], &actualNotification)
	assert.NoError(t, err)

	assert.Equal(t, notificationID, actualNotification.ID)
	assert.Equal(t, model.LikeMessageFromTopicType, actualNotification.Type)

	assert.Equal(t, expectedNotification.SenderMessageID, actualNotification.SenderMessageID)
	assert.Equal(t, expectedNotification.LikerAccountID, actualNotification.LikerAccountID)
	assert.Equal(t, expectedNotification.TopicID, actualNotification.TopicID)
	assert.Equal(t, expectedNotification.SenderMessageText, actualNotification.SenderMessageText)
	assert.Equal(t, expectedNotification.TopicName, actualNotification.TopicName)
	assert.Equal(t, expectedNotification.LikerLogin, actualNotification.LikerLogin)
}

func TestCreateMessageFromTopicNotification(t *testing.T) {
	testConfig.PrepareDB()
	testConfig.notificationClient.WebsocketConnect("user1")
	defer testConfig.notificationClient.WebsocketDisconnect()

	expectedNotification := notification.NewMessageFromTopicNotificationRequest{
		TopicOwnerAccountID: 1,
		SenderMessageID:     100,
		SenderAccountID:     2,
		TopicID:             10,
		SenderMessageText:   "Test message",
		TopicName:           "Test topic",
		SenderLogin:         "user2",
	}
	notificationID, err := testConfig.notificationClient.CreateNotification(
		model.MessageFromTopicType, expectedNotification,
	)

	assert.NoError(t, err)

	notifications, err := testConfig.notificationClient.GetNotifications("user1")
	assert.NoError(t, err)
	assert.NotEmpty(t, notifications)

	var actualNotification model.NewMessageFromTopicNotification
	err = json.Unmarshal(notifications[0], &actualNotification)
	assert.NoError(t, err)

	assert.Equal(t, notificationID, actualNotification.ID)
	assert.Equal(t, model.MessageFromTopicType, actualNotification.Type)

	assert.Equal(t, expectedNotification.SenderMessageID, actualNotification.SenderMessageID)
	assert.Equal(t, expectedNotification.SenderAccountID, actualNotification.SenderAccountID)
	assert.Equal(t, expectedNotification.TopicID, actualNotification.TopicID)
	assert.Equal(t, expectedNotification.SenderMessageText, actualNotification.SenderMessageText)
	assert.Equal(t, expectedNotification.TopicName, actualNotification.TopicName)
	assert.Equal(t, expectedNotification.SenderLogin, actualNotification.SenderLogin)
}

func TestCreateMessageReplyFromTopicNotification(t *testing.T) {
	testConfig.PrepareDB()
	testConfig.notificationClient.WebsocketConnect("user1")
	defer testConfig.notificationClient.WebsocketDisconnect()

	expectedNotification := notification.NewMessageReplyFromTopicNotificationRequest{
		ReplyMessageOwnerAccountID: 1,
		SenderMessageID:            100,
		SenderAccountID:            2,
		TopicID:                    10,
		SenderMessageText:          "Test message",
		TopicName:                  "Test topic",
		SenderLogin:                "user2",
	}
	notificationID, err := testConfig.notificationClient.CreateNotification(
		model.MessageReplyFromTopicType, expectedNotification,
	)

	assert.NoError(t, err)

	notifications, err := testConfig.notificationClient.GetNotifications("user1")
	assert.NoError(t, err)
	assert.NotEmpty(t, notifications)

	var actualNotification model.NewMessageReplyFromTopicNotification
	err = json.Unmarshal(notifications[0], &actualNotification)
	assert.NoError(t, err)

	assert.Equal(t, notificationID, actualNotification.ID)
	assert.Equal(t, model.MessageReplyFromTopicType, actualNotification.Type)

	assert.Equal(t, expectedNotification.SenderMessageID, actualNotification.SenderMessageID)
	assert.Equal(t, expectedNotification.SenderAccountID, actualNotification.SenderAccountID)
	assert.Equal(t, expectedNotification.TopicID, actualNotification.TopicID)
	assert.Equal(t, expectedNotification.SenderMessageText, actualNotification.SenderMessageText)
	assert.Equal(t, expectedNotification.TopicName, actualNotification.TopicName)
	assert.Equal(t, expectedNotification.SenderLogin, actualNotification.SenderLogin)
}

func TestCreateTopicClosedNotification(t *testing.T) {
	testConfig.PrepareDB()
	testConfig.notificationClient.WebsocketConnect("user1")
	defer testConfig.notificationClient.WebsocketDisconnect()

	expectedNotification := notification.TopicClosedNotificationRequest{
		TopicOwnerAccountID: 1,
		AdminAccountID:      2,
		TopicID:             10,
		TopicName:           "Test topic",
		AdminLogin:          "user2",
	}
	notificationID, err := testConfig.notificationClient.CreateNotification(
		model.TopicClosedType, expectedNotification,
	)

	assert.NoError(t, err)

	notifications, err := testConfig.notificationClient.GetNotifications("user1")
	assert.NoError(t, err)
	assert.NotEmpty(t, notifications)

	var actualNotification model.TopicClosedNotification
	err = json.Unmarshal(notifications[0], &actualNotification)
	assert.NoError(t, err)

	assert.Equal(t, notificationID, actualNotification.ID)
	assert.Equal(t, model.TopicClosedType, actualNotification.Type)

	assert.Equal(t, expectedNotification.AdminAccountID, actualNotification.AdminAccountID)
	assert.Equal(t, expectedNotification.TopicID, actualNotification.TopicID)
	assert.Equal(t, expectedNotification.TopicName, actualNotification.TopicName)
	assert.Equal(t, expectedNotification.AdminLogin, actualNotification.AdminLogin)
}

func TestCreateResponseToSupportRequestNotification(t *testing.T) {
	testConfig.PrepareDB()
	testConfig.notificationClient.WebsocketConnect("user1")
	defer testConfig.notificationClient.WebsocketDisconnect()

	expectedNotification := notification.ResponseToSupportRequestNotificationRequest{
		RequesterAccountID: 1,
		SupportRequestID:   2,
	}
	notificationID, err := testConfig.notificationClient.CreateNotification(
		model.ResponseToSupportRequestType, expectedNotification,
	)

	assert.NoError(t, err)

	notifications, err := testConfig.notificationClient.GetNotifications("user1")
	assert.NoError(t, err)
	assert.NotEmpty(t, notifications)

	var actualNotification model.ResponseToSupportRequestNotification
	err = json.Unmarshal(notifications[0], &actualNotification)
	assert.NoError(t, err)

	assert.Equal(t, notificationID, actualNotification.ID)
	assert.Equal(t, model.ResponseToSupportRequestType, actualNotification.Type)

	assert.Equal(t, expectedNotification.SupportRequestID, actualNotification.SupportRequestID)
}

func TestCreateStatusReceivedNotification(t *testing.T) {
	testConfig.PrepareDB()
	testConfig.notificationClient.WebsocketConnect("user1")
	defer testConfig.notificationClient.WebsocketDisconnect()

	expectedNotification := notification.StatusReceivedNotificationRequest{
		ReceiverAccountID: 1,
		StatusName:        "Test status",
	}
	notificationID, err := testConfig.notificationClient.CreateNotification(
		model.StatusReceivedType, expectedNotification,
	)

	assert.NoError(t, err)

	notifications, err := testConfig.notificationClient.GetNotifications("user1")
	assert.NoError(t, err)
	assert.NotEmpty(t, notifications)

	var actualNotification model.StatusReceivedNotification
	err = json.Unmarshal(notifications[0], &actualNotification)
	assert.NoError(t, err)

	assert.Equal(t, notificationID, actualNotification.ID)
	assert.Equal(t, model.StatusReceivedType, actualNotification.Type)

	assert.Equal(t, expectedNotification.StatusName, actualNotification.StatusName)
}

func TestCreateFrameReceivedNotification(t *testing.T) {
	testConfig.PrepareDB()
	testConfig.notificationClient.WebsocketConnect("user1")
	defer testConfig.notificationClient.WebsocketDisconnect()

	expectedNotification := notification.FrameReceivedNotificationRequest{
		ReceiverAccountID: 1,
		FrameName:         "Test Frame",
	}
	notificationID, err := testConfig.notificationClient.CreateNotification(
		model.FrameReceivedType, expectedNotification,
	)

	assert.NoError(t, err)

	notifications, err := testConfig.notificationClient.GetNotifications("user1")
	assert.NoError(t, err)
	assert.NotEmpty(t, notifications)

	var actualNotification model.FrameReceivedNotification
	err = json.Unmarshal(notifications[0], &actualNotification)
	assert.NoError(t, err)

	assert.Equal(t, notificationID, actualNotification.ID)
	assert.Equal(t, model.FrameReceivedType, actualNotification.Type)

	assert.Equal(t, expectedNotification.FrameName, actualNotification.FrameName)
}

func TestCreateMessageFromDialogNotification(t *testing.T) {
	testConfig.PrepareDB()
	testConfig.notificationClient.WebsocketConnect("user1")
	defer testConfig.notificationClient.WebsocketDisconnect()

	expectedNotification := notification.NewMessageFromDialogNotificationRequest{
		AccountID:       1,
		MessageID:       100,
		SenderAccountID: 2,
		DialogID:        10,
		MessageText:     "Test message",
		SenderLogin:     "user2",
	}
	notificationID, err := testConfig.notificationClient.CreateNotification(
		model.MessageFromDialogType, expectedNotification,
	)

	assert.NoError(t, err)

	notifications, err := testConfig.notificationClient.GetNotifications("user1")
	assert.NoError(t, err)
	assert.NotEmpty(t, notifications)

	var actualNotification model.NewMessageFromDialogNotification
	err = json.Unmarshal(notifications[0], &actualNotification)
	assert.NoError(t, err)

	assert.Equal(t, notificationID, actualNotification.ID)
	assert.Equal(t, model.MessageFromDialogType, actualNotification.Type)

	assert.Equal(t, expectedNotification.MessageID, actualNotification.MessageID)
	assert.Equal(t, expectedNotification.SenderAccountID, actualNotification.SenderAccountID)
	assert.Equal(t, expectedNotification.DialogID, actualNotification.DialogID)
	assert.Equal(t, expectedNotification.MessageText, actualNotification.MessageText)
	assert.Equal(t, expectedNotification.SenderLogin, actualNotification.SenderLogin)
}

func TestCreateMentionFromTopicNotification(t *testing.T) {
	testConfig.PrepareDB()
	testConfig.notificationClient.WebsocketConnect("user1")
	defer testConfig.notificationClient.WebsocketDisconnect()

	expectedNotification := notification.NewMentionFromTopicNotificationRequest{
		MentionedAccountID: 1,
		MessageID:          100,
		SenderAccountID:    2,
		SenderMessageText:  "Test message",
		TopicName:          "Test topic",
		SenderLogin:        "user2",
	}
	notificationID, err := testConfig.notificationClient.CreateNotification(
		model.MentionFromTopicType, expectedNotification,
	)

	assert.NoError(t, err)

	notifications, err := testConfig.notificationClient.GetNotifications("user1")
	assert.NoError(t, err)
	assert.NotEmpty(t, notifications)

	var actualNotification model.NewMentionFromTopicNotification
	err = json.Unmarshal(notifications[0], &actualNotification)
	assert.NoError(t, err)

	assert.Equal(t, notificationID, actualNotification.ID)
	assert.Equal(t, model.MentionFromTopicType, actualNotification.Type)

	assert.Equal(t, expectedNotification.MessageID, actualNotification.MessageID)
	assert.Equal(t, expectedNotification.SenderAccountID, actualNotification.SenderAccountID)
	assert.Equal(t, expectedNotification.SenderMessageText, actualNotification.SenderMessageText)
	assert.Equal(t, expectedNotification.TopicName, actualNotification.TopicName)
	assert.Equal(t, expectedNotification.SenderLogin, actualNotification.SenderLogin)
}

func TestCreateWarningFromAdminNotification(t *testing.T) {
	testConfig.PrepareDB()
	testConfig.notificationClient.WebsocketConnect("user1")
	defer testConfig.notificationClient.WebsocketDisconnect()

	expectedNotification := notification.NewWarningFromAdminNotificationRequest{
		AccountID:      1,
		AdminAccountID: 2,
		WarningText:    "Test warning",
		AdminLogin:     "user2",
	}
	notificationID, err := testConfig.notificationClient.CreateNotification(
		model.WarningFromAdminType, expectedNotification,
	)

	assert.NoError(t, err)

	notifications, err := testConfig.notificationClient.GetNotifications("user1")
	assert.NoError(t, err)
	assert.NotEmpty(t, notifications)

	var actualNotification model.NewWarningFromAdminNotification
	err = json.Unmarshal(notifications[0], &actualNotification)
	assert.NoError(t, err)

	assert.Equal(t, notificationID, actualNotification.ID)
	assert.Equal(t, model.WarningFromAdminType, actualNotification.Type)

	assert.Equal(t, expectedNotification.AdminAccountID, actualNotification.AdminAccountID)
	assert.Equal(t, expectedNotification.WarningText, actualNotification.WarningText)
	assert.Equal(t, expectedNotification.AdminLogin, actualNotification.AdminLogin)
}

func TestWebsocketNotifications(t *testing.T) {
	testConfig.PrepareDB()
	testConfig.notificationClient.WebsocketConnect("user1")
	defer testConfig.notificationClient.WebsocketDisconnect()

	mentionNotification := notification.NewMentionFromTopicNotificationRequest{
		MentionedAccountID: 1,
		MessageID:          101,
		SenderAccountID:    2,
		SenderMessageText:  "Hey @user1, check this",
		TopicName:          "Mention test",
		SenderLogin:        "user2",
	}
	_, err := testConfig.notificationClient.CreateNotification(
		model.MentionFromTopicType, mentionNotification,
	)
	assert.NoError(t, err)

	msg, err := testConfig.notificationClient.WebsocketReadMessage()
	assert.NoError(t, err)
	assert.NotEmpty(t, msg)

	var wsNotification service.NotificationWsBody
	var payload service.NewMentionFromTopicNotificationPayload

	err = json.Unmarshal(msg, &wsNotification)
	assert.NoError(t, err)
	assert.Equal(t, model.MentionFromTopicType, wsNotification.Type)

	err = json.Unmarshal(wsNotification.Payload, &payload)
	assert.NoError(t, err)
	assert.Equal(t, mentionNotification.MessageID, payload.MessageID)
	assert.Equal(t, mentionNotification.SenderMessageText, payload.MessageText)
	assert.Equal(t, mentionNotification.SenderLogin, payload.MentionLogin)
	assert.Equal(t, mentionNotification.SenderAccountID, payload.MentionAccountID)
	assert.Equal(t, mentionNotification.TopicName, payload.TopicName)
}

func TestSetupFilters(t *testing.T) {
	testConfig.PrepareDB()
	testConfig.notificationClient.WebsocketConnect("user1")
	defer testConfig.notificationClient.WebsocketDisconnect()

	expectedFilters := []model.NotificationType{
		model.MessageFromTopicType,
		model.MessageReplyFromTopicType,
		model.LikeMessageFromTopicType,
		model.TopicClosedType,
		model.ResponseToSupportRequestType,
		model.StatusReceivedType,
		model.FrameReceivedType,
		model.MessageFromDialogType,
		model.MentionFromTopicType,
		model.WarningFromAdminType,
	}

	actualFilters, err := testConfig.notificationClient.GetFilters("user1")
	assert.NoError(t, err)
	assert.ElementsMatch(t, expectedFilters, actualFilters.EnabledFilters)

	expectedFilters = []model.NotificationType{
		model.LikeMessageFromTopicType,
		model.TopicClosedType,
		model.MentionFromTopicType,
	}

	err = testConfig.notificationClient.SetFilters("user1", expectedFilters)
	assert.NoError(t, err, err)

	actualFilters, err = testConfig.notificationClient.GetFilters("user1")
	assert.NoError(t, err, err)
	assert.ElementsMatch(t, expectedFilters, actualFilters.EnabledFilters)
}
