package unit

import (
	"bytes"
	"encoding/json"
	"errors"
	topicAmqpHandler "forum-thread/internal/controller/amqp/handler/topic"
	"mime/multipart"

	"github.com/rabbitmq/amqp091-go"
	"strconv"

	"io"
	"net/http"
)

import (
	messageAmqpHandler "forum-thread/internal/controller/amqp/handler/message"
	subthreadAmqpHandler "forum-thread/internal/controller/amqp/handler/subthread"
	httpAccountStatistic "forum-thread/internal/controller/http/handler/account-statistic"
	httpMessage "forum-thread/internal/controller/http/handler/message"
	httpSubthread "forum-thread/internal/controller/http/handler/subthread"
	httpThread "forum-thread/internal/controller/http/handler/thread"
	httpTopic "forum-thread/internal/controller/http/handler/topic"
)

func NewThreadClient(host, port string) *ClientThread {
	return &ClientThread{
		client:  &http.Client{},
		baseURL: "http://" + host + ":" + port + "/api/thread",
	}
}

type ClientThread struct {
	client  *http.Client
	baseURL string
}

func (c *ClientThread) CreateThread(
	threadName,
	threadDescription,
	threadColor string,
	allowedStatuses []string,
) (*httpThread.CreateThreadResponse, error) {
	body := httpThread.CreateThreadBody{
		ThreadName:        threadName,
		ThreadDescription: threadDescription,
		ThreadColor:       threadColor,
		AllowedStatuses:   allowedStatuses,
	}

	response, err := c.post("/create", body)
	if err != nil {
		return nil, err
	}

	var createThreadResponse httpThread.CreateThreadResponse
	err = json.Unmarshal(response, &createThreadResponse)
	if err != nil {
		return nil, err
	}

	return &createThreadResponse, err
}

func (c *ClientThread) AllThread() (*httpThread.AllThreadResponse, error) {
	response, err := c.get("/all")
	if err != nil {
		return nil, err
	}

	var allThreadResponse httpThread.AllThreadResponse
	err = json.Unmarshal(response, &allThreadResponse)
	if err != nil {
		return nil, err
	}

	return &allThreadResponse, err
}

func (c *ClientThread) CreateSubthread(
	threadID int,
	threadName,
	subthreadName,
	subthreadDescription string,
) (*httpSubthread.CreateSubthreadResponse, error) {
	body := httpSubthread.CreateSubthreadBody{
		ThreadID:             threadID,
		ThreadName:           threadName,
		SubthreadName:        subthreadName,
		SubthreadDescription: subthreadDescription,
	}

	response, err := c.post("/subthread/create", body)
	if err != nil {
		return nil, err
	}

	var createSubthreadResponse httpSubthread.CreateSubthreadResponse
	err = json.Unmarshal(response, &createSubthreadResponse)
	if err != nil {
		return nil, err
	}

	return &createSubthreadResponse, err
}

func (c *ClientThread) AddViewToSubthread(
	subthreadID int,
) error {
	body := httpSubthread.AddViewToSubthreadBody{
		SubthreadID: subthreadID,
	}

	_, err := c.post("/subthread/view", body)
	if err != nil {
		return err
	}

	return nil
}

func (c *ClientThread) AddViewToSubthreadPostprocessing(
	subthreadID int,
) error {
	body, _ := json.Marshal(subthreadAmqpHandler.AddViewToSubthreadPostprocessingBody{
		SubthreadID: subthreadID,
	})
	event := amqp091.Delivery{
		Body: body,
	}
	err := subthreadAmqpHandler.AddViewToSubthreadPostprocessing(testConfig.subthreadService)(event)
	if err != nil {
		return err
	}
	return nil
}

func (c *ClientThread) SubthreadsByThreadID(
	threadID int,
) (*httpSubthread.SubthreadsByThreadIDResponse, error) {
	response, err := c.get("/subthread/" + strconv.Itoa(threadID))
	if err != nil {
		return nil, err
	}

	var subthreadsByThreadIDResponse httpSubthread.SubthreadsByThreadIDResponse
	err = json.Unmarshal(response, &subthreadsByThreadIDResponse)
	if err != nil {
		return nil, err
	}

	return &subthreadsByThreadIDResponse, err
}

