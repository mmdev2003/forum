package unit

import (
	"github.com/stretchr/testify/assert"
	"strconv"
	"testing"
)

func ControllerCreateThreads(t *testing.T) []TestThread {
	expectedThreads := []TestThread{
		{1, "thread1", "desc1", []string{"Writer", "Admin"}},
		{2, "thread2", "desc2", []string{"Writer", "Admin"}},
	}

	for _, thread := range expectedThreads {
		_, err := testConfig.threadClient.CreateThread(
			thread.ThreadName,
			thread.ThreadDescription,
			"green",
			thread.AllowedStatuses,
		)
		assert.NoError(t, err)
	}
	return expectedThreads
}

func ControllerCreateThread(t *testing.T) TestThread {
	testThread := TestThread{1, "thread1", "desc1", []string{"Writer", "Admin"}}
	_, err := testConfig.threadClient.CreateThread(
		testThread.ThreadName,
		testThread.ThreadDescription,
		"green",
		testThread.AllowedStatuses,
	)
	assert.NoError(t, err)
	return testThread
}

func ControllerCreateSubthreads(t *testing.T) []TestSubthread {
	expectedSubthreads := []TestSubthread{
		{1, 1, "thread1", "subthread1", "desc1"},
		{2, 1, "thread1", "subthread2", "desc2"},
		{3, 2, "thread2", "subthread3", "desc3"},
		{4, 2, "thread2", "subthread4", "desc4"},
	}
	for _, subthread := range expectedSubthreads {
		_, err := testConfig.threadClient.CreateSubthread(
			subthread.ThreadID,
			subthread.ThreadName,
			subthread.SubthreadName,
			subthread.SubthreadDescription,
		)
		assert.NoError(t, err)
	}
	return expectedSubthreads
}

func ControllerCreateSubthread(t *testing.T) TestSubthread {
	testSubthread := TestSubthread{1, 1, "thread1", "subthread1", "desc1"}
	_, err := testConfig.threadClient.CreateSubthread(
		testSubthread.ThreadID,
		testSubthread.ThreadName,
		testSubthread.SubthreadName,
		testSubthread.SubthreadDescription,
	)
	assert.NoError(t, err)
	return testSubthread
}

func ControllerCreateTopics(t *testing.T) []TestTopic {
	expectedTopics := []TestTopic{
		{1, 1, 1, 1, "subthread1", "thread1", "topic1", "user1"},
		{2, 1, 1, 2, "subthread1", "thread1", "topic2", "user2"},
		{3, 3, 2, 1, "subthread3", "thread2", "topic3", "user1"},
		{4, 3, 2, 2, "subthread3", "thread2", "topic4", "user2"},
	}
	for _, topic := range expectedTopics {
		_, err := testConfig.threadClient.CreateTopic(
			topic.SubthreadID,
			topic.ThreadID,
			topic.SubthreadName,
			topic.ThreadName,
			topic.TopicName,
			topic.TopicOwnerLogin,
			"user"+strconv.Itoa(topic.TopicOwnerAccountID),
		)
		assert.NoError(t, err)
	}
	return expectedTopics
}

func ControllerCreateTopicsPostprocessing(t *testing.T) []TestTopic {
	expectedTopics := []TestTopic{
		{1, 1, 1, 1, "subthread1", "thread1", "topic1", "user1"},
		{2, 1, 1, 2, "subthread1", "thread1", "topic2", "user2"},
		{3, 3, 2, 1, "subthread3", "thread2", "topic3", "user1"},
		{4, 3, 2, 2, "subthread3", "thread2", "topic4", "user2"},
	}

	for _, topic := range expectedTopics {
		err := testConfig.threadClient.CreateTopicPostprocessing(topic.TopicOwnerAccountID)
		assert.NoError(t, err)
	}

	return expectedTopics
}

func ControllerCreateTopic(t *testing.T) TestTopic {
	testTopic := TestTopic{1, 1, 1, 1, "subthread1", "thread1", "topic1", "user1"}

	_, err := testConfig.threadClient.CreateTopic(
		testTopic.SubthreadID,
		testTopic.ThreadID,
		testTopic.SubthreadName,
		testTopic.ThreadName,
		testTopic.TopicName,
		testTopic.TopicOwnerLogin,
		"user"+strconv.Itoa(testTopic.TopicOwnerAccountID),
	)
	assert.NoError(t, err)

	return testTopic
}

func ControllerCreateTopicPostprocessing(t *testing.T) TestTopic {
	testTopic := TestTopic{1, 1, 1, 1, "subthread1", "thread1", "topic1", "user1"}

	err := testConfig.threadClient.CreateTopicPostprocessing(testTopic.TopicOwnerAccountID)
	assert.NoError(t, err)

	return testTopic
}

func ControllerSendMessagesToTopic(t *testing.T) []TestMessage {
	expectedMessages := []TestMessage{
		{1, 1, 1, 1, "user1", "Всем привет"},
		{2, 1, 1, 2, "user2", "Привет я владимир"},
		{3, 2, 1, 1, "user1", "алоха товариши приветики"},
		{4, 2, 1, 2, "user2", "может быть скажем приветствуем всем вам"},
	}

	for _, message := range expectedMessages {
		_, err := testConfig.threadClient.SendMessageToTopic(
			message.SubthreadID,
			message.TopicID,
			0,
			1,
			1,
			message.MessageOwnerLogin,
			"thread1",
			"subthread1",
			"topic1",
			message.MessageText,
			nil,
			nil,
			"user"+strconv.Itoa(message.MessageOwnerAccountID),
		)
		assert.NoError(t, err)
	}

	return expectedMessages
}

func ControllerSendMessagesToTopicPostprocessing(t *testing.T) []TestMessage {
	expectedMessages := []TestMessage{
		{1, 1, 1, 1, "user1", "Всем привет"},
		{2, 1, 1, 2, "user2", "Привет я владимир"},
		{3, 2, 1, 1, "user1", "алоха товариши приветики"},
		{4, 2, 1, 2, "user2", "может быть скажем приветствуем всем вам"},
	}

	for _, message := range expectedMessages {
		err := testConfig.threadClient.SendMessageToTopicPostprocessing(
			message.SubthreadID,
			message.TopicID,
			0,
			1,
			1,
			message.MessageOwnerAccountID,
			message.MessageOwnerLogin,
			"topic",
			message.MessageText,
			nil,
			nil,
			nil,
			nil,
		)
		assert.NoError(t, err)
	}

	return expectedMessages
}

func ControllerSendMessageToTopic(t *testing.T) TestMessage {
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
		"user"+strconv.Itoa(testMessage.MessageOwnerAccountID),
	)
	assert.NoError(t, err)

	return testMessage
}

func ControllerSendMessageToTopicPostprocessing(t *testing.T) TestMessage {
	testMessage := TestMessage{1, 1, 1, 1, "user1", "Всем привет"}
	err := testConfig.threadClient.SendMessageToTopicPostprocessing(
		testMessage.SubthreadID,
		testMessage.TopicID,
		0,
		0,
		1,
		testMessage.MessageOwnerAccountID,
		testMessage.MessageOwnerLogin,
		"topic",
		testMessage.MessageText,
		nil,
		nil,
		nil,
		nil,
	)
	assert.NoError(t, err)

	return testMessage
}

func ControllerCreateAccountStatistic(t *testing.T) {
	err := testConfig.threadClient.CreateAccountStatistic(1)
	assert.NoError(t, err)
}
