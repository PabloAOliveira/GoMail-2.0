package models

type EmailTask struct {
	ID        int    `json:"id,omitempty"`
	To        string `json:"to" binding:"required"`
	Subject   string `json:"subject" binding:"required"`
	Body      string `json:"body" binding:"required"`
	Status    string `json:"status,omitempty"`
	CreatedAt string `json:"created_at,omitempty"`
	UpdatedAt string `json:"updated_at,omitempty"`
}

const (
	StatusEnqueued = "enqueued"
	StatusSending  = "sending"
	StatusSent     = "sent"
	StatusFailed   = "failed"
)