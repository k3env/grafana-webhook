package model

type TelegramResponse struct {
	Ok          bool                   `json:"ok"`
	Result      TelegramResponseResult `json:"result"`
	ErrorCode   int                    `json:"error_code"`
	Description string                 `json:"description"`
}

type TelegramResponseResult struct {
	MessageID int                      `json:"message_id"`
	From      TelegramResponseFrom     `json:"from"`
	Chat      TelegramResponseChat     `json:"chat"`
	Date      int                      `json:"date"`
	Text      string                   `json:"text"`
	Entities  []TelegramResponseEntity `json:"entities"`
}

type TelegramResponseFrom struct {
	ID        int64  `json:"id"`
	IsBot     bool   `json:"is_bot"`
	FirstName string `json:"first_name"`
	Username  string `json:"username"`
}

type TelegramResponseChat struct {
	ID        int    `json:"id"`
	FirstName string `json:"first_name"`
	Username  string `json:"username"`
	Type      string `json:"type"`
}

type TelegramResponseEntity struct {
	Offset int    `json:"offset"`
	Length int    `json:"length"`
	Type   string `json:"type"`
}
