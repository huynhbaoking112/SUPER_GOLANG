package dto

import (
	"time"
)

type GenericEvent struct {
	EventID       string      `json:"eventId"`
	Topic         string      `json:"topic"`
	SourceService string      `json:"sourceService"`
	Timestamp     time.Time   `json:"timestamp"`
	Payload       interface{} `json:"payload"`
}

type WorkspaceCreatedPayload struct {
	WorkspaceID   string `json:"workspaceId"`
	WorkspaceName string `json:"workspaceName"`
	WorkspaceSlug string `json:"workspaceSlug"`
	CreatedByID   string `json:"createdById"`
}

type UserLoginPayload struct {
	UserID string `json:"userId"`
}
