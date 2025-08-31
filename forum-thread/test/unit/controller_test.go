package unit

import (
	"context"
	"forum-thread/internal/model"
	"github.com/stretchr/testify/assert"
	"sort"
	"strconv"
	"testing"
)

func TestControllerCreateThread(t *testing.T) {
	testConfig.PrepareDB()

	expectedThreads := ControllerCreateThreads(t)

	threads, err := testConfig.threadClient.AllThread()
	assert.NoError(t, err)
	assert.Equal(t, 2, len(threads.Threads))

	for i := range threads.Threads {
		assert.Equal(t, expectedThreads[i].ID, threads.Threads[i].ID)
		assert.Equal(t, expectedThreads[i].ThreadName, threads.Threads[i].ThreadName)
		assert.Equal(t, expectedThreads[i].ThreadDescription, threads.Threads[i].ThreadDescription)
		assert.Equal(t, expectedThreads[i].AllowedStatuses, threads.Threads[i].AllowedStatuses)
		assert.Equal(t, "green", threads.Threads[i].ThreadColor)
	}
}

func TestControllerCreateSubthread(t *testing.T) {
	testConfig.PrepareDB()

	expectedThreads := ControllerCreateThreads(t)
	expectedSubthreads := ControllerCreateSubthreads(t)

	for threadIndex, thread := range expectedThreads {
		subthreads, err := testConfig.threadClient.SubthreadsByThreadID(thread.ID)
		assert.NoError(t, err)
		assert.Equal(t, 2, len(subthreads.Subthreads))

		for i := 0; i < 2; i++ {
			expectedSubthread := expectedSubthreads[threadIndex*2+i]

			assert.Equal(t, expectedSubthread.ID, subthreads.Subthreads[i].ID)
			assert.Equal(t, expectedSubthread.SubthreadName, subthreads.Subthreads[i].SubthreadName)
			assert.Equal(t, expectedSubthread.ThreadName, subthreads.Subthreads[i].ThreadName)
			assert.Equal(t, expectedSubthread.SubthreadDescription, subthreads.Subthreads[i].SubthreadDescription)
		}
	}
}

func TestControllerAddViewToSubthread(t *testing.T) {
	testConfig.PrepareDB()

	expectedThreads := ControllerCreateThreads(t)
	expectedSubthreads := ControllerCreateSubthreads(t)

	for _, expectedSubthread := range expectedSubthreads {
		for _ = range 2 {
			err := testConfig.threadClient.AddViewToSubthread(expectedSubthread.ID)
			assert.NoError(t, err)
			err = testConfig.threadClient.AddViewToSubthreadPostprocessing(expectedSubthread.ID)
			assert.NoError(t, err)
		}
	}

	for _, thread := range expectedThreads {
		subthreads, err := testConfig.threadClient.SubthreadsByThreadID(thread.ID)
		assert.NoError(t, err)

		for _, subthread := range subthreads.Subthreads {
			assert.Equal(t, 2, subthread.SubthreadViewCount)
		}
	}
}

