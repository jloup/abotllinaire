package api

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"strings"

	"github.com/jloup/abotllinaire/app/db"
)

type ApiFBMessengerMessage struct {
	In struct {
		Seed     string
		SenderId string
	}

	Out string

	Error
}

func (a *ApiFBMessengerMessage) GetRawOut() interface{} {
	if a.err != nil {
		return nil
	}

	return a.Out
}

func (a *ApiFBMessengerMessage) GetJSONOut() json.RawMessage {
	if a.err != nil {
		return nil
	}

	var out json.RawMessage
	var err error
	out, err = json.Marshal(a.Out)

	if err != nil {
		a.Error = Error{INTERNAL_ERROR, err}
	}

	return out
}

func (a *ApiFBMessengerMessage) SetSeedParam(param string) {
	if a.err != nil {
		return
	}

	a.In.Seed = strings.TrimSpace(param)
}

func (a *ApiFBMessengerMessage) SetSenderIdParam(param string) {
	if a.err != nil {
		return
	}

	if param == "" {
		a.Error = Error{INTERNAL_ERROR, fmt.Errorf("sender id is empty")}
		return
	}

	a.In.SenderId = param
}

func (a *ApiFBMessengerMessage) Run() {
	if a.err != nil {
		return
	}

	a.Out = "ok"

	if a.In.Seed == "" {
		SendFacebookMessage(a.In.SenderId, _USER_ERROR_MESSAGE)
		a.Error = Error{INTERNAL_ERROR, fmt.Errorf("empty seed")}
		return
	}

	temperature := rand.Float32()*(_BOT_PARAMETERS.TemperatureMax-_BOT_PARAMETERS.TemperatureMin) + _BOT_PARAMETERS.TemperatureMin

	verses, err := workerPool.Request(float64(temperature), uint16(len(a.In.Seed))+_BOT_PARAMETERS.PoemLen, a.In.Seed)
	if err != nil {
		SendFacebookMessage(a.In.SenderId, _USER_ERROR_MESSAGE)
		a.Error = Error{INTERNAL_ERROR, err}
		return
	}

	SendFacebookMessage(a.In.SenderId, strings.Join(verses, "\n"))

	poem := db.Poem{Seed: a.In.Seed, Content: strings.Join(verses, "\n"), FacebookUserId: a.In.SenderId, Temperature: temperature}

	err = db.InsertPoem(&poem)
	if err != nil {
		a.Error = Error{INTERNAL_ERROR, err}
	}

}

type ApiFBMessengerChallenge struct {
	In struct {
		VerifyToken string
		Challenge   string
	}
	Out string
	Error
}

func (a *ApiFBMessengerChallenge) GetRawOut() interface{} {
	if a.err != nil {
		return nil
	}

	return a.Out
}

func (a *ApiFBMessengerChallenge) GetJSONOut() json.RawMessage {
	if a.err != nil {
		return nil
	}

	var out json.RawMessage
	var err error
	out, err = json.Marshal(a.Out)

	if err != nil {
		a.Error = Error{INTERNAL_ERROR, err}
	}

	return out
}

func (a *ApiFBMessengerChallenge) SetVerifyTokenParam(param string) {
	if a.err != nil {
		return
	}

	if param == "" {
		a.Error = Error{INTERNAL_ERROR, fmt.Errorf("empty param")}
		return
	}

	a.In.VerifyToken = param
}

func (a *ApiFBMessengerChallenge) SetChallengeParam(param string) {
	if a.err != nil {
		return
	}

	if param == "" {
		a.Error = Error{INTERNAL_ERROR, fmt.Errorf("empty param")}
		return
	}

	a.In.Challenge = param
}

func (a *ApiFBMessengerChallenge) Run() {
	if a.err != nil {
		return
	}

	if a.In.VerifyToken != _WEBHOOK_TOKEN {
		a.Out = "wrong validation token"
	} else {
		a.Out = a.In.Challenge
	}

}
