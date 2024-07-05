package web

import (
	"encoding/json"
	"github.com/k3env/notification-webhook/model"
	"github.com/k3env/notification-webhook/tg"
	"github.com/rs/zerolog/log"
	"io"
	"net/http"
)

type WebhookHandler struct {
	sender *tg.TelegramSender
}

func NewWebhookHandler(sender *tg.TelegramSender) *WebhookHandler {
	return &WebhookHandler{sender: sender}
}

func (h *WebhookHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	body, _ := io.ReadAll(r.Body)
	webhook := model.WebhookRequest{}
	err := json.Unmarshal(body, &webhook)
	if err != nil {
		log.Error().Err(err).Msg("Failed to unmarshal JSON")
	}
	for _, alert := range webhook.Alerts {
		err := h.sender.Send(&alert)
		if err != nil {
			log.Error().Err(err).Msg("Failed to send alert")
		}
	}
	w.Write([]byte("ok"))
}
