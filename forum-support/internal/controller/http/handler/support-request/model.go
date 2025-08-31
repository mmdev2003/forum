package support_request

type CreateSupportRequestRequest struct {
	Title       string `json:"title"`
	Description string `json:"description"`
}

type CreateSupportRequestResponse struct {
	SupportRequestID int `json:"supportRequestID"`
}