func (c *ClientThread) CreateTopic(
	subthreadID,
	threadID int,
	subthreadName,
	threadName,
	topicName,
	topicOwnerLogin,
	accessToken string,
) (*httpTopic.CreateTopicResponse, error) {
	body := httpTopic.CreateTopicBody{
		SubthreadID:     subthreadID,
		ThreadID:        threadID,
		SubthreadName:   subthreadName,
		ThreadName:      threadName,
		TopicName:       topicName,
		TopicOwnerLogin: topicOwnerLogin,
	}

	response, err := c.postWithAccessToken("/topic/create", accessToken, body)
	if err != nil {
		return nil, err
	}

	var createTopicResponse httpTopic.CreateTopicResponse
	err = json.Unmarshal(response, &createTopicResponse)
	if err != nil {
		return nil, err
	}

	return &createTopicResponse, err
}

func (c *ClientThread) CreateTopicPostprocessing(
	topicOwnerAccountID int,
) error {
	body, err := json.Marshal(topicAmqpHandler.CreateTopicPostprocessingBody{
		TopicOwnerAccountID: topicOwnerAccountID,
	})
	event := amqp091.Delivery{
		Body: body,
	}
	err = topicAmqpHandler.CreateTopicPostprocessing(testConfig.topicService)(event)
	if err != nil {
		return err
	}

	return nil
}

func (c *ClientThread) AddViewToTopic(
	topicID int,
) error {
	body := httpTopic.AddViewToTopicBody{
		TopicID: topicID,
	}

	_, err := c.post("/topic/view", body)
	if err != nil {
		return err
	}

	return nil
}

func (c *ClientThread) AddViewToTopicPostprocessing(
	topicID int,
) error {
	body, err := json.Marshal(topicAmqpHandler.AddViewToTopicPostprocessingBody{
		TopicID: topicID,
	})
	event := amqp091.Delivery{
		Body: body,
	}
	err = topicAmqpHandler.AddViewToTopicPostprocessing(testConfig.topicService)(event)
	if err != nil {
		return err
	}

	return nil
}

func (c *ClientThread) CloseTopic(
	topicOwnerAccountID,
	adminAccountID,
	topicID int,
	topicName,
	adminLogin,
	accessToken string,
) error {
	body := httpTopic.CloseTopicBody{
		TopicOwnerAccountID: topicOwnerAccountID,
		AdminAccountID:      adminAccountID,
		TopicID:             topicID,
		TopicName:           topicName,
		AdminLogin:          adminLogin,
	}

	_, err := c.postWithAccessToken("/topic/close", accessToken, body)
	if err != nil {
		return err
	}

	return nil
}

func (c *ClientThread) ApproveTopic(
	topicID int,
	accessToken string,
) error {
	_, err := c.postWithAccessToken("/topic/approve/"+strconv.Itoa(topicID), accessToken, nil)
	if err != nil {
		return err
	}

	return nil
}

func (c *ClientThread) RejectTopic(
	topicID int,
	accessToken string,
) error {
	_, err := c.postWithAccessToken("/topic/reject/"+strconv.Itoa(topicID), accessToken, nil)
	if err != nil {
		return err
	}

	return nil
}

func (c *ClientThread) CloseTopicPostprocessing(
	topicID int,
) error {
	body, err := json.Marshal(topicAmqpHandler.CloseTopicPostprocessingBody{
		TopicID: topicID,
	})
	event := amqp091.Delivery{
		Body: body,
	}
	err = topicAmqpHandler.CloseTopicPostprocessing(testConfig.topicService)(event)
	if err != nil {
		return err
	}

	return nil
}

func (c *ClientThread) ChangeTopicPriority(
	topicID,
	topicPriority int,
) error {
	body := httpTopic.ChangeTopicPriorityBody{
		TopicID:       topicID,
		TopicPriority: topicPriority,
	}

	_, err := c.post("/topic/priority/change", body)
	if err != nil {
		return err
	}

	return nil
}

func (c *ClientThread) ChangeTopicPriorityPostprocessing(
	topicID,
	topicPriority int,
) error {
	body, err := json.Marshal(topicAmqpHandler.ChangeTopicPriorityPostprocessingBody{
		TopicID:       topicID,
		TopicPriority: topicPriority,
	})
	event := amqp091.Delivery{
		Body: body,
	}
	err = topicAmqpHandler.ChangeTopicPriorityPostprocessing(testConfig.topicService)(event)
	if err != nil {
		return err
	}

	return nil
}