func TestControllerCreateTopic(t *testing.T) {
	testConfig.PrepareDB()

	ControllerCreateAccountStatistic(t)
	_ = ControllerCreateThreads(t)
	_ = ControllerCreateSubthreads(t)
	_ = ControllerCreateTopics(t)
	expectedTopics := ControllerCreateTopicsPostprocessing(t)

	for _, expectedTopic := range expectedTopics {
		err := testConfig.threadClient.ApproveTopic(expectedTopic.ID, "admin1")
		assert.NoError(t, err)
	}

	subthreadIDsWithTopics := []int{1, 3}
	for _, subthreadID := range subthreadIDsWithTopics {

		topics, err := testConfig.threadClient.TopicsBySubthreadID(subthreadID)
		assert.NoError(t, err)
		assert.Equal(t, 2, len(topics.Topics))

		var expectedForSubthread []TestTopic
		for _, topic := range expectedTopics {
			if topic.SubthreadID == subthreadID {
				expectedForSubthread = append(expectedForSubthread, topic)
			}
		}

		for i, topic := range topics.Topics {
			expected := expectedForSubthread[i]
			assert.Equal(t, expected.ID, topic.ID)
			assert.Equal(t, expected.SubthreadID, topic.SubthreadID)
			assert.Equal(t, expected.ThreadID, topic.ThreadID)
			assert.Equal(t, expected.TopicOwnerAccountID, topic.TopicOwnerAccountID)
			assert.Equal(t, expected.SubthreadName, topic.SubthreadName)
			assert.Equal(t, expected.ThreadName, topic.ThreadName)
			assert.Equal(t, expected.TopicName, topic.TopicName)
			assert.Equal(t, expected.TopicOwnerLogin, topic.TopicOwnerLogin)
			assert.Equal(t, model.ApprovedTopicStatus, topic.TopicModerationStatus)
		}
	}

	subthreadIDsWithoutTopics := []int{2, 4}
	for _, subthreadID := range subthreadIDsWithoutTopics {
		topics, err := testConfig.threadClient.TopicsBySubthreadID(subthreadID)
		assert.NoError(t, err)
		assert.Empty(t, topics)
	}

	accountStatistic, err := testConfig.threadClient.StatisticByAccountID(1)
	assert.NoError(t, err)

	assert.Equal(t, 4, accountStatistic.AccountStatistic.CreatedTopicsCount)
}

func TestControllerApproveTopic(t *testing.T) {
	testConfig.PrepareDB()

	ControllerCreateAccountStatistic(t)
	_ = ControllerCreateThread(t)
	_ = ControllerCreateSubthread(t)
	_ = ControllerCreateTopic(t)
	expectedTopic := ControllerCreateTopicPostprocessing(t)

	topics, err := testConfig.threadClient.TopicsByAccountID(1)
	assert.NoError(t, err)

	assert.Equal(t, model.OnModerationTopicStatus, topics.Topics[0].TopicModerationStatus)

	err = testConfig.threadClient.ApproveTopic(expectedTopic.ID, "admin1")
	assert.NoError(t, err)

	topics, err = testConfig.threadClient.TopicsByAccountID(1)
	assert.NoError(t, err)

	assert.Equal(t, model.ApprovedTopicStatus, topics.Topics[0].TopicModerationStatus)
}

func TestControllerCreateTopicWithoutPermissionInThread(t *testing.T) {
	testConfig.PrepareDB()

	ControllerCreateAccountStatistic(t)
	_ = ControllerCreateThread(t)
	_ = ControllerCreateSubthread(t)
	testTopic := TestTopic{1, 1, 1, 1, "subthread1", "thread1", "topic1", "user1"}

	_, err := testConfig.threadClient.CreateTopic(
		testTopic.SubthreadID,
		testTopic.ThreadID,
		testTopic.SubthreadName,
		testTopic.ThreadName,
		testTopic.TopicName,
		testTopic.TopicOwnerLogin,
		"user3",
	)
	assert.Error(t, err)
}

func TestControllerCreateTopicWithoutPermissionInSubthread(t *testing.T) {
	testConfig.PrepareDB()

	ControllerCreateAccountStatistic(t)
	_ = ControllerCreateThread(t)
	_ = ControllerCreateSubthread(t)
	testTopic := TestTopic{1, 1, 1, 1, "subthread1", "thread1", "topic1", "user1"}

	_, err := testConfig.threadClient.CreateTopic(
		testTopic.SubthreadID,
		testTopic.ThreadID,
		testTopic.SubthreadName,
		testTopic.ThreadName,
		testTopic.TopicName,
		testTopic.TopicOwnerLogin,
		"user4",
	)
	assert.Error(t, err)
}

