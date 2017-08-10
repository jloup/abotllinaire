package api

import (
	"encoding/json"
	"fmt"
	"strings"
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

	response, err := DispatchUserMessage(a.In.Seed, a.In.SenderId)
	if err != nil {
		a.Error = Error{INTERNAL_ERROR, err}
	}

	if response != "" {
		SendFacebookMessage(a.In.SenderId, response)
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