func (c *ClientThread) TopicsBySubthreadID(
	subthreadID int,
) (*httpTopic.TopicsBySubthreadIDResponse, error) {
	response, err := c.get("/topic/subthreadID/" + strconv.Itoa(subthreadID))
	if err != nil {
		return nil, err
	}

	var topicsBySubthreadIDResponse httpTopic.TopicsBySubthreadIDResponse
	err = json.Unmarshal(response, &topicsBySubthreadIDResponse)
	if err != nil {
		return nil, err
	}

	return &topicsBySubthreadIDResponse, err
}

func (c *ClientThread) TopicsByAccountID(
	accountID int,
) (*httpTopic.TopicsByAccountIDResponse, error) {
	response, err := c.get("/topic/accountID/" + strconv.Itoa(accountID))
	if err != nil {
		return nil, err
	}

	var topicsByAccountIDResponse httpTopic.TopicsByAccountIDResponse
	err = json.Unmarshal(response, &topicsByAccountIDResponse)
	if err != nil {
		return nil, err
	}

	return &topicsByAccountIDResponse, err
}

func (c *ClientThread) TopicsOnModeration(
	accessToken string,
) (*httpTopic.TopicsOnModerationIDResponse, error) {
	response, err := c.getWithAccessToken("/topic/moderation", accessToken)
	if err != nil {
		return nil, err
	}

	var topicsOnModerationIDResponse httpTopic.TopicsOnModerationIDResponse
	err = json.Unmarshal(response, &topicsOnModerationIDResponse)
	if err != nil {
		return nil, err
	}

	return &topicsOnModerationIDResponse, err
}

func (c *ClientThread) SendMessageToTopic(
	subthreadID,
	topicID,
	replyToMessageID,
	replyMessageOwnerAccountID,
	topicOwnerAccountID int,
	senderLogin,
	threadName,
	subthreadName,
	topicName,
	senderMessageText string,
	files [][]byte,
	fileNames []string,
	accessToken string,
) (*httpMessage.SendMessageToTopicResponse, error) {

	var body bytes.Buffer
	writer := multipart.NewWriter(&body)

	payload := httpMessage.SendMessageToTopicBody{
		SubthreadID:                subthreadID,
		TopicID:                    topicID,
		ReplyToMessageID:           replyToMessageID,
		ReplyMessageOwnerAccountID: replyMessageOwnerAccountID,
		TopicOwnerAccountID:        topicOwnerAccountID,
		SenderLogin:                senderLogin,
		SenderMessageText:          senderMessageText,
		ThreadName:                 threadName,
		SubthreadName:              subthreadName,
		TopicName:                  topicName,
	}

	jsonData, err := json.Marshal(payload)
	if err != nil {
		return nil, err
	}
	err = writer.WriteField("body", string(jsonData))
	if err != nil {
		return nil, err
	}

	for i := range fileNames {
		formFile, err := writer.CreateFormFile("files", fileNames[i])
		if err != nil {
			return nil, err
		}

		_, err = io.Copy(formFile, bytes.NewReader(files[i]))
		if err != nil {
			return nil, err
		}
	}

	err = writer.Close()
	if err != nil {
		return nil, err
	}

	response, err := c.multipartFormDataWithAccessToken("/message/send", writer.FormDataContentType(), accessToken, &body)
	if err != nil {
		return nil, err
	}

	var sendMessageToTopicResponse httpMessage.SendMessageToTopicResponse
	err = json.Unmarshal(response, &sendMessageToTopicResponse)
	if err != nil {
		return nil, err
	}

	return &sendMessageToTopicResponse, err
}

