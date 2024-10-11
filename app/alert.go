package app

import (
	"bharvest.io/axelmon/log"
	"bharvest.io/axelmon/server"
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"net/http"
	"strings"
	"time"
)

type alertMsg struct {
	tg  bool
	slk bool

	resolved bool
	message  string
	uniqueId string

	tgChannel  string
	tgKey      string
	tgMentions string

	slkHook     string
	slkMentions string

	resendDuration *time.Duration
}

type notifyDest uint8

const (
	tg notifyDest = iota
	slk
)

func shouldNotify(msg *alertMsg, dest notifyDest) bool {
	server.GlobalState.Alerts.NotifyMux.Lock()
	defer server.GlobalState.Alerts.NotifyMux.Unlock()
	var whichMap map[string]time.Time
	var service string

	if server.GlobalState.Alerts.AllAlarms == nil {
		server.GlobalState.Alerts.AllAlarms = make(map[string]time.Time)
	}

	switch dest {
	case tg:
		whichMap = server.GlobalState.Alerts.SentTgAlarms
		service = "Telegram"
	case slk:
		whichMap = server.GlobalState.Alerts.SentSlkAlarms
		service = "Slack"
	default:
		panic("unhandled default case")
	}

	switch {
	case !whichMap[msg.message].IsZero() && !msg.resolved:
		// already sent this alert
		if msg.resendDuration != nil && whichMap[msg.message].Add(*msg.resendDuration).Before(time.Now()) {
			log.Warn("it's not resolved. resending...")
			break
		}
		return false
	case !whichMap[msg.message].IsZero() && msg.resolved:
		// alarm is cleared
		delete(whichMap, msg.message)
		log.Info(fmt.Sprintf("🟢 Resolved     alarm (%s) - notifying %s", msg.message, service))
		return true
	case msg.resolved:
		// duplicated resolved messages
		return false
	}

	//log.Info(fmt.Sprintf("new alarm  (%s) - notifying %s", msg.message, service))
	whichMap[msg.message] = time.Now()
	return true
}

func notifySlack(msg *alertMsg) (err error) {
	if !msg.slk {
		return
	}
	if !shouldNotify(msg, slk) {
		return nil
	}
	data, err := json.Marshal(buildSlackMessage(msg))
	if err != nil {
		return
	}

	req, err := http.NewRequest("POST", msg.slkHook, bytes.NewBuffer(data))
	if err != nil {
		return
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return
	}
	_ = resp.Body.Close()

	if resp.StatusCode != 200 {
		return fmt.Errorf("could not notify slack. got %d response", resp.StatusCode)
	}

	return
}

type SlackMessage struct {
	Text        string       `json:"text"`
	Attachments []Attachment `json:"attachments"`
}

type Attachment struct {
	Text      string `json:"text"`
	Color     string `json:"color"`
	Title     string `json:"title"`
	TitleLink string `json:"title_link"`
}

func buildSlackMessage(msg *alertMsg) *SlackMessage {
	color := "danger"
	prefix := "🛑 "

	if msg.resolved {
		msg.message = "OK: " + msg.message
		prefix = "🟢 Healthy: "
		color = "good"
	}
	return &SlackMessage{
		Text: msg.message,
		Attachments: []Attachment{
			{
				Title: fmt.Sprintf("🤖 Axelmon %s %s", prefix, msg.slkMentions),
				Color: color,
			},
		},
	}
}

func notifyTg(msg *alertMsg) (err error) {
	if !msg.tg {
		return nil
	}
	if !shouldNotify(msg, tg) {
		return nil
	}
	prefix := "🛑 "
	if msg.resolved {
		prefix = "🟢 Healthy: "
	}

	bot, err := tgbotapi.NewBotAPI(msg.tgKey)
	if err != nil {
		log.Error(errors.New(fmt.Sprintf("notify telegram: %v", err)))
		return
	}

	mc := tgbotapi.NewMessageToChannel(msg.tgChannel, fmt.Sprintf("%s - %s %s\n%s", "🤖 Axelmon", prefix, msg.message, msg.tgMentions))
	_, err = bot.Send(mc)
	if err != nil {
		log.Error(errors.New(fmt.Sprintf("telegram send: %v", err)))
	}
	return err
}

func (c *Config) alert(message string, resolved, notSend bool) {
	prefix := "🛑 "
	if resolved {
		prefix = "🟢 Healthy "
	}

	log.Info(fmt.Sprintf("%s %s", prefix, message))

	if !notSend {
		c.alertMux.RLock()
		mentions := ""
		for _, m := range c.Alerts.Slack.Mentions {
			if m[:1] != "@" {
				mentions = fmt.Sprintf("%s <@%s>", mentions, m)
			}
		}

		resendDuration := time.Hour * 24
		if c.Alerts.ResendDuration != nil {
			resendDuration = time.Duration(*c.Alerts.ResendDuration)
		}
		a := &alertMsg{
			tg:             c.Alerts.Tg.Enabled,
			slk:            c.Alerts.Slack.Enabled,
			resolved:       resolved,
			message:        message,
			tgChannel:      c.Alerts.Tg.ChatID,
			tgKey:          c.Alerts.Tg.Token,
			tgMentions:     strings.Join(c.Alerts.Tg.Mentions, " "),
			slkHook:        c.Alerts.Slack.Webhook,
			slkMentions:    mentions,
			resendDuration: &resendDuration,
		}
		c.alertChan <- a
		c.alertMux.RUnlock()
	}

	server.GlobalState.Alerts.NotifyMux.Lock()
	defer server.GlobalState.Alerts.NotifyMux.Unlock()
	if server.GlobalState.Alerts.AllAlarms == nil {
		server.GlobalState.Alerts.AllAlarms = make(map[string]time.Time)
	}
	if server.GlobalState.Alerts.SentTgAlarms == nil {
		server.GlobalState.Alerts.SentTgAlarms = make(map[string]time.Time)
	}
	if server.GlobalState.Alerts.SentSlkAlarms == nil {
		server.GlobalState.Alerts.SentSlkAlarms = make(map[string]time.Time)
	}

	if resolved && !server.GlobalState.Alerts.AllAlarms[message].IsZero() {
		delete(server.GlobalState.Alerts.AllAlarms, message)
		return
	} else if resolved {
		return
	}
	server.GlobalState.Alerts.AllAlarms[message] = time.Now()

}
