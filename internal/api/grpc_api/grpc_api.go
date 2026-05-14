package grpc_api

import (
	"context"

	"github.com/google/uuid"
	"github.com/itmo-lite-chat/messages_svc/internal/domain"
	pb "github.com/itmo-lite-chat/proto-registry/gen/services/messages_service/messages/v1"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type API struct {
	pb.UnimplementedMessagesServiceServer
	service Service
}

func NewAPI(service Service) *API {
	return &API{service: service}
}

func (a *API) SendMessage(ctx context.Context, req *pb.SendMessageRequest) (*pb.SendMessageResponse, error) {
	chatID, err := uuid.Parse(req.GetChatId())
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid chatID")
	}

	senderID, err := uuid.Parse(req.GetSenderId())
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid sender_id")
	}

	msg := &domain.Message{
		ChatID:   chatID,
		SenderID: senderID,
		Content: domain.Content{
			Type: externalContentTypeToInternal(req.GetContent().GetType()),
			Body: req.GetContent().GetBody(),
			Metadata: domain.Metadata{
				FileSizeBytes: req.GetContent().GetMetadata().GetFileSizeBytes(),
				Width:         req.GetContent().GetMetadata().GetWidth(),
				Height:        req.GetContent().GetMetadata().GetHeight(),
				Duration:      req.GetContent().GetMetadata().GetDuration().AsDuration(),
			},
		},
		ReplyToID: req.ReplyToId,
	}

	savedMsg, err := a.service.SendMessage(ctx, msg)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to process message: %v", err)
	}

	return &pb.SendMessageResponse{
		Message: internalMessageToExternal(savedMsg),
	}, nil
}

func (a *API) GetChatHistory(ctx context.Context, req *pb.GetChatHistoryRequest) (*pb.GetChatHistoryResponse, error) {
	chatID, err := uuid.Parse(req.GetChatId())
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid uuid")
	}
	domainMsgs, err := a.service.GetChatHistory(ctx, chatID, req.GetLimit(), req.BeforeMessageId)
	if err != nil {
		return nil, err
	}

	var pbMsgs []*pb.Message
	for _, m := range domainMsgs {
		pbMsgs = append(pbMsgs, internalMessageToExternal(m))
	}

	return &pb.GetChatHistoryResponse{
		Messages: pbMsgs,
	}, nil
}

func (a *API) GetLastMessages(ctx context.Context, req *pb.GetLastMessagesRequest) (*pb.GetLastMessagesResponse, error) {
	chatIDs := make([]uuid.UUID, 0, len(req.GetChatIds()))
	for _, rawChatID := range req.GetChatIds() {
		chatID, err := uuid.Parse(rawChatID)
		if err != nil {
			return nil, status.Error(codes.InvalidArgument, "invalid chat_id")
		}
		chatIDs = append(chatIDs, chatID)
	}

	lastMessages, err := a.service.GetLastMessages(ctx, chatIDs)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to get last messages: %v", err)
	}

	resp := &pb.GetLastMessagesResponse{
		LastMessages: make([]*pb.LastMessage, 0, len(lastMessages)),
	}
	for chatID, msg := range lastMessages {
		resp.LastMessages = append(resp.LastMessages, &pb.LastMessage{
			ChatId:  chatID.String(),
			Message: internalMessageToExternal(msg),
		})
	}

	return resp, nil
}

func (a *API) EditMessage(ctx context.Context, req *pb.EditMessageRequest) (*pb.EditMessageResponse, error) {
	senderID, err := uuid.Parse(req.GetSenderId())
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid sender_id format")
	}

	content := domain.Content{
		Type: externalContentTypeToInternal(req.GetContent().GetType()),
		Body: req.GetContent().GetBody(),
		Metadata: domain.Metadata{
			FileSizeBytes: req.GetContent().GetMetadata().GetFileSizeBytes(),
			Duration:      req.GetContent().GetMetadata().GetDuration().AsDuration(),
			Width:         req.GetContent().GetMetadata().GetWidth(),
			Height:        req.GetContent().GetMetadata().GetHeight(),
		},
	}

	updatedMsg, err := a.service.EditMessage(ctx, req.GetMessageId(), senderID, content)

	if err != nil {
		return nil, status.Errorf(codes.Internal, "edit failed: %v", err)
	}

	return &pb.EditMessageResponse{
		Message: internalMessageToExternal(updatedMsg),
	}, nil
}

func (a *API) DeleteMessage(ctx context.Context, req *pb.DeleteMessageRequest) (*pb.DeleteMessageResponse, error) {
	userID, err := uuid.Parse(req.GetUserId())
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid user_id format")
	}

	err = a.service.DeleteMessage(ctx, req.GetMessageId(), userID.String())
	if err != nil {
		return nil, status.Errorf(codes.Internal, "delete failed: %v", err)
	}
	return &pb.DeleteMessageResponse{}, nil
}

func (a *API) GetUnreadCount(ctx context.Context, req *pb.GetUnreadCountRequest) (*pb.GetUnreadCountResponse, error) {
	chatID, err := uuid.Parse(req.GetChatId())
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid chat_id")
	}

	count, err := a.service.GetUnreadCount(ctx, chatID, req.GetLastReadMessageId())
	if err != nil {
		return nil, status.Errorf(codes.Internal, "count failed: %v", err)
	}
	return &pb.GetUnreadCountResponse{UnreadCount: count}, nil
}

func (a *API) GetUnreadCounts(ctx context.Context, req *pb.GetUnreadCountsRequest) (*pb.GetUnreadCountsResponse, error) {
	resp := &pb.GetUnreadCountsResponse{
		UnreadCounts: make([]*pb.UnreadCount, 0, len(req.GetItems())),
	}

	for _, item := range req.GetItems() {
		chatID, err := uuid.Parse(item.GetChatId())
		if err != nil {
			return nil, status.Error(codes.InvalidArgument, "invalid chat_id")
		}

		count, err := a.service.GetUnreadCount(ctx, chatID, item.GetLastReadMessageId())
		if err != nil {
			return nil, status.Errorf(codes.Internal, "count failed: %v", err)
		}

		resp.UnreadCounts = append(resp.UnreadCounts, &pb.UnreadCount{
			ChatId:      item.GetChatId(),
			UnreadCount: count,
		})
	}

	return resp, nil
}
