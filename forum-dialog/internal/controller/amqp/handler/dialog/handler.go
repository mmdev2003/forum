package dialog

import (
	"context"
	"encoding/json"
	"forum-dialog/internal/model"
	"github.com/rabbitmq/amqp091-go"
)

func SendMessageToDialog(
	dialogService model.IDialogService,
) func(event amqp091.Delivery) error {
	return func(event amqp091.Delivery) error {
		ctx := context.Background()

		var body model.DialogWsMessage
		if err := json.Unmarshal(event.Body, &body); err != nil {
			return err
		}

		_, err := dialogService.AddMessageToDialog(ctx,
			body.DialogID,
			body.FromAccountID,
			body.ToAccountID,
			body.Text,
			body.FilesURLs,
		)
		if err != nil {
			return err
		}
		return nil
	}
}