func TestControllerCreateTopicMaxPerDay(t *testing.T) {
	testConfig.PrepareDB()

	ControllerCreateAccountStatistic(t)
	_ = ControllerCreateThread(t)
	_ = ControllerCreateSubthread(t)
	expectedTopics := []TestTopic{
		{1, 1, 1, 1, "subthread1", "thread1", "topic1", "user1"},
		{2, 1, 1, 1, "subthread1", "thread1", "topic2", "user1"},
		{3, 3, 2, 1, "subthread3", "thread2", "topic3", "user1"},
		{4, 3, 2, 1, "subthread3", "thread2", "topic4", "user1"},
	}
	for i, topic := range expectedTopics {
		_, err := testConfig.threadClient.CreateTopic(
			topic.SubthreadID,
			topic.ThreadID,
			topic.SubthreadName,
			topic.ThreadName,
			topic.TopicName,
			topic.TopicOwnerLogin,
			"user1",
		)
		if i <= 1 {
			assert.NoError(t, err)
		} else {
			assert.Error(t, err)
		}
	}
}

func TestControllerRejectTopic(t *testing.T) {
	testConfig.PrepareDB()

	ControllerCreateAccountStatistic(t)
	_ = ControllerCreateThread(t)
	_ = ControllerCreateSubthread(t)
	_ = ControllerCreateTopic(t)
	expectedTopic := ControllerCreateTopicPostprocessing(t)

	err := testConfig.threadClient.RejectTopic(expectedTopic.ID, "admin1")
	assert.NoError(t, err)

	topics, err := testConfig.threadClient.TopicsByAccountID(1)
	assert.NoError(t, err)

	assert.Equal(t, model.RejectedTopicStatus, topics.Topics[0].TopicModerationStatus)
}

func TestControllerTopicsByAccountID(t *testing.T) {
	testConfig.PrepareDB()

	ControllerCreateAccountStatistic(t)
	_ = ControllerCreateThread(t)
	_ = ControllerCreateSubthread(t)
	_ = ControllerCreateTopic(t)
	expectedTopic := ControllerCreateTopicPostprocessing(t)

	topics, err := testConfig.threadClient.TopicsByAccountID(1)
	assert.NoError(t, err)

	assert.Equal(t, 1, len(topics.Topics))
	assert.Equal(t, expectedTopic.ID, topics.Topics[0].ID)
	assert.Equal(t, expectedTopic.SubthreadID, topics.Topics[0].SubthreadID)
	assert.Equal(t, expectedTopic.ThreadID, topics.Topics[0].ThreadID)
	assert.Equal(t, expectedTopic.TopicOwnerAccountID, topics.Topics[0].TopicOwnerAccountID)
	assert.Equal(t, expectedTopic.SubthreadName, topics.Topics[0].SubthreadName)
	assert.Equal(t, expectedTopic.ThreadName, topics.Topics[0].ThreadName)
	assert.Equal(t, expectedTopic.TopicName, topics.Topics[0].TopicName)
	assert.Equal(t, expectedTopic.TopicOwnerLogin, topics.Topics[0].TopicOwnerLogin)
}

func TestControllerAddViewToTopic(t *testing.T) {
	testConfig.PrepareDB()

	_ = ControllerCreateThread(t)
	testSubthread := ControllerCreateSubthread(t)
	testTopic := ControllerCreateTopic(t)

	err := testConfig.threadClient.ApproveTopic(testTopic.ID, "admin1")
	assert.NoError(t, err)

	for _ = range 2 {
		err := testConfig.threadClient.AddViewToTopic(testTopic.ID)
		assert.NoError(t, err)
		err = testConfig.threadClient.AddViewToTopicPostprocessing(testTopic.ID)
		assert.NoError(t, err)
	}

	topics, err := testConfig.threadClient.TopicsBySubthreadID(testSubthread.ID)
	assert.NoError(t, err)

	assert.Equal(t, 2, topics.Topics[0].TopicViewCount)
}