func (c *ClientThread) SendMessageToTopicPostprocessing(
	subthreadID,
	topicID,
	replyToMessageID,
	replyMessageOwnerAccountID,
	topicOwnerAccountID,
	senderAccountID int,
	senderLogin,
	topicName,
	senderMessageText string,
	filesURLs []string,
	filesNames []string,
	filesExtensions []string,
	filesSizes []int,
) error {
	body, err := json.Marshal(messageAmqpHandler.SendMessageToTopicPostprocessingBody{
		SubthreadID:                subthreadID,
		TopicID:                    topicID,
		ReplyToMessageID:           replyToMessageID,
		ReplyMessageOwnerAccountID: replyMessageOwnerAccountID,
		TopicOwnerAccountID:        topicOwnerAccountID,
		SenderAccountID:            senderAccountID,
		SenderLogin:                senderLogin,
		SenderMessageText:          senderMessageText,
		TopicName:                  topicName,
		FilesURLs:                  filesURLs,
		FilesNames:                 filesNames,
		FilesExtensions:            filesExtensions,
		FilesSizes:                 filesSizes,
	})
	event := amqp091.Delivery{
		Body: body,
	}

	err = messageAmqpHandler.SendMessageToTopicPostprocessing(testConfig.messageService)(event)
	if err != nil {
		return err
	}

	return nil
}

func (c *ClientThread) LikeMessage(
	topicID,
	messageOwnerAccountID,
	likerAccountID,
	likeMessageID,
	likeTypeID int,
	likerLogin,
	topicName,
	likeMessageText string,
) error {
	body := httpMessage.LikeMessageBody{
		TopicID:               topicID,
		MessageOwnerAccountID: messageOwnerAccountID,
		LikerAccountID:        likerAccountID,
		LikeMessageID:         likeMessageID,
		LikeTypeID:            likeTypeID,
		LikerLogin:            likerLogin,
		TopicName:             topicName,
		LikeMessageText:       likeMessageText,
	}

	_, err := c.post("/message/like", body)
	if err != nil {
		return err
	}

	return nil
}

func (c *ClientThread) UnlikeMessage(
	likeMessageID int,
	accessToken string,
) error {
	body := httpMessage.UnlikeMessageBody{
		LikeMessageID: likeMessageID,
	}

	_, err := c.postWithAccessToken("/message/unlike", accessToken, body)
	if err != nil {
		return err
	}

	return nil
}

func (c *ClientThread) LikeMessagePostprocessing(
	topicID,
	messageOwnerAccountID,
	likerAccountID,
	likeMessageID,
	likeTypeID int,
	likerLogin,
	topicName,
	likeMessageText string,
) error {
	body, err := json.Marshal(messageAmqpHandler.LikeMessagePostprocessingBody{
		TopicID:               topicID,
		MessageOwnerAccountID: messageOwnerAccountID,
		LikerAccountID:        likerAccountID,
		LikeMessageID:         likeMessageID,
		LikeTypeID:            likeTypeID,
		LikerLogin:            likerLogin,
		TopicName:             topicName,
		LikeMessageText:       likeMessageText,
	})
	event := amqp091.Delivery{
		Body: body,
	}

	err = messageAmqpHandler.LikeMessagePostprocessing(testConfig.messageService)(event)
	if err != nil {
		return err
	}

	return nil
}

func (c *ClientThread) ReportMessage(
	messageID,
	accountID int,
	reportText string,
) error {
	body := httpMessage.ReportMessageBody{
		MessageID:  messageID,
		AccountID:  accountID,
		ReportText: reportText,
	}

	_, err := c.post("/message/report", body)
	if err != nil {
		return err
	}

	return nil
}

func (c *ClientThread) EditMessage(
	messageID int,
	messageText string,
) error {
	body := httpMessage.EditMessageBody{
		MessageID:   messageID,
		MessageText: messageText,
	}

	_, err := c.post("/message/edit", body)
	if err != nil {
		return err
	}

	return nil
}

func (c *ClientThread) ReportMessagePostprocessing(
	messageID,
	accountID int,
	reportText string,
) error {
	body, err := json.Marshal(messageAmqpHandler.ReportMessagePostprocessingBody{
		MessageID:  messageID,
		AccountID:  accountID,
		ReportText: reportText,
	})
	event := amqp091.Delivery{
		Body: body,
	}

	err = messageAmqpHandler.ReportMessagePostprocessing(testConfig.messageService)(event)
	if err != nil {
		return err
	}

	return nil
}

func (c *ClientThread) MessagesByTopicID(
	topicID int,
	accessToken string,
) (*httpMessage.MessagesByTopicIDResponse, error) {
	response, err := c.getWithAccessToken("/message/topicID/"+strconv.Itoa(topicID), accessToken)
	if err != nil {
		return nil, err
	}

	var messagesByTopicIDResponse httpMessage.MessagesByTopicIDResponse
	err = json.Unmarshal(response, &messagesByTopicIDResponse)
	if err != nil {
		return nil, err
	}

	return &messagesByTopicIDResponse, err
}

