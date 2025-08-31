package subthread

import (
	"context"
	"encoding/json"
	"forum-thread/internal/model"
	"github.com/rabbitmq/amqp091-go"
)

func AddViewToSubthreadPostprocessing(
	subthreadService model.ISubthreadService,
) func(event amqp091.Delivery) error {
	return func(event amqp091.Delivery) error {
		ctx := context.Background()

		var body AddViewToSubthreadPostprocessingBody
		if err := json.Unmarshal(event.Body, &body); err != nil {
			return err
		}

		err := subthreadService.AddViewToSubthreadPostprocessing(ctx,
			body.SubthreadID,
		)
		if err != nil {
			return err
		}
		return nil
	}
}
