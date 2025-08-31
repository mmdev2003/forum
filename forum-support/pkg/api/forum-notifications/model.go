package forum_notifications

type (
	ResponseToSupportRequestNotificationRequest struct {
		RequesterAccountID int `json:"requesterAccountID"`
		SupportRequestID   int `json:"supportRequestID"`
	}

	SupportRequestClosedNotificationRequest struct {
		RequesterAccountID int `json:"requesterAccountID"`
		SupportRequestID   int `json:"supportRequestID"`
	}

	StatusReceivedNotificationRequest struct {
		ReceiverAccountID int    `json:"receiverAccountID"`
		StatusName        string `json:"statusMame"`
	}
)
