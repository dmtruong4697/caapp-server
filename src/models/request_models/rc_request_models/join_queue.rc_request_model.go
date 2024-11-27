package rcrequestmodels

type RCQueueRequest struct {
	UserID       uint   `json:"user_id"`
	Gender       string `json:"gender"`
	TargetGender string `json:"target_gender"`
	RequestType  string `json:"request_type"` // 0: join, 1: quit
}
