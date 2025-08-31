package unit

type TestThread struct {
	ID                int
	ThreadName        string
	ThreadDescription string
	AllowedStatuses   []string
}

type TestSubthread struct {
	ID                   int
	ThreadID             int
	ThreadName           string
	SubthreadName        string
	SubthreadDescription string
}
type TestTopic struct {
	ID                  int
	SubthreadID         int
	ThreadID            int
	TopicOwnerAccountID int
	SubthreadName       string
	ThreadName          string
	TopicName           string
	TopicOwnerLogin     string
}
type TestMessage struct {
	ID                    int
	TopicID               int
	SubthreadID           int
	MessageOwnerAccountID int
	MessageOwnerLogin     string
	MessageText           string
}
