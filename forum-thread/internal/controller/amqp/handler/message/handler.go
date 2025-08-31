package message

import (
	"context"
	"encoding/json"
	"forum-thread/internal/model"
	"github.com/rabbitmq/amqp091-go"
)

func SendMessageToTopicPostprocessing(
	messageService model.IMessageService,
) func(event amqp091.Delivery) error {
	return func(event amqp091.Delivery) error {
		ctx := context.Background()

		var body SendMessageToTopicPostprocessingBody
		if err := json.Unmarshal(event.Body, &body); err != nil {
			return err
		}

		_, err := messageService.SendMessageToTopicPostprocessing(ctx,
			body.SubthreadID,
			body.TopicID,
			body.ReplyToMessageID,
			body.ReplyMessageOwnerAccountID,
			body.TopicOwnerAccountID,
			body.SenderAccountID,
			body.SenderLogin,
			body.TopicName,
			body.SenderMessageText,
			body.FilesURLs,
			body.FilesNames,
			body.FilesExtensions,
			body.FilesSizes,
		)
		if err != nil {
			return err
		}
		return nil
	}
}

func LikeMessagePostprocessing(
	messageService model.IMessageService,
) func(event amqp091.Delivery) error {
	return func(event amqp091.Delivery) error {
		ctx := context.Background()

		var body LikeMessagePostprocessingBody
		if err := json.Unmarshal(event.Body, &body); err != nil {
			return err
		}

		err := messageService.LikeMessagePostprocessing(ctx,
			body.TopicID,
			body.MessageOwnerAccountID,
			body.LikerAccountID,
			body.LikeMessageID,
			body.LikeTypeID,
			body.LikerLogin,
			body.TopicName,
			body.LikeMessageText,
		)
		if err != nil {
			return err
		}
		return nil
	}
}

func ReportMessagePostprocessing(
	messageService model.IMessageService,
) func(event amqp091.Delivery) error {
	return func(event amqp091.Delivery) error {
		ctx := context.Background()

		var body ReportMessagePostprocessingBody
		if err := json.Unmarshal(event.Body, &body); err != nil {
			return err
		}

		err := messageService.ReportMessagePostprocessing(ctx,
			body.AccountID,
			body.MessageID,
			body.ReportText,
		)
		if err != nil {
			return err
		}
		return nil
	}
}
