package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"time"
)

type GetUpdatesT struct {
	Ok     bool                `json:"ok"`
	Result []GetUpdatesResultT `json:"result"`
}

type GetUpdatesResultT struct {
	UpdateID int                `json:"update_id"`
	Message  GetUpdatesMessageT `json:"message,omitempty"`
}

type GetUpdatesMessageT struct {
	MessageID int `json:"message_id"`
	From      struct {
		ID        int    `json:"id"`
		IsBot     bool   `json:"is_bot"`
		FirstName string `json:"first_name"`
		LastName  string `json:"last_name"`
		Username  string `json:"username"`
	} `json:"from"`
	Chat struct {
		ID        int    `json:"id"`
		FirstName string `json:"first_name"`
		LastName  string `json:"last_name"`
		Username  string `json:"username"`
		Type      string `json:"type"`
	} `json:"chat"`
	Date int    `json:"date"`
	Text string `json:"text"`
}

const telegramToken = "INSERT_YOUR_BOT_TOKEN"
const mainChatID = "INSERT_YOUR_CHAT_ID"
const telegramBaseUrl = "https://api.telegram.org/bot"

const methodGetMe = "getMe"
const methodGetUpdates = "getUpdates"
const methodSendMessage = "sendMessage"

func main() {
	var length, lastUpdateID int
	getUpdates := GetUpdatesT{}
	sendMessageUrl := getUrlByMethod(methodSendMessage)

	for {
		body := getBodyByUrl(getUrlByMethod(methodGetUpdates) + "?offset=" + strconv.Itoa(lastUpdateID))
		err := json.Unmarshal(body, &getUpdates)
		if err != nil {
			fmt.Printf("Error in unmarshal getUpdates: %s", err.Error())
			return
		}

		length = len(getUpdates.Result) - 1
		if length >= 0 {
			if lastUpdateID <= getUpdates.Result[length].UpdateID {
				for _, elem := range getUpdates.Result {
					getBodyByUrl(fmt.Sprintf("%s?chat_id=%s&text=%s", sendMessageUrl, mainChatID, elem.Message.Text))
				}
				lastUpdateID = getUpdates.Result[length].UpdateID + 1
			}
		}
		time.Sleep(time.Second * 30)
	}
}

// create URL with selected method
func getUrlByMethod(methodName string) string {
	return telegramBaseUrl + telegramToken + "/" + methodName
}

// send GET request with created URL
func getBodyByUrl(url string) []byte {
	response, err := http.Get(url)
	if err != nil {
		fmt.Println(err)
	}
	defer response.Body.Close()

	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		fmt.Println(err)
	}

	return body
}
