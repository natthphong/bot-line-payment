package line_message

import (
	"github.com/line/line-bot-sdk-go/v8/linebot/messaging_api"
)

type SendMessageFunc func(messages []messaging_api.MessageInterface, userId, retryToken string, notification bool) (*messaging_api.PushMessageResponse, error)

func SendMessage(bot *messaging_api.MessagingApiAPI) SendMessageFunc {
	return func(messages []messaging_api.MessageInterface, userId, retryToken string, notification bool) (*messaging_api.PushMessageResponse, error) {
		req := &messaging_api.PushMessageRequest{
			To:                   userId,
			Messages:             messages,
			NotificationDisabled: true,
		}
		return bot.PushMessage(req, retryToken)
	}

}

type SendMessageAllUserFunc func(messages []messaging_api.MessageInterface, retryToken string, notification bool) (*map[string]interface{}, error)

func SendMessageAllUser(bot *messaging_api.MessagingApiAPI) SendMessageAllUserFunc {
	return func(messages []messaging_api.MessageInterface, retryToken string, notification bool) (*map[string]interface{}, error) {
		req := &messaging_api.BroadcastRequest{
			Messages:             messages,
			NotificationDisabled: true,
		}
		return bot.Broadcast(req, retryToken)
	}
}
