package rcrequestmodels

type RCReportRequest struct {
	ChannelID uint   `json:"channel_id"`
	Content   string `json:"content"`
	Type      string `json:"type"`
}
