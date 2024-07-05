package model

type TelegramSendRequest struct {
	Text      string `json:"text"`
	ChatID    string `json:"chat_id"`
	ParseMode string `json:"parse_mode"`
}

type TelegramReplyParameters struct {
	MessageID int `json:"message_id"`
}

type TelegramReplyRequest struct {
	Text            string                  `json:"text"`
	ChatID          string                  `json:"chat_id"`
	ParseMode       string                  `json:"parse_mode"`
	ReplyParameters TelegramReplyParameters `json:"reply_parameters"`
}
