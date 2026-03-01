package domain

type ContentType = string

const (
	ContentTypeUnspecified ContentType = ""
	ContentTypeText        ContentType = "text"
	ContentTypeImage       ContentType = "image"
	ContentTypeFile        ContentType = "file"
	ContentTypeSticker     ContentType = "sticker"
)
