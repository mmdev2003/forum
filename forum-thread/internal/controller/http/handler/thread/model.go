package thread

import "forum-thread/internal/model"

type CreateThreadBody struct {
	ThreadName        string   `json:"threadName" validate:"required"`
	ThreadDescription string   `json:"threadDescription" validate:"required"`
	ThreadColor       string   `json:"threadColor" validate:"required"`
	AllowedStatuses   []string `json:"allowedStatuses" validate:"required"`
}

type CreateThreadResponse struct {
	ThreadID int `json:"threadID"`
}

type AllThreadResponse struct {
	Threads []*model.Thread `json:"threads"`
}