func TestControllerCloseTopic(t *testing.T) {
	testConfig.PrepareDB()

	_ = ControllerCreateThread(t)
	testSubthread := ControllerCreateSubthread(t)
	testTopic := ControllerCreateTopic(t)

	err := testConfig.threadClient.ApproveTopic(testTopic.ID, "admin1")
	assert.NoError(t, err)

	topics, err := testConfig.threadClient.TopicsBySubthreadID(testSubthread.ID)
	assert.NoError(t, err)

	assert.Equal(t, false, topics.Topics[0].TopicIsClosed)

	err = testConfig.threadClient.CloseTopic(
		1,
		1,
		testTopic.ID,
		"topic",
		"admin",
		"admin1",
	)
	assert.NoError(t, err)

	err = testConfig.threadClient.CloseTopicPostprocessing(testTopic.ID)
	assert.NoError(t, err)

	topics, err = testConfig.threadClient.TopicsBySubthreadID(testSubthread.ID)
	assert.NoError(t, err)

	assert.Equal(t, true, topics.Topics[0].TopicIsClosed)
}

func TestControllerChangeTopicPriority(t *testing.T) {
	_ = context.Background()
}

func TestControllerSendMessagesToTopic(t *testing.T) {
	testConfig.PrepareDB()

	err := testConfig.threadClient.CreateAccountStatistic(1)
	assert.NoError(t, err)

	_ = ControllerCreateThreads(t)
	_ = ControllerCreateSubthreads(t)
	_ = ControllerCreateTopics(t)
	expectedTopics := ControllerCreateTopicsPostprocessing(t)
	_ = ControllerSendMessagesToTopic(t)
	expectedMessages := ControllerSendMessagesToTopicPostprocessing(t)

	for _, expectedTopic := range expectedTopics {
		err := testConfig.threadClient.ApproveTopic(expectedTopic.ID, "admin1")
		assert.NoError(t, err)
	}

	topicIDsWithTMessages := []int{1, 2}
	for _, topicID := range topicIDsWithTMessages {
		response, err := testConfig.threadClient.MessagesByTopicID(topicID, "user1")
		assert.NoError(t, err)
		assert.Equal(t, 2, len(response.Messages))

		var expectedForTopic []TestMessage
		for _, message := range expectedMessages {
			if message.TopicID == topicID {
				expectedForTopic = append(expectedForTopic, message)
			}
		}

		for i, message := range response.Messages {
			expected := expectedForTopic[i]
			assert.Equal(t, expected.ID, message.ID)
			assert.Equal(t, expected.TopicID, message.TopicID)
			assert.Equal(t, expected.MessageOwnerAccountID, message.MessageOwnerAccountID)
			assert.Equal(t, expected.MessageOwnerLogin, message.MessageOwnerLogin)
			assert.Equal(t, expected.MessageText, message.MessageText)
		}
	}

	accountStatistic, err := testConfig.threadClient.StatisticByAccountID(1)
	assert.NoError(t, err)
	assert.Equal(t, 2, accountStatistic.AccountStatistic.SentMessagesToTopicsCount)

	subthreads, err := testConfig.threadClient.SubthreadsByThreadID(1)
	assert.NoError(t, err)
	assert.Equal(t, 2, subthreads.Subthreads[0].SubthreadMessageCount)
	assert.Equal(t, expectedMessages[3].MessageText, subthreads.Subthreads[0].SubthreadLastMessageText)
	assert.Equal(t, expectedMessages[3].MessageOwnerLogin, subthreads.Subthreads[0].SubthreadLastMessageLogin)

	topics, err := testConfig.threadClient.TopicsBySubthreadID(1)
	assert.NoError(t, err)
	assert.Equal(t, 2, topics.Topics[0].TopicMessageCount)
	assert.Equal(t, expectedMessages[1].MessageText, topics.Topics[0].TopicLastMessageText)
	assert.Equal(t, expectedMessages[1].MessageOwnerLogin, topics.Topics[0].TopicLastMessageLogin)

	messages, err := testConfig.threadClient.MessagesByText("приет")
	assert.NoError(t, err)

	sort.Slice(messages.Messages, func(i, j int) bool {
		return messages.Messages[i].ID < messages.Messages[j].ID
	})

	for i, message := range messages.Messages {
		assert.Equal(t, expectedMessages[i].ID, message.ID)
		assert.Equal(t, expectedMessages[i].TopicID, message.TopicID)
		assert.Equal(t, expectedMessages[i].MessageOwnerAccountID, message.AccountID)
		assert.Equal(t, expectedMessages[i].MessageOwnerLogin, message.Login)
		assert.Equal(t, expectedMessages[i].MessageText, message.Text)
	}

}

