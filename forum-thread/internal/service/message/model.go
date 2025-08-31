package message

type SendMessageToTopicPostprocessingBody struct {
	SubthreadID                int      `json:"subthreadID"`
	TopicID                    int      `json:"topicID"`
	ReplyToMessageID           int      `json:"replyToMessageID"`
	ReplyMessageOwnerAccountID int      `json:"replyMessageOwnerAccountID"`
	SenderAccountID            int      `json:"senderAccountID"`
	TopicOwnerAccountID        int      `json:"topicOwnerAccountID"`
	SenderLogin                string   `json:"senderLogin"`
	SenderMessageText          string   `json:"senderMessageText"`
	TopicName                  string   `json:"topicName"`
	FilesURLs                  []string `json:"filesURLs"`
	FilesNames                 []string `json:"filesNames"`
	FilesExtensions            []string `json:"filesExtensions"`
	FilesSizes                 []int    `json:"filesSizes"`
}

type LikeMessagePostprocessingBody struct {
	TopicID               int    `json:"topicID"`
	MessageOwnerAccountID int    `json:"messageOwnerAccountID"`
	LikerAccountID        int    `json:"likerAccountID"`
	LikeMessageID         int    `json:"likeMessageID"`
	LikeTypeID            int    `json:"likeTypeID"`
	LikerLogin            string `json:"likerLogin"`
	TopicName             string `json:"topicName"`
	LikeMessageText       string `json:"likeMessageText"`
}

type ReportMessagePostprocessingBody struct {
	AccountID  int    `json:"accountID"`
	MessageID  int    `json:"messageID"`
	ReportText string `json:"reportText"`
}