func (c *ClientThread) MessagesByAccountID(
	accountID int,
) (*httpMessage.MessagesByAccountIDResponse, error) {
	response, err := c.get("/message/accountID/" + strconv.Itoa(accountID))
	if err != nil {
		return nil, err
	}

	var messagesByAccountIDResponse httpMessage.MessagesByAccountIDResponse
	err = json.Unmarshal(response, &messagesByAccountIDResponse)
	if err != nil {
		return nil, err
	}

	return &messagesByAccountIDResponse, err
}

func (c *ClientThread) MessagesByText(
	text string,
) (*httpMessage.MessagesByTextResponse, error) {
	response, err := c.get("/message/search/" + text)
	if err != nil {
		return nil, err
	}

	var messagesByTextResponse httpMessage.MessagesByTextResponse
	err = json.Unmarshal(response, &messagesByTextResponse)
	if err != nil {
		return nil, err
	}

	return &messagesByTextResponse, err
}

func (c *ClientThread) DownloadFile(
	fileURL string,
) ([]byte, error) {
	response, err := c.get("/file/download/" + fileURL)
	if err != nil {
		return nil, err
	}

	return response, err
}

func (c *ClientThread) CreateAccountStatistic(
	accountID int,
) error {
	body := httpAccountStatistic.CreateAccountStatisticBody{
		AccountID: accountID,
	}

	_, err := c.post("/statistic/create", body)
	if err != nil {
		return err
	}

	return nil
}

func (c *ClientThread) StatisticByAccountID(
	accountID int,
) (*httpAccountStatistic.StatisticByAccountIDResponse, error) {
	accountIDStr := strconv.Itoa(accountID)
	response, err := c.get("/statistic/" + accountIDStr)
	if err != nil {
		return nil, err
	}

	var statisticByAccountIDResponse httpAccountStatistic.StatisticByAccountIDResponse
	err = json.Unmarshal(response, &statisticByAccountIDResponse)
	if err != nil {
		return nil, err
	}

	return &statisticByAccountIDResponse, err
}

func (c *ClientThread) post(path string, body any) ([]byte, error) {
	bytesBody, _ := json.Marshal(body)

	req, err := http.NewRequest("POST", c.baseURL+path, bytes.NewBuffer(bytesBody))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 300 {
		return nil, errors.New("status code is not 200")
	}

	response, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return response, nil
}

func (c *ClientThread) postWithAccessToken(path, accessToken string, body any) ([]byte, error) {
	bytesBody, _ := json.Marshal(body)

	req, err := http.NewRequest("POST", c.baseURL+path, bytes.NewBuffer(bytesBody))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	req.AddCookie(&http.Cookie{Name: "Access-Token", Value: accessToken})

	resp, err := c.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 300 {
		return nil, errors.New("status code is not 200")
	}

	response, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return response, nil
}

func (c *ClientThread) get(path string) ([]byte, error) {
	req, err := http.NewRequest("GET", c.baseURL+path, nil)
	if err != nil {
		return nil, err
	}

	resp, err := c.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 300 {
		return nil, errors.New("status code is not 200")
	}

	response, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return response, nil
}

func (c *ClientThread) getWithAccessToken(path, accessToken string) ([]byte, error) {
	req, err := http.NewRequest("GET", c.baseURL+path, nil)
	if err != nil {
		return nil, err
	}
	req.AddCookie(&http.Cookie{Name: "Access-Token", Value: accessToken})

	resp, err := c.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 300 {
		return nil, errors.New("status code is not 200")
	}

	response, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return response, nil
}

func (c *ClientThread) multipartFormDataWithAccessToken(path, contentType, accessToken string, body *bytes.Buffer) ([]byte, error) {

	req, err := http.NewRequest("POST", c.baseURL+path, body)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", contentType)
	req.AddCookie(&http.Cookie{Name: "Access-Token", Value: accessToken})

	resp, err := c.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode >= 300 {
		return nil, errors.New("status code is not 200")
	}

	response, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return response, nil
}
