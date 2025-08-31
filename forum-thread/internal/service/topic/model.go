package subthread

type CreateTopicPostprocessingBody struct {
	TopicOwnerAccountID int `json:"topicOwnerAccountID"`
}

type AddViewToTopicPostprocessingBody struct {
	TopicID int `json:"topicID"`
}

type CloseTopicPostprocessingBody struct {
	TopicOwnerAccountID int    `json:"topicOwnerAccountID"`
	AdminAccountID      int    `json:"adminAccountID"`
	TopicID             int    `json:"topicID"`
	TopicName           string `json:"topicName"`
	AdminLogin          string `json:"adminLogin"`
}

type ChangeTopicPriorityPostprocessingBody struct {
	SubthreadID   int `json:"subthreadID"`
	TopicID       int `json:"topicID"`
	TopicPriority int `json:"topicPriority"`
}
