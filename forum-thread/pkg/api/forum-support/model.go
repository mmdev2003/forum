package forum_support

type ReportMessageBody struct {
	AccountID  int    `json:"accountID"`
	MessageID  int    `json:"messageID"`
	ReportText string `json:"reportText"`
}