func TestControllerSendMessagesToTopicWithFile(t *testing.T) {
	testConfig.PrepareDB()

	err := testConfig.threadClient.CreateAccountStatistic(1)
	assert.NoError(t, err)

	_ = ControllerCreateThread(t)
	_ = ControllerCreateSubthread(t)
	_ = ControllerCreateTopic(t)
	_ = ControllerCreateTopicPostprocessing(t)
	testMessage := TestMessage{1, 1, 1, 1, "user1", "Всем привет"}
	files := [][]byte{[]byte("test_test1"), []byte("test2")}
	filesFullNames := []string{"test1.txt", "test2.txt"}
	filesNames := []string{"test1", "test2"}
	filesExtensions := []string{"txt", "txt"}
	filesSizes := []int{len(files[0]), len(files[1])}

	sendResponse, err := testConfig.threadClient.SendMessageToTopic(
		testMessage.SubthreadID,
		testMessage.TopicID,
		0,
		0,
		1,
		testMessage.MessageOwnerLogin,
		"thread1",
		"subthread1",
		"topic",
		testMessage.MessageText,
		files,
		filesFullNames,
		"user1",
	)
	assert.NoError(t, err)

	err = testConfig.threadClient.SendMessageToTopicPostprocessing(
		testMessage.SubthreadID,
		testMessage.TopicID,
		0,
		1,
		1,
		testMessage.MessageOwnerAccountID,
		testMessage.MessageOwnerLogin,
		"topic",
		testMessage.MessageText,
		sendResponse.FilesURLs,
		filesNames,
		filesExtensions,
		filesSizes,
	)
	assert.NoError(t, err)

	messages, err := testConfig.threadClient.MessagesByTopicID(testMessage.TopicID, "user2")
	assert.NoError(t, err)
	assert.Equal(t, 2, len(messages.Files))

	for i := range files {
		assert.Equal(t, len(files[i]), messages.Files[i].Size)
		assert.Equal(t, filesNames[i], messages.Files[i].Name)
		assert.Equal(t, filesExtensions[i], messages.Files[i].Extension)
		assert.Equal(t, sendResponse.FilesURLs[i], messages.Files[i].URL)
		assert.Equal(t, testMessage.ID, messages.Files[i].MessageID)
	}

	file1, err := testConfig.threadClient.DownloadFile(messages.Files[0].URL)
	assert.NoError(t, err)
	assert.Equal(t, files[0], file1)

	file2, err := testConfig.threadClient.DownloadFile(messages.Files[1].URL)
	assert.NoError(t, err)
	assert.Equal(t, files[1], file2)
}

