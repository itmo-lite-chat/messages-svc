package domain

import "time"

type Metadata struct {
	FileSizeBytes int32         `json:"file_size_bytes,omitempty" bson:"file_size_bytes,omitempty"`
	Duration      time.Duration `json:"duration,omitempty" bson:"duration,omitempty"`
	Width         int32         `json:"width,omitempty" bson:"width,omitempty"`
	Height        int32         `json:"height,omitempty" bson:"height,omitempty"`
}

type Content struct {
	Type     ContentType `json:"type" bson:"type"`
	Body     string      `json:"body" bson:"body"`
	Metadata Metadata    `json:"metadata" bson:"metadata"`
}

type Message struct {
	MessageID MessageID  `json:"id" bson:"_id,omitempty"`
	ChatID    ChatID     `json:"chat_id" bson:"chat_id"`
	SenderID  UserID     `json:"sender_id" bson:"sender_id"`
	Content   Content    `json:"content" bson:"content"`
	Timestamp time.Time  `json:"timestamp" bson:"timestamp"`
	CreatedAt time.Time  `json:"created_at" bson:"created_at"`
	UpdatedAt *time.Time `json:"updated_at" bson:"updated_at"`
	DeletedAt *time.Time `json:"deleted_at" bson:"deleted_at"`
	ReplyToID *MessageID `json:"reply_to" bson:"reply_to_id"`
}
