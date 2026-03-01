package grpc_api

import (
	"github.com/itmo-lite-chat/messages_svc/internal/domain"
	pb "github.com/itmo-lite-chat/proto-registry/gen/services/messages_service/messages/v1"
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