func TestControllerLikeMessage(t *testing.T) {
	testConfig.PrepareDB()

	ControllerCreateAccountStatistic(t)
	_ = ControllerCreateThread(t)
	_ = ControllerCreateSubthread(t)
	testTopic := ControllerCreateTopic(t)
	_ = ControllerCreateTopicPostprocessing(t)
	_ = ControllerSendMessageToTopic(t)
	testMessage := ControllerSendMessageToTopicPostprocessing(t)

	accountIDs := []int{1, 2}
	for _, likerAccountID := range accountIDs {
		err := testConfig.threadClient.LikeMessage(
			testTopic.ID,
			testMessage.MessageOwnerAccountID,
			likerAccountID,
			testMessage.ID,
			1,
			"user",
			"topic",
			"message",
		)
		assert.NoError(t, err)
		err = testConfig.threadClient.LikeMessagePostprocessing(
			testTopic.ID,
			testMessage.MessageOwnerAccountID,
			likerAccountID,
			testMessage.ID,
			1,
			"user",
			"topic",
			"message",
		)
		assert.NoError(t, err)
	}

	for _ = range accountIDs {
		response, err := testConfig.threadClient.MessagesByTopicID(testTopic.ID, "user1")
		assert.NoError(t, err)
		assert.Equal(t, 1, len(response.Likes))
		assert.Equal(t, 2, response.Messages[0].MessageLikeCount)

		for _, like := range response.Likes {
			assert.Equal(t, testTopic.ID, like.TopicID)
			assert.Equal(t, testMessage.ID, like.MessageID)
			assert.Equal(t, 1, like.LikeTypeID)
		}
	}
}

func TestControllerUnlikeMessage(t *testing.T) {
	testConfig.PrepareDB()

	ControllerCreateAccountStatistic(t)
	_ = ControllerCreateThread(t)
	_ = ControllerCreateSubthread(t)
	testTopic := ControllerCreateTopic(t)
	_ = ControllerCreateTopicPostprocessing(t)
	_ = ControllerSendMessageToTopic(t)
	testMessage := ControllerSendMessageToTopicPostprocessing(t)

	accountIDs := []int{1, 2}
	for _, likerAccountID := range accountIDs {
		err := testConfig.threadClient.LikeMessage(
			testTopic.ID,
			testMessage.MessageOwnerAccountID,
			likerAccountID,
			testMessage.ID,
			1,
			"user",
			"topic",
			"message",
		)
		assert.NoError(t, err)
		err = testConfig.threadClient.LikeMessagePostprocessing(
			testTopic.ID,
			testMessage.MessageOwnerAccountID,
			likerAccountID,
			testMessage.ID,
			1,
			"user",
			"topic",
			"message",
		)
		assert.NoError(t, err)

		err = testConfig.threadClient.UnlikeMessage(
			testMessage.ID,
			"user"+strconv.Itoa(likerAccountID),
		)
		assert.NoError(t, err)
	}

	for _ = range accountIDs {
		response, err := testConfig.threadClient.MessagesByTopicID(testTopic.ID, "user1")
		assert.NoError(t, err)
		assert.Equal(t, 0, len(response.Likes))

		assert.Equal(t, 0, response.Messages[0].MessageLikeCount)
	}
}

func TestControllerReportMessage(t *testing.T) {
	testConfig.PrepareDB()

	ControllerCreateAccountStatistic(t)
	_ = ControllerCreateThread(t)
	_ = ControllerCreateSubthread(t)
	testTopic := ControllerCreateTopic(t)
	_ = ControllerCreateTopicPostprocessing(t)
	_ = ControllerSendMessageToTopic(t)
	testMessage := ControllerSendMessageToTopicPostprocessing(t)

	accountIDs := []int{1, 2, 3}
	for _, accountID := range accountIDs {
		err := testConfig.threadClient.ReportMessage(testMessage.ID, accountID, "report")
		assert.NoError(t, err)
		err = testConfig.threadClient.ReportMessagePostprocessing(testMessage.ID, accountID, "report")
		assert.NoError(t, err)
	}

	response, err := testConfig.threadClient.MessagesByTopicID(testTopic.ID, "user1")
	assert.NoError(t, err)

	assert.Equal(t, 3, response.Messages[0].MessageReportCount)
}

