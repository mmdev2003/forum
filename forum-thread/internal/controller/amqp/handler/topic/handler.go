package topic

import (
	"encoding/json"
	"forum-thread/internal/model"
	"github.com/rabbitmq/amqp091-go"
	"golang.org/x/net/context"
)

func CreateTopicPostprocessing(
	topicService model.ITopicService,
) func(event amqp091.Delivery) error {
	return func(event amqp091.Delivery) error {
		ctx := context.Background()

		var body CreateTopicPostprocessingBody
		if err := json.Unmarshal(event.Body, &body); err != nil {
			return err
		}

		err := topicService.CreateTopicPostprocessing(ctx,
			body.TopicOwnerAccountID,
		)
		if err != nil {
			return err
		}
		return nil
	}
}

func AddViewToTopicPostprocessing(
	topicService model.ITopicService,
) func(event amqp091.Delivery) error {
	return func(event amqp091.Delivery) error {
		ctx := context.Background()

		var body AddViewToTopicPostprocessingBody
		if err := json.Unmarshal(event.Body, &body); err != nil {
			return err
		}

		err := topicService.AddViewToTopicPostprocessing(ctx,
			body.TopicID,
		)
		if err != nil {
			return err
		}
		return nil
	}
}

func CloseTopicPostprocessing(
	topicService model.ITopicService,
) func(event amqp091.Delivery) error {
	return func(event amqp091.Delivery) error {
		ctx := context.Background()

		var body CloseTopicPostprocessingBody
		if err := json.Unmarshal(event.Body, &body); err != nil {
			return err
		}

		err := topicService.CloseTopicPostprocessing(ctx,
			body.TopicOwnerAccountID,
			body.AdminAccountID,
			body.TopicID,
			body.TopicName,
			body.AdminLogin,
		)
		if err != nil {
			return err
		}
		return nil
	}
}

func ChangeTopicPriorityPostprocessing(
	topicService model.ITopicService,
) func(event amqp091.Delivery) error {
	return func(event amqp091.Delivery) error {
		ctx := context.Background()

		var body ChangeTopicPriorityPostprocessingBody
		if err := json.Unmarshal(event.Body, &body); err != nil {
			return err
		}

		err := topicService.ChangeTopicPriorityPostprocessing(ctx,
			body.SubthreadID,
			body.TopicID,
			body.TopicPriority,
		)
		if err != nil {
			return err
		}
		return nil
	}
}
