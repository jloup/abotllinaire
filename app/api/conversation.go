package api

import (
	"fmt"
	"math/rand"
	"strings"

	"github.com/jloup/abotllinaire/app/conversation"
	"github.com/jloup/abotllinaire/app/db"
	_ "github.com/labstack/gommon/log"
)

var Intents [][2]string = [][2]string{
	{"subject", "TU? 1G(Parler) DE SUBJECT"},
	{"subject", "Que 1G(penser) TU DE SUBJECT"},
	{"subject", "Qui (?:est|sont) LE? SUBJECT"},
	{"subject", "Qu'est ce que LE? SUBJECT"},
	{"subject", "Que TU 1G(inspirer) LE? SUBJECT"},
	{"subject", "ce que TU 1G(penser) DE SUBJECT"},
	{"subject", "et DE SUBJECT"},
	{"subject", "et LE? SUBJECT"},
	{"freestyle", "(?:1G(composer)|1G(raconter)) un poème"},
	{"freestyle", "2G(Ecrire) un poème"},
	{"freestyle", "2G(Dire) des mots"},
	{"followup", "^(?:encore|again)"},
	{"followup", "j'en veux plus"},
	{"greeting", "^(?P<greetings>Bonjour|Salut|Hi|Coucou|Hello|Hola|Ola)"},
}

var IntentCollection []conversation.Intent

func init() {
	IntentCollection = conversation.NewIntentCollection(Intents)
}

func DispatchUserMessage(msg string, FBUserId string) (string, error) {
	intent, meta := conversation.FindIntent(IntentCollection, msg)

	lastAction, err := db.GetLastUserAction(FBUserId)
	if err != nil {
		return "", fmt.Errorf("cannot retrieve user last action %v", err)
	}

	log.Infof("user last action %v\n", lastAction)

	var response string
	var action db.UserAction

	if intent == "followup" {
		switch lastAction.Type {
		case db.SubjectAction:
			intent = "subject"
			meta["subject"] = lastAction.Meta
		case db.FreestyleAction:
			intent = "freestyle"
		default:
			intent = "unknown"
		}
	}

	switch intent {
	case "subject":
		response, err = SubjectResponse(msg, meta["subject"], meta["pronoun"], FBUserId)
		action.Type = db.SubjectAction
		action.Meta = meta["subject"]
	case "freestyle":
		response, err = FreestyleResponse(msg, FBUserId)
		action.Type = db.FreestyleAction
	case "greeting":
		response, err = GreetingsResponse(msg, meta["greetings"], FBUserId)
		action.Type = db.GreetingsAction
	default:
		response, err = NotRecognizedResponse(msg, FBUserId)
		action.Type = db.NoneAction
	}

	if err != nil {
		return response, err
	}

	err = db.SetLastUserAction(FBUserId, action)

	return response, err
}

func SubjectResponse(msg, subject, pronoun, FBUserId string) (string, error) {
	return SearchVerse(subject, 2, _BOT_PARAMETERS.VerseFilePath, _BOT_PARAMETERS.VerseFilePathLower)
}

func FreestyleResponse(msg, FBUserId string) (string, error) {
	temperature := rand.Float32()*(_BOT_PARAMETERS.TemperatureMax-_BOT_PARAMETERS.TemperatureMin) + _BOT_PARAMETERS.TemperatureMin

	poemLen := _BOT_PARAMETERS.MaxPoemLen

	verses, err := workerPool.Request(float64(temperature), poemLen, msg)
	if err != nil {
		return "", err
	}

	content := strings.Join(verses, "\n")

	poem := db.Poem{Seed: "none", Content: content, FacebookUserId: FBUserId, Temperature: temperature}

	err = db.InsertPoem(&poem)

	return content, err
}

func GreetingsResponse(msg, greetings, FBUserId string) (string, error) {
	if greetings == "" {
		return "Bonjour.", nil
	}

	return fmt.Sprintf("%s%s.", strings.ToUpper(string(greetings[0])), greetings[1:]), nil
}

func NotRecognizedResponse(msg string, FBUserID string) (string, error) {
	return "Hum, ce que vous me dites là n'est pas très inspirant...", nil
}
