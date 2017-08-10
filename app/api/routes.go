package api

import (
	"github.com/jloup/abotllinaire/app/db"
	"github.com/labstack/echo"
)

type Route struct {
	Method      string
	Middlewares []echo.MiddlewareFunc
	Handler     echo.HandlerFunc
	Path        string
}

type Group struct {
	Root        string
	Middlewares []echo.MiddlewareFunc
	Routes      []Route
}

var Groups = []Group{
	{
		"",
		nil,
		[]Route{
			{echo.GET, nil, Get_MakeVerses, "/verses/create"},
			{echo.GET, nil, Get_FacebookVerses, "/fb/verses"},
		},
	},
	{
		"/fb",
		nil,
		[]Route{
			{echo.GET, nil, Get_FBHook, "/hook"},
			{echo.POST, []echo.MiddlewareFunc{FacebookSignatureAuth, UnmarshalFacebookMessage}, Post_FBHook, "/hook"},
		},
	},
}

func Get_FBHook(c echo.Context) error {
	query := &ApiFBMessengerChallenge{}

	query.SetVerifyTokenParam(c.Request().URL.Query().Get("hub.verify_token"))
	query.SetChallengeParam(c.Request().URL.Query().Get("hub.challenge"))

	return RunApiQueryResponseRaw(c, query)
}

func Post_FBHook(c echo.Context) error {

	in := c.Get("fbmsg").(FacebookMessengerMsg)

	for _, entry := range in.Entry {
		for _, msg := range entry.Messaging {

			exists, err := db.FacebookSeqExists(msg.Sender.Id, msg.Message.Seq)
			if err != nil {
				log.Errorf("db error %v", err)
				continue
			}

			if exists {
				log.WithField("seq", db.BuildSeqKey(msg.Sender.Id, msg.Message.Seq)).Infof("already responded to this message")
				continue
			}

			db.SetFacebookSeq(msg.Sender.Id, msg.Message.Seq, db.SeqReceived)

			query := &ApiFBMessengerMessage{}

			query.SetSenderIdParam(msg.Sender.Id)
			query.SetSeedParam(msg.Message.Text)
			query.Run()

			db.SetFacebookSeq(msg.Sender.Id, msg.Message.Seq, db.SeqResponded)

			if query.Err() != nil {
				log.Error(query.Err())
			}
		}
	}
	return WriteResponseRaw(200, []byte("OK"), "text/plain", c)
}

// /fb/verses?fromId=[verseId]&count=[versesCount]
func Get_FacebookVerses(c echo.Context) error {
	query := &ApiGetFacebookVerses{}

	if c.Request().URL.Query().Get("count") != "" {
		query.SetCountParam(c.Request().URL.Query().Get("count"))
	} else {
		query.SetCountParam("10")
	}

	if c.Request().URL.Query().Get("fromId") != "" {
		query.SetFromIdParam(c.Request().URL.Query().Get("fromId"))
	}

	return RunApiQuery(c, query)
}

// /verses/create?temp=[temperature]&seed=[primetext]&l=[length]
func Get_MakeVerses(c echo.Context) error {
	query := &ApiMakeVerses{}

	query.SetTemperatureParam(c.Request().URL.Query().Get("temp"))
	//query.SetLengthParam(c.Request().URL.Query().Get("l"))
	query.SetLengthParam("500")
	query.SetSeedParam(c.Request().URL.Query().Get("seed"))

	return RunApiQuery(c, query)
}
