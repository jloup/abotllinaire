package api

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha1"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/labstack/echo"
)

var _WEBHOOK_TOKEN string
var _APP_SECRET string
var _FB_SEND_ENDPOINT string
var _PAGE_TOKEN string

type FacebookConfig struct {
	WebhookToken string
	AppSecret    string
	PageToken    string
	SendEndpoint string
}

const _USER_ERROR_MESSAGE = "Hum, ce que vous me dites là n'est pas très inspirant..."

type BotConfig struct {
	TemperatureMin float32
	TemperatureMax float32
	PoemLen        uint16
}

var _BOT_PARAMETERS BotConfig

func SetFacebookCredentials(config FacebookConfig) {
	_WEBHOOK_TOKEN = config.WebhookToken
	_APP_SECRET = config.AppSecret
	_PAGE_TOKEN = config.PageToken
	_FB_SEND_ENDPOINT = config.SendEndpoint
}

func SetBotParameters(config BotConfig) {
	_BOT_PARAMETERS = config
}

// each facebook request is signed
func verifyRequetsSignature(body []byte, signature string) bool {
	mac := hmac.New(sha1.New, []byte(_APP_SECRET))
	mac.Write(body)
	if fmt.Sprintf("%x", mac.Sum(nil)) != signature {
		return false
	}
	return true
}

type FacebookMessengerMsg struct {
	Object string
	Entry  []struct {
		Id        string
		Time      int
		Messaging []struct {
			Sender struct {
				Id string
			}
			Recipient struct {
				Id string
			}
			Timestamp int
			Message   struct {
				Mid  string
				Seq  int
				Text string
			}
		}
	}
}

// middleware to unmarshal facebook request to FacebookMessengerMsg
func UnmarshalFacebookMessage(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {

		defer c.Request().Body.Close()

		var msg FacebookMessengerMsg

		dec := json.NewDecoder(c.Request().Body)
		err := dec.Decode(&msg)

		if err != nil {
			log.Error(err)
			return echo.ErrUnauthorized
		}

		c.Set("fbmsg", msg)

		return next(c)
	}
}

type FacebookMessengerSend struct {
	Recipient struct {
		Id string `json:"id"`
	} `json:"recipient"`
	Message struct {
		Text string `json:"text"`
	} `json:"message"`
}

func SendFacebookMessage(recipient string, message string) error {
	client := &http.Client{
		Timeout: time.Second * 60,
	}

	msg := FacebookMessengerSend{}
	msg.Recipient.Id = recipient
	msg.Message.Text = message

	b, err := json.Marshal(msg)
	if err != nil {
		return err
	}

	r := bytes.NewReader(b)

	response, err := client.Post(fmt.Sprintf("%s?access_token=%s", _FB_SEND_ENDPOINT, _PAGE_TOKEN), "application/json", r)
	if err != nil || response.StatusCode != 200 {
		bb, _ := ioutil.ReadAll(response.Body)
		log.Error("error", string(bb))
		return err
	}

	return nil
}

// middle ware to check the request is actually signed from facebook
func FacebookSignatureAuth(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {

		defer c.Request().Body.Close()

		limited := io.LimitReader(c.Request().Body, 1000000)
		body, err := ioutil.ReadAll(limited)

		log.Info(string(body))

		if err != nil {
			return echo.ErrUnauthorized
		}

		c.Request().Body = ioutil.NopCloser(bytes.NewReader(body))

		if c.Request().Header.Get("X-Hub-Signature") == "" || !verifyRequetsSignature(body, c.Request().Header.Get("X-Hub-Signature")[5:]) {
			return echo.ErrUnauthorized
		}
		return next(c)
	}
}