func TestControllerSendMessageToTopicWithoutPermissionInThread(t *testing.T) {
	testConfig.PrepareDB()

	ControllerCreateAccountStatistic(t)
	_ = ControllerCreateThread(t)
	_ = ControllerCreateSubthread(t)
	_ = ControllerCreateTopic(t)
	_ = ControllerCreateTopicPostprocessing(t)
	testMessage := TestMessage{1, 1, 1, 1, "user1", "Всем привет"}

	_, err := testConfig.threadClient.SendMessageToTopic(
		testMessage.SubthreadID,
		testMessage.TopicID,
		0,
		1,
		1,
		testMessage.MessageOwnerLogin,
		"thread1",
		"subthread1",
		"topic1",
		testMessage.MessageText,
		nil,
		nil,
		"user3",
	)
	assert.Error(t, err)
}

func TestControllerSendMessageToTopicWithoutPermissionInSubthread(t *testing.T) {
	testConfig.PrepareDB()

	ControllerCreateAccountStatistic(t)
	_ = ControllerCreateThread(t)
	_ = ControllerCreateSubthread(t)
	_ = ControllerCreateTopic(t)
	_ = ControllerCreateTopicPostprocessing(t)
	testMessage := TestMessage{1, 1, 1, 1, "user1", "Всем привет"}

	_, err := testConfig.threadClient.SendMessageToTopic(
		testMessage.SubthreadID,
		testMessage.TopicID,
		0,
		1,
		1,
		testMessage.MessageOwnerLogin,
		"thread1",
		"subthread1",
		"topic1",
		testMessage.MessageText,
		nil,
		nil,
		"user4",
	)
	assert.Error(t, err)
}

func TestControllerEditMessage(t *testing.T) {
	testConfig.PrepareDB()

	ControllerCreateAccountStatistic(t)
	_ = ControllerCreateThread(t)
	_ = ControllerCreateSubthread(t)
	testTopic := ControllerCreateTopic(t)
	_ = ControllerCreateTopicPostprocessing(t)
	_ = ControllerSendMessageToTopic(t)
	testMessage := ControllerSendMessageToTopicPostprocessing(t)

	err := testConfig.threadClient.EditMessage(testMessage.ID, "report")
	assert.NoError(t, err)

	response, err := testConfig.threadClient.MessagesByTopicID(testTopic.ID, "user1")
	assert.NoError(t, err)

	assert.Equal(t, "report", response.Messages[0].MessageText)
}

func TestControllerMessagesByAccountID(t *testing.T) {
	testConfig.PrepareDB()

	ControllerCreateAccountStatistic(t)
	_ = ControllerCreateThread(t)
	_ = ControllerCreateSubthread(t)
	_ = ControllerCreateTopic(t)
	_ = ControllerCreateTopicPostprocessing(t)
	_ = ControllerSendMessageToTopic(t)
	testMessage := ControllerSendMessageToTopicPostprocessing(t)

	messages, err := testConfig.threadClient.MessagesByAccountID(1)
	assert.NoError(t, err)

	assert.Equal(t, 1, len(messages.Messages))
	assert.Equal(t, testMessage.ID, messages.Messages[0].ID)
	assert.Equal(t, testMessage.MessageOwnerLogin, messages.Messages[0].MessageOwnerLogin)
	assert.Equal(t, testMessage.MessageText, messages.Messages[0].MessageText)
	assert.Equal(t, testMessage.MessageOwnerAccountID, messages.Messages[0].MessageOwnerAccountID)
}

func TestControllerCreateAccountStatistic(t *testing.T) {
	testConfig.PrepareDB()

	ControllerCreateAccountStatistic(t)

	accountStatistic, err := testConfig.threadClient.StatisticByAccountID(1)
	assert.NoError(t, err)
	assert.Equal(t, 0, accountStatistic.AccountStatistic.SentMessagesToTopicsCount)
	assert.Equal(t, 0, accountStatistic.AccountStatistic.CreatedTopicsCount)
	assert.Equal(t, 1, accountStatistic.AccountStatistic.AccountID)
}
