package tg

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/rs/zerolog/log"
)

type TG struct {
	enable bool
	token string
	chat_id string
}
type Body struct {
	ChatID   string `json:"chat_id"`
	Text     string `json:"text"`
}

var tg TG
var tgQueue chan func()

func SetTg(enable bool, token string, chat_id string) {
	if !enable {
		return
	}

	// Set TG
	// It is singleton
	tgQueue = make(chan func())
	tg = TG{
		enable,
		token,
		chat_id,
	}

	// For thread safe
	go func() {
		for tg := range tgQueue {
			tg()
		}
	}()
}

func enqueue(tg func()) {
	tgQueue <- tg
}

func SendMsg(msg string) {
	if !tg.enable {
		return
	}

	url := fmt.Sprintf("https://api.telegram.org/bot%s/sendMessage?parse_mode=MarkdownV2", tg.token)

	body := Body{
		tg.chat_id,
		msg,
	}
	bodyBytes, err := json.Marshal(body)
	if err != nil {
		log.Error().Stack().Err(err).Msg("")
		return
	}
	buff := bytes.NewBuffer(bodyBytes)

	req, err := http.NewRequest(
		"POST",
		url,
		buff,
	)
	if err != nil {
		log.Error().Stack().Err(err).Msg("")
		return
	}
	req.Header.Set("Content-type", "application/json")

	client := &http.Client{}

	tg := func() {
		resp, err := client.Do(req)
		if err != nil {
			log.Error().Stack().Err(err).Msg("")
		}
		defer resp.Body.Close()

		fmt.Println(resp.Status)
		if !(200 <= resp.StatusCode && resp.StatusCode <= 299) {
			log.Error().Msg("Fail to seding msg from tg module")
		}
	}
	enqueue(tg)
}
