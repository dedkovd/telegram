package telegram

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

const botApiUrl string = "https://api.telegram.org"

type ResponceResult interface{}

type Responce struct {
	Ok          bool
	Description string
	Result      ResponceResult
}

type UserResponce struct {
	Responce
	Result *User
}

type ChatMemberResponce struct {
	Responce
	Result *ChatMember
}

type MessageResponce struct {
	Responce
	Result *TMessage
}

type UpdateResponse struct {
	Responce
	Result []*Update
}

type User struct {
	Id         int
	First_Name string
	Last_Name  string
	Username   string
}

type ChatMember struct {
	ChatUser *User  `json:"user"`
	Status   string `json:"status"`
}

type TChat struct {
	Id         int
	Type       string
	Title      string
	Username   string
	First_Name string
	Last_Name  string
}

type TMessage struct {
	Message_id int
	From       *User
	Date       int
	Chat       *TChat
	Text       string
}

type Update struct {
	Update_id int
	Message   *TMessage
}

type MessageParams struct {
	Chat_id int    `json:"chat_id"`
	Text    string `json:"text"`
}

type ChatMemberParams struct {
	Chat_id int `json:"chat_id"`
	User_id int `json:"user_id"`
}

type GetUpdatesParams struct {
	Offset  int `json:"offset"`
	Limit   int `json:"limit"`
	Timeout int `json:"timeout"`
}

func responceFromJson(pJson []byte, res interface{}) {
	err := json.Unmarshal(pJson, res)

	if err != nil {
		log.Fatal(err)
	}
}

func dataToJson(data interface{}) []byte {
	res, err := json.Marshal(data)

	if err != nil {
		log.Fatal(err)
	}

	return res
}

type BotApi struct {
	BotToken     string
	AllowedChats []int
	Roulette     []int
}

func (ba *BotApi) runMethod(method string, params interface{}) []byte {
	url := fmt.Sprintf("%s/bot%s/%s", botApiUrl, ba.BotToken, method)
	data := dataToJson(params)

	resp, err := http.Post(url, "application/json", bytes.NewBuffer(data))

	if err != nil {
		log.Fatal(err)
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}

	resp.Body.Close()

	return body
}

func (ba *BotApi) CheckMessage(m *TMessage) bool {
	if m == nil {
		log.Println("Nil message")
		return false
	}
	for _, v := range ba.AllowedChats {
		if v == m.Chat.Id {
			return true
		}
	}

	log.Printf("Not allowed chat: %d\n", m.Chat.Id)
	return false
}

func (ba *BotApi) GetMe() User {
	r := ba.runMethod("getMe", nil)
	var res UserResponce
	responceFromJson(r, &res)

	if !res.Ok {
		log.Println("Not OK\n", res.Description)
	}

	return *res.Result
}

func (ba *BotApi) SendMessage(chat_id int, text string) TMessage {
	m := MessageParams{Chat_id: chat_id, Text: text}

	r := ba.runMethod("sendMessage", m)
	var mr MessageResponce
	responceFromJson(r, &mr)

	if !mr.Ok {
		log.Println("Not OK\n", mr.Description)
	}

	return *mr.Result
}

func (ba *BotApi) KickChatMember(chat_id int, user_id int) bool {
	r := ba.runMethod("kickChatMember", ChatMemberParams{chat_id, user_id})
	var resp Responce
	responceFromJson(r, &resp)

	if !resp.Ok {
		log.Println("Not OK\n", resp.Description)
		ba.SendMessage(chat_id, resp.Description)
	}

	return resp.Ok
}

func (ba *BotApi) GetChatMember(chat_id int, user_id int) ChatMember {
	r := ba.runMethod("getChatMember", ChatMemberParams{chat_id, user_id})
	var resp ChatMemberResponce
	responceFromJson(r, &resp)

	if !resp.Ok {
		log.Println("Not OK\n", resp.Description)
	}

	return *resp.Result
}

func (ba *BotApi) GetUpdates(offset int) []*Update {
	p := GetUpdatesParams{Offset: offset, Limit: 100}

	r := ba.runMethod("getUpdates", p)

	var ur UpdateResponse
	responceFromJson(r, &ur)

	if !ur.Ok {
		log.Println("Not OK", ur.Description)
	}

	return ur.Result
}
