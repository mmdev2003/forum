package model

import (
	"time"
)

type Thread struct {
	ID                int       `db:"id"`
	ThreadName        string    `db:"thread_name"`
	ThreadDescription string    `db:"thread_description"`
	ThreadColor       string    `db:"thread_color"`
	AllowedStatuses   []string  `db:"allowed_statuses"`
	CreatedAt         time.Time `db:"created_at"`
	UpdatedAt         time.Time `db:"updated_at"`
}

type Subthread struct {
	ID                        int       `db:"id"`
	ThreadID                  int       `db:"thread_id"`
	ThreadName                string    `db:"thread_name"`
	SubthreadName             string    `db:"subthread_name"`
	SubthreadDescription      string    `db:"subthread_description"`
	SubthreadMessageCount     int       `db:"subthread_message_count"`
	SubthreadViewCount        int       `db:"subthread_view_count"`
	SubthreadLastMessageLogin string    `db:"subthread_last_message_login"`
	SubthreadLastMessageText  string    `db:"subthread_last_message_text"`
	SubthreadLastMessageAt    time.Time `db:"subthread_last_message_at"`
	CreatedAt                 time.Time `db:"created_at"`
	UpdatedAt                 time.Time `db:"updated_at"`
}

type Topic struct {
	ID                    int       `db:"id"`
	ThreadID              int       `db:"thread_id"`
	ThreadName            string    `db:"thread_name"`
	SubthreadID           int       `db:"subthread_id"`
	SubthreadName         string    `db:"subthread_name"`
	TopicName             string    `db:"topic_name"`
	TopicOwnerAccountID   int       `db:"topic_owner_account_id"`
	TopicModerationStatus string    `db:"topic_moderation_status"`
	TopicOwnerLogin       string    `db:"topic_owner_login"`
	TopicLastMessageLogin string    `db:"topic_last_message_login"`
	TopicLastMessageText  string    `db:"topic_last_message_text"`
	TopicMessageCount     int       `db:"topic_message_count"`
	TopicViewCount        int       `db:"topic_view_count"`
	TopicIsClosed         bool      `db:"topic_is_closed"`
	TopicPriority         int       `db:"topic_priority"`
	TopicIsAuthor         bool      `db:"topic_is_author"`
	TopicAvatarURL        string    `db:"topic_avatar_url"`
	TopicLastMessageAt    time.Time `db:"topic_last_message_at"`
	CreatedAt             time.Time `db:"created_at"`
	UpdatedAt             time.Time `db:"updated_at"`
}

type Message struct {
	ID                    int       `db:"id"`
	TopicID               int       `db:"topic_id"`
	MessageOwnerAccountID int       `db:"message_owner_account_id"`
	MessageOwnerLogin     string    `db:"message_owner_login"`
	MessageText           string    `db:"message_text"`
	MessageFileIDs        []int     `db:"message_file_ids"`
	MessageReplyToID      int       `db:"message_reply_to_id"`
	MessageLikeCount      int       `db:"message_like_count"`
	MessageReplyCount     int       `db:"message_reply_count"`
	MessageReportCount    int       `db:"message_report_count"`
	CreatedAt             time.Time `db:"created_at"`
	UpdatedAt             time.Time `db:"updated_at"`
}

type MessageSearch struct {
	ID        int    `json:"id"`
	TopicID   int    `json:"topicID"`
	AccountID int    `json:"accountID"`
	Login     string `json:"login"`
	Text      string `json:"text"`
}

type Like struct {
	ID         int       `db:"id"`
	TopicID    int       `db:"topic_id"`
	MessageID  int       `db:"message_id"`
	AccountID  int       `db:"account_id"`
	LikeTypeID int       `db:"like_type_id"`
	CreatedAt  time.Time `db:"created_at"`
	UpdatedAt  time.Time `db:"updated_at"`
}

type File struct {
	ID        int       `db:"id"`
	MessageID int       `db:"message_id"`
	URL       string    `db:"url"`
	Name      string    `db:"name"`
	Size      int       `db:"size"`
	Extension string    `db:"extension"`
	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`
}

type AccountStatistic struct {
	ID                        int       `db:"id"`
	AccountID                 int       `db:"account_id"`
	CreatedTopicsCount        int       `db:"created_topics_count"`
	SentMessagesToTopicsCount int       `db:"sent_messages_to_topics_count"`
	CreatedAt                 time.Time `db:"created_at"`
	UpdatedAt                 time.Time `db:"updated_at"`
}
