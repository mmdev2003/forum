package subthread

import "forum-thread/internal/model"

type CreateSubthreadBody struct {
	ThreadID             int    `json:"threadID" validate:"required"`
	ThreadName           string `json:"threadName" validate:"required"`
	SubthreadName        string `json:"subthreadName" validate:"required"`
	SubthreadDescription string `json:"subthreadDescription" validate:"required"`
}

type CreateSubthreadResponse struct {
	SubthreadID int `json:"subthreadID"`
}

type AddViewToSubthreadBody struct {
	SubthreadID int `json:"subthreadID" validate:"required"`
}

type SubthreadsByThreadIDResponse struct {
	Subthreads []*model.Subthread
}
