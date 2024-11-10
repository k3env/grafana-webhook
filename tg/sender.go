package tg

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/k3env/grafana-webhook/config"
	"github.com/k3env/grafana-webhook/model"
	"github.com/rs/zerolog/log"
	"io"
	"net/http"
	"strings"
	"text/template"
)

type TelegramSender struct {
	config    config.TelegramConfig
	fingers   map[string]int
	templates map[string]*template.Template
}

func NewTelegramSender(config config.TelegramConfig, templates map[string]*template.Template) *TelegramSender {
	return &TelegramSender{config: config, fingers: make(map[string]int), templates: templates}
}

func (s *TelegramSender) Send(alert *model.Alert) error {
	var t *template.Template
	var ok bool
	var tName string
	switch alert.Status {
	case "resolved":
		tName = alert.Labels["tpl_resolved"]
		t, ok = s.templates[alert.Labels["tpl_resolved"]]
		break
	case "firing":
		tName = alert.Labels["tpl_firing"]
		t, ok = s.templates[alert.Labels["tpl_firing"]]
		break
	default:
		t, ok = s.templates["templates.generic"]
	}

	if !ok {
		log.Warn().Str("template", tName).Msg("No template found, using default")
		t = s.templates["templates.generic"]
		//continue
	}
	swbuf := bytes.NewBuffer(nil)
	err := t.Execute(swbuf, alert)
	if err != nil {
		return err
	}

	text := swbuf.String()
	if s.config.ParseMode == "MarkdownV2" {
		specials := []string{".", "+", "-", "|", "!"}
		for _, special := range specials {
			text = strings.ReplaceAll(text, special, "\\"+special)
		}
	}

	c := http.Client{}

	var bb []byte

	chat, ok := s.fingers[alert.Fingerprint]
	if alert.Status == "resolved" && ok {
		body := model.TelegramReplyRequest{
			Text:      text,
			ChatID:    s.config.Receiver,
			ParseMode: s.config.ParseMode,
			ReplyParameters: struct {
				MessageID int `json:"message_id"`
			}(struct{ MessageID int }{
				MessageID: chat,
			}),
		}
		bb, _ = json.Marshal(body)
		delete(s.fingers, alert.Fingerprint)
	} else {
		body := model.TelegramSendRequest{
			Text:      text,
			ChatID:    s.config.Receiver,
			ParseMode: s.config.ParseMode,
		}
		bb, _ = json.Marshal(body)
	}

	url := fmt.Sprintf("https://api.telegram.org/bot%s/sendMessage", s.config.Token)
	res, err := c.Post(url, "application/json", strings.NewReader(string(bb)))
	if err != nil {
		return err
	}

	reb, err := io.ReadAll(res.Body)
	var resJson model.TelegramResponse
	json.Unmarshal(reb, &resJson)
	if alert.Status == "firing" {
		s.fingers[alert.Fingerprint] = resJson.Result.MessageID
	}

	return nil
}
