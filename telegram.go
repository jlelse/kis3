package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"path"
	"strconv"
)

type telegram struct {
	botToken string
}

type telegramUpdate struct {
	Message struct {
		Chat struct {
			Id int `json:"id"`
		} `json:"chat"`
		Id   int    `json:"message_id"`
		Text string `json:"text"`
	} `json:"message"`
}

func initTelegram() {
	if appConfig.TGBotToken == "" {
		fmt.Println("Telegram not configured.")
		return
	}
	tg := &telegram{
		botToken: appConfig.TGBotToken,
	}
	username, err := tg.getBotUsername()
	if err != nil {
		fmt.Println("Failed to setup Telegram:", err)
		return
	}
	err = tg.setTelegramHook()
	if err != nil {
		fmt.Println("Failed to setup Telegram webhook:", err)
		return
	}
	fmt.Println("Authorized Telegram bot on account", username)
	app.telegram = tg
}

func initTelegramRouter() {
	app.router.HandleFunc(path.Join("/telegramHook", appConfig.TGHookSecret), TelegramHookHandler)
}

func TelegramHookHandler(w http.ResponseWriter, r *http.Request) {
	tgUpdate := &telegramUpdate{}
	err := json.NewDecoder(r.Body).Decode(tgUpdate)
	if err != nil {
		http.Error(w, "Failed to decode body", http.StatusBadRequest)
		return
	}
	go respondToTelegramUpdate(tgUpdate)
	return
}

var telegramBaseUrl = "https://api.telegram.org/bot"

func (t *telegram) getBotUsername() (string, error) {
	tgUrl, err := url.Parse(telegramBaseUrl + t.botToken + "/getMe")
	if err != nil {
		return "", errors.New("failed to create Telegram request")
	}
	req, _ := http.NewRequest(http.MethodPost, tgUrl.String(), nil)
	resp, err := http.DefaultClient.Do(req)
	if err != nil || resp.StatusCode != 200 {
		return "", errors.New("failed to get Telegram bot info")
	}
	tgBotInfo := &struct {
		Ok     bool `json:"ok"`
		Result struct {
			Id       int    `json:"id"`
			Username string `json:"username"`
		} `json:"result"`
	}{}
	err = json.NewDecoder(resp.Body).Decode(tgBotInfo)
	_ = resp.Body.Close()
	if err != nil || !tgBotInfo.Ok {
		return "", errors.New("failed to parse Telegram bot info")
	}
	// If getMe returns no username, but only an ID for whatever reason
	if len(tgBotInfo.Result.Username) == 0 {
		tgBotInfo.Result.Username = strconv.Itoa(tgBotInfo.Result.Id)
	}
	return tgBotInfo.Result.Username, nil
}

func (t *telegram) setTelegramHook() error {
	if len(appConfig.BaseUrl) < 1 {
		return errors.New("base URL not configured")
	}
	hookUrl, e := url.Parse(appConfig.BaseUrl)
	if e != nil {
		return errors.New("failed to parse base URL")
	}
	hookUrl.Path = path.Join(hookUrl.Path, path.Join("telegramHook", appConfig.TGHookSecret))
	params := url.Values{}
	params.Add("url", hookUrl.String())
	params.Add("allowed_updates", "[\"message\"]")
	tgUrl, err := url.Parse(telegramBaseUrl + t.botToken + "/setWebhook")
	if err != nil {
		return errors.New("failed to create Telegram request")
	}
	tgUrl.RawQuery = params.Encode()
	req, _ := http.NewRequest(http.MethodPost, tgUrl.String(), nil)
	resp, err := http.DefaultClient.Do(req)
	if err != nil || resp.StatusCode != 200 {
		return errors.New("failed to set Telegram webhook")
	}
	fmt.Println("Telegram webhook URL:", hookUrl.String())
	return nil
}

func (t *telegram) sendMessage(chat int, message string) error {
	return t.replyToMessage(chat, message, 0)
}

func (t *telegram) replyToMessage(chat int, message string, replyTo int) error {
	params := url.Values{}
	params.Add("chat_id", strconv.Itoa(chat))
	if replyTo != 0 {
		params.Add("reply_to_message_id", strconv.Itoa(replyTo))
	}
	params.Add("text", message)
	tgUrl, err := url.Parse(telegramBaseUrl + t.botToken + "/sendMessage")
	if err != nil {
		return errors.New("failed to create Telegram request")
	}
	tgUrl.RawQuery = params.Encode()
	req, _ := http.NewRequest(http.MethodPost, tgUrl.String(), nil)
	resp, err := http.DefaultClient.Do(req)
	if err != nil || resp.StatusCode != 200 {
		return errors.New("failed to send Telegram message")
	}
	return nil
}
