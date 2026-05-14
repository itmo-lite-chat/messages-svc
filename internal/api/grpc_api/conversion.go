package grpc_api

import (
	"github.com/itmo-lite-chat/messages_svc/internal/domain"
	pb "github.com/itmo-lite-chat/proto-registry/gen/services/messages_service/messages/v1"
	"google.golang.org/protobuf/types/known/durationpb"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func externalContentTypeToInternal(tp pb.ContentType) domain.ContentType {
	switch tp {
	case pb.ContentType_CONTENT_TYPE_TEXT:
		return domain.ContentTypeText
	case pb.ContentType_CONTENT_TYPE_IMAGE:
		return domain.ContentTypeImage
	case pb.ContentType_CONTENT_TYPE_FILE:
		return domain.ContentTypeFile
	case pb.ContentType_CONTENT_TYPE_STICKER:
		return domain.ContentTypeSticker
	default:
		return domain.ContentTypeUnspecified
	}
}

func internalContentTypeToExternal(tp domain.ContentType) pb.ContentType {
	switch tp {
	case domain.ContentTypeText:
		return pb.ContentType_CONTENT_TYPE_TEXT
	case domain.ContentTypeImage:
		return pb.ContentType_CONTENT_TYPE_IMAGE
	case domain.ContentTypeFile:
		return pb.ContentType_CONTENT_TYPE_FILE
	case domain.ContentTypeSticker:
		return pb.ContentType_CONTENT_TYPE_STICKER
	default:
		return pb.ContentType_CONTENT_TYPE_UNSPECIFIED
	}
}

func internalMessageToExternal(m *domain.Message) *pb.Message {
	if m == nil {
		return nil
	}

	pbMsg := &pb.Message{
		MessageId: m.MessageID,
		ChatId:    m.ChatID.String(),
		SenderId:  m.SenderID.String(),
		Content: &pb.Content{
			Type: internalContentTypeToExternal(m.Content.Type),
			Body: m.Content.Body,
			Metadata: &pb.Metadata{
				FileSizeBytes: m.Content.Metadata.FileSizeBytes,
				Width:         m.Content.Metadata.Width,
				Height:        m.Content.Metadata.Height,
				Duration:      durationpb.New(m.Content.Metadata.Duration),
			},
		},
		CreatedAt: timestamppb.New(m.CreatedAt),
		ReplyToId: m.ReplyToID,
	}

	if m.UpdatedAt != nil {
		pbMsg.UpdatedAt = timestamppb.New(*m.UpdatedAt)
	}
	if m.DeletedAt != nil {
		pbMsg.DeletedAt = timestamppb.New(*m.DeletedAt)
	}

	return pbMsg
}
